package backup

import (
	"errors"
	"net/http"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockService struct{}

func (s *mockService) CreateBackup(o openapi.BackupRequest, apikey string) (interface{}, int, error) {
	if apikey == "test1" {
		return openapi.Backup{
			Location:   "s3://bucket/loc/backup-1.dmp",
			Bucket:     "bucket",
			Status:     "running",
			StartTime:  time.Now(),
			Identifier: "abcd"}, http.StatusCreated, nil
	}
	return nil, http.StatusConflict, errors.New("backup failed")
}

func (s *mockService) GetBackupByID(uuid, apikey string) (interface{}, int, error) {
	if apikey == "test1" {
		return &openapi.Backup{
			Location:   "/loc/backup-1.dmp",
			Bucket:     "bucket",
			Status:     "Succeeded",
			StartTime:  time.Now(),
			Identifier: "abcd",
		}, http.StatusOK, nil
	}
	return nil, http.StatusNotFound, errors.New("failed")
}
