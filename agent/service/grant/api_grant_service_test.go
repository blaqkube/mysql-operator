package grant

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
	testService MysqlGrantServicer
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.testService = NewMysqlGrantService(s.db)
}

func (s *Suite) Test_CreateReadWriteGrant() {
	s.mock.ExpectExec(regexp.QuoteMeta(
		"GRANT ALL PRIVILEGES ON pong.* TO 'me'@'%'",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := s.testService.CreateGrantByUserDatabase(openapi.Grant{AccessMode: "readWrite"}, "me", "pong", "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_CreateErrorGrant() {
	s.mock.ExpectExec(regexp.QuoteMeta(
		"GRANT ALL PRIVILEGES ON pong.* TO 'me'@'%'",
	)).
		WithArgs().
		WillReturnError(errors.New("BaBoom"))

	_, err := s.testService.CreateGrantByUserDatabase(openapi.Grant{AccessMode: "readWrite"}, "me", "pong", "test1")
	require.Error(s.T(), err)

}

func (s *Suite) Test_CreateReadOnlyGrant() {
	s.mock.ExpectExec(regexp.QuoteMeta(
		"GRANT SELECT ON pong.* TO 'me'@'%'",
	)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := s.testService.CreateGrantByUserDatabase(openapi.Grant{AccessMode: "readOnly"}, "me", "pong", "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_GetNoneGrant() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SHOW GRANTS FOR 'me'@'%'")).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"GRANTS"}).
			AddRow("GRANT USAGE ON *.* TO `me`@`%`"))

	_, err := s.testService.GetGrantByUserDatabase("me", "pong", "test1")
	require.NoError(s.T(), err)

}

func (s *Suite) Test_GetReadOnlyGrant() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SHOW GRANTS FOR 'me'@'%'")).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"GRANTS"}).
			AddRow("GRANT USAGE ON *.* TO `me`@`%`").
			AddRow("GRANT SELECT ON `pong`.* TO `me`@`%`"))

	_, err := s.testService.GetGrantByUserDatabase("me", "pong", "test1")
	require.NoError(s.T(), err)

}

func (s *Suite) Test_GetReadWriteGrant() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SHOW GRANTS FOR 'me'@'%'")).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"GRANTS FOR greg@%"}).
			AddRow("GRANT USAGE ON *.* TO `me`@`%`").
			AddRow("GRANT ALL PRIVILEGES ON `pong`.* TO `me`@`%`"))
	_, err := s.testService.GetGrantByUserDatabase("me", "pong", "test1")
	require.NoError(s.T(), err)

}

func (s *Suite) Test_GetErrorOnGrant() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SHOW GRANTS FOR 'me'@'%'")).
		WithArgs().
		WillReturnError(errors.New("BaBoom"))

	_, err := s.testService.GetGrantByUserDatabase("me", "pong", "test1")
	require.Error(s.T(), err)

}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
