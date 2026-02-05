package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/momokii/go-cli-notes/cmd/cli/tui/components"
	"github.com/momokii/go-cli-notes/internal/model"
)

// NoteDetailTab represents the different tabs in note detail view
type NoteDetailTab int

const (
	NoteContentTab NoteDetailTab = iota
	NoteTagsTab
	NoteLinksTab
	NoteBacklinksTab
)

// String returns the string representation of a tab
func (t NoteDetailTab) String() string {
	switch t {
	case NoteContentTab:
		return "Content"
	case NoteTagsTab:
		return "Tags"
	case NoteLinksTab:
		return "Links"
	case NoteBacklinksTab:
		return "Backlinks"
	default:
		return "Unknown"
	}
}

// NoteDetailModel is the model for viewing a single note
type NoteDetailModel struct {
	client        *client.APIClient
	authState     *client.AuthState
	note          *model.Note
	noteID        uuid.UUID
	loading       bool
	err           error
	currentTab    NoteDetailTab
	tags          []*model.Tag
	tagsErr       error
	links         []*model.LinkDetail
	linksErr      error
	backlinks     []*model.LinkDetail
	backlinksErr  error
	showConfirm   bool
	confirmDialog components.ConfirmDialog
	width         int
	height        int
	// Tag management fields
	selectedTagIndex     int
	availableTags        []*model.Tag
	availableTagsLoading bool
	availableTagsErr     error
	showAddTagForm       bool
	addTagInput          components.TextInput
	addTagFilter         string
	filteredAvailableTags []*model.Tag
	selectedAvailableIndex int
}

// NewNoteDetailModel creates a new note detail model
func NewNoteDetailModel(apiClient *client.APIClient, authState *client.AuthState) NoteDetailModel {
	addTagInput := components.NewTextInput()
	addTagInput.SetPlaceholder("Type tag name or select from list...")
	addTagInput.SetWidth(40)

	return NoteDetailModel{
		client:               apiClient,
		authState:            authState,
		loading:              true,
		currentTab:           NoteContentTab,
		width:                80,
		height:               24,
		selectedTagIndex:     -1, // No tag selected initially
		availableTagsLoading: false,
		showAddTagForm:       false,
		addTagInput:          addTagInput,
		addTagFilter:         "",
		selectedAvailableIndex: -1,
	}
}

// SetNoteID sets the note ID to fetch
func (m NoteDetailModel) SetNoteID(id uuid.UUID) (NoteDetailModel, tea.Cmd) {
	// Clear previous state
	m.note = nil
	m.noteID = id
	m.err = nil
	m.tags = nil
	m.tagsErr = nil
	m.links = nil
	m.linksErr = nil
	m.backlinks = nil
	m.backlinksErr = nil
	m.loading = true
	m.currentTab = NoteContentTab // Reset to content tab
	m.showConfirm = false
	// Reset tag management state
	m.selectedTagIndex = -1
	m.availableTags = nil
	m.availableTagsLoading = false
	m.availableTagsErr = nil
	m.showAddTagForm = false
	m.addTagFilter = ""
	m.filteredAvailableTags = nil
	m.selectedAvailableIndex = -1
	return m, m.fetchNoteCmd()
}

// Init initializes the note detail model
func (m NoteDetailModel) Init() tea.Cmd {
	if m.noteID != uuid.Nil {
		return m.fetchNoteCmd()
	}
	return nil
}

// fetchNoteCmd returns a command that fetches the note
func (m NoteDetailModel) fetchNoteCmd() tea.Cmd {
	noteID := m.noteID
	return func() tea.Msg {
		note, err := m.client.GetNote(noteID)
		if err != nil {
			return NoteDetailErrMsg{Err: err}
		}
		return NoteDetailFetchedMsg{Note: note}
	}
}

