package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-lib/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups/status,verbs=get;update;patch

// Reconcile implement the reconciliation loop for backups
func (r *BackupReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("backup", req.NamespacedName)
	log.Info("Reconciling Backup")
	// your logic here
	backup := &mysqlv1alpha1.Backup{}
	err := r.Client.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if backup.Status.LastCondition != "" {
		return reconcile.Result{}, nil
	}
	// Check if this Pod already exists
	pod := &corev1.Pod{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: backup.Spec.Instance + "-0", Namespace: backup.Namespace}, pod)
	if err != nil {
		t := metav1.Now()
		condition := status.Condition{
			Type:               status.ConditionType("podmonitor"),
			Status:             corev1.ConditionTrue,
			Reason:             status.ConditionReason("Failed"),
			Message:            fmt.Sprintf("Cannot find pod %s-0; error: %v", backup.Spec.Instance, err),
			LastTransitionTime: t,
		}
		backup.Status.LastCondition = "Failed"
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	store := &mysqlv1alpha1.Store{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: backup.Spec.Store, Namespace: backup.Namespace}, store)
	if err != nil {
		t := metav1.Now()
		condition := status.Condition{
			Type:               status.ConditionType("store"),
			Status:             corev1.ConditionTrue,
			Reason:             status.ConditionReason("Failed"),
			Message:            fmt.Sprintf("Error accessing store %s: %v", backup.Spec.Store, err),
			LastTransitionTime: t,
		}
		backup.Status.LastCondition = "Failed"
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	cfg := agent.NewConfiguration()
	cfg.BasePath = "http://" + pod.Status.PodIP + ":8080"
	api := agent.NewAPIClient(cfg)
	payload := agent.Backup{
		S3access: agent.S3Info{
			Bucket: store.Spec.S3.Bucket,
			Path:   store.Spec.S3.Path,
			AwsConfig: agent.AwsConfig{
				AwsAccessKeyId:     store.Spec.S3Backup.AWSConfig.AccessKey,
				AwsSecretAccessKey: store.Spec.S3Backup.AWSConfig.SecretKey,
				Region:             store.Spec.S3Backup.AWSConfig.Region,
			},
		},
	}
	b, _, err := api.MysqlApi.CreateBackup(context.TODO(), payload, nil)
	if err != nil {
		t := metav1.Now()
		condition := status.Condition{
			Type:               status.ConditionType("backup"),
			Status:             corev1.ConditionTrue,
			Reason:             status.ConditionReason("Failed"),
			Message:            fmt.Sprintf("Error accessing api: %v", err),
			LastTransitionTime: t,
		}
		backup.Status.LastCondition = "Failed"
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	t := metav1.Now()
	details := mysqlv1alpha1.BackupDetails{
		Location:   b.Location,
		BackupTime: &metav1.Time{Time: b.Timestamp},
	}
	condition := status.Condition{
		Type:               status.ConditionType("backup"),
		Status:             corev1.ConditionTrue,
		Reason:             status.ConditionReason(b.Status),
		Message:            b.Message,
		LastTransitionTime: t,
	}

	backup.Status.LastCondition = b.Status
	backup.Status.Conditions = append(backup.Status.Conditions, condition)
	backup.Status.Details = &details
	err = r.Client.Status().Update(context.TODO(), backup)
	if err != nil {
		return reconcile.Result{}, err
	}
	go r.MonitorBackup(req.NamespacedName, api.MysqlApi, b.Timestamp.Format(time.RFC3339))
	return reconcile.Result{}, nil
}

// MonitorBackup watch backup progress and update results
func (r *BackupReconciler) MonitorBackup(n types.NamespacedName, a *agent.MysqlApiService, backupName string) {
	log := r.Log.WithValues("Request.Namespace", n.Namespace, "Request.Name", n.Name)
	endTime := time.Now().Add(60 * time.Second)
	succeeded := false
	backup := &mysqlv1alpha1.Backup{}
	err := r.Client.Get(context.TODO(), n, backup)
	if err != nil {
		log.Info(fmt.Sprintf("Error querying backup: %v", err))
		return
	}
	log.Info(fmt.Sprintf("Starting to check for backup, current status %s", backup.Status.LastCondition))
	for time.Now().Before(endTime) && !succeeded {
		log.Info(fmt.Sprintf("Loop..."))
		b, _, err := a.GetBackupByName(context.TODO(), backupName, nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		lastCondition := backup.Status.LastCondition
		if b.Status != lastCondition {
			t := metav1.Now()
			backup.Status.LastCondition = b.Status
			condition := status.Condition{
				Type:               status.ConditionType("backup"),
				Status:             corev1.ConditionTrue,
				Reason:             status.ConditionReason(b.Status),
				Message:            b.Message,
				LastTransitionTime: t,
			}
			backup.Status.Conditions = append(backup.Status.Conditions, condition)
			err = r.Client.Status().Update(context.TODO(), backup)
			if err != nil {
				backup.Status.LastCondition = lastCondition
			}
		}
		if backup.Status.LastCondition == "Available" || backup.Status.LastCondition == "Failed" {
			succeeded = true
			break
		}
		time.Sleep(2 * time.Second)
	}
	if !succeeded {
		t := metav1.Now()
		backup.Status.LastCondition = "Failed"
		condition := status.Condition{
			Type:               status.ConditionType("backup"),
			Status:             corev1.ConditionTrue,
			Reason:             status.ConditionReason("Failed"),
			Message:            "Backup did not finish in the expected time",
			LastTransitionTime: t,
		}
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(context.TODO(), backup)
		if err != nil {
			log.Info(fmt.Sprintf("Could not update status %v", err))
		}
	}
	log.Info("Monitor backup is now over...")
}

// SetupWithManager configure type of events the manager should watch
func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Backup{}).
		Complete(r)
}
