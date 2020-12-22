package backup

import (
	"testing"

	"github.com/blaqkube/mysql-operator/agent/backend"
	"github.com/blaqkube/mysql-operator/agent/backend/mock"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BackupServiceSuite struct {
	suite.Suite
	Service *Service
}

func (s *BackupServiceSuite) SetupSuite() {
	backup := mock.NewBackup()
	storages := map[string]backend.Storage{
		"s3":        mock.NewStorage(),
		"blackhole": mock.NewStorage(),
	}
	s.Service = NewService(backup, storages)
}

func (s *BackupServiceSuite) Test_GetBackups() {
	_, _, err := s.Service.GetBackups("apikey")
	require.NoError(s.T(), err)
}

func (s *BackupServiceSuite) Test_CreateBackup() {
	_, err := s.Service.CreateBackup(
		openapi.BackupRequest{Backend: "s3", Bucket: "bucket", Location: "file"},
		"apikey",
	)
	require.NoError(s.T(), err)
}

func TestBackupSuite(t *testing.T) {
	suite.Run(t, &BackupServiceSuite{})
}
