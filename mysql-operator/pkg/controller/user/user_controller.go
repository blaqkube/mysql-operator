package user

import (
	"context"
	"fmt"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/pkg/apis/mysql/v1alpha1"
	agent "github.com/blaqkube/mysql-operator/pkg/client-agent"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_user")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new User Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileUser{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("user-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource User
	err = c.Watch(&source.Kind{Type: &mysqlv1alpha1.User{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner User
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &mysqlv1alpha1.User{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileUser implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileUser{}

// ReconcileUser reconciles a User object
type ReconcileUser struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a User object and makes changes based on the state read
// and what is in the User.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileUser) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling User")

	// Fetch the User instance
	instance := &mysqlv1alpha1.User{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	pod := &corev1.Pod{}
	err = r.client.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      instance.Spec.Instance + "-0",
			Namespace: instance.Namespace,
		},
		pod,
	)
	if err != nil {
		time := metav1.Now()
		condition := mysqlv1alpha1.ConditionStatus{
			LastProbeTime: &time,
			Status:        "Failed",
			Message:       fmt.Sprintf("Cannot find pod %s-0; error: %v", instance.Spec.Instance, err),
		}
		instance.Status.LastCondition = "Failed"
		instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	cfg := agent.NewConfiguration()
	cfg.BasePath = "http://" + pod.Status.PodIP + ":8080"
	// cfg.BasePath = "http://localhost:8080"
	api := agent.NewAPIClient(cfg)
	user := agent.User{
		Username: instance.Spec.Username,
		Password: instance.Spec.Password,
	}
	_, _, err = api.MysqlApi.CreateUser(context.TODO(), user, nil)
	if err != nil {
		time := metav1.Now()
		condition := mysqlv1alpha1.ConditionStatus{
			LastProbeTime: &time,
			Status:        "Failed",
			Message:       fmt.Sprintf("Error accessing api: %v", err),
		}
		instance.Status.LastCondition = "Failed"
		instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	t := metav1.Now()
	condition := mysqlv1alpha1.ConditionStatus{
		LastProbeTime: &t,
		Status:        "Succeeded",
		Message:       "User " + instance.Spec.Username + " created",
	}
	instance.Status.LastCondition = "Succeeded"
	instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
