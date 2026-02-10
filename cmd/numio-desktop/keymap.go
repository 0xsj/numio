package main

import (
	"fyne.io/fyne/v2"
)

// ════════════════════════════════════════════════════════════════
// KEY TRANSLATION (for simple/non-vim mode)
// ════════════════════════════════════════════════════════════════

// translateKey converts Fyne key names to our internal key format.
func translateKey(name fyne.KeyName) string {
	switch name {
	// Special keys
	case fyne.KeyReturn, fyne.KeyEnter:
		return "Enter"
	case fyne.KeyTab:
		return "Tab"
	case fyne.KeyBackspace:
		return "Backspace"
	case fyne.KeyDelete:
		return "Delete"
	case fyne.KeyEscape:
		return "Escape"
	case fyne.KeySpace:
		return " "

	// Arrow keys
	case fyne.KeyUp:
		return "ArrowUp"
	case fyne.KeyDown:
		return "ArrowDown"
	case fyne.KeyLeft:
		return "ArrowLeft"
	case fyne.KeyRight:
		return "ArrowRight"

	// Navigation
	case fyne.KeyHome:
		return "Home"
	case fyne.KeyEnd:
		return "End"
	case fyne.KeyPageUp:
		return "PageUp"
	case fyne.KeyPageDown:
		return "PageDown"

	// Function keys
	case fyne.KeyF1:
		return "F1"
	case fyne.KeyF2:
		return "F2"
	case fyne.KeyF3:
		return "F3"
	case fyne.KeyF4:
		return "F4"
	case fyne.KeyF5:
		return "F5"
	case fyne.KeyF6:
		return "F6"
	case fyne.KeyF7:
		return "F7"
	case fyne.KeyF8:
		return "F8"
	case fyne.KeyF9:
		return "F9"
	case fyne.KeyF10:
		return "F10"
	case fyne.KeyF11:
		return "F11"
	case fyne.KeyF12:
		return "F12"

	default:
		// Regular character keys
		return string(name)
	}
}

// ════════════════════════════════════════════════════════════════
// KEYMAP-COMPATIBLE KEY TRANSLATION
// ════════════════════════════════════════════════════════════════

// translateKeyForKeymap converts Fyne key names to the keymap package's
// lowercase format (e.g. "esc", "enter", "up", "ctrl+r").
func translateKeyForKeymap(name fyne.KeyName) string {
	switch name {
	case fyne.KeyReturn, fyne.KeyEnter:
		return "enter"
	case fyne.KeyTab:
		return "tab"
	case fyne.KeyBackspace:
		return "backspace"
	case fyne.KeyDelete:
		return "delete"
	case fyne.KeyEscape:
		return "esc"
	case fyne.KeySpace:
		return "space"

	case fyne.KeyUp:
		return "up"
	case fyne.KeyDown:
		return "down"
	case fyne.KeyLeft:
		return "left"
	case fyne.KeyRight:
		return "right"

	case fyne.KeyHome:
		return "home"
	case fyne.KeyEnd:
		return "end"
	case fyne.KeyPageUp:
		return "pgup"
	case fyne.KeyPageDown:
		return "pgdown"

	case fyne.KeyF1:
		return "f1"
	case fyne.KeyF2:
		return "f2"
	case fyne.KeyF3:
		return "f3"
	case fyne.KeyF4:
		return "f4"
	case fyne.KeyF5:
		return "f5"
	case fyne.KeyF6:
		return "f6"
	case fyne.KeyF7:
		return "f7"
	case fyne.KeyF8:
		return "f8"
	case fyne.KeyF9:
		return "f9"
	case fyne.KeyF10:
		return "f10"
	case fyne.KeyF11:
		return "f11"
	case fyne.KeyF12:
		return "f12"

	default:
		return ""
	}
}
