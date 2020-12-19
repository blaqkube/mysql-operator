package backup

import (
	"encoding/json"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
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
			Name: "CreateBackup",
			Method: strings.ToUpper("Post"),
			Pattern: "/backup",
			HandlerFunc: c.CreateBackup,
		},
		{
			Name: "GetBackups",
			Method: strings.ToUpper("Get"),
			Pattern: "/backup",
			HandlerFunc: c.GetBackups,
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
	result, err := c.service.CreateBackup(*request, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	code := http.StatusCreated
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
