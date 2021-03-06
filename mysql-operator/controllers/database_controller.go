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

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases/finalizers,verbs=update

// Reconcile implement the reconciliation loop for databases
func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("database", req.NamespacedName)
	log.Info("Running a reconcile loop")

	// your logic here
	database := &mysqlv1alpha1.Database{}
	if err := r.Get(ctx, req.NamespacedName, database); err != nil {
		log.Info("Unable to fetch database from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if database.Status.Reason == mysqlv1alpha1.DatabaseSucceeded {
		return ctrl.Result{}, nil
	}

	dm := &DatabaseManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	err := dm.CreateDatabase(database)
	if err != nil {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
		}
		switch err {
		case ErrInstanceNotFound:
			condition.Reason = mysqlv1alpha1.DatabaseInstanceAccessError
			condition.Message = "Could not find the instance"
		case ErrInstanceNotReady:
			condition.Reason = mysqlv1alpha1.DatabaseInstanceNotReady
			condition.Message = "The instance is not available"
		case ErrPodNotFound:
			condition.Reason = mysqlv1alpha1.DatabaseAgentNotFound
			condition.Message = "Could not find the agent"
		default:
			condition.Reason = mysqlv1alpha1.DatabaseAgentFailed
			condition.Message = fmt.Sprintf("Unexpected failure with agent: %v", err)
		}
		return dm.setDatabaseCondition(database, condition)
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.DatabaseSucceeded,
		Message:            fmt.Sprintf("Database %s successfully created", database.Spec.Name),
	}
	return dm.setDatabaseCondition(database, condition)
}

// SetupWithManager configure type of events the manager should watch
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Database{}).
		Complete(r)
}
