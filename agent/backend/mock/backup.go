package mock

import (
	"github.com/stretchr/testify/mock"
)

// Backup provides a mock for the database backup
type Backup struct {
	mock.Mock
}

// NewBackup instanciate a backup interface
func NewBackup() *Backup {
	return &Backup{}
}

// Run runs a backup and store it as the filename
func (m *Backup) Run(filename string) error {
	return nil
}
