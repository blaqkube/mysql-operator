package controllers

import (
	"context"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxOperationConditions = 10
)

// OperationManager provides methods to manage operations
type OperationManager struct {
	Context     context.Context
	Reconciler  *OperationReconciler
	TimeManager *TimeManager
}

func (om *OperationManager) setOperationCondition(operation *mysqlv1alpha1.Operation, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == operation.Status.Reason {
		c := len(operation.Status.Conditions) - 1
		d := om.TimeManager.Next(operation.Status.Conditions[c].LastTransitionTime.Time)
		if condition.Reason != mysqlv1alpha1.DatabaseSucceeded {
			return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
		}
		return ctrl.Result{}, nil
	}
	operation.Status.Ready = condition.Status
	operation.Status.Reason = condition.Reason
	operation.Status.Message = condition.Message
	conditions := append(operation.Status.Conditions, condition)
	if len(conditions) > maxOperationConditions {
		conditions = conditions[1:]
	}
	operation.Status.Conditions = conditions
	log := om.Reconciler.Log.WithValues("namespace", operation.Namespace, "chat", operation.Name)
	log.Info("Updating chat with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := om.Reconciler.Status().Update(om.Context, operation); err != nil {
		log.Error(err, "Unable to update chat")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// NoOp is a No Operation function
func (om *OperationManager) NoOp() error {
	return nil
}
