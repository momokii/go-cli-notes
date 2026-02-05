# Knowledge Garden CLI

A terminal-based Personal Knowledge Management (PKM) system with wiki-style links, full-text search, and knowledge graph visualization.

## Features

- **Note Management**: Create, edit, delete, and organize notes
- **Wiki-Style Links**: Connect notes using `[[Note Title]]` syntax
- **Full-Text Search**: Fast PostgreSQL-based full-text search
- **Knowledge Graph**: Visualize connections between your notes
- **Tags**: Organize notes with tags for easy filtering
- **Daily Notes**: Automatic daily journal entries
- **Analytics**: Track your writing habits and activity
- **CLI & API**: Use via command-line or REST API

## Architecture

```
┌─────────────┐
│   CLI/UI    │  ← Terminal Interface or REST API
└──────┬──────┘
       │
┌──────▼──────┐
│ REST API    │  ← Fiber v2 HTTP Server
└──────┬──────┘
       │
┌──────▼──────┐
│  Services   │  ← Business Logic Layer
└──────┬──────┘
       │
┌──────▼──────┐
│Repositories │  ← Data Access Layer (pgx/v5)
└──────┬──────┘
       │
┌──────▼──────┐
│ PostgreSQL  │  ← Database with FTS & JSONB
└─────────────┘
```

## Quick Start

### Installation

The easiest way to install `kg-cli` is using the web-based installer:

```bash
curl -fsSL https://raw.githubusercontent.com/momokii/go-cli-notes/main/scripts/install.sh | bash
```

This will:
- Automatically detect your platform (Linux, macOS, Windows)
- Download the pre-built binary from GitHub Releases
- Install to `~/.local/bin` (no sudo required) or `/usr/local/bin`
- Set up your API configuration (local, cloud, or custom URL)

**Supported Platforms:**
- Linux (amd64, arm64)
- macOS (Intel, Apple Silicon)
- Windows (amd64)

For detailed installation instructions, troubleshooting, and alternative installation methods, see [docs/INSTALL.md](docs/INSTALL.md).

### Development Setup

#### Prerequisites

- **Go**: 1.24 or later
- **PostgreSQL**: 15 or later
- **Docker** & **Docker Compose** (optional, for containerized setup)

#### Option 1: Using Docker (Recommended for Development)

```bash
# Clone the repository
git clone https://github.com/momokii/go-cli-notes.git
cd go-cli-notes

# Start the API and database
docker compose up -d

# The API will be available at http://localhost:8080
```

#### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/momokii/go-cli-notes.git
cd go-cli-notes

# Install dependencies
go mod download

# Build the CLI
go build -o kg-cli ./cmd/cli

# (Optional) Build and run the API
go build -o api ./cmd/api
./api
```

### First-Time Setup

```bash
# 1. Register a new account
./kg-cli register
# Prompts for: Username, Email, Password

# 2. Login
./kg-cli login
# Prompts for: Email, Password

# 3. Check status
./kg-cli status
```

## CLI Usage

### Global Commands

```bash
# Show help
./kg-cli --help
./kg-cli [command] --help

# Authentication
./kg-cli register          # Register a new account
./kg-cli login             # Login to your account
./kg-cli logout            # Logout from your account
./kg-cli status            # Show authentication and connection status
```

### Note Management

```bash
# List all notes
./kg-cli note list

# List notes with pagination
./kg-cli note list --page 1 --limit 10

# Search notes
./kg-cli note search "golang"

# Create a new note
./kg-cli note create --title "My First Note" --content "This is the content"

# Create a note with specific type
./kg-cli note create --title "Meeting Notes" --content "Project discussion" --type meeting

# Get a specific note
./kg-cli note get <note-id>

# Get or create today's daily note
./kg-cli note daily

# Get or create a daily note for a specific date
./kg-cli note daily 2026-01-04
```

### Note Types

- `note` - Regular notes (default)
- `daily` - Daily journal entries
- `meeting` - Meeting notes
- `idea` - Quick ideas and thoughts

### Wiki-Style Links

Connect notes using double brackets in your note content:

```markdown
This is related to [[Another Note]].

You can also reference [[Go Programming]] concepts.
```

**How Links Work:**

Links are created **automatically** when you use the `[[Note Title]]` syntax in your note content:

1. **Target note must exist first** - Links are only created if the referenced note already exists
2. **Links are created on save** - When you save or update a note, the system parses the content and creates links to any existing notes
3. **Bidirectional** - Links work both ways (see "links" and "backlinks" commands below)

**Example workflow:**

```bash
# 1. Create the first note
./kg-cli note create --title "Go Tips" --content "Useful Go programming tips"
NOTE_A_ID=<output-id>

