package backup

import (
	"testing"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BackupSuite struct {
	suite.Suite
	testService MysqlBackupServicer
}

func (s *BackupSuite) SetupSuite() {
	s.testService = NewMysqlBackupService(&mockBackupPrimitive{})
}

func (s *BackupSuite) Test_Delete() {
	_, err := s.testService.DeleteBackup("1", "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_GetByName() {
	_, _, err := s.testService.GetBackupByName("2006-01-02T15:04:05Z", "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_GetByNameWithZeroTime() {
	_, _, err := s.testService.GetBackupByName("0001-01-01T00:00:00Z", "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_GetByNameWithError() {
	_, _, err := s.testService.GetBackupByName("123", "2")
	require.Error(s.T(), err)
}

func (s *BackupSuite) Test_CreateBackup() {
	_, err := s.testService.CreateBackup(openapi.Backup{}, "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_CreateBackupWithError() {
	_, err := s.testService.CreateBackup(openapi.Backup{Location: "/123"}, "2")
	require.Error(s.T(), err)
}

func (s *BackupSuite) Test_PullS3() {
	m := &mockBackupPrimitive{}
	err := m.PullS3File(&openapi.Backup{}, "/", "myfile.dmp")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_PushS3() {
	m := &mockBackupPrimitive{}
	err := m.PushS3File(&openapi.Backup{}, "myfile.dmp")
	require.NoError(s.T(), err)
}

func TestBackupSuite(t *testing.T) {
	suite.Run(t, &BackupSuite{})
}
