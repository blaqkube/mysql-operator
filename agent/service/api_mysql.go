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
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/mysql"
	"github.com/blaqkube/mysql-operator/agent/service/backup"
	"github.com/blaqkube/mysql-operator/agent/service/database"
	"github.com/blaqkube/mysql-operator/agent/service/user"
)

// A MysqlApiController binds http requests to an api service and writes the service results to the http response
type MysqlApiController struct {
	backup   backup.MysqlBackupRouter
	database database.MysqlDatabaseRouter
	user     user.MysqlUserRouter
}

// NewMysqlApiController creates a default api controller
func NewMysqlApiController(
	bck mysql.S3MysqlBackup,
) MysqlApiRouter {
	b := backup.NewMysqlBackupService(bck)
	d := database.NewMysqlDatabaseService()
	u := user.NewMysqlUserService()
	return &MysqlApiController{
		backup:   backup.NewMysqlBackupController(b),
		database: database.NewMysqlDatabaseController(d),
		user:     user.NewMysqlUserController(u),
	}
}

// Routes returns all of the api route for the MysqlApiController
func (c *MysqlApiController) Routes() openapi.Routes {
	routes := openapi.Routes{}
	routes = append(routes, c.backup.Routes()...)
	routes = append(routes, c.database.Routes()...)
	routes = append(routes, c.user.Routes()...)
	return routes
}
