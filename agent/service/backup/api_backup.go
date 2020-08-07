package backup

import (
	"encoding/json"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gorilla/mux"
)

// A MysqlApiController binds http requests to an api service and writes the service results to the http response
type MysqlBackupController struct {
	service MysqlBackupServicer
}

// NewMysqlApiController creates a default api controller
func NewMysqlBackupController(s MysqlBackupServicer) MysqlBackupRouter {
	return &MysqlBackupController{service: s}
}

// Routes returns all of the api route for the MysqlApiController
func (c *MysqlBackupController) Routes() openapi.Routes {
	return openapi.Routes{
		{
			"CreateBackup",
			strings.ToUpper("Post"),
			"/backup",
			c.CreateBackup,
		},
		{
			"DeleteBackup",
			strings.ToUpper("Delete"),
			"/backup/{backup}",
			c.DeleteBackup,
		},
		{
			"GetBackupByName",
			strings.ToUpper("Get"),
			"/backup/{backup}",
			c.GetBackupByName,
		},
	}
}

// CreateBackup - create an on-demand backup
func (c *MysqlBackupController) CreateBackup(w http.ResponseWriter, r *http.Request) {
	backup := &openapi.Backup{}
	if err := json.NewDecoder(r.Body).Decode(&backup); err != nil {
		w.WriteHeader(500)
		return
	}

	apiKey := r.Header.Get("apiKey")
	result, err := c.service.CreateBackup(*backup, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	code := http.StatusCreated
	openapi.EncodeJSONResponse(result, &code, w)
}

// DeleteBackup - Deletes a backup
func (c *MysqlBackupController) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	backup := params["backup"]
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.DeleteBackup(backup, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}

// GetBackupByName - Get backup properties
func (c *MysqlBackupController) GetBackupByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	backup := params["backup"]
	apiKey := r.Header.Get("apiKey")
	result, code, err := c.service.GetBackupByName(backup, apiKey)
	if err != nil && code != 0 {
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, &code, w)
}