// fetchTagsCmd returns a command that fetches note tags
func (m NoteDetailModel) fetchTagsCmd() tea.Cmd {
	noteID := m.noteID
	return func() tea.Msg {
		tags, err := m.client.GetNoteTags(noteID)
		if err != nil {
			return NoteDetailTagsErrMsg{Err: err}
		}
		return NoteDetailTagsMsg{Tags: tags}
	}
}

// fetchLinksCmd returns a command that fetches note links
func (m NoteDetailModel) fetchLinksCmd() tea.Cmd {
	noteID := m.noteID
	return func() tea.Msg {
		links, err := m.client.GetLinks(noteID)
		if err != nil {
			return NoteDetailLinksErrMsg{Err: err}
		}
		return NoteDetailLinksMsg{Links: links}
	}
}

// fetchBacklinksCmd returns a command that fetches note backlinks
func (m NoteDetailModel) fetchBacklinksCmd() tea.Cmd {
	noteID := m.noteID
	return func() tea.Msg {
		backlinks, err := m.client.GetBacklinks(noteID)
		if err != nil {
			return NoteDetailBacklinksErrMsg{Err: err}
		}
		return NoteDetailBacklinksMsg{Backlinks: backlinks}
	}
}

// fetchTagsCmdWithID returns a command that fetches note tags with a specific ID
func (m NoteDetailModel) fetchTagsCmdWithID(noteID uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		tags, err := m.client.GetNoteTags(noteID)
		if err != nil {
			return NoteDetailTagsErrMsg{Err: err}
		}
		return NoteDetailTagsMsg{Tags: tags}
	}
}

// fetchLinksCmdWithID returns a command that fetches note links with a specific ID
func (m NoteDetailModel) fetchLinksCmdWithID(noteID uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		links, err := m.client.GetLinks(noteID)
		if err != nil {
			return NoteDetailLinksErrMsg{Err: err}
		}
		return NoteDetailLinksMsg{Links: links}
	}
}

// fetchBacklinksCmdWithID returns a command that fetches note backlinks with a specific ID
func (m NoteDetailModel) fetchBacklinksCmdWithID(noteID uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		backlinks, err := m.client.GetBacklinks(noteID)
		if err != nil {
			return NoteDetailBacklinksErrMsg{Err: err}
		}
		return NoteDetailBacklinksMsg{Backlinks: backlinks}
	}
}

// fetchAllAvailableTagsCmd returns a command that fetches all available tags
func (m NoteDetailModel) fetchAllAvailableTagsCmd() tea.Cmd {
	return func() tea.Msg {
		tags, err := m.client.GetTags()
		if err != nil {
			return NoteAvailableTagsErrMsg{Err: err}
		}
		return NoteAvailableTagsMsg{Tags: tags}
	}
}

// addTagToNoteCmd returns a command that adds a tag to the note
func (m NoteDetailModel) addTagToNoteCmd(tagID uuid.UUID) tea.Cmd {
	// Capture noteID to prevent closure issues with model copying
	noteID := m.noteID
	return func() tea.Msg {
		if err := m.client.AddTagToNote(noteID, tagID); err != nil {
			return NoteDetailTagsErrMsg{Err: err}
		}
		return NoteTagAddedMsg{TagID: tagID}
	}
}

// removeTagFromNoteCmd returns a command that removes a tag from the note
func (m NoteDetailModel) removeTagFromNoteCmd(tagID uuid.UUID) tea.Cmd {
	// Capture noteID to prevent closure issues with model copying
	noteID := m.noteID
	return func() tea.Msg {
		if err := m.client.RemoveTagFromNote(noteID, tagID); err != nil {
			return NoteDetailTagsErrMsg{Err: err}
		}
		return NoteTagRemovedMsg{TagID: tagID}
	}
}

