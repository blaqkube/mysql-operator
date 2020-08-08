package mysql

import (
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DumpSuite struct {
	suite.Suite
	testService S3MysqlBackup
}

func (s *DumpSuite) SetupSuite() {
	s.testService = NewS3MysqlBackup()
}

func (s *DumpSuite) Test_InitializeBackup() {
	backup := openapi.Backup{}
	b, err := s.testService.InitializeBackup(backup)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Pending", b.Status, "Status should be Started")
}
