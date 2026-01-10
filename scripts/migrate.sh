#!/bin/bash
# Migration helper script for local development

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
COMMAND=${1:-"help"}
DIR="./migrations"

# Print usage
usage() {
    cat << EOF
Usage: ./scripts/migrate.sh <command>

Commands:
    up                   Apply all pending migrations
    down                 Rollback the most recent migration
    status               Show migration status
    redo                 Rollback and re-apply the most recent migration
    create NAME [type]   Create a new migration file (default type: sql)
    bootstrap            Mark existing migrations as applied (for databases set up before migration system)
    help                 Show this help message

Environment Variables:
    DB_HOST             Database host (default: localhost)
    DB_PORT             Database port (default: 5432)
    DB_USER             Database user (default: kg_user)
    DB_PASSWORD         Database password (required)
    DB_NAME             Database name (default: knowledge_garden)
    DB_SSL_MODE         SSL mode (default: disable)

Examples:
    ./scripts/migrate.sh status
    ./scripts/migrate.sh up
    ./scripts/migrate.sh down
    ./scripts/migrate.sh bootstrap
    ./scripts/migrate.sh create add_users_table sql

EOF
}

# Check if .env file exists and load it
if [ -f .env ]; then
    echo -e "${YELLOW}Loading .env file...${NC}"
    set -a
    # Use '.' instead of 'source' for POSIX compatibility (works with sh, bash, dash, etc.)
    . ./.env
    set +a
fi

# Check if DB_PASSWORD is set
if [ -z "$DB_PASSWORD" ] && [ "$COMMAND" != "help" ] && [ "$COMMAND" != "create" ]; then
    echo -e "${RED}Error: DB_PASSWORD environment variable is required${NC}"
    echo ""
    echo "You can either:"
    echo "  1. Export DB_PASSWORD: export DB_PASSWORD=your_password"
    echo "  2. Create a .env file with DB_PASSWORD=your_password"
    exit 1
fi

# Build the migrate binary first
echo -e "${YELLOW}Building migrate binary...${NC}"
cd "$(dirname "$0")/.."
go build -o /tmp/kg-migrate ./migrations

# Execute command
case $COMMAND in
    up)
        echo -e "${GREEN}Running migrations...${NC}"
        /tmp/kg-migrate up
        ;;
    down)
        echo -e "${YELLOW}Rolling back last migration...${NC}"
        /tmp/kg-migrate down
        ;;
    status)
        echo -e "${GREEN}Migration status:${NC}"
        /tmp/kg-migrate status
        ;;
    redo)
        echo -e "${YELLOW}Redoing last migration...${NC}"
        /tmp/kg-migrate redo
        ;;
    bootstrap)
        echo -e "${YELLOW}Bootstrapping migration tracking...${NC}"
        echo -e "${YELLOW}This will mark all migrations as applied without running SQL.${NC}"
        /tmp/kg-migrate bootstrap
        ;;
    create)
        if [ -z "$2" ]; then
            echo -e "${RED}Error: Migration name is required${NC}"
            echo "Usage: ./scripts/migrate.sh create <name> [type]"
            exit 1
        fi
        MIGRATION_NAME=$2
        MIGRATION_TYPE=${3:-sql}
        echo -e "${GREEN}Creating migration: ${MIGRATION_NAME} (type: ${MIGRATION_TYPE})${NC}"
        /tmp/kg-migrate create "$MIGRATION_NAME" "$MIGRATION_TYPE"
        ;;
    help|--help|-h)
        usage
        ;;
    *)
        echo -e "${RED}Error: Unknown command '$COMMAND'${NC}"
        echo ""
        usage
        exit 1
        ;;
esac
