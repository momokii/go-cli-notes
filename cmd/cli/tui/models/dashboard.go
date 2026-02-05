package models

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/momokii/go-cli-notes/internal/model"
)

// DashboardModel is the model for the dashboard view
type DashboardModel struct {
	client      *client.APIClient
	authState   *client.AuthState
	stats       *model.UserStats
	activity    []*model.Activity
	trending    []*model.TrendingNote
	loading     bool
	err         error
	width       int
	height      int
	lastUpdate  time.Time
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(apiClient *client.APIClient, authState *client.AuthState) DashboardModel {
	return DashboardModel{
		client:    apiClient,
		authState: authState,
		loading:   true,
		width:     80,
		height:    24,
	}
}

// Init initializes the dashboard model
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchStatsCmd(),
		m.fetchActivityCmd(),
		m.fetchTrendingCmd(),
	)
}

// fetchStatsCmd returns a command that fetches user stats
func (m DashboardModel) fetchStatsCmd() tea.Cmd {
	return func() tea.Msg {
		stats, err := m.client.GetStats()
		if err != nil {
			return dashboardErrMsg{err}
		}
		return dashboardStatsMsg{stats}
	}
}

// fetchActivityCmd returns a command that fetches recent activity
func (m DashboardModel) fetchActivityCmd() tea.Cmd {
	return func() tea.Msg {
		activity, err := m.client.GetRecentActivity(5)
		if err != nil {
			return dashboardErrMsg{err}
		}
		return dashboardActivityMsg{activity}
	}
}

// fetchTrendingCmd returns a command that fetches trending notes
func (m DashboardModel) fetchTrendingCmd() tea.Cmd {
	return func() tea.Msg {
		trending, err := m.client.GetTrendingNotes(5)
		if err != nil {
			return dashboardErrMsg{err}
		}
		return dashboardTrendingMsg{trending}
	}
}

// Update handles messages for the dashboard model
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			// Request to switch to help view - main model will handle this
			return m, func() tea.Msg {
				return ShowHelpMsg{}
			}
		case "esc":
			// Already on dashboard, ignore
			return m, nil
		case "n":
			// New note - Phase C
			return m, nil
		case "s":
			// Search - Phase D
			return m, nil
		case "l":
			// List notes
			return m, func() tea.Msg {
				return ShowNoteListMsg{}
			}
		case "a":
			// Activity feed - Phase D
			return m, nil
		}

	case dashboardStatsMsg:
		m.stats = msg.stats
		m.loading = false
		return m, nil

	case dashboardActivityMsg:
		m.activity = msg.activity
		return m, nil

	case dashboardTrendingMsg:
		m.trending = msg.trending
		return m, nil

	case dashboardErrMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the dashboard
func (m DashboardModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#45475a")). // Dark gray
		Padding(0, 1).
		MarginBottom(1)

	// Build content
	var content string

	// Statistics section
	content += titleStyle.Render("STATISTICS")
	content += "\n"
	statsBox := m.renderStats(labelStyle, valueStyle)
	content += boxStyle.Width(m.width).Render(statsBox)
	content += "\n\n"

	// Recent Activity section - always show (even when empty)
	content += titleStyle.Render("RECENT ACTIVITY")
	content += "\n"
	activityBox := m.renderActivity(labelStyle, valueStyle, mutedStyle)
	content += boxStyle.Width(m.width).Render(activityBox)
	content += "\n\n"

	// Trending Notes section
	if len(m.trending) > 0 {
		content += titleStyle.Render("TRENDING NOTES")
		content += "\n"
		trendingBox := m.renderTrending(labelStyle, valueStyle, mutedStyle)
		content += boxStyle.Width(m.width).Render(trendingBox)
		content += "\n\n"
	}

	// Quick Actions
	content += m.renderQuickActions()

	return content
}

// renderLoading renders the loading state
func (m DashboardModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Loading dashboard...")
}

// renderError renders the error state
func (m DashboardModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render(fmt.Sprintf("Error loading dashboard: %v", m.err))
}

// renderStats renders the statistics section
func (m DashboardModel) renderStats(labelStyle, valueStyle lipgloss.Style) string {
	var stats string

	if m.stats == nil {
		return "No statistics available"
	}

	// Format stats
	stats += labelStyle.Render("Total Notes:")
	stats += valueStyle.Render(fmt.Sprintf("%d\n", m.stats.TotalNotes))

	stats += labelStyle.Render("Total Tags:")
	stats += valueStyle.Render(fmt.Sprintf("%d\n", m.stats.TotalTags))

	stats += labelStyle.Render("Total Links:")
	stats += valueStyle.Render(fmt.Sprintf("%d\n", m.stats.TotalLinks))

	stats += labelStyle.Render("Total Words:")
	stats += valueStyle.Render(fmt.Sprintf("%d\n", m.stats.TotalWords))

	stats += labelStyle.Render("Created Today:")
	stats += valueStyle.Render(fmt.Sprintf("%d\n", m.stats.NotesCreatedToday))

	stats += labelStyle.Render("Created This Week:")
	stats += valueStyle.Render(fmt.Sprintf("%d\n", m.stats.NotesCreatedWeek))

	if m.stats.LastActivity != nil {
		stats += labelStyle.Render("Last Activity:")
		stats += valueStyle.Render(formatTimeAgo(*m.stats.LastActivity) + "\n")
	}

	return stats
}

