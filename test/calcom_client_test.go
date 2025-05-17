package test

import (
	"os"
	"testing"
	"time"

	"github.com/yourusername/cal-chatbot/internal/calcom"
)

// TestCalcomClient is a test suite for the Cal.com API client
func TestCalcomClient(t *testing.T) {
	// Skip tests if environment variables are not set
	if os.Getenv("CALCOM_API_KEY") == "" || os.Getenv("CALCOM_USERNAME") == "" {
		t.Skip("CALCOM_API_KEY or CALCOM_USERNAME not set, skipping Cal.com API tests")
	}

	t.Run("NewClient", func(t *testing.T) {
		client, err := calcom.NewClient()
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		if client == nil {
			t.Fatal("Client is nil")
		}
	})

	t.Run("GetEvents", func(t *testing.T) {
		client, err := calcom.NewClient()
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		// Use a test email
		events, err := client.GetEvents("test@example.com")
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		// We just check that the call completes, not the actual events
		// as this varies depending on the account
		t.Logf("Retrieved %d events", len(events))
	})

	t.Run("GetAvailableSlots", func(t *testing.T) {
		client, err := calcom.NewClient()
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		// Use event type ID 1 (default) and a range of dates
		now := time.Now()
		end := now.Add(24 * 7 * time.Hour) // One week from now
		slots, err := client.GetAvailableSlots(1, now, end)
		if err != nil {
			t.Fatalf("Failed to get available slots: %v", err)
		}

		t.Logf("Retrieved %d available slots", len(slots))
	})

	// Note: We don't test booking, cancellation, or rescheduling
	// as those would create real events in the Cal.com account
}
