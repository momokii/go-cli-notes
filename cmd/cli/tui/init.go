package tui

import (
	"fmt"

	"github.com/momokii/go-cli-notes/cmd/cli/client"
)

// InitTUI initializes and validates the TUI session
// Returns an error if the session is invalid or expired
func InitTUI(apiClient *client.APIClient, authState *client.AuthState) error {
	// Check if user is authenticated (has access token)
	if !authState.IsAuthenticated() {
		return fmt.Errorf("not authenticated. please run 'kg-cli login' first")
	}

	// Validate the token is still fresh by making an API call
	// We use the stats endpoint as it's lightweight
	if !apiClient.IsAuthenticated() {
		return fmt.Errorf("no valid session found. please run 'kg-cli login'")
	}

	// Verify token is valid by making a test API call
	stats, err := apiClient.GetStats()
	if err != nil {
		return fmt.Errorf("session validation failed: %w\nplease run 'kg-cli login'", err)
	}

	// If we got here, the session is valid
	_ = stats // We don't need the stats, just validation

	return nil
}

// ValidateSession checks if the current session is still valid
// Returns true if valid, false otherwise
func ValidateSession(apiClient *client.APIClient, authState *client.AuthState) bool {
	// Check if auth state exists
	if !authState.IsAuthenticated() {
		return false
	}

	// Check if API client has a valid token
	if !apiClient.IsAuthenticated() {
		return false
	}

	// Validate token with API call
	_, err := apiClient.GetStats()
	return err == nil
}
