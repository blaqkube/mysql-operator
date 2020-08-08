package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDatabaseService is a service that implents the logic for the MysqlDatabaseServicer
// This service should implement the business logic for every endpoint for the MysqlDatabase API.
// Include any external packages or services that will be required by this service.
type MysqlDatabaseService struct {
}

// NewMysqlDatabaseService creates a default api service
func NewMysqlDatabaseService() MysqlDatabaseServicer {
	return &MysqlDatabaseService{}
}

// CreateDatabase - create an on-demand database
func (s *MysqlDatabaseService) CreateDatabase(body map[string]interface{}, apiKey string) (interface{}, error) {
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

// DeleteDatabase - Deletes a database
func (s *MysqlDatabaseService) DeleteDatabase(database string, apiKey string) (interface{}, error) {
	// TODO - update DeleteDatabase with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeleteDatabase' not implemented")
}

// GetDatabaseByName - Get Database properties
func (s *MysqlDatabaseService) GetDatabaseByName(database string, apiKey string) (interface{}, error) {
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
func (s *MysqlDatabaseService) GetDatabases(apiKey string) (interface{}, error) {
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
