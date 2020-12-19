package backend

import (
	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// Backup provides the interfaces required to start backup an instance
type Backup interface {
	Run(string) error
}

// Instance provides the interfaces required to start an instance
type Instance interface {
	Check(retry int) error
	Initialize() error
}

// Storage defines an interface to externalize stores
type Storage interface {
	Pull(backup *openapi.Backup, filename string) error
	Push(backup *openapi.Backup, filename string) error
	Delete(backup *openapi.Backup) error
}
