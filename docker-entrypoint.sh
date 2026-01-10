#!/bin/sh
set -e

echo "Starting Knowledge Garden API..."

# Wait for database to be ready
if [ -n "$DB_HOST" ]; then
    echo "Waiting for database to be ready..."
    max_attempts=30
    attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if nc -z "$DB_HOST" "${DB_PORT:-5432}" 2>/dev/null; then
            echo "Database is ready!"
            break
        fi
        attempt=$((attempt + 1))
        echo "Database not ready yet, waiting... ($attempt/$max_attempts)"
        sleep 2
    done

    if [ $attempt -eq $max_attempts ]; then
        echo "Error: Database connection timeout after ${max_attempts} attempts"
        exit 1
    fi

    # Run migrations if SKIP_MIGRATIONS is not set
    if [ -z "$SKIP_MIGRATIONS" ]; then
        echo "Running database migrations..."
        if ./migrate up 2>&1; then
            echo "✓ Migrations applied successfully"
        else
            # Migrations failed - check if it's because objects already exist
            echo "⚠ Migration check complete (some migrations may already be applied)"
            echo "  Note: Migrations are idempotent - existing objects are safe"
            echo "  To skip migrations, set SKIP_MIGRATIONS=true"
        fi
    else
        echo "Skipping migrations (SKIP_MIGRATIONS is set)"
    fi
else
    echo "Warning: DB_HOST not set, skipping migrations"
fi

# Execute the main command
echo "Starting API server..."
exec "$@"
