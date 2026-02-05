package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextInput wraps the Bubbles textinput component
type TextInput struct {
	textInput textinput.Model
	focused   bool
	width     int
}

// NewTextInput creates a new text input component
func NewTextInput() TextInput {
	ti := textinput.New()
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Faint(true)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa"))
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))

	return TextInput{
		textInput: ti,
		focused:   false,
		width:     40,
	}
}

// SetPlaceholder sets the placeholder text
func (t *TextInput) SetPlaceholder(text string) {
	t.textInput.Placeholder = text
}

// SetPrompt sets the prompt text (label shown before input)
func (t *TextInput) SetPrompt(prompt string) {
	t.textInput.Prompt = prompt
}

// SetValue sets the current value
func (t *TextInput) SetValue(value string) {
	t.textInput.SetValue(value)
}

// Value returns the current value
func (t *TextInput) Value() string {
	return t.textInput.Value()
}

// SetWidth sets the width of the text input
func (t *TextInput) SetWidth(width int) {
	t.width = width
	t.textInput.Width = width
}

// SetEchoMode sets the echo mode (normal, password, etc.)
func (t *TextInput) SetEchoMode(mode textinput.EchoMode) {
	t.textInput.EchoMode = mode
}

// SetCharLimit sets the maximum number of characters
func (t *TextInput) SetCharLimit(limit int) {
	t.textInput.CharLimit = limit
}

// Focus sets focus on the text input
func (t *TextInput) Focus() {
	t.focused = true
	t.textInput.Focus()
	// Update cursor style when focused
	t.textInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa"))
}

// Blur removes focus from the text input
func (t *TextInput) Blur() {
	t.focused = false
	t.textInput.Blur()
}

// Focused returns whether the text input is focused
func (t *TextInput) Focused() bool {
	return t.focused
}

// Reset clears the input value
func (t *TextInput) Reset() {
	t.textInput.Reset()
}

// Update handles messages for the text input
func (t *TextInput) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.textInput, cmd = t.textInput.Update(msg)
	return cmd
}

// View renders the text input
func (t *TextInput) View() string {
	return t.textInput.View()
}

// Error returns the current error state (if any)
func (t *TextInput) Error() string {
	// Bubbles textinput doesn't have built-in error state
	// Errors are handled at the form level
	return ""
}

// Validate checks if the input is valid
func (t *TextInput) Validate(required bool, minLength, maxLength int) error {
	value := t.Value()
	if required && value == "" {
		return &ValidationError{Field: "input", Message: "This field is required"}
	}
	if minLength > 0 && len(value) < minLength {
		return &ValidationError{Field: "input", Message: "Too short"}
	}
	if maxLength > 0 && len(value) > maxLength {
		return &ValidationError{Field: "input", Message: "Too long"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
