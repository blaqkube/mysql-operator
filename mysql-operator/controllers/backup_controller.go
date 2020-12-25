package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	uuid "github.com/hashicorp/go-uuid"
)

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Properties StatefulSetProperties
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mysql.blaqkube.io,resources=backups/finalizers,verbs=update

// Reconcile implement the reconciliation loop for backups
func (r *BackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("backup", req.NamespacedName)
	log.Info("Running a reconcile loop")

	// your logic here
	backup := &mysqlv1alpha1.Backup{}
	err := r.Client.Get(ctx, req.NamespacedName, backup)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if backup.Status.Reason != "" {
		return reconcile.Result{}, nil
	}
	// Check if this Pod already exists
	pod := &corev1.Pod{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: backup.Spec.Instance + "-0", Namespace: backup.Namespace}, pod)
	if err != nil {
		t := metav1.Now()
		condition := metav1.Condition{
			Type:               "podmonitor",
			Status:             metav1.ConditionTrue,
			Reason:             "Failed",
			Message:            fmt.Sprintf("Cannot find pod %s-0; error: %v", backup.Spec.Instance, err),
			LastTransitionTime: t,
		}
		backup.Status.Reason = "Failed"
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	store := &mysqlv1alpha1.Store{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: backup.Spec.Store, Namespace: backup.Namespace}, store)
	if err != nil {
		t := metav1.Now()
		condition := metav1.Condition{
			Type:               "store",
			Status:             metav1.ConditionTrue,
			Reason:             "Failed",
			Message:            fmt.Sprintf("Error accessing store %s: %v", backup.Spec.Store, err),
			LastTransitionTime: t,
		}
		backup.Status.Reason = "Failed"
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	cfg := agent.NewConfiguration()
	keys, err := r.GetEnvVars(ctx, *store)
	envs := []agent.EnvVar{}
	for i := range keys {
		env := agent.EnvVar{Name: i, Value: keys[i]}
		envs = append(envs, env)
	}
	cfg.BasePath = "http://" + pod.Status.PodIP + ":8080"
	api := agent.NewAPIClient(cfg)
	id, _ := uuid.GenerateUUID()
	location := fmt.Sprintf("%s/%s.dmp", store.Spec.Prefix, id)
	payload := agent.BackupRequest{
		Bucket:   store.Spec.Bucket,
		Location: location,
		Envs:     envs,
	}
	b, _, err := api.MysqlApi.CreateBackup(context.TODO(), payload, nil)
	if err != nil {
		t := metav1.Now()
		condition := metav1.Condition{
			Type:               "backup",
			Status:             metav1.ConditionTrue,
			Reason:             "Failed",
			Message:            fmt.Sprintf("Error accessing api: %v", err),
			LastTransitionTime: t,
		}
		backup.Status.Reason = "Failed"
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(ctx, backup)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	t := metav1.Now()
	details := mysqlv1alpha1.BackupDetails{
		Location:   b.Location,
		BackupTime: &metav1.Time{Time: b.StartTime},
	}
	condition := metav1.Condition{
		Type:               "backup",
		Status:             metav1.ConditionTrue,
		Reason:             b.Status,
		LastTransitionTime: t,
	}

	backup.Status.Reason = b.Status
	backup.Status.Conditions = append(backup.Status.Conditions, condition)
	backup.Status.Details = &details
	err = r.Client.Status().Update(context.TODO(), backup)
	if err != nil {
		return reconcile.Result{}, err
	}
	go r.MonitorBackup(req.NamespacedName, api.MysqlApi, b.StartTime.Format(time.RFC3339))
	return ctrl.Result{}, nil
}

// MonitorBackup watch backup progress and update results
func (r *BackupReconciler) MonitorBackup(n types.NamespacedName, a *agent.MysqlApiService, backupName string) {
	log := r.Log.WithValues("Request.Namespace", n.Namespace, "Request.Name", n.Name)
	endTime := time.Now().Add(60 * time.Second)
	succeeded := false
	backup := &mysqlv1alpha1.Backup{}
	err := r.Client.Get(context.TODO(), n, backup)
	if err != nil {
		log.Info(fmt.Sprintf("Error querying backup: %v", err))
		return
	}
	log.Info(fmt.Sprintf("Starting to check for backup, current status %s", backup.Status.Reason))
	for time.Now().Before(endTime) && !succeeded {
		log.Info(fmt.Sprintf("Loop..."))
		b, _, err := a.GetBackups(context.TODO(), nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		reason := backup.Status.Reason
		if b.Items[0].Status != reason {
			t := metav1.Now()
			backup.Status.Reason = b.Items[0].Status
			condition := metav1.Condition{
				Type:               "backup",
				Status:             metav1.ConditionTrue,
				Reason:             b.Items[0].Status,
				LastTransitionTime: t,
			}
			backup.Status.Conditions = append(backup.Status.Conditions, condition)
			err = r.Client.Status().Update(context.TODO(), backup)
			if err != nil {
				backup.Status.Reason = reason
			}
		}
		if backup.Status.Reason == "Available" || backup.Status.Reason == "Failed" {
			succeeded = true
			break
		}
		time.Sleep(2 * time.Second)
	}
	if !succeeded {
		t := metav1.Now()
		backup.Status.Reason = "Failed"
		condition := metav1.Condition{
			Type:               "backup",
			Status:             metav1.ConditionTrue,
			Reason:             "Failed",
			Message:            "Backup did not finish in the expected time",
			LastTransitionTime: t,
		}
		backup.Status.Conditions = append(backup.Status.Conditions, condition)
		err = r.Client.Status().Update(context.TODO(), backup)
		if err != nil {
			log.Info(fmt.Sprintf("Could not update status %v", err))
		}
	}
	log.Info("Monitor backup is now over...")
}

// SetupWithManager configure type of events the manager should watch
func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mysqlv1alpha1.Backup{}).
		Complete(r)
}

