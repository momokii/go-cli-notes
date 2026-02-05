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

// TagListModel is the model for the tag list view
type TagListModel struct {
	client         *client.APIClient
	authState      *client.AuthState
	tags           []*model.TagWithCount
	loading        bool
	err            error
	selectedIndex  int
	showCreateForm bool
	showEditForm   bool
	showConfirm    bool
	editingTagID   uuid.UUID
	createForm     components.TextInput
	editForm       components.TextInput
	confirmDialog  components.ConfirmDialog
	width          int
	height         int
}

// NewTagListModel creates a new tag list model
func NewTagListModel(apiClient *client.APIClient, authState *client.AuthState) TagListModel {
	createForm := components.NewTextInput()
	createForm.SetPlaceholder("Tag name...")
	createForm.SetWidth(30)

	editForm := components.NewTextInput()
	editForm.SetPlaceholder("New tag name...")
	editForm.SetWidth(30)

	return TagListModel{
		client:         apiClient,
		authState:      authState,
		loading:        true,
		selectedIndex:  0,
		createForm:     createForm,
		editForm:       editForm,
		width:          80,
		height:         24,
	}
}

// Init initializes the tag list model
func (m TagListModel) Init() tea.Cmd {
	return m.fetchTagsCmd()
}

// fetchTagsCmd returns a command that fetches all tags
func (m TagListModel) fetchTagsCmd() tea.Cmd {
	return func() tea.Msg {
		// Get tags - note: API returns []*Tag, we'll get TagWithCount if available
		tags, err := m.client.GetTags()
		if err != nil {
			return TagListErrMsg{err}
		}

		// Convert to TagWithCount (note count will be 0 for basic Tag)
		// If API doesn't return counts, we'll display tags without counts
		var tagsWithCount []*model.TagWithCount
		for _, tag := range tags {
			tagsWithCount = append(tagsWithCount, &model.TagWithCount{
				ID:        tag.ID,
				UserID:    tag.UserID,
				Name:      tag.Name,
				Color:     tag.Color,
				CreatedAt: tag.CreatedAt,
				NoteCount: 0, // Would need separate API call to get counts
			})
		}

		return TagsFetchedMsg{Tags: tagsWithCount}
	}
}

// Update handles messages for the tag list model
func (m TagListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle confirmation dialog first
		if m.showConfirm {
			m.confirmDialog.Update(msg)
			if m.confirmDialog.IsYesSelected() {
				return m, m.deleteTagCmd()
			} else if m.confirmDialog.IsNoSelected() || msg.String() == "esc" {
				m.showConfirm = false
				return m, nil
			}
			return m, nil
		}

		// Handle create form
		if m.showCreateForm {
			switch msg.String() {
			case "enter":
				if m.createForm.Value() != "" {
					return m, m.createTagCmd(m.createForm.Value())
				}
			case "esc":
				m.showCreateForm = false
				m.createForm.SetValue("")
				m.createForm.Blur()
				return m, nil
			}
			cmd := m.createForm.Update(msg)
			return m, cmd
		}

		// Handle edit form
		if m.showEditForm {
			switch msg.String() {
			case "enter":
				if m.editForm.Value() != "" {
					return m, m.updateTagCmd(m.editingTagID, m.editForm.Value())
				}
			case "esc":
				m.showEditForm = false
				m.editForm.SetValue("")
				m.editForm.Blur()
				m.editingTagID = uuid.Nil
				return m, nil
			}
			cmd := m.editForm.Update(msg)
			return m, cmd
		}

		// Normal key handling
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			return m, func() tea.Msg {
				return ShowHelpMsg{}
			}
		case "esc":
			return m, func() tea.Msg {
				return ShowDashboardMsg{}
			}
		case "j", "down":
			if m.selectedIndex < len(m.tags)-1 {
				m.selectedIndex++
			}
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "c":
			// Create new tag
			m.showCreateForm = true
			m.createForm.Focus()
			return m, nil
		case "e":
			// Edit selected tag
			if len(m.tags) > 0 && m.selectedIndex >= 0 {
				m.showEditForm = true
				m.editingTagID = m.tags[m.selectedIndex].ID
				m.editForm.SetValue(m.tags[m.selectedIndex].Name)
				m.editForm.Focus()
			}
			return m, nil
		case "d":
			// Delete selected tag
			if len(m.tags) > 0 && m.selectedIndex >= 0 {
				m.showConfirm = true
				m.confirmDialog = components.NewConfirmDialog("Delete this tag?")
				m.confirmDialog.SetSubtext("This will remove the tag from all notes.")
				m.confirmDialog.Focus()
			}
			return m, nil
		case "enter":
			// Filter notes by selected tag
			if len(m.tags) > 0 && m.selectedIndex >= 0 {
				return m, func() tea.Msg {
					return FilterNotesByTagMsg{TagID: m.tags[m.selectedIndex].ID, TagName: m.tags[m.selectedIndex].Name}
				}
			}
		}

	case TagsFetchedMsg:
		m.tags = msg.Tags
		m.loading = false
		if len(m.tags) > 0 && m.selectedIndex >= len(m.tags) {
			m.selectedIndex = len(m.tags) - 1
		}
		return m, nil

	case TagCreatedMsg:
		m.showCreateForm = false
		m.createForm.SetValue("")
		m.createForm.Blur()
		return m, m.fetchTagsCmd()

	case TagUpdatedMsg:
		m.showEditForm = false
		m.editForm.SetValue("")
		m.editForm.Blur()
		m.editingTagID = uuid.Nil
		return m, m.fetchTagsCmd()

	case TagDeletedMsg:
		m.showConfirm = false
		if m.selectedIndex >= len(m.tags) {
			m.selectedIndex = len(m.tags) - 1
		}
		return m, m.fetchTagsCmd()

	case TagListErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.createForm.SetWidth(msg.Width - 20)
		m.editForm.SetWidth(msg.Width - 20)
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// createTagCmd returns a command that creates a new tag
func (m TagListModel) createTagCmd(name string) tea.Cmd {
	return func() tea.Msg {
		tag, err := m.client.CreateTag(name)
		if err != nil {
			return TagListErrMsg{Err: err}
		}
		return TagCreatedMsg{Tag: tag}
	}
}

