package service

import openapi "github.com/blaqkube/mysql-operator/agent/go"

// Router defines the required methods for binding the api requests to a responses for the MysqlApi
// The Router implementation should parse necessary information from the http request,
// pass the data to a Servicer to perform the required actions, then write the service results to the http response.
type Router interface {
	Routes() openapi.Routes
}
