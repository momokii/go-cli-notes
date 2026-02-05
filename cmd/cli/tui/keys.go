package tui

import (
	"github.com/charmbracelet/bubbletea"
)

// KeyBinding defines a keyboard shortcut with its description
type KeyBinding struct {
	Keys   string // Comma-separated key names
	Action string // What action it performs
	Help   string // Help text to display
}

// GlobalKeyBindings are key bindings that work across all views
var GlobalKeyBindings = []KeyBinding{
	{Keys: "q", Action: "quit", Help: "q:quit"},
	{Keys: "?", Action: "help", Help: "?:help"},
	{Keys: "esc", Action: "back", Help: "esc:back"},
	{Keys: "/", Action: "search", Help: "/:search"},
	{Keys: "n", Action: "new", Help: "n:new"},
	{Keys: "ctrl+c", Action: "force_quit", Help: "ctrl+c:force quit"},
}

// NavigationKeyBindings are keys for navigating lists and menus
var NavigationKeyBindings = []KeyBinding{
	{Keys: "h,←", Action: "left", Help: "h/←:left"},
	{Keys: "j,↓", Action: "down", Help: "j/↓:down"},
	{Keys: "k,↑", Action: "up", Help: "k/↑:up"},
	{Keys: "l,→", Action: "right", Help: "l/→:right"},
	{Keys: "enter", Action: "select", Help: "enter:select"},
	{Keys: "0", Action: "top", Help: "0:top"},
	{Keys: "G", Action: "bottom", Help: "G:bottom"},
}

// NoteListKeyBindings are keys specific to the note list view
var NoteListKeyBindings = []KeyBinding{
	{Keys: "h,←", Action: "left", Help: "h/←:left"},
	{Keys: "j,↓", Action: "down", Help: "j/↓:down"},
	{Keys: "k,↑", Action: "up", Help: "k/↑:up"},
	{Keys: "l,→", Action: "right", Help: "l/→:right"},
	{Keys: "enter", Action: "select", Help: "enter:select"},
	{Keys: "/", Action: "filter", Help: "/:filter"},
	{Keys: "f", Action: "filter_menu", Help: "f:filter"},
	{Keys: "s", Action: "sort", Help: "s:sort"},
	{Keys: "t", Action: "tag_filter", Help: "t:tag"},
}

// NoteDetailKeyBindings are keys for viewing a note
var NoteDetailKeyBindings = []KeyBinding{
	{Keys: "h,←", Action: "left", Help: "h/←:left"},
	{Keys: "j,↓", Action: "down", Help: "j/↓:down"},
	{Keys: "k,↑", Action: "up", Help: "k/↑:up"},
	{Keys: "l,→", Action: "right", Help: "l/→:right"},
	{Keys: "enter", Action: "select", Help: "enter:select"},
	{Keys: "e", Action: "edit", Help: "e:edit"},
	{Keys: "l", Action: "links", Help: "l:links"},
	{Keys: "b", Action: "backlinks", Help: "b:backlinks"},
	{Keys: "t", Action: "tags", Help: "t:tags"},
	{Keys: "tab", Action: "next_tab", Help: "tab:next"},
	{Keys: "shift+tab", Action: "prev_tab", Help: "shift+tab:prev"},
}

// NoteEditKeyBindings are keys for editing a note
var NoteEditKeyBindings = []KeyBinding{
	{Keys: "ctrl+s", Action: "save", Help: "ctrl+s:save"},
	{Keys: "esc", Action: "cancel", Help: "esc:cancel"},
	{Keys: "tab", Action: "next_field", Help: "tab:next"},
	{Keys: "shift+tab", Action: "prev_field", Help: "shift+tab:prev"},
}

// TagListKeyBindings are keys for the tag list view
var TagListKeyBindings = []KeyBinding{
	{Keys: "h,←", Action: "left", Help: "h/←:left"},
	{Keys: "j,↓", Action: "down", Help: "j/↓:down"},
	{Keys: "k,↑", Action: "up", Help: "k/↑:up"},
	{Keys: "l,→", Action: "right", Help: "l/→:right"},
	{Keys: "enter", Action: "select", Help: "enter:select"},
	{Keys: "c", Action: "create", Help: "c:create"},
	{Keys: "e", Action: "edit", Help: "e:edit"},
	{Keys: "d", Action: "delete", Help: "d:delete"},
}

// PaginationKeyBindings are keys for pagination
var PaginationKeyBindings = []KeyBinding{
	{Keys: "ctrl+n", Action: "next_page", Help: "ctrl+n:next"},
	{Keys: "]", Action: "next_page", Help: "]:next"},
	{Keys: "ctrl+p", Action: "prev_page", Help: "ctrl+p:prev"},
	{Keys: "[", Action: "prev_page", Help: "[:prev"},
}

// GetKeyHelp returns help text for a given set of key bindings
func GetKeyHelp(bindings []KeyBinding) string {
	var help string
	for _, kb := range bindings {
		if kb.Help != "" {
			if help != "" {
				help += " "
			}
			help += kb.Help
		}
	}
	return help
}

// GetViewKeyHelp returns the appropriate help text for a given view
func GetViewKeyHelp(view View) string {
	switch view {
	case DashboardView:
		return "n:new s:search l:list a:activity ?:help q:quit"
	case NoteListView:
		return GetKeyHelp(NoteListKeyBindings) + " q:back ?:help"
	case NoteDetailView:
		return GetKeyHelp(NoteDetailKeyBindings) + " q:back ?:help"
	case NoteCreateView:
		return "tab:next enter:create esc:cancel ?:help"
	case TagListView:
		return GetKeyHelp(TagListKeyBindings) + " q:back ?:help"
	case SearchView:
		return "enter:search ↑↓:nav q:back ?:help"
	case ActivityView:
		return "↑↓:scroll enter:open f:filter q:back ?:help"
	case GraphView:
		return "↑↓←→:nav enter:view d:details q:back ?:help"
	case HelpView:
		return "↑↓:scroll q:close esc:close"
	default:
		return "q:quit ?:help"
	}
}

// IsQuitKey checks if a key message is a quit command
func IsQuitKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "q", "ctrl+c":
		return true
	default:
		return false
	}
}

// IsHelpKey checks if a key message is a help command
func IsHelpKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "?", "f1":
		return true
	default:
		return false
	}
}

// IsBackKey checks if a key message is a back/escape command
func IsBackKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "esc":
		return true
	default:
		return false
	}
}

// IsNavUp checks if a key message is an up navigation
func IsNavUp(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "k", "up":
		return true
	default:
		return false
	}
}

// IsNavDown checks if a key message is a down navigation
func IsNavDown(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "j", "down":
		return true
	default:
		return false
	}
}

// IsNavLeft checks if a key message is a left navigation
func IsNavLeft(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "h", "left":
		return true
	default:
		return false
	}
}

// IsNavRight checks if a key message is a right navigation
func IsNavRight(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "l", "right":
		return true
	default:
		return false
	}
}

// IsSelectKey checks if a key message is a select command
func IsSelectKey(msg tea.KeyMsg) bool {
	return msg.String() == "enter"
}
