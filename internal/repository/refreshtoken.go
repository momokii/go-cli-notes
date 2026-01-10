package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/momokii/go-cli-notes/internal/model"
)

// RefreshTokenRepository handles refresh token data operations
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// refreshTokenRepository implements RefreshTokenRepository
type refreshTokenRepository struct {
	db *DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create inserts a new refresh token
func (r *refreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, token_hash, expires_at, created_at, is_revoked
	`

	now := time.Now()
	token.ID = uuid.New()
	token.CreatedAt = now

	err := r.db.Pool.QueryRow(ctx, query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.IsRevoked,
	)

	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}

	return nil
}

// FindByTokenHash finds a refresh token by its hash
func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, is_revoked
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	token := &model.RefreshToken{}
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.IsRevoked,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find refresh token: %w", err)
	}

	return token, nil
}

// Revoke marks a refresh token as revoked
func (r *refreshTokenRepository) Revoke(ctx context.Context, tokenID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE user_id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens: %w", err)
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens
func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := r.db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}

	return nil
}
