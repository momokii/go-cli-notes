# Knowledge Garden CLI - TUI User Guide

## Overview

The Knowledge Garden CLI includes an interactive Terminal User Interface (TUI) that provides a feature-rich, vim-style navigation experience for managing your notes, tags, and knowledge graph.

## Installation & Setup

### Prerequisites

- Go 1.21 or later
- Valid Knowledge Garden account
- Terminal emulator with 80x24 minimum size

### Installation

```bash
# Clone the repository
git clone https://github.com/momokii/go-cli-notes.git
cd go-cli-notes

# Build the CLI
go build -o kg-cli ./cmd/cli

# (Optional) Install to your path
sudo cp kg-cli /usr/local/bin/
```

### Initial Setup

1. **Login to your account:**
   ```bash
   kg-cli login
   ```

2. **Launch the TUI:**
   ```bash
   kg-cli tui
   ```

## Quick Start

Once you're in the TUI:

| Key | Action |
|-----|--------|
| `?` | Show help screen |
| `n` | Create a new note |
| `/` | Search notes |
| `t` | View tags |
| `a` | View activity feed |
| `g` | View knowledge graph |
| `ESC` | Go back to dashboard |
| `q` | Quit TUI |

## Navigation

### Vim-Style Navigation

The TUI supports vim-style key bindings:

| Key | Action |
|-----|--------|
| `h` / `←` | Left / Previous |
| `j` / `↓` | Down / Next item |
| `k` / `↑` | Up / Previous item |
| `l` / `→` | Right / Select |
| `0` | Go to top of list |
| `G` | Go to bottom of list |
| `Enter` | Open / Select item |
| `ESC` | Go back / Cancel |

### Global Shortcuts

These shortcuts work from anywhere in the TUI:

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit TUI |
| `?` / `F1` | Show help |
| `/` | Quick search |
| `n` | New note |
| `ESC` | Return to dashboard |

## Views

### Dashboard

The dashboard provides an overview of your knowledge garden:

- **Statistics**: Note count, tag count, link count, word count
- **Recent Activity**: Latest actions on your notes
- **Trending Notes**: Most viewed notes

**Dashboard Shortcuts:**
| Key | Action |
|-----|--------|
| `l` | View note list |
| `n` | Create new note |
| `/` | Search |
| `t` | View tags |
| `a` | Activity feed |
| `g` | Knowledge graph |

### Note List

Browse and search through all your notes.

**Note List Shortcuts:**
| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Open selected note |
| `/` | Start new search |
| `n` | Create new note |
| `Ctrl+N` | Next page |
| `Ctrl+P` | Previous page |

### Note Detail

View and edit individual notes.

**Note View Shortcuts:**
| Key | Action |
|-----|--------|
| `TAB` | Switch between tabs (Content/Tags/Links/Backlinks) |
| `e` | Edit note |
| `d` | Delete note (in Content/Links/Backlinks tabs) or Remove selected tag (in Tags tab) |
| `a` | Add tag to note (in Tags tab only) |
| `↑` / `↓` or `j` / `k` | Navigate tags in Tags tab |
| `ESC` | Go back |

**Tags Tab Shortcuts:**
| Key | Action |
|-----|--------|
| `a` | Add tag to note (opens tag selector) |
| `d` | Remove selected tag |
| `↑` / `↓` or `j` / `k` | Select tag |
| `TAB` | Switch to next tab |

**Add Tag Form Shortcuts:**
| Key | Action |
|-----|--------|
| `Enter` | Add selected tag |
| `↑` / `↓` | Navigate tag list |
| `ESC` | Cancel |
| Type | Filter tags by name |

**Note Edit Shortcuts:**
| Key | Action |
|-----|--------|
| `Ctrl+S` | Save note |
| `Ctrl+C` | Cancel edit |
| `ESC` | Cancel edit |

### Tag List

Manage your tags and view notes by tag.

**Tag List Shortcuts:**
| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `c` | Create new tag |
| `e` | Edit selected tag |
| `d` | Delete selected tag |
| `Enter` | View notes with this tag |

### Search

Full-text search across all notes with highlighting.

**Search Shortcuts:**
| Key | Action |
|-----|--------|
| `/` | Start new search |
| `Enter` | Open selected result |
| `Ctrl+N` | Next page of results |
| `Ctrl+P` | Previous page of results |
| `ESC` | Clear search / Go back |

### Activity Feed

View recent activity on your notes.

**Activity Feed Shortcuts:**
| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Open affected note |
| `Ctrl+N` | Next page |
| `Ctrl+P` | Previous page |

### Knowledge Graph

Visualize connections between your notes in ASCII format.

