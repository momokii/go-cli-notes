package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// GetTokenExpiry returns the expiry time of the access token
// Returns zero time if the token is invalid or cannot be parsed
func (s *AuthState) GetTokenExpiry() time.Time {
	if s.AccessToken == "" {
		return time.Time{}
	}

	// Parse the token (without verifying signature for getting expiry)
	// We're just reading the claims, not validating the token
	parts := strings.Split(s.AccessToken, ".")
	if len(parts) != 3 {
		return time.Time{}
	}

	// Decode the payload (base64url encoded)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}
	}

	// Parse claims
	claims := make(map[string]interface{})
	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}
	}

	// Get exp claim
	if exp, ok := claims["exp"].(float64); ok {
		return time.Unix(int64(exp), 0)
	}

	return time.Time{}
}

// TimeUntilExpiry returns the duration until the token expires
// Returns negative duration if already expired
func (s *AuthState) TimeUntilExpiry() time.Duration {
	expiry := s.GetTokenExpiry()
	if expiry.IsZero() {
		return 0
	}
	return time.Until(expiry)
}

// IsExpiringSoon returns true if the token will expire within the given duration
func (s *AuthState) IsExpiringSoon(within time.Duration) bool {
	until := s.TimeUntilExpiry()
	return until > 0 && until <= within
}

// IsExpired returns true if the token is already expired
func (s *AuthState) IsExpired() bool {
	return s.TimeUntilExpiry() < 0
}
