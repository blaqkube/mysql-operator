package backup

import (
	"database/sql"
	"net/http"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/mysql"
)

// MysqlBackupService is a service that implents the logic for the MysqlBackupServicer
// This service should implement the business logic for every endpoint for the MysqlBackup API.
// Include any external packages or services that will be required by this service.
type MysqlBackupService struct {
	DB       *sql.DB
	S3Backup mysql.S3MysqlBackup
}

// NewMysqlBackupService creates a MySQL backup service
func NewMysqlBackupService(
	db *sql.DB,
	s3 mysql.S3MysqlBackup,
) MysqlBackupServicer {
	return &MysqlBackupService{
		S3Backup: s3,
		DB:       db,
	}
}

// CreateBackup - create an on-demand backup
func (s *MysqlBackupService) CreateBackup(backup openapi.Backup, apiKey string) (interface{}, error) {
	// TODO - update CreateBackup with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	// mysqldump --all-databases --single-transaction -h 127.0.0.1 > mysql.backup.sql
	my := mysql.NewS3MysqlBackup()
	b, err := my.InitializeBackup(backup)
	if err != nil {
		return nil, err
	}
	go my.ExecuteBackup(*b)
	return b, nil
}

// DeleteBackup - Deletes a backup
func (s *MysqlBackupService) DeleteBackup(backup string, apiKey string) (interface{}, error) {
	// TODO - update DeleteBackup with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	b := openapi.Message{Code: int32(http.StatusNotImplemented), Message: "Not Implemented"}
	return b, nil
}

// GetBackupByName - Get backup properties
func (s *MysqlBackupService) GetBackupByName(backup string, apiKey string) (interface{}, int, error) {
	// TODO - update GetBackupByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	t, err := time.Parse(time.RFC3339, backup)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	b, err := s.S3Backup.GetBackup(t)
	if err != nil {
		return nil, http.StatusNotFound, nil
	}
	return b, http.StatusOK, nil
}
