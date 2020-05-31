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
	"errors"
	"fmt"
	"net/http"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"
)

// MysqlApiService is a service that implents the logic for the MysqlApiServicer
// This service should implement the business logic for every endpoint for the MysqlApi API.
// Include any external packages or services that will be required by this service.
type MysqlApiService struct {
}

// NewMysqlApiService creates a default api service
func NewMysqlApiService() MysqlApiServicer {
	return &MysqlApiService{}
}

// CreateBackup - create an on-demand backup
func (s *MysqlApiService) CreateBackup(backup openapi.Backup, apiKey string) (interface{}, error) {
	// TODO - update CreateBackup with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	// mysqldump --all-databases --single-transaction -h 127.0.0.1 > mysql.backup.sql
	b, err := InitializeBackup(backup)
	if err != nil {
		return nil, err
	}
	go ExecuteBackup(*b)
	return b, nil
}

// CreateDatabase - create an on-demand database
func (s *MysqlApiService) CreateDatabase(body map[string]interface{}, apiKey string) (interface{}, error) {
	// TODO - update CreateDatabase with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	fmt.Printf("Connect to database\n")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	if w, ok := body["name"].(string); ok {
		_, err = db.Exec("create database " + w)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return nil, err
		}
		return body, nil
	}
	return nil, errors.New("Unknown Name")
}

// CreateUser - create an on-demand user
func (s *MysqlApiService) CreateUser(user openapi.User, apiKey string) (interface{}, error) {
	// TODO - update CreateUser with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	fmt.Printf("Connect to database\n")
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	sql := fmt.Sprintf(
		"create user '%s'@'%%' identified by '%s'",
		user.Username,
		user.Password,
	)
	fmt.Println(sql)
	_, err = db.Exec(sql)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	for _, v := range user.Grants {
		sql = fmt.Sprintf("grant GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%'", v.Database, user.Username)
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("Error granting priviles; %v\n", err)
			return nil, err
		}
	}
	return user, nil
}

// DeleteBackup - Deletes a backup
func (s *MysqlApiService) DeleteBackup(backup string, apiKey string) (interface{}, error) {
	// TODO - update DeleteBackup with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeleteBackup' not implemented")
}

// DeleteDatabase - Deletes a database
func (s *MysqlApiService) DeleteDatabase(database string, apiKey string) (interface{}, error) {
	// TODO - update DeleteDatabase with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeleteDatabase' not implemented")
}

// DeleteUser - Deletes a user
func (s *MysqlApiService) DeleteUser(user string, apiKey string) (interface{}, error) {
	// TODO - update DeleteUser with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeleteUser' not implemented")
}

// GetBackupByName - Get backup properties
func (s *MysqlApiService) GetBackupByName(backup string, apiKey string) (interface{}, int, error) {
	// TODO - update GetBackupByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	t, err := time.Parse(time.RFC3339, backup)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	b, ok := backups[t]
	if !ok {
		return nil, http.StatusNotFound, nil
	}
	return b, http.StatusOK, nil
}

// GetDatabaseByName - Get Database properties
func (s *MysqlApiService) GetDatabaseByName(database string, apiKey string) (interface{}, error) {
	// TODO - update GetDatabaseByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var name string
	err = db.QueryRow("SELECT schema_name FROM information_schema.schemata where schema_name=?", database).Scan(&name)
	if err != nil {
		return nil, err
	}
	return map[string]string{"name": name}, nil
}

// GetDatabases - list all databases
func (s *MysqlApiService) GetDatabases(apiKey string) (interface{}, error) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		return nil, err
	}
	results, err := db.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		return nil, err
	}
	databases := []map[string]string{}
	for results.Next() {
		var name string
		err = results.Scan(&name)
		database := map[string]string{"name": name}
		databases = append(databases, database)
	}
	return databases, nil
}

// GetUserByName - Get user properties
func (s *MysqlApiService) GetUserByName(user string, apiKey string) (interface{}, error) {
	// TODO - update GetUserByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	var name string
	err = db.QueryRow("SELECT User FROM mysql.user where User=?", user).Scan(&name)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	if user != name {
		return nil, errors.New("User not found")
	}
	return openapi.User{Username: name}, nil
}

// GetUsers - list all users
func (s *MysqlApiService) GetUsers(apiKey string) (interface{}, error) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
	defer db.Close()
	if err != nil {
		return nil, err
	}
	results, err := db.Query("SELECT User FROM mysql.user WHERE Host='%'")
	if err != nil {
		return nil, err
	}
	users := []openapi.User{}
	for results.Next() {
		var name string
		err = results.Scan(&name)
		user := openapi.User{Username: name}
		users = append(users, user)
	}
	return users, nil
}
