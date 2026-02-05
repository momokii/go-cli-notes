package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/momokii/go-cli-notes/cmd/cli/tui/components"
	"github.com/momokii/go-cli-notes/cmd/cli/tui/models"
)

// MainModel is the root model for the TUI application
type MainModel struct {
	// Shared state
	client    *client.APIClient
	authState *client.AuthState
	userInfo  string

	// View management
	currentView View
	prevView    View
	quitting    bool

	// Child models
	helpModel       models.HelpModel
	dashboardModel  models.DashboardModel
	noteListModel   models.NoteListModel
	noteDetailModel models.NoteDetailModel
	noteCreateModel models.NoteCreateModel
	tagListModel    models.TagListModel
	searchModel     models.SearchModel
	activityModel   models.ActivityModel
	graphModel      models.GraphModel

	// Track initialization of child models
	dashboardInitialized  bool
	noteListInitialized   bool
	noteDetailInitialized bool
	noteCreateInitialized bool
	tagListInitialized    bool
	searchInitialized     bool
	activityInitialized   bool
	graphInitialized      bool

	// Shared components
	statusBar *components.StatusBar

	// Session management
	sessionValid         bool
	lastSessionCheck     time.Time
	sessionCheckInterval time.Duration
	sessionWarningShown  bool
	sessionExpiryWarning time.Duration // Warning threshold (e.g., 5 minutes)

	// Error handling
	currentError    error
	clearErrorAfter time.Duration

	// Dimensions
	width  int
	height int
}

// NewMainModel creates a new main TUI model
func NewMainModel(apiClient *client.APIClient, authState *client.AuthState) MainModel {
	// Get user info from auth state
	userInfo := ""
	if authState != nil && authState.Email != "" {
		userInfo = authState.Email
	}

	sb := components.NewStatusBar()
	sb.SetUserInfo(userInfo)

	return MainModel{
		client:                apiClient,
		authState:             authState,
		userInfo:              userInfo,
		currentView:           DashboardView, // Start with dashboard
		prevView:              DashboardView,
		quitting:              false,
		helpModel:             models.NewHelpModel(),
		dashboardModel:        models.NewDashboardModel(apiClient, authState),
		noteListModel:         models.NewNoteListModel(apiClient, authState),
		noteDetailModel:       models.NewNoteDetailModel(apiClient, authState),
		noteCreateModel:       models.NewNoteCreateModel(apiClient, authState),
		tagListModel:          models.NewTagListModel(apiClient, authState),
		searchModel:           models.NewSearchModel(apiClient, authState),
		activityModel:         models.NewActivityModel(apiClient, authState),
		graphModel:            models.NewGraphModel(apiClient, authState),
		dashboardInitialized:  false,
		noteListInitialized:   false,
		noteDetailInitialized: false,
		noteCreateInitialized: false,
		tagListInitialized:    false,
		searchInitialized:     false,
		activityInitialized:   false,
		graphInitialized:      false,
		statusBar:             sb,
		sessionValid:          true,
		lastSessionCheck:      time.Now(),
		sessionCheckInterval:  5 * time.Minute,
		sessionWarningShown:   false,
		sessionExpiryWarning:  5 * time.Minute, // Show warning 5 minutes before expiry
		currentError:          nil,
		clearErrorAfter:       5 * time.Second,
		width:                 80,
		height:                24,
	}
}

// Init initializes the main model
func (m MainModel) Init() tea.Cmd {
	// Start with initial session validation and periodic checks
	return tea.Batch(
		m.checkSessionCmd(),
		m.dashboardModel.Init(),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return clearErrorMsg{}
		}),
	)
}

