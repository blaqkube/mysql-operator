package user

import (
	"database/sql"
	"errors"
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
	testService MysqlUserServicer
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.testService = NewMysqlUserService(s.db)
}

func (s *Suite) Test_CreateExistingUser() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT user FROM mysql.user where user=?")).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"user"}).
			AddRow("me"))

	_, err := s.testService.CreateUser(openapi.User{Username: "me"}, "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_CreateMissingUser() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT user FROM mysql.user where user=?")).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)
	s.mock.ExpectExec(regexp.QuoteMeta(
		"create user 'me'@'%' identified by 'me'",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := s.testService.CreateUser(openapi.User{Username: "me", Password: "me"}, "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_CreateUserWithError1() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT user FROM mysql.user where user=?")).
		WithArgs(name).
		WillReturnError(errors.New("error"))

	_, err := s.testService.CreateUser(openapi.User{Username: "me", Password: "me"}, "test1")
	require.Error(s.T(), err)
}

func (s *Suite) Test_CreateUserWithError2() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT user FROM mysql.user where user=?")).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)
	s.mock.ExpectExec(regexp.QuoteMeta(
		"create user 'me'@'%' identified by 'me'",
	)).
		WithArgs().
		WillReturnError(errors.New("error"))

	_, err := s.testService.CreateUser(openapi.User{Username: "me", Password: "me"}, "test1")
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetMissingUserByName() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT User FROM mysql.user where User=?")).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)

	_, err := s.testService.GetUserByName("me", "test1")
	require.Error(s.T(), err)
}

func (s *Suite) Test_GetExistingUserByName() {
	name := "me"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT User FROM mysql.user where User=?")).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"user"}).
			AddRow("me"))

	_, err := s.testService.GetUserByName("me", "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetAllUsers() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT User FROM mysql.user")).
		WillReturnRows(sqlmock.NewRows([]string{"user"}).
			AddRow("me"))

	_, err := s.testService.GetUsers("test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetUsersWithError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT User FROM mysql.user")).
		WillReturnError(errors.New("query failed"))

	_, err := s.testService.GetUsers("test1")
	require.Error(s.T(), err)
}

func (s *Suite) Test_DeleteUser() {
	_, err := s.testService.DeleteUser("me", "test1")
	require.Error(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
