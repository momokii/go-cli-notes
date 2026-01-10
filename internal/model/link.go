package model

import (
	"time"

	"github.com/google/uuid"
)

// Link represents a bidirectional link between notes
type Link struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	SourceNoteID   uuid.UUID  `json:"source_note_id" db:"source_note_id"`
	TargetNoteID   uuid.UUID  `json:"target_note_id" db:"target_note_id"`
	LinkContext    *string    `json:"link_context,omitempty" db:"link_context"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	TargetNote     *Note      `json:"target_note,omitempty"` // Populated when needed
	SourceNote     *Note      `json:"source_note,omitempty"` // Populated when needed
}

// CreateLinkRequest represents a manual link creation request
type CreateLinkRequest struct {
	TargetNoteID string `json:"target_note_id" validate:"required,uuid"`
}

// GraphNeighbor represents a node in the knowledge graph
type GraphNeighbor struct {
	Note        *Note   `json:"note"`
	LinkContext *string `json:"link_context,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// GraphPath represents a path between two notes
type GraphPath struct {
	Path []*Note `json:"path"` // Ordered list of notes from source to target
}

// GraphStats represents statistics about the knowledge graph
type GraphStats struct {
	TotalNotes    int64 `json:"total_notes"`
	TotalLinks    int64 `json:"total_links"`
	Connected     int64 `json:"connected"`     // Notes with at least one link
	Orphans       int64 `json:"orphans"`        // Notes with no links
	MaxDepth      int   `json:"max_depth"`      // Longest shortest path
	AverageDegree float64 `json:"average_degree"` // Average links per note
}

// GraphNode represents a node in the knowledge graph
type GraphNode struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Type   NoteType  `json:"type"`
	TagIDs []string  `json:"tag_ids,omitempty"`
}

// GraphEdge represents an edge in the knowledge graph
type GraphEdge struct {
	Source    uuid.UUID `json:"source"`
	Target    uuid.UUID `json:"target"`
	Context   *string   `json:"context,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// GraphResponse represents the full knowledge graph
type GraphResponse struct {
	Nodes []*GraphNode `json:"nodes"`
	Edges []*GraphEdge `json:"edges"`
	Stats *GraphStats  `json:"stats,omitempty"`
}

// LinkDetail represents a link with source and target note details
type LinkDetail struct {
	ID           uuid.UUID  `json:"id"`
	SourceID     uuid.UUID  `json:"source_id"`
	TargetID     uuid.UUID  `json:"target_id"`
	LinkContext  *string    `json:"link_context,omitempty"`
	SourceNote   *Note      `json:"source_note,omitempty"`
	TargetNote   *Note      `json:"target_note,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

