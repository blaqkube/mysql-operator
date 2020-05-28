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
	"encoding/json"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gorilla/mux"
)

// A MysqlApiController binds http requests to an api service and writes the service results to the http response
type MysqlApiController struct {
	service MysqlApiServicer
}

// NewMysqlApiController creates a default api controller
func NewMysqlApiController(s MysqlApiServicer) openapi.Router {
	return &MysqlApiController{service: s}
}

// Routes returns all of the api route for the MysqlApiController
func (c *MysqlApiController) Routes() openapi.Routes {
	return openapi.Routes{
		{
			"CreateBackup",
			strings.ToUpper("Post"),
			"/backup",
			c.CreateBackup,
		},
		{
			"CreateDatabase",
			strings.ToUpper("Post"),
			"/database",
			c.CreateDatabase,
		},
		{
			"CreateUser",
			strings.ToUpper("Post"),
			"/user",
			c.CreateUser,
		},
		{
			"DeleteBackup",
			strings.ToUpper("Delete"),
			"/backup/{backup}",
			c.DeleteBackup,
		},
		{
			"DeleteDatabase",
			strings.ToUpper("Delete"),
			"/database/{database}",
			c.DeleteDatabase,
		},
		{
			"DeleteUser",
			strings.ToUpper("Delete"),
			"/user/{user}",
			c.DeleteUser,
		},
		{
			"GetBackupByName",
			strings.ToUpper("Get"),
			"/backup/{backup}",
			c.GetBackupByName,
		},
		{
			"GetDatabaseByName",
			strings.ToUpper("Get"),
			"/database/{database}",
			c.GetDatabaseByName,
		},
		{
			"GetDatabases",
			strings.ToUpper("Get"),
			"/database",
			c.GetDatabases,
		},
		{
			"GetUserByName",
			strings.ToUpper("Get"),
			"/user/{user}",
			c.GetUserByName,
		},
	}
}

// CreateBackup - create an on-demand backup
func (c *MysqlApiController) CreateBackup(w http.ResponseWriter, r *http.Request) {
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

// CreateDatabase - create an on-demand database
func (c *MysqlApiController) CreateDatabase(w http.ResponseWriter, r *http.Request) {
	body := &map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(500)
		return
	}

	apiKey := r.Header.Get("apiKey")
	result, err := c.service.CreateDatabase(*body, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}

// CreateUser - create an on-demand user
func (c *MysqlApiController) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &openapi.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(500)
		return
	}

	apiKey := r.Header.Get("apiKey")
	result, err := c.service.CreateUser(*user, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}

// DeleteBackup - Deletes a backup
func (c *MysqlApiController) DeleteBackup(w http.ResponseWriter, r *http.Request) {
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

// DeleteDatabase - Deletes a database
func (c *MysqlApiController) DeleteDatabase(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	database := params["database"]
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.DeleteDatabase(database, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}

// DeleteUser - Deletes a user
func (c *MysqlApiController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.DeleteUser(user, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}

// GetBackupByName - Get backup properties
func (c *MysqlApiController) GetBackupByName(w http.ResponseWriter, r *http.Request) {
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

// GetDatabaseByName - Get Database properties
func (c *MysqlApiController) GetDatabaseByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	database := params["database"]
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.GetDatabaseByName(database, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}

// GetDatabases - list all databases
func (c *MysqlApiController) GetDatabases(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.GetDatabases(apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, nil, w)
}

// GetUserByName - Get user properties
func (c *MysqlApiController) GetUserByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.GetUserByName(user, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	openapi.EncodeJSONResponse(result, nil, w)
}
