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

// GrantReconciler reconciles a Grant object
type GrantReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=grants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=grants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=grants/finalizers,verbs=update

// Reconcile implement the reconciliation loop for users
func (r *GrantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("grant", req.NamespacedName)
	log.Info("Running a reconcile loop")

	// Fetch the Grant instance
	grant := &mysqlv1alpha1.Grant{}
	if err := r.Get(ctx, req.NamespacedName, grant); err != nil {
		log.Info("Unable to fetch grant from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if grant.Status.Reason == mysqlv1alpha1.GrantSucceeded {
		return ctrl.Result{}, nil
	}

	um := &GrantManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	err := um.CreateGrant(grant)
	if err != nil {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
		}
		switch err {
		case ErrNotImplemented:
			condition.Reason = mysqlv1alpha1.GrantNotImplemented
			condition.Message = "This feature is not implemented yet"
		case ErrUserDatabaseMismatch:
			condition.Reason = mysqlv1alpha1.GrantUserDatabaseMismatch
			condition.Message = "Instance from user and database does not match"
		case ErrDatabaseNotFound:
			condition.Reason = mysqlv1alpha1.GrantDatabaseAccessError
			condition.Message = "Could not find the database"
		case ErrDatabaseNotReady:
			condition.Reason = mysqlv1alpha1.GrantDatabaseNotReady
			condition.Message = "The database is not ready"
		case ErrUserNotFound:
			condition.Reason = mysqlv1alpha1.GrantUserAccessError
			condition.Message = "Could not find the user"
		case ErrUserNotReady:
			condition.Reason = mysqlv1alpha1.GrantUserNotReady
			condition.Message = "The user is not ready"
		case ErrInstanceNotFound:
			condition.Reason = mysqlv1alpha1.GrantInstanceAccessError
			condition.Message = "Could not find the instance"
		case ErrInstanceNotReady:
			condition.Reason = mysqlv1alpha1.GrantInstanceNotReady
			condition.Message = "The Instance is not ready"
		case ErrPodNotFound:
			condition.Reason = mysqlv1alpha1.GrantAgentNotFound
			condition.Message = "Could not find the agent"
		default:
			condition.Reason = mysqlv1alpha1.GrantAgentFailed
			condition.Message = fmt.Sprintf("Unexpected failure with agent: %v", err)
		}
		return um.setGrantCondition(grant, condition)
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.GrantSucceeded,
		Message:            fmt.Sprintf("Grant %s on %s for %s successful", grant.Spec.AccessMode, grant.Spec.User, grant.Spec.Database),
	}
	return um.setGrantCondition(grant, condition)
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Grant{}).
		Complete(r)
}
