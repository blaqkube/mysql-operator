package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InstanceSuite struct {
	suite.Suite
	Service *Instance
}

func (s *InstanceSuite) SetupTest() {
	s.Service = NewInstance()
}

func (s *InstanceSuite) TestInstanceSuccess() {

	err := s.Service.Check(1)
	assert.NoError(s.T(), err, "No Error")

	err = s.Service.Initialize()
	assert.NoError(s.T(), err, "No Error")
}

func TestInstanceSuite(t *testing.T) {
	suite.Run(t, &InstanceSuite{})
}