// Update handles messages for the note detail model
func (m NoteDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle confirmation dialog first
		if m.showConfirm {
			m.confirmDialog.Update(msg)
			// Check if user confirmed
			if m.confirmDialog.IsYesSelected() {
				return m, m.deleteNoteCmd()
			} else if m.confirmDialog.IsNoSelected() || msg.String() == "esc" {
				m.showConfirm = false
				return m, nil
			}
			return m, nil
		}

		// Handle add tag form
		if m.showAddTagForm {
			switch msg.String() {
			case "enter":
				// Check if we have a selected tag from the list
				if m.selectedAvailableIndex >= 0 && m.selectedAvailableIndex < len(m.filteredAvailableTags) {
					// Add the selected tag
					tagID := m.filteredAvailableTags[m.selectedAvailableIndex].ID
					m.showAddTagForm = false
					m.addTagInput.SetValue("")
					m.addTagInput.Blur()
					m.addTagFilter = ""
					m.filteredAvailableTags = nil
					m.selectedAvailableIndex = -1
					return m, m.addTagToNoteCmd(tagID)
				} else if m.addTagInput.Value() != "" {
					// Try to find tag by name
					tagName := m.addTagInput.Value()
					for _, tag := range m.availableTags {
						if tag.Name == tagName {
							m.showAddTagForm = false
							m.addTagInput.SetValue("")
							m.addTagInput.Blur()
							m.addTagFilter = ""
							m.filteredAvailableTags = nil
							m.selectedAvailableIndex = -1
							return m, m.addTagToNoteCmd(tag.ID)
						}
					}
					// Tag not found, show error (for now, just close the form)
					m.showAddTagForm = false
					m.addTagInput.SetValue("")
					m.addTagInput.Blur()
					return m, nil
				}
			case "esc":
				m.showAddTagForm = false
				m.addTagInput.SetValue("")
				m.addTagInput.Blur()
				m.addTagFilter = ""
				m.filteredAvailableTags = nil
				m.selectedAvailableIndex = -1
				return m, nil
			case "up", "k":
				// Navigate up in available tags list
				if m.selectedAvailableIndex > 0 {
					m.selectedAvailableIndex--
				} else if len(m.filteredAvailableTags) > 0 {
					m.selectedAvailableIndex = len(m.filteredAvailableTags) - 1
				}
				return m, nil
			case "down", "j":
				// Navigate down in available tags list
				if m.selectedAvailableIndex < len(m.filteredAvailableTags)-1 {
					m.selectedAvailableIndex++
				} else if len(m.filteredAvailableTags) > 0 {
					m.selectedAvailableIndex = 0
				}
				return m, nil
			}
			// Update the input field
			cmd := m.addTagInput.Update(msg)
			// Update filter based on input
			m.addTagFilter = m.addTagInput.Value()
			m.updateFilteredAvailableTags()
			// Reset selection when typing
			if len(m.filteredAvailableTags) > 0 {
				m.selectedAvailableIndex = 0
			} else {
				m.selectedAvailableIndex = -1
			}
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			return m, func() tea.Msg {
				return ShowHelpMsg{}
			}
		case "esc":
			// Go back to dashboard
			return m, func() tea.Msg {
				return ShowDashboardMsg{}
			}
		case "e":
			// Edit note - Phase C
			return m, func() tea.Msg {
				return EditNoteMsg{NoteID: m.noteID}
			}
		case "d":
			// Delete behavior depends on current tab
			if m.currentTab == NoteTagsTab && m.selectedTagIndex >= 0 && len(m.tags) > 0 {
				// Remove the selected tag
				tagID := m.tags[m.selectedTagIndex].ID
				m.selectedTagIndex = -1
				return m, m.removeTagFromNoteCmd(tagID)
			} else {
				// Delete note - show confirmation
				m.showConfirm = true
				m.confirmDialog = components.NewConfirmDialog("Delete this note?")
				m.confirmDialog.SetSubtext("This action cannot be undone.")
				m.confirmDialog.Focus()
				return m, nil
			}
		case "a":
			// Add tag - only works in tags tab
			if m.currentTab == NoteTagsTab {
				m.showAddTagForm = true
				m.addTagInput.Focus()
				m.selectedAvailableIndex = -1
				// Fetch available tags if not already loaded
				if len(m.availableTags) == 0 && !m.availableTagsLoading {
					m.availableTagsLoading = true
					return m, m.fetchAllAvailableTagsCmd()
				}
				return m, nil
			}
		case "tab", "l", "right":
			// Next tab
			m.currentTab = (m.currentTab + 1) % 4
			// Reset tag selection when switching tabs
			if m.currentTab != NoteTagsTab {
				m.selectedTagIndex = -1
			}
			// Fetch data for the tab if needed
			switch m.currentTab {
			case NoteTagsTab:
				if len(m.tags) == 0 && !m.loading {
					cmds = append(cmds, m.fetchTagsCmd())
				}
			case NoteLinksTab:
				if len(m.links) == 0 && !m.loading {
					cmds = append(cmds, m.fetchLinksCmd())
				}
			case NoteBacklinksTab:
				if len(m.backlinks) == 0 && !m.loading {
					cmds = append(cmds, m.fetchBacklinksCmd())
				}
			}
		case "shift+tab", "h", "left":
			// Previous tab
			m.currentTab = (m.currentTab - 1 + 4) % 4
			// Reset tag selection when switching tabs
			if m.currentTab != NoteTagsTab {
				m.selectedTagIndex = -1
			}
		case "up", "k":
			// Navigate up in tags list (only in tags tab)
			if m.currentTab == NoteTagsTab && m.selectedTagIndex > 0 {
				m.selectedTagIndex--
			} else if m.currentTab == NoteTagsTab && m.selectedTagIndex == -1 && len(m.tags) > 0 {
				m.selectedTagIndex = len(m.tags) - 1
			}
		case "down", "j":
			// Navigate down in tags list (only in tags tab)
			if m.currentTab == NoteTagsTab && m.selectedTagIndex < len(m.tags)-1 {
				m.selectedTagIndex++
			} else if m.currentTab == NoteTagsTab && m.selectedTagIndex == -1 && len(m.tags) > 0 {
				m.selectedTagIndex = 0
			}
		}

	case NoteDetailFetchedMsg:
		m.note = msg.Note
		m.loading = false
		// FIX: Use note ID from fetched note to ensure it's valid
		// Capture in local variable to avoid closure issues
		noteID := msg.Note.ID
		return m, tea.Batch(
			m.fetchTagsCmdWithID(noteID),
			m.fetchLinksCmdWithID(noteID),
			m.fetchBacklinksCmdWithID(noteID),
		)

	case NoteDetailTagsMsg:
		m.tags = msg.Tags
		m.tagsErr = nil // Clear error on success
		// Reset selection if out of bounds
		if m.selectedTagIndex >= len(m.tags) {
			m.selectedTagIndex = len(m.tags) - 1
		}
		return m, nil

	case NoteDetailLinksMsg:
		m.links = msg.Links
		m.linksErr = nil // Clear error on success
		return m, nil

	case NoteDetailBacklinksMsg:
		m.backlinks = msg.Backlinks
		m.backlinksErr = nil // Clear error on success
		return m, nil

	case NoteAvailableTagsMsg:
		m.availableTags = msg.Tags
		m.availableTagsLoading = false
		m.availableTagsErr = nil
		m.updateFilteredAvailableTags()
		return m, nil

	case NoteTagAddedMsg:
		// Refresh tags after adding
		return m, m.fetchTagsCmd()

	case NoteTagRemovedMsg:
		// Refresh tags after removing
		return m, m.fetchTagsCmd()

	case NoteDetailErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case NoteDetailTagsErrMsg:
		m.tagsErr = msg.Err
		m.loading = false
		return m, nil

	case NoteAvailableTagsErrMsg:
		m.availableTagsErr = msg.Err
		m.availableTagsLoading = false
		return m, nil

	case NoteDetailLinksErrMsg:
		m.linksErr = msg.Err
		m.loading = false
		return m, nil

	case NoteDetailBacklinksErrMsg:
		m.backlinksErr = msg.Err
		m.loading = false
		return m, nil

	case NoteDeletedMsg:
		// Note was deleted, go back to dashboard
		return m, func() tea.Msg {
			return ShowDashboardMsg{}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.addTagInput.SetWidth(msg.Width - 20)
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// updateFilteredAvailableTags updates the filtered list of available tags
func (m *NoteDetailModel) updateFilteredAvailableTags() {
	if m.addTagFilter == "" {
		// Show all available tags that aren't already on the note
		m.filteredAvailableTags = nil
		for _, availableTag := range m.availableTags {
			alreadyOnNote := false
			for _, noteTag := range m.tags {
				if noteTag.ID == availableTag.ID {
					alreadyOnNote = true
					break
				}
			}
			if !alreadyOnNote {
				m.filteredAvailableTags = append(m.filteredAvailableTags, availableTag)
			}
		}
	} else {
		// Filter by name and exclude tags already on the note
		m.filteredAvailableTags = nil
		for _, availableTag := range m.availableTags {
			// Check if tag matches filter
			if !containsIgnoreCase(availableTag.Name, m.addTagFilter) {
				continue
			}
			// Check if tag is already on the note
			alreadyOnNote := false
			for _, noteTag := range m.tags {
				if noteTag.ID == availableTag.ID {
					alreadyOnNote = true
					break
				}
			}
			if !alreadyOnNote {
				m.filteredAvailableTags = append(m.filteredAvailableTags, availableTag)
			}
		}
	}
	// Reset selection if out of bounds
	if m.selectedAvailableIndex >= len(m.filteredAvailableTags) {
		if len(m.filteredAvailableTags) > 0 {
			m.selectedAvailableIndex = len(m.filteredAvailableTags) - 1
		} else {
			m.selectedAvailableIndex = -1
		}
	}
}

// containsIgnoreCase checks if a string contains a substring (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// IsInputFocused returns whether the add tag form is focused
// This allows the main TUI to skip global key handlers when typing
func (m NoteDetailModel) IsInputFocused() bool {
	return m.showAddTagForm && m.addTagInput.Focused()
}

// GetCurrentTab returns the current active tab
// This allows the main TUI to check if global keys should be blocked
func (m NoteDetailModel) GetCurrentTab() NoteDetailTab {
	return m.currentTab
}

// deleteNoteCmd returns a command that deletes the note
func (m NoteDetailModel) deleteNoteCmd() tea.Cmd {
	// Use the actual note ID from the fetched note, not the field
	// m.noteID may not be reliable if SetNoteID had issues
	if m.note == nil {
		return func() tea.Msg {
			return NoteDetailErrMsg{Err: fmt.Errorf("no note loaded to delete")}
		}
	}
	noteID := m.note.ID
	return func() tea.Msg {
		if err := m.client.DeleteNote(noteID); err != nil {
			return NoteDetailErrMsg{Err: err}
		}
		return NoteDeletedMsg{}
	}
}

// GetNote returns the current note (for accessing from parent models)
func (m NoteDetailModel) GetNote() *model.Note {
	return m.note
}

// View renders the note detail view
func (m NoteDetailModel) View() string {
	if m.showConfirm {
		// Show confirmation dialog centered
		return "\n\n" + m.renderContent() + "\n\n" + m.confirmDialog.View()
	}

	if m.showAddTagForm {
		// Show add tag form below the content
		return m.renderContent() + "\n\n" + m.renderAddTagForm()
	}

	if m.loading {
		return m.renderLoading()
	}

	// Only show global error if note itself failed to load
	if m.err != nil && m.note == nil {
		return m.renderError()
	}

	return m.renderContent()
}

// renderLoading renders the loading state
func (m NoteDetailModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Loading note...")
}

// renderError renders the error state
func (m NoteDetailModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
}

// renderAddTagForm renders the add tag form
func (m NoteDetailModel) renderAddTagForm() string {
	formStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	var content string
	content += formStyle.Render("Add Tag to Note") + "\n\n"
	content += m.addTagInput.View()
	content += "\n\n"

	// Show available tags
	if m.availableTagsLoading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Gray
			Faint(true)
		content += loadingStyle.Render("Loading available tags...")
	} else if m.availableTagsErr != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Red
			Faint(true)
		content += errorStyle.Render(fmt.Sprintf("Error: %v", m.availableTagsErr))
	} else if len(m.filteredAvailableTags) == 0 {
		mutedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Gray
			Faint(true)
		if m.addTagFilter == "" {
			content += mutedStyle.Render("(No more tags available to add)")
		} else {
			content += mutedStyle.Render("(No matching tags found)")
		}
	} else {
		// Show filtered available tags
		for i, tag := range m.filteredAvailableTags {
			if i == m.selectedAvailableIndex {
				content += selectedStyle.Render("→ #" + tag.Name)
			} else {
				content += tagStyle.Render("  #" + tag.Name)
			}
			content += "\n"
		}
	}

	content += "\n"
	content += hintStyle.Render("↑↓:select Enter:add ESC:cancel")

	return content
}

