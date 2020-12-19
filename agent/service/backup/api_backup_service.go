package backup

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/blaqkube/mysql-operator/agent/backend"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	uuid "github.com/hashicorp/go-uuid"
)

const (
	// StatusWaiting defines the status when there is no backup running
	StatusWaiting = "Waiting"

	// StatusUnknown defines an unknown status
	StatusUnknown = "Unknown"

	// StatusRunning defines the status of a backup that is running
	StatusRunning = "Running"

	// StatusFailed defines the status of a backup that has failed
	StatusFailed = "Failed"

	// StatusSucceeded defines the status of a backup that has succeeded
	StatusSucceeded = "Succeeded"
)

// Service is a service that implements the logic for the MysqlBackupServicer
// This service should implement the business logic for every endpoint for the MysqlBackup API.
// Include any external packages or services that will be required by this service.
type Service struct {
	Backup    backend.Backup
	CurrState *string
	LastState *string
	M         sync.Mutex
	States    map[string]openapi.Backup
	Status    string
	Storage   backend.Storage
}

// NewService creates a backup service
func NewService(backup backend.Backup, storage backend.Storage) *Service {
	return &Service{
		Backup:  backup,
		Status:  StatusWaiting,
		States:  map[string]openapi.Backup{},
		Storage: storage,
	}
}

// CreateBackup - create an on-demand backup
func (s *Service) CreateBackup(request openapi.BackupRequest, apiKey string) (interface{}, error) {
	s.M.Lock()
	defer s.M.Unlock()
	if s.Status != StatusWaiting {
		return nil, fmt.Errorf("State %s", s.Status)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	s.LastState = s.CurrState
	s.States[id] = openapi.Backup{
		Identifier: id,
		Bucket:     request.Bucket,
		Location:   request.Location,
		Status:     StatusWaiting,
		StartTime:  time.Now(),
	}
	go runBackup(s, s.States[id])
	s.CurrState = &id
	s.Status = StatusRunning
	return &id, nil
}

// GetBackups - Get backup properties
func (s *Service) GetBackups(apiKey string) (interface{}, int, error) {
	s.M.Lock()
	defer s.M.Unlock()
	backups := []openapi.Backup{}
	size := int32(0)
	if s.CurrState != nil {
		size++
		backups = append(backups, s.States[*s.CurrState])
	}
	if s.LastState != nil {
		size++
		backups = append(backups, s.States[*s.LastState])
	}
	return &openapi.BackupList{Size: size, Items: backups}, http.StatusOK, nil
}

// runBackup is the routine that runs the backup
func runBackup(b *Service, backup openapi.Backup) {
	b.Backup.Run(fmt.Sprintf("%s.dmp", backup.Identifier))
	b.Storage.Push(&backup, fmt.Sprintf("%s.dmp", backup.Identifier))

	b.M.Lock()
	defer b.M.Unlock()
	b.Status = StatusWaiting
	s := b.States[backup.Identifier]
	s.Status = StatusSucceeded
	t := time.Now()
	s.EndTime = &t
	b.States[backup.Identifier] = s
}
