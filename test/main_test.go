package test

import (
	"os"
	"testing"
)

// TestMain is the entry point for all tests
func TestMain(m *testing.M) {
	// Setup before tests
	setup()

	// Run tests
	exitCode := m.Run()

	// Teardown after tests
	teardown()

	os.Exit(exitCode)
}

// setup runs before all tests
func setup() {
	// You could load a test .env file here
	// or set up any other test requirements
	loadTestEnv()
}

// teardown runs after all tests
func teardown() {
	// Clean up any resources
}

// loadTestEnv loads environment variables for testing
func loadTestEnv() {
	// For testing, you can set environment variables directly
	// Or load them from a test .env file

	// Skip if already set (e.g., from CI environment)
	if os.Getenv("OPENAI_API_KEY") == "" {
		os.Setenv("OPENAI_API_KEY", "test_openai_key")
	}

	if os.Getenv("CALCOM_API_KEY") == "" {
		os.Setenv("CALCOM_API_KEY", "test_calcom_key")
	}

	if os.Getenv("CALCOM_USERNAME") == "" {
		os.Setenv("CALCOM_USERNAME", "test_username")
	}

	// Set a test port
	os.Setenv("PORT", "8081")

	// Enable debug mode for tests
	os.Setenv("DEBUG", "true")
}
