package controllers

import (
	"fmt"
	"time"

	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// BackupScheduling defines the schedule associated with instance backups
	BackupScheduling string = "backup"
	// MaintenanceScheduling defines the schedule associated with maintenance except for backup
	MaintenanceScheduling string = "maintenance"
	// MaintenanceUnscheduling defines the schedule associated with maintenance removal
	MaintenanceUnscheduling string = "removal"
)

// Crontab provides a simple struct to manage cron EntryID for instances
type Crontab struct {
	Cron        *cron.Cron
	Incarnation string
}

func NewCrontab() Crontab {
	crontab := Crontab{
		Cron:        cron.New(),
		Incarnation: uuid.New().String(),
	}
	crontab.Cron.Start()
	return crontab
}

func (c *Crontab) isScheduled(instance mysqlv1alpha1.Instance, scheduleType string) bool {
	if c.Incarnation != instance.Status.Schedules.Incarnation {
		return false
	}
	schedule := mysqlv1alpha1.ScheduleEntry{}
	switch scheduleType {
	case BackupScheduling:
		schedule = instance.Status.Schedules.Backup
	case MaintenanceScheduling:
		schedule = instance.Status.Schedules.Maintenance
	case MaintenanceUnscheduling:
		schedule = instance.Status.Schedules.MaintenanceOff
	}
	entry := schedule.EntryID
	if entry == -1 {
		return false
	}
	return true
}

func (c *Crontab) unSchedule(instance *mysqlv1alpha1.Instance, scheduleType string) bool {

	if !c.isScheduled(*instance, scheduleType) {
		return false
	}

	switch scheduleType {
	case BackupScheduling:
		entry := instance.Status.Schedules.Backup.EntryID
		c.Cron.Remove(cron.EntryID(entry))
		instance.Status.Schedules.Backup = mysqlv1alpha1.ScheduleEntry{
			EntryID: -1,
		}
	case MaintenanceScheduling:
		entry := instance.Status.Schedules.Maintenance.EntryID
		c.Cron.Remove(cron.EntryID(entry))
		instance.Status.Schedules.Maintenance = mysqlv1alpha1.ScheduleEntry{
			EntryID: -1,
		}
	case MaintenanceUnscheduling:
		entry := instance.Status.Schedules.MaintenanceOff.EntryID
		c.Cron.Remove(cron.EntryID(entry))
		instance.Status.Schedules.MaintenanceOff = mysqlv1alpha1.ScheduleEntry{
			EntryID: -1,
		}
	}
	return true
}

func (c *Crontab) reScheduleAll(client client.Client, instance *mysqlv1alpha1.Instance, log logr.Logger, scheme *runtime.Scheme) bool {
	changed := false
	restarted := false
	sc := []string{
		BackupScheduling,
		MaintenanceScheduling,
		MaintenanceUnscheduling,
	}
	if instance.Status.Schedules.Incarnation != c.Incarnation {
		instance.Status.Schedules.Incarnation = c.Incarnation
		restarted = true
		changed = true
	}
	nn := types.NamespacedName{
		Namespace: instance.ObjectMeta.Namespace,
		Name:      instance.ObjectMeta.Name,
	}
	for _, v := range sc {
		switch v {
		case BackupScheduling:
			if instance.Spec.BackupSchedule.Schedule != "" {
				if restarted {
					cmd := NewBackupJob(client, nn, log, scheme)
					c.schedule(log, instance, v, instance.Spec.BackupSchedule.Schedule, cmd)
					changed = true
				}
				if !restarted && instance.Spec.BackupSchedule.Schedule != instance.Status.Schedules.Backup.Schedule {
					c.unSchedule(instance, BackupScheduling)
					cmd := NewBackupJob(client, nn, log, scheme)
					c.schedule(log, instance, v, instance.Spec.BackupSchedule.Schedule, cmd)
					changed = true
				}
			}
		case MaintenanceScheduling:
			if instance.Spec.MaintenanceSchedule.Schedule != "" {
				if restarted {
					cmd := NewMaintenanceJob(client, nn, log, scheme, c)
					c.schedule(log, instance, v, instance.Spec.BackupSchedule.Schedule, cmd)
					changed = true
				}
				if !restarted && instance.Spec.BackupSchedule.Schedule != instance.Status.Schedules.Backup.Schedule {
					c.unSchedule(instance, MaintenanceScheduling)
					cmd := NewMaintenanceJob(client, nn, log, scheme, c)
					c.schedule(log, instance, v, instance.Spec.BackupSchedule.Schedule, cmd)
					changed = true
				}
			}
		case MaintenanceUnscheduling:
			if instance.Status.MaintenanceMode == true {
				if restarted {
					if instance.Status.Schedules.MaintenanceEndTime == nil || instance.Status.Schedules.MaintenanceEndTime.Time.Before(time.Now()) {
						instance.Status.MaintenanceMode = false
						instance.Status.Schedules.MaintenanceEndTime = nil
						instance.Status.Schedules.MaintenanceOff = mysqlv1alpha1.ScheduleEntry{
							EntryID: -1,
						}
					} else {
						cmd := NewUnMaintenanceJob(client, nn, log, scheme, c)
						c.schedule(
							log,
							instance,
							v,
							fmt.Sprintf(
								"%s *",
								instance.Status.Schedules.MaintenanceEndTime.Time.Add(time.Minute).Format("4 15 2 1"),
							),
							cmd,
						)
					}
					changed = true
				}
			}
		}
	}
	return changed
}

