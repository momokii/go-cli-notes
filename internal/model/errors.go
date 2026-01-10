package model

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrForbidden     = errors.New("forbidden")
	ErrValidation    = errors.New("validation failed")
	ErrDuplicate     = errors.New("resource already exists")
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

// Common API errors
var (
	ErrAPIUserNotFound     = NewAPIError("USER_NOT_FOUND", "User not found")
	ErrAPIInvalidEmail     = NewAPIError("INVALID_EMAIL", "Invalid email address")
	ErrAPIInvalidPassword  = NewAPIError("INVALID_PASSWORD", "Invalid password")
	ErrAPIEmailExists      = NewAPIError("EMAIL_EXISTS", "Email already registered")
	ErrAPIUsernameExists   = NewAPIError("USERNAME_EXISTS", "Username already taken")
	ErrAPINoteNotFound     = NewAPIError("NOTE_NOT_FOUND", "Note not found")
	ErrAPITagNotFound      = NewAPIError("TAG_NOT_FOUND", "Tag not found")
	ErrAPIInvalidToken     = NewAPIError("INVALID_TOKEN", "Invalid or expired token")
	ErrAPIUnauthorized     = NewAPIError("UNAUTHORIZED", "Unauthorized access")
	ErrAPIValidationFailed = NewAPIError("VALIDATION_FAILED", "Input validation failed")
)
