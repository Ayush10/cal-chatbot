package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

func HandleFunctionCall(c *Client, ctx context.Context, functionCall *openai.FunctionCall, messages []openai.ChatCompletionMessage) (string, error) {
	log.Printf("[INFO] HandleFunctionCall called for function: %s", functionCall.Name)
	var result interface{}
	var err error

	switch functionCall.Name {
	case "bookMeeting":
		result, err = c.bookMeeting(functionCall.Arguments)
	case "listEvents":
		result, err = c.listEvents(functionCall.Arguments)
	case "cancelEvent":
		result, err = c.cancelEvent(functionCall.Arguments)
	case "checkAvailability":
		result, err = c.checkAvailability(functionCall.Arguments)
	case "rescheduleEvent":
		result, err = c.rescheduleEvent(functionCall.Arguments)
	case "createEventType":
		result, err = c.createEventType(functionCall.Arguments)
	case "listEventTypes":
		result, err = c.listEventTypes(functionCall.Arguments)
	default:
		log.Printf("[ERROR] HandleFunctionCall: unknown function: %s", functionCall.Name)
		return "", fmt.Errorf("unknown function: %s", functionCall.Name)
	}

	if err != nil {
		log.Printf("[ERROR] HandleFunctionCall: function execution error for %s: %v", functionCall.Name, err)
		return "", fmt.Errorf("function execution error: %v", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] HandleFunctionCall: failed to marshal result for %s: %v", functionCall.Name, err)
		return "", fmt.Errorf("failed to marshal function result: %v", err)
	}

	log.Printf("[INFO] HandleFunctionCall: function %s executed successfully", functionCall.Name)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleFunction,
		Name:    functionCall.Name,
		Content: string(resultJSON),
	})

	resp, err := c.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    c.model,
			Messages: messages,
		},
	)
	if err != nil {
		log.Printf("[ERROR] HandleFunctionCall: failed to generate response after function call %s: %v", functionCall.Name, err)
		return "", fmt.Errorf("failed to generate response after function call: %v", err)
	}

	log.Printf("[INFO] HandleFunctionCall: returning response for function %s", functionCall.Name)
	return resp.Choices[0].Message.Content, nil
}
