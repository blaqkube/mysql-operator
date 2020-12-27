package backup

import (
	"fmt"
	"net/http"
	"testing"
	"time"

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
	s.Service.M.Lock()
	defer s.Service.M.Unlock()
	key := "abcd"
	s.Service.States[key] = openapi.Backup{
		Bucket:     "bucket",
		Location:   "location",
		Identifier: "abcd",
		StartTime:  time.Now(),
		Status:     "Succeeded",
	}
	s.Service.CurrState = key
}

func (s *BackupServiceSuite) Test_GetBackupByIDSucceed() {
	b, code, err := s.Service.GetBackupByID("abcd", "apikey")
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, http.StatusOK)
	switch v := b.(type) {
	case *openapi.Backup:
		require.Equal(s.T(), v.Bucket, "bucket")
	default:
		require.Equal(s.T(), fmt.Sprintf("%T", b), "unknown type")
	}
}

func (s *BackupServiceSuite) Test_GetBackupByIDFailed() {
	b, code, err := s.Service.GetBackupByID("abce", "apikey")
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, http.StatusNotFound)
	switch v := b.(type) {
	case *openapi.Backup:
		require.Equal(s.T(), v.Bucket, "")
	case nil:
		require.Nil(s.T(), b)
	default:
		require.Equal(s.T(), fmt.Sprintf("%T", b), "unknown type")
	}
}

func (s *BackupServiceSuite) Test_GetBackupsSucceed() {
	b, code, err := s.Service.GetBackups("apikey")
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, http.StatusOK)
	switch v := b.(type) {
	case *openapi.BackupList:
		require.Equal(s.T(), v.Size, int32(2))
		require.Equal(s.T(), v.Items[0].Bucket, "bucket")
		require.Equal(s.T(), v.Items[1].Bucket, "bucket")
	default:
		require.Equal(s.T(), fmt.Sprintf("%T", b), "unknown type")
	}
}

func (s *BackupServiceSuite) Test_CreateBackup() {
	b, code, err := s.Service.CreateBackup(
		openapi.BackupRequest{Backend: "s3", Bucket: "bucket", Location: "file"},
		"apikey",
	)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, code)
	switch v := b.(type) {
	case *openapi.Backup:
		require.Equal(s.T(), v.Bucket, "bucket")
	default:
		require.Equal(s.T(), fmt.Sprintf("%T", b), "unknown type")
	}
}

func TestBackupSuite(t *testing.T) {
	suite.Run(t, &BackupServiceSuite{})
}
