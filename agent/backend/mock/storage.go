package mock

import (
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/stretchr/testify/mock"
)

// NewStorage takes a S3 connection and creates a default storage
func NewStorage() *Storage {
	return &Storage{}
}

// Storage is the default storage for S3
type Storage struct {
	mock.Mock
}

// Push pushes a file
func (s *Storage) Push(backup *openapi.BackupRequest, filename string) error {
	return nil
}

// Pull pull a file from S3, using a different location if necessary
func (s *Storage) Pull(backup *openapi.BackupRequest, filename string) error {
	return nil
}

// Delete deletes a file from S3
func (s *Storage) Delete(backup *openapi.BackupRequest) error {
	return nil
}
