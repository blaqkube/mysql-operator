package helpers

import (
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"

	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
)

func TestPushFileToS3(t *testing.T) {

	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Endpoint:         aws.String(ts.URL),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Create session and bucket
	s := session.New(s3Config)
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String("bucket"),
	}
	s3Client := s3.New(s)
	_, err := s3Client.CreateBucket(cparams)
	assert.Equal(t, nil, err, "should succeed")

	// Perform the 1st test
	c := &S3DefaultTool{
		session:    s,
		defaultDir: "/tmp",
	}
	err = c.PushFileToS3("s3.go", "bucket", "/greg/demo.txt")
	assert.Equal(t, nil, err, "should succeed")

	// Perform the 2nd test
	err = c.PushFileToS3("doesnotexist.go", "bucket", "/greg/demo.txt")
	assert.NotEqual(t, nil, err, "should succeed")

	// Check the file has been pushed
	lparams := &s3.ListObjectsInput{
		Bucket: aws.String("bucket"),
	}
	resp, _ := s3Client.ListObjects(lparams)
	assert.Equal(t, "greg/demo.txt", *resp.Contents[0].Key, "should succeed")
}

func TestS3Access(t *testing.T) {
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Endpoint:         aws.String(ts.URL),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Create session and bucket
	s := session.New(s3Config)
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String("test"),
	}
	s3Client := s3.New(s)
	_, err := s3Client.CreateBucket(cparams)
	assert.Equal(t, nil, err, "should succeed")

	// Perform the 1st test
	c := &S3DefaultTool{
		session:    s,
		defaultDir: "/tmp",
	}
	err = c.TestS3Access("test", "/greg")
	assert.Equal(t, nil, err, "should succeed")

	// Perform the 2nd test
	c.defaultDir = "/rooot"
	err = c.TestS3Access("test", "/greg")
	assert.Error(t, err, "should succeed")

	// Check the file has been pushed
	lparams := &s3.ListObjectsInput{
		Bucket: aws.String("test"),
	}
	resp, _ := s3Client.ListObjects(lparams)
	assert.Equal(t, "greg/manifest.txt", *resp.Contents[0].Key, "should succeed")
}
