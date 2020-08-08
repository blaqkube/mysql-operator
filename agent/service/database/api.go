package database

import (
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// MysqlDatabaseRouter defines the required methods for binding the api requests to a responses for the MysqlDatabase
// The MysqlDatabaseRouter implementation should parse necessary information from the http request,
// pass the data to a MysqlDatabaseServicer to perform the required actions, then write the service results to the http response.
type MysqlDatabaseRouter interface {
	Routes() openapi.Routes
	CreateDatabase(http.ResponseWriter, *http.Request)
	DeleteDatabase(http.ResponseWriter, *http.Request)
	GetDatabaseByName(http.ResponseWriter, *http.Request)
	GetDatabases(http.ResponseWriter, *http.Request)
}

// MysqlDatabaseServicer defines the api actions for the MysqlDatabase service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type MysqlDatabaseServicer interface {
	CreateDatabase(map[string]interface{}, string) (interface{}, error)
	DeleteDatabase(string, string) (interface{}, error)
	GetDatabaseByName(string, string) (interface{}, error)
	GetDatabases(string) (interface{}, error)
}
