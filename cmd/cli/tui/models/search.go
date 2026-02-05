package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/momokii/go-cli-notes/cmd/cli/tui/components"
	"github.com/momokii/go-cli-notes/internal/model"
)

// SearchModel is the model for the search view
type SearchModel struct {
	client        *client.APIClient
	authState     *client.AuthState
	query         string
	results       []*model.SearchResult
	loading       bool
	err           error
	selectedIndex int
	input         components.TextInput
	paginator     components.Paginator
	width         int
	height        int
}

// NewSearchModel creates a new search model
func NewSearchModel(apiClient *client.APIClient, authState *client.AuthState) SearchModel {
	input := components.NewTextInput()
	input.SetPlaceholder("Search notes...")
	input.SetWidth(40)

	paginator := components.NewPaginator()
	paginator.SetPerPage(10)

	return SearchModel{
		client:    apiClient,
		authState: authState,
		input:     input,
		paginator: paginator,
		width:     80,
		height:    24,
	}
}

// Init initializes the search model
func (m SearchModel) Init() tea.Cmd {
	m.input.Focus()
	return nil
}

// FocusInput focuses the search input and returns the updated model
// This is needed because models are stored as values in the main model
func (m SearchModel) FocusInput() SearchModel {
	m.input.Focus()
	return m
}

// SetQuery sets the search query (for external navigation)
func (m SearchModel) SetQuery(query string) tea.Cmd {
	m.query = query
	m.input.SetValue(query)
	return m.searchCmd(query)
}

// Update handles messages for the search model
func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input mode
		if m.input.Focused() {
			switch msg.String() {
			case "enter":
				query := m.input.Value()
				if query != "" {
					m.query = query
					m.input.Blur()
					return m, m.searchCmd(query)
				}
			case "esc":
				m.input.Blur()
				if m.query == "" {
					// No search yet, go back
					return m, func() tea.Msg {
						return ShowDashboardMsg{}
					}
				}
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			default:
				cmd := m.input.Update(msg)
				return m, cmd
			}
		}

		// Handle results mode
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			return m, func() tea.Msg {
				return ShowHelpMsg{}
			}
		case "esc":
			// Go back to input or dashboard
			if m.query != "" {
				m.query = ""
				m.results = nil
				m.input.SetValue("")
				m.input.Focus()
				return m, nil
			}
			return m, func() tea.Msg {
				return ShowDashboardMsg{}
			}
		case "/":
			// Focus search input
			m.input.Focus()
			return m, nil
		case "j", "down":
			if len(m.results) > 0 && m.selectedIndex < len(m.results)-1 {
				m.selectedIndex++
			}
		case "k", "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "enter":
			// Open selected note
			if len(m.results) > 0 && m.selectedIndex >= 0 {
				noteID := m.results[m.selectedIndex].Note.ID
				return m, func() tea.Msg {
					return OpenNoteMsg{NoteID: noteID}
				}
			}
		case "ctrl+n", "right":
			// Next page
			if m.paginator.CanGoNext() {
				m.paginator.NextPage()
				return m, m.searchPageCmd(m.query, m.paginator.CurrentPage())
			}
		case "ctrl+p", "left":
			// Previous page
			if m.paginator.CanGoPrev() {
				m.paginator.PrevPage()
				return m, m.searchPageCmd(m.query, m.paginator.CurrentPage())
			}
		}

	case SearchResultsMsg:
		m.results = msg.Results
		m.loading = false
		m.selectedIndex = 0

		// Update paginator
		totalItems := len(msg.Results)
		if msg.Pagination != nil {
			totalItems = int(msg.Pagination.Total)
		}
		m.paginator.SetTotalItems(totalItems)
		m.paginator.SetPage(msg.CurrentPage)

		return m, nil

	case SearchErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(msg.Width - 20)
		return m, nil
	}

	// Update paginator (in case it receives paginator messages)
	_ = m.paginator.Update(msg)

	return m, tea.Batch(cmds...)
}

// searchCmd returns a command that performs a search
func (m SearchModel) searchCmd(query string) tea.Cmd {
	m.loading = true
	m.query = query
	return func() tea.Msg {
		resp, err := m.client.SearchNotes(query, 1, 10)
		if err != nil {
			return SearchErrMsg{Err: err}
		}
		return SearchResultsMsg{
			Query:       resp.Query,
			Results:     resp.Results,
			Pagination:  resp.Pagination,
			CurrentPage: 1,
		}
	}
}

