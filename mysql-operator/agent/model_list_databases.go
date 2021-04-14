package agent

// ListDatabases struct for ListDatabases
type ListDatabases struct {
	Size  int32      `json:"size,omitempty"`
	Items []Database `json:"items,omitempty"`
}
