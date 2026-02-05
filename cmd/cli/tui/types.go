// Package tui provides the Terminal User Interface for kg-cli
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// View represents the different screens/views in the TUI
type View int

const (
	// DashboardView shows statistics, recent activity, and quick actions
	DashboardView View = iota
	// NoteListView displays a paginated list of notes
	NoteListView
	// NoteDetailView shows a single note with edit capabilities
	NoteDetailView
	// NoteCreateView is the form for creating a new note
	NoteCreateView
	// NoteEditView is the form for editing an existing note
	NoteEditView
	// TagListView displays a list of all tags
	TagListView
	// SearchView allows searching notes
	SearchView
	// ActivityView shows recent activity feed
	ActivityView
	// GraphView displays the knowledge graph
	GraphView
	// HelpView shows keyboard shortcuts and help
	HelpView
)

// String returns the string representation of a View
func (v View) String() string {
	switch v {
	case DashboardView:
		return "Dashboard"
	case NoteListView:
		return "Notes"
	case NoteDetailView:
		return "Note Detail"
	case NoteCreateView:
		return "Create Note"
	case NoteEditView:
		return "Edit Note"
	case TagListView:
		return "Tags"
	case SearchView:
		return "Search"
	case ActivityView:
		return "Activity"
	case GraphView:
		return "Knowledge Graph"
	case HelpView:
		return "Help"
	default:
		return "Unknown"
	}
}

// TUI messages are custom tea.Msg types for TUI-specific events

// switchViewMsg signals to switch to a different view
type switchViewMsg struct {
	View View
}

// viewBackMsg signals to go back to the previous view
type viewBackMsg struct{}

// sessionExpiredMsg signals that the session has expired
type sessionExpiredMsg struct{}

// sessionValidMsg signals that the session is still valid
type sessionValidMsg struct{}

// sessionExpiringSoonMsg warns that the session will expire soon
type sessionExpiringSoonMsg struct {
	TimeRemaining string // Human-readable time remaining
}

// errorMsg signals an error occurred
type errorMsg struct {
	Error error
}

// clearErrorMsg signals to clear the current error
type clearErrorMsg struct{}

// loadingMsg signals that a loading operation should start
type loadingMsg struct{}

// loadedMsg signals that loading is complete
type loadedMsg struct{}

// quitMsg signals to quit the TUI
type quitMsg struct {
	Quit bool
}

// NewSwitchViewCmd creates a command to switch views
func NewSwitchViewCmd(view View) tea.Cmd {
	return func() tea.Msg {
		return switchViewMsg{View: view}
	}
}

// NewViewBackCmd creates a command to go back to the previous view
func NewViewBackCmd() tea.Cmd {
	return func() tea.Msg {
		return viewBackMsg{}
	}
}

// NewErrorMsg creates a command to show an error
func NewErrorMsg(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg{Error: err}
	}
}

// NewClearErrorCmd creates a command to clear the error
func NewClearErrorCmd() tea.Cmd {
	return func() tea.Msg {
		return clearErrorMsg{}
	}
}

// NewLoadingCmd creates a command to start loading
func NewLoadingCmd() tea.Cmd {
	return func() tea.Msg {
		return loadingMsg{}
	}
}

// NewLoadedCmd creates a command to finish loading
func NewLoadedCmd() tea.Cmd {
	return func() tea.Msg {
		return loadedMsg{}
	}
}

// NewQuitCmd creates a command to quit the TUI
func NewQuitCmd() tea.Cmd {
	return func() tea.Msg {
		return quitMsg{Quit: true}
	}
}
