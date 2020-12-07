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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Properties StatefulSetProperties
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=instances/status,verbs=get;update;patch

// Reconcile implement the reconciliation loop for instances
func (r *InstanceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("instance", req.NamespacedName)

	// your logic here
	var instance mysqlv1alpha1.Instance
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		log.Error(err, "unable to fetch Store")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var store mysqlv1alpha1.Store
	NamespacedStore := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Spec.Restore.Store}
	if instance.Status.LastCondition == "" || instance.Status.LastCondition == "Waiting for store" {
		if instance.Spec.Restore.Store != "" {
			if err := r.Get(ctx, NamespacedStore, &store); err != nil {
				log.Error(err, "unable to fetch Store")
				instance.Status.LastCondition = "Waiting for store"
				if err := r.Status().Update(ctx, &instance); err != nil {
					log.Error(err, "unable to update instance status")
					return ctrl.Result{}, err
				}
				return ctrl.Result{Requeue: true, RequeueAfter: time.Duration(30 * time.Second)}, nil
			}
		}
		return r.CreateOrUpdateStafefulSet(&instance, &store, instance.Spec.Restore.FilePath)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager configure type of events the manager should watch
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Instance{}).
		Complete(r)
}
