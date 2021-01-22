package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Properties *StatefulSetProperties
	Crontab    Crontab
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

	if r.Crontab.reScheduleAll(r.Client, instance, r.Log, r.Scheme) {
		if err := r.Status().Update(ctx, instance); err != nil {
			log.Info(fmt.Sprintf("Error rescheduling jobs, err: %v", err))
			return ctrl.Result{}, nil
		}
		log.Info(fmt.Sprintf("Success rescheduling jobs, backup: %d, maintenance: %d", instance.Status.Backup.EntryID, instance.Status.Maintenance.EntryID))
		return ctrl.Result{}, nil
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
		Owns(&mysqlv1alpha1.Backup{}).
		Complete(r)
}