// renderContent renders the note detail content
func (m NoteDetailModel) renderContent() string {
	if m.note == nil {
		return ""
	}

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	metaStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true).
		MarginBottom(1)

	tabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	tabActiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true).
		Padding(0, 1)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")). // Light text
		MarginTop(1)

	var content string

	// Title and metadata
	content += titleStyle.Render(m.note.Title)
	content += "\n"
	content += metaStyle.Render(m.renderMetadata())
	content += "\n"

	// Tabs
	tabs := []NoteDetailTab{NoteContentTab, NoteTagsTab, NoteLinksTab, NoteBacklinksTab}
	var tabViews []string
	for _, tab := range tabs {
		if tab == m.currentTab {
			tabViews = append(tabViews, tabActiveStyle.Render(tab.String()))
		} else {
			tabViews = append(tabViews, tabStyle.Render(tab.String()))
		}
	}
	content += lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)
	content += "\n"

	// Tab content
	content += contentStyle.Render(m.renderTabContent())

	// Action hints at bottom - dynamic based on current tab
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true).
		MarginTop(1)

	var hints string
	if m.currentTab == NoteTagsTab {
		hints = "a:add tag d:remove tag ↑↓:select TAB:tabs e:edit ESC:back"
	} else {
		hints = "TAB:tabs e:edit d:delete ESC:back"
	}
	content += "\n" + hintStyle.Render(hints)

	return content
}

