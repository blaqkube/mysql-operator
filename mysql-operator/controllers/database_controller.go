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
	"github.com/operator-framework/operator-lib/status"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *DatabaseReconciler) GetPod(database *mysqlv1alpha1.Database) (string, error) {
	ctx := context.Background()
	pod := &corev1.Pod{}
	named := types.NamespacedName{
		Name:      database.Spec.Instance + "-0",
		Namespace: database.Namespace,
	}

	condition := status.Condition{
		Type:    status.ConditionType("mysql-agent"),
		Status:  corev1.ConditionTrue,
		Reason:  status.ConditionReason("pod not found"),
		Message: "could not get pod for the backup instance",
	}
	database.Status.LastCondition = "mysql-agent unreachable"
	url := ""
	err := r.Client.Get(ctx, named, pod)
	if err == nil {
		condition = status.Condition{
			Type:    status.ConditionType("mysql-agent"),
			Status:  corev1.ConditionTrue,
			Reason:  status.ConditionReason("pod found"),
			Message: "agent URL found",
		}
		database.Status.LastCondition = "mysql-agent found"
		url = "http://" + pod.Status.PodIP + ":8080"
	}
	(&database.Status.Conditions).SetCondition(condition)
	if err := r.Client.Status().Update(ctx, database); err != nil {
		return "", err
	}
	return url, err
}

func (r *DatabaseReconciler) CreateDatabase(database *mysqlv1alpha1.Database, url string) error {
	ctx := context.Background()

	cfg := agent.NewConfiguration()
	cfg.BasePath = url
	payload := agent.Database{
		Name: database.Spec.Name,
	}
	api := agent.NewAPIClient(cfg)

	condition := status.Condition{
		Type:    status.ConditionType("database"),
		Status:  corev1.ConditionTrue,
		Reason:  status.ConditionReason("creation failed"),
		Message: "could not create the database",
	}
	database.Status.LastCondition = "database creation failed"
	_, _, err := api.MysqlApi.CreateDatabase(ctx, payload, nil)
	if err == nil {
		condition = status.Condition{
			Type:    status.ConditionType("database"),
			Status:  corev1.ConditionTrue,
			Reason:  status.ConditionReason("creation succeeded"),
			Message: "database creation has succeeded",
		}
		database.Status.LastCondition = "database creation succeeded"
	}
	(&database.Status.Conditions).SetCondition(condition)
	return r.Client.Status().Update(ctx, database)
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=databases/status,verbs=get;update;patch
func (r *DatabaseReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("database", req.NamespacedName)

	// your logic here
	database := &mysqlv1alpha1.Database{}
	if err := r.Client.Get(ctx, req.NamespacedName, database); err != nil {
		log.Error(err, "unable to fetch CronJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if database.Status.LastCondition != "" {
		return ctrl.Result{}, nil
	}

	// Check if this Pod already exists
	url, err := r.GetPod(database)
	if err != nil {
		return ctrl.Result{}, nil
	}
	r.CreateDatabase(database, url)
	return ctrl.Result{}, nil

}

func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Database{}).
		Complete(r)
}
