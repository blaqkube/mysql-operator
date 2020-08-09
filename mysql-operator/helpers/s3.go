package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Tool interface {
	PushFileToS3(localFile, bucket, key string) error
}

type S3DefaultTool struct {
	session *session.Session
}

func NewS3DefaultTool(s *session.Session) S3Tool {
	return &S3DefaultTool{
		session: s,
	}
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
	file.Read(buffer)

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
	err := ioutil.WriteFile("/tmp/manifest.txt", d1, 0644)
	if err != nil {
		return err
	}
	err = s.PushFileToS3("/tmp/manifest.txt", bucket, fmt.Sprintf("%s/%s", directory, "/manifest.txt"))
	return err
}
