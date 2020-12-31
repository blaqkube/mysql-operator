package s3

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// NewStorage takes a S3 connection and creates a default storage
func NewStorage() *Storage {
	return &Storage{}
}

// Storage is the default storage for S3
type Storage struct {
}

// Push pushes a file to S3
func (s *Storage) Push(request *openapi.BackupRequest, filename string) error {
	for _, v := range request.Envs {
		os.Setenv(v.Name, v.Value)
	}
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Could not open session, error: %v", err)
		return err
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Could not open file %s, error: %v", filename, err)
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(request.Bucket),
		Key:                aws.String(request.Location),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		log.Printf("Error pushing %s to %s:%s, error: %v", filename, request.Bucket, request.Location, err)
	}
	return err
}

// Pull pull a file from S3, using a different location if necessary
func (s *Storage) Pull(request *openapi.BackupRequest, filename string) error {
	for _, v := range request.Envs {
		os.Setenv(v.Name, v.Value)
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	l := request.Location
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(request.Bucket),
			Key:    aws.String(l),
		})
	return err
}

// Delete deletes a file from S3
func (s *Storage) Delete(request *openapi.BackupRequest) error {
	for _, v := range request.Envs {
		os.Setenv(v.Name, v.Value)
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	objectsToDelete := []*s3.ObjectIdentifier{
		{Key: aws.String(request.Location)},
	}
	deleteArray := s3.Delete{Objects: objectsToDelete}
	_, err = s3.New(sess).DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(request.Bucket),
		Delete: &deleteArray,
	})
	return err
}
