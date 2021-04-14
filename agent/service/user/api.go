package user

import (
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// MysqlUserRouter defines the required methods for binding the api requests to a responses for the MysqlApi
// The MysqlApiRouter implementation should parse necessary information from the http request,
// pass the data to a MysqlApiServicer to perform the required actions, then write the service results to the http response.
type MysqlUserRouter interface {
	Routes() openapi.Routes
	CreateUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	GetUserByName(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
}

// MysqlUserServicer defines the api actions for the MysqlApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type MysqlUserServicer interface {
	CreateUser(openapi.User, string) (interface{}, error)
	DeleteUser(string, string) (interface{}, error)
	GetUserByName(string, string) (interface{}, error)
	GetUsers(string) (interface{}, error)
}
