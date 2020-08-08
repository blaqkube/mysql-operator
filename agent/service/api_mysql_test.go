package service

import (
	"testing"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	c := NewMysqlApiController()

	next := openapi.NewRouter(c)
	p, err := next.GetRoute("CreateBackup").GetPathRegexp()
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, "^/backup[/]?$", p, "Should succeed")
	m, err := next.GetRoute("CreateBackup").GetMethods()
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, []string{"POST"}, m, "Should succeed")

	p, err = next.GetRoute("CreateDatabase").GetPathRegexp()
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, "^/database[/]?$", p, "Should succeed")
	m, err = next.GetRoute("CreateDatabase").GetMethods()
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, []string{"POST"}, m, "Should succeed")

	p, err = next.GetRoute("CreateUser").GetPathRegexp()
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, "^/user[/]?$", p, "Should succeed")
	m, err = next.GetRoute("CreateUser").GetMethods()
	assert.Equal(t, nil, err, "Should succeed")
	assert.Equal(t, []string{"POST"}, m, "Should succeed")

}
