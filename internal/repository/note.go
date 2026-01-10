package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/momokii/go-cli-notes/internal/model"
)

// NoteRepository handles note data operations
type NoteRepository struct {
	db *DB
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *DB) NoteRepository {
	return NoteRepository{db: db}
}

// Create inserts a new note
func (r *NoteRepository) Create(ctx context.Context, note *model.Note) error {
	query := `
		INSERT INTO notes (id, user_id, title, content, note_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, title, content, note_type, word_count, reading_time_minutes,
		          is_deleted, deleted_at, created_at, updated_at, last_accessed_at, access_count, metadata
	`

	now := time.Now()
	note.ID = uuid.New()
	note.CreatedAt = now
	note.UpdatedAt = now

	err := r.db.Pool.QueryRow(ctx, query,
		note.ID,
		note.UserID,
		note.Title,
		note.Content,
		note.NoteType,
		note.CreatedAt,
		note.UpdatedAt,
	).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.NoteType,
		&note.WordCount,
		&note.ReadingTimeMinutes,
		&note.IsDeleted,
		&note.DeletedAt,
		&note.CreatedAt,
		&note.UpdatedAt,
		&note.LastAccessedAt,
		&note.AccessCount,
		&note.Metadata,
	)

	if err != nil {
		return fmt.Errorf("create note: %w", err)
	}

	return nil
}

// FindByID finds a note by ID (with user scoping)
func (r *NoteRepository) FindByID(ctx context.Context, userID, id uuid.UUID) (*model.Note, error) {
	query := `
		SELECT id, user_id, title, content, note_type, word_count, reading_time_minutes,
		       is_deleted, deleted_at, created_at, updated_at, last_accessed_at, access_count, metadata
		FROM notes
		WHERE id = $1 AND user_id = $2 AND is_deleted = false
	`

	note := &model.Note{}
	err := r.db.Pool.QueryRow(ctx, query, id, userID).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.NoteType,
		&note.WordCount,
		&note.ReadingTimeMinutes,
		&note.IsDeleted,
		&note.DeletedAt,
		&note.CreatedAt,
		&note.UpdatedAt,
		&note.LastAccessedAt,
		&note.AccessCount,
		&note.Metadata,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find note by id: %w", err)
	}

	return note, nil
}

// FindByTitle finds a note by title and user
func (r *NoteRepository) FindByTitle(ctx context.Context, userID uuid.UUID, title string) (*model.Note, error) {
	query := `
		SELECT id, user_id, title, content, note_type, word_count, reading_time_minutes,
		       is_deleted, deleted_at, created_at, updated_at, last_accessed_at, access_count, metadata
		FROM notes
		WHERE user_id = $1 AND title = $2 AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT 1
	`

	note := &model.Note{}
	err := r.db.Pool.QueryRow(ctx, query, userID, title).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.NoteType,
		&note.WordCount,
		&note.ReadingTimeMinutes,
		&note.IsDeleted,
		&note.DeletedAt,
		&note.CreatedAt,
		&note.UpdatedAt,
		&note.LastAccessedAt,
		&note.AccessCount,
		&note.Metadata,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find note by title: %w", err)
	}

	return note, nil
}

// List lists notes for a user with pagination
func (r *NoteRepository) List(ctx context.Context, userID uuid.UUID, filter model.NoteFilter) ([]*model.Note, int64, error) {
	// Build the base query
	baseQuery := `
		SELECT id, user_id, title, content, note_type, word_count, reading_time_minutes,
		       is_deleted, deleted_at, created_at, updated_at, last_accessed_at, access_count, metadata
		FROM notes
		WHERE user_id = $1 AND is_deleted = false
	`

	countQuery := `
		SELECT COUNT(*)
		FROM notes
		WHERE user_id = $1 AND is_deleted = false
	`

	args := []any{userID}
	argPos := 2

	// Add filters
	if filter.NoteType != nil {
		baseQuery += fmt.Sprintf(" AND note_type = $%d", argPos)
		countQuery += fmt.Sprintf(" AND note_type = $%d", argPos)
		args = append(args, *filter.NoteType)
		argPos++
	}

	if filter.TagID != nil {
		baseQuery += fmt.Sprintf(" AND id IN (SELECT note_id FROM note_tags WHERE tag_id = $%d)", argPos)
		countQuery += fmt.Sprintf(" AND id IN (SELECT note_id FROM note_tags WHERE tag_id = $%d)", argPos)
		args = append(args, *filter.TagID)
		argPos++
	}

	if filter.Search != "" {
		baseQuery += fmt.Sprintf(" AND content_tsv @@ plainto_tsquery('english', $%d)", argPos)
		countQuery += fmt.Sprintf(" AND content_tsv @@ plainto_tsquery('english', $%d)", argPos)
		args = append(args, filter.Search)
		argPos++
	}

	// Get total count (use same args as base query, before pagination)
	var total int64
	countArgs := args
	countErr := r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if countErr != nil {
		return nil, 0, fmt.Errorf("count notes: %w", countErr)
	}

	// Add sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (filter.Page - 1) * limit
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	// Execute query
	rows, err := r.db.Pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list notes: %w", err)
	}
	defer rows.Close()

	notes := []*model.Note{}
	for rows.Next() {
		note := &model.Note{}
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.NoteType,
			&note.WordCount,
			&note.ReadingTimeMinutes,
			&note.IsDeleted,
			&note.DeletedAt,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.LastAccessedAt,
			&note.AccessCount,
			&note.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("iterate notes: %w", rows.Err())
	}

	return notes, total, nil
}

// Update updates a note
func (r *NoteRepository) Update(ctx context.Context, note *model.Note) error {
	query := `
		UPDATE notes
		SET title = COALESCE($1, title),
		    content = COALESCE($2, content),
		    updated_at = NOW()
		WHERE id = $3 AND user_id = $4 AND is_deleted = false
		RETURNING id, user_id, title, content, note_type, word_count, reading_time_minutes,
		          is_deleted, deleted_at, created_at, updated_at, last_accessed_at, access_count, metadata
	`

	err := r.db.Pool.QueryRow(ctx, query,
		note.Title,
		note.Content,
		note.ID,
		note.UserID,
	).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.NoteType,
		&note.WordCount,
		&note.ReadingTimeMinutes,
		&note.IsDeleted,
		&note.DeletedAt,
		&note.CreatedAt,
		&note.UpdatedAt,
		&note.LastAccessedAt,
		&note.AccessCount,
		&note.Metadata,
	)

	if err == pgx.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}

	return nil
}

// Delete soft deletes a note
func (r *NoteRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	query := `
		UPDATE notes
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_deleted = false
	`

	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Restore restores a soft deleted note
func (r *NoteRepository) Restore(ctx context.Context, userID, id uuid.UUID) error {
	query := `
		UPDATE notes
		SET is_deleted = false, deleted_at = NULL
		WHERE id = $1 AND user_id = $2 AND is_deleted = true
	`

	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("restore note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateAccessCount updates the access count and last accessed time
func (r *NoteRepository) UpdateAccessCount(ctx context.Context, userID, id uuid.UUID) error {
	query := `
		UPDATE notes
		SET access_count = access_count + 1,
		    last_accessed_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("update access count: %w", err)
	}

	return nil
}
