package mocks

import (
	"time"

	"github.com/yourusername/cal-chatbot/internal/models"
)

// MockCalcomClient is a mock implementation of the Cal.com client
type MockCalcomClient struct {
	// Mock return values
	Events         []models.Event
	AvailableSlots []time.Time
	BookedEvent    *models.Event
	Err            error
}

// NewMockCalcomClient creates a new mock Cal.com client
func NewMockCalcomClient() *MockCalcomClient {
	// Create a default mock client with some test data
	return &MockCalcomClient{
		Events: []models.Event{
			{
				ID:        "event-1",
				Title:     "Test Meeting 1",
				StartTime: time.Now().Add(24 * time.Hour),
				EndTime:   time.Now().Add(25 * time.Hour),
				Status:    "confirmed",
			},
			{
				ID:        "event-2",
				Title:     "Test Meeting 2",
				StartTime: time.Now().Add(48 * time.Hour),
				EndTime:   time.Now().Add(49 * time.Hour),
				Status:    "confirmed",
			},
		},
		AvailableSlots: []time.Time{
			time.Now().Add(72 * time.Hour),
			time.Now().Add(96 * time.Hour),
		},
		BookedEvent: &models.Event{
			ID:        "new-event",
			Title:     "New Meeting",
			StartTime: time.Now().Add(120 * time.Hour),
			EndTime:   time.Now().Add(121 * time.Hour),
			Status:    "confirmed",
		},
		Err: nil,
	}
}

// GetEvents mocks the GetEvents method
func (m *MockCalcomClient) GetEvents(email string) ([]models.Event, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Events, nil
}

// GetAvailableSlots mocks the GetAvailableSlots method
func (m *MockCalcomClient) GetAvailableSlots(eventTypeID int, startDate, endDate time.Time) ([]time.Time, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.AvailableSlots, nil
}

// BookEvent mocks the BookEvent method
func (m *MockCalcomClient) BookEvent(booking models.BookingRequest) (*models.Event, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.BookedEvent, nil
}

// CancelEvent mocks the CancelEvent method
func (m *MockCalcomClient) CancelEvent(eventID string) error {
	return m.Err
}

// RescheduleEvent mocks the RescheduleEvent method
func (m *MockCalcomClient) RescheduleEvent(eventID string, newStartTime, newEndTime time.Time) (*models.Event, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.BookedEvent, nil
}