// searchPageCmd returns a command that performs a search for a specific page
func (m SearchModel) searchPageCmd(query string, page int) tea.Cmd {
	m.loading = true
	return func() tea.Msg {
		resp, err := m.client.SearchNotes(query, page, 10)
		if err != nil {
			return SearchErrMsg{Err: err}
		}
		return SearchResultsMsg{
			Query:       resp.Query,
			Results:     resp.Results,
			Pagination:  resp.Pagination,
			CurrentPage: page,
		}
	}
}

// View renders the search view
func (m SearchModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	return m.renderContent()
}

// renderLoading renders the loading state
func (m SearchModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Searching...")
}

// renderError renders the error state
func (m SearchModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render("Error: " + m.err.Error())
}

// renderContent renders the search content
func (m SearchModel) renderContent() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	resultStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	snippetStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
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
	content += titleStyle.Render("SEARCH") + "\n\n"

	// Search input
	content += labelStyle.Render("Query: ")
	content += m.input.View()
	content += "\n\n"

	// Results
	if m.query == "" {
		content += mutedStyle.Render("Enter a search query and press Enter")
		content += "\n\n"
		content += hintStyle.Render("/:focus ESC:back ?:help")
		return content
	}

	if len(m.results) == 0 {
		content += mutedStyle.Render(fmt.Sprintf("No results for \"%s\"", m.query))
		content += "\n\n"
		content += hintStyle.Render("j/k:navigate /:new search ESC:back ?:help")
		return content
	}

	// Display all results (API already handles pagination)
	// Don't apply local pagination since results are already paginated by API
	for i := range m.results {
		result := m.results[i]
		var line string

		if i == m.selectedIndex {
			line = "→ "
		} else {
			line = "  "
		}

		// Title with highlighting
		title := result.Note.Title
		if title == "" {
			title = "(untitled)"
		}
		title = m.highlightText(title, m.query)

		if i == m.selectedIndex {
			line += selectedStyle.Render(title)
		} else {
			line += resultStyle.Render(title)
		}

		content += line + "\n"

		// Snippet
		if result.Snippet != "" {
			snippet := m.highlightText(result.Snippet, m.query)
			content += "    " + snippetStyle.Render(snippet) + "\n"
		}
	}

	// Paginator
	content += "\n" + m.paginator.View()

	// Hints
	content += "\n" + hintStyle.Render("Enter:open j/k:navigate Ctrl+N/P:page /:new search ESC:back")

	return content
}

// highlightText highlights search terms in text
func (m SearchModel) highlightText(text, query string) string {
	if query == "" {
		return text
	}

	// Split query into words
	words := strings.Fields(query)
	if len(words) == 0 {
		return text
	}

	// Use the first word for highlighting (simple approach)
	term := strings.ToLower(words[0])
	lowerText := strings.ToLower(text)

	var result strings.Builder
	lastIndex := 0

	for {
		idx := strings.Index(lowerText[lastIndex:], term)
		if idx == -1 {
			result.WriteString(text[lastIndex:])
			break
		}

		idx += lastIndex

		// Add text before match
		result.WriteString(text[lastIndex:idx])

		// Add highlighted match
		matchEnd := idx + len(term)
		if matchEnd > len(text) {
			matchEnd = len(text)
		}
		result.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f9e2af")). // Yellow
			Bold(true).
			Render(text[idx:matchEnd]))

		lastIndex = matchEnd
	}

	return result.String()
}

// IsInputFocused returns whether the search input is focused
// FIX: Check actual component focus state, not just view type
func (m SearchModel) IsInputFocused() bool {
	return m.input.Focused()
}

// BlurInput removes focus from the search input
func (m SearchModel) BlurInput() SearchModel {
	m.input.Blur()
	return m
}

// Message types for search

type SearchResultsMsg struct {
	Query       string
	Results     []*model.SearchResult
	Pagination  *model.Pagination
	CurrentPage int
}

type SearchErrMsg struct {
	Err error
}
