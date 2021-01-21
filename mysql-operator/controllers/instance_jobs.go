package controllers

import (
	"context"
	"fmt"

	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
)

// BackupJob is a struct that manages Jobs for backups
type BackupJob struct {
	client.Client
	Instance types.NamespacedName
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

// NewBackupJob creates a BackupJob to schedule it
func NewBackupJob(client client.Client, instance types.NamespacedName, log logr.Logger, scheme *runtime.Scheme) *BackupJob {
	return &BackupJob{
		Client:   client,
		Instance: instance,
		Log:      log,
		Scheme:   scheme,
	}
}

// Run implement the Job interface to use with Cron AddFunc()
func (b *BackupJob) Run() {
	ctx := context.Background()
	instance := mysqlv1alpha1.Instance{}
	if err := b.Client.Get(ctx, b.Instance, &instance); err != nil {
		b.Log.Info(fmt.Sprintf("job for %s/%s failed. Could not access instance...", b.Instance.Namespace, b.Instance.Name))
		return
	}
	b.Log.Info(fmt.Sprintf("job for %s/%s succeeded...", b.Instance.Namespace, b.Instance.Name))
	backupName := fmt.Sprintf("%s-backup-%s", instance.Name, time.Now().Format("20060102-150405"))
	backup := &mysqlv1alpha1.Backup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupName,
			Namespace: instance.Namespace,
		},
		Spec: mysqlv1alpha1.BackupSpec{
			Store:    instance.Spec.BackupSchedule.Store,
			Instance: instance.Name,
		},
	}
	if err := controllerutil.SetControllerReference(&instance, backup, b.Scheme); err != nil {
		b.Log.Info(fmt.Sprintf("Error registering backup %s/%s with instance %s", instance.Namespace, backupName, instance.Name))
		return
	}
	if err := b.Client.Create(ctx, backup); err != nil {
		b.Log.Info(fmt.Sprintf("Error creating backup %s/%s for instance %s", instance.Namespace, backupName, instance.Name))
		return
	}
	b.Log.Info(fmt.Sprintf("Backup %s/%s for instance %s successfully created", instance.Namespace, backupName, instance.Name))
}

// MaintenanceJob is a struct that manages Jobs for maintenance
type MaintenanceJob struct {
	client.Client
	Instance types.NamespacedName
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

// NewMaintenanceJob creates a MaintenanceJob to schedule it
func NewMaintenanceJob(client client.Client, instance types.NamespacedName, log logr.Logger, scheme *runtime.Scheme) *MaintenanceJob {
	return &MaintenanceJob{
		Client:   client,
		Instance: instance,
		Log:      log,
		Scheme:   scheme,
	}
}

// Run implement the Job interface to use with Cron AddFunc()
func (b *MaintenanceJob) Run() {
	ctx := context.Background()
	instance := mysqlv1alpha1.Instance{}
	if err := b.Client.Get(ctx, b.Instance, &instance); err != nil {
		b.Log.Info(fmt.Sprintf("job for %s/%s failed. Could not access instance...", b.Instance.Namespace, b.Instance.Name))
		return
	}
	if instance.Status.MaintenanceMode == true && instance.Status.Schedules.MaintenanceEndTime.Time.After(time.Now()) {
		b.Log.Info(fmt.Sprintf("job for %s/%s already in maintenance, do not reschedule...", b.Instance.Namespace, b.Instance.Name))
		return
	}
	instance.Status.MaintenanceMode = true
	t := time.Now().Add(time.Duration(instance.Spec.MaintenanceSchedule.Duration) * time.Minute)
	instance.Status.Schedules.MaintenanceEndTime = &metav1.Time{Time: t}
	t = t.Add(time.Minute)
	nn := types.NamespacedName{
		Namespace: instance.ObjectMeta.Namespace,
		Name:      instance.ObjectMeta.Name,
	}
	cmd := NewUnMaintenanceJob(b.Client, nn, b.Log, b.Scheme)
	crontab.schedule(b.Log, &instance, MaintenanceUnscheduling, fmt.Sprintf("%s *", t.Format("4 15 2 1")), cmd)
	if err := b.Client.Status().Update(ctx, &instance); err != nil {
		b.Log.Info(fmt.Sprintf("Error updating Status.Maintenance, err: %v", err))
		return
	}
	b.Log.Info(fmt.Sprintf("Maintenance Mode for %s/%s enabled...", b.Instance.Namespace, b.Instance.Name))
}

// MaintenanceJob is a struct that manages Jobs for maintenance
type UnMaintenanceJob struct {
	client.Client
	Instance types.NamespacedName
	Log      logr.Logger
	Scheme   *runtime.Scheme
}

// NewUnMaintenanceJob creates a MaintenanceJob to schedule it
func NewUnMaintenanceJob(client client.Client, instance types.NamespacedName, log logr.Logger, scheme *runtime.Scheme) *UnMaintenanceJob {
	return &UnMaintenanceJob{
		Client:   client,
		Instance: instance,
		Log:      log,
		Scheme:   scheme,
	}
}

// Run implement the Job interface to use with Cron AddFunc()
func (b *UnMaintenanceJob) Run() {
	ctx := context.Background()
	instance := mysqlv1alpha1.Instance{}
	if err := b.Client.Get(ctx, b.Instance, &instance); err != nil {
		b.Log.Info(fmt.Sprintf("job for %s/%s failed. Could not access instance...", b.Instance.Namespace, b.Instance.Name))
		return
	}
	if instance.Status.MaintenanceMode == true && instance.Status.Schedules.MaintenanceEndTime.Time.After(time.Now()) {
		b.Log.Info(fmt.Sprintf("job for %s/%s race condition with maintenance, should be have been rescheduled", b.Instance.Namespace, b.Instance.Name))
		return
	}
	if instance.Status.MaintenanceMode == true {
		instance.Status.MaintenanceMode = false
		instance.Status.Schedules.MaintenanceEndTime = nil
		crontab.unSchedule(&instance, MaintenanceUnscheduling)
	}
	if err := b.Client.Status().Update(ctx, &instance); err != nil {
		b.Log.Info(fmt.Sprintf("Error updating Status.Maintenance, err: %v", err))
		return
	}
	b.Log.Info(fmt.Sprintf("Maintenance Mode for %s/%s disabled...", b.Instance.Namespace, b.Instance.Name))
}
