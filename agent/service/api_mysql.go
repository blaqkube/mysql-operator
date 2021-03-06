package service

import (
	"database/sql"

	// "github.com/blaqkube/mysql-operator/agent/backend/mysql"
	"github.com/blaqkube/mysql-operator/agent/backend"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/service/backup"
	"github.com/blaqkube/mysql-operator/agent/service/database"
	"github.com/blaqkube/mysql-operator/agent/service/grant"
	"github.com/blaqkube/mysql-operator/agent/service/user"
)

// A MysqlAPIController binds http requests to an api service and writes the service results to the http response
type MysqlAPIController struct {
	backup   backup.Router
	database database.MysqlDatabaseRouter
	user     user.MysqlUserRouter
	grant    grant.MysqlGrantRouter
}

// NewMysqlAPIController creates a default api controller
func NewMysqlAPIController(
	db *sql.DB,
	bck backend.Backup,
	strs map[string]backend.Storage,
) Router {
	b := backup.NewService(bck, strs)
	d := database.NewMysqlDatabaseService(db)
	u := user.NewMysqlUserService(db)
	g := grant.NewMysqlGrantService(db)
	return &MysqlAPIController{
		backup:   backup.NewController(b),
		database: database.NewMysqlDatabaseController(d),
		user:     user.NewMysqlUserController(u),
		grant:    grant.NewMysqlGrantController(g),
	}
}

// Routes returns all of the api route for the MysqlApiController
func (c *MysqlAPIController) Routes() openapi.Routes {
	routes := openapi.Routes{}
	routes = append(routes, c.backup.Routes()...)
	routes = append(routes, c.database.Routes()...)
	routes = append(routes, c.user.Routes()...)
	routes = append(routes, c.grant.Routes()...)
	return routes
}
