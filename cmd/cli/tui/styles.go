package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for dark theme
var (
	// Base colors
	colorForeground  = lipgloss.Color("#cdd6f4")  // Light text
	colorBackground  = lipgloss.Color("#1e1e2e")  // Dark background
	colorPrimary     = lipgloss.Color("#89b4fa")  // Blue
	colorSecondary   = lipgloss.Color("#fab387")  // Orange
	colorAccent      = lipgloss.Color("#a6e3a1")  // Green
	colorError       = lipgloss.Color("#f38ba8")  // Red
	colorWarning     = lipgloss.Color("#f9e2af")  // Yellow
	colorMuted       = lipgloss.Color("#6c7086")  // Gray
	colorBorder      = lipgloss.Color("#45475a")  // Dark gray border
	colorSelected    = lipgloss.Color("#313244")  // Selected background
	colorHeader      = lipgloss.Color("#181825")  // Header background
)

// Styles for various UI elements

// BaseStyle is the base style for all text
var BaseStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Background(colorBackground)

// PrimaryStyle is for primary actions and emphasis
var PrimaryStyle = lipgloss.NewStyle().
	Foreground(colorBackground).
	Background(colorPrimary).
	Bold(true)

// SecondaryStyle is for secondary emphasis
var SecondaryStyle = lipgloss.NewStyle().
	Foreground(colorSecondary).
	Bold(true)

// AccentStyle is for success states and positive feedback
var AccentStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// ErrorStyle is for errors
var ErrorStyle = lipgloss.NewStyle().
	Foreground(colorError).
	Bold(true)

// WarningStyle is for warnings
var WarningStyle = lipgloss.NewStyle().
	Foreground(colorWarning).
	Bold(true)

// MutedStyle is for secondary text
var MutedStyle = lipgloss.NewStyle().
	Foreground(colorMuted)

// BorderStyle is for borders
var BorderStyle = lipgloss.NewStyle().
	Foreground(colorBorder)

// SelectedStyle is for selected items
var SelectedStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Background(colorSelected).
	Bold(true)

// HeaderStyle is for the top header bar
var HeaderStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Background(colorHeader).
	Bold(true).
	Padding(0, 1)

// StatusBarStyle is for the bottom status bar
var StatusBarStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Background(colorSelected).
	Padding(0, 1)

// KeyBindingStyle is for displaying key bindings
var KeyBindingStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// KeyDescStyle is for key binding descriptions
var KeyDescStyle = lipgloss.NewStyle().
	Foreground(colorMuted)

// TitleStyle is for titles and headings
var TitleStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true).
	MarginBottom(1)

// SubtitleStyle is for subtitles
var SubtitleStyle = lipgloss.NewStyle().
	Foreground(colorSecondary).
	Bold(true).
	MarginBottom(1)

// ItemStyle is for list items
var ItemStyle = lipgloss.NewStyle().
	Foreground(colorForeground).
	Padding(0, 1)

// DimStyle is for dimmed text
var DimStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Faint(true)

// BoxStyle is for drawing boxes/borders
var BoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder).
	Padding(1)

// ModalStyle is for modal dialogs
var ModalStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(colorPrimary).
	Background(colorBackground).
	Padding(2).
	Width(60)

// HelpKeyStyle is for key names in help screen
var HelpKeyStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true).
	Width(12)

// HelpDescStyle is for key descriptions in help screen
var HelpDescStyle = lipgloss.NewStyle().
	Foreground(colorForeground)

// LinkStyle is for clickable links
var LinkStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Underline(true)

// CodeStyle is for inline code
var CodeStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Background(colorSelected).
	Padding(0, 1)

// MetadataStyle is for metadata (dates, counts, etc.)
var MetadataStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Faint(true)

// TagStyle is for tag display
var TagStyle = lipgloss.NewStyle().
	Foreground(colorBackground).
	Background(colorSecondary).
	Bold(true).
	Padding(0, 1)

// ActiveTabStyle is for active tab in tabbed interfaces
var ActiveTabStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true).
	Underline(true)

// InactiveTabStyle is for inactive tabs
var InactiveTabStyle = lipgloss.NewStyle().
	Foreground(colorMuted)

// ScrollbarStyle is for custom scrollbars
var ScrollbarStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Faint(true)

// HelpSectionStyle is for help screen sections
var HelpSectionStyle = lipgloss.NewStyle().
	Foreground(colorSecondary).
	Bold(true).
	MarginTop(1).
	MarginBottom(0)

