package controllers

import (
	"context"
	"fmt"
	"net/http"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
)

const (
	maxDatabaseConditions = 10
)

// DatabaseManager provides methods to manage database subcomponents
type DatabaseManager struct {
	Context     context.Context
	Reconciler  *DatabaseReconciler
	TimeManager *TimeManager
}

func (dm *DatabaseManager) setDatabaseCondition(database *mysqlv1alpha1.Database, condition metav1.Condition) (ctrl.Result, error) {
	if condition.Reason == database.Status.Reason {
		c := len(database.Status.Conditions) - 1
		d := dm.TimeManager.Next(database.Status.Conditions[c].LastTransitionTime.Time)
		if condition.Reason != mysqlv1alpha1.DatabaseSucceeded {
			return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
		}
		return ctrl.Result{}, nil
	}
	database.Status.Ready = condition.Status
	database.Status.Reason = condition.Reason
	database.Status.Message = condition.Message
	conditions := append(database.Status.Conditions, condition)
	if len(conditions) > maxDatabaseConditions {
		conditions = conditions[1:]
	}
	database.Status.Conditions = conditions
	log := dm.Reconciler.Log.WithValues("namespace", database.Namespace, "database", database.Name)
	log.Info("Updating database with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := dm.Reconciler.Status().Update(dm.Context, database); err != nil {
		log.Error(err, "Unable to update database")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// CreateDatabase is the script that creates a database
func (dm *DatabaseManager) CreateDatabase(database *mysqlv1alpha1.Database) error {
	log := dm.Reconciler.Log.WithValues("namespace", database.Namespace, "database", database.Name)

	a := &APIReconciler{
		Client: dm.Reconciler.Client,
		Log:    dm.Reconciler.Log,
	}
	api, err := a.GetAPI(
		dm.Context,
		types.NamespacedName{
			Name:      database.Spec.Instance,
			Namespace: database.Namespace,
		},
	)
	if err != nil {
		return err
	}
	payload := agent.Database{
		Name: database.Spec.Name,
	}

	_, response, err := api.MysqlApi.CreateDatabase(dm.Context, payload, nil)
	if err != nil || response == nil {
		msg := "NoResponse"
		if err != nil {
			msg = err.Error()
		}
		log.Info(fmt.Sprintf("Could not access agent, error: %s", msg))
		return ErrAgentAccessFailed
	}
	if response.StatusCode != http.StatusCreated {
		log.Info("Agent returned unexpected response", "httpcode", response.StatusCode)
		return ErrAgentRequestFailed
	}
	return nil
}