// Update handles messages for the main model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			m.statusBar.ShowError(fmt.Sprintf("Recovered from panic: %v", r))
			m.currentError = fmt.Errorf("panic: %v", r)
		}
	}()

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle exit keys FIRST - these should always work regardless of focus
		switch msg.String() {
		case "q", "ctrl+c":
			if !m.quitting {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			// Go back to previous view
			if m.currentView == HelpView {
				m.currentView = m.prevView
			} else if m.currentView != DashboardView {
				m.cleanupView(m.currentView)
				m.prevView = m.currentView
				m.currentView = DashboardView
				// Clear any status notifications when returning to dashboard
				m.statusBar.ClearError()
				if !m.dashboardInitialized {
					m.dashboardInitialized = true
					initCmd := m.dashboardModel.Init()
					m.updateStatusBar()
					return m, initCmd
				}
			}
			m.updateStatusBar()
			return m, nil
		}

		// Check if current view has focused input component
		// If so, skip navigation-only handlers to allow typing
		if m.isInputFocused() {
			// Forward directly to child model, don't process navigation keys
			break
		}

		// Handle other global keys (navigation, view switching, etc.)
		switch msg.String() {
		case "?":
			// Toggle help
			if m.currentView != HelpView {
				m.prevView = m.currentView
				m.currentView = HelpView
				m.helpModel.GotoTop()
			} else {
				m.currentView = m.prevView
			}
			m.updateStatusBar()
			return m, nil

		case "/":
			// Quick search
			m.cleanupView(m.currentView)
			m.prevView = m.currentView
			m.currentView = SearchView
			// FIX: Ensure input is focused when entering search view
			m.searchModel = m.searchModel.FocusInput()
			if !m.searchInitialized {
				m.searchInitialized = true
				initCmd := m.searchModel.Init()
				m.updateStatusBar()
				return m, initCmd
			}
			m.updateStatusBar()
			return m, nil

		case "s":
			// Quick search (same as "/" - documented in help)
			m.cleanupView(m.currentView)
			m.prevView = m.currentView
			m.currentView = SearchView
			// FIX: Ensure input is focused when entering search view
			m.searchModel = m.searchModel.FocusInput()
			if !m.searchInitialized {
				m.searchInitialized = true
				initCmd := m.searchModel.Init()
				m.updateStatusBar()
				return m, initCmd
			}
			m.updateStatusBar()
			return m, nil

		case "t":
			// Tags view
			m.cleanupView(m.currentView)
			m.prevView = m.currentView
			m.currentView = TagListView
			if !m.tagListInitialized {
				m.tagListInitialized = true
				initCmd := m.tagListModel.Init()
				m.updateStatusBar()
				return m, initCmd
			}
			m.updateStatusBar()
			return m, nil

		case "a":
			// Activity view - skip if we're in Note Detail Tags tab (where 'a' is for adding tags)
			if m.currentView != NoteDetailView || m.noteDetailModel.GetCurrentTab() != models.NoteTagsTab {
				m.cleanupView(m.currentView)
				m.prevView = m.currentView
				m.currentView = ActivityView
				if !m.activityInitialized {
					m.activityInitialized = true
					initCmd := m.activityModel.Init()
					m.updateStatusBar()
					return m, initCmd
				}
				m.updateStatusBar()
				return m, nil
			}
			// If we're in Note Detail Tags tab, fall through - let the child model handle it

		case "g":
			// Graph view
			m.cleanupView(m.currentView)
			m.prevView = m.currentView
			m.currentView = GraphView
			if !m.graphInitialized {
				m.graphInitialized = true
				initCmd := m.graphModel.Init()
				m.updateStatusBar()
				return m, initCmd
			}
			m.updateStatusBar()
			return m, nil

		case "n":
			// Quick new note
			m.cleanupView(m.currentView)
			m.prevView = m.currentView
			m.currentView = NoteCreateView
			// Create a fresh model to clear previous input, then focus it
			m.noteCreateModel = models.NewNoteCreateModel(m.client, m.authState)
			m.noteCreateModel = m.noteCreateModel.FocusForm() // Focus the form
			m.noteCreateInitialized = false
			if !m.noteCreateInitialized {
				m.noteCreateInitialized = true
				initCmd := m.noteCreateModel.Init()
				m.updateStatusBar()
				return m, initCmd
			}
			m.updateStatusBar()
			return m, nil
		}

	// Handle quit message
	case quitMsg:
		if msg.Quit {
			m.quitting = true
			return m, tea.Quit
		}

	// Handle view switching
	case switchViewMsg:
		// Cleanup the previous view before switching
		m.cleanupView(m.currentView)

		m.prevView = m.currentView
		m.currentView = msg.View
		m.updateStatusBar()
		return m, nil

	// Handle view requests from child models
	case models.ShowHelpMsg:
		m.prevView = m.currentView
		m.currentView = HelpView
		m.helpModel.GotoTop()
		m.updateStatusBar()
		return m, nil

	case models.ShowNoteListMsg:
		m.cleanupView(m.currentView)
		m.prevView = m.currentView
		m.currentView = NoteListView
		if !m.noteListInitialized {
			m.noteListInitialized = true
			m.updateStatusBar()
			return m, m.noteListModel.Init()
		}
		m.updateStatusBar()
		return m, nil

	case models.OpenNoteMsg:
		// Navigate to note detail view
		m.cleanupView(m.currentView)
		m.prevView = m.currentView
		m.currentView = NoteDetailView
		// FIX: Capture both the modified model and the fetch command
		// Previously, SetNoteID() only returned the command but modified a copy of the model,
		// causing noteID to remain nil and all API calls to use UUID 00000000-0000-0000-0000-000000000000
		var cmd tea.Cmd
		m.noteDetailModel, cmd = m.noteDetailModel.SetNoteID(msg.NoteID)
		if !m.noteDetailInitialized {
			m.noteDetailInitialized = true
		}
		m.updateStatusBar()
		return m, cmd

	case models.EditNoteMsg:
		// Navigate to edit view
		m.cleanupView(m.currentView)
		m.prevView = m.currentView
		m.currentView = NoteEditView

		// Fetch the note if we don't have it or if it's a different note
		var cmds []tea.Cmd

		note := m.noteDetailModel.GetNote()
		if note == nil || note.ID != msg.NoteID {
			// Need to fetch the note first
			var cmd tea.Cmd
			m.noteDetailModel, cmd = m.noteDetailModel.SetNoteID(msg.NoteID)
			cmds = append(cmds, cmd)
			// After note is fetched, retry EditNoteMsg to set edit mode
			cmds = append(cmds, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return models.EditNoteMsg{NoteID: msg.NoteID}
			}))
		} else {
			// Note already loaded, set edit mode and capture the returned model
			m.noteCreateModel, _ = m.noteCreateModel.SetEditMode(note)
		}

		if !m.noteCreateInitialized {
			m.noteCreateInitialized = true
		}
		m.updateStatusBar()
		return m, tea.Batch(cmds...)

	case models.ShowDashboardMsg:
		m.prevView = m.currentView
		m.currentView = DashboardView
		// Clear any notifications when showing dashboard
		m.statusBar.ClearError() // This clears both error and info messages

		// Initialize dashboard if needed to prevent stuck loading state
		var cmd tea.Cmd
		if !m.dashboardInitialized {
			m.dashboardInitialized = true
			cmd = m.dashboardModel.Init()
		}

		m.updateStatusBar()
		return m, cmd

	// Handle note created message
	case models.NoteCreatedMsg:
		m.statusBar.ShowInfo("Note created successfully")
		// Clear the loading state and form before transitioning
		m.noteCreateModel = m.noteCreateModel.BlurForm()
		// Go to dashboard after a brief delay so user sees the success message
		return m, tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
			return models.ShowDashboardMsg{}
		})

	// Handle note updated message
	case models.NoteUpdatedMsg:
		m.statusBar.ShowInfo("Note updated successfully")
		// Clear the form and go to dashboard
		m.noteCreateModel = m.noteCreateModel.BlurForm()
		return m, tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
			return models.ShowDashboardMsg{}
		})

	// Handle note deleted message
	case models.NoteDeletedMsg:
		m.statusBar.ShowInfo("Note deleted")
		// Use ShowDashboardMsg flow to properly clear notifications and handle initialization
		// Also auto-clear the "Note deleted" notification after brief delay
		return m, tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
			return models.ShowDashboardMsg{}
		})

	// Handle note create error
	case models.NoteCreateErrMsg:
		m.statusBar.ShowError(msg.Err.Error())
		return m, nil

	// Handle tag operations
	case models.TagCreatedMsg:
		m.statusBar.ShowInfo("Tag created")
		// Clear notification after a brief delay
		return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return clearErrorMsg{}
		})

	case models.TagUpdatedMsg:
		m.statusBar.ShowInfo("Tag updated")
		// Clear notification after a brief delay
		return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return clearErrorMsg{}
		})

	case models.TagDeletedMsg:
		m.statusBar.ShowInfo("Tag deleted")
		// Clear notification after a brief delay
		return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
			return clearErrorMsg{}
		})

	case models.TagListErrMsg:
		m.statusBar.ShowError(msg.Err.Error())
		return m, nil

	case models.SearchErrMsg:
		m.statusBar.ShowError(msg.Err.Error())
		return m, nil

	// Handle activity operations
	case models.ActivityErrMsg:
		m.statusBar.ShowError(msg.Err.Error())
		return m, nil

	case models.GraphErrMsg:
		m.statusBar.ShowError(msg.Err.Error())
		return m, nil

	// Handle filter notes by tag message
	case models.FilterNotesByTagMsg:
		m.prevView = m.currentView
		m.currentView = NoteListView
		// Set the tag filter on the note list model
		m.noteListModel = m.noteListModel.SetTagFilter(msg.TagID.String())
		// Reset initialization state so it fetches with new filter
		m.noteListInitialized = true
		initCmd := m.noteListModel.Init()
		m.statusBar.ShowInfo(fmt.Sprintf("Filtered by tag: %s", msg.TagName))
		m.updateStatusBar()
		return m, initCmd

	// Handle going back
	case viewBackMsg:
		if m.currentView == HelpView {
			m.currentView = m.prevView
		} else if m.currentView != DashboardView {
			m.currentView = DashboardView
		}
		m.updateStatusBar()
		return m, nil

	// Handle session expiration
	case sessionExpiredMsg:
		m.sessionValid = false
		m.quitting = true
		// Display session expired message before quitting
		return m, tea.Sequence(
			tea.Printf("\n"+FatalStyle.Render("Session Expired")+"\n"),
			tea.Quit,
		)

	// Handle session validation
	case sessionValidMsg:
		m.sessionValid = true
		m.lastSessionCheck = time.Now()
		return m, nil

	// Handle session expiring soon warning
	case sessionExpiringSoonMsg:
		if !m.sessionWarningShown {
			m.sessionWarningShown = true
			m.statusBar.ShowError(fmt.Sprintf("Session expiring in %s - please save your work", msg.TimeRemaining))
		}
		return m, nil

	// Handle error messages
	case errorMsg:
		m.currentError = msg.Error
		m.statusBar.ShowError(msg.Error.Error())
		return m, tea.Tick(m.clearErrorAfter, func(t time.Time) tea.Msg {
			return clearErrorMsg{}
		})

	// Handle clearing errors
	case clearErrorMsg:
		m.currentError = nil
		m.statusBar.ClearError()
		return m, nil

	// Handle window resize
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusBar.SetWidth(msg.Width)
		m.helpModel.SetSize(msg.Width, msg.Height)
		// Also update child models' internal sizes
		var model tea.Model
		var _ tea.Cmd
		model, _ = m.dashboardModel.Update(msg)
		m.dashboardModel = model.(models.DashboardModel)
		model, _ = m.noteListModel.Update(msg)
		m.noteListModel = model.(models.NoteListModel)
		model, _ = m.noteDetailModel.Update(msg)
		m.noteDetailModel = model.(models.NoteDetailModel)
		model, _ = m.noteCreateModel.Update(msg)
		m.noteCreateModel = model.(models.NoteCreateModel)
		model, _ = m.tagListModel.Update(msg)
		m.tagListModel = model.(models.TagListModel)
		model, _ = m.searchModel.Update(msg)
		m.searchModel = model.(models.SearchModel)
		model, _ = m.activityModel.Update(msg)
		m.activityModel = model.(models.ActivityModel)
		model, _ = m.graphModel.Update(msg)
		m.graphModel = model.(models.GraphModel)
		return m, nil

	// Handle tea.Quit (from child models)
	case tea.QuitMsg:
		m.quitting = true
		return m, tea.Quit
	}

	// Update child models based on current view
	var cmd tea.Cmd
	var model tea.Model

	switch m.currentView {
	case HelpView:
		model, cmd = m.helpModel.Update(msg)
		m.helpModel = model.(models.HelpModel)

		// Check if child model wants to quit
		// If cmd is not nil and is the quit command, handle it
		if cmd != nil && m.quitting {
			m.quitting = true
			return m, cmd
		}
		// Also check if the child model returned a quit command
		// by seeing if it responds to keys with quit
		if keyMsg, ok := msg.(tea.KeyMsg); ok && (keyMsg.String() == "q" || keyMsg.String() == "esc" || keyMsg.String() == "ctrl+c") {
			m.quitting = true
			return m, tea.Quit
		}

	case DashboardView:
		// Let dashboard handle its own messages
		model, cmd = m.dashboardModel.Update(msg)
		m.dashboardModel = model.(models.DashboardModel)

	case NoteListView:
		// Let note list handle its own messages
		model, cmd = m.noteListModel.Update(msg)
		m.noteListModel = model.(models.NoteListModel)

	case NoteDetailView:
		// Let note detail handle its own messages
		model, cmd = m.noteDetailModel.Update(msg)
		m.noteDetailModel = model.(models.NoteDetailModel)

	case NoteCreateView, NoteEditView:
		// Let note create/edit handle its own messages
		model, cmd = m.noteCreateModel.Update(msg)
		m.noteCreateModel = model.(models.NoteCreateModel)

	case TagListView:
		// Let tag list handle its own messages
		model, cmd = m.tagListModel.Update(msg)
		m.tagListModel = model.(models.TagListModel)

	case SearchView:
		// Let search handle its own messages
		model, cmd = m.searchModel.Update(msg)
		m.searchModel = model.(models.SearchModel)

	case ActivityView:
		// Let activity handle its own messages
		model, cmd = m.activityModel.Update(msg)
		m.activityModel = model.(models.ActivityModel)

	case GraphView:
		// Let graph handle its own messages
		model, cmd = m.graphModel.Update(msg)
		m.graphModel = model.(models.GraphModel)

	default:
		// Unknown view, do nothing
	}

	cmds = append(cmds, cmd)

	// Schedule periodic session check
	if time.Since(m.lastSessionCheck) > m.sessionCheckInterval {
		cmds = append(cmds, m.checkSessionCmd())
	}

	return m, tea.Batch(cmds...)
}

