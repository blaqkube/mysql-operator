package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type StorageSuite struct {
	suite.Suite
	Storage *Storage
}

func (s *StorageSuite) SetupTest() {
	s.Storage = NewStorage()
}

func (s *StorageSuite) TestStorageSuccess() {

	b := openapi.Backup{
		Bucket:   "bucket",
		Location: "location",
		Envs:     []openapi.EnvVar{},
	}

	err := s.Storage.Push(&b, "test.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Pull(&b, "test2.txt")
	assert.NoError(s.T(), err, "No Error")

	err = s.Storage.Delete(&b)
	assert.NoError(s.T(), err, "No Error")
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, &StorageSuite{})
}
