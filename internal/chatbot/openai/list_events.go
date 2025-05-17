package openai

import (
	"encoding/json"
	"fmt"
	"log"
)

// listEvents handles the listEvents function call
func (c *Client) listEvents(args string) (interface{}, error) {
	log.Printf("[INFO] listEvents called with args: %s", args)
	var params struct {
		Email string `json:"email"`
	}

	if err := json.Unmarshal([]byte(args), &params); err != nil {
		log.Printf("[ERROR] listEvents: failed to parse args: %v", err)
		return nil, fmt.Errorf("failed to parse list events parameters: %v", err)
	}

	log.Printf("[INFO] listEvents: fetching events for email: %s", params.Email)
	events, err := c.calcomClient.GetEvents(params.Email)
	if err != nil {
		log.Printf("[ERROR] listEvents: failed to list events: %v", err)
		return nil, fmt.Errorf("failed to list events: %v", err)
	}

	log.Printf("[INFO] listEvents: events fetched successfully for %s", params.Email)
	if len(events) == 0 {
		return "You have no scheduled events.", nil
	}

	// Format the events for user-friendly display
	result := "Here are your scheduled events:\n"
	for _, event := range events {
		result += fmt.Sprintf("- %s: %s to %s\n  %s\n", event.Title, event.StartTime.Format("2006-01-02 15:04"), event.EndTime.Format("15:04"), event.Description)
	}
	return result, nil
}
