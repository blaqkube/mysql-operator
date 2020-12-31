package blackhole

import (
	"fmt"
	"os"
"log"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/gobuffalo/packr/v2"
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
		log.Printf("Could not open file %s, error: %v", filename, err)
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
		log.Printf("Could not create file %s, error: %v", filename, err)
		return err
	}
	defer file.Close()
	dumps := packr.New("dumps", "./dumps")
	blue, err := dumps.FindString("blue.sql")
	if err != nil {
		return err
	}

	_, err = file.WriteString(blue)
	if err != nil {
		log.Printf("Could not write file %s, error: %v", filename, err)
		return err
	}
	err = file.Sync()
	return err
}

// Delete deletes a file from the blackhole
func (s *Storage) Delete(request *openapi.BackupRequest) error {
	return nil
}
