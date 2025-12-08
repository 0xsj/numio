// internal/tui/keymap.go

package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode represents the editor mode (vim-style).
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeVisual
)

// String returns the mode name.
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	default:
		return "UNKNOWN"
	}
}

// KeyMap defines all key bindings for the application.
type KeyMap struct {
	// Navigation (Normal mode)
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Home       key.Binding
	End        key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	GotoTop    key.Binding
	GotoBottom key.Binding
	WordNext   key.Binding
	WordPrev   key.Binding

	// Mode switching
	InsertMode       key.Binding
	InsertModeAppend key.Binding
	InsertModeStart  key.Binding
	InsertModeEnd    key.Binding
	NormalMode       key.Binding
	VisualMode       key.Binding

	// Editing (Normal mode)
	NewLineBelow key.Binding
	NewLineAbove key.Binding
	DeleteLine   key.Binding
	DeleteChar   key.Binding
	Yank         key.Binding
	Paste        key.Binding
	Undo         key.Binding
	Redo         key.Binding

	// Editing (Insert mode)
	Backspace key.Binding
	Delete    key.Binding
	Enter     key.Binding
	Tab       key.Binding

	// General
	Quit         key.Binding
	ForceQuit    key.Binding
	Save         key.Binding
	Help         key.Binding
	Refresh      key.Binding
	ToggleWrap   key.Binding
	ToggleLines  key.Binding
	ToggleHeader key.Binding
	ToggleDebug  key.Binding
}

// DefaultKeyMap returns the default vim-style key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "Move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "Move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "Move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "Move right"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "0"),
			key.WithHelp("Home/0", "Start of line"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "$"),
			key.WithHelp("End/$", "End of line"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("PgUp/^u", "Page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("PgDn/^d", "Page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("gg", "Go to top"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "Go to bottom"),
		),
		WordNext: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "Next word"),
		),
		WordPrev: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "Previous word"),
		),

		// Mode switching
		InsertMode: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "Insert mode"),
		),
		InsertModeAppend: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Append mode"),
		),
		InsertModeStart: key.NewBinding(
			key.WithKeys("I"),
			key.WithHelp("I", "Insert at start"),
		),
		InsertModeEnd: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "Append at end"),
		),
		NormalMode: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "Normal mode"),
		),
		VisualMode: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "Visual mode"),
		),

		// Editing (Normal mode)
		NewLineBelow: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "New line below"),
		),
		NewLineAbove: key.NewBinding(
			key.WithKeys("O"),
			key.WithHelp("O", "New line above"),
		),
		DeleteLine: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("dd", "Delete line"),
		),
		DeleteChar: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "Delete char"),
		),
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("yy", "Yank line"),
		),
		Paste: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Paste"),
		),
		Undo: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "Undo"),
		),
		Redo: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("^r", "Redo"),
		),

		// Editing (Insert mode)
		Backspace: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("Backspace", "Delete back"),
		),
		Delete: key.NewBinding(
			key.WithKeys("delete"),
			key.WithHelp("Delete", "Delete forward"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "New line"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "Insert tab"),
		),

		// General (work in all modes)
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "Quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("^c", "Force quit"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("^s", "Save file"),
		),
		Help: key.NewBinding(
			key.WithKeys("?", "f1"),
			key.WithHelp("?/F1", "Toggle help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("^r", "Refresh rates"),
		),
		ToggleWrap: key.NewBinding(
			key.WithKeys("W"),
			key.WithHelp("W", "Toggle wrap"),
		),
		ToggleLines: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "Toggle line numbers"),
		),
		ToggleHeader: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "Toggle header"),
		),
		ToggleDebug: key.NewBinding(
			key.WithKeys("f12"),
			key.WithHelp("F12", "Toggle debug"),
		),
	}
}

// NormalModeKeys returns keys active in normal mode.
func (k KeyMap) NormalModeKeys() []key.Binding {
	return []key.Binding{
		k.Up, k.Down, k.Left, k.Right,
		k.Home, k.End, k.PageUp, k.PageDown,
		k.GotoTop, k.GotoBottom, k.WordNext, k.WordPrev,
		k.InsertMode, k.InsertModeAppend, k.InsertModeStart, k.InsertModeEnd,
		k.VisualMode,
		k.NewLineBelow, k.NewLineAbove, k.DeleteLine, k.DeleteChar,
		k.Yank, k.Paste, k.Undo, k.Redo,
		k.Quit, k.Save, k.Help, k.Refresh,
		k.ToggleWrap, k.ToggleLines, k.ToggleHeader, k.ToggleDebug,
	}
}

// InsertModeKeys returns keys active in insert mode.
func (k KeyMap) InsertModeKeys() []key.Binding {
	return []key.Binding{
		k.NormalMode,
		k.Backspace, k.Delete, k.Enter, k.Tab,
		k.Up, k.Down, k.Left, k.Right,
		k.Home, k.End,
		k.Save, k.Help,
	}
}

// VisualModeKeys returns keys active in visual mode.
func (k KeyMap) VisualModeKeys() []key.Binding {
	return []key.Binding{
		k.NormalMode,
		k.Up, k.Down, k.Left, k.Right,
		k.Home, k.End, k.PageUp, k.PageDown,
		k.Yank, k.DeleteChar,
		k.Help,
	}
}

// HelpBindings returns key bindings organized for help display.
type HelpBindings struct {
	Section  string
	Bindings []key.Binding
}

// GetHelpBindings returns all bindings organized by section.
func (k KeyMap) GetHelpBindings() []HelpBindings {
	return []HelpBindings{
		{
			Section: "Navigation",
			Bindings: []key.Binding{
				k.Up, k.Down, k.Left, k.Right,
				k.Home, k.End,
				k.PageUp, k.PageDown,
			},
		},
		{
			Section: "Editing",
			Bindings: []key.Binding{
				k.InsertMode, k.InsertModeAppend,
				k.NewLineBelow, k.NewLineAbove,
				k.DeleteLine, k.DeleteChar,
			},
		},
		{
			Section: "General",
			Bindings: []key.Binding{
				k.ToggleWrap, k.ToggleLines, k.ToggleHeader,
				k.Save, k.Refresh, k.ToggleDebug,
				k.Help, k.Quit,
			},
		},
	}
}

// IsQuit checks if the key should quit the application.
func (k KeyMap) IsQuit(msg tea.KeyMsg, mode Mode) bool {
	// Force quit always works
	if key.Matches(msg, k.ForceQuit) {
		return true
	}
	// Regular quit only in normal mode
	if mode == ModeNormal && key.Matches(msg, k.Quit) {
		return true
	}
	return false
}