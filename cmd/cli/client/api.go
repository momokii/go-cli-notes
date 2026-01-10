package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/momokii/go-cli-notes/internal/model"
)

// APIClient handles communication with the API
type APIClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
	refreshToken string
}

// AuthResponse holds authentication tokens
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// apiErrorResponse represents an error response from the API
type apiErrorResponse struct {
	Error string `json:"error"`
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SetTokens sets the authentication tokens
func (c *APIClient) SetTokens(accessToken, refreshToken string) {
	c.token = accessToken
	c.refreshToken = refreshToken
}

// GetToken returns the current access token
func (c *APIClient) GetToken() string {
	return c.token
}

// IsAuthenticated returns true if the client has a valid token
func (c *APIClient) IsAuthenticated() bool {
	return c.token != ""
}

// makeRequest makes an HTTP request with authentication
func (c *APIClient) makeRequest(method, path string, body interface{}, authenticated bool) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if authenticated && c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.httpClient.Do(req)
}

// formatAPIError converts an API error response into a user-friendly message
func formatAPIError(statusCode int, body []byte) error {
	// Try to parse as JSON error response
	var apiErr apiErrorResponse
	if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
		// Clean up the error message
		errMsg := apiErr.Error

		// Remove common prefixes
		errMsg = strings.TrimPrefix(errMsg, "Internal server error: ")
		errMsg = strings.TrimPrefix(errMsg, "validation failed: ")

		// Capitalize first letter
		if len(errMsg) > 0 {
			errMsg = strings.ToUpper(string(errMsg[0])) + errMsg[1:]
		}

		// Map common errors to friendly messages
		// Order matters - check more specific patterns first
		switch {
		case strings.Contains(strings.ToLower(errMsg), "username is required"):
			return fmt.Errorf("username is required")
		case strings.Contains(strings.ToLower(errMsg), "email is required"):
			return fmt.Errorf("email is required")
		case strings.Contains(strings.ToLower(errMsg), "password is required"):
			return fmt.Errorf("password is required")
		// Check for the exact API response first
		case errMsg == "Invalid email or password" || errMsg == "invalid email or password":
			return fmt.Errorf("invalid email or password")
		case strings.Contains(strings.ToLower(errMsg), "invalid credentials"):
			return fmt.Errorf("invalid email or password")
		case strings.Contains(strings.ToLower(errMsg), "invalid email") && !strings.Contains(strings.ToLower(errMsg), "credentials") && !strings.Contains(strings.ToLower(errMsg), "invalid email or password"):
			return fmt.Errorf("invalid email format")
		case strings.Contains(strings.ToLower(errMsg), "unauthorized"):
			return fmt.Errorf("not authenticated. Please run 'kg-cli login'")
		case strings.Contains(strings.ToLower(errMsg), "user already exists"):
			return fmt.Errorf("user with this email or username already exists")
		case strings.Contains(strings.ToLower(errMsg), "username already exists"):
			return fmt.Errorf("username already exists")
		case strings.Contains(strings.ToLower(errMsg), "not found"):
			return fmt.Errorf("resource not found")
		default:
			// Return cleaned error message
			return fmt.Errorf("%s", strings.TrimSuffix(errMsg, "; "))
		}
	}

	// Fallback for non-JSON errors
	return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
}

// decodeResponse decodes a JSON response
func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return formatAPIError(resp.StatusCode, body)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// Register registers a new user
func (c *APIClient) Register(username, email, password string) error {
	payload := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}

	resp, err := c.makeRequest("POST", "/api/v1/auth/register", payload, false)
	if err != nil {
		return err
	}

	return decodeResponse(resp, nil)
}

// Login authenticates a user
func (c *APIClient) Login(email, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	resp, err := c.makeRequest("POST", "/api/v1/auth/login", payload, false)
	if err != nil {
		return nil, err
	}

	var authResp AuthResponse
	if err := decodeResponse(resp, &authResp); err != nil {
		return nil, err
	}

	c.SetTokens(authResp.AccessToken, authResp.RefreshToken)
	return &authResp, nil
}

// RefreshToken refreshes the access token
func (c *APIClient) RefreshToken() error {
	if c.refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	payload := map[string]string{
		"refresh_token": c.refreshToken,
	}

	resp, err := c.makeRequest("POST", "/api/v1/auth/refresh", payload, false)
	if err != nil {
		return err
	}

	var authResp AuthResponse
	if err := decodeResponse(resp, &authResp); err != nil {
		return err
	}

	c.SetTokens(authResp.AccessToken, authResp.RefreshToken)
	return nil
}

