package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents a single field in a form
type FormField struct {
	ID          string
	Label       string
	InputType   FieldType // TextInput or Textarea
	textInput   TextInput
	textarea    Textarea
	Required    bool
	MinLength   int
	MaxLength   int
	MinWords    int // For textarea only
	Placeholder string
	Error       string
}

// FieldType represents the type of input field
type FieldType int

const (
	FieldInput FieldType = iota
	FieldTextarea
)

// NewFormField creates a new form field
func NewFormField(id, label string, fieldType FieldType) FormField {
	f := FormField{
		ID:        id,
		Label:     label,
		InputType: fieldType,
		Required:  false,
	}

	if fieldType == FieldInput {
		f.textInput = NewTextInput()
		f.textInput.SetPrompt(label + ": ")
	} else {
		f.textarea = NewTextarea()
		f.textarea.SetPlaceholder(label)
	}

	return f
}

// SetValue sets the value of the field
func (f *FormField) SetValue(value string) {
	if f.InputType == FieldInput {
		f.textInput.SetValue(value)
	} else {
		f.textarea.SetValue(value)
	}
}

// Value returns the current value of the field
func (f *FormField) Value() string {
	if f.InputType == FieldInput {
		return f.textInput.Value()
	}
	return f.textarea.Value()
}

// Focus sets focus on this field
func (f *FormField) Focus() {
	if f.InputType == FieldInput {
		f.textInput.Focus()
	} else {
		f.textarea.Focus()
	}
}

// Blur removes focus from this field
func (f *FormField) Blur() {
	if f.InputType == FieldInput {
		f.textInput.Blur()
	} else {
		f.textarea.Blur()
	}
}

// Focused returns whether this field is focused
func (f *FormField) Focused() bool {
	if f.InputType == FieldInput {
		return f.textInput.Focused()
	}
	return f.textarea.Focused()
}

// Validate validates the field value
func (f *FormField) Validate() error {
	if f.InputType == FieldInput {
		return f.textInput.Validate(f.Required, f.MinLength, f.MaxLength)
	}
	return f.textarea.Validate(f.Required, f.MinWords)
}

// SetPlaceholder sets the placeholder for the field
func (f *FormField) SetPlaceholder(text string) {
	if f.InputType == FieldInput {
		f.textInput.SetPlaceholder(text)
	} else {
		f.textarea.SetPlaceholder(text)
	}
}

// Form represents a multi-field form
type Form struct {
	fields     []FormField
	currentIdx int
	width      int
	focused    bool
	submitText string
	cancelText string
}

// NewForm creates a new form
func NewForm() Form {
	return Form{
		fields:     []FormField{},
		currentIdx: 0,
		width:      60,
		focused:    false,
		submitText: "Submit",
		cancelText: "Cancel",
	}
}

// AddField adds a field to the form
func (f *Form) AddField(field FormField) {
	f.fields = append(f.fields, field)
}

// SetFields sets all fields at once
func (f *Form) SetFields(fields []FormField) {
	f.fields = fields
}

// Fields returns the form fields
func (f *Form) Fields() []FormField {
	return f.fields
}

// CurrentField returns the currently focused field
func (f *Form) CurrentField() *FormField {
	if len(f.fields) == 0 {
		return nil
	}
	return &f.fields[f.currentIdx]
}

// SetCurrentIndex sets the current field index
func (f *Form) SetCurrentIndex(idx int) {
	if idx >= 0 && idx < len(f.fields) {
		// Blur current field
		if f.Focused() {
			f.fields[f.currentIdx].Blur()
		}
		f.currentIdx = idx
		// Focus new field
		if f.Focused() {
			f.fields[f.currentIdx].Focus()
		}
	}
}

// CurrentIndex returns the current field index
func (f *Form) CurrentIndex() int {
	return f.currentIdx
}

// SetWidth sets the form width
func (f *Form) SetWidth(width int) {
	f.width = width
	for _, field := range f.fields {
		if field.InputType == FieldInput {
			field.textInput.SetWidth(width - 20) // Account for label
		} else {
			field.textarea.SetWidth(width)
		}
	}
}

// SetSubmitText sets the submit button text
func (f *Form) SetSubmitText(text string) {
	f.submitText = text
}

// SetCancelText sets the cancel button text
func (f *Form) SetCancelText(text string) {
	f.cancelText = text
}

// Focus sets focus on the form
func (f *Form) Focus() {
	f.focused = true
	if len(f.fields) > 0 {
		f.fields[f.currentIdx].Focus()
	}
}

// Blur removes focus from the form
func (f *Form) Blur() {
	f.focused = false
	if len(f.fields) > 0 {
		f.fields[f.currentIdx].Blur()
	}
}

// Focused returns whether the form is focused
func (f *Form) Focused() bool {
	return f.focused
}

// Validate validates all fields in the form
func (f *Form) Validate() map[string]error {
	errors := make(map[string]error)
	for i := range f.fields {
		if err := f.fields[i].Validate(); err != nil {
			errors[f.fields[i].ID] = err
		}
	}
	return errors
}

// Values returns a map of field IDs to their values
func (f *Form) Values() map[string]string {
	values := make(map[string]string)
	for i := range f.fields {
		values[f.fields[i].ID] = f.fields[i].Value()
	}
	return values
}

// Update handles messages for the form
func (f *Form) Update(msg tea.Msg) tea.Cmd {
	if !f.focused || len(f.fields) == 0 {
		return nil
	}

	// Handle keyboard navigation between fields
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			// Move to next field
			if f.currentIdx < len(f.fields)-1 {
				f.SetCurrentIndex(f.currentIdx + 1)
			}
			return nil
		case "shift+tab", "up":
			// Move to previous field
			if f.currentIdx > 0 {
				f.SetCurrentIndex(f.currentIdx - 1)
			}
			return nil
		}
	}

	// Update current field
	field := &f.fields[f.currentIdx]
	if field.InputType == FieldInput {
		return field.textInput.Update(msg)
	}
	return field.textarea.Update(msg)
}

// View renders the form
func (f *Form) View() string {
	if len(f.fields) == 0 {
		return ""
	}

	var content string

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Italic(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	for _, field := range f.fields {
		// Render field
		if field.InputType == FieldInput {
			content += field.textInput.View()
		} else {
			// For textarea, show label above it
			label := labelStyle.Render(field.Label)
			if field.Required {
				label += labelStyle.Render("*")
			}
			content += label + "\n"
			content += field.textarea.View()
		}

		// Show error if any
		if field.Error != "" {
			content += "\n" + errorStyle.Render("⚠ "+field.Error)
		}

		content += "\n\n"
	}

	// Add hints at the bottom
	hints := "TAB:next Shift+TAB:prev"
	if f.submitText != "" {
		hints += " Enter:" + f.submitText
	}
	if f.cancelText != "" {
		hints += " ESC:" + f.cancelText
	}
	content += hintStyle.Render(hints)

	return content
}
