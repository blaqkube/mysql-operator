package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWSConfig is an AWS configuration that can be used to access a store
type AWSConfig struct {
	AccessKey string `json:"aws_access_key_id"`
	SecretKey string `json:"aws_secret_access_key"`
	Region    string `json:"region"`
}

type S3Tool interface {
	PushFileToS3(localFile, bucket, key string) error
	TestS3Access(bucket, directory string) error
}

type StoreInitializer interface {
	New(*AWSConfig) (S3Tool, error)
}

type StoreDefaultInitialize struct{}

func NewStoreDefaultInitialize() StoreInitializer {
	return &StoreDefaultInitialize{}
}

type S3DefaultTool struct {
	session    *session.Session
	defaultDir string
}

func (i *StoreDefaultInitialize) New(c *AWSConfig) (S3Tool, error) {
	defaultDir := "/tmp"
	a := &aws.Config{
		Region:      aws.String(c.Region),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	}

	sess, err := session.NewSession(a)
	if err != nil {
		return nil, err
	}

	return &S3DefaultTool{
		session:    sess,
		defaultDir: defaultDir,
	}, nil
}

func (s *S3DefaultTool) PushFileToS3(localFile, bucket, key string) error {
	file, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)

	_, err = s3.New(s.session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})
	return err
}

func (s *S3DefaultTool) TestS3Access(bucket, directory string) error {
	d1 := []byte("content\n")
	err := ioutil.WriteFile(
		fmt.Sprintf("%s/%s", s.defaultDir, "manifest.txt"),
		d1,
		0644,
	)
	if err != nil {
		return err
	}
	err = s.PushFileToS3(
		fmt.Sprintf("%s/%s", s.defaultDir, "manifest.txt"),
		bucket,
		fmt.Sprintf("%s/%s", directory, "manifest.txt"),
	)
	return err
}
