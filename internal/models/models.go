package models

// APIResponse is a general response structure from Cal.com API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
