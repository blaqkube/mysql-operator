package agent

// Message struct for Message
type Message struct {
	Code    int32  `json:"code"`
	Message string `json:"message,omitempty"`
}
