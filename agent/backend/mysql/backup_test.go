package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BackupSuite struct {
	suite.Suite
	backupService *Backup
}

func (s *BackupSuite) SetupSuite() {
	s.backupService = &Backup{}
}

func (s *BackupSuite) TestBackup() {
	s.backupService.Exec = "true"
	err := s.backupService.Run("backup.dmp")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) TestFailedBackup() {
	s.backupService.Exec = "false"
	err := s.backupService.Run("backup.dmp")
	require.Error(s.T(), err)
	require.Equal(s.T(), "exit status 1", err.Error())
}

func TestBackupSuite(t *testing.T) {
	suite.Run(t, &BackupSuite{})
}
