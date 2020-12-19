package main

import (
	"testing"

	"github.com/blaqkube/mysql-operator/agent/backend"
	"github.com/stretchr/testify/suite"
)

type InstanceSuite struct {
	suite.Suite
	Service *backend.Instance
}

func (s *InstanceSuite) SetupSuite() {
}

func (s *InstanceSuite) TestSuccess() {

}

func TestInstanceSuite(t *testing.T) {
	suite.Run(t, &InstanceSuite{})
}
