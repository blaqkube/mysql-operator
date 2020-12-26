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

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users/finalizers,verbs=update

// Reconcile implement the reconciliation loop for users
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("user", req.NamespacedName)
	log.Info("Running a reconcile loop")

	// Fetch the User instance
	user := &mysqlv1alpha1.User{}
	if err := r.Get(ctx, req.NamespacedName, user); err != nil {
		log.Info("Unable to fetch user from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if user.Status.Reason == mysqlv1alpha1.UserSucceeded {
		return ctrl.Result{}, nil
	}

	um := &UserManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	err := um.CreateUser(user)
	if err != nil {
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
		}
		switch err {
		case ErrMissingPassword:
			condition.Reason = mysqlv1alpha1.UserPasswordError
			condition.Message = "Missing password from user"
		case ErrKeyNotFound, ErrConfigMapNotFound, ErrSecretNotFound:
			condition.Reason = mysqlv1alpha1.UserPasswordAccessError
			condition.Message = "Could not find the passsord from its definition"
		case ErrInstanceNotFound:
			condition.Reason = mysqlv1alpha1.UserInstanceAccessError
			condition.Message = "Could not find the instance"
		case ErrInstanceNotReady:
			condition.Reason = mysqlv1alpha1.UserInstanceNotReady
			condition.Message = "The Instance is not ready"
		case ErrPodNotFound:
			condition.Reason = mysqlv1alpha1.UserAgentNotFound
			condition.Message = "Could not find the agent"
		default:
			condition.Reason = mysqlv1alpha1.UserAgentFailed
			condition.Message = fmt.Sprintf("Unexpected failure with agent: %v", err)
		}
		return um.setUserCondition(user, condition)
	}
	condition := metav1.Condition{
		Type:               "available",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             mysqlv1alpha1.UserSucceeded,
		Message:            fmt.Sprintf("User %s successfully created", user.Spec.Username),
	}
	return um.setUserCondition(user, condition)
}

// SetupWithManager configure type of events the manager should watch
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.User{}).
		Complete(r)
}
