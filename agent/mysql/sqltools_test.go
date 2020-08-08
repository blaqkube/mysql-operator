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

type Suite struct {
	suite.Suite
	db          *sql.DB
	mock        sqlmock.Sqlmock
	testService DBTools
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.testService = NewDBTools(s.db)
}

func (s *Suite) Test_CheckDB() {
	err := s.testService.CheckDB(1)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_CreateExporter() {
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
	err := s.testService.CreateExporter()
	require.NoError(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
	suite.Run(t, &DumpSuite{})
}