// updateTagCmd returns a command that updates a tag
func (m TagListModel) updateTagCmd(id uuid.UUID, name string) tea.Cmd {
	return func() tea.Msg {
		tag, err := m.client.UpdateTag(id, name)
		if err != nil {
			return TagListErrMsg{Err: err}
		}
		return TagUpdatedMsg{Tag: tag}
	}
}

// deleteTagCmd returns a command that deletes the selected tag
func (m TagListModel) deleteTagCmd() tea.Cmd {
	if len(m.tags) == 0 || m.selectedIndex < 0 {
		return nil
	}
	tagID := m.tags[m.selectedIndex].ID
	return func() tea.Msg {
		if err := m.client.DeleteTag(tagID); err != nil {
			return TagListErrMsg{err}
		}
		return TagDeletedMsg{}
	}
}

// View renders the tag list view
func (m TagListModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	if m.showConfirm {
		return m.renderList() + "\n\n" + m.confirmDialog.View()
	}

	if m.showCreateForm {
		return m.renderList() + "\n\n" + m.renderCreateForm()
	}

	if m.showEditForm {
		return m.renderList() + "\n\n" + m.renderEditForm()
	}

	return m.renderList()
}

// renderLoading renders the loading state
func (m TagListModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Loading tags...")
}

// renderError renders the error state
func (m TagListModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render("Error: " + m.err.Error())
}

// renderList renders the tag list
func (m TagListModel) renderList() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true).
		MarginTop(1)

	var content string

	// Title
	content += titleStyle.Render("TAGS") + "\n\n"

	if len(m.tags) == 0 {
		content += mutedStyle.Render("(no tags)")
		content += "\n\n"
		content += hintStyle.Render("c:create ESC:back ?:help")
		return content
	}

	// Tags list
	for i, tag := range m.tags {
		var line string
		if i == m.selectedIndex {
			line = selectedStyle.Render("→ " + tag.Name)
			if tag.NoteCount > 0 {
				line += " (" + formatCount(tag.NoteCount) + ")"
			}
		} else {
			line = tagStyle.Render("  "+tag.Name)
			if tag.NoteCount > 0 {
				line += mutedStyle.Render(" (" + formatCount(tag.NoteCount) + ")")
			}
		}
		content += line + "\n"
	}

	// Hints
	content += "\n" + hintStyle.Render("c:create e:edit d:delete Enter:view notes j/k:navigate ESC:back")

	return content
}

// renderCreateForm renders the create form
func (m TagListModel) renderCreateForm() string {
	formStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	var content string
	content += formStyle.Render("Create New Tag") + "\n\n"
	content += m.createForm.View()
	content += "\n\n"
	content += hintStyle.Render("Enter:create ESC:cancel")

	return content
}

// renderEditForm renders the edit form
func (m TagListModel) renderEditForm() string {
	formStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	var content string
	content += formStyle.Render("Edit Tag") + "\n\n"
	content += m.editForm.View()
	content += "\n\n"
	content += hintStyle.Render("Enter:save ESC:cancel")

	return content
}

// formatCount formats a number for display
func formatCount(n int) string {
	if n == 1 {
		return "1 note"
	}
	// Simple int to string conversion
	if n < 10 {
		return string(rune('0'+n)) + " notes"
	}
	// For larger numbers, use a simple conversion
	return fmt.Sprintf("%d notes", n)
}

// IsInputFocused returns whether the create/edit form is focused
// This allows the main TUI to skip global key handlers when typing
func (m TagListModel) IsInputFocused() bool {
	// Check if create form is shown and focused
	if m.showCreateForm && m.createForm.Focused() {
		return true
	}
	// Check if edit form is shown and focused
	if m.showEditForm && m.editForm.Focused() {
		return true
	}
	return false
}

// BlurInput removes focus from create/edit forms and hides them
// Called when leaving the tag view to prevent stuck focus state
func (m TagListModel) BlurInput() TagListModel {
	// Clear form visibility flags to prevent forms from showing when returning
	m.showCreateForm = false
	m.showEditForm = false
	// Also blur the forms to remove focus
	m.createForm.Blur()
	m.editForm.Blur()
	// Clear form values
	m.createForm.SetValue("")
	m.editForm.SetValue("")
	// Clear editing state
	m.editingTagID = uuid.Nil
	// Hide confirm dialog
	m.showConfirm = false
	return m
}

// Message types for tag list

type TagsFetchedMsg struct {
	Tags []*model.TagWithCount
}

type TagCreatedMsg struct {
	Tag *model.Tag
}

type TagUpdatedMsg struct {
	Tag *model.Tag
}

type TagDeletedMsg struct{}

type TagListErrMsg struct {
	Err error
}

// FilterNotesByTagMsg is a message to filter notes by tag
type FilterNotesByTagMsg struct {
	TagID   uuid.UUID
	TagName string
}
