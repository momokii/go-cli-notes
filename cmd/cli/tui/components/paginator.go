package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Paginator wraps the Bubbles paginator component
type Paginator struct {
	paginator paginator.Model
	width     int
}

// NewPaginator creates a new paginator component
func NewPaginator() Paginator {
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 20
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa")).Render("•") // Blue
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render("•") // Gray

	return Paginator{
		paginator: p,
		width:     80,
	}
}

// SetTotalItems sets the total number of items and calculates total pages
func (p *Paginator) SetTotalItems(total int) {
	p.paginator.TotalPages = (total + p.paginator.PerPage - 1) / p.paginator.PerPage
	if p.paginator.TotalPages == 0 {
		p.paginator.TotalPages = 1
	}
}

// SetPerPage sets the number of items per page
func (p *Paginator) SetPerPage(perPage int) {
	p.paginator.PerPage = perPage
}

// SetPage sets the current page (1-indexed)
func (p *Paginator) SetPage(page int) {
	if page < 1 {
		page = 1
	}
	if page > p.paginator.TotalPages {
		page = p.paginator.TotalPages
	}
	p.paginator.Page = page - 1 // Bubbles uses 0-indexed pages
}

// CurrentPage returns the current page number (1-indexed)
func (p *Paginator) CurrentPage() int {
	return p.paginator.Page + 1
}

// TotalPages returns the total number of pages
func (p *Paginator) TotalPages() int {
	return p.paginator.TotalPages
}

// OnPage returns true if the paginator is on the specified page (1-indexed)
func (p *Paginator) OnPage(page int) bool {
	return p.paginator.Page == page-1
}

// NextPage moves to the next page
func (p *Paginator) NextPage() {
	p.paginator.Page++
}

// PrevPage moves to the previous page
func (p *Paginator) PrevPage() {
	if p.paginator.Page > 0 {
		p.paginator.Page--
	}
}

// FirstPage moves to the first page
func (p *Paginator) FirstPage() {
	p.paginator.Page = 0
}

// LastPage moves to the last page
func (p *Paginator) LastPage() {
	if p.paginator.TotalPages > 0 {
		p.paginator.Page = p.paginator.TotalPages - 1
	}
}

// Update handles messages for the paginator
func (p *Paginator) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	p.paginator, cmd = p.paginator.Update(msg)
	return cmd
}

// SetWidth sets the width for rendering
func (p *Paginator) SetWidth(width int) {
	p.width = width
}

// View renders the paginator
func (p *Paginator) View() string {
	if p.paginator.TotalPages <= 1 {
		return ""
	}

	// Build custom pagination display
	current := p.CurrentPage()
	total := p.TotalPages()

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")). // Gray
		Faint(true)

	info := style.Render(fmt.Sprintf("Page %d of %d", current, total))

	// Add navigation hints
	navHints := []string{}
	if current > 1 {
		navHints = append(navHints, "Ctrl+P:prev")
	}
	if current < total {
		navHints = append(navHints, "Ctrl+N:next")
	}

	var navText string
	if len(navHints) > 0 {
		navStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89b4fa")) // Blue
		navText = " | " + navStyle.Render(navHints[0])
		if len(navHints) > 1 {
			navText += " " + navStyle.Render(navHints[1])
		}
	}

	return info + navText
}

// ItemsOnPage returns the start and end indices for the current page
func (p *Paginator) ItemsOnPage(totalItems int) (start, end int) {
	currentPage := p.CurrentPage()
	perPage := p.paginator.PerPage

	start = (currentPage - 1) * perPage
	end = start + perPage

	if end > totalItems {
		end = totalItems
	}

	return start, end
}

// CanGoNext returns true if there's a next page
func (p *Paginator) CanGoNext() bool {
	return p.paginator.Page < p.paginator.TotalPages-1
}

// CanGoPrev returns true if there's a previous page
func (p *Paginator) CanGoPrev() bool {
	return p.paginator.Page > 0
}
