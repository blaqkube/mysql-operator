package mysql

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InstanceSuite struct {
	suite.Suite
	db          *sql.DB
	mock        sqlmock.Sqlmock
	testService *Instance
}

func (s *InstanceSuite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.testService = NewInstance(s.db)
}

func (s *InstanceSuite) Test_Check() {
	err := s.testService.Check(1)
	require.NoError(s.T(), err)
}

func (s *InstanceSuite) Test_Initialize() {
	s.mock.ExpectExec(regexp.QuoteMeta(
		"create user if not exists 'exporter'@'localhost' identified by 'exporter' WITH MAX_USER_CONNECTIONS 3",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectExec(regexp.QuoteMeta(
		"create user if not exists 'exporter'@'::1' identified by 'exporter' WITH MAX_USER_CONNECTIONS 3",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectExec(regexp.QuoteMeta(
		"GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'localhost'",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectExec(regexp.QuoteMeta(
		"GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'::1'",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))
	err := s.testService.Initialize()
	require.NoError(s.T(), err)
}

func TestInstanceSuite(t *testing.T) {
	suite.Run(t, &InstanceSuite{})
}
