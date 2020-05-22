/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore 
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"errors"
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
func (s *MysqlApiService) CreateBackup(backup Backup, apiKey string) (interface{}, error) {
	// TODO - update CreateBackup with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'CreateBackup' not implemented")
}

// DeleteBackup - Deletes a backup
func (s *MysqlApiService) DeleteBackup(backup string, apiKey string) (interface{}, error) {
	// TODO - update DeleteBackup with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeleteBackup' not implemented")
}

// GetBackupByName - Get backup properties
func (s *MysqlApiService) GetBackupByName(backup string, apiKey string) (interface{}, error) {
	// TODO - update GetBackupByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'GetBackupByName' not implemented")
}
