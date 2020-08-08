package backup

import (
	"errors"
	"time"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

type mockBackupPrimitive struct{}

func (m *mockBackupPrimitive) InitializeBackup(o openapi.Backup) (*openapi.Backup, error) {
	if o.Location == "/123" {
		return nil, errors.New("error")
	}
	return &o, nil
}

func (m *mockBackupPrimitive) ExecuteBackup(openapi.Backup) {
	return
}

func (m *mockBackupPrimitive) GetBackup(t time.Time) (*openapi.Backup, error) {
	if t.IsZero() {
		return nil, errors.New("error")
	}
	return nil, nil
}

func (m *mockBackupPrimitive) PullS3File(backup *openapi.Backup, location, filename string) error {
	return nil
}

func (m *mockBackupPrimitive) PushS3File(backup *openapi.Backup, filename string) error {
	return nil
}
