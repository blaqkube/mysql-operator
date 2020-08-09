package helpers

import "errors"

type StoreMockInitialize struct{}

func NewStoreMockInitialize() StoreInitializer {
	return &StoreMockInitialize{}
}

type S3MockTool struct{}

func (i *StoreMockInitialize) New(c *AWSConfig) (S3Tool, error) {
	if c != nil && c.Region == "fail" {
		return &S3MockTool{}, errors.New("fail")
	}
	return &S3MockTool{}, nil
}

func (s *S3MockTool) PushFileToS3(localFile, bucket, key string) error {
	if bucket == "fail" {
		return errors.New("fail")
	}
	return nil
}

func (s *S3MockTool) TestS3Access(bucket, directory string) error {
	if bucket == "fail" {
		return errors.New("fail")
	}
	return nil
}
