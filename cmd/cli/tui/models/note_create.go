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

// NoteCreateMode represents whether we're creating or editing
type NoteCreateMode int

const (
	ModeCreate NoteCreateMode = iota
	ModeEdit
)

// NoteCreateModel is the model for creating/editing notes
type NoteCreateModel struct {
	client     *client.APIClient
	authState  *client.AuthState
	mode       NoteCreateMode
	noteID     uuid.UUID // For edit mode
	form       components.Form
	loading    bool
	err        error
	saved      bool
	hasChanges bool // Track unsaved changes
	width      int
	height     int
}

// NewNoteCreateModel creates a new note create model
func NewNoteCreateModel(apiClient *client.APIClient, authState *client.AuthState) NoteCreateModel {
	form := components.NewForm()
	form.SetSubmitText("Save")
	form.SetCancelText("Cancel")

	// Title field
	titleField := components.NewFormField("title", "Title", components.FieldInput)
	titleField.Required = true
	titleField.MinLength = 1
	titleField.MaxLength = 500
	titleField.SetPlaceholder("Enter note title...")
	form.AddField(titleField)

	// Content field
	contentField := components.NewFormField("content", "Content", components.FieldTextarea)
	contentField.Required = true
	contentField.MinWords = 1
	contentField.MaxLength = 100000
	contentField.SetPlaceholder("Enter note content...")
	form.AddField(contentField)

	return NoteCreateModel{
		client:    apiClient,
		authState: authState,
		mode:      ModeCreate,
		noteID:    uuid.Nil,
		form:      form,
		loading:   false,
		width:     80,
		height:    24,
	}
}

// SetEditMode sets the model to edit mode with existing note data
// Returns the modified model (value receiver pattern)
func (m NoteCreateModel) SetEditMode(note *model.Note) (NoteCreateModel, tea.Cmd) {
	m.mode = ModeEdit
	m.noteID = note.ID
	m.form.Fields()[0].SetValue(note.Title)    // Title
	m.form.Fields()[1].SetValue(note.Content)  // Content
	m.form.SetSubmitText("Update")
	m.hasChanges = false
	// Focus the form so user can edit
	m = m.FocusForm()
	return m, nil
}

// Init initializes the note create model
func (m NoteCreateModel) Init() tea.Cmd {
	m.form.Focus()
	return nil
}

// FocusForm focuses the form and returns the updated model
// This is needed because models are stored as values in the main model
func (m NoteCreateModel) FocusForm() NoteCreateModel {
	m.form.Focus()
	return m
}

// Update handles messages for the note create model
func (m NoteCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle form submission
		if msg.String() == "enter" && m.form.Focused() {
			// Validate and submit
			if err := m.validateForm(); err != nil {
				m.form.Fields()[m.form.CurrentIndex()].Error = err.Error()
				return m, nil
			}

			// Submit form
			if m.mode == ModeCreate {
				return m, m.createNoteCmd()
			}
			return m, m.updateNoteCmd()
		}

		// ESC to cancel - always allow exiting (discards unsaved changes)
		if msg.String() == "esc" {
			if m.hasChanges && !m.saved {
				// TODO: Phase E - add proper unsaved changes dialog
				// For now, just discard changes and go back
			}
			return m, func() tea.Msg {
				return ShowDashboardMsg{}
			}
		}

		// Track changes
		if m.form.Focused() && (msg.String() == "backspace" || len(msg.String()) == 1) {
			m.hasChanges = true
		}

	case NoteCreatedMsg:
		m.saved = true
		m.loading = false
		m.hasChanges = false // RESET: Allow ESC to work
		m.form.Blur()        // FIX: Remove focus so global keys work
		return m, func() tea.Msg {
			return ShowDashboardMsg{}
		}

	case NoteUpdatedMsg:
		m.saved = true
		m.loading = false
		m.hasChanges = false // RESET: Allow ESC to work
		m.form.Blur()        // FIX: Remove focus so global keys work
		return m, func() tea.Msg {
			return OpenNoteMsg{NoteID: m.noteID}
		}

	case NoteCreateErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.form.SetWidth(msg.Width - 10) // Leave margin
		return m, nil
	}

	// Update form
	cmd := m.form.Update(msg)
	return m, cmd
}

// validateForm validates the form fields
func (m NoteCreateModel) validateForm() error {
	values := m.form.Values()

	// Validate title
	title := values["title"]
	if title == "" {
		return &ValidationError{Field: "title", Message: "Title is required"}
	}
	if len(title) > 500 {
		return &ValidationError{Field: "title", Message: "Title too long (max 500)"}
	}

	// Validate content
	content := values["content"]
	if content == "" {
		return &ValidationError{Field: "content", Message: "Content is required"}
	}
	if len(content) > 100000 {
		return &ValidationError{Field: "content", Message: "Content too long"}
	}

	return nil
}

// createNoteCmd returns a command that creates a new note
func (m NoteCreateModel) createNoteCmd() tea.Cmd {
	m.loading = true
	values := m.form.Values()

	return func() tea.Msg {
		req := &model.CreateNoteRequest{
			Title:    values["title"],
			Content:  values["content"],
			NoteType: "note", // Default to note type
		}

		note, err := m.client.CreateNote(req)
		if err != nil {
			return NoteCreateErrMsg{Err: err}
		}

		return NoteCreatedMsg{NoteID: note.ID}
	}
}

// updateNoteCmd returns a command that updates an existing note
func (m NoteCreateModel) updateNoteCmd() tea.Cmd {
	m.loading = true
	values := m.form.Values()

	return func() tea.Msg {
		title := values["title"]
		content := values["content"]

		req := &model.UpdateNoteRequest{
			Title:   &title,
			Content: &content,
		}

		if err := m.client.UpdateNote(m.noteID, req); err != nil {
			return NoteCreateErrMsg{Err: err}
		}

		return NoteUpdatedMsg{NoteID: m.noteID}
	}
}

// View renders the note create/edit view
func (m NoteCreateModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	var content string

	// Header
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true)

	if m.mode == ModeCreate {
		content += titleStyle.Render("CREATE NEW NOTE")
	} else {
		content += titleStyle.Render("EDIT NOTE")
	}

	if m.hasChanges && !m.saved {
		content += " " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Red
			Render("(unsaved)")
	}

	content += "\n\n"

	// Form
	content += m.form.View()

	return content
}

// renderLoading renders the loading state
func (m NoteCreateModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	if m.mode == ModeCreate {
		return style.Render("Creating note...")
	}
	return style.Render("Updating note...")
}

// renderError renders the error state
func (m NoteCreateModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render(fmt.Sprintf("Error: %v\n\nESC to go back", m.err))
}

// IsInputFocused returns whether the form has focused input
// FIX: Check if form is focused and has a focused field
func (m NoteCreateModel) IsInputFocused() bool {
	// Check if form itself is focused
	if !m.form.Focused() {
		return false
	}
	// Also check if any field is focused
	for _, field := range m.form.Fields() {
		if field.Focused() {
			return true
		}
	}
	return false
}

// BlurForm removes focus from the form
func (m NoteCreateModel) BlurForm() NoteCreateModel {
	m.form.Blur()
	return m
}

// Message types for note create/edit

type NoteCreatedMsg struct {
	NoteID uuid.UUID
}

type NoteUpdatedMsg struct {
	NoteID uuid.UUID
}

type NoteCreateErrMsg struct {
	Err error
}

// ValidationError represents a form validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
