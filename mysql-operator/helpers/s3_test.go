package helpers

import (
	"net/http/httptest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/tree/master/backend/s3mem"
	"github.com/testify/assert"
)

func simpleTest(t *testing) {
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
	s := session.New(s3Config)
	c := NewS3DefaultTool(s)
	err := c.PushFileToS3("s3.go", "me", "/greg")
	assert.Equal(t, nil, err, "should succeed")

}
