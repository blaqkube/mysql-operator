package controllers

import (
	"context"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxBackupConditions = 10
)

// BackupManager provides methods to manage the backup subcomponents
type BackupManager struct {
	Context    context.Context
	Reconciler *BackupReconciler
}

func (bm *BackupManager) setBackupCondition(backup *mysqlv1alpha1.Backup, condition metav1.Condition) (ctrl.Result, error) {
	backup.Status.Ready = condition.Status
	backup.Status.Reason = condition.Reason
	backup.Status.Message = condition.Message
	conditions := append(backup.Status.Conditions, condition)
	if len(conditions) > maxBackupConditions {
		conditions = conditions[1:]
	}
	backup.Status.Conditions = conditions
	log := bm.Reconciler.Log.WithValues("namespace", backup.Namespace, "store", backup.Name)
	log.Info("Updating store with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := bm.Reconciler.Status().Update(bm.Context, backup); err != nil {
		log.Error(err, "Unable to update store")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
