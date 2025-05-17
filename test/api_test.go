package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cal-chatbot/internal/api"
	"github.com/yourusername/cal-chatbot/internal/chatbot"
	"github.com/yourusername/cal-chatbot/internal/models"
)

// TestAPIHandlers is a test suite for the API handlers
func TestAPIHandlers(t *testing.T) {
	// Skip tests if environment variables are not set
	if os.Getenv("OPENAI_API_KEY") == "" || os.Getenv("CALCOM_API_KEY") == "" || os.Getenv("CALCOM_USERNAME") == "" {
		t.Skip("Required environment variables not set, skipping API tests")
	}

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new chatbot instance
	bot, err := chatbot.NewChatbot()
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	// Create API handlers
	handler := api.NewHandler(bot)

	t.Run("HandleHealth", func(t *testing.T) {
		// Create a test router and register the health route
		router := gin.New()
		router.GET("/api/health", handler.HandleHealth)

		// Create a test request
		req, _ := http.NewRequest("GET", "/api/health", nil)
		resp := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(resp, req)

		// Check the response status code
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}

		// Parse the response body
		var response gin.H
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}

		// Check the response content
		if response["status"] != "healthy" {
			t.Fatalf("Expected status 'healthy', got %v", response["status"])
		}
	})

	t.Run("HandleChat", func(t *testing.T) {
		// Create a test router and register the chat route
		router := gin.New()
		router.POST("/api/chat", handler.HandleChat)

		// Create a test request with a chat message
		chatRequest := models.ChatRequest{
			Message: "Hello, how are you?",
		}
		requestBody, _ := json.Marshal(chatRequest)
		req, _ := http.NewRequest("POST", "/api/chat", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(resp, req)

		// Check the response status code
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}

		// Parse the response body
		var response models.ChatResponse
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}

		// Check that the response message is not empty
		if response.Message == "" {
			t.Fatal("Empty response message")
		}

		t.Logf("Response message: %s", response.Message)
	})
}
