package controllers

import (
	mysqlv1alpha1 "github.com/blaqkube/mysql-operator/mysql-operator/api/v1alpha1"
	"github.com/robfig/cron/v3"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MockCrontab provides a simple struct to mock cron EntryID
type MockCrontab struct {
	Incarnation string
}

// NewMockCrontabCrontab initialize a predictable incarnation to ease tests
func NewMockCrontabCrontab() Crontab {
	crontab := &MockCrontab{
		Incarnation: "00000000-0000-0000-0000-000000000001",
	}
	return crontab
}

func (c *MockCrontab) isScheduled(instance mysqlv1alpha1.Instance, scheduleType string) bool {
	return false
}

func (c *MockCrontab) unSchedule(instance *mysqlv1alpha1.Instance, scheduleType string) bool {

	return false
}

func (c *MockCrontab) reScheduleAll(client client.Client, instance *mysqlv1alpha1.Instance, log logr.Logger, scheme *runtime.Scheme) bool {
	if instance.Status.Schedules.Incarnation == "00000000-0000-0000-0000-000000000001" {
		return false
	}
	instance.Status.Schedules.Incarnation = "00000000-0000-0000-0000-000000000001"
	return true
}

func (c *MockCrontab) schedule(log logr.Logger, instance *mysqlv1alpha1.Instance, scheduleType string, schedule string, cmd cron.Job) bool {
	return false
}
