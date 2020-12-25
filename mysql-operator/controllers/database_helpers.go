package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
)

const (
	defaultAgentPort      = 8080
	maxDatabaseConditions = 10
)

var (
	// ErrPodNotFound is reported when a pod is not found
	ErrPodNotFound = errors.New("PodNotFound")

	// ErrAgentAccessFailed is reported when the agent could not be accessed
	ErrAgentAccessFailed = errors.New("AgentAccessFailed")

	// ErrAgentRequestFailed is reported when the agent request fails
	ErrAgentRequestFailed = errors.New("AgentRequestFailed")
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
		return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
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
	log.Info("Updating store with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := dm.Reconciler.Status().Update(dm.Context, database); err != nil {
		log.Error(err, "Unable to update database")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// CreateDatabase is the script that creates a database
func (dm *DatabaseManager) CreateDatabase(database *mysqlv1alpha1.Database) error {
	log := dm.Reconciler.Log.WithValues("namespace", database.Namespace, "database", database.Name)
	pod := &corev1.Pod{}
	podName := types.NamespacedName{
		Name:      database.Spec.Instance + "-0",
		Namespace: database.Namespace,
	}
	if err := dm.Reconciler.Client.Get(dm.Context, podName, pod); err != nil {
		log.Info("Could not access pod", "pod", podName.Name)
		return ErrPodNotFound
	}
	url := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, defaultAgentPort)

	cfg := agent.NewConfiguration()
	cfg.BasePath = url
	payload := agent.Database{
		Name: database.Spec.Name,
	}
	api := agent.NewAPIClient(cfg)

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
