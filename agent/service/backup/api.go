package backup

import (
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// Router defines the required methods for binding the api requests to a responses for the MysqlBackup
// The Router implementation should parse necessary information from the http request,
// pass the data to a Servicer to perform the required actions, then write the service results to the http response.
type Router interface {
	Routes() openapi.Routes
	CreateBackup(http.ResponseWriter, *http.Request)
	GetBackupByID(http.ResponseWriter, *http.Request)
	GetBackups(http.ResponseWriter, *http.Request)
}

// Servicer defines the api actions for the Backup service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type Servicer interface {
	CreateBackup(openapi.BackupRequest, string) (interface{}, int, error)
	GetBackupByID(string, string) (interface{}, int, error)
	GetBackups(string) (interface{}, int, error)
}
