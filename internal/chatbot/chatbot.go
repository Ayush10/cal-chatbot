package chatbot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/cal-chatbot/internal/calcom"
	openai "github.com/yourusername/cal-chatbot/internal/chatbot/openai"
	"github.com/yourusername/cal-chatbot/internal/models"
)

const historyDir = "history"

// Chatbot represents the chatbot instance
type Chatbot struct {
	openaiClient *openai.Client
	calcomClient *calcom.Client
	model        string
}

// NewChatbot creates a new chatbot instance
func NewChatbot() (*Chatbot, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4-turbo"
	}

	calcomClient, err := calcom.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cal.com client: %v", err)
	}

	return &Chatbot{
		openaiClient: openai.NewClient(apiKey, calcomClient, model),
		calcomClient: calcomClient,
		model:        model,
	}, nil
}

// CheckOpenAIConnection checks if the OpenAI API key is valid and the connection is successful
func (c *Chatbot) CheckOpenAIConnection() error {
	return c.openaiClient.CheckConnection()
}

// ProcessMessage delegates to the OpenAI client
func (c *Chatbot) ProcessMessage(ctx context.Context, messages []models.ChatMessage) (string, error) {
	return c.openaiClient.ProcessMessage(ctx, messages)
}

// SaveMessage appends a message to the conversation history file
func SaveMessage(conversationID, role, message string) error {
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return err
	}
	filePath := filepath.Join(historyDir, conversationID+".txt")
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	line := fmt.Sprintf("[%s] [%s]: %s\n", time.Now().Format(time.RFC3339), role, message)
	_, err = f.WriteString(line)
	return err
}

// LoadHistory loads all messages for a conversation
func LoadHistory(conversationID string) ([]string, error) {
	filePath := filepath.Join(historyDir, conversationID+".txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

// SearchHistory searches all conversation files for a term and returns matching conversation IDs
func SearchHistory(term string) ([]string, error) {
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return nil, err
	}
	var matches []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		filePath := filepath.Join(historyDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(string(data)), strings.ToLower(term)) {
			matches = append(matches, strings.TrimSuffix(entry.Name(), ".txt"))
		}
	}
	return matches, nil
}
