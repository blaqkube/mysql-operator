package user

import (
	"encoding/json"
	"net/http"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gorilla/mux"
)

// A MysqlUserController binds http requests to an api service and writes the service results to the http response
type MysqlUserController struct {
	service MysqlUserServicer
}

// NewMysqlUserController creates a default api controller
func NewMysqlUserController(s MysqlUserServicer) MysqlUserRouter {
	return &MysqlUserController{
		service: s,
	}
}

// Routes returns all of the api route for the MysqlUserController
func (c *MysqlUserController) Routes() openapi.Routes {
	routes := openapi.Routes{
		{
			Name:        "CreateUser",
			Method:      strings.ToUpper("Post"),
			Pattern:     "/user",
			HandlerFunc: c.CreateUser,
		},
		{
			Name:        "DeleteUser",
			Method:      strings.ToUpper("Delete"),
			Pattern:     "/user/{user}",
			HandlerFunc: c.DeleteUser,
		},
		{
			Name:        "GetUserByName",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/user/{user}",
			HandlerFunc: c.GetUserByName,
		},
		{
			Name:        "GetUsers",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/user",
			HandlerFunc: c.GetUsers,
		},
	}
	return routes
}

// CreateUser - create an on-demand user
func (c *MysqlUserController) CreateUser(w http.ResponseWriter, r *http.Request) {
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
	statusCode := http.StatusCreated
	openapi.EncodeJSONResponse(result, &statusCode, w)
}

// DeleteUser - Deletes a user
func (c *MysqlUserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

// GetUserByName - Get user properties
func (c *MysqlUserController) GetUserByName(w http.ResponseWriter, r *http.Request) {
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

// GetUsers - list all users
func (c *MysqlUserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("apiKey")
	result, err := c.service.GetUsers(apiKey)
	if err != nil {
		statusCode := int(http.StatusInternalServerError)
		openapi.EncodeJSONResponse(result, &statusCode, w)
		return
	}
	openapi.EncodeJSONResponse(result, nil, w)
}
