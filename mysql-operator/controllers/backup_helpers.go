package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/blaqkube/mysql-operator/mysql-operator/agent"
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	maxBackupConditions   = 10
	backupPollingInterval = 30 * time.Second
)

var (
	// ErrBackupFailed is reported the backup has failed
	ErrBackupFailed = errors.New("BackupFailed")

	// ErrBackupRunning is reported the backup is always running
	ErrBackupRunning = errors.New("BackupRunning")
)

// BackupManager provides methods to manage the backup subcomponents
type BackupManager struct {
	Context     context.Context
	Reconciler  *BackupReconciler
	TimeManager *TimeManager
}

func (bm *BackupManager) setBackupCondition(backup *mysqlv1alpha1.Backup, condition metav1.Condition, details *mysqlv1alpha1.BackupDetails) (ctrl.Result, error) {
	if condition.Reason == backup.Status.Reason {
		if condition.Reason == mysqlv1alpha1.BackupRunning {
			return ctrl.Result{Requeue: true, RequeueAfter: backupPollingInterval}, nil
		}
		c := len(backup.Status.Conditions) - 1
		d := bm.TimeManager.Next(backup.Status.Conditions[c].LastTransitionTime.Time)
		return ctrl.Result{Requeue: true, RequeueAfter: d}, nil
	}

	backup.Status.Ready = condition.Status
	backup.Status.Reason = condition.Reason
	backup.Status.Message = condition.Message
	conditions := append(backup.Status.Conditions, condition)
	if len(conditions) > maxBackupConditions {
		conditions = conditions[1:]
	}
	backup.Status.Conditions = conditions
	log := bm.Reconciler.Log.WithValues("namespace", backup.Namespace, "backup", backup.Name)
	log.Info("Updating backup with new Status", "Reason", condition.Reason, "Message", condition.Message)
	if err := bm.Reconciler.Status().Update(bm.Context, backup); err != nil {
		log.Error(err, "Unable to update the backup")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// MonitorBackup watch backup progress and update results
func (bm *BackupManager) MonitorBackup(backup *mysqlv1alpha1.Backup) (*mysqlv1alpha1.BackupDetails, error) {
	_ = bm.Reconciler.Log.WithValues("Namespace", backup.Namespace, "backup", backup.Name)
	b := backup.Status.Details

	a := &APIReconciler{
		Client: bm.Reconciler.Client,
		Log:    bm.Reconciler.Log,
	}

	api, err := a.GetAPI(
		bm.Context,
		types.NamespacedName{
			Name:      backup.Spec.Instance,
			Namespace: backup.Namespace,
		},
	)
	if err != nil {
		return b, err
	}
	data, code, err := api.MysqlApi.GetBackupByID(bm.Context, backup.Status.Details.Identifier, nil)
	if err != nil || code.StatusCode != http.StatusOK || data.Status == "Failed" {
		v := metav1.Now()
		b.EndTime = &v
		return b, ErrBackupFailed
	}
	if data.Status == "Running" {
		return b, ErrBackupRunning
	}
	details := &mysqlv1alpha1.BackupDetails{
		Identifier: data.Identifier,
		Bucket:     data.Bucket,
		StartTime:  &metav1.Time{Time: data.StartTime},
		Location:   data.Location,
	}
	if data.EndTime != nil {
		b.EndTime = &metav1.Time{Time: *data.EndTime}
	}
	if data.Status == "Succeeded" {
		return details, nil
	}
	return details, ErrNotImplemented
}

// CreateBackup is the script that creates a user
func (bm *BackupManager) CreateBackup(backup *mysqlv1alpha1.Backup) (*mysqlv1alpha1.BackupDetails, error) {
	log := bm.Reconciler.Log.WithValues("namespace", backup.Namespace, "backup", backup.Name)
	a := &APIReconciler{
		Client: bm.Reconciler.Client,
		Log:    bm.Reconciler.Log,
	}
	store, err := a.GetStore(
		bm.Context,
		types.NamespacedName{
			Name:      backup.Spec.Store,
			Namespace: backup.Namespace,
		},
	)
	if err != nil {
		return nil, err
	}
	api, err := a.GetAPI(
		bm.Context,
		types.NamespacedName{
			Name:      backup.Spec.Instance,
			Namespace: backup.Namespace,
		},
	)
	if err != nil {
		return nil, err
	}

	em := &EnvManager{
		Client: bm.Reconciler.Client,
		Log:    bm.Reconciler.Log,
	}
	envs, err := em.GetEnvVars(bm.Context, *store)
	if err != nil {
		return nil, err
	}
	agentEnvs := []agent.EnvVar{}
	for k, v := range envs {
		agentEnvs = append(agentEnvs, agent.EnvVar{Name: k, Value: v})
	}

	payload := agent.BackupRequest{
		Backend:  string(store.Spec.Backend),
		Bucket:   store.Spec.Bucket,
		Location: fmt.Sprintf("%s/%s-%s.sql", store.Spec.Prefix, backup.Spec.Instance, time.Now().Format("20060102-150405")),
		Envs:     agentEnvs,
	}

	b, response, err := api.MysqlApi.CreateBackup(bm.Context, payload, nil)
	if err != nil || response == nil {
		msg := "NoResponse"
		if err != nil {
			msg = err.Error()
		}
		log.Info(fmt.Sprintf("Could not access agent, error: %s", msg))
		return nil, ErrAgentAccessFailed
	}
	if response.StatusCode != http.StatusCreated {
		log.Info("Agent returned unexpected response", "httpcode", response.StatusCode)
		return nil, ErrAgentRequestFailed
	}
	return &mysqlv1alpha1.BackupDetails{
		Identifier: b.Identifier,
		Bucket:     b.Bucket,
		StartTime:  &metav1.Time{Time: b.StartTime},
		Location:   b.Location,
	}, nil
}

// GetEnvVars returns the environment variables for the store
func (bm *BackupManager) GetEnvVars(store mysqlv1alpha1.Store) (map[string]string, error) {
	em := &EnvManager{
		Client: bm.Reconciler.Client,
		Log:    bm.Reconciler.Log,
	}
	return em.GetEnvVars(bm.Context, store)
}
