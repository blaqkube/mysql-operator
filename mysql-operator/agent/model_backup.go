package agent

import (
	"time"
)

// Backup output for a backup request
type Backup struct {
	Identifier string     `json:"identifier"`
	Bucket     string     `json:"bucket"`
	Location   string     `json:"location"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	// backup status
	Status string `json:"status"`
}
