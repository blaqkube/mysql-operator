package backup

import (
	"errors"
	"net/http"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockService struct{}

func (s *mockService) CreateBackup(o openapi.BackupRequest, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.Backup{
			Location:   "s3://bucket/loc/backup-1.dmp",
			Bucket:     "bucket",
			Status:     "running",
			StartTime:  time.Now(),
			Identifier: "abcd"}, nil
	}
	return nil, errors.New("backup failed")
}

func (s *mockService) GetBackups(apikey string) (interface{}, int, error) {
	if apikey == "test1" {
		return &openapi.BackupList{
			Size: 1,
			Items: []openapi.Backup{
				{
					Location:   "s3://bucket/loc/backup-1.dmp",
					Bucket:     "bucket",
					Status:     "succeeded",
					StartTime:  time.Now(),
					Identifier: "abcd",
				},
			},
		}, http.StatusOK, nil
	}
	return nil, http.StatusNotFound, errors.New("failed")
}
