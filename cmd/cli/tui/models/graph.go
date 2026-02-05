package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/momokii/go-cli-notes/internal/model"
)

// GraphModel is the model for the knowledge graph view
type GraphModel struct {
	client     *client.APIClient
	authState  *client.AuthState
	graph      *model.GraphResponse
	loading    bool
	err        error
	selected   int
	expanded   map[uuid.UUID]bool // Track which nodes are expanded
	width      int
	height     int
	maxNodes   int
}

// NewGraphModel creates a new graph model
func NewGraphModel(apiClient *client.APIClient, authState *client.AuthState) GraphModel {
	return GraphModel{
		client:   apiClient,
		authState: authState,
		loading:  true,
		expanded: make(map[uuid.UUID]bool),
		width:    80,
		height:   24,
		maxNodes: 20, // Initial view shows 20 nodes
	}
}

// Init initializes the graph model
func (m GraphModel) Init() tea.Cmd {
	return m.fetchGraphCmd()
}

// fetchGraphCmd returns a command that fetches the graph
func (m GraphModel) fetchGraphCmd() tea.Cmd {
	return func() tea.Msg {
		graph, err := m.client.GetGraph()
		if err != nil {
			return GraphErrMsg{Err: err}
		}
		return GraphFetchedMsg{Graph: graph}
	}
}

// Update handles messages for the graph model
func (m GraphModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.graph != nil && m.selected < len(m.graph.Nodes)-1 {
				m.selected++
			}
		case "k", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "+", "=":
			// Show more nodes
			if m.maxNodes < 100 {
				m.maxNodes += 10
			}
		case "-", "_":
			// Show fewer nodes
			if m.maxNodes > 10 {
				m.maxNodes -= 10
			}
		case "enter":
			// Open selected note
			if m.graph != nil && len(m.graph.Nodes) > 0 && m.selected >= 0 {
				noteID := m.graph.Nodes[m.selected].ID
				return m, func() tea.Msg {
					return OpenNoteMsg{NoteID: noteID}
				}
			}
		case " ":
			// Toggle expand/collapse selected node
			if m.graph != nil && len(m.graph.Nodes) > 0 && m.selected >= 0 {
				nodeID := m.graph.Nodes[m.selected].ID
				if m.expanded[nodeID] {
					delete(m.expanded, nodeID)
				} else {
					m.expanded[nodeID] = true
				}
			}
		}

	case GraphFetchedMsg:
		m.graph = msg.Graph
		m.loading = false
		if len(m.graph.Nodes) > 0 {
			// Auto-expand first few nodes
			for i := 0; i < min(3, len(m.graph.Nodes)); i++ {
				m.expanded[m.graph.Nodes[i].ID] = true
			}
		}
		return m, nil

	case GraphErrMsg:
		m.err = msg.Err
		m.loading = false
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the graph view
func (m GraphModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.err != nil {
		return m.renderError()
	}

	return m.renderContent()
}

// renderLoading renders the loading state
func (m GraphModel) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	return style.Render("Loading knowledge graph...")
}

// renderError renders the error state
func (m GraphModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")). // Red
		Bold(true)

	return errorStyle.Render("Error: " + m.err.Error())
}

// renderContent renders the graph content
func (m GraphModel) renderContent() string {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")). // Orange
		Bold(true).
		MarginBottom(1)

	nodeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")) // Light text

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Dark
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)

	linkStyle := lipgloss.NewStyle().
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
	content += titleStyle.Render("KNOWLEDGE GRAPH") + "\n\n"

	if m.graph == nil || len(m.graph.Nodes) == 0 {
		content += mutedStyle.Render("(no notes in knowledge graph)")
		content += "\n\n"
		content += hintStyle.Render("+:more nodes -:fewer nodes ESC:back ?:help")
		return content
	}

	// Stats
	if m.graph.Stats != nil {
		statsText := fmt.Sprintf("Nodes: %d | Links: %d",
			m.graph.Stats.TotalNotes,
			m.graph.Stats.TotalLinks)
		content += mutedStyle.Render(statsText) + "\n\n"
	}

	// Build adjacency list for connections
	connections := m.buildConnections()

	// Display nodes (with connections)
	displayCount := min(m.maxNodes, len(m.graph.Nodes))
	for i := 0; i < displayCount; i++ {
		node := m.graph.Nodes[i]
		isSelected := i == m.selected
		isExpanded := m.expanded[node.ID]

		// Node indicator
		var indicator string
		if isSelected {
			indicator = "→ "
		} else {
			indicator = "  "
		}

		// Expand/collapse indicator
		var expandChar string
		if hasConnections(connections, node.ID) {
			if isExpanded {
				expandChar = "[-]"
			} else {
				expandChar = "[+]"
			}
		} else {
			expandChar = "   "
		}

		// Node title
		title := node.Title
		if title == "" {
			title = "(untitled)"
		}
		if len(title) > 50 {
			title = title[:47] + "..."
		}

		nodeLine := fmt.Sprintf("%s%s %s", indicator, expandChar, title)

		if isSelected {
			content += selectedStyle.Render(nodeLine)
		} else {
			content += nodeStyle.Render(nodeLine)
		}

		content += "\n"

		// Show connections if expanded
		if isExpanded {
			for _, edge := range connections[node.ID] {
				targetTitle := m.getNodeTitle(edge.Target)
				content += linkStyle.Render(fmt.Sprintf("    └─→ %s", targetTitle)) + "\n"
			}
		}
	}

	if len(m.graph.Nodes) > displayCount {
		content += fmt.Sprintf("\n... and %d more (press + to show more)", len(m.graph.Nodes)-displayCount)
	}

	// Hints
	content += "\n" + hintStyle.Render("j/k:navigate Enter:open +/-:zoom Space:expand ESC:back ?:help")

	return content
}

// buildConnections builds an adjacency list of connections for each node
func (m GraphModel) buildConnections() map[uuid.UUID][]*model.GraphEdge {
	connections := make(map[uuid.UUID][]*model.GraphEdge)

	if m.graph == nil {
		return connections
	}

	for _, edge := range m.graph.Edges {
		connections[edge.Source] = append(connections[edge.Source], edge)
	}

	return connections
}

// getNodeTitle returns the title of a node by ID
func (m GraphModel) getNodeTitle(id uuid.UUID) string {
	if m.graph == nil {
		return "(unknown)"
	}

	for _, node := range m.graph.Nodes {
		if node.ID == id {
			title := node.Title
			if title == "" {
				title = "(untitled)"
			}
			if len(title) > 40 {
				title = title[:37] + "..."
			}
			return title
		}
	}

	return "(unknown)"
}

// hasConnections checks if a node has any outgoing connections
func hasConnections(connections map[uuid.UUID][]*model.GraphEdge, nodeID uuid.UUID) bool {
	return len(connections[nodeID]) > 0
}

// Message types for graph

type GraphFetchedMsg struct {
	Graph *model.GraphResponse
}

type GraphErrMsg struct {
	Err error
}
