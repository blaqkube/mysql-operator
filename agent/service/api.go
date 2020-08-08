/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package service

import (
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// MysqlApiRouter defines the required methods for binding the api requests to a responses for the MysqlApi
// The MysqlApiRouter implementation should parse necessary information from the http request,
// pass the data to a MysqlApiServicer to perform the required actions, then write the service results to the http response.
type MysqlApiRouter interface {
	CreateDatabase(http.ResponseWriter, *http.Request)
	CreateUser(http.ResponseWriter, *http.Request)
	DeleteDatabase(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	GetDatabaseByName(http.ResponseWriter, *http.Request)
	GetDatabases(http.ResponseWriter, *http.Request)
	GetUserByName(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
}

// MysqlApiServicer defines the api actions for the MysqlApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type MysqlApiServicer interface {
	CreateDatabase(map[string]interface{}, string) (interface{}, error)
	CreateUser(openapi.User, string) (interface{}, error)
	DeleteDatabase(string, string) (interface{}, error)
	DeleteUser(string, string) (interface{}, error)
	GetDatabaseByName(string, string) (interface{}, error)
	GetDatabases(string) (interface{}, error)
	GetUserByName(string, string) (interface{}, error)
	GetUsers(string) (interface{}, error)
}
