package main

import (
	"fyne.io/fyne/v2"
)

// ════════════════════════════════════════════════════════════════
// KEY EVENT
// ════════════════════════════════════════════════════════════════

// KeyEvent represents a processed key event for the editor.
type KeyEvent struct {
	Key       string
	Modifiers []string
}

// HasModifier checks if a modifier is present.
func (k KeyEvent) HasModifier(mod string) bool {
	for _, m := range k.Modifiers {
		if m == mod {
			return true
		}
	}
	return false
}

// IsCtrl returns true if Ctrl is held.
func (k KeyEvent) IsCtrl() bool {
	return k.HasModifier("ctrl")
}

// IsAlt returns true if Alt/Option is held.
func (k KeyEvent) IsAlt() bool {
	return k.HasModifier("alt")
}

// IsShift returns true if Shift is held.
func (k KeyEvent) IsShift() bool {
	return k.HasModifier("shift")
}

// IsCmd returns true if Cmd/Super is held.
func (k KeyEvent) IsCmd() bool {
	return k.HasModifier("cmd")
}

// ════════════════════════════════════════════════════════════════
// KEY TRANSLATION
// ════════════════════════════════════════════════════════════════

// TranslateKeyEvent converts a Fyne key event to our KeyEvent.
func TranslateKeyEvent(ev *fyne.KeyEvent) KeyEvent {
	return KeyEvent{
		Key:       translateKey(ev.Name),
		Modifiers: []string{}, // Fyne doesn't provide modifiers in KeyEvent
	}
}

// TranslateKeyEventWithMods converts a Fyne key event with modifiers.
func TranslateKeyEventWithMods(ev *fyne.KeyEvent, ctrl, alt, shift, super bool) KeyEvent {
	mods := make([]string, 0, 4)
	if ctrl {
		mods = append(mods, "ctrl")
	}
	if alt {
		mods = append(mods, "alt")
	}
	if shift {
		mods = append(mods, "shift")
	}
	if super {
		mods = append(mods, "cmd")
	}

	return KeyEvent{
		Key:       translateKey(ev.Name),
		Modifiers: mods,
	}
}

// translateKey converts Fyne key names to our key format.
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

// TranslateRune converts a typed rune to a KeyEvent.
func TranslateRune(r rune) KeyEvent {
	return KeyEvent{
		Key:       string(r),
		Modifiers: []string{},
	}
}

// TranslateRuneWithMods converts a typed rune with modifiers.
func TranslateRuneWithMods(r rune, ctrl, alt, shift, super bool) KeyEvent {
	mods := make([]string, 0, 4)
	if ctrl {
		mods = append(mods, "ctrl")
	}
	if alt {
		mods = append(mods, "alt")
	}
	if shift {
		mods = append(mods, "shift")
	}
	if super {
		mods = append(mods, "cmd")
	}

	return KeyEvent{
		Key:       string(r),
		Modifiers: mods,
	}
}
