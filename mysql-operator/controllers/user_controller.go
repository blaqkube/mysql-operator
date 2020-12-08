package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-lib/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users/status,verbs=get;update;patch

// Reconcile implement the reconciliation loop for users
func (r *UserReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("user", req.NamespacedName)

	log.Info("Reconciling User")

	// Fetch the User instance
	user := &mysqlv1alpha1.User{}
	err := r.Client.Get(ctx, req.NamespacedName, user)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Check if this Pod already exists
	pod := &corev1.Pod{}
	err = r.Client.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      user.Spec.Instance + "-0",
			Namespace: user.Namespace,
		},
		pod,
	)
	if err != nil {
		t := metav1.Now()
		condition := status.Condition{
			Type:               status.ConditionType("podmonitor"),
			Status:             corev1.ConditionTrue,
			Reason:             status.ConditionReason("Failed"),
			Message:            fmt.Sprintf("Cannot find pod %s-0; error: %v", user.Spec.Instance, err),
			LastTransitionTime: t,
		}
		user.Status.LastCondition = "Failed"
		user.Status.Conditions = append(user.Status.Conditions, condition)
		err = r.Client.Status().Update(context.TODO(), user)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	cfg := agent.NewConfiguration()
	cfg.BasePath = "http://" + pod.Status.PodIP + ":8080"
	api := agent.NewAPIClient(cfg)
	u := agent.User{
		Username: user.Spec.Username,
		Password: user.Spec.Password,
		Grants:   []agent.Grant{},
	}
	for _, v := range user.Spec.Grants {
		u.Grants = append(u.Grants, agent.Grant{
			Database:   v.Database,
			AccessMode: v.AccessMode,
		})
	}
	_, _, err = api.MysqlApi.CreateUser(ctx, u, nil)
	if err != nil {
		t := metav1.Now()
		condition := status.Condition{
			Type:               status.ConditionType("user"),
			Status:             corev1.ConditionTrue,
			Reason:             status.ConditionReason("Failed"),
			Message:            fmt.Sprintf("Cannot create user %s, error: %v", u.Username, err),
			LastTransitionTime: t,
		}
		user.Status.LastCondition = "Failed"
		user.Status.Conditions = append(user.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, user)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	t := metav1.Now()
	condition := status.Condition{
		Type:               status.ConditionType("user"),
		Status:             corev1.ConditionTrue,
		Reason:             status.ConditionReason("Succeeded"),
		Message:            fmt.Sprintf("User %s created", u.Username),
		LastTransitionTime: t,
	}
	user.Status.LastCondition = "Failed"
	user.Status.Conditions = append(user.Status.Conditions, condition)
	err = r.Client.Status().Update(ctx, user)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager configure type of events the manager should watch
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.User{}).
		Complete(r)
}