# 2. Create a second note that references the first
./kg-cli note create --title "My Learning Journey" \
  --content "Today I learned about Go. See [[Go Tips]] for more."
NOTE_B_ID=<output-id>

# 3. View the links from Note B (should show "Go Tips")
./kg-cli note links $NOTE_B_ID

# 4. View backlinks to Note A (should show "My Learning Journey")
./kg-cli note backlinks $NOTE_A_ID
```

**Note:** If you create a note with `[[Some Note]]` but "Some Note" doesn't exist yet, no link will be created. You can create the target note later and then update your original note to create the link.

### Tag Management

```bash
# List all tags
./kg-cli tag list

# Create a new tag
./kg-cli tag create "programming"

# Update a tag
./kg-cli tag update <tag-id> "new-name"

# Delete a tag
./kg-cli tag delete <tag-id>

# Add a tag to a note (by name or ID)
./kg-cli tag add "programming" <note-id>

# Remove a tag from a note
./kg-cli tag remove "programming" <note-id>

# View tags on a specific note
./kg-cli note tags <note-id>

# List notes with a specific tag
./kg-cli note list --tag "programming"
```

### Getting Started: Tags and Links Workflow

Here's a practical example of how to use tags and links together to build your knowledge garden:

```bash
# 1. Create some tags for organization
./kg-cli tag create "project"
./kg-cli tag create "research"
./kg-cli tag create "tutorial"

# 2. Create a project note
./kg-cli note create --title "Build CLI Tool" \
  --content "Building a CLI tool for personal knowledge management"
PROJECT_ID=<output-id>

# 3. Tag the project note
./kg-cli tag add "project" $PROJECT_ID

# 4. Create research notes and link them to the project
./kg-cli note create --title "Cobra Framework Research" \
  --content "Cobra is a library for creating CLI applications. Will use for [[Build CLI Tool]]."
RESEARCH_ID=<output-id>
./kg-cli tag add "research" $RESEARCH_ID

# 5. Create a tutorial note with examples
./kg-cli note create --title "CLI Tutorial" \
  --content "Step-by-step tutorial for building a CLI. See [[Cobra Framework Research]] for details."
TUTORIAL_ID=<output-id>
./kg-cli tag add "tutorial" $TUTORIAL_ID

# 6. View all notes in your project
./kg-cli note list --tag "project"

# 7. See what references your research note (backlinks)
./kg-cli note backlinks $RESEARCH_ID

# 8. View tags on any note
./kg-cli note tags $PROJECT_ID
```

**Key Concepts:**
- **Tags** = Categories and organization (e.g., "project", "research")
- **Links** = Connections and relationships between specific notes
- Use tags to group notes by type or topic
- Use links to connect related ideas and create knowledge pathways

For comprehensive workflows and best practices, see [WORKFLOWS.md](WORKFLOWS.md).

### Search

```bash
# Basic search
./kg-cli note search "golang"

# Search with pagination
./kg-cli note search "golang" --page 1 --limit 20
```

### Analytics & Statistics

```bash
# Show user statistics
./kg-cli stats

# Show recent activity
./kg-cli activity

# Show recent activity (limit 20)
./kg-cli activity --limit 20

# Show trending notes
./kg-cli trending

# Show trending notes (limit 10)
./kg-cli trending --limit 10
```

### Terminal User Interface (TUI)

The Knowledge Garden CLI includes an interactive Terminal User Interface (TUI) for a rich, visual experience.

```bash
# Launch the TUI
./kg-cli tui
```

**TUI Features:**
- **Dashboard**: Overview of your knowledge garden with statistics and recent activity
- **Note Browser**: Browse, search, and view notes with vim-style navigation
- **Note Editor**: Create and edit notes directly in the terminal
- **Tag Manager**: Create, edit, and delete tags
- **Search**: Full-text search with result highlighting
- **Activity Feed**: View your recent actions
- **Knowledge Graph**: ASCII visualization of note connections

**Key Bindings:**
- `?` - Show help
- `n` - New note
- `/` - Search
- `t` - Tags
- `a` - Activity
- `g` - Knowledge graph
- `j`/`k` - Navigate up/down
- `Enter` - Open/Select
- `ESC` - Go back
- `q` - Quit

**Requirements:**
- Terminal size: 80x24 minimum
- Valid authentication session (run `kg-cli login` first)

For detailed documentation, see [docs/TUI_USER_GUIDE.md](docs/TUI_USER_GUIDE.md)

### Configuration

The CLI stores configuration in `~/.config/kg-cli/config.yaml`:

```yaml
api:
  base_url: "http://localhost:8080"
  timeout: 30

