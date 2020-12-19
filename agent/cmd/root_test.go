package cmd

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/blaqkube/mysql-operator/agent/backend/mock"
)

func Test_ExecuteCommand(t *testing.T) {
	db, _, _ := sqlmock.New()
	storage := mock.NewStorage()
	backup := mock.NewBackup()
	instance := mock.NewInstance()
	Execute(backup, db, instance, storage)
}
