package openapi

type Message struct {
	Code int32 `json:"code"`

	Message string `json:"message,omitempty"`
}
