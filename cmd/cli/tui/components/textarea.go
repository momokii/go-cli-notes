package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Textarea wraps the Bubbles textarea component for multi-line input
type Textarea struct {
	textarea textarea.Model
	focused  bool
	width    int
	height   int
}

// NewTextarea creates a new textarea component
func NewTextarea() Textarea {
	ta := textarea.New()
	ta.Placeholder = "Enter content here..."
	ta.ShowLineNumbers = false
	ta.CharLimit = 0 // No limit by default
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("#313244"))
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))

	return Textarea{
		textarea: ta,
		focused:  false,
		width:    60,
		height:   10,
	}
}

// SetPlaceholder sets the placeholder text
func (t *Textarea) SetPlaceholder(text string) {
	t.textarea.Placeholder = text
}

// SetValue sets the current value
func (t *Textarea) SetValue(value string) {
	t.textarea.SetValue(value)
}

// Value returns the current value
func (t *Textarea) Value() string {
	return t.textarea.Value()
}

// SetWidth sets the width of the textarea
func (t *Textarea) SetWidth(width int) {
	t.width = width
	t.textarea.SetWidth(width)
}

// SetHeight sets the height of the textarea
func (t *Textarea) SetHeight(height int) {
	t.height = height
	t.textarea.SetHeight(height)
}

// SetSize sets both width and height
func (t *Textarea) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.textarea.SetWidth(width)
	t.textarea.SetHeight(height)
}

// SetCharLimit sets the maximum number of characters (0 = no limit)
func (t *Textarea) SetCharLimit(limit int) {
	t.textarea.CharLimit = limit
}

// Focus sets focus on the textarea
func (t *Textarea) Focus() {
	t.focused = true
	t.textarea.Focus()
	// Update focused styles
	t.textarea.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("#313244"))
}

// Blur removes focus from the textarea
func (t *Textarea) Blur() {
	t.focused = false
	t.textarea.Blur()
}

// Focused returns whether the textarea is focused
func (t *Textarea) Focused() bool {
	return t.focused
}

// Reset clears the textarea value
func (t *Textarea) Reset() {
	t.textarea.Reset()
}

// WordCount returns the number of words in the textarea
func (t *Textarea) WordCount() int {
	// Simple word count implementation
	value := t.Value()
	if value == "" {
		return 0
	}
	count := 0
	inWord := false
	for _, c := range value {
		if c == ' ' || c == '\t' || c == '\n' {
			if inWord {
				count++
				inWord = false
			}
		} else {
			inWord = true
		}
	}
	if inWord {
		count++
	}
	return count
}

// LineCount returns the number of lines in the textarea
func (t *Textarea) LineCount() int {
	return len(t.textarea.Value())
}

// Update handles messages for the textarea
func (t *Textarea) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.textarea, cmd = t.textarea.Update(msg)
	return cmd
}

// View renders the textarea
func (t *Textarea) View() string {
	return t.textarea.View()
}

// Validate checks if the textarea content is valid
func (t *Textarea) Validate(required bool, minWords int) error {
	value := t.Value()
	if required && value == "" {
		return &ValidationError{Field: "content", Message: "Content is required"}
	}
	if minWords > 0 && t.WordCount() < minWords {
		return &ValidationError{Field: "content", Message: "Too few words"}
	}
	return nil
}

// SetCursor moves the cursor to a specific position (character offset)
func (t *Textarea) SetCursor(offset int) {
	t.textarea.SetCursor(offset)
}
