package calcom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yourusername/cal-chatbot/internal/models"
)

// Client represents a Cal.com API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	username   string
}

// NewClient creates a new Cal.com API client
func NewClient() (*Client, error) {
	apiKey := os.Getenv("CALCOM_API_KEY")
	fmt.Println("[DEBUG] CALCOM_API_KEY:", apiKey)
	baseURL := os.Getenv("CALCOM_API_URL")
	fmt.Println("[DEBUG] CALCOM_API_URL:", baseURL)
	username := os.Getenv("CALCOM_USERNAME")
	fmt.Println("[DEBUG] CALCOM_USERNAME:", username)
	if apiKey == "" {
		return nil, fmt.Errorf("CALCOM_API_KEY environment variable is not set")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL:  baseURL,
		apiKey:   apiKey,
		username: username,
	}, nil
}

// makeRequest makes an HTTP request to the Cal.com API
func (c *Client) makeRequest(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Append apiKey as a query parameter
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	if strings.Contains(url, "?") {
		url += "&apiKey=" + c.apiKey
	} else {
		url += "?apiKey=" + c.apiKey
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s (status code: %d)", respBody, resp.StatusCode)
	}

	return respBody, nil
}

// GetEvents retrieves all events for a user
func (c *Client) GetEvents(email string) ([]models.Event, error) {
	path := "/bookings"
	if email != "" {
		path = fmt.Sprintf("%s?email=%s", path, email)
	}

	respBody, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Debug: log the raw API response
	fmt.Printf("[DEBUG] Cal.com GetEvents raw response for email '%s': %s\n", email, string(respBody))

	var response struct {
		Bookings []models.Event `json:"bookings"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events: %v", err)
	}

	return response.Bookings, nil
}

// GetAvailableSlots retrieves available time slots for a specific event type
func (c *Client) GetAvailableSlots(eventTypeID int, startDate, endDate time.Time) ([]time.Time, error) {
	path := fmt.Sprintf("/availability/%s/%d", c.username, eventTypeID)
	query := struct {
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}{
		StartTime: startDate,
		EndTime:   endDate,
	}

	respBody, err := c.makeRequest(http.MethodPost, path, query)
	if err != nil {
		return nil, err
	}

	var response struct {
		Available []time.Time `json:"available"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal available slots: %v", err)
	}

	return response.Available, nil
}

// BookEvent books a new event
func (c *Client) BookEvent(booking models.BookingRequest) (*models.Event, error) {
	fmt.Printf("[DEBUG] BookEvent called with payload: %+v\n", booking)

	// Build payload according to Cal.com API reference
	payload := map[string]interface{}{
		"eventTypeId": booking.EventTypeID,
		"start":       booking.Start.Format(time.RFC3339),
		"end":         booking.End.Format(time.RFC3339),
		"responses": map[string]interface{}{
			"name":  booking.Name,
			"email": booking.Email,
			"location": map[string]interface{}{
				"value":       booking.Location,
				"optionValue": "",
			},
		},
		"timeZone":    "UTC", // You may want to make this dynamic
		"language":    "en",
		"title":       booking.Title, // Add Title to BookingRequest if not present
		"description": booking.Notes, // Use Notes as description
		"status":      "PENDING",
		"metadata":    map[string]interface{}{},
	}

	respBody, err := c.makeRequest(http.MethodPost, "/bookings", payload)
	fmt.Printf("[DEBUG] Cal.com raw response: %s\n", string(respBody))
	if err != nil {
		fmt.Printf("[ERROR] BookEvent failed: %v\n", err)
		return nil, err
	}

	var response struct {
		Booking models.Event `json:"booking"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Printf("[ERROR] Failed to unmarshal booking response: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal booking response: %v", err)
	}

	return &response.Booking, nil
}

// CancelEvent cancels an existing event
func (c *Client) CancelEvent(eventID string) error {
	path := fmt.Sprintf("/bookings/%s/cancel", eventID)
	_, err := c.makeRequest(http.MethodPost, path, nil)
	return err
}

// RescheduleEvent reschedules an existing event
func (c *Client) RescheduleEvent(eventID string, newStartTime, newEndTime time.Time) (*models.Event, error) {
	path := fmt.Sprintf("/bookings/%s/reschedule", eventID)
	body := struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}{
		Start: newStartTime,
		End:   newEndTime,
	}

	respBody, err := c.makeRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Booking models.Event `json:"booking"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rescheduled booking: %v", err)
	}

	return &response.Booking, nil
}

