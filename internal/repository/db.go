package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the pgxpool for database operations
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB creates a new database connection pool
func NewDB(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*DB, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	// Set pool configuration
	config.MaxConns = int32(maxOpenConns)
	config.MaxConnIdleTime = connMaxLifetime
	config.MaxConnLifetime = connMaxLifetime
	config.HealthCheckPeriod = 1 * time.Minute

	// Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Ping checks if the database is accessible
func (db *DB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Repository holds all repositories
type Repository struct {
	User          UserRepository
	Note          NoteRepository
	Tag           TagRepository
	Link          LinkRepository
	Activity      ActivityRepository
	RefreshToken  RefreshTokenRepository
}

// NewRepository creates a new repository with all sub-repositories
func NewRepository(db *DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Note:         NewNoteRepository(db),
		Tag:          NewTagRepository(db),
		Link:         NewLinkRepository(db),
		Activity:     NewActivityRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
	}
}
