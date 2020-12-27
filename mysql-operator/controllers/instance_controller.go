package controllers

import (
	"context"
	"fmt"
	"sync"

	// "time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/google/uuid"

	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// BackupJob is a struct that manages Jobs for backups
type BackupJob struct {
	client.Client
	Instance types.NamespacedName
	Log      logr.Logger
}

// NewBackupJob creates a BackupJob to schedule it
func NewBackupJob(client client.Client, instance types.NamespacedName, log logr.Logger) *BackupJob {
	return &BackupJob{
		Client:   client,
		Instance: instance,
		Log:      log,
	}
}

// Run implement the Job interface to use with Cron AddFunc()
func (b *BackupJob) Run() {
	ctx := context.Background()
	instance := mysqlv1alpha1.Instance{}
	if err := b.Client.Get(ctx, b.Instance, &instance); err != nil {
		b.Log.Info(fmt.Sprintf("job for %s/%s failed. Could not access instance...", b.Instance.Namespace, b.Instance.Name))
		return
	}
	b.Log.Info(fmt.Sprintf("job for %s/%s succeeded...", b.Instance.Namespace, b.Instance.Name))
}

// Crontab provides a simple struct to manage cron EntryID for instances
type Crontab struct {
	M           sync.Mutex
	Schedulers  map[int]string
	Cron        *cron.Cron
	Incarnation string
}

var (
	crontab = &Crontab{}
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Properties *StatefulSetProperties
}

func (c *Crontab) isBackupScheduleRunning(instance mysqlv1alpha1.Instance) bool {
	entry := instance.Status.BackupSchedule.EntryID
	if entry == -1 {
		return false
	}
	if c.Cron == nil {
		c.Schedulers = map[int]string{}
		c.Cron = cron.New()
		c.Cron.Start()
		c.Incarnation = uuid.New().String()
		return false
	}
	c.M.Lock()
	defer c.M.Unlock()
	if c.Incarnation != instance.Status.BackupSchedule.Incarnation {
		return false
	}
	v, ok := c.Schedulers[entry]
	if !ok {
		return false
	}
	if v != fmt.Sprintf("%s/%s", instance.Namespace, instance.Name) {
		return false
	}
	if instance.Spec.BackupSchedule.Schedule != instance.Spec.BackupSchedule.Schedule {
		c.Cron.Remove(cron.EntryID(entry))
		return false
	}
	return true
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances/finalizers,verbs=update

// Reconcile implement the reconciliation loop for instances
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("instance", req.NamespacedName)
	log.Info("Running a reconcile loop")

	instance := &mysqlv1alpha1.Instance{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Info("Unable to fetch instance from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !crontab.isBackupScheduleRunning(*instance) {
		if instance.Spec.BackupSchedule.Schedule != "" {
			crontab.M.Lock()
			defer crontab.M.Unlock()
			entry, err := crontab.Cron.AddJob(instance.Spec.BackupSchedule.Schedule, NewBackupJob(r.Client, req.NamespacedName, log))
			if err == nil {
				log.Info("Backup scheduler job for instance added")
				crontab.Schedulers[int(entry)] = fmt.Sprintf("%s/%s", req.NamespacedName.Namespace, req.NamespacedName.Name)
				instance.Status.BackupSchedule = mysqlv1alpha1.BackupScheduleStatus{
					Schedule:    instance.Spec.BackupSchedule.Schedule,
					EntryID:     int(entry),
					Incarnation: crontab.Incarnation,
				}
				if err := r.Status().Update(ctx, instance); err != nil {
					log.Info(fmt.Sprintf("Error updating crontab, err: %v", err))
					return ctrl.Result{}, nil
				}
				log.Info("crontab updated with success...")
				return ctrl.Result{}, nil
			}
			if err != nil {
				log.Info("Backup scheduler job for instance failed, error: %v", err)
			}
		}
	}

	im := &InstanceManager{
		Context:     ctx,
		Reconciler:  r,
		Properties:  r.Properties,
		TimeManager: NewTimeManager(),
	}
	secret, err := im.getExporterSecret(instance)
	if err != nil && !errors.IsNotFound(err) {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceExporterSecretInaccessible,
			Message:            fmt.Sprintf("The secret exporter could not be accessed: %v", err),
		}
		return im.setInstanceCondition(instance, condition)
	}
	if err != nil {
		return im.createExporterSecret(instance)
	}
	if secret.UID != instance.Status.ExporterSecret.UID {
		im.deleteExporterSecret(instance, secret)
	}

	sts, stsErr := im.getStatefulSet(instance)
	if stsErr != nil && !errors.IsNotFound(stsErr) {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceStatefulSetInaccessible,
			Message:            fmt.Sprintf("The statefulset could not be accessed: %v", stsErr),
		}
		return im.setInstanceCondition(instance, condition)
	}
	// TODO: detect changes on the store when it is needed and start the instance accordingly
	store := &mysqlv1alpha1.Store{}
	location := ""
	if instance.Spec.Restore.Store != "" {
		location = instance.Spec.Restore.Location
		NamespacedStore := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Spec.Restore.Store}
		if err := r.Get(ctx, NamespacedStore, store); err != nil {
			log.Info(fmt.Sprintf("Unable to fetch store, error: %v", err), "store", store.Name)
			condition := metav1.Condition{
				Type:               "available",
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.Now(),
				Reason:             mysqlv1alpha1.InstanceStoreInaccessible,
				Message:            fmt.Sprintf("Cannot access the restore store: %v", err),
			}
			return im.setInstanceCondition(instance, condition)
		}
		if store.Status.Ready == metav1.ConditionFalse {
			log.Error(err, "Store exists but is not available", "store", store.Name)
			condition := metav1.Condition{
				Type:               "available",
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.Now(),
				Reason:             mysqlv1alpha1.InstanceStoreNotReady,
				Message:            fmt.Sprintf("Store %s ready is False or Unknown, Reason: %s", store.Name, store.Status.Reason),
			}
			return im.setInstanceCondition(instance, condition)
		}
	}
	if stsErr != nil {
		return im.createStatefulSet(instance, store, location)
	}
	// TODO: Check the StatefulSet matches the requirements
	if sts.UID != instance.Status.StatefulSet.UID {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceStatefulSetUpdated,
			Message: fmt.Sprintf("The statefulset has been updated %s from %s",
				sts.UID,
				sts.ResourceVersion,
			),
		}
		return im.setInstanceCondition(instance, condition)
	}
	if sts.Status.ReadyReplicas != sts.Status.Replicas {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.InstanceStatefulSetWaiting,
			Message:            fmt.Sprintf("Waiting for statefulset to become ready"),
		}
		return im.setInstanceCondition(instance, condition)
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.InstanceStatefulSetReady,
		Message:            fmt.Sprintf("The statefulset is ready"),
	}
	return im.setInstanceCondition(instance, condition)
}

// SetupWithManager configure type of events the manager should watch
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Instance{}).
		Owns(&corev1.Secret{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
