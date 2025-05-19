package models

// ChatRequest represents a chat message from the user
// Deprecated: use Messages for full conversation context
// type ChatRequest struct {
// 	Message string `json:"message"`
// 	UserID  string `json:"userId,omitempty"`
// }
type ChatMessage struct {
	Role       string                 `json:"role"`
	Content    string                 `json:"content"`
	Booking    map[string]interface{} `json:"booking,omitempty"`
	ListEvents map[string]interface{} `json:"listEvents,omitempty"`
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
	UserID   string        `json:"userId,omitempty"`
}

// ChatResponse represents the chatbot's response
type ChatResponse struct {
	Message string `json:"message"`
}