**Graph Shortcuts:**
| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Space` | Expand/collapse connections |
| `+` | Show more nodes |
| `-` | Show fewer nodes |
| `Enter` | Open selected note |

## Creating Notes

1. Press `n` from anywhere to create a new note
2. Enter the note title
3. Enter the note content (supports Markdown)
4. Add tags by typing tag names
5. Press `Ctrl+S` to save or `ESC` to cancel

## Linking Notes

Create connections between notes using wiki-style links:

```
[[Note Title]] - Link to another note by title
```

Backlinks are automatically created when you link to notes.

## Search

The TUI supports full-text search across:

- Note titles
- Note content
- Tags

**Search Tips:**
- Use specific terms for better results
- Results are highlighted with yellow
- Navigate results with `j`/`k`
- Press `Enter` to open a result

## Keyboard Reference

### Complete Key Binding Reference

| Key | Global | Dashboard | Notes | Note Detail | Tags | Search | Activity | Graph |
|-----|--------|-----------|-------|-------------|-------|--------|----------|-------|
| `q` | Quit | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `?` | Help | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `ESC` | Back | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `/` | Search | ✓ | ✓ | - | - | ✓ | - | - |
| `n` | New note | ✓ | ✓ | - | - | - | - | - |
| `t` | Tags | ✓ | - | - | - | - | - | - |
| `a` | Add tag | - | - | - | ✓ | - | - | - |
| `a` | Activity | ✓ | - | - | - | - | - | - |
| `g` | Graph | ✓ | - | - | - | - | - | - |
| `j` | Down | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `k` | Up | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `Enter` | Open | - | ✓ | - | ✓ | ✓ | ✓ | ✓ |
| `TAB` | Tab | - | - | - | ✓ | - | - | - |
| `c` | Create | - | - | - | - | ✓ | - | - |
| `e` | Edit | - | ✓ | ✓ | - | ✓ | - | - |
| `d` | Delete | - | ✓ | ✓ | - | - | - | - |
| `+` | More | - | - | - | - | - | - | ✓ |
| `-` | Fewer | - | - | - | - | - | - | ✓ |

## Tips & Tricks

### Productivity Tips

1. **Quick Search**: Press `/` from anywhere to instantly search
2. **Fast Navigation**: Use `j`/`k` for faster list navigation
3. **Quick Actions**: Press `n` from any view to create a note
4. **Keyboard Shortcuts**: Learn the single-key shortcuts for speed

### Efficiency Tips

1. **Batch Operations**: Navigate with `j`/`k`, then press `Enter` to open multiple items
2. **Search First**: Use `/` to find notes quickly instead of browsing
3. **Use Tags**: Organize with tags and filter by tag with `t`
4. **Graph Navigation**: Use `g` to explore note connections

### Editing Tips

1. **Auto-Save**: Edits are not auto-saved - press `Ctrl+S` to save
2. **Cancel Edit**: Press `ESC` or `Ctrl+C` to cancel without saving
3. **Markdown**: Full Markdown support in note content
4. **Wiki Links**: Use `[[Note Title]]` to create bidirectional links

## Troubleshooting

### Session Expired

**Problem**: "Session expired. Please run 'kg-cli login'"

**Solution**: Your authentication token expired. Exit the TUI and run:
```bash
kg-cli login
```

### Terminal Too Small

**Problem**: "Terminal is too small for the TUI interface"

**Solution**: Resize your terminal to at least 80x24 characters

### Connection Error

**Problem**: "Failed to fetch data" or "Network error"

**Solution**: Check your internet connection and API URL configuration

### Blank Screen

**Problem**: TUI launches but screen is blank

**Solution**:
1. Ensure terminal supports colors
2. Try a different terminal emulator
3. Check `kg-cli status` for authentication issues

## Configuration

### Config File Location

`~/.config/kg-cli/config.yaml`

### Available Settings

```yaml
api:
  base_url: "https://cli-notes-api.kelanach.xyz/"
  timeout: 30

editor:
  command: "vim"  # Your preferred editor (not used in TUI)
```

## Advanced Features

### Wiki-Style Links

Create connections between notes using `[[Note Title]]` syntax. These create bidirectional links automatically.

### Backlinks

When note A links to note B, note B automatically shows note A in its backlinks section.

### Knowledge Graph

The graph view shows:
- **Nodes**: Your notes
- **Edges**: Links between notes
- **Expansion**: Press `Space` to expand/collapse connections
- **Zoom**: Use `+`/`-` to show more/fewer nodes

### Activity Tracking

All actions are tracked:
- Created notes
- Updated notes
- Deleted notes
- Viewed notes
- Search queries

## Support

For issues, questions, or contributions:
- GitHub: https://github.com/momokii/go-cli-notes
- Documentation: https://github.com/momokii/go-cli-notes/wiki

## Version

This guide is for TUI version 1.0.0

---

**Happy note-taking! 📝**
