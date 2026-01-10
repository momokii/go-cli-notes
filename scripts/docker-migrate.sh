#!/bin/bash
# Docker migration helper script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
COMMAND=${1:-"help"}
CONTAINER_NAME=${CONTAINER_NAME:-"kg-api"}

# Print usage
usage() {
    cat << EOF
Usage: ./scripts/docker-migrate.sh <command>

Commands:
    up                   Apply all pending migrations
    down                 Rollback the most recent migration
    status               Show migration status
    redo                 Rollback and re-apply the most recent migration
    create NAME [type]   Create a new migration file (default type: sql)
    help                 Show this help message

Environment Variables:
    CONTAINER_NAME       Container name (default: go-cli-notes-api)

Examples:
    ./scripts/docker-migrate.sh status
    ./scripts/docker-migrate.sh up
    ./scripts/docker-migrate.sh down

Note: This script requires the API container to be running.

EOF
}

# Check if container is running
check_container() {
    if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo -e "${RED}Error: Container '${CONTAINER_NAME}' is not running${NC}"
        echo ""
        echo "Start the container first:"
        echo "  docker compose up -d api"
        echo ""
        echo "Or set CONTAINER_NAME if using a different name:"
        echo "  CONTAINER_NAME=my-api ./scripts/docker-migrate.sh status"
        exit 1
    fi
}

# Execute command in container
exec_in_container() {
    docker exec "${CONTAINER_NAME}" ./migrate "$@"
}

# Execute command
case $COMMAND in
    up)
        check_container
        echo -e "${GREEN}Running migrations in container...${NC}"
        exec_in_container up
        ;;
    down)
        check_container
        echo -e "${YELLOW}Rolling back last migration in container...${NC}"
        exec_in_container down
        ;;
    status)
        check_container
        echo -e "${GREEN}Migration status in container:${NC}"
        exec_in_container status
        ;;
    redo)
        check_container
        echo -e "${YELLOW}Redoing last migration in container...${NC}"
        exec_in_container redo
        ;;
    create)
        echo -e "${RED}Error: Cannot create migrations inside Docker container${NC}"
        echo ""
        echo "Use local migration script instead:"
        echo "  ./scripts/migrate.sh create $2 ${3:-sql}"
        exit 1
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
