package grant

import (
	"errors"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockService struct{}

func (s *mockService) CreateGrantByUserDatabase(o openapi.Grant, user, database, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return o, nil
	}
	return nil, errors.New("user failed")
}

func (s *mockService) GetGrantByUserDatabase(user, database, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.Grant{
			AccessMode: "readWrite",
		}, nil
	}
	return nil, errors.New("failed")
}
