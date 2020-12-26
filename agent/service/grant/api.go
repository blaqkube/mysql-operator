package grant

import (
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

const (
	// NoneAccessMode the user does not have access to the database
	NoneAccessMode = "none"
	// ReadWriteAccessMode the user does not have readWrite access to the database
	ReadWriteAccessMode = "readWrite"
	// ReadOnlyAccessMode the user does not have readOnly access to the database
	ReadOnlyAccessMode  = "readOnly"
)

// MysqlGrantRouter defines the required methods for binding the api requests to a responses for the MysqlApi
// The MysqlGrantRouter implementation should parse necessary information from the http request,
// pass the data to a MysqlGrantServicer to perform the required actions, then write the service results to the http response.
type MysqlGrantRouter interface {
	Routes() openapi.Routes
	CreateGrantByUserDatabase(http.ResponseWriter, *http.Request)
	GetGrantByUserDatabase(http.ResponseWriter, *http.Request)
}

// MysqlGrantServicer defines the api actions for the MysqlApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type MysqlGrantServicer interface {
	CreateGrantByUserDatabase(openapi.Grant, string, string, string) (interface{}, error)
	GetGrantByUserDatabase(string, string, string) (interface{}, error)
}
