package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cal-chatbot/internal/api"
	"github.com/yourusername/cal-chatbot/internal/models"
)

// MockChatbot is a simplified mock chatbot for testing
type MockChatbot struct{}

// ProcessMessage mocks the chatbot's message processing
func (m *MockChatbot) ProcessMessage(ctx interface{}, message string) (string, error) {
	// Return different responses based on the message
	if message == "help me book a meeting" {
		return "I'd be happy to help you book a meeting. What date and time works for you?", nil
	} else if message == "show me my scheduled events" {
		return "You have 2 scheduled events: 'Test Meeting 1' tomorrow and 'Test Meeting 2' in 2 days.", nil
	} else if message == "cancel my event at 3pm today" {
		return "I've canceled your event at 3pm today.", nil
	} else {
		return "I'm here to help you manage your calendar. You can ask me to book a meeting, show your scheduled events, or cancel an event.", nil
	}
}

// TestIntegration tests the entire flow from API to response
func TestIntegration(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock chatbot
	mockChatbot := &MockChatbot{}

	// Create API handlers with mock chatbot
	handler := api.NewHandler(mockChatbot)

	// Create a test router
	router := gin.New()
	router.POST("/api/chat", handler.HandleChat)

	testCases := []struct {
		name           string
		message        string
		expectedStatus int
		expectedPhrase string
	}{
		{
			name:           "Book Meeting",
			message:        "help me book a meeting",
			expectedStatus: http.StatusOK,
			expectedPhrase: "help you book a meeting",
		},
		{
			name:           "List Events",
			message:        "show me my scheduled events",
			expectedStatus: http.StatusOK,
			expectedPhrase: "scheduled events",
		},
		{
			name:           "Cancel Event",
			message:        "cancel my event at 3pm today",
			expectedStatus: http.StatusOK,
			expectedPhrase: "canceled your event",
		},
		{
			name:           "General Question",
			message:        "what can you do?",
			expectedStatus: http.StatusOK,
			expectedPhrase: "manage your calendar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			chatRequest := models.ChatRequest{
				Message: tc.message,
			}
			requestBody, _ := json.Marshal(chatRequest)

			// Create request
			req, _ := http.NewRequest("POST", "/api/chat", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.Code)
			}

			// Parse response
			var response models.ChatResponse
			err := json.Unmarshal(resp.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("Failed to parse response: %v", err)
			}

			// Check that response contains expected phrase
			if response.Message == "" {
				t.Error("Empty response message")
			} else if !contains(response.Message, tc.expectedPhrase) {
				t.Errorf("Expected response to contain '%s', got '%s'", tc.expectedPhrase, response.Message)
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
