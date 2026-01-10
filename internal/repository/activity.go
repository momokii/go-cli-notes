package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/momokii/go-cli-notes/internal/model"
)

// ActivityRepository handles activity log operations
type ActivityRepository struct {
	db *DB
}

// NewActivityRepository creates a new activity repository
func NewActivityRepository(db *DB) ActivityRepository {
	return ActivityRepository{db: db}
}

// Create logs a new activity
func (r *ActivityRepository) Create(ctx context.Context, activity *model.Activity) error {
	query := `
		INSERT INTO activity_log (id, user_id, note_id, action, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, note_id, action, metadata, created_at
	`

	now := time.Now()
	activity.ID = uuid.New()
	activity.CreatedAt = now

	// Handle nil note_id
	var noteID pgtype.UUID
	if activity.NoteID != nil {
		noteID = pgtype.UUID{Bytes: *activity.NoteID, Valid: true}
	} else {
		noteID = pgtype.UUID{Valid: false}
	}

	err := r.db.Pool.QueryRow(ctx, query,
		activity.ID,
		activity.UserID,
		noteID,
		activity.Action,
		activity.Metadata,
		activity.CreatedAt,
	).Scan(
		&activity.ID,
		&activity.UserID,
		&noteID,
		&activity.Action,
		&activity.Metadata,
		&activity.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("create activity: %w", err)
	}

	// Convert back to uuid.UUID pointer
	if noteID.Valid {
		u := uuid.UUID(noteID.Bytes)
		activity.NoteID = &u
	}

	return nil
}

// GetRecent gets recent activities for a user
func (r *ActivityRepository) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*model.Activity, error) {
	query := `
		SELECT id, user_id, note_id, action, metadata, created_at
		FROM activity_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get recent activities: %w", err)
	}
	defer rows.Close()

	activities := []*model.Activity{}
	for rows.Next() {
		var noteID pgtype.UUID
		activity := &model.Activity{}
		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&noteID,
			&activity.Action,
			&activity.Metadata,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan activity: %w", err)
		}

		// Convert to uuid.UUID pointer
		if noteID.Valid {
			u := uuid.UUID(noteID.Bytes)
			activity.NoteID = &u
		}

		activities = append(activities, activity)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate activities: %w", rows.Err())
	}

	return activities, nil
}

// GetLastActivity gets the last activity timestamp for a user
func (r *ActivityRepository) GetLastActivity(ctx context.Context, userID uuid.UUID) (*time.Time, error) {
	query := `
		SELECT created_at
		FROM activity_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var createdAt time.Time
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&createdAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get last activity: %w", err)
	}

	return &createdAt, nil
}

// GetUserStats gets statistics for a user
func (r *ActivityRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*model.UserStats, error) {
	stats := &model.UserStats{}

	// Get total notes (non-deleted)
	err := r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notes WHERE user_id = $1 AND is_deleted = false
	`, userID).Scan(&stats.TotalNotes)
	if err != nil {
		return nil, fmt.Errorf("get total notes: %w", err)
	}

	// Get total tags
	err = r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM tags WHERE user_id = $1
	`, userID).Scan(&stats.TotalTags)
	if err != nil {
		return nil, fmt.Errorf("get total tags: %w", err)
	}

	// Get total links
	err = r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM links WHERE user_id = $1
	`, userID).Scan(&stats.TotalLinks)
	if err != nil {
		return nil, fmt.Errorf("get total links: %w", err)
	}

	// Get total words
	err = r.db.Pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(word_count), 0) FROM notes WHERE user_id = $1 AND is_deleted = false
	`, userID).Scan(&stats.TotalWords)
	if err != nil {
		return nil, fmt.Errorf("get total words: %w", err)
	}

	// Get notes created today
	err = r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notes
		WHERE user_id = $1 AND is_deleted = false
		AND DATE(created_at) = CURRENT_DATE
	`, userID).Scan(&stats.NotesCreatedToday)
	if err != nil {
		return nil, fmt.Errorf("get notes created today: %w", err)
	}

	// Get notes created this week
	err = r.db.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notes
		WHERE user_id = $1 AND is_deleted = false
		AND created_at >= DATE_TRUNC('week', CURRENT_DATE)
	`, userID).Scan(&stats.NotesCreatedWeek)
	if err != nil {
		return nil, fmt.Errorf("get notes created this week: %w", err)
	}

	// Get last activity
	stats.LastActivity, _ = r.GetLastActivity(ctx, userID)

	return stats, nil
}

// GetTrendingNotes gets frequently accessed notes
func (r *ActivityRepository) GetTrendingNotes(ctx context.Context, userID uuid.UUID, limit int) ([]*model.TrendingNote, error) {
	query := `
		SELECT id, user_id, title, access_count, last_accessed_at
		FROM notes
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY access_count DESC, last_accessed_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get trending notes: %w", err)
	}
	defer rows.Close()

	trending := []*model.TrendingNote{}
	for rows.Next() {
		note := &model.Note{}
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.AccessCount,
			&note.LastAccessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan trending note: %w", err)
		}

		trending = append(trending, &model.TrendingNote{
			Note:        note,
			AccessCount: note.AccessCount,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate trending notes: %w", rows.Err())
	}

	return trending, nil
}

// GetForgottenNotes gets notes that haven't been accessed in a while
func (r *ActivityRepository) GetForgottenNotes(ctx context.Context, userID uuid.UUID, days int, limit int) ([]*model.ForgottenNote, error) {
	query := `
		SELECT id, user_id, title, COALESCE(last_accessed_at, created_at) as last_accessed
		FROM notes
		WHERE user_id = $1 AND is_deleted = false
		AND (last_accessed_at < NOW() - INTERVAL '1 day' * $2 OR last_accessed_at IS NULL)
		ORDER BY last_accessed ASC
		LIMIT $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, days, limit)
	if err != nil {
		return nil, fmt.Errorf("get forgotten notes: %w", err)
	}
	defer rows.Close()

	forgotten := []*model.ForgottenNote{}
	for rows.Next() {
		note := &model.Note{}
		var lastAccessed time.Time
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&lastAccessed,
		)
		if err != nil {
			return nil, fmt.Errorf("scan forgotten note: %w", err)
		}

		daysSince := int(time.Since(lastAccessed).Hours() / 24)

		forgotten = append(forgotten, &model.ForgottenNote{
			Note:            note,
			LastAccessedAt:  lastAccessed,
			DaysSinceAccess: daysSince,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate forgotten notes: %w", rows.Err())
	}

	return forgotten, nil
}
