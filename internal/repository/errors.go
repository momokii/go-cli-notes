package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")
)

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	return err == ErrNotFound || err == pgx.ErrNoRows
}