// View renders the main TUI
func (m MainModel) View() string {
	// Panic recovery to prevent rendering crashes
	defer func() {
		if r := recover(); r != nil {
			// Can't show in status bar here, but log the panic
			// The next Update cycle will handle it
		}
	}()

	if m.quitting {
		return ""
	}

	// Render header
	header := RenderHeader("Knowledge Garden - TUI", m.userInfo, m.width)

	// Render current view content
	var content string
	switch m.currentView {
	case HelpView:
		content = m.helpModel.View()
	case DashboardView:
		content = m.dashboardModel.View()
	case NoteListView:
		content = m.noteListModel.View()
	case NoteDetailView:
		content = m.noteDetailModel.View()
	case NoteCreateView, NoteEditView:
		content = m.noteCreateModel.View()
	case TagListView:
		content = m.tagListModel.View()
	case SearchView:
		content = m.searchModel.View()
	case ActivityView:
		content = m.activityModel.View()
	case GraphView:
		content = m.graphModel.View()
	default:
		// Unknown view
		content := "\n  Unknown view\n  Press ? for help\n"
		content = DimStyle.Render(content)
	}

	// Calculate content height
	headerLines := 2
	statusBarLines := 1
	availableHeight := m.height - headerLines - statusBarLines

	// Ensure content fits
	contentLines := countLines(content)
	for contentLines < availableHeight {
		content += "\n"
		contentLines++
	}

	// Render status bar
	statusBar := m.statusBar.View()

	// Combine everything
	return header + "\n" + content + "\n" + statusBar
}