// renderMetadata renders note metadata
func (m NoteDetailModel) renderMetadata() string {
	if m.note == nil {
		return ""
	}

	noteType := string(m.note.NoteType)
	if noteType == "" {
		noteType = "note"
	}

	var info string
	info += fmt.Sprintf("Type: %s", noteType)
	info += fmt.Sprintf(" | Words: %d", m.note.WordCount)
	info += fmt.Sprintf(" | Created: %s", formatTimeAgo(m.note.CreatedAt))
	if !m.note.UpdatedAt.IsZero() && m.note.UpdatedAt != m.note.CreatedAt {
		info += fmt.Sprintf(" | Updated: %s", formatTimeAgo(m.note.UpdatedAt))
	}
	if m.note.AccessCount > 0 {
		info += fmt.Sprintf(" | Views: %d", m.note.AccessCount)
	}

	return info
}

// renderTabContent renders the content for the current tab
func (m NoteDetailModel) renderTabContent() string {
	if m.note == nil {
		return ""
	}

	switch m.currentTab {
	case NoteContentTab:
		return m.renderContentTab()
	case NoteTagsTab:
		return m.renderTagsTab()
	case NoteLinksTab:
		return m.renderLinksTab()
	case NoteBacklinksTab:
		return m.renderBacklinksTab()
	default:
		return ""
	}
}

