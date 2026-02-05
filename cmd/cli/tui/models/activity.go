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

// ActivityModel is the model for the activity feed view
type ActivityModel struct {
	client        *client.APIClient
	authState     *client.AuthState
	activities    []*model.Activity
	loading       bool
	err           error
	selectedIndex int
	paginator     components.Paginator
	width         int
	height        int
}

// NewActivityModel creates a new activity model
func NewActivityModel(apiClient *client.APIClient, authState *client.AuthState) ActivityModel {
	paginator := components.NewPaginator()
	paginator.SetPerPage(20)

	return ActivityModel{
		client:    apiClient,
		authState: authState,
		paginator: paginator,
		width:     80,
		height:    24,
	}
}

// Init initializes the activity model
func (m ActivityModel) Init() tea.Cmd {
	return m.fetchActivityCmd()
}

// fetchActivityCmd returns a command that fetches recent activity
func (m ActivityModel) fetchActivityCmd() tea.Cmd {
	return func() tea.Msg {
		activities, err := m.client.GetRecentActivity(50) // Get last 50 activities
		if err != nil {
			return ActivityErrMsg{Err: err}
		}
		return ActivityFetchedMsg{Activities: activities}
	}
}

// Update handles messages for the activity model
func (m ActivityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			return m, func() tea.Msg {
				return ShowHelpMsg{}
			}
		case "esc":
			return m, func() tea.Msg {
				return ShowDashboardMsg{}
			}
		case "j", "down":
			_, endIdx := m.paginator.ItemsOnPage(len(m.activities))
			if m.selectedIndex < endIdx-1 {
				m.selectedIndex++
			}
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "ctrl+n", "right":
			// Next page
			if m.paginator.CanGoNext() {
				m.paginator.NextPage()
				startIdx, _ := m.paginator.ItemsOnPage(len(m.activities))
				m.selectedIndex = startIdx
			}
		case "ctrl+p", "left":
			// Previous page
			if m.paginator.CanGoPrev() {
				m.paginator.PrevPage()
				startIdx, _ := m.paginator.ItemsOnPage(len(m.activities))
				m.selectedIndex = startIdx
			}
		case "enter":
			// Open the note associated with this activity
			if len(m.activities) > 0 && m.selectedIndex >= 0 {
				activity := m.activities[m.selectedIndex]
				if activity.NoteID != nil && *activity.NoteID != uuid.Nil {
					return m, func() tea.Msg {
						return OpenNoteMsg{NoteID: *activity.NoteID}
					}
				}
			}
		}

	case ActivityFetchedMsg:
		m.activities = msg.Activities
		m.loading = false
		m.selectedIndex = 0
		m.paginator.SetTotalItems(len(msg.Activities))
		return m, nil

	case ActivityErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	// Update paginator
	_ = m.paginator.Update(msg)

	return m, nil
}

// View renders the activity feed view
func (m ActivityModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	return m.renderContent()
}

// renderLoading renders the loading state
func (m ActivityModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Loading activity...")
}

// renderError renders the error state
func (m ActivityModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render("Error: " + m.err.Error())
}

// renderContent renders the activity content
func (m ActivityModel) renderContent() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	activityStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	timestampStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	actionCreateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")). // Green
	Bold(true)

	actionUpdateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
	Bold(true)

	actionDeleteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
	Bold(true)

	actionViewStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cba6f7")). // Purple
	Faint(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true).
		MarginTop(1)

	var content string

	// Title
	content += titleStyle.Render("RECENT ACTIVITY") + "\n\n"

	if len(m.activities) == 0 {
		content += mutedStyle.Render("(no recent activity)")
		content += "\n\n"
		content += hintStyle.Render("ESC:back ?:help")
		return content
	}

	// Get current page range
	startIdx, endIdx := m.paginator.ItemsOnPage(len(m.activities))

	for i := startIdx; i < endIdx; i++ {
		activity := m.activities[i]
		var line string

		// Selection indicator
		if i == m.selectedIndex {
			line = "→ "
		} else {
			line = "  "
		}

		// Timestamp
		timestamp := formatTimeAgo(activity.CreatedAt)
		line += timestampStyle.Render("["+timestamp+"]") + " "

		// Action
		actionStyle := getActivityActionStyle(activity.Action, actionCreateStyle, actionUpdateStyle, actionDeleteStyle, actionViewStyle)
		actionText := actionStyle.Render(formatActivityAction(activity.Action))

		// Note info
		var noteInfo string
		if activity.NoteID != nil && *activity.NoteID != uuid.Nil {
			noteInfo = " note"
		}

		// Build full line
		fullLine := fmt.Sprintf("%s%s%s", line, actionText, noteInfo)

		if i == m.selectedIndex {
			content += selectedStyle.Render(fullLine)
		} else {
			content += activityStyle.Render(fullLine)
		}

		content += "\n"
	}

	// Paginator
	perPage := 20
	if len(m.activities) > perPage {
		content += "\n" + m.paginator.View()
	}

	// Hints
	content += "\n" + hintStyle.Render("j/k:navigate Enter:open Ctrl+N/P:page ESC:back ?:help")

	return content
}

// getActivityActionStyle returns the appropriate style for an action type
func getActivityActionStyle(action model.ActionType, create, update, delete, view lipgloss.Style) lipgloss.Style {
	switch action {
	case model.ActionCreate:
		return create
	case model.ActionUpdate:
		return update
	case model.ActionDelete:
		return delete
	case model.ActionView, model.ActionSearch:
		return view
	default:
		return view
	}
}

// formatActivityAction formats an activity action for display
func formatActivityAction(action model.ActionType) string {
	switch action {
	case model.ActionCreate:
		return "Created"
	case model.ActionUpdate:
		return "Updated"
	case model.ActionDelete:
		return "Deleted"
	case model.ActionView:
		return "Viewed"
	case model.ActionSearch:
		return "Searched"
	case model.ActionLogin:
		return "Logged in"
	case model.ActionLogout:
		return "Logged out"
	default:
		return string(action)
	}
}

// Message types for activity feed

type ActivityFetchedMsg struct {
	Activities []*model.Activity
}

type ActivityErrMsg struct {
	Err error
}
