package openai

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/yourusername/cal-chatbot/internal/models"
)

// parseTimeFromText tries to extract a time (e.g., 3pm, 14:00) from a string and returns it as time.Time (today's date)
func parseTimeFromText(text string) (time.Time, error) {
	// This is a simple implementation; you may want to improve it for more formats
	re := regexp.MustCompile(`(\d{1,2})(?::(\d{2}))?\s*(am|pm)?`)
	matches := re.FindStringSubmatch(strings.ToLower(text))
	if len(matches) == 0 {
		return time.Time{}, fmt.Errorf("no time found")
	}
	hour := 0
	minute := 0
	fmt.Sscanf(matches[1], "%d", &hour)
	if matches[2] != "" {
		fmt.Sscanf(matches[2], "%d", &minute)
	}
	if matches[3] == "pm" && hour < 12 {
		hour += 12
	}
	if matches[3] == "am" && hour == 12 {
		hour = 0
	}
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
}

// cancelEvent handles the cancelEvent function call
func (c *Client) cancelEvent(args string) (interface{}, error) {
	log.Printf("[INFO] cancelEvent called with args: %s", args)
	var params struct {
		EventID  string `json:"eventId"`
		TimeText string `json:"timeText,omitempty"`
	}

	if err := json.Unmarshal([]byte(args), &params); err != nil {
		log.Printf("[ERROR] cancelEvent: failed to parse args: %v", err)
		return nil, fmt.Errorf("failed to parse cancel event parameters: %v", err)
	}

	// If EventID is provided, cancel directly
	if params.EventID != "" {
		log.Printf("[INFO] cancelEvent: cancelling event with ID: %s", params.EventID)
		err := c.calcomClient.CancelEvent(params.EventID)
		if err != nil {
			log.Printf("[ERROR] cancelEvent: failed to cancel event: %v", err)
			return nil, fmt.Errorf("failed to cancel event: %v", err)
		}
		log.Printf("[INFO] cancelEvent: event cancelled successfully: %s", params.EventID)
		return map[string]bool{"success": true}, nil
	}

	// If no EventID, try to parse time from args and find the event at that time
	if params.TimeText != "" {
		cancelTime, err := parseTimeFromText(params.TimeText)
		if err != nil {
			return nil, fmt.Errorf("Could not parse time from your request.")
		}
		events, err := c.calcomClient.GetEvents("") // Get all events
		if err != nil {
			return nil, fmt.Errorf("Could not retrieve events to find the one to cancel.")
		}
		for _, event := range events {
			if event.StartTime.Hour() == cancelTime.Hour() && event.StartTime.Minute() == cancelTime.Minute() && event.StartTime.Day() == cancelTime.Day() {
				err := c.calcomClient.CancelEvent(event.ID)
				if err != nil {
					return nil, fmt.Errorf("Failed to cancel event at %s: %v", cancelTime.Format("15:04"), err)
				}
				return fmt.Sprintf("Event '%s' at %s has been canceled.", event.Title, cancelTime.Format("15:04")), nil
			}
		}
		return nil, fmt.Errorf("I couldn't find an event at %s to cancel.", cancelTime.Format("15:04"))
	}

	return nil, fmt.Errorf("Please specify the event ID or the time of the event you want to cancel.")
}

