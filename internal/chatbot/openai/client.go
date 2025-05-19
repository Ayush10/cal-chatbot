package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yourusername/cal-chatbot/internal/calcom"
	"github.com/yourusername/cal-chatbot/internal/models"

	// NOTE: handleFunctionCall will be implemented in handlers.go and imported here.
	// import "./handlers"
	goopenai "github.com/sashabaranov/go-openai"
)

type Client struct {
	openaiClient *goopenai.Client
	calcomClient *calcom.Client
	model        string
}

// ProcessMessage handles a user message and returns a response
func (c *Client) ProcessMessage(ctx context.Context, messages []models.ChatMessage) (string, error) {
	log.Printf("[INFO] ProcessMessage called with %d messages", len(messages))

	// Check for direct booking intent in the last user message
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		if lastMsg.Role == "user" && lastMsg.Booking != nil {
			log.Printf("[INFO] Direct booking detected in user message, bypassing LLM.")
			bookingBytes, err := json.Marshal(lastMsg.Booking)
			if err != nil {
				return "Sorry, I couldn't process your booking details.", err
			}
			result, err := c.bookMeeting(string(bookingBytes))
			if err != nil {
				return "Sorry, I couldn't book your meeting: " + err.Error(), err
			}
			resultJSON, _ := json.Marshal(result)
			return string(resultJSON), nil
		}
	}

	openaiMessages := ConvertToOpenAIMessages(messages)
	functionDefinitions := getFunctionDefinitions()

	// Build the request, only set Functions/FunctionCall if functions are present
	req := goopenai.ChatCompletionRequest{
		Model:    c.model,
		Messages: openaiMessages,
	}
	if len(functionDefinitions) > 0 {
		req.Functions = functionDefinitions
		req.FunctionCall = "auto"
	}

	resp, err := c.openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("[ERROR] ProcessMessage: failed to generate response: %v", err)
		return "", fmt.Errorf("failed to generate response: %v", err)
	}
	assistantMessage := resp.Choices[0].Message
	if assistantMessage.FunctionCall != nil {
		log.Printf("[INFO] ProcessMessage: function call detected: %s", assistantMessage.FunctionCall.Name)
		return HandleFunctionCall(c, ctx, assistantMessage.FunctionCall, openaiMessages)
	}
	log.Printf("[INFO] ProcessMessage: returning assistant message content")
	return assistantMessage.Content, nil
}

// getFunctionDefinitions returns the OpenAI function definitions
func getFunctionDefinitions() []goopenai.FunctionDefinition {
	return []goopenai.FunctionDefinition{
		{
			Name:        "bookMeeting",
			Description: "Book a new meeting or event in the user's Cal.com calendar.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"eventTypeId": map[string]interface{}{
						"type":        "integer",
						"description": "The ID of the event type to book.",
					},
					"startTime": map[string]interface{}{
						"type":        "string",
						"description": "The start time of the event in RFC3339 format.",
					},
					"endTime": map[string]interface{}{
						"type":        "string",
						"description": "The end time of the event in RFC3339 format.",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The name of the attendee.",
					},
					"email": map[string]interface{}{
						"type":        "string",
						"description": "The email of the attendee.",
					},
					"notes": map[string]interface{}{
						"type":        "string",
						"description": "Optional notes or description for the event.",
					},
				},
				"required": []string{"eventTypeId", "startTime", "endTime", "name", "email"},
			},
		},
		{
			Name:        "listEvents",
			Description: "List all scheduled events for a user.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"email": map[string]interface{}{
						"type":        "string",
						"description": "The email address to filter events (optional).",
					},
				},
			},
		},
		{
			Name:        "cancelEvent",
			Description: "Cancel an existing event by its event ID.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"eventId": map[string]interface{}{
						"type":        "string",
						"description": "The ID of the event to cancel.",
					},
				},
				"required": []string{"eventId"},
			},
		},
		{
			Name:        "rescheduleEvent",
			Description: "Reschedule an existing event to a new time.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"eventId": map[string]interface{}{
						"type":        "string",
						"description": "The ID of the event to reschedule.",
					},
					"newStartTime": map[string]interface{}{
						"type":        "string",
						"description": "The new start time in RFC3339 format.",
					},
					"newEndTime": map[string]interface{}{
						"type":        "string",
						"description": "The new end time in RFC3339 format.",
					},
				},
				"required": []string{"eventId", "newStartTime", "newEndTime"},
			},
		},
		{
			Name:        "checkAvailability",
			Description: "Check available time slots for a specific event type.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"eventTypeId": map[string]interface{}{
						"type":        "integer",
						"description": "The ID of the event type.",
					},
					"startDate": map[string]interface{}{
						"type":        "string",
						"description": "The start date (YYYY-MM-DD).",
					},
					"endDate": map[string]interface{}{
						"type":        "string",
						"description": "The end date (YYYY-MM-DD).",
					},
				},
				"required": []string{"eventTypeId", "startDate", "endDate"},
			},
		},
		{
			Name:        "createEventType",
			Description: "Create a new event type for the user.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "The title of the event type.",
					},
					"slug": map[string]interface{}{
						"type":        "string",
						"description": "A unique slug for the event type.",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "A description for the event type.",
					},
					"length": map[string]interface{}{
						"type":        "integer",
						"description": "Length of the event in minutes.",
					},
					"lengthUnit": map[string]interface{}{
						"type":        "string",
						"description": "Unit for the length (should be 'minutes').",
					},
				},
				"required": []string{"title", "slug", "length", "lengthUnit"},
			},
		},
		{
			Name:        "listEventTypes",
			Description: "List all event types available for booking.",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

// ConvertToOpenAIMessages converts []models.ChatMessage to []goopenai.ChatCompletionMessage
func ConvertToOpenAIMessages(messages []models.ChatMessage) []goopenai.ChatCompletionMessage {
	var openaiMessages []goopenai.ChatCompletionMessage
	for _, m := range messages {
		role := m.Role
		if role != goopenai.ChatMessageRoleSystem && role != goopenai.ChatMessageRoleUser && role != goopenai.ChatMessageRoleAssistant && role != goopenai.ChatMessageRoleFunction {
			role = goopenai.ChatMessageRoleUser
		}
		openaiMessages = append(openaiMessages, goopenai.ChatCompletionMessage{
			Role:    role,
			Content: m.Content,
		})
	}
	return openaiMessages
}

// NewClient creates a new OpenAI client wrapper
func NewClient(apiKey string, calcomClient *calcom.Client, model string) *Client {
	return &Client{
		openaiClient: goopenai.NewClient(apiKey),
		calcomClient: calcomClient,
		model:        model,
	}
}

// CheckConnection checks if the OpenAI API key is valid and the connection is successful
func (c *Client) CheckConnection() error {
	ctx := context.Background()
	_, err := c.openaiClient.CreateChatCompletion(
		ctx,
		goopenai.ChatCompletionRequest{
			Model: c.model,
			Messages: []goopenai.ChatCompletionMessage{{
				Role:    goopenai.ChatMessageRoleUser,
				Content: "ping",
			}},
		},
	)
	if err != nil {
		log.Printf("[ERROR] CheckConnection: failed to connect to OpenAI: %v", err)
	}
	return err
}
