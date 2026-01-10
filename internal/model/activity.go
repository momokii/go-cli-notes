package model

import (
	"time"

	"github.com/google/uuid"
)

// ActionType represents the type of activity action
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionView   ActionType = "view"
	ActionSearch ActionType = "search"
	ActionDelete ActionType = "delete"
	ActionLogin  ActionType = "login"
	ActionLogout ActionType = "logout"
)

// Activity represents a user activity log entry
type Activity struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	UserID    uuid.UUID          `json:"user_id" db:"user_id"`
	NoteID    *uuid.UUID         `json:"note_id,omitempty" db:"note_id"`
	Action    ActionType         `json:"action" db:"action"`
	Metadata  ActivityMetadata   `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
}

// ActivityMetadata represents flexible metadata for activities
type ActivityMetadata map[string]any

// LogActivityRequest represents a request to log an activity
type LogActivityRequest struct {
	NoteID   *uuid.UUID
	Action   ActionType
	Metadata ActivityMetadata
}

// UserStats represents user statistics
type UserStats struct {
	TotalNotes       int64     `json:"total_notes"`
	TotalTags        int64     `json:"total_tags"`
	TotalLinks       int64     `json:"total_links"`
	TotalWords       int64     `json:"total_words"`
	NotesCreatedToday int64    `json:"notes_created_today"`
	NotesCreatedWeek int64     `json:"notes_created_week"`
	LastActivity     *time.Time `json:"last_activity,omitempty"`
}

// TrendingNote represents a note that's trending (frequently accessed)
type TrendingNote struct {
	Note         *Note `json:"note"`
	AccessCount  int   `json:"access_count"`
	RecentAccess int   `json:"recent_access"` // Access count in recent period
}

// ForgottenNote represents a note that hasn't been accessed in a while
type ForgottenNote struct {
	Note              *Note     `json:"note"`
	LastAccessedAt    time.Time `json:"last_accessed_at"`
	DaysSinceAccess   int       `json:"days_since_access"`
}
