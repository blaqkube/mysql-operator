package controllers

import (
	"context"

	"github.com/blaqkube/mysql-operator/mysql-operator/helpers"
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	BackupStore helpers.StoreInitializer
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=stores/status,verbs=get;update;patch

// Reconcile implement the reconciliation loop for stores
func (r *StoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("store", req.NamespacedName)

	// your logic here
	var store mysqlv1alpha1.Store
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Error(err, "unable to fetch Store")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if store.Status.LastCondition == "" {
		store.Status.LastCondition = "Pending"
		if err := r.Status().Update(ctx, &store); err != nil {
			log.Error(err, "unable to update store status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if store.Status.LastCondition == "Pending" {
		if store.Spec.Backend == nil || *store.Spec.Backend == "s3" {

			s, err := r.BackupStore.New(store.Spec.S3Backup.AWSConfig)
			if err != nil {
				store.Status.LastCondition = "Error"
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Error(err, "unable to update store status")
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
			err = s.TestS3Access(store.Spec.S3Backup.Bucket, "/validation")
			if err != nil {
				store.Status.LastCondition = "Error"
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Error(err, "unable to update store status")
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
			store.Status.LastCondition = "Success"
			if err := r.Status().Update(ctx, &store); err != nil {
				log.Error(err, "unable to update store status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
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