// renderContentTab renders the note content
func (m NoteDetailModel) renderContentTab() string {
	if m.note.Content == "" {
		return "(no content)"
	}
	return m.note.Content
}

// renderTagsTab renders the tags tab with interactive selection
func (m NoteDetailModel) renderTagsTab() string {
	// Show empty state first (takes priority over errors)
	if len(m.tags) == 0 {
		mutedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Gray
			Faint(true)
		return mutedStyle.Render("(no tags - press 'a' to add)")
	}

	// Show error if tags fetch failed (only shown if we expected data but got error)
	if m.tagsErr != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Red
			Faint(true)
		return errorStyle.Render(fmt.Sprintf("Error loading tags: %v", m.tagsErr))
	}

	// Show tags with selection
	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	var content string
	for i, tag := range m.tags {
		if i == m.selectedTagIndex {
			content += selectedStyle.Render("#" + tag.Name + " [d=remove]")
		} else {
			content += tagStyle.Render("#" + tag.Name)
		}
		if i < len(m.tags)-1 {
			content += "\n"
		}
	}
	return content
}

// renderLinksTab renders the links tab
func (m NoteDetailModel) renderLinksTab() string {
	// Show empty state first (takes priority over errors)
	if len(m.links) == 0 {
		mutedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Gray
			Faint(true)
		return mutedStyle.Render("(no links from this note)")
	}

	// Show error if links fetch failed (only shown if we expected data but got error)
	if m.linksErr != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Red
			Faint(true)
		return errorStyle.Render(fmt.Sprintf("Error loading links: %v", m.linksErr))
	}

	// Show links
	var content string
	for _, link := range m.links {
		if link.TargetNote != nil {
			content += fmt.Sprintf("→ %s\n", link.TargetNote.Title)
		}
	}
	return content
}

