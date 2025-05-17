package models

import "time"

// Event represents a Cal.com event
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Status      string    `json:"status"`
	Location    string    `json:"location,omitempty"`
}

// BookingRequest represents the parameters needed to book a new event
type BookingRequest struct {
	EventTypeID int       `json:"eventTypeId"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Notes       string    `json:"notes,omitempty"`
	Location    string    `json:"location,omitempty"`
	Title       string    `json:"title,omitempty"`
}

// AvailabilityRequest represents the parameters to check availability
type AvailabilityRequest struct {
	EventTypeID int       `json:"eventTypeId"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}

// EventTypeCreateRequest represents the parameters needed to create a new event type
type EventTypeCreateRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Length      int    `json:"length"`
	LengthUnit  string `json:"lengthUnit"`
}

// EventType represents event types fetched from Cal.com for random selection
type EventType struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Length      int    `json:"length"`
	LengthUnit  string `json:"lengthUnit"`
}
