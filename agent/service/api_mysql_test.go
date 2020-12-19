package service

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	bmock "github.com/blaqkube/mysql-operator/agent/backend/mock"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	db          *sql.DB
	mock        sqlmock.Sqlmock
	testService Router
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	storage := bmock.NewStorage()
	backup := bmock.NewBackup()
	require.NoError(s.T(), err)

	s.testService = NewMysqlAPIController(s.db, backup, storage)
}

func (s *Suite) Test_Routes() {
	next := openapi.NewRouter(s.testService)
	p, err := next.GetRoute("CreateBackup").GetPathRegexp()
	assert.Equal(s.T(), nil, err, "Should succeed")
	assert.Equal(s.T(), "^/backup[/]?$", p, "Should succeed")
	m, err := next.GetRoute("CreateBackup").GetMethods()
	assert.Equal(s.T(), nil, err, "Should succeed")
	assert.Equal(s.T(), []string{"POST"}, m, "Should succeed")

	p, err = next.GetRoute("CreateDatabase").GetPathRegexp()
	assert.Equal(s.T(), nil, err, "Should succeed")
	assert.Equal(s.T(), "^/database[/]?$", p, "Should succeed")
	m, err = next.GetRoute("CreateDatabase").GetMethods()
	assert.Equal(s.T(), nil, err, "Should succeed")
	assert.Equal(s.T(), []string{"POST"}, m, "Should succeed")

	p, err = next.GetRoute("CreateUser").GetPathRegexp()
	assert.Equal(s.T(), nil, err, "Should succeed")
	assert.Equal(s.T(), "^/user[/]?$", p, "Should succeed")
	m, err = next.GetRoute("CreateUser").GetMethods()
	assert.Equal(s.T(), nil, err, "Should succeed")
	assert.Equal(s.T(), []string{"POST"}, m, "Should succeed")

}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