// CreateEventType creates a new event type
func (c *Client) CreateEventType(req models.EventTypeCreateRequest) (map[string]interface{}, error) {
	respBody, err := c.makeRequest(http.MethodPost, "/event-types", req)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event type creation response: %v", err)
	}
	return response, nil
}

// GetEventTypes fetches all event types for the user
func (c *Client) GetEventTypes() ([]models.EventType, error) {
	respBody, err := c.makeRequest(http.MethodGet, "/event-types", nil)
	if err != nil {
		return nil, err
	}
	var response struct {
		EventTypes []models.EventType `json:"eventTypes"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event types: %v", err)
	}
	return response.EventTypes, nil
}

// FindAllEventTypes fetches all event types
func (c *Client) FindAllEventTypes() ([]models.EventType, error) {
	respBody, err := c.makeRequest(http.MethodGet, "/event-types", nil)
	if err != nil {
		return nil, err
	}
	var response struct {
		EventTypes []models.EventType `json:"eventTypes"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event types: %v", err)
	}
	return response.EventTypes, nil
}

// FindAllSchedules fetches all schedules
func (c *Client) FindAllSchedules() ([]map[string]interface{}, error) {
	respBody, err := c.makeRequest(http.MethodGet, "/schedules", nil)
	if err != nil {
		return nil, err
	}
	var response struct {
		Schedules []map[string]interface{} `json:"schedules"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedules: %v", err)
	}
	return response.Schedules, nil
}

// CreateSchedule creates a new schedule
func (c *Client) CreateSchedule(name, timeZone string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"name":     name,
		"timeZone": timeZone,
	}
	respBody, err := c.makeRequest(http.MethodPost, "/schedules", payload)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule creation response: %v", err)
	}
	return response, nil
}

// GetBookableSlots fetches all bookable slots between a datetime range
func (c *Client) GetBookableSlots(start, end string) (map[string][]map[string]interface{}, error) {
	params := fmt.Sprintf("?start=%s&end=%s", start, end)
	respBody, err := c.makeRequest(http.MethodGet, "/slots"+params, nil)
	if err != nil {
		return nil, err
	}
	var response struct {
		Slots map[string][]map[string]interface{} `json:"slots"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal slots: %v", err)
	}
	return response.Slots, nil
}

// RemoveSchedule deletes a schedule by ID
func (c *Client) RemoveSchedule(scheduleID string) error {
	path := fmt.Sprintf("/schedules/%s", scheduleID)
	_, err := c.makeRequest(http.MethodDelete, path, nil)
	return err
}

// EditSchedule edits an existing schedule by ID
func (c *Client) EditSchedule(scheduleID string, updates map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/schedules/%s", scheduleID)
	respBody, err := c.makeRequest(http.MethodPatch, path, updates)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule edit response: %v", err)
	}
	return response, nil
}

// FindBooking fetches a booking by ID
func (c *Client) FindBooking(bookingID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/bookings/%s", bookingID)
	respBody, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal booking: %v", err)
	}
	return response, nil
}

// EditBooking edits an existing booking by ID
func (c *Client) EditBooking(bookingID string, updates map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/bookings/%s", bookingID)
	respBody, err := c.makeRequest(http.MethodPatch, path, updates)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal booking edit response: %v", err)
	}
	return response, nil
}

// CancelBooking cancels a booking by ID
func (c *Client) CancelBooking(bookingID string) error {
	path := fmt.Sprintf("/bookings/%s/cancel", bookingID)
	_, err := c.makeRequest(http.MethodPost, path, nil)
	return err
}
