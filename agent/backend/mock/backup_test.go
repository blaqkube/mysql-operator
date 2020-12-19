package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BackupSuite struct {
	suite.Suite
	Service *Backup
}

func (s *BackupSuite) SetupTest() {
	s.Service = NewBackup()
}

func (s *BackupSuite) TestBackupSuccess() {

	err := s.Service.Run("key")
	assert.NoError(s.T(), err, "No Error")
}

func TestBackupSuite(t *testing.T) {
	suite.Run(t, &BackupSuite{})
}
