package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/momokii/go-cli-notes/cmd/cli/tui/components"
	"github.com/momokii/go-cli-notes/internal/model"
)

// NoteListModel is the model for the note list view
type NoteListModel struct {
	client    *client.APIClient
	authState *client.AuthState
	notes     []*model.Note
	total     int64
	page      int
	limit     int
	loading   bool
	err       error
	table     components.Table
	paginator components.Paginator
	width     int
	height    int
	// Filter state (for Phase D)
	search    string
	tagFilter *string
}

// NewNoteListModel creates a new note list model
func NewNoteListModel(apiClient *client.APIClient, authState *client.AuthState) NoteListModel {
	table := components.NewTable()
	paginator := components.NewPaginator()

	return NoteListModel{
		client:    apiClient,
		authState: authState,
		notes:     []*model.Note{},
		page:      1,
		limit:     20,
		loading:   true,
		table:     table,
		paginator: paginator,
		width:     80,
		height:    24,
	}
}

// SetTagFilter sets a tag filter and resets the model state
func (m NoteListModel) SetTagFilter(tagID string) NoteListModel {
	m.tagFilter = &tagID
	m.page = 1                 // Reset to first page
	m.notes = []*model.Note{}  // Clear cached notes
	m.loading = true           // Will fetch on Init
	return m
}

// Init initializes the note list model
func (m NoteListModel) Init() tea.Cmd {
	return m.fetchNotesCmd()
}

// fetchNotesCmd returns a command that fetches notes
func (m NoteListModel) fetchNotesCmd() tea.Cmd {
	return func() tea.Msg {
		filter := model.NoteFilter{
			Page:  m.page,
			Limit: m.limit,
		}
		if m.search != "" {
			filter.Search = m.search
		}
		if m.tagFilter != nil {
			filter.TagID = m.tagFilter
		}

		notes, total, err := m.client.ListNotes(filter)
		if err != nil {
			return noteListErrMsg{err}
		}
		return noteListFetchedMsg{notes: notes, total: total}
	}
}

// Update handles messages for the note list model
func (m NoteListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			// Request help view
			return m, func() tea.Msg {
				return ShowHelpMsg{}
			}
		case "esc":
			// Go back to dashboard
			return m, func() tea.Msg {
				return ShowDashboardMsg{}
			}
		case "j", "down":
			// Move down
			m.table.CursorDown()
			return m, nil
		case "k", "up":
			// Move up
			m.table.CursorUp()
			return m, nil
		case "g":
			// Go to top
			m.table.Top()
			return m, nil
		case "G":
			// Go to bottom
			m.table.Bottom()
			return m, nil
		case "enter", " ":
			// Open selected note
			if selected := m.table.SelectedItem(); selected != nil {
				noteID, err := uuid.Parse(selected.ID)
				if err == nil {
					return m, func() tea.Msg {
						return OpenNoteMsg{NoteID: noteID}
					}
				}
			}
			return m, nil
		case "ctrl+n":
			// Next page
			if m.paginator.CanGoNext() {
				m.page++
				m.paginator.NextPage()
				return m, m.fetchNotesCmd()
			}
			return m, nil
		case "ctrl+p":
			// Previous page
			if m.paginator.CanGoPrev() {
				m.page--
				m.paginator.PrevPage()
				return m, m.fetchNotesCmd()
			}
			return m, nil
		}

	case noteListFetchedMsg:
		m.notes = msg.notes
		m.total = msg.total
		m.loading = false

		// Update table rows
		rows := make([]components.TableRow, len(msg.notes))
		for i, note := range msg.notes {
			rows[i] = components.TableRow{
				ID:          note.ID.String(),
				Title:       note.Title,
				Description: m.formatNoteDescription(note),
				Metadata:    m.formatNoteMetadata(note),
			}
		}
		m.table.SetItems(rows)

		// Update paginator
		m.paginator.SetTotalItems(int(msg.total))
		m.paginator.SetPerPage(m.limit)
		m.paginator.SetPage(m.page)

		// Update sizes
		m.table.SetSize(m.width, m.height-3) // Leave room for header/paginator
		m.paginator.SetWidth(m.width)

		return m, nil

	case noteListErrMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetSize(msg.Width, msg.Height-3)
		m.paginator.SetWidth(msg.Width)
		return m, nil
	}

	// Update child components
	var cmd tea.Cmd
	var model tea.Model

	model, cmd = m.table.Update(msg)
	m.table = *(model.(*components.Table))
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the note list
func (m NoteListModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	var content string

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true)

	content += headerStyle.Render("NOTES")
	if m.total > 0 {
		content += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Gray
			Render(fmt.Sprintf(" (%d total)", m.total))
	}
	content += "\n\n"

	// Table
	content += m.table.View()
	content += "\n"

	// Paginator
	if m.total > int64(m.limit) {
		paginatorView := m.paginator.View()
		if paginatorView != "" {
			content += "\n" + paginatorView
		}
	}

	// Quick actions hint
	content += "\n"
	content += m.renderQuickActions()

	return content
}

// renderLoading renders the loading state
func (m NoteListModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Loading notes...")
}

// renderError renders the error state
func (m NoteListModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render(fmt.Sprintf("Error loading notes: %v", m.err))
}

// renderQuickActions renders the quick actions hint
func (m NoteListModel) renderQuickActions() string {
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	return hintStyle.Render("↑↓:nav Enter:open Ctrl+N/P:page ?:help ESC:back q:quit")
}

// formatNoteDescription formats the note description for the table
func (m NoteListModel) formatNoteDescription(note *model.Note) string {
	// Show tags if available
	if len(note.Tags) > 0 {
		tags := ""
		for i, tag := range note.Tags {
			if i > 0 {
				tags += " "
			}
			tags += "#" + tag.Name
		}
		return tags
	}

	// Show preview of content
	if len(note.Content) > 60 {
		return note.Content[:60] + "..."
	}
	return note.Content
}

// formatNoteMetadata formats the note metadata for the table
func (m NoteListModel) formatNoteMetadata(note *model.Note) string {
	// Show type and access info
	noteType := string(note.NoteType)
	if noteType == "" {
		noteType = "note"
	}

	var info string
	info += noteType

	// Show access count if any
	if note.AccessCount > 0 {
		info += fmt.Sprintf(" • %d views", note.AccessCount)
	}

	// Show time since last update
	info += " • " + formatTimeAgo(note.UpdatedAt)

	return info
}

// Message types for note list

type noteListFetchedMsg struct {
	notes []*model.Note
	total int64
}

type noteListErrMsg struct {
	err error
}

// View request messages
type ShowDashboardMsg struct{}
type OpenNoteMsg struct {
	NoteID uuid.UUID
}
