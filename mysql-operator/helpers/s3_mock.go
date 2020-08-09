package helpers

type StoreMockInitialize struct{}

func NewStoreMockInitialize() StoreInitializer {
	return &StoreMockInitialize{}
}

type S3MockTool struct{}

func (i *StoreMockInitialize) New(c *AWSConfig) (S3Tool, error) {
	return &S3MockTool{}, nil
}

func (s *S3MockTool) PushFileToS3(localFile, bucket, key string) error {
	return nil
}

func (s *S3MockTool) TestS3Access(bucket, directory string) error {
	return nil
}