editor:
  external_editor: "vim"

preferences:
  default_note_type: "note"
  auto_save_interval: 30
  theme: "dark"
```

### Environment Variables

You can override configuration using environment variables:

```bash
export KG_CLI_API_BASE_URL="http://localhost:8080"
export KG_CLI_API_TIMEOUT="30"
export KG_CLI_EDITOR="vim"
```

## REST API

### Authentication

#### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","email":"user@example.com","password":"SecurePass123"}'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"SecurePass123"}'
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "bearer",
  "expires_in": 3600
}
```

### Notes API

#### List Notes
```bash
curl http://localhost:8080/api/v1/notes \
  -H "Authorization: Bearer <access_token>"
```

#### Create Note
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Note",
    "content": "Note content here",
    "note_type": "note"
  }'
```

#### Get Note
```bash
curl http://localhost:8080/api/v1/notes/<note-id> \
  -H "Authorization: Bearer <access_token>"
```

#### Update Note
```bash
curl -X PUT http://localhost:8080/api/v1/notes/<note-id> \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content"
  }'
```

#### Delete Note
```bash
curl -X DELETE http://localhost:8080/api/v1/notes/<note-id> \
  -H "Authorization: Bearer <access_token>"
```

### Search API

```bash
curl "http://localhost:8080/api/v1/search?q=golang&page=1&limit=20" \
  -H "Authorization: Bearer <access_token>"
```

### Knowledge Graph API

```bash
curl http://localhost:8080/api/v1/notes/graph \
  -H "Authorization: Bearer <access_token>"
```

Response:
```json
{
  "nodes": [
    {
      "id": "uuid",
      "title": "Note Title",
      "type": "note"
    }
  ],
  "edges": [
    {
      "source": "uuid-1",
      "target": "uuid-2",
      "context": "link context text",
      "created_at": "2026-01-04T12:00:00Z"
    }
  ]
}
```

### Tags API

#### List Tags
```bash
curl http://localhost:8080/api/v1/tags \
  -H "Authorization: Bearer <access_token>"
```

#### Create Tag
```bash
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "programming"}'
```

### Analytics API

#### User Statistics
```bash
curl http://localhost:8080/api/v1/stats \
  -H "Authorization: Bearer <access_token>"
```

#### Recent Activity
```bash
curl "http://localhost:8080/api/v1/activity/recent?limit=10" \
  -H "Authorization: Bearer <access_token>"
```

#### Trending Notes
```bash
curl "http://localhost:8080/api/v1/notes/trending?limit=5" \
  -H "Authorization: Bearer <access_token>"
```

#### Forgotten Notes
```bash
curl "http://localhost:8080/api/v1/notes/forgotten?days=30&limit=5" \
  -H "Authorization: Bearer <access_token>"
```

## Development

### Project Structure

```
go-cli-notes/
├── cmd/
│   ├── api/                # REST API server
│   │   └── main.go
│   └── cli/               # CLI application
│       ├── main.go         # Entry point
│       ├── config.go       # Configuration management
│       ├── client/         # HTTP client for API
│       │   ├── api.go
│       │   └── auth.go
│       ├── note.go         # Note commands
│       ├── tag.go          # Tag commands
│       └── stats.go        # Stats commands
├── internal/
│   ├── api/
│   │   ├── handler/        # HTTP request handlers
│   │   ├── middleware/     # Middleware (auth, logger, etc.)
│   │   └── router/         # Route definitions
│   ├── config/            # Configuration structs
│   ├── model/             # Data models
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic
│   └── util/              # Utilities (JWT, password, etc.)
├── migrations/            # Database migrations
├── docker-compose.yml     # Docker services
├── Dockerfile.api         # API container image
└── go.mod               # Go module definition
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migration management. Migrations are tracked in the `goose_db_version` table and can be run automatically or manually.

**Important:**
- All migration methods automatically load environment variables from the `.env` file if it exists
- All migrations are **idempotent** - safe to re-run even if tables already exist
- Migrations use `IF NOT EXISTS` / `IF EXISTS` to prevent errors

#### For Existing Databases

If you have an existing database (created before the migration system was implemented), use the **bootstrap** command to mark all migrations as applied without running the SQL:

```bash
go run ./migrations bootstrap
# or
./scripts/migrate.sh bootstrap
```

This tells the migration system that your database is already up-to-date, so future migrations will run correctly.

#### Local Development

