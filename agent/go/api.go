/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore 
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"
)


// MysqlApiRouter defines the required methods for binding the api requests to a responses for the MysqlApi
// The MysqlApiRouter implementation should parse necessary information from the http request, 
// pass the data to a MysqlApiServicer to perform the required actions, then write the service results to the http response.
type MysqlApiRouter interface { 
	CreateBackup(http.ResponseWriter, *http.Request)
	DeleteBackup(http.ResponseWriter, *http.Request)
	GetBackupByName(http.ResponseWriter, *http.Request)
}


// MysqlApiServicer defines the api actions for the MysqlApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type MysqlApiServicer interface { 
	CreateBackup(Backup, string) (interface{}, error)
	DeleteBackup(string, string) (interface{}, error)
	GetBackupByName(string, string) (interface{}, error)
}