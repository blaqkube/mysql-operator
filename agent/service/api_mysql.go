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
	"database/sql"

	// "github.com/blaqkube/mysql-operator/agent/backend/mysql"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	// "github.com/blaqkube/mysql-operator/agent/service/backup"
	"github.com/blaqkube/mysql-operator/agent/service/database"
	"github.com/blaqkube/mysql-operator/agent/service/user"
)

// A MysqlAPIController binds http requests to an api service and writes the service results to the http response
type MysqlAPIController struct {
	// backup   backup.MysqlBackupRouter
	database database.MysqlDatabaseRouter
	user     user.MysqlUserRouter
}

// NewMysqlAPIController creates a default api controller
func NewMysqlAPIController(
	db *sql.DB,
	// bck mysql.S3MysqlBackup,
) MysqlApiRouter {
	// b := backup.NewMysqlBackupService(bck)
	d := database.NewMysqlDatabaseService(db)
	u := user.NewMysqlUserService(db)
	return &MysqlAPIController{
		// backup:   backup.NewMysqlBackupController(b),
		database: database.NewMysqlDatabaseController(d),
		user:     user.NewMysqlUserController(u),
	}
}

// Routes returns all of the api route for the MysqlApiController
func (c *MysqlAPIController) Routes() openapi.Routes {
	routes := openapi.Routes{}
	// routes = append(routes, c.backup.Routes()...)
	routes = append(routes, c.database.Routes()...)
	routes = append(routes, c.user.Routes()...)
	return routes
}