// Logout logs out the user
func (c *APIClient) Logout() error {
	resp, err := c.makeRequest("POST", "/api/v1/auth/logout", nil, true)
	if err != nil {
		return err
	}

	c.token = ""
	c.refreshToken = ""
	return decodeResponse(resp, nil)
}

// CreateNote creates a new note
func (c *APIClient) CreateNote(req *model.CreateNoteRequest) (*model.Note, error) {
	resp, err := c.makeRequest("POST", "/api/v1/notes", req, true)
	if err != nil {
		return nil, err
	}

	var note model.Note
	if err := decodeResponse(resp, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// ListNotes lists notes with optional filters
func (c *APIClient) ListNotes(filter model.NoteFilter) ([]*model.Note, int64, error) {
	// Build query string
	path := "/api/v1/notes?page=" + fmt.Sprint(filter.Page) + "&limit=" + fmt.Sprint(filter.Limit)
	if filter.SortBy != "" {
		path += "&sort_by=" + filter.SortBy
	}
	if filter.Search != "" {
		path += "&search=" + filter.Search
	}
	if filter.TagID != nil && *filter.TagID != "" {
		path += "&tag=" + *filter.TagID
	}
	if filter.NoteType != nil {
		path += "&type=" + string(*filter.NoteType)
	}

	resp, err := c.makeRequest("GET", path, nil, true)
	if err != nil {
		return nil, 0, err
	}

	var result struct {
		Notes      []*model.Note `json:"notes"`
		Pagination struct {
			Page       int   `json:"page"`
			Limit      int   `json:"limit"`
			Total      int64 `json:"total"`
			TotalPages int   `json:"total_pages"`
		} `json:"pagination"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, 0, err
	}

	return result.Notes, result.Pagination.Total, nil
}

// GetNote retrieves a single note by ID
func (c *APIClient) GetNote(id uuid.UUID) (*model.Note, error) {
	resp, err := c.makeRequest("GET", "/api/v1/notes/"+id.String(), nil, true)
	if err != nil {
		return nil, err
	}

	var note model.Note
	if err := decodeResponse(resp, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// UpdateNote updates an existing note
func (c *APIClient) UpdateNote(id uuid.UUID, req *model.UpdateNoteRequest) error {
	resp, err := c.makeRequest("PUT", "/api/v1/notes/"+id.String(), req, true)
	if err != nil {
		return err
	}

	return decodeResponse(resp, nil)
}

// DeleteNote deletes a note
func (c *APIClient) DeleteNote(id uuid.UUID) error {
	resp, err := c.makeRequest("DELETE", "/api/v1/notes/"+id.String(), nil, true)
	if err != nil {
		return err
	}

	return decodeResponse(resp, nil)
}

// SearchNotes searches notes using full-text search
func (c *APIClient) SearchNotes(query string, page, limit int) (*model.SearchResponse, error) {
	path := fmt.Sprintf("/api/v1/search?q=%s&page=%d&limit=%d", query, page, limit)

	resp, err := c.makeRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var searchResp model.SearchResponse
	if err := decodeResponse(resp, &searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// GetTags retrieves all tags
func (c *APIClient) GetTags() ([]*model.Tag, error) {
	resp, err := c.makeRequest("GET", "/api/v1/tags", nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tags       []*model.Tag `json:"tags"`
		Pagination struct {
			Page       int `json:"page"`
			Limit      int `json:"limit"`
			Total      int64 `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"pagination"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Tags, nil
}

// CreateTag creates a new tag
func (c *APIClient) CreateTag(name string) (*model.Tag, error) {
	payload := map[string]string{"name": name}

	resp, err := c.makeRequest("POST", "/api/v1/tags", payload, true)
	if err != nil {
		return nil, err
	}

	var tag model.Tag
	if err := decodeResponse(resp, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// UpdateTag updates an existing tag
func (c *APIClient) UpdateTag(id uuid.UUID, name string) (*model.Tag, error) {
	payload := map[string]string{"name": name}

	resp, err := c.makeRequest("PUT", "/api/v1/tags/"+id.String(), payload, true)
	if err != nil {
		return nil, err
	}

	var tag model.Tag
	if err := decodeResponse(resp, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

// DeleteTag deletes a tag
func (c *APIClient) DeleteTag(id uuid.UUID) error {
	resp, err := c.makeRequest("DELETE", "/api/v1/tags/"+id.String(), nil, true)
	if err != nil {
		return err
	}

	return decodeResponse(resp, nil)
}

// GetNoteTags retrieves tags for a specific note
func (c *APIClient) GetNoteTags(noteID uuid.UUID) ([]*model.Tag, error) {
	resp, err := c.makeRequest("GET", "/api/v1/notes/"+noteID.String()+"/tags", nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tags []*model.Tag `json:"tags"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Tags, nil
}

// AddTagToNote adds a tag to a note
func (c *APIClient) AddTagToNote(noteID, tagID uuid.UUID) error {
	resp, err := c.makeRequest("POST", "/api/v1/notes/"+noteID.String()+"/tags/"+tagID.String(), nil, true)
	if err != nil {
		return err
	}

	return decodeResponse(resp, nil)
}

// RemoveTagFromNote removes a tag from a note
func (c *APIClient) RemoveTagFromNote(noteID, tagID uuid.UUID) error {
	resp, err := c.makeRequest("DELETE", "/api/v1/notes/"+noteID.String()+"/tags/"+tagID.String(), nil, true)
	if err != nil {
		return err
	}

	return decodeResponse(resp, nil)
}

// GetGraph retrieves the knowledge graph
func (c *APIClient) GetGraph() (*model.GraphResponse, error) {
	resp, err := c.makeRequest("GET", "/api/v1/notes/graph", nil, true)
	if err != nil {
		return nil, err
	}

	var graph model.GraphResponse
	if err := decodeResponse(resp, &graph); err != nil {
		return nil, err
	}

	return &graph, nil
}

// GetLinks retrieves outgoing links from a note
func (c *APIClient) GetLinks(id uuid.UUID) ([]*model.LinkDetail, error) {
	resp, err := c.makeRequest("GET", "/api/v1/notes/"+id.String()+"/links", nil, true)
	if err != nil {
		return nil, err
	}

	var links []*model.LinkDetail
	if err := decodeResponse(resp, &links); err != nil {
		return nil, err
	}

	return links, nil
}

// GetBacklinks retrieves backlinks to a note
func (c *APIClient) GetBacklinks(id uuid.UUID) ([]*model.LinkDetail, error) {
	resp, err := c.makeRequest("GET", "/api/v1/notes/"+id.String()+"/backlinks", nil, true)
	if err != nil {
		return nil, err
	}

	var links []*model.LinkDetail
	if err := decodeResponse(resp, &links); err != nil {
		return nil, err
	}

	return links, nil
}

// GetDailyNote gets or creates a daily note for a given date
func (c *APIClient) GetDailyNote(date string) (*model.Note, bool, error) {
	resp, err := c.makeRequest("GET", "/api/v1/notes/daily/"+date, nil, true)
	if err != nil {
		return nil, false, err
	}

	var result struct {
		Note       *model.Note `json:"note"`
		IsCreated  bool        `json:"is_created"`
		Date       string      `json:"date"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, false, err
	}

	return result.Note, result.IsCreated, nil
}

// GetStats retrieves user statistics
func (c *APIClient) GetStats() (*model.UserStats, error) {
	resp, err := c.makeRequest("GET", "/api/v1/stats", nil, true)
	if err != nil {
		return nil, err
	}

	var stats model.UserStats
	if err := decodeResponse(resp, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetRecentActivity retrieves recent activity
func (c *APIClient) GetRecentActivity(limit int) ([]*model.Activity, error) {
	path := fmt.Sprintf("/api/v1/activity/recent?limit=%d", limit)

	resp, err := c.makeRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Activities []*model.Activity `json:"activities"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Activities, nil
}

// GetTrendingNotes retrieves trending notes
func (c *APIClient) GetTrendingNotes(limit int) ([]*model.TrendingNote, error) {
	path := fmt.Sprintf("/api/v1/notes/trending?limit=%d", limit)

	resp, err := c.makeRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Trending []*model.TrendingNote `json:"trending"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Trending, nil
}

// GetForgottenNotes retrieves forgotten notes
func (c *APIClient) GetForgottenNotes(days, limit int) ([]*model.ForgottenNote, error) {
	path := fmt.Sprintf("/api/v1/notes/forgotten?days=%d&limit=%d", days, limit)

	resp, err := c.makeRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var result struct {
		Forgotten []*model.ForgottenNote `json:"forgotten"`
	}

	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Forgotten, nil
}