func (c *Crontab) schedule(log logr.Logger, instance *mysqlv1alpha1.Instance, scheduleType string, schedule string, cmd cron.Job) bool {
	switch scheduleType {
	case BackupScheduling:
		entry := instance.Status.Schedules.Backup.EntryID
		if entry != -1 && instance.Status.Schedules.Backup.Schedule != schedule {
			c.Cron.Remove(cron.EntryID(entry))
			instance.Status.Schedules.Backup = mysqlv1alpha1.ScheduleEntry{
				EntryID: -1,
			}
			entry = -1
		}
		if entry == -1 {
			eid, err := c.Cron.AddJob(schedule, cmd)
			if err != nil {
				log.Info(
					fmt.Sprintf("Error scheduling backup job: %s", err.Error()),
				)
			}
			instance.Status.Schedules.Backup = mysqlv1alpha1.ScheduleEntry{
				EntryID:  int(eid),
				Schedule: schedule,
			}
			return true
		}
	case MaintenanceScheduling:
		entry := instance.Status.Schedules.Backup.EntryID
		if entry != -1 && instance.Status.Schedules.Maintenance.Schedule != schedule {
			c.Cron.Remove(cron.EntryID(entry))
			instance.Status.Schedules.Maintenance = mysqlv1alpha1.ScheduleEntry{
				EntryID: -1,
			}
			entry = -1
		}
		if entry == -1 {
			eid, err := c.Cron.AddJob(schedule, cmd)
			if err != nil {
				log.Info(
					fmt.Sprintf("Error scheduling maintenance job: %s", err.Error()),
				)
			}
			instance.Status.Schedules.Maintenance = mysqlv1alpha1.ScheduleEntry{
				EntryID:  int(eid),
				Schedule: schedule,
			}
			return true
		}
	case MaintenanceUnscheduling:
		entry := instance.Status.Schedules.Backup.EntryID
		if entry != -1 && instance.Status.Schedules.MaintenanceOff.Schedule != schedule {
			c.Cron.Remove(cron.EntryID(entry))
			instance.Status.Schedules.MaintenanceOff = mysqlv1alpha1.ScheduleEntry{
				EntryID: -1,
			}
			entry = -1
		}
		if entry == -1 {
			eid, err := c.Cron.AddJob(schedule, cmd)
			if err != nil {
				log.Info(
					fmt.Sprintf("Error scheduling unmaintenance job: %s", err.Error()),
				)
			}
			instance.Status.Schedules.MaintenanceOff = mysqlv1alpha1.ScheduleEntry{
				EntryID:  int(eid),
				Schedule: schedule,
			}
			return true
		}
	}
	return false
}
