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

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups/status,verbs=get;update;patch

func (r *BackupReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("backup", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

/*
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
	// cfg.BasePath = "http://localhost:8080"
	api := agent.NewAPIClient(cfg)
	backup := agent.Backup{
		S3access: agent.S3Info{
			Bucket: store.Spec.S3Access.Bucket,
			Path:   store.Spec.S3Access.Path,
			AwsConfig: agent.AwsConfig{
				AwsAccessKeyId:     store.Spec.S3Access.AWSConfig.AccessKey,
				AwsSecretAccessKey: store.Spec.S3Access.AWSConfig.SecretKey,
				Region:             store.Spec.S3Access.AWSConfig.Region,
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
	t := metav1.Now()
	details := mysqlv1alpha1.BackupDetails{
		Location:   b.Location,
		BackupTime: &metav1.Time{Time: b.Timestamp},
	}
	condition := mysqlv1alpha1.ConditionStatus{
		LastProbeTime: &t,
		Status:        b.Status,
		Message:       b.Message,
	}
	instance.Status.LastCondition = b.Status
	instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
	instance.Status.Details = &details
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	go r.MonitorBackup(request.NamespacedName, api.MysqlApi, b.Timestamp.Format(time.RFC3339))
	return reconcile.Result{}, nil
}

func (r *ReconcileBackup) MonitorBackup(n types.NamespacedName, a *agent.MysqlApiService, backup string) {
	reqLogger := log.WithValues("Request.Namespace", n.Namespace, "Request.Name", n.Name)
	endTime := time.Now().Add(60 * time.Second)
	succeeded := false
	instance := &mysqlv1alpha1.Backup{}
	err := r.client.Get(context.TODO(), n, instance)
	if err != nil {
		reqLogger.Info(fmt.Sprintf("Error querying backup: %v", err))
		return
	}
	reqLogger.Info(fmt.Sprintf("Starting to check for backup, current status %s", instance.Status.LastCondition))
	for time.Now().Before(endTime) && !succeeded {
		reqLogger.Info(fmt.Sprintf("Loop..."))
		b, _, err := a.GetBackupByName(context.TODO(), backup, nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		lastCondition := instance.Status.LastCondition
		if b.Status != lastCondition {
			time := metav1.Now()
			instance.Status.LastCondition = b.Status
			condition := mysqlv1alpha1.ConditionStatus{
				LastProbeTime: &time,
				Status:        b.Status,
				Message:       b.Message,
			}
			instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
			err = r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				instance.Status.LastCondition = lastCondition
			}
		}
		if instance.Status.LastCondition == "Available" || instance.Status.LastCondition == "Failed" {
			succeeded = true
			break
		}
		time.Sleep(2 * time.Second)
	}
	if !succeeded {
		time := metav1.Now()
		instance.Status.LastCondition = "Failed"
		condition := mysqlv1alpha1.ConditionStatus{
			LastProbeTime: &time,
			Status:        "Failed",
			Message:       "Backup did not finish in the expected time",
		}
		instance.Status.Conditions = []mysqlv1alpha1.ConditionStatus{condition}
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Could not update status %v", err))
		}
	}
	reqLogger.Info("Monitor backup is now over...")
}
*/

func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Backup{}).
		Complete(r)
}