// renderActivity renders the recent activity section
func (m DashboardModel) renderActivity(labelStyle, valueStyle, mutedStyle lipgloss.Style) string {
	if len(m.activity) == 0 {
		return mutedStyle.Render("No recent activity")
	}

	var activity string
	for _, act := range m.activity {
		timestamp := formatTimeAgo(act.CreatedAt)
		activity += labelStyle.Render(timestamp)
		activity += " "
		activity += valueStyle.Render(formatActivity(act))
		activity += "\n"
	}

	return activity
}

// renderTrending renders the trending notes section
func (m DashboardModel) renderTrending(labelStyle, valueStyle, mutedStyle lipgloss.Style) string {
	if len(m.trending) == 0 {
		return mutedStyle.Render("No trending notes")
	}

	var trending string
	for i, note := range m.trending {
		if note.Note == nil {
			continue
		}

		trending += labelStyle.Render(fmt.Sprintf("%d.", i+1))
		trending += " "
		trending += valueStyle.Render(truncateText(note.Note.Title, 40))
		trending += "\n"

		// Show access count
		trending += mutedStyle.Render(fmt.Sprintf("   Accessed %d times recently\n", note.RecentAccess))
	}

	return trending
}

// renderQuickActions renders the quick actions section
func (m DashboardModel) renderQuickActions() string {
	quickActionsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")). // Green
		Bold(true)

	actions := "Quick Actions: "
	actions += quickActionsStyle.Render("[N]ew Note")
	actions += " "
	actions += quickActionsStyle.Render("[S]earch")
	actions += " "
	actions += quickActionsStyle.Render("[L]ist Notes")
	actions += " "
	actions += quickActionsStyle.Render("[T]ags")
	actions += " "
	actions += quickActionsStyle.Render("[A]ctivity Feed")

	return actions
}

// formatActivity formats an activity for display
func formatActivity(act *model.Activity) string {
	switch act.Action {
	case "create":
		if act.NoteID != nil {
			return fmt.Sprintf("Created note: %s", getNoteTitleFromMetadata(act.Metadata))
		}
		return "Created a note"
	case "update":
		if act.NoteID != nil {
			return fmt.Sprintf("Updated note: %s", getNoteTitleFromMetadata(act.Metadata))
		}
		return "Updated a note"
	case "delete":
		if act.NoteID != nil {
			return fmt.Sprintf("Deleted note: %s", getNoteTitleFromMetadata(act.Metadata))
		}
		return "Deleted a note"
	case "link":
		return "Created a link between notes"
	default:
		return string(act.Action)
	}
}

// getNoteTitleFromMetadata extracts the note title from activity metadata
func getNoteTitleFromMetadata(metadata model.ActivityMetadata) string {
	if title, ok := metadata["note_title"].(string); ok {
		return title
	}
	if title, ok := metadata["title"].(string); ok {
		return title
	}
	return "Unknown Note"
}

// formatTimeAgo formats a time as "X time ago"
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "Just now"
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	return t.Format("Jan 2, 2006")
}

// truncateText truncates text to a maximum length
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// Message types for dashboard

type dashboardStatsMsg struct {
	stats *model.UserStats
}

type dashboardActivityMsg struct {
	activity []*model.Activity
}

type dashboardTrendingMsg struct {
	trending []*model.TrendingNote
}

type dashboardErrMsg struct {
	err error
}

// View request messages - these are sent to the main model to request view changes
type ShowHelpMsg struct{}
type ShowNoteListMsg struct{}

// Ensure DashboardModel implements list.Item for potential list usage
var _ list.Item = (*dashboardActivityItem)(nil)

type dashboardActivityItem struct {
	activity *model.Activity
}

func (i dashboardActivityItem) Title() string {
	return formatTimeAgo(i.activity.CreatedAt)
}

func (i dashboardActivityItem) Description() string {
	return formatActivity(i.activity)
}

func (i dashboardActivityItem) FilterValue() string {
	return formatActivity(i.activity) + " " + formatTimeAgo(i.activity.CreatedAt)
}
