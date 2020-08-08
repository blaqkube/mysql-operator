package database

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	openapi "github.com/blaqkube/mysql-operator/agent/go"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	db          *sql.DB
	mock        sqlmock.Sqlmock
	testService MysqlDatabaseServicer
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.testService = NewMysqlDatabaseService(s.db)
}

func (s *Suite) Test_CreateDatabase() {
	name := "me"
	s.mock.ExpectExec(regexp.QuoteMeta(
		"create database me")).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := s.testService.CreateDatabase(openapi.Database{Name: name}, "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetMissingDatabaseByName() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT schema_name FROM information_schema.schemata where schema_name=?")).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)

	_, err := s.testService.GetDatabaseByName("me", "test1")
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetDatabaseByName() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT schema_name FROM information_schema.schemata where schema_name=?")).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"schema_name"}).
			AddRow("me"))

	db, err := s.testService.GetDatabaseByName("me", "test1")
	require.NoError(s.T(), err)
	switch t := db.(type) {
	case *openapi.Database:
		require.Equal(s.T(), "me", t.Name)
	default:
		require.Equal(s.T(), "type", fmt.Sprintf("%T", db))
	}
}

func (s *Suite) Test_GetAllDatabases() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT schema_name FROM information_schema.schemata")).
		WillReturnRows(sqlmock.NewRows([]string{"database"}).
			AddRow("me"))

	db, err := s.testService.GetDatabases("test1")
	require.NoError(s.T(), err)
	switch t := db.(type) {
	case *openapi.ListDatabases:
		require.Equal(s.T(), "me", t.Items[0].Name)
	default:
		require.Equal(s.T(), "type", fmt.Sprintf("%T", db))
	}
}

func (s *Suite) Test_GetDatabasesWithError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT schema_name FROM information_schema.schemata")).
		WillReturnError(errors.New("query failed"))

	_, err := s.testService.GetDatabases("test1")
	require.Error(s.T(), err)
}

func (s *Suite) Test_DeleteDatabase() {
	_, err := s.testService.DeleteDatabase("me", "test1")
	require.Error(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
