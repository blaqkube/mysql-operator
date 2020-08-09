/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

func (r *StoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("store", req.NamespacedName)

	// your logic here
	var store mysqlv1alpha1.Store
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Error(err, "unable to fetch Store")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if store.Status.Status == "" {
		store.Status.Status = "Pending"
		if err := r.Status().Update(ctx, &store); err != nil {
			log.Error(err, "unable to update store status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	if store.Status.Status == "Pending" {
		if store.Spec.Backend == nil || *store.Spec.Backend == "s3" {

			s, err := r.BackupStore.New(store.Spec.S3Backup.AWSConfig)
			if err != nil {
				store.Status.Status = "Error"
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Error(err, "unable to update store status")
					return ctrl.Result{}, err
				}
			}
			err = s.TestS3Access("test", "/validation")
			if err != nil {
				store.Status.Status = "Error"
				if err := r.Status().Update(ctx, &store); err != nil {
					log.Error(err, "unable to update store status")
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, err
			}
			store.Status.Status = "Success"
			if err := r.Status().Update(ctx, &store); err != nil {
				log.Error(err, "unable to update store status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}
	return ctrl.Result{}, nil
}

func (r *StoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Store{}).
		Complete(r)
}
