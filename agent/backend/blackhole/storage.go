package blackhole

import (
	"fmt"
	"os"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// NewStorage takes a S3 connection and creates a default storage
func NewStorage() *Storage {
	return &Storage{}
}

// Storage is the default storage for S3
type Storage struct {
}

// Push pushes a file to blaqhole bucket
func (s *Storage) Push(request *openapi.BackupRequest, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	var size int64 = fileInfo.Size()
	fmt.Printf(
		"Copying file %s (size: %d) to %s:%s",
		filename,
		size,
		request.Bucket,
		request.Location,
	)
	return nil
}

// Pull pull a file from the blackhole
func (s *Storage) Pull(request *openapi.BackupRequest, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("select '1';\n")
	if err != nil {
		return err
	}
	err = file.Sync()
	return err
}

// Delete deletes a file from the blackhole
func (s *Storage) Delete(request *openapi.BackupRequest) error {
	return nil
}
