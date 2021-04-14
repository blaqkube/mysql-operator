package openapi

// BackupList - The List of backups
type BackupList struct {
	Size int32 `json:"size,omitempty"`

	Items []Backup `json:"items,omitempty"`
}
