package backup

import (
	"errors"
	"net/http"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockService struct{}

func (s *mockService) CreateBackup(o openapi.Backup, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return openapi.Backup{
			Location:  "s3://bucket/loc/backup-1.dmp",
			Timestamp: time.Now(),
			S3access: openapi.S3Info{
				Bucket: "bucket",
				Path:   "/loc",
				AwsConfig: openapi.AwsConfig{
					AwsAccessKeyId:     "keyid",
					AwsSecretAccessKey: "secret",
					Region:             "us-east-1",
				},
			},
			Status:  "success",
			Message: "backup has succeeded",
		}, nil
	}
	if apikey == "test2" {
		return nil, errors.New("backup failed")
	}
	return nil, nil
}

func (s *mockService) DeleteBackup(backup, apikey string) (interface{}, error) {
	if apikey == "test1" {
		return nil, errors.New("not implemented")
	}
	return nil, errors.New("not implemented")
}

func (s *mockService) GetBackupByName(backup, apikey string) (interface{}, int, error) {
	if apikey == "test1" {
		return openapi.Backup{
			Location:  "s3://bucket/loc/backup-1.dmp",
			Timestamp: time.Now(),
			S3access: openapi.S3Info{
				Bucket: "bucket",
				Path:   "/loc",
				AwsConfig: openapi.AwsConfig{
					AwsAccessKeyId:     "keyid",
					AwsSecretAccessKey: "secret",
					Region:             "us-east-1",
				},
			},
			Status:  "success",
			Message: "backup has succeeded",
		}, http.StatusOK, nil
	}
	if apikey == "test2" {
		return nil, http.StatusNotFound, nil
	}
	return nil, http.StatusBadRequest, nil
}
