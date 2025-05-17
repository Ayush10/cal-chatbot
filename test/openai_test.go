package test

import (
	"context"
	"os"
	"testing"

	"github.com/yourusername/cal-chatbot/internal/chatbot"
)

// TestOpenAIIntegration is a test suite for the OpenAI integration
func TestOpenAIIntegration(t *testing.T) {
	// Skip tests if OPENAI_API_KEY is not set
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping OpenAI tests")
	}

	// Skip tests if Cal.com environment variables are not set
	if os.Getenv("CALCOM_API_KEY") == "" || os.Getenv("CALCOM_USERNAME") == "" {
		t.Skip("CALCOM_API_KEY or CALCOM_USERNAME not set, skipping OpenAI tests that require Cal.com")
	}

	t.Run("NewChatbot", func(t *testing.T) {
		bot, err := chatbot.NewChatbot()
		if err != nil {
			t.Fatalf("Failed to create chatbot: %v", err)
		}
		if bot == nil {
			t.Fatal("Chatbot is nil")
		}
	})

	t.Run("ProcessMessage", func(t *testing.T) {
		bot, err := chatbot.NewChatbot()
		if err != nil {
			t.Fatalf("Failed to create chatbot: %v", err)
		}

		// Test a simple message that should not trigger a function call
		ctx := context.Background()
		response, err := bot.ProcessMessage(ctx, "Hello, how are you?")
		if err != nil {
			t.Fatalf("Failed to process message: %v", err)
		}

		if response == "" {
			t.Fatal("Empty response from chatbot")
		}

		t.Logf("Response: %s", response)
	})
}
