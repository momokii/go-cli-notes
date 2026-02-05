package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableRow represents a single row in the table
type TableRow struct {
	Title       string
	Description string
	Metadata    string // Additional info to display (e.g., date, tags)
	ID          string // Unique identifier for the row
}

// Table wraps the Bubbles list component for displaying tabular data
type Table struct {
	list     list.Model
	width    int
	height   int
	showDesc bool // Whether to show description column
}

// NewTable creates a new table component
func NewTable() Table {
	// Define list styles
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Background color
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#cdd6f4"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Faint(true)
	delegate.SetSpacing(0)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)

	return Table{
		list:     l,
		width:    80,
		height:   20,
		showDesc: true,
	}
}

// SetSize sets the size of the table
func (t *Table) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.list.SetSize(width, height)
}

// SetShowDescription sets whether to show the description
func (t *Table) SetShowDescription(show bool) {
	t.showDesc = show
	t.list.SetDelegate(t.getDelegate())
}

// getDelegate returns the delegate based on current settings
func (t *Table) getDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = t.showDesc
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1e1e2e")). // Background color
		Background(lipgloss.Color("#89b4fa")). // Blue
		Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.Copy().Foreground(lipgloss.Color("#cdd6f4"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Faint(true)
	delegate.SetSpacing(0)
	return delegate
}

// SetItems sets the items in the table
func (t *Table) SetItems(rows []TableRow) {
	t.SetRows(rows)
}

// SetItems converts table rows to list items and sets them
func (t *Table) SetRows(rows []TableRow) {
	items := make([]list.Item, len(rows))
	for i, row := range rows {
		items[i] = tableItem{row: row}
	}
	t.list.SetItems(items)
}

// AppendItem adds an item to the table
func (t *Table) AppendItem(row TableRow) {
	t.list.InsertItem(len(t.list.Items()), tableItem{row: row})
}

// ClearItems removes all items from the table
func (t *Table) ClearItems() {
	t.list.SetItems([]list.Item{})
}

// SelectedItem returns the selected table row
func (t *Table) SelectedItem() *TableRow {
	selected := t.list.SelectedItem()
	if selected == nil {
		return nil
	}
	if item, ok := selected.(tableItem); ok {
		return &item.row
	}
	return nil
}

// SelectedIndex returns the index of the selected row
func (t *Table) SelectedIndex() int {
	return t.list.Index()
}

// SetCursor sets the cursor to the specified index
func (t *Table) SetCursor(index int) {
	t.list.Select(index)
}

// CursorUp moves the selection up
func (t *Table) CursorUp() {
	if t.list.Index() > 0 {
		t.list.Select(t.list.Index() - 1)
	}
}

// CursorDown moves the selection down
func (t *Table) CursorDown() {
	items := t.list.Items()
	if t.list.Index() < len(items)-1 {
		t.list.Select(t.list.Index() + 1)
	}
}

// Top moves the cursor to the first item
func (t *Table) Top() {
	t.list.Select(0)
}

// Bottom moves the cursor to the last item
func (t *Table) Bottom() {
	items := t.list.Items()
	if len(items) > 0 {
		t.list.Select(len(items) - 1)
	}
}

// ItemsCount returns the number of items in the table
func (t *Table) ItemsCount() int {
	return len(t.list.Items())
}

// IsEmpty returns true if the table has no items
func (t *Table) IsEmpty() bool {
	return len(t.list.Items()) == 0
}

// Init initializes the table component
func (t *Table) Init() tea.Cmd {
	return nil
}

// Update handles messages for the table
func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.list, cmd = t.list.Update(msg)
	return t, cmd
}

// View renders the table
func (t *Table) View() string {
	if len(t.list.Items()) == 0 {
		return t.emptyView()
	}
	return t.list.View()
}

// emptyView renders the table when empty
func (t *Table) emptyView() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true).
		Italic(true)

	emptyText := "No items found"
	return style.Render(emptyText)
}

// tableItem implements list.Item for TableRow
type tableItem struct {
	row TableRow
}

func (i tableItem) Title() string {
	return i.row.Title
}

func (i tableItem) Description() string {
	return i.row.Description
}

func (i tableItem) FilterValue() string {
	return i.row.Title + " " + i.row.Description + " " + i.row.Metadata
}
