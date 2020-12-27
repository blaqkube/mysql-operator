package controllers

import (
	"context"
	"fmt"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/blaqkube/mysql-operator/agent/backend"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Storages map[string]backend.Storage
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores/finalizers,verbs=update

// Reconcile implement the reconciliation loop for stores
func (r *StoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("store", req.NamespacedName)
	log.Info("Running a reconcile loop")

	// TODO: Reconciler should be able to
	// - detect a change in the ConfigMap or Secret and reload the associated data
	// - Retry on regular basis in the event of a failure
	// - Update databases status when store moves to success
	var store mysqlv1alpha1.Store
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Info("Unable to fetch store from kubernetes")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	sm := &StoreManager{
		Context:     ctx,
		Reconciler:  r,
		TimeManager: NewTimeManager(),
	}

	if store.Status.Reason == "" || store.Status.CheckRequested == true {
		store.Status.CheckRequested = false
		condition := metav1.Condition{
			Type:               "available",
			Status:             metav1.ConditionUnknown,
			LastTransitionTime: metav1.Now(),
			Reason:             mysqlv1alpha1.StoreCheckRequested,
			Message:            "A new check has been requested",
		}
		return sm.setStoreCondition(&store, condition)
	}

	if store.Status.Reason == mysqlv1alpha1.StoreCheckRequested {
		storage := "s3"
		if store.Spec.Backend != "" {
			storage = string(store.Spec.Backend)
		}
		if storage == "s3" || storage == "blackhole" || storage == "gcp" {
			store.Status.CheckRequested = false
			envs, err := sm.GetEnvVars(store)
			if err != nil {
				condition := metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.StoreCheckFailed,
					Message:            "Cannot access values for envs",
				}
				return sm.setStoreCondition(&store, condition)
			}
			e := []openapi.EnvVar{}
			for k := range envs {
				e = append(e, openapi.EnvVar{Name: k, Value: envs[k]})
			}
			filename, err := initTestFile()
			if err != nil {
				condition := metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.StoreCheckFailed,
					Message:            fmt.Sprintf("Cannot initialize local file, error: %v", err),
				}
				return sm.setStoreCondition(&store, condition)
			}
			request := &openapi.BackupRequest{
				Bucket:   store.Spec.Bucket,
				Location: "/blaqkube/.mysql-operator.out",
				Envs:     e,
			}
			log.Info("Checking access for bucket", "bucket", request.Bucket)
			err = r.Storages[storage].Push(request, *filename)
			if err == nil {
				err = r.Storages[storage].Delete(request)
			}
			if err != nil {
				condition := metav1.Condition{
					Type:               "available",
					Status:             metav1.ConditionFalse,
					LastTransitionTime: metav1.Now(),
					Reason:             mysqlv1alpha1.StoreCheckFailed,
					Message:            fmt.Sprintf("Cannot write to bucket, error: %v", err),
				}
				return sm.setStoreCondition(&store, condition)
			}
			condition := metav1.Condition{
				Type:               "available",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             mysqlv1alpha1.StoreCheckSucceeded,
				Message:            "The check has succeeded",
			}
			return sm.setStoreCondition(&store, condition)
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager configure type of events the manager should watch
func (r *StoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Store{}).
		Complete(r)
}