// GetEnvVars returns the environment variables for the store
func (r *BackupReconciler) GetEnvVars(ctx context.Context, store mysqlv1alpha1.Store) (map[string]string, error) {
	output := map[string]string{}
	configMaps := map[string]corev1.ConfigMap{}
	secrets := map[string]corev1.Secret{}
	for _, envVar := range store.Spec.Envs {
		if envVar.Name == "" {
			return nil, errors.New("MissingVariable")
		}
		if envVar.Value != "" {
			output[envVar.Name] = envVar.Value
			continue
		}
		if envVar.ValueFrom != nil {
			value := ""
			switch {
			case envVar.ValueFrom.ConfigMapKeyRef != nil:
				cm := envVar.ValueFrom.ConfigMapKeyRef
				namespace := store.Namespace
				name := cm.Name
				key := cm.Key
				optional := cm.Optional != nil && *cm.Optional
				configMap, ok := configMaps[name]
				if !ok {
					configMap := corev1.ConfigMap{}
					if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &configMap); err != nil {
						if !optional {
							return nil, err
						}
					} else {
						configMaps[name] = configMap
						ok = true
					}
				}
				if ok {
					value, ok = configMap.Data[key]
					if !ok && !optional {
						return nil, errors.New("MissingVariable")
					}
				}
			case envVar.ValueFrom.SecretKeyRef != nil:
				s := envVar.ValueFrom.SecretKeyRef
				namespace := store.Namespace
				name := s.Name
				key := s.Key
				optional := s.Optional != nil && *s.Optional
				secret, ok := secrets[name]
				if !ok {
					secret := corev1.Secret{}
					if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &secret); err != nil {
						if !optional {
							return nil, err
						}
					} else {
						secrets[name] = secret
						ok = true
					}
				}
				if ok {
					valueBytes, ok := secret.Data[key]
					if !ok && !optional {
						return nil, errors.New("MissingVariable")
					}
					if ok {
						value = string(valueBytes)
					}
				}
			}
			output[envVar.Name] = value
			continue
		}
		return nil, errors.New("MissingVariable")
	}
	return output, nil
}
