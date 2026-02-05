package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/momokii/go-cli-notes/cmd/cli/client"
)

// Run starts the TUI application
// Returns an error if initialization fails or if the program exits with an error
func Run(apiClient *client.APIClient, authState *client.AuthState) error {
	// Validate session before starting TUI
	if err := InitTUI(apiClient, authState); err != nil {
		return fmt.Errorf("session validation failed: %w", err)
	}

	// Create the main model
	mainModel := NewMainModel(apiClient, authState)

	// Create the Bubbletea program
	p := tea.NewProgram(
		mainModel,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Start the program
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	// Check if the session was valid on exit
	if m, ok := finalModel.(MainModel); ok {
		if !m.IsSessionValidForTesting() {
			fmt.Println("\n" + FatalStyle.Render("Session Expired"))
			fmt.Println(MutedStyle.Render("Please run 'kg-cli login' to refresh your session"))
			return fmt.Errorf("session expired")
		}
	}

	return nil
}

// CheckTerminalSize checks if the terminal is large enough for the TUI
// Returns true if the terminal is sufficient size, false otherwise
func CheckTerminalSize() (bool, int, int) {
	width, height, err := getTerminalSize()
	if err != nil {
		// If we can't determine size, assume it's fine
		return true, 80, 24
	}

	// Minimum dimensions: 80x24
	minWidth := 80
	minHeight := 24

	if width < minWidth || height < minHeight {
		return false, width, height
	}

	return true, width, height
}

// getTerminalSize returns the current terminal size
func getTerminalSize() (int, int, error) {
	// Use the tea implementation to get terminal size
	// This is a simplified version - the actual size will be detected by Bubbletea
	return 80, 24, nil
}

// ShowTerminalSizeWarning displays a warning if the terminal is too small
func ShowTerminalSizeWarning(width, height int) {
	warning := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")). // Yellow
		Bold(true).
		Render("⚠ Terminal Size Warning")

	content := fmt.Sprintf(
		`Your terminal is too small for the TUI interface.

Current size: %dx%d
Minimum size: 80x24

Please resize your terminal or use a larger window.

Press any key to continue anyway (some features may not work properly).`,
		width, height,
	)

	fmt.Println(warning)
	fmt.Println(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render(content))
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// ExitWithMessage exits the TUI with a message
func ExitWithMessage(message string, code int) {
	ClearScreen()
	if message != "" {
		fmt.Println(message)
	}
	os.Exit(code)
}
