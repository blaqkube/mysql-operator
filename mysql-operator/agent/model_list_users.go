package agent

// ListUsers struct for ListUsers
type ListUsers struct {
	Size  int32  `json:"size,omitempty"`
	Items []User `json:"items,omitempty"`
}
