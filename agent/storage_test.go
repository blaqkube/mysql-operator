package main

import (
	"testing"

	"github.com/blaqkube/mysql-operator/agent/backend"
	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	Service backend.Storage
}

func (s *StorageSuite) SetupSuite() {
}

func (s *StorageSuite) TestSuccess() {
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, &StorageSuite{})
}
