package backup

import (
	"context"
	"fmt"
	"time"

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

var log = logf.Log.WithName("controller_backup")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Backup Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileBackup{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("backup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Backup
	err = c.Watch(&source.Kind{Type: &mysqlv1alpha1.Backup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileBackup implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileBackup{}

// ReconcileBackup reconciles a Backup object
type ReconcileBackup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Backup object and makes changes based on the state read
// and what is in the Backup.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Backup")

	// Fetch the Backup instance
	instance := &mysqlv1alpha1.Backup{}
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
	if instance.Status.LastCondition != "" {
		return reconcile.Result{}, nil
	}
	// Check if this Pod already exists
	pod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Instance + "-0", Namespace: instance.Namespace}, pod)
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
	store := &mysqlv1alpha1.Store{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Store, Namespace: instance.Namespace}, store)
	if err != nil {
		time := metav1.Now()
		condition := mysqlv1alpha1.ConditionStatus{
			LastProbeTime: &time,
			Status:        "Failed",
			Message:       fmt.Sprintf("Error accessing store %s: %v", instance.Spec.Store, err),
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
	api := agent.NewAPIClient(cfg)
	backup := agent.Backup{
		S3access: agent.S3Info{
			Bucket: store.Spec.S3Access.Bucket,
			Path:   store.Spec.S3Access.Path,
			Credentials: agent.S3Credentials{
				AwsAccessKeyId:     store.Spec.S3Access.Credentials.AccessKey,
				AwsSecretAccessKey: store.Spec.S3Access.Credentials.SecretKey,
				Region:             store.Spec.S3Access.Credentials.Region,
			},
		},
	}
	b, _, err := api.MysqlApi.CreateBackup(context.TODO(), backup, nil)
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
	time := metav1.Now()
	condition := mysqlv1alpha1.ConditionStatus{
		LastProbeTime: &time,
		Status:        b.Status,
		Message:       b.Message,
	}
	instance.Status.LastCondition = b.Status
	instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	go r.MonitorBackup(request.NamespacedName)
	return reconcile.Result{}, nil
}

func (r *ReconcileBackup) MonitorBackup(n types.NamespacedName) {
	reqLogger := log.WithValues("Request.Namespace", n.Namespace, "Request.Name", n.Name)
	endTime := time.Now().Add(30 * time.Second)
	for time.Now().Before(endTime) {
		time.Sleep(time.Second)
	}
	instance := &mysqlv1alpha1.Backup{}
	err := r.client.Get(context.TODO(), n, instance)
	if err != nil {
		reqLogger.Info(fmt.Sprintf("Error querying backup: %v", err))
		return
	}
	instance.Status.LastCondition = "Zzzzzzzzzzz"
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Info(fmt.Sprintf("Error updating backup: %v", err))
		return
	}
	reqLogger.Info("Monitor backup with success...")
}
