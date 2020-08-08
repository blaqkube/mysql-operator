package backup

import (
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// MysqlBackupRouter defines the required methods for binding the api requests to a responses for the MysqlBackup
// The MysqlBackupRouter implementation should parse necessary information from the http request,
// pass the data to a MysqlBackupServicer to perform the required actions, then write the service results to the http response.
type MysqlBackupRouter interface {
	Routes() openapi.Routes
	CreateBackup(http.ResponseWriter, *http.Request)
	DeleteBackup(http.ResponseWriter, *http.Request)
	GetBackupByName(http.ResponseWriter, *http.Request)
}

// MysqlApiServicer defines the api actions for the MysqlApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type MysqlBackupServicer interface {
	CreateBackup(openapi.Backup, string) (interface{}, error)
	DeleteBackup(string, string) (interface{}, error)
	GetBackupByName(string, string) (interface{}, int, error)
}