// SeparatorStyle is for separators
var SeparatorStyle = lipgloss.NewStyle().
	Foreground(colorBorder).
	SetString("─")

// SpinnerStyle is for loading spinners
var SpinnerStyle = lipgloss.NewStyle().
	Foreground(colorAccent)

// SuccessStyle is for success messages
var SuccessStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// InfoStyle is for informational messages
var InfoStyle = lipgloss.NewStyle().
	Foreground(colorPrimary)

// FatalStyle is for fatal errors
var FatalStyle = lipgloss.NewStyle().
	Foreground(colorError).
	Background(colorSelected).
	Bold(true).
	Padding(1).
	Width(80).
	Align(lipgloss.Center)

// StyledKey returns a styled key binding string
func StyledKey(key string) string {
	return KeyBindingStyle.Render(key)
}

// StyledKeyWithDesc returns a styled key binding with description
func StyledKeyWithDesc(key, desc string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		HelpKeyStyle.Render(key),
		HelpDescStyle.Render(desc),
	)
}

// RenderHeader renders the top header bar
func RenderHeader(title, userInfo string, width int) string {
	left := TitleStyle.Render(title)
	right := lipgloss.NewStyle().Foreground(colorMuted).Render(userInfo)

	// Calculate spacing to push right content to the edge
	spacer := width - lipgloss.Width(left) - lipgloss.Width(right)
	if spacer < 0 {
		spacer = 0
	}

	return HeaderStyle.Width(width).Render(lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		lipgloss.NewStyle().Width(spacer).Render(""),
		right,
	))
}

// RenderStatusBar renders the bottom status bar
func RenderStatusBar(view string, keyHelp string, userInfo string, width int) string {
	left := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render(view)
	center := KeyDescStyle.Render(keyHelp)
	right := lipgloss.NewStyle().Foreground(colorMuted).Render(userInfo)

	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)
	rightWidth := lipgloss.Width(right)

	// Calculate spacing
	leftSpacer := (width - leftWidth - centerWidth - rightWidth) / 2
	rightSpacer := width - leftWidth - leftSpacer - centerWidth - rightWidth

	if leftSpacer < 2 {
		leftSpacer = 2
	}
	if rightSpacer < 2 {
		rightSpacer = 2
	}

	return StatusBarStyle.Width(width).Render(lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		lipgloss.NewStyle().Width(leftSpacer).Render(""),
		center,
		lipgloss.NewStyle().Width(rightSpacer).Render(""),
		right,
	))
}

// RenderBox renders content in a bordered box
func RenderBox(content string, title string) string {
	if title != "" {
		// Create a new style with title
		titleStyle := lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)
		titledContent := titleStyle.Render(title) + "\n" + content
		return BoxStyle.Render(titledContent)
	}
	return BoxStyle.Render(content)
}

// RenderModal renders a modal dialog
func RenderModal(content string, title string) string {
	if title != "" {
		// Create a new style with title
		titleStyle := lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)
		titledContent := titleStyle.Render(title) + "\n" + content
		return ModalStyle.Render(titledContent)
	}
	return ModalStyle.Render(content)
}

// RenderError renders an error message
func RenderError(err string) string {
	return ErrorStyle.Render("Error: " + err)
}

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return SuccessStyle.Render("✓ " + msg)
}

// RenderWarning renders a warning message
func RenderWarning(msg string) string {
	return WarningStyle.Render("⚠ " + msg)
}

// RenderInfo renders an info message
func RenderInfo(msg string) string {
	return InfoStyle.Render("ℹ " + msg)
}

// TruncateText truncates text to fit within max width
func TruncateText(text string, maxWidth int) string {
	if lipgloss.Width(text) <= maxWidth {
		return text
	}
	return lipgloss.NewStyle().MaxWidth(maxWidth).Render(text) + "…"
}

// GetViewportWidth returns the width for a viewport given total width and margins
func GetViewportWidth(totalWidth, margins int) int {
	width := totalWidth - (margins * 2)
	if width < 20 {
		return 20
	}
	return width
}

// GetViewportHeight returns the height for a viewport given total height and header/footer usage
func GetViewportHeight(totalHeight, headerHeight, footerHeight, margins int) int {
	height := totalHeight - headerHeight - footerHeight - (margins * 2)
	if height < 10 {
		return 10
	}
	return height
}
