package cmd

import (
	"testing"

	"github.com/blaqkube/mysql-operator/agent/backend/mock"
)

func Test_ExecuteCommand(t *testing.T) {
	storage := mock.NewStorage()
	backup := mock.NewBackup()
	instance := mock.NewInstance()
	Execute(storage, backup, instance)
}
