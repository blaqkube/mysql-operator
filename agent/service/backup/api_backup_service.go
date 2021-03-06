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
	CurrState string
	LastState string
	M         sync.Mutex
	States    map[string]openapi.Backup
	Status    string
	Storages  map[string]backend.Storage
}

// NewService creates a backup service
func NewService(backup backend.Backup, storages map[string]backend.Storage) *Service {
	return &Service{
		Backup:   backup,
		Status:   StatusWaiting,
		States:   map[string]openapi.Backup{},
		Storages: storages,
	}
}

// CreateBackup - create an on-demand backup
func (s *Service) CreateBackup(request openapi.BackupRequest, apiKey string) (interface{}, int, error) {
	s.M.Lock()
	defer s.M.Unlock()
	if s.Status != StatusWaiting {
		return &openapi.Backup{}, http.StatusConflict, fmt.Errorf("State %s", s.Status)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return &openapi.Backup{}, http.StatusInternalServerError, err
	}
	s.LastState = s.CurrState
	backup := openapi.Backup{
		Identifier: id,
		Bucket:     request.Bucket,
		Location:   request.Location,
		Status:     StatusWaiting,
		StartTime:  time.Now(),
	}
	s.States[id] = backup
	go runBackup(s, request, s.States[id])
	s.CurrState = id
	s.Status = StatusRunning
	return &backup, http.StatusCreated, nil
}

// GetBackupByID - Get backup from UUID
func (s *Service) GetBackupByID(uuid, apiKey string) (interface{}, int, error) {
	s.M.Lock()
	defer s.M.Unlock()
	backup, ok := s.States[uuid]
	if ok {
		return &backup, http.StatusOK, nil
	}
	return &openapi.Backup{}, http.StatusNotFound, nil
}

// GetBackups - Get backups
func (s *Service) GetBackups(apikey string) (interface{}, int, error) {
	s.M.Lock()
	defer s.M.Unlock()
	size := int32(0)
	backups := []openapi.Backup{}
	for _, v := range s.States {
		backups = append(backups, v)
		size++
	}
	return &openapi.BackupList{
		Size: size,
		Items: backups,
	}, http.StatusOK, nil
}

// runBackup is the routine that runs the backup
func runBackup(b *Service, request openapi.BackupRequest, backup openapi.Backup) {
	b.Backup.Run(fmt.Sprintf("%s.dmp", backup.Identifier))
	st := request.Backend
	if st == "" {
		st = "s3"
	}
	b.Storages[request.Backend].Push(&request, fmt.Sprintf("%s.dmp", backup.Identifier))
	b.M.Lock()
	defer b.M.Unlock()
	b.Status = StatusWaiting
	s := b.States[backup.Identifier]
	s.Status = StatusSucceeded
	t := time.Now()
	s.EndTime = &t
	b.States[backup.Identifier] = s
}