// updateStatusBar updates the status bar based on current view
func (m *MainModel) updateStatusBar() {
	m.statusBar.SetViewInfo(m.currentView.String(), GetViewKeyHelp(m.currentView))
	m.statusBar.SetUserInfo(m.userInfo)
	m.statusBar.SetWidth(m.width)
}

// checkSessionCmd returns a command that checks if the session is still valid
func (m MainModel) checkSessionCmd() tea.Cmd {
	return func() tea.Msg {
		// First check if token is expiring soon (using JWT expiry time)
		if m.authState.IsExpiringSoon(m.sessionExpiryWarning) {
			until := m.authState.TimeUntilExpiry()
			timeRemaining := formatDuration(until)
			return sessionExpiringSoonMsg{TimeRemaining: timeRemaining}
		}

		// Then check if token is already expired
		if m.authState.IsExpired() {
			return sessionExpiredMsg{}
		}

		// Finally validate with API call
		if ValidateSession(m.client, m.authState) {
			return sessionValidMsg{}
		}
		return sessionExpiredMsg{}
	}
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than 1 minute"
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}
	hours := int(d.Hours())
	if hours == 1 {
		return "about 1 hour"
	}
	return fmt.Sprintf("about %d hours", hours)
}

// cleanupView cleans up resources when leaving a view
// This helps free memory for views with large data structures
func (m *MainModel) cleanupView(view View) {
	switch view {
	case DashboardView:
		// Clear dashboard to force refresh on next visit
		m.dashboardModel = models.NewDashboardModel(m.client, m.authState)
		m.dashboardInitialized = false
	case NoteListView:
		// Clear note list to force refresh on next visit
		m.noteListModel = models.NewNoteListModel(m.client, m.authState)
		m.noteListInitialized = false
	case NoteCreateView:
		// Clear create note form and remove focus
		m.noteCreateModel = m.noteCreateModel.BlurForm()
		m.noteCreateInitialized = false
	case NoteEditView:
		// Clear edit note form and remove focus
		m.noteCreateModel = m.noteCreateModel.BlurForm()
		m.noteCreateInitialized = false
	case GraphView:
		// Clear graph data (can be large with many nodes/edges)
		m.graphModel = models.NewGraphModel(m.client, m.authState)
		m.graphInitialized = false
	case ActivityView:
		// Clear activity feed (can accumulate over time)
		m.activityModel = models.NewActivityModel(m.client, m.authState)
		m.activityInitialized = false
	case SearchView:
		// Clear search and remove focus
		m.searchModel = m.searchModel.BlurInput()
		m.searchInitialized = false
	case TagListView:
		// Clear tag list and remove focus from forms
		m.tagListModel = m.tagListModel.BlurInput()
		m.tagListInitialized = false
		// Note: We don't clear NoteDetailView
		// as it is commonly used and clearing it would disrupt user workflow
	}
}

// countLines counts the number of lines in a string
func countLines(s string) int {
	count := 0
	for _, c := range s {
		if c == '\n' {
			count++
		}
	}
	if len(s) > 0 {
		count++
	}
	return count
}

// GetViewForTesting returns the current view (for testing purposes)
func (m MainModel) GetViewForTesting() View {
	return m.currentView
}

// IsSessionValidForTesting returns whether the session is valid (for testing)
func (m MainModel) IsSessionValidForTesting() bool {
	return m.sessionValid
}

// isInputFocused checks if the current view has a focused input component
// FIX: This prevents global navigation keys from consuming typing input
func (m MainModel) isInputFocused() bool {
	switch m.currentView {
	case NoteCreateView, NoteEditView:
		return m.noteCreateModel.IsInputFocused()
	case SearchView:
		return m.searchModel.IsInputFocused()
	case TagListView:
		return m.tagListModel.IsInputFocused()
	case NoteDetailView:
		return m.noteDetailModel.IsInputFocused()
	default:
		return false
	}
}
