package controllers

import (
	"context"
	"fmt"
	"os"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxStoreConditions = 10
)

// StoreManager provides methods to manage the store subcomponents
type StoreManager struct {
	Context     context.Context
	Reconciler  *StoreReconciler
	TimeManager *TimeManager
}

func (sm *StoreManager) setStoreCondition(store *mysqlv1alpha1.Store, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == store.Status.Reason {
		c := len(store.Status.Conditions) - 1
		d := sm.TimeManager.Next(store.Status.Conditions[c].LastTransitionTime.Time)
		return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
	}
	store.Status.Ready = condition.Status
	store.Status.Reason = condition.Reason
	store.Status.Message = condition.Message
	conditions := append(store.Status.Conditions, condition)
	if len(conditions) > maxStoreConditions {
		conditions = conditions[1:]
	}
	store.Status.Conditions = conditions
	log := sm.Reconciler.Log.WithValues("namespace", store.Namespace, "store", store.Name)
	log.Info("Updating store with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := sm.Reconciler.Status().Update(sm.Context, store); err != nil {
		log.Error(err, "Unable to update store")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func initTestFile() (*string, error) {
	name := ".mysql-operator.out"
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "[blaqkube]")
	return &name, err
}

// GetEnvVars returns the environment variables for the store
func (sm *StoreManager) GetEnvVars(store mysqlv1alpha1.Store) (map[string]string, error) {
	em := &EnvManager{
		Client: sm.Reconciler.Client,
		Log:    sm.Reconciler.Log,
	}
	return em.GetEnvVars(sm.Context, store)
}
