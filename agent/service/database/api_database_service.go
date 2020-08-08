package database

import (
	"database/sql"
	"errors"
	"fmt"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"
)

// MysqlDatabaseService is a service that implents the logic for the MysqlDatabaseServicer
// This service should implement the business logic for every endpoint for the MysqlDatabase API.
// Include any external packages or services that will be required by this service.
type MysqlDatabaseService struct {
	DB *sql.DB
}

// NewMysqlDatabaseService creates a default api service
func NewMysqlDatabaseService(db *sql.DB) MysqlDatabaseServicer {
	return &MysqlDatabaseService{
		DB: db,
	}
}

// CreateDatabase - create an on-demand database
func (s *MysqlDatabaseService) CreateDatabase(database openapi.Database, apiKey string) (interface{}, error) {
	// TODO - update CreateDatabase with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	fmt.Printf("Connect to database\n")
	_, err := s.DB.Exec("create database " + database.Name)
	return database, err
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
	var name string
	err := s.DB.QueryRow("SELECT schema_name FROM information_schema.schemata where schema_name=?", database).Scan(&name)
	if err != nil {
		return nil, err
	}
	return &openapi.Database{Name: name}, nil
}

// GetDatabases - list all databases
func (s *MysqlDatabaseService) GetDatabases(apiKey string) (interface{}, error) {
	results, err := s.DB.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		return nil, err
	}
	databases := []openapi.Database{}
	count := int32(0)
	for results.Next() {
		var name string
		err = results.Scan(&name)
		database := openapi.Database{Name: name}
		databases = append(databases, database)
		count++
	}
	return &openapi.ListDatabases{Size: count, Items: databases}, nil
}
