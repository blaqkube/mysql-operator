package grant

import (
	"encoding/json"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gorilla/mux"
)

// A MysqlGrantController binds http requests to an api service and writes the service results to the http response
type MysqlGrantController struct {
	service MysqlGrantServicer
}

// NewMysqlGrantController creates a default api controller
func NewMysqlGrantController(s MysqlGrantServicer) MysqlGrantRouter {
	return &MysqlGrantController{
		service: s,
	}
}

// Routes returns all of the api route for the MysqlUserController
func (c *MysqlGrantController) Routes() openapi.Routes {
	routes := openapi.Routes{
		{
			Name: "CreateGrantByUserDatabase",
			Method: strings.ToUpper("Post"),
			Pattern: "/user/{user}/database/{database}/grant",
			HandlerFunc: c.CreateGrantByUserDatabase,
		},
		{
			Name: "GetGrantByUserDatabase",
			Method: strings.ToUpper("Get"),
			Pattern: "/user/{user}/database/{database}/grant",
			HandlerFunc: c.GetGrantByUserDatabase,
		},
	}
	return routes
}

// CreateGrantByUserDatabase - create an on-demand grant
func (c *MysqlGrantController) CreateGrantByUserDatabase(w http.ResponseWriter, r *http.Request) {
	grant := &openapi.Grant{}
	if err := json.NewDecoder(r.Body).Decode(&grant); err != nil {
		w.WriteHeader(500)
		return
	}

	apiKey := r.Header.Get("apiKey")
	params := mux.Vars(r)
	user := params["user"]
	database := params["database"]
	result, err := c.service.CreateGrantByUserDatabase(*grant, user, database, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	statusCode := http.StatusCreated
	openapi.EncodeJSONResponse(result, &statusCode, w)
}

// GetGrantByUserDatabase - Get grant properties
func (c *MysqlGrantController) GetGrantByUserDatabase(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	database := params["database"]
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.GetGrantByUserDatabase(user, database, apiKey)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	openapi.EncodeJSONResponse(result, nil, w)
}
