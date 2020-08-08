/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package database

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gorilla/mux"
)

// A MysqlDatabaseController binds http requests to an api service and writes the service results to the http response
type MysqlDatabaseController struct {
	service MysqlDatabaseServicer
}

// NewMysqlDatabaseController creates a default api controller
func NewMysqlDatabaseController(s MysqlDatabaseServicer) MysqlDatabaseRouter {
	return &MysqlDatabaseController{
		service: s,
	}
}

// Routes returns all of the api route for the MysqlDatabaseController
func (c *MysqlDatabaseController) Routes() openapi.Routes {
	routes := openapi.Routes{
		{
			"CreateDatabase",
			strings.ToUpper("Post"),
			"/database",
			c.CreateDatabase,
		},
		{
			"DeleteDatabase",
			strings.ToUpper("Delete"),
			"/database/{database}",
			c.DeleteDatabase,
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
	}
	return routes
}

// CreateDatabase - create an on-demand database
func (c *MysqlDatabaseController) CreateDatabase(w http.ResponseWriter, r *http.Request) {
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

// DeleteDatabase - Deletes a database
func (c *MysqlDatabaseController) DeleteDatabase(w http.ResponseWriter, r *http.Request) {
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

// GetDatabaseByName - Get Database properties
func (c *MysqlDatabaseController) GetDatabaseByName(w http.ResponseWriter, r *http.Request) {
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
func (c *MysqlDatabaseController) GetDatabases(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.GetDatabases(apiKey)
	if err != nil {
		fmt.Printf("Error: %v", err)
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, nil, w)
}