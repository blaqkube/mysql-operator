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
	"errors"
	"net/http"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
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
	return nil, errors.New("service method 'CreateDatabase' not implemented")
}

// CreateUser - create an on-demand user
func (s *MysqlApiService) CreateUser(user openapi.User, apiKey string) (interface{}, error) {
	// TODO - update CreateUser with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'CreateUser' not implemented")
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
	return nil, errors.New("service method 'GetDatabaseByName' not implemented")
}

// GetUserByName - Get user properties
func (s *MysqlApiService) GetUserByName(user string, apiKey string) (interface{}, error) {
	// TODO - update GetUserByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'GetUserByName' not implemented")
}
