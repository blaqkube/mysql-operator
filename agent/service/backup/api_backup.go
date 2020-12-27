package backup

import (
	"encoding/json"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gorilla/mux"
)

// Controller binds http requests to an api service and writes the service results to the http response
type Controller struct {
	service Servicer
}

// NewController creates a default api controller
func NewController(s Servicer) Router {
	return &Controller{service: s}
}

// Routes returns all of the api route for the ApiController
func (c *Controller) Routes() openapi.Routes {
	return openapi.Routes{
		{
			Name:        "CreateBackup",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/backup",
			HandlerFunc: c.CreateBackup,
		},
		{
			Name:        "GetBackups",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/backup",
			HandlerFunc: c.GetBackups,
		},
		{
			Name:        "GetBackupByID",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/backup/{uuid}",
			HandlerFunc: c.GetBackupByID,
		},
	}
}

// CreateBackup - create an on-demand backup
func (c *Controller) CreateBackup(w http.ResponseWriter, r *http.Request) {
	request := &openapi.BackupRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(500)
		return
	}

	apiKey := r.Header.Get("apiKey")
	result, code, err := c.service.CreateBackup(*request, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, &code, w)
}

// GetBackupByID - Get backup from UUID
func (c *Controller) GetBackupByID(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("apiKey")
	params := mux.Vars(r)
	uuid := params["uuid"]
	result, code, err := c.service.GetBackupByID(uuid, apiKey)
	if err != nil && code != 0 {
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, &code, w)
}

// GetBackups - Get backups
func (c *Controller) GetBackups(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("apiKey")
	result, code, err := c.service.GetBackups(apiKey)
	if err != nil && code != 0 {
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, &code, w)
}
