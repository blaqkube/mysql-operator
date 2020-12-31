package cmd

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blaqkube/mysql-operator/agent/backend"
	"github.com/blaqkube/mysql-operator/agent/backend/mock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateExporterUser(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)
	storages := map[string]backend.Storage{"s3": mock.NewStorage()}
	backup := mock.NewBackup()
	instance := mock.NewInstance()
	resources = &Backend{
		Backup:   backup,
		DB:       db,
		Instance: instance,
		Storages: storages,
	}
	sqlMock.ExpectExec(
		`CREATE USER 'exporter'@'%' IDENTIFIED BY 'exporter' WITH MAX_USER_CONNECTIONS 3`,
	).WillReturnResult(sqlmock.NewResult(0, 0))
	sqlMock.ExpectExec(regexp.QuoteMeta(
		`GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'%'`,
	)).WillReturnResult(sqlmock.NewResult(0, 0))
	sqlMock.ExpectCommit()

	err = createExporterUser("exporter", "exporter")
	assert.NoError(t, err)

}
