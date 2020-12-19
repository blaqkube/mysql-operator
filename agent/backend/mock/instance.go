package mock

import (
	"github.com/stretchr/testify/mock"
)

// NewInstance takes a S3 connection and creates a default storage
func NewInstance() *Instance {
	return &Instance{}
}

// Instance is the default storage for S3
type Instance struct {
	mock.Mock
}

// Check wait for the instance to start
func (db *Instance) Check(retry int) error {
	return nil
}

// Initialize creates a set of queries to connect to the database from tools
func (db *Instance) Initialize() error {
	return nil
}