// renderBacklinksTab renders the backlinks tab
func (m NoteDetailModel) renderBacklinksTab() string {
	// Show empty state first (takes priority over errors)
	if len(m.backlinks) == 0 {
		mutedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")). // Gray
			Faint(true)
		return mutedStyle.Render("(no backlinks to this note)")
	}

	// Show error if backlinks fetch failed (only shown if we expected data but got error)
	if m.backlinksErr != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Red
			Faint(true)
		return errorStyle.Render(fmt.Sprintf("Error loading backlinks: %v", m.backlinksErr))
	}

	// Show backlinks
	var content string
	for _, link := range m.backlinks {
		if link.SourceNote != nil {
			content += fmt.Sprintf("← %s\n", link.SourceNote.Title)
		}
	}
	return content
}

// Message types for note detail

type NoteDetailFetchedMsg struct {
	Note *model.Note
}

type NoteDetailTagsMsg struct {
	Tags []*model.Tag
}

type NoteDetailLinksMsg struct {
	Links []*model.LinkDetail
}

type NoteDetailBacklinksMsg struct {
	Backlinks []*model.LinkDetail
}

type NoteDetailErrMsg struct {
	Err error
}

// Per-tab error messages
type NoteDetailTagsErrMsg struct {
	Err error
}

type NoteDetailLinksErrMsg struct {
	Err error
}

type NoteDetailBacklinksErrMsg struct {
	Err error
}

type NoteDeletedMsg struct{}

// Available tags messages
type NoteAvailableTagsMsg struct {
	Tags []*model.Tag
}

type NoteAvailableTagsErrMsg struct {
	Err error
}

type NoteTagAddedMsg struct {
	TagID uuid.UUID
}

type NoteTagRemovedMsg struct {
	TagID uuid.UUID
}

// View request messages
type EditNoteMsg struct {
	NoteID uuid.UUID
}
