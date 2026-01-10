# CLI Guide - Complete Reference

This guide provides comprehensive documentation for all CLI commands, syntax, and features.

## Table of Contents

- [Installation](#installation)
- [Authentication](#authentication)
- [Configuration](#configuration)
- [Note Commands](#note-commands)
- [Tag Commands](#tag-commands)
- [Search](#search)
- [Analytics](#analytics)
- [Wiki-Style Links](#wiki-style-links)
- [Examples](#examples)

---

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/momokii/go-cli-notes.git
cd go-cli-notes

# Build the CLI
go build -o kg-cli ./cmd/cli

# (Optional) Install globally
sudo cp kg-cli /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/momokii/go-cli-notes/cmd/cli@latest
```

### Verify Installation

```bash
kg-cli --version
kg-cli --help
```

---

## Authentication

### Register

Create a new account.

**Syntax:**
```bash
kg-cli register
```

**Interactive prompts:**
- `Username:` Your username (3-100 characters, alphanumeric only)
- `Email:` Your email address
- `Password:` Your password (minimum 8 characters)

**Example:**
```bash
$ kg-cli register
Username: johndoe
Email: user@example.com
Password: ********

Registration successful! Please login with `kg-cli login`
```

### Login

Authenticate with your credentials.

**Syntax:**
```bash
kg-cli login
```

**Interactive prompts:**
- `Email:` Your email address
- `Password:` Your password

**Example:**
```bash
$ kg-cli login
Email: user@example.com
Password: ********

Login successful!
```

### Logout

Log out and clear stored credentials.

**Syntax:**
```bash
kg-cli logout
```

**Example:**
```bash
$ kg-cli logout
Logged out successfully
```

### Status

Check your authentication and connection status.

**Syntax:**
```bash
kg-cli status
```

**Example:**
```bash
$ kg-cli status
Knowledge Garden CLI Status
==========================
API URL: http://localhost:8080
Status: Authenticated
Email: user@example.com
```

---

## Configuration

### Config File Location

The CLI stores configuration in:
- `~/.config/kg-cli/config.yaml` (Linux/Mac)
- `%USERPROFILE%\.config\kg-cli\config.yaml` (Windows)

### Default Configuration

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

Override configuration using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `KG_CLI_API_BASE_URL` | API server URL | `http://localhost:8080` |
| `KG_CLI_API_TIMEOUT` | Request timeout (seconds) | `30` |
| `KG_CLI_EDITOR` | External editor | `$EDITOR` or `vi` |

---

## Note Commands

### List Notes

Display all notes with pagination and filtering.

**Syntax:**
```bash
kg-cli note list [flags]
```

**Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--page` | `-p` | Page number | `1` |
| `--limit` | `-l` | Notes per page (1-100) | `20` |
| `--search` | `-s` | Search query | - |
| `--tag` | `-t` | Filter by tag name or ID | - |

**Examples:**
```bash
# List first 20 notes (default)
kg-cli note list

# List page 2
kg-cli note list --page 2

# List 50 notes per page
kg-cli note list --limit 50

# Search in notes
kg-cli note list --search "golang"

# Filter by tag name
kg-cli note list --tag "programming"

# Filter by tag ID
kg-cli note list --tag 123e4567-e89b-12d3-a456-426614174000

# Combine filters
kg-cli note list --search "golang" --tag "programming" --limit 10
```

### Get Note

Display a specific note by ID.

**Syntax:**
```bash
kg-cli note get <note-id>
```

**Arguments:**
- `note-id` - The UUID of the note (required)

**Example:**
```bash
$ kg-cli note get 123e4567-e89b-12d3-a456-426614174000
Title: My Go Project
Type: note
Word Count: 150
Reading Time: 1 min
Created: 2026-01-04 10:30:00
Updated: 2026-01-04 11:15:00

Content:
This is my Go project...
```

### Create Note

Create a new note.

**Syntax:**
```bash
kg-cli note create [flags]
```

**Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--title` | `-t` | Note title (required) | - |
| `--content` | `-c` | Note content | Empty string |
| `--type` | `-T` | Note type | `note` |

**Note Types:**
- `note` - Regular notes
- `daily` - Daily journal entries
- `meeting` - Meeting notes
- `idea` - Quick ideas and thoughts

**Examples:**
```bash
# Create note with title and content (most common)
kg-cli note create -t "My First Note" -c "This is my note content"

# Create with long flags
kg-cli note create --title "Project Idea" --content "Build a CLI app"

# Create note with title only
kg-cli note create -t "My First Note"

# Create a meeting note
kg-cli note create --title "Team Standup" --type meeting

# Create with short flags
kg-cli note create -t "Quick thought" -T idea
```

### Search Notes

Search notes using full-text search.

**Syntax:**
```bash
kg-cli note search <query> [flags]
```

**Arguments:**
- `query` - Search query (required)

**Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--page` | `-p` | Page number | `1` |
| `--limit` | `-l` | Results per page (1-100) | `20` |

**Examples:**
```bash
# Basic search
kg-cli note search "golang"

# Search with pagination
kg-cli note search "golang" --page 2 --limit 10

# Search for phrases
kg-cli note search "full-text search"
```

### Daily Note

Get or create a daily note for a specific date.

**Syntax:**
```bash
kg-cli note daily [date]
```

**Arguments:**
- `date` - Date in YYYY-MM-DD format, or "today" (optional, default: today)

**Examples:**
```bash
# Get or create today's daily note
kg-cli note daily

# Get or create a specific date's daily note
kg-cli note daily 2026-01-04

# Get tomorrow's daily note
kg-cli note daily 2026-01-05
```

### Update Note

Update an existing note's title or content. **Interactive mode is enabled by default** - it shows current values and prompts for changes.

**Syntax:**
```bash
kg-cli note update <note-id> [flags]
```

**Arguments:**
- `note-id` - The UUID of the note (required)

**Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--title` | `-t` | New note title (skips interactive mode) | - |
| `--content` | `-c` | New note content (skips interactive mode) | - |

**Interactive Mode (Default)**

When no flags are provided, the update command enters interactive mode:

```bash
$ kg-cli note update 123e4567-e89b-12d3-a456-426614174000
Updating note: My Note Title
Current values shown - leave empty to keep existing value

Current title: My Note Title
New title (press Enter to keep current):

Edit content? (Y/n):
[Opens your $EDITOR with current content]

Save changes? (Y/n): y
Note updated successfully!
```

**How Interactive Mode Works:**

1. Shows current note title
2. Prompts for new title (press Enter to keep current)
3. Asks if you want to edit content (opens `$EDITOR` - vi, nano, vim, etc.)
4. Shows all changes and asks for confirmation
5. Saves only if you confirm

**Setting Your Editor:**

```bash
# Set your preferred editor
export EDITOR=vim
export EDITOR=nano
export EDITOR=code

# Then run update (uses your editor)
kg-cli note update <note-id>
```

**Flag-Based Update (For Automation)**

Use flags for quick updates or scripting:

```bash
# Update title only
kg-cli note update <uuid> --title "New Title"

# Update content only
kg-cli note update <uuid> --content "New content here"

# Update both title and content
kg-cli note update <uuid> -t "Updated Title" -c "Updated content"

# Add wiki-style links in content
kg-cli note update <uuid> --content "See [[Related Note]] for more info"
```

**Examples:**

```bash
# Interactive update (recommended for manual edits)
kg-cli note update 123e4567-e89b-12d3-a456-426614174000

# Quick title change (skips interactive mode)
kg-cli note update 123e4567-e89b-12d3-a456-426614174000 -t "New Title"

# Quick content replacement (skips interactive mode)
kg-cli note update 123e4567-e89b-12d3-a456-426614174000 -c "New content"

# Using with $EDITOR set
export EDITOR=nano
kg-cli note update <note-id>  # Opens nano editor
```

### Delete Note

Delete a note permanently (with confirmation prompt).

**Syntax:**
```bash
kg-cli note delete <note-id>
```

**Arguments:**
- `note-id` - The UUID of the note (required)

**Example:**
```bash
$ kg-cli note delete 123e4567-e89b-12d3-a456-426614174000
Are you sure you want to delete this note? (y/N): y

Note deleted successfully!
```

**Note:** Deleted notes cannot be recovered. The deletion requires confirmation to prevent accidental deletion.

### View Links

View outgoing wiki-style links from a note.

**Syntax:**
```bash
kg-cli note links <note-id>
```

**Arguments:**
- `note-id` - The UUID of the note (required)

**Example:**
```bash
$ kg-cli note links 123e4567-e89b-12d3-a456-426614174000
Found 2 outgoing link(s):

To: Go Fiber Framework Research (ID: 456e7890-e89b-12d3-a456-426614174001)
Context: See [[Go Fiber Framework Research]] for more details
Created: 2026-01-04 10:30
---
To: PostgreSQL Setup (ID: 567e8901-e89b-12d3-a456-426614174002)
Context: Database setup in [[PostgreSQL Setup]]
Created: 2026-01-04 10:35
---
```

### View Backlinks

View wiki-style backlinks (links from other notes pointing to this note).

**Syntax:**
```bash
kg-cli note backlinks <note-id>
```

**Arguments:**
- `note-id` - The UUID of the note (required)

**Example:**
```bash
$ kg-cli note backlinks 123e4567-e89b-12d3-a456-426614174000
Found 1 backlink(s):

From: Project Overview (ID: 789e9012-e89b-12d3-a456-426614174003)
Context: Main project is [[Go CLI Project]]
Created: 2026-01-04 11:00
---
```

### View Tags on Note

View all tags associated with a specific note.

**Syntax:**
```bash
kg-cli note tags <note-id>
```

**Arguments:**
- `note-id` - The UUID of the note (required)

**Example:**
```bash
$ kg-cli note tags 123e4567-e89b-12d3-a456-426614174000
Found 2 tag(s):

ID: 123e4567-e89b-12d3-a456-426614174001
Name: programming
---
ID: 123e4567-e89b-12d3-a456-426614174002
Name: golang
---
```

---

## Tag Commands

### List Tags

Display all tags.

**Syntax:**
```bash
kg-cli tag list
```

**Example:**
```bash
$ kg-cli tag list
Found 3 tag(s):

ID: 123e4567-e89b-12d3-a456-426614174000
Name: programming
---
ID: 123e4567-e89b-12d3-a456-426614174001
Name: golang
---
ID: 123e4567-e89b-12d3-a456-426614174002
Name: ideas
---
```

### Create Tag

Create a new tag.

**Syntax:**
```bash
kg-cli tag create <name>
```

**Arguments:**
- `name` - Tag name (required)

**Examples:**
```bash
# Create a tag
kg-cli tag create "programming"

# Create a multi-word tag
kg-cli tag create "web development"

# Create hierarchical tags
kg-cli tag create "dev/backend"
kg-cli tag create "dev/frontend"
kg-cli tag create "learning/golang"
```

### Update Tag

Update an existing tag's name.

**Syntax:**
```bash
kg-cli tag update <id> <new-name>
```

**Arguments:**
- `id` - Tag UUID (required)
- `new-name` - New tag name (required)

**Example:**
```bash
$ kg-cli tag update 123e4567-e89b-12d3-a456-426614174000 "golang-dev"
Tag updated successfully!
ID: 123e4567-e89b-12d3-a456-426614174000
New Name: golang-dev
```

### Delete Tag

Delete a tag (with confirmation prompt).

**Syntax:**
```bash
kg-cli tag delete <id>
```

**Arguments:**
- `id` - Tag UUID (required)

**Example:**
```bash
$ kg-cli tag delete 123e4567-e89b-12d3-a456-426614174000
Are you sure you want to delete this tag? (y/N): y

Tag deleted successfully!
```

### Add Tag to Note

Add a tag to a note. Supports both tag ID and tag name.

**Syntax:**
```bash
kg-cli tag add <tag-id-or-name> <note-id>
```

**Arguments:**
- `tag-id-or-name` - Tag UUID or tag name (required)
- `note-id` - Note UUID (required)

**Examples:**
```bash
# Add tag by name (recommended)
kg-cli tag add "programming" 123e4567-e89b-12d3-a456-426614174000

# Add tag by ID
kg-cli tag add 123e4567-e89b-12d3-a456-426614174001 123e4567-e89b-12d3-a456-426614174000
```

### Remove Tag from Note

Remove a tag from a note. Supports both tag ID and tag name.

**Syntax:**
```bash
kg-cli tag remove <tag-id-or-name> <note-id>
```

**Arguments:**
- `tag-id-or-name` - Tag UUID or tag name (required)
- `note-id` - Note UUID (required)

**Examples:**
```bash
# Remove tag by name (recommended)
kg-cli tag remove "programming" 123e4567-e89b-12d3-a456-426614174000

# Remove tag by ID
kg-cli tag remove 123e4567-e89b-12d3-a456-426614174001 123e4567-e89b-12d3-a456-426614174000
```

---

## Analytics

### Stats

Display user statistics.

**Syntax:**
```bash
kg-cli stats
```

**Example:**
```bash
$ kg-cli stats
Knowledge Garden Statistics
==========================
Total Notes: 42
Total Tags: 8
Total Links: 15
Total Words: 12,450
Notes Created Today: 3
Notes Created This Week: 12
Last Activity: Sun, 04 Jan 2026 14:30:00 UTC
```

### Activity

Display recent activity.

**Syntax:**
```bash
kg-cli activity [flags]
```

**Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--limit` | `-l` | Number of activities to show | `10` |

**Examples:**
```bash
# Show last 10 activities
kg-cli activity

# Show last 20 activities
kg-cli activity --limit 20

# Show last 50 activities
kg-cli activity -l 50
```

### Trending

Display frequently accessed notes.

**Syntax:**
```bash
kg-cli trending [flags]
```

**Flags:**
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--limit` | `-l` | Number of trending notes to show | `5` |

**Examples:**
```bash
# Show top 5 trending notes
kg-cli trending

# Show top 10 trending notes
kg-cli trending --limit 10
```

---

## Wiki-Style Links

### Syntax

Connect notes using double brackets:

```markdown
This note is related to [[Another Note]].

You can reference multiple notes: [[Note 1]], [[Note 2]], and [[Note 3]].

Links work with spaces in titles: [[My Long Note Title]].
```

### How Links Work (Step-by-Step)

**Important:** Links are created **automatically** when you use `[[Note Title]]` syntax. The target note must already exist for the link to be created.

**Step 1: Create the target note first**
```bash
kg-cli note create --title "Go Programming Tips" \
  --content "Here are some useful Go programming tips..."
# Save the note ID: TARGET_ID=<output-id>
```

**Step 2: Create another note that references it**
```bash
kg-cli note create --title "My Learning Journey" \
  --content "Today I learned Go. See [[Go Programming Tips]] for more details."
# Save the note ID: SOURCE_ID=<output-id>
```

**Step 3: Verify the link was created**
```bash
# View outgoing links from My Learning Journey
kg-cli note links $SOURCE_ID
# Output should show: "To: Go Programming Tips"

# View backlinks to Go Programming Tips
kg-cli note backlinks $TARGET_ID
# Output should show: "From: My Learning Journey"
```

### Key Points

1. **Target note must exist first** - Links are only created if the referenced note already exists
2. **Links are created on save** - When you create or update a note, the system parses `[[...]]` syntax and creates links
3. **Links are bidirectional** - When A links to B, B has a backlink from A
4. **Links update automatically** - When you update note content, old links are removed and new ones are created

### Troubleshooting

**No links found?**
- Make sure the target note exists (check with `kg-cli note list`)
- Note titles must match exactly (case-insensitive)
- Try updating the note to trigger link parsing: `kg-cli note update <source-id>`

**Link to non-existent note?**
- Create the target note first
- Then update your source note to create the link

### Viewing Links via CLI

```bash
# View outgoing links from a note
kg-cli note links <note-id>

# View backlinks to a note
kg-cli note backlinks <note-id>
```

### Viewing Links via API

```bash
# Get outgoing links from a note
curl http://localhost:8080/api/v1/notes/<id>/links \
  -H "Authorization: Bearer <token>"

# Get backlinks to a note
curl http://localhost:8080/api/v1/notes/<id>/backlinks \
  -H "Authorization: Bearer <token>"

# Get the full knowledge graph
curl http://localhost:8080/api/v1/notes/graph \
  -H "Authorization: Bearer <token>"
```

---

## Examples

### Daily Workflow

```bash
# 1. Check status
kg-cli status

# 2. Create today's daily note
kg-cli note daily

# 3. List recent notes
kg-cli note list --limit 10

# 4. Check activity
kg-cli activity

# 5. View stats
kg-cli stats
```

### Research Project Workflow

```bash
# 1. Create project note
kg-cli note create --title "Go CLI Project" \
  --content "Building a CLI tool for note-taking"

# 2. Create research notes
kg-cli note create --title "Go Fiber Framework Research" \
  --type note \
  --content "Fiber is inspired by Express..."

kg-cli note create --title "PostgreSQL FTS Setup" \
  --type note \
  --content "Full-text search configuration..."

# 3. Link related notes
kg-cli note update <id> --content "See [[Go Fiber Framework Research]]"

# 4. Create tags
kg-cli tag create "golang"
kg-cli tag create "research"

# 5. Search notes
kg-cli note search "Fiber"
kg-cli note search "database"
```

### Meeting Notes Workflow

```bash
# 1. Create meeting note
kg-cli note create \
  --title "Sprint Planning Meeting" \
  --type meeting \
  --content "Date: 2026-01-04
Attendees: John, Jane, Bob

Agenda:
1. Review sprint goals
2. Plan tasks
3. Assign responsibilities

Action items:
- John: Setup database
- Jane: Implement authentication
- Bob: Design UI"

# 2. Create follow-up note
kg-cli note create \
  --title "Action Items from Sprint Planning" \
  --type note \
  --content "See [[Sprint Planning Meeting]] for details"
```

### Learning Journal Workflow

```bash
# 1. Get today's daily note
kg-cli note daily

# 2. Create learning notes
kg-cli note create \
  --title "Learned about Go Context" \
  --content "Context is used for request-scoped data..."

kg-cli note create \
  --title "Understanding Goroutines" \
  --content "Goroutines are lightweight threads..."

# 3. Search for review
kg-cli note search "goroutine"

# 4. Check progress
kg-cli stats
```

---

## Tips and Tricks

### Quick Login

Store credentials for quick login:
```bash
# Create alias for login
alias kg-login='echo "user@example.com\npassword" | kg-cli login'
```

### Aliases for Common Commands

```bash
# Add to ~/.bashrc or ~/.zshrc
alias kgnotes='kg-cli note list'
alias kgstats='kg-cli stats'
alias kgtoday='kg-cli note daily'
alias kgsearch='kg-cli note search'
```

### Scripting with CLI

```bash
#!/bin/bash
# Backup script

# 1. Login
TOKEN=$(echo "user@example.com\npassword" | kg-cli login | grep -oP 'access_token.*' | cut -d: -f2)

# 2. Export all notes
kg-cli note list --limit 1000 > notes_backup.txt

# 3. Get stats
kg-cli stats > stats_backup.txt
```

### Keyboard Shortcuts (Future)

The planned TUI version will include:
- `n` - New note
- `/` - Search
- `g` - Knowledge graph
- `s` - Settings
- `q` - Quit
- `Ctrl+S` - Save
- `Esc` - Back

---

## Error Messages

### Common Errors and Solutions

**"not authenticated. Please run 'kg-cli login' first"**
- Solution: Run `kg-cli login` to authenticate

**"API error (status 401): unauthorized"**
- Solution: Your session expired. Run `kg-cli login` again

**"Failed to connect to API"**
- Solution: Check that the API server is running at `http://localhost:8080`
- Or set the correct API URL: `export KG_CLI_API_BASE_URL=http://your-api:8080`

**"no notes found"**
- Solution: You haven't created any notes yet. Use `kg-cli note create` to create your first note

---

## Advanced Usage

### JSON Output

For scripting, you can parse the JSON output:

```bash
# Get notes as JSON
kg-cli note list | jq '.notes[] | .title'

# Count notes
kg-cli note list | jq '.total'

# Get specific fields
kg-cli stats | jq '.total_words'
```

### Integration with Other Tools

```bash
# Use with fzf (fuzzy finder)
kg-cli note list | jq '.notes[] | .title' | fzf

# Open note in editor
NOTE_ID=$(kg-cli note list | jq '.notes[0].id')
kg-cli note get $NOTE_ID > /tmp/note.md
vim /tmp/note.md
```

### Batch Operations

```bash
# Create multiple notes from a file
while IFS= read -r title; do
  kg-cli note create --title "$title"
done < notes.txt

# Search and process results
kg-cli note search "project" | jq -r '.results[].note.id' | while read id; do
  kg-cli note get "$id"
done
```

---

## Getting Help

### Built-in Help

```bash
# Global help
kg-cli --help

# Command-specific help
kg-cli note --help
kg-cli note create --help
kg-cli tag --help
kg-cli stats --help
```

### Version Information

```bash
kg-cli --version
```

### Debug Mode

For debugging, set the environment variable:

```bash
# Enable verbose logging
export KG_CLI_DEBUG=true
kg-cli note list
```

---

## Next Steps

- Read the [README.md](README.md) for project overview
- Check the [API Documentation](#rest-api) for API usage
- Explore the [Knowledge Graph](#knowledge-graph-api) feature
- Learn about [Wiki-Style Links](#wiki-style-links)

For more information, visit: https://github.com/momokii/go-cli-notes
