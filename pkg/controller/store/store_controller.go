package store

import (
	"context"
	"fmt"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/pkg/apis/mysql/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_store")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Store Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStore{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("store-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Store
	err = c.Watch(&source.Kind{Type: &mysqlv1alpha1.Store{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileStore implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileStore{}

// ReconcileStore reconciles a Store object
type ReconcileStore struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Store object and makes changes based on the state read
// and what is in the Store.Spec
func (r *ReconcileStore) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Store")

	// Fetch the Store instance
	instance := &mysqlv1alpha1.Store{}
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
	if instance.Status.LastCondition == "" {
		instance.Status.LastCondition = "Pending"
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Store updated successfully - don't requeue
		return reconcile.Result{}, nil
	}
	if instance.Status.LastCondition == "Pending" {
		time := metav1.Now()
		condition := mysqlv1alpha1.ConditionStatus{
			LastProbeTime: &time,
			Status:        "Pending",
			Message:       "",
		}
		err = TestS3Connection(
			instance.Spec.S3Access.Credentials.AccessKey,
			instance.Spec.S3Access.Credentials.SecretKey,
			instance.Spec.S3Access.Credentials.Region,
			instance.Spec.S3Access.Bucket,
			instance.Spec.S3Access.Path)
		if err != nil {
			condition.Status = "Error"
			condition.Message = fmt.Sprintf("%v", err)
			instance.Status.LastCondition = "Error"
			instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Store updated successfully - don't requeue
			return reconcile.Result{}, nil
		}
		instance.Status.LastCondition = "Success"
		condition.Status = "Success"
		condition.Message = fmt.Sprintf("File %s/manifest.txt successfully written in s3://%s", instance.Spec.S3Access.Path, instance.Spec.S3Access.Bucket)
		instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Store updated successfully - don't requeue
		return reconcile.Result{}, nil
	}
	reqLogger.Info("Skip reconcile: store exists and LastCondition updated")
	return reconcile.Result{}, nil
}
