package user

import (
	"errors"
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockService struct{}

func (s *mockService) CreateUser(o openapi.User, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return o, nil
	}
	return nil, errors.New("user failed")
}

func (s *mockService) DeleteUser(user, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.Message{
			Code:    int32(http.StatusOK),
			Message: "user deleted",
		}, nil
	}
	return nil, errors.New("failed")
}

func (s *mockService) GetUserByName(user, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.User{
			Username: user,
			Password: "****",
		}, nil
	}
	return nil, errors.New("failed")
}

func (s *mockService) GetUsers(apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.ListUsers{
			Size: int32(1),
			Items: []openapi.User{{
				Username: "me",
				Password: "****",
			}},
		}, nil
	}
	return nil, errors.New("failed")
}
