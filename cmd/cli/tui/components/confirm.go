package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmDialog represents a confirmation dialog
type ConfirmDialog struct {
	question string
	subtext  string
	yesText  string
	noText   string
	selected ConfirmChoice // yes or no
	focused  bool
	width    int
}

// ConfirmChoice represents the user's choice
type ConfirmChoice int

const (
	ConfirmNo ConfirmChoice = iota
	ConfirmYes
)

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(question string) ConfirmDialog {
	return ConfirmDialog{
		question: question,
		subtext:  "",
		yesText:  "Yes",
		noText:   "No",
		selected: ConfirmNo, // Default to No for safety
		focused:  false,
		width:    60,
	}
}

// SetSubtext sets additional context text below the question
func (c *ConfirmDialog) SetSubtext(text string) {
	c.subtext = text
}

// SetYesText sets the text for the Yes option
func (c *ConfirmDialog) SetYesText(text string) {
	c.yesText = text
}

// SetNoText sets the text for the No option
func (c *ConfirmDialog) SetNoText(text string) {
	c.noText = text
}

// SetWidth sets the width of the dialog
func (c *ConfirmDialog) SetWidth(width int) {
	c.width = width
}

// SetSelected sets the current selection
func (c *ConfirmDialog) SetSelected(choice ConfirmChoice) {
	c.selected = choice
}

// Selected returns the current selection
func (c *ConfirmDialog) Selected() ConfirmChoice {
	return c.selected
}

// IsYesSelected returns true if Yes is selected
func (c *ConfirmDialog) IsYesSelected() bool {
	return c.selected == ConfirmYes
}

// IsNoSelected returns true if No is selected
func (c *ConfirmDialog) IsNoSelected() bool {
	return c.selected == ConfirmNo
}

// Focus sets focus on the dialog
func (c *ConfirmDialog) Focus() {
	c.focused = true
}

// Blur removes focus from the dialog
func (c *ConfirmDialog) Blur() {
	c.focused = false
}

// Update handles messages for the confirm dialog
func (c *ConfirmDialog) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !c.focused {
			return nil
		}

		switch msg.String() {
		case "h", "left":
			c.selected = ConfirmNo
		case "l", "right":
			c.selected = ConfirmYes
		case "y", "Y":
			c.selected = ConfirmYes
			return tea.Quit
		case "n", "N":
			c.selected = ConfirmNo
			return tea.Quit
		case "enter", " ":
			return tea.Quit
		case "esc":
			c.selected = ConfirmNo
			return tea.Quit
		case "q":
			c.selected = ConfirmNo
			return tea.Quit
		}
	}

	return nil
}

// View renders the confirm dialog
func (c *ConfirmDialog) View() string {
	// Define styles
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f38ba8")). // Red for warning
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	subtextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		MarginBottom(1)

	yesStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")). // Green
		Bold(true).
		Padding(0, 1)

	noStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Padding(0, 1)

	yesSelectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark background
		Background(lipgloss.Color("#a6e3a1")). // Green
		Bold(true).
		Padding(0, 1)

	noSelectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark background
		Background(lipgloss.Color("#6c7086")). // Gray
		Bold(true).
		Padding(0, 1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Faint(true).
		MarginTop(1)

	// Build content
	content := titleStyle.Render(c.question)

	if c.subtext != "" {
		content += "\n" + subtextStyle.Render(c.subtext)
	}

	content += "\n"

	// Render options
	var yesBtn, noBtn string
	if c.selected == ConfirmYes {
		yesBtn = yesSelectedStyle.Render(c.yesText)
		noBtn = noStyle.Render(c.noText)
	} else {
		yesBtn = yesStyle.Render(c.yesText)
		noBtn = noSelectedStyle.Render(c.noText)
	}

	options := lipgloss.JoinHorizontal(lipgloss.Top, yesBtn, " ", noBtn)
	content += "\n" + options

	content += "\n" + hintStyle.Render("←→:select y:yes n:no enter:confirm")

	return boxStyle.Render(content)
}

// Result returns a ConfirmResult with the user's choice
func (c *ConfirmDialog) Result() ConfirmResult {
	return ConfirmResult{
		Confirmed: c.selected == ConfirmYes,
		Choice:    c.selected,
	}
}

// ConfirmResult represents the result of a confirmation dialog
type ConfirmResult struct {
	Confirmed bool
	Choice    ConfirmChoice
}
