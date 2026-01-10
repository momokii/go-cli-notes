package model

import (
	"time"

	"github.com/google/uuid"
)

// NoteType represents the type of note
type NoteType string

const (
	NoteTypeNote    NoteType = "note"
	NoteTypeDaily   NoteType = "daily"
	NoteTypeMeeting NoteType = "meeting"
	NoteTypeIdea    NoteType = "idea"
)

// Note represents a note in the system
type Note struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	Title                string     `json:"title" db:"title"`
	Content              string     `json:"content" db:"content"`
	NoteType             NoteType   `json:"note_type" db:"note_type"`
	WordCount            int        `json:"word_count" db:"word_count"`
	ReadingTimeMinutes   int        `json:"reading_time_minutes" db:"reading_time_minutes"`
	IsDeleted            bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	LastAccessedAt       *time.Time `json:"last_accessed_at,omitempty" db:"last_accessed_at"`
	AccessCount          int        `json:"access_count" db:"access_count"`
	Metadata             Metadata   `json:"metadata" db:"metadata"`
	Tags                 []*Tag     `json:"tags,omitempty"` // Populated when needed
}

// Metadata represents flexible JSONB metadata for notes
type Metadata map[string]any

// CreateNoteRequest represents a note creation request
type CreateNoteRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=500"`
	Content  string   `json:"content" validate:"max=100000"` // Large limit for markdown
	NoteType NoteType `json:"note_type" validate:"omitempty,oneof=note daily meeting idea"`
}

// UpdateNoteRequest represents a note update request
type UpdateNoteRequest struct {
	Title   *string `json:"title" validate:"omitempty,min=1,max=500"`
	Content *string `json:"content" validate:"omitempty,max=100000"`
}

// ListNotesRequest represents a note list request with filters
type ListNotesRequest struct {
	Page      int      `query:"page" validate:"min=1"`
	Limit     int      `query:"limit" validate:"min=1,max=100"`
	Type      NoteType `query:"type" validate:"omitempty,oneof=note daily meeting idea"`
	TagID     *string  `query:"tag_id"`
	Search    string   `query:"search"`
	SortBy    string   `query:"sort_by" validate:"omitempty,oneof=created_at updated_at title access_count"`
	SortOrder string   `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NoteListResponse represents a paginated list of notes
type NoteListResponse struct {
	Notes      []*Note     `json:"notes"`
	Pagination *Pagination `json:"pagination"`
}

// NoteFilter represents filters for listing notes (used by repository)
type NoteFilter struct {
	Page      int
	Limit     int
	NoteType  *NoteType
	TagID     *string
	Search    string
	SortBy    string
	SortOrder string
}
