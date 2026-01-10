package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/momokii/go-cli-notes/internal/model"
)

// LinkRepository handles link data operations
type LinkRepository struct {
	db *DB
}

// NewLinkRepository creates a new link repository
func NewLinkRepository(db *DB) LinkRepository {
	return LinkRepository{db: db}
}

// Create inserts a new link
func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	query := `
		INSERT INTO links (id, user_id, source_note_id, target_note_id, link_context, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (source_note_id, target_note_id) DO NOTHING
		RETURNING id, user_id, source_note_id, target_note_id, link_context, created_at
	`

	now := time.Now()
	link.ID = uuid.New()
	link.CreatedAt = now

	err := r.db.Pool.QueryRow(ctx, query,
		link.ID,
		link.UserID,
		link.SourceNoteID,
		link.TargetNoteID,
		link.LinkContext,
		link.CreatedAt,
	).Scan(
		&link.ID,
		&link.UserID,
		&link.SourceNoteID,
		&link.TargetNoteID,
		&link.LinkContext,
		&link.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		// Link already exists
		return nil
	}
	if err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

// GetBySource gets all outgoing links from a note
func (r *LinkRepository) GetBySource(ctx context.Context, userID, noteID uuid.UUID) ([]*model.Link, error) {
	query := `
		SELECT l.id, l.user_id, l.source_note_id, l.target_note_id, l.link_context, l.created_at
		FROM links l
		WHERE l.user_id = $1 AND l.source_note_id = $2
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("get links by source: %w", err)
	}
	defer rows.Close()

	links := []*model.Link{}
	for rows.Next() {
		link := &model.Link{}
		err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.SourceNoteID,
			&link.TargetNoteID,
			&link.LinkContext,
			&link.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan link: %w", err)
		}
		links = append(links, link)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate links: %w", rows.Err())
	}

	return links, nil
}

// GetByTarget gets all incoming links to a note (backlinks)
func (r *LinkRepository) GetByTarget(ctx context.Context, userID, noteID uuid.UUID) ([]*model.Link, error) {
	query := `
		SELECT l.id, l.user_id, l.source_note_id, l.target_note_id, l.link_context, l.created_at
		FROM links l
		WHERE l.user_id = $1 AND l.target_note_id = $2
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, noteID)
	if err != nil {
		return nil, fmt.Errorf("get links by target: %w", err)
	}
	defer rows.Close()

	links := []*model.Link{}
	for rows.Next() {
		link := &model.Link{}
		err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.SourceNoteID,
			&link.TargetNoteID,
			&link.LinkContext,
			&link.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan link: %w", err)
		}
		links = append(links, link)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate links: %w", rows.Err())
	}

	return links, nil
}

// Delete deletes a link
func (r *LinkRepository) Delete(ctx context.Context, userID, sourceID, targetID uuid.UUID) error {
	query := `
		DELETE FROM links
		WHERE user_id = $1 AND source_note_id = $2 AND target_note_id = $3
	`

	result, err := r.db.Pool.Exec(ctx, query, userID, sourceID, targetID)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteByNote deletes all links associated with a note (both incoming and outgoing)
func (r *LinkRepository) DeleteByNote(ctx context.Context, userID, noteID uuid.UUID) error {
	query := `
		DELETE FROM links
		WHERE user_id = $1 AND (source_note_id = $2 OR target_note_id = $2)
	`

	_, err := r.db.Pool.Exec(ctx, query, userID, noteID)
	if err != nil {
		return fmt.Errorf("delete links by note: %w", err)
	}

	return nil
}
