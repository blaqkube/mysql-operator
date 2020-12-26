package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
)

const (
	defaultAgentPort = 8080
)

var (
	// ErrDatabaseNotFound is reported when a database is not found
	ErrDatabaseNotFound = errors.New("DatabaseNotFound")

	// ErrDatabaseNotReady is reported when a database is not ready
	ErrDatabaseNotReady = errors.New("DatabaseNotReady")

	// ErrInstanceNotFound is reported when an instance is not found
	ErrInstanceNotFound = errors.New("InstanceNotFound")

	// ErrInstanceNotReady is reported when an instance is not ready
	ErrInstanceNotReady = errors.New("InstanceNotReady")

	// ErrPodNotFound is reported when a pod is not found
	ErrPodNotFound = errors.New("PodNotFound")

	// ErrAgentAccessFailed is reported when the agent could not be accessed
	ErrAgentAccessFailed = errors.New("AgentAccessFailed")

	// ErrAgentRequestFailed is reported when the agent request fails
	ErrAgentRequestFailed = errors.New("AgentRequestFailed")
)

// APIReconciler reconciles an object
type APIReconciler struct {
	client.Client
	Log logr.Logger
}

// GetAPI is the script that gets the API from an instance name and namespace
func (a *APIReconciler) GetAPI(ctx context.Context, instanceName types.NamespacedName) (*agent.APIClient, error) {
	log := a.Log.WithValues("namespace", instanceName.Namespace, "instance", instanceName.Name)

	instance := &mysqlv1alpha1.Instance{}
	if err := a.Client.Get(ctx, instanceName, instance); err != nil {
		log.Info("Unable to fetch instance", "instance", instanceName.Name)
		return nil, ErrInstanceNotFound
	}
	if instance.Status.Reason != mysqlv1alpha1.InstanceStatefulSetReady {
		log.Info("Instance is not ready yet")
		return nil, ErrInstanceNotReady
	}
	pod := &corev1.Pod{}
	podName := types.NamespacedName{
		Name:      instanceName.Name + "-0",
		Namespace: instanceName.Namespace,
	}
	if err := a.Client.Get(ctx, podName, pod); err != nil {
		log.Info("Could not access pod", "pod", podName.Name)
		return nil, ErrPodNotFound
	}
	url := fmt.Sprintf("http://%s:%d", pod.Status.PodIP, defaultAgentPort)

	cfg := agent.NewConfiguration()
	cfg.BasePath = url
	return agent.NewAPIClient(cfg), nil

}
