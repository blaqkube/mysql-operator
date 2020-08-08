package database

import (
	"errors"
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockService struct{}

func (s *mockService) CreateDatabase(o openapi.Database, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return o, nil
	}
	return nil, errors.New("database failed")
}

func (s *mockService) DeleteDatabase(database, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.Message{
			Code:    int32(http.StatusOK),
			Message: "database deleted",
		}, nil
	}
	return nil, errors.New("failed")
}

func (s *mockService) GetDatabaseByName(database, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.Database{
			Name: database,
		}, nil
	}
	return nil, errors.New("failed")
}

func (s *mockService) GetDatabases(apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.ListDatabases{
			Size: int32(1),
			Items: []openapi.Database{{
				Name: "me",
			}},
		}, nil
	}
	return &openapi.Message{Code: int32(http.StatusInternalServerError), Message: "Internal Server Error"}, errors.New("failed")
}