// checkAvailability handles the checkAvailability function call
func (c *Client) checkAvailability(args string) (interface{}, error) {
	log.Printf("[INFO] checkAvailability called with args: %s", args)
	var params struct {
		EventTypeID int    `json:"eventTypeId"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
	}

	if err := json.Unmarshal([]byte(args), &params); err != nil {
		log.Printf("[ERROR] checkAvailability: failed to parse args: %v", err)
		return nil, fmt.Errorf("failed to parse availability parameters: %v", err)
	}

	startDate, err := time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		log.Printf("[ERROR] checkAvailability: invalid start date format: %v", err)
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		log.Printf("[ERROR] checkAvailability: invalid end date format: %v", err)
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	endDate = endDate.Add(24 * time.Hour)

	log.Printf("[INFO] checkAvailability: checking slots for eventTypeId=%d from %s to %s", params.EventTypeID, params.StartDate, params.EndDate)
	slots, err := c.calcomClient.GetAvailableSlots(params.EventTypeID, startDate, endDate)
	if err != nil {
		log.Printf("[ERROR] checkAvailability: failed to check availability: %v", err)
		return nil, fmt.Errorf("failed to check availability: %v", err)
	}

	log.Printf("[INFO] checkAvailability: slots fetched successfully for eventTypeId=%d", params.EventTypeID)
	return map[string]interface{}{
		"availableSlots": slots,
	}, nil
}

// rescheduleEvent handles the rescheduleEvent function call
func (c *Client) rescheduleEvent(args string) (interface{}, error) {
	log.Printf("[INFO] rescheduleEvent called with args: %s", args)
	var params struct {
		EventID      string `json:"eventId"`
		NewStartTime string `json:"newStartTime"`
		NewEndTime   string `json:"newEndTime"`
	}

	if err := json.Unmarshal([]byte(args), &params); err != nil {
		log.Printf("[ERROR] rescheduleEvent: failed to parse args: %v", err)
		return nil, fmt.Errorf("failed to parse reschedule parameters: %v", err)
	}

	newStartTime, err := time.Parse(time.RFC3339, params.NewStartTime)
	if err != nil {
		log.Printf("[ERROR] rescheduleEvent: invalid new start time format: %v", err)
		return nil, fmt.Errorf("invalid new start time format: %v", err)
	}

	newEndTime, err := time.Parse(time.RFC3339, params.NewEndTime)
	if err != nil {
		log.Printf("[ERROR] rescheduleEvent: invalid new end time format: %v", err)
		return nil, fmt.Errorf("invalid new end time format: %v", err)
	}

	log.Printf("[INFO] rescheduleEvent: rescheduling event %s to %s - %s", params.EventID, params.NewStartTime, params.NewEndTime)
	event, err := c.calcomClient.RescheduleEvent(params.EventID, newStartTime, newEndTime)
	if err != nil {
		log.Printf("[ERROR] rescheduleEvent: failed to reschedule event: %v", err)
		return nil, fmt.Errorf("failed to reschedule event: %v", err)
	}

	log.Printf("[INFO] rescheduleEvent: event rescheduled successfully: %+v", event)
	return event, nil
}

// createEventType handles the createEventType function call
func (c *Client) createEventType(args string) (interface{}, error) {
	log.Printf("[INFO] createEventType called with args: %s", args)
	var params struct {
		Title       string `json:"title"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Length      int    `json:"length"`
		LengthUnit  string `json:"lengthUnit"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		log.Printf("[ERROR] createEventType: failed to parse args: %v", err)
		return nil, fmt.Errorf("failed to parse event type creation parameters: %v", err)
	}
	request := models.EventTypeCreateRequest{
		Title:       params.Title,
		Slug:        params.Slug,
		Description: params.Description,
		Length:      params.Length,
		LengthUnit:  params.LengthUnit,
	}
	log.Printf("[INFO] createEventType: creating event type %+v", request)
	result, err := c.calcomClient.CreateEventType(request)
	if err != nil {
		log.Printf("[ERROR] createEventType: failed to create event type: %v", err)
		return nil, fmt.Errorf("failed to create event type: %v", err)
	}
	log.Printf("[INFO] createEventType: event type created successfully: %+v", result)
	return result, nil
}

// containsIgnoreCase checks if substr is in s, case-insensitive
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// listEventTypes handles the listEventTypes function call
func (c *Client) listEventTypes(args string) (interface{}, error) {
	log.Printf("[INFO] listEventTypes called")
	eventTypes, err := c.calcomClient.GetEventTypes()
	if err != nil {
		log.Printf("[ERROR] listEventTypes: failed to fetch event types: %v", err)
		return nil, fmt.Errorf("failed to fetch event types: %v", err)
	}
	if len(eventTypes) == 0 {
		return "You have no event types set up.", nil
	}
	result := "Here are your available event types:\n"
	for _, et := range eventTypes {
		result += fmt.Sprintf("- %s (%s): %s [%d %s]\n", et.Title, et.Slug, et.Description, et.Length, et.LengthUnit)
	}
	return result, nil
}
