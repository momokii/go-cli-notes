package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/momokii/go-cli-notes/internal/model"
)

// TagRepository handles tag data operations
type TagRepository struct {
	db *DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *DB) TagRepository {
	return TagRepository{db: db}
}

// Create inserts a new tag
func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	query := `
		INSERT INTO tags (id, user_id, name, color, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, color, created_at
	`

	now := time.Now()
	tag.ID = uuid.New()
	tag.CreatedAt = now

	err := r.db.Pool.QueryRow(ctx, query,
		tag.ID,
		tag.UserID,
		tag.Name,
		tag.Color,
		tag.CreatedAt,
	).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Name,
		&tag.Color,
		&tag.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}

	return nil
}

// FindByID finds a tag by ID
func (r *TagRepository) FindByID(ctx context.Context, userID, id uuid.UUID) (*model.Tag, error) {
	query := `
		SELECT id, user_id, name, color, created_at
		FROM tags
		WHERE id = $1 AND user_id = $2
	`

	tag := &model.Tag{}
	err := r.db.Pool.QueryRow(ctx, query, id, userID).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Name,
		&tag.Color,
		&tag.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find tag by id: %w", err)
	}

	return tag, nil
}

// FindByName finds a tag by name for a user
func (r *TagRepository) FindByName(ctx context.Context, userID uuid.UUID, name string) (*model.Tag, error) {
	query := `
		SELECT id, user_id, name, color, created_at
		FROM tags
		WHERE user_id = $1 AND name = $2
	`

	tag := &model.Tag{}
	err := r.db.Pool.QueryRow(ctx, query, userID, name).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Name,
		&tag.Color,
		&tag.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find tag by name: %w", err)
	}

	return tag, nil
}

// List lists all tags for a user
func (r *TagRepository) List(ctx context.Context, userID uuid.UUID) ([]*model.Tag, error) {
	query := `
		SELECT id, user_id, name, color, created_at
		FROM tags
		WHERE user_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	tags := []*model.Tag{}
	for rows.Next() {
		tag := &model.Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.UserID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate tags: %w", rows.Err())
	}

	return tags, nil
}

// ListWithNoteCount lists tags with note count for a user
func (r *TagRepository) ListWithNoteCount(ctx context.Context, userID uuid.UUID, page, limit int) ([]*model.TagWithCount, int64, error) {
	// Get total count
	var total int64
	countErr := r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM tags WHERE user_id = $1", userID).Scan(&total)
	if countErr != nil {
		return nil, 0, fmt.Errorf("count tags: %w", countErr)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get tags with note count
	query := `
		SELECT t.id, t.user_id, t.name, t.color, t.created_at, COUNT(nt.note_id) as note_count
		FROM tags t
		LEFT JOIN note_tags nt ON t.id = nt.tag_id
		WHERE t.user_id = $1
		GROUP BY t.id, t.user_id, t.name, t.color, t.created_at
		ORDER BY t.name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list tags with count: %w", err)
	}
	defer rows.Close()

	tags := []*model.TagWithCount{}
	for rows.Next() {
		tag := &model.TagWithCount{}
		err := rows.Scan(
			&tag.ID,
			&tag.UserID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
			&tag.NoteCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan tag with count: %w", err)
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("iterate tags with count: %w", rows.Err())
	}

	return tags, total, nil
}

// Update updates a tag
func (r *TagRepository) Update(ctx context.Context, tag *model.Tag) error {
	query := `
		UPDATE tags
		SET name = COALESCE($1, name),
		    color = COALESCE($2, color)
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, color, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		tag.Name,
		tag.Color,
		tag.ID,
		tag.UserID,
	).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Name,
		&tag.Color,
		&tag.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update tag: %w", err)
	}

	return nil
}

// Delete deletes a tag
func (r *TagRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	query := `DELETE FROM tags WHERE id = $1 AND user_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// AddToNote adds a tag to a note
func (r *TagRepository) AddToNote(ctx context.Context, noteID, tagID uuid.UUID) error {
	query := `
		INSERT INTO note_tags (note_id, tag_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (note_id, tag_id) DO NOTHING
	`

	_, err := r.db.Pool.Exec(ctx, query, noteID, tagID)
	if err != nil {
		return fmt.Errorf("add tag to note: %w", err)
	}

	return nil
}

// RemoveFromNote removes a tag from a note
func (r *TagRepository) RemoveFromNote(ctx context.Context, noteID, tagID uuid.UUID) error {
	query := `DELETE FROM note_tags WHERE note_id = $1 AND tag_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, noteID, tagID)
	if err != nil {
		return fmt.Errorf("remove tag from note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// GetByNote gets all tags for a note
func (r *TagRepository) GetByNote(ctx context.Context, noteID uuid.UUID) ([]*model.Tag, error) {
	query := `
		SELECT t.id, t.user_id, t.name, t.color, t.created_at
		FROM tags t
		INNER JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = $1
		ORDER BY t.name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, noteID)
	if err != nil {
		return nil, fmt.Errorf("get tags by note: %w", err)
	}
	defer rows.Close()

	tags := []*model.Tag{}
	for rows.Next() {
		tag := &model.Tag{}
		err := rows.Scan(
			&tag.ID,
			&tag.UserID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate tags: %w", rows.Err())
	}

	return tags, nil
}
