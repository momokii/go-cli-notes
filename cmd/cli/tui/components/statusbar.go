package components

import (
	"github.com/charmbracelet/lipgloss"
)

// StatusBar represents the status bar at the bottom of the TUI
type StatusBar struct {
	viewName  string // Current view name (e.g., "Help", "Notes")
	keyHelp   string // Key bindings help text
	userInfo  string // User info to display
	width     int
	showError bool
	errorMsg  string
	showInfo  bool
	infoMsg   string
}

// NewStatusBar creates a new status bar
func NewStatusBar() *StatusBar {
	return &StatusBar{
		viewName:  "Help",
		keyHelp:   "↑↓:scroll q:close esc:close",
		userInfo:  "",
		width:     80,
		showError: false,
		showInfo:  false,
	}
}

// SetViewInfo sets the view name and key help
func (s *StatusBar) SetViewInfo(viewName, keyHelp string) {
	s.viewName = viewName
	s.keyHelp = keyHelp
}

// SetUserInfo sets the user info to display
func (s *StatusBar) SetUserInfo(userInfo string) {
	s.userInfo = userInfo
}

// SetWidth sets the width of the status bar
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// ShowError displays an error message in the status bar
func (s *StatusBar) ShowError(msg string) {
	s.showError = true
	s.errorMsg = msg
}

// ShowInfo displays an info message in the status bar
func (s *StatusBar) ShowInfo(msg string) {
	s.showInfo = true
	s.infoMsg = msg
}

// ClearError clears the error message
func (s *StatusBar) ClearError() {
	s.showError = false
	s.errorMsg = ""
	s.showInfo = false
	s.infoMsg = ""
}

// View renders the status bar
func (s *StatusBar) View() string {
	if s.showError && s.errorMsg != "" {
		// Show error message in red
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8")). // Red
			Bold(true)
		errorMsg := errorStyle.Render("⚠ " + s.errorMsg)
		return renderStatusBar("Error", errorMsg, s.userInfo, s.width)
	}

	if s.showInfo && s.infoMsg != "" {
		// Show info message in green
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1")). // Green
			Bold(true)
		infoMsg := infoStyle.Render("✓ " + s.infoMsg)
		return renderStatusBar("Info", infoMsg, s.userInfo, s.width)
	}

	return renderStatusBar(s.viewName, s.keyHelp, s.userInfo, s.width)
}

// renderStatusBar renders the status bar with the given content
func renderStatusBar(view string, keyHelp string, userInfo string, width int) string {
	// Define colors
	colorPrimary := lipgloss.Color("#89b4fa")   // Blue
	colorMuted := lipgloss.Color("#6c7086")    // Gray
	colorSelected := lipgloss.Color("#313244")  // Selected background

	// Create styles
	leftStyle := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	centerStyle := lipgloss.NewStyle().Foreground(colorMuted)
	rightStyle := lipgloss.NewStyle().Foreground(colorMuted)

	left := leftStyle.Render(view)
	center := centerStyle.Render(keyHelp)
	right := rightStyle.Render(userInfo)

	// Calculate spacing
	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)
	rightWidth := lipgloss.Width(right)

	leftSpacer := (width - leftWidth - centerWidth - rightWidth) / 2
	rightSpacer := width - leftWidth - leftSpacer - centerWidth - rightWidth

	if leftSpacer < 2 {
		leftSpacer = 2
	}
	if rightSpacer < 2 {
		rightSpacer = 2
	}

	// Build status bar
	statusBar := lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		lipgloss.NewStyle().Width(leftSpacer).Render(""),
		center,
		lipgloss.NewStyle().Width(rightSpacer).Render(""),
		right,
	)

	// Apply background style
	bgStyle := lipgloss.NewStyle().
		Foreground(colorMuted).
		Background(colorSelected).
		Padding(0, 1)

	return bgStyle.Width(width).Render(statusBar)
}
