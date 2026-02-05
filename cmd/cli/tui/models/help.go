package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// HelpModel is the model for the help screen
type HelpModel struct {
	viewport viewport.Model
	width    int
	height   int
	// Styles passed from parent
	styles   helpStyles
}

// helpStyles contains the styles needed for help rendering
type helpStyles struct {
	SectionStyle lipgloss.Style
	KeyStyle     lipgloss.Style
	DescStyle    lipgloss.Style
	CodeStyle    lipgloss.Style
}

// NewHelpModel creates a new help model
func NewHelpModel() HelpModel {
	vp := viewport.New(0, 0)

	// Define styles inline to avoid import cycle
	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginTop(1).
		MarginBottom(0)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")). // Green
		Bold(true).
		Width(12)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	codeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")). // Green
		Background(lipgloss.Color("#313244")).  // Selected background
		Padding(0, 1)

	return HelpModel{
		viewport: vp,
		styles: helpStyles{
			SectionStyle: sectionStyle,
			KeyStyle:     keyStyle,
			DescStyle:    descStyle,
			CodeStyle:    codeStyle,
		},
	}
}

// Init initializes the help model
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the help model
func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 3
		m.viewport.GotoTop()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the help screen
func (m HelpModel) View() string {
	content := m.getHelpContent()
	m.viewport.SetContent(content)
	return m.viewport.View()
}

// getHelpContent returns the help text content
func (m HelpModel) getHelpContent() string {
	joinHorizontal := lipgloss.JoinHorizontal

	content := `
` + m.styles.SectionStyle.Render("KEY BINDINGS") + `

` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("q"),
		m.styles.DescStyle.Render("Quit TUI"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("? / F1"),
		m.styles.DescStyle.Render("Show this help screen"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("ESC"),
		m.styles.DescStyle.Render("Go back / Cancel current operation"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("/"),
		m.styles.DescStyle.Render("Quick search (from any view)"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("n"),
		m.styles.DescStyle.Render("Create new note (from any view)"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("Ctrl+C"),
		m.styles.DescStyle.Render("Force quit (no confirmation)"),
	) + `

` + m.styles.SectionStyle.Render("NAVIGATION") + `

` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("h / ←"),
		m.styles.DescStyle.Render("Move left / Previous item"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("j / ↓"),
		m.styles.DescStyle.Render("Move down / Next item"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("k / ↑"),
		m.styles.DescStyle.Render("Move up / Previous item"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("l / →"),
		m.styles.DescStyle.Render("Move right / Select item"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("Enter"),
		m.styles.DescStyle.Render("Open selected item / Confirm"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("0"),
		m.styles.DescStyle.Render("Go to top of list"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("G"),
		m.styles.DescStyle.Render("Go to bottom of list"),
	) + `

` + m.styles.SectionStyle.Render("DASHBOARD") + `

` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("n"),
		m.styles.DescStyle.Render("Create new note"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("s"),
		m.styles.DescStyle.Render("Go to search"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("l"),
		m.styles.DescStyle.Render("View all notes"),
	) + `
` + joinHorizontal(lipgloss.Top,
		m.styles.KeyStyle.Render("a"),
		m.styles.DescStyle.Render("View activity feed"),
	) + `

` + m.styles.SectionStyle.Render("TIPS") + `

• Press ` + m.styles.CodeStyle.Render("?") + ` anytime to see this help
• Key hints are shown in the status bar (bottom of screen)
• Vim navigation (h/j/k/l) works alongside arrow keys
• Session expiry will auto-exit TUI to protect your data
• All changes are saved before session expiry
• Use ` + m.styles.CodeStyle.Render("ESC") + ` to go back from any view
• Use ` + m.styles.CodeStyle.Render("q") + ` to quit TUI from any view

`

	return content
}

// SetSize sets the size of the help viewport
func (m *HelpModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 4
	m.viewport.Height = height - 3
}

// GotoTop scrolls to the top of the help content
func (m *HelpModel) GotoTop() {
	m.viewport.GotoTop()
}
