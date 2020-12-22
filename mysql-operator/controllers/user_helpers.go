package controllers

import (
	"context"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxUserConditions = 10
)

// UserManager provides methods to manage user subcomponents
type UserManager struct {
	Context    context.Context
	Reconciler *UserReconciler
}

func (sm *UserManager) setUserCondition(user *mysqlv1alpha1.User, condition metav1.Condition) (ctrl.Result, error) {
	user.Status.Ready = condition.Status
	user.Status.Reason = condition.Reason
	user.Status.Message = condition.Message
	conditions := append(user.Status.Conditions, condition)
	if len(conditions) > maxUserConditions {
		conditions = conditions[1:]
	}
	user.Status.Conditions = conditions
	log := sm.Reconciler.Log.WithValues("namespace", user.Namespace, "store", user.Name)
	log.Info("Updating store with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := sm.Reconciler.Status().Update(sm.Context, user); err != nil {
		log.Error(err, "Unable to update store")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
