package openapi

type BackupRequest struct {
	Backend string `json:"backend"`

	Bucket string `json:"bucket"`

	Location string `json:"location"`

	Envs []EnvVar `json:"envs,omitempty"`
}