```bash
# Method 1: Using the helper script (recommended)
# Works with bash, sh, dash, zsh, etc.
./scripts/migrate.sh status              # Check migration status
./scripts/migrate.sh up                  # Apply all pending migrations
./scripts/migrate.sh down                # Rollback the most recent migration
./scripts/migrate.sh redo                # Rollback and re-apply the most recent migration
./scripts/migrate.sh bootstrap           # Mark existing migrations as applied

# Method 2: Using the Makefile
make migrate-status
make migrate-up
make migrate-down
make migrate-redo

# Method 3: Using Go directly (also auto-loads .env)
go run ./migrations status
go run ./migrations up
go run ./migrations down
go run ./migrations bootstrap           # For existing databases

# Create a new migration
./scripts/migrate.sh create add_users_table sql
```

#### Docker / Production

```bash
# 1. Run migrations first (before starting API)
./scripts/migrate.sh up

# 2. Start services
docker compose up -d

# Run migrations manually
./scripts/docker-migrate.sh status
./scripts/docker-migrate.sh up
./scripts/docker-migrate.sh down
```

**Note:** Migrations must be run manually before starting the API:
- API container no longer runs migrations automatically
- Use `./scripts/migrate.sh up` to apply migrations
- Use `./scripts/migrate.sh bootstrap` for existing databases
- Migrations are idempotent - safe to re-run

**Environment Variables for Migrations:**
- `DATABASE_URL` - Full database connection URL (for NeonDB, Supabase, etc.) - takes precedence
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: kg_user)
- `DB_PASSWORD` - Database password (required)
- `DB_NAME` - Database name (default: knowledge_garden)
- `DB_SSL_MODE` - SSL mode (default: disable)

#### Migration Files

Migration files are stored in the `migrations/` directory with the naming convention:
- `YYYYMMDDHHMMSS_description.sql` - Single file with up/down migrations

Each migration file uses Goose's SQL format with `-- +goose Up` and `-- +goose Down` markers:

```sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
```

Example:
- `20250110120000_init_schema.sql`
- `20250110120001_add_fts.sql`

## Deployment

### Building for Production

```bash
# Build API binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Build CLI binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kg-cli ./cmd/cli
```

### Docker Deployment

```bash
# Build and start all services
docker compose up -d

# View logs
docker compose logs -f api

# Stop services
docker compose down
```

### Configuration

Set environment variables for production:

```bash
# Database - Option 1: Full database URL (recommended for NeonDB, Supabase, etc.)
export DATABASE_URL="postgresql://user:password@host:port/dbname?sslmode=require"

# Database - Option 2: Individual variables (for local PostgreSQL)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=kg_user
export DB_PASSWORD=secure_password
export DB_NAME=kg_db

# JWT
export JWT_SECRET=your-secret-key-here
export JWT_ACCESS_EXPIRATION=3600
export JWT_REFRESH_EXPIRATION=86400

# Server
export SERVER_ADDRESS=0.0.0.0:8080
export SERVER_READ_TIMEOUT=30s
export SERVER_WRITE_TIMEOUT=30s
```

**Note:** If `DATABASE_URL` is set, it takes precedence over individual `DB_*` variables.

## Troubleshooting

### Common Issues

**1. "connection refused" error**
- Ensure PostgreSQL is running
- Check the database connection string in config

**2. NeonDB / Cloud Database Connection**

For NeonDB, Supabase, or other cloud PostgreSQL services, use the `DATABASE_URL` environment variable:

```bash
# In .env file
DATABASE_URL=postgresql://neondb_owner:password@ep-xxx-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
```

Or in docker-compose.yml:
```yaml
services:
  api:
    environment:
      DATABASE_URL: ${DATABASE_URL}
```

**Common NeonDB issues:**
- Make sure to use the "pooler" connection string (not direct host)
- Include `sslmode=require` in the connection string
- The connection string should start with `postgresql://` not `postgres://`

**2. "unauthorized" errors**
- Login again using `./kg-cli login`
- Check that your access token hasn't expired

**3. Build errors**
- Run `go mod tidy` to update dependencies
- Ensure you have Go 1.24 or later installed

### Getting Help

```bash
# Show help for all commands
./kg-cli --help

# Show help for specific command
./kg-cli note --help
./kg-cli tag --help
./kg-cli stats --help
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Roadmap

- [ ] Bubbletea TUI for interactive terminal interface
- [ ] Export notes to Markdown, PDF, HTML
- [ ] Import from Markdown files
- [ ] Offline mode with local caching
- [ ] Mobile app (via SSH)
- [ ] Collaborative editing features
- [ ] Plugin system
- [ ] AI-powered suggestions

## Acknowledgments

Built with:
- [Fiber v2](https://gofiber.io/) - Web framework
- [PostgreSQL](https://www.postgresql.org/) - Database
- [pgx/v5](https://github.com/jackc/pgx) - PostgreSQL driver
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration
