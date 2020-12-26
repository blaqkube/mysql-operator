package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Properties StatefulSetProperties
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="apps",resources=statefulsets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups/finalizers,verbs=update

// Reconcile implement the reconciliation loop for backups
func (r *BackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("backup", req.NamespacedName)
	log.Info("Running a reconcile loop")

	// Fetch the Backup instance
	backup := &mysqlv1alpha1.Backup{}
	if err := r.Get(ctx, req.NamespacedName, backup); err != nil {
		log.Info("Unable to fetch backup from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if backup.Status.Reason == mysqlv1alpha1.BackupSucceeded ||
		backup.Status.Reason == mysqlv1alpha1.BackupFailed ||
		backup.Status.Reason == mysqlv1alpha1.BackupNotImplemented {
		return ctrl.Result{}, nil
	}

	bm := &BackupManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	if backup.Status.Reason == mysqlv1alpha1.BackupRunning {
		b, err := bm.MonitorBackup(backup)
		if err != nil {
			condition := metav1.Condition{
				Type:               "available",
				Status:             metav1.ConditionFalse,
				LastTransitionTime: metav1.Now(),
				Reason:             mysqlv1alpha1.BackupNotImplemented,
				Message:            "Monitoring not implemented for now",
			}
			return bm.setBackupCondition(backup, condition, b)
		}
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.BackupNotImplemented,
			Message:            "Monitoring not implemented for now",
		}
		return bm.setBackupCondition(backup, condition, b)
	}

	b, err := bm.CreateBackup(backup)
	if err != nil {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
		}
		switch err {
		case ErrStoreNotFound:
			condition.Reason = mysqlv1alpha1.BackupStoreAccessError
			condition.Message = "Backup store not found"
		case ErrStoreNotReady:
			condition.Reason = mysqlv1alpha1.BackupStoreNotReady
			condition.Message = "Backup store not ready"
		case ErrInstanceNotFound:
			condition.Reason = mysqlv1alpha1.BackupInstanceAccessError
			condition.Message = "Backup instance not found"
		case ErrInstanceNotReady:
			condition.Reason = mysqlv1alpha1.BackupInstanceNotReady
			condition.Message = "Backup instance not ready"
		case ErrPodNotFound, ErrAgentAccessFailed:
			condition.Reason = mysqlv1alpha1.BackupAgentNotFound
			condition.Message = "Backup agent not found"
		case ErrAgentRequestFailed:
			condition.Reason = mysqlv1alpha1.BackupAgentFailed
			condition.Message = "Backup request failed"
		case ErrMissingVariable:
			condition.Reason = mysqlv1alpha1.BackupMissingVariable
			condition.Message = "Backup environment variable missing from store"
		default:
			condition.Reason = mysqlv1alpha1.BackupAgentFailed
			condition.Message = fmt.Sprintf("Unexpected failure with agent: %v", err)
		}
		return bm.setBackupCondition(backup, condition, nil)
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.BackupRunning,
		Message:            fmt.Sprintf("Backup started on %s with success. Now monitoring progress", backup.Spec.Instance),
	}

	return bm.setBackupCondition(
		backup,
		condition,
		&mysqlv1alpha1.BackupDetails{
			Identifier: b.Identifier,
			Location:   b.Location,
			StartTime:  b.StartTime,
		})
}

// SetupWithManager configure type of events the manager should watch
func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Backup{}).
		Complete(r)
}
