package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// AuthState holds the authentication state
type AuthState struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
}

const authFileName = "auth.json"

// getAuthFilePath returns the path to the auth file
func getAuthFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "kg-cli")
	return filepath.Join(configDir, authFileName), nil
}

// LoadAuthState loads the authentication state from disk
func LoadAuthState() (*AuthState, error) {
	authPath, err := getAuthFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(authPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No auth file exists, return empty state
			return &AuthState{}, nil
		}
		return nil, fmt.Errorf("read auth file: %w", err)
	}

	var state AuthState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal auth state: %w", err)
	}

	return &state, nil
}

// SaveAuthState saves the authentication state to disk
func SaveAuthState(state *AuthState) error {
	authPath, err := getAuthFilePath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(authPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal auth state: %w", err)
	}

	// Write with restricted permissions (owner read/write only)
	if err := os.WriteFile(authPath, data, 0600); err != nil {
		return fmt.Errorf("write auth file: %w", err)
	}

	return nil
}

// ClearAuthState removes the authentication state from disk
func ClearAuthState() error {
	authPath, err := getAuthFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(authPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove auth file: %w", err)
	}

	return nil
}

// IsAuthenticated returns true if the user is authenticated
func (s *AuthState) IsAuthenticated() bool {
	return s.AccessToken != ""
}

// GetUserID parses and returns the user ID
func (s *AuthState) GetUserID() (uuid.UUID, error) {
	if s.UserID == "" {
		return uuid.Nil, errors.New("no user ID in auth state")
	}

	return uuid.Parse(s.UserID)
}

// ApplyToClient applies the auth state to an API client
func (s *AuthState) ApplyToClient(client *APIClient) {
	if s.AccessToken != "" {
		client.SetTokens(s.AccessToken, s.RefreshToken)
	}
}
