package model

import (
	"time"

	"github.com/google/uuid"
)

// Tag represents a tag in the system
type Tag struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Color     *string   `json:"color,omitempty" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TagWithCount represents a tag with note count
type TagWithCount struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Color     *string   `json:"color,omitempty" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	NoteCount int       `json:"note_count" db:"note_count"`
}

// CreateTagRequest represents a tag creation request
type CreateTagRequest struct {
	Name  string  `json:"name" validate:"required,min=1,max=100"`
	Color *string `json:"color" validate:"omitempty,len=7"` // Hex color, e.g., #00ADD8
}

// UpdateTagRequest represents a tag update request
type UpdateTagRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=1,max=100"`
	Color *string `json:"color" validate:"omitempty,len=7"`
}

// AddTagRequest represents a request to add a tag to a note
type AddTagRequest struct {
	TagID string `json:"tag_id" validate:"required,uuid"`
}

// TagListResponse represents a paginated list of tags
type TagListResponse struct {
	Tags       []*Tag      `json:"tags"`
	Pagination *Pagination `json:"pagination"`
}
