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

	"github.com/go-logr/logr"
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

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=users/status,verbs=get;update;patch

// Reconcile implement the reconciliation loop for users
func (r *UserReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("user", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

/*
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
		Grants:   []agent.Grant{},
	}
	for _, v := range instance.Spec.Grants {
		user.Grants = append(user.Grants, agent.Grant{
			Database:   v.Database,
			AccessMode: v.AccessMode,
		})
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
*/

// SetupWithManager configure type of events the manager should watch
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.User{}).
		Complete(r)
}
