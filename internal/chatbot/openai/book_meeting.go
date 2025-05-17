package openai

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/yourusername/cal-chatbot/internal/models"
)

// bookMeeting handles the bookMeeting function call
func (c *Client) bookMeeting(args string) (interface{}, error) {
	log.Printf("[INFO] bookMeeting called with args: %s", args)
	var params struct {
		EventTypeID int    `json:"eventTypeId"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		Notes       string `json:"notes,omitempty"`
	}

	if err := json.Unmarshal([]byte(args), &params); err != nil {
		log.Printf("[ERROR] bookMeeting: failed to parse args: %v", err)
		return nil, fmt.Errorf("failed to parse booking parameters: %v", err)
	}

	// If email is empty or contains 'random' or 'placeholder', generate a random email
	if params.Email == "" || containsIgnoreCase(params.Email, "random") || containsIgnoreCase(params.Email, "placeholder") {
		params.Email = fmt.Sprintf("randomuser+%d@example.com", time.Now().Unix())
		log.Printf("[INFO] bookMeeting: generated random email: %s", params.Email)
	}

	// If event type is 0 or contains 'random', select a random event type
	if params.EventTypeID == 0 || containsIgnoreCase(fmt.Sprint(params.EventTypeID), "random") {
		eventTypes, err := c.calcomClient.GetEventTypes()
		if err != nil || len(eventTypes) == 0 {
			log.Printf("[ERROR] bookMeeting: could not fetch event types for random selection: %v", err)
			return nil, fmt.Errorf("could not fetch event types for random selection")
		}
		rand.Seed(time.Now().UnixNano())
		randomType := eventTypes[rand.Intn(len(eventTypes))]
		params.EventTypeID = randomType.ID
		log.Printf("[INFO] bookMeeting: selected random event type ID: %d", params.EventTypeID)
	}

	startTime, err := time.Parse(time.RFC3339, params.StartTime)
	if err != nil {
		log.Printf("[ERROR] bookMeeting: invalid start time format: %v", err)
		return nil, fmt.Errorf("invalid start time format: %v", err)
	}

	endTime, err := time.Parse(time.RFC3339, params.EndTime)
	if err != nil {
		log.Printf("[ERROR] bookMeeting: invalid end time format: %v", err)
		return nil, fmt.Errorf("invalid end time format: %v", err)
	}

	// Set a default title if not provided
	title := "Cal.com Meeting"
	if params.Name != "" {
		title = fmt.Sprintf("Meeting with %s", params.Name)
	}

	// Prevent double-booking: check for overlapping events
	events, err := c.calcomClient.GetEvents(params.Email)
	if err == nil {
		for _, event := range events {
			if (startTime.Before(event.EndTime) && endTime.After(event.StartTime)) || startTime.Equal(event.StartTime) {
				return nil, fmt.Errorf("You already have an event scheduled at that time: %s (%s to %s)", event.Title, event.StartTime.Format("2006-01-02 15:04"), event.EndTime.Format("15:04"))
			}
		}
	}

	log.Printf("[INFO] bookMeeting: booking event for %s (%s) from %s to %s", params.Name, params.Email, params.StartTime, params.EndTime)
	event, err := c.calcomClient.BookEvent(models.BookingRequest{
		EventTypeID: params.EventTypeID,
		Start:       startTime,
		End:         endTime,
		Name:        params.Name,
		Email:       params.Email,
		Notes:       params.Notes,
		Title:       title,
	})
	if err != nil {
		log.Printf("[ERROR] bookMeeting: failed to book event: %v", err)
		return nil, fmt.Errorf("failed to book event: %v", err)
	}

	log.Printf("[INFO] bookMeeting: event booked successfully: %+v", event)
	return event, nil
}
