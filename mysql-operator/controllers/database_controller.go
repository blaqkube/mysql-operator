package controllers

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// GetPod figures out the pod IP to access the agent
func (r *DatabaseReconciler) GetPod(ctx context.Context, database *mysqlv1alpha1.Database) (string, error) {
	pod := &corev1.Pod{}
	named := types.NamespacedName{
		Name:      database.Spec.Instance + "-0",
		Namespace: database.Namespace,
	}

	condition := metav1.Condition{
		Type:    "mysql-agent",
		Status:  metav1.ConditionTrue,
		Reason:  "pod not found",
		Message: "could not get pod for the backup instance",
	}
	database.Status.LastCondition = "mysql-agent unreachable"
	url := ""
	err := r.Client.Get(ctx, named, pod)
	if err == nil {
		condition = metav1.Condition{
			Type:    "mysql-agent",
			Status:  metav1.ConditionTrue,
			Reason:  "pod found",
			Message: "agent URL found",
		}
		database.Status.LastCondition = "mysql-agent found"
		url = "http://" + pod.Status.PodIP + ":8080"
	}
	conditions := append(database.Status.Conditions, condition)
	database.Status.Conditions = conditions
	if err := r.Client.Status().Update(ctx, database); err != nil {
		return "", err
	}
	return url, err
}

// CreateDatabase is the script that creates a database
func (r *DatabaseReconciler) CreateDatabase(ctx context.Context, database *mysqlv1alpha1.Database, url string) error {

	cfg := agent.NewConfiguration()
	cfg.BasePath = url
	payload := agent.Database{
		Name: database.Spec.Name,
	}
	api := agent.NewAPIClient(cfg)

	condition := metav1.Condition{
		Type:    "database",
		Status:  metav1.ConditionTrue,
		Reason:  "creation failed",
		Message: "could not create the database",
	}
	database.Status.LastCondition = "database creation failed"
	_, _, err := api.MysqlApi.CreateDatabase(ctx, payload, nil)
	if err == nil {
		condition = metav1.Condition{
			Type:    "database",
			Status:  metav1.ConditionTrue,
			Reason:  "creation succeeded",
			Message: "database creation has succeeded",
		}
		database.Status.LastCondition = "database creation succeeded"
	}
	conditions := append(database.Status.Conditions, condition)
	database.Status.Conditions = conditions
	return r.Client.Status().Update(ctx, database)
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases/finalizers,verbs=update

// Reconcile implement the reconciliation loop for databases
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("database", req.NamespacedName)

	// your logic here
	database := &mysqlv1alpha1.Database{}
	if err := r.Client.Get(ctx, req.NamespacedName, database); err != nil {
		log.Error(err, "unable to fetch CronJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if database.Status.LastCondition != "" {
		return ctrl.Result{}, nil
	}

	// Check if this Pod already exists
	url, err := r.GetPod(ctx, database)
	if err != nil {
		return ctrl.Result{}, nil
	}
	r.CreateDatabase(ctx, database, url)
	return ctrl.Result{}, nil

}

// SetupWithManager configure type of events the manager should watch
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Database{}).
		Complete(r)
}
