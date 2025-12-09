// internal/tui/keymap/keymap.go

package keymap

// KeyMap holds all key bindings for all modes.
type KeyMap struct {
	// Mode-specific bindings
	Normal   *BindingMap
	Insert   *BindingMap
	Visual   *BindingMap
	Operator *BindingMap // When waiting for motion after d, y, c

	// Motion state for tracking counts and pending operators
	State *MotionState

	// Current mode
	CurrentMode Mode
}

// New creates a new KeyMap with empty bindings.
func New() *KeyMap {
	return &KeyMap{
		Normal:      NewBindingMap(),
		Insert:      NewBindingMap(),
		Visual:      NewBindingMap(),
		Operator:    NewBindingMap(),
		State:       NewMotionState(),
		CurrentMode: ModeNormal,
	}
}

// Default creates a KeyMap with default vim-style bindings.
func Default() *KeyMap {
	km := New()
	km.LoadDefaults()
	return km
}

// LoadDefaults loads the default vim-style key bindings.
func (km *KeyMap) LoadDefaults() {
	km.loadNormalDefaults()
	km.loadInsertDefaults()
	km.loadVisualDefaults()
	km.loadOperatorDefaults()
}

func (km *KeyMap) loadNormalDefaults() {
	n := km.Normal

	// Mode switching
	n.Bind("i", ActionInsertMode)
	n.Bind("a", ActionAppendMode)
	n.Bind("v", ActionVisualMode)

	// Cursor movement
	n.Bind("h", ActionMoveLeft)
	n.Bind("j", ActionMoveDown)
	n.Bind("k", ActionMoveUp)
	n.Bind("l", ActionMoveRight)
	n.Bind("left", ActionMoveLeft)
	n.Bind("down", ActionMoveDown)
	n.Bind("up", ActionMoveUp)
	n.Bind("right", ActionMoveRight)

	// Word movement
	n.Bind("w", ActionMoveWordNext)
	n.Bind("b", ActionMoveWordPrev)

	// Line movement
	n.Bind("0", ActionGotoLineStart)
	n.Bind("^", ActionGotoLineStart)
	n.Bind("$", ActionGotoLineEnd)
	n.Bind("home", ActionGotoLineStart)
	n.Bind("end", ActionGotoLineEnd)

	// Document movement
	n.Bind("gg", ActionGotoTop)
	n.Bind("G", ActionGotoBottom)
	n.Bind("ctrl+u", ActionPageUp)
	n.Bind("ctrl+d", ActionPageDown)
	n.Bind("pgup", ActionPageUp)
	n.Bind("pgdown", ActionPageDown)

	// Editing
	n.Bind("x", ActionDeleteChar)
	n.Bind("X", ActionDeleteCharBack)
	n.Bind("dd", ActionDeleteLine)
	n.Bind("D", ActionDeleteToEnd)
	n.Bind("yy", ActionYankLine)
	n.Bind("Y", ActionYankLine)
	n.Bind("p", ActionPaste)
	n.Bind("P", ActionPasteAbove)
	n.Bind("u", ActionUndo)
	n.Bind("ctrl+r", ActionRedo)
	n.Bind("J", ActionJoinLines)

	// Line operations
	n.Bind("o", ActionOpenBelow)
	n.Bind("O", ActionOpenAbove)

	// Operators (trigger operator-pending mode)
	n.Bind("d", ActionOperatorDelete)
	n.Bind("y", ActionOperatorYank)
	n.Bind("c", ActionOperatorChange)

	// General
	n.Bind("q", ActionQuit)
	n.Bind("ctrl+c", ActionForceQuit)
	n.Bind("ctrl+q", ActionForceQuit)
	n.Bind("ctrl+s", ActionSave)
	n.Bind("ZZ", ActionSaveQuit)
	n.Bind("ZQ", ActionForceQuit)

	// Help & UI
	n.Bind("?", ActionToggleHelp)
	n.Bind("f1", ActionToggleHelp)
	n.Bind("ctrl+l", ActionToggleLineNumbers)
}

func (km *KeyMap) loadInsertDefaults() {
	i := km.Insert

	// Exit insert mode
	i.Bind("esc", ActionNormalMode)
	i.Bind("ctrl+c", ActionNormalMode)
	i.Bind("ctrl+[", ActionNormalMode)

	// Editing
	i.Bind("backspace", ActionBackspace)
	i.Bind("delete", ActionDelete)
	i.Bind("enter", ActionInsertNewline)
	i.Bind("tab", ActionInsertTab)

	// Navigation (still works in insert mode)
	i.Bind("up", ActionMoveUp)
	i.Bind("down", ActionMoveDown)
	i.Bind("left", ActionMoveLeft)
	i.Bind("right", ActionMoveRight)
	i.Bind("home", ActionGotoLineStart)
	i.Bind("end", ActionGotoLineEnd)

	// Save without leaving insert mode
	i.Bind("ctrl+s", ActionSave)
}

func (km *KeyMap) loadVisualDefaults() {
	v := km.Visual

	// Exit visual mode
	v.Bind("esc", ActionNormalMode)
	v.Bind("ctrl+c", ActionNormalMode)
	v.Bind("v", ActionNormalMode)

	// Movement
	v.Bind("h", ActionMoveLeft)
	v.Bind("j", ActionMoveDown)
	v.Bind("k", ActionMoveUp)
	v.Bind("l", ActionMoveRight)
	v.Bind("left", ActionMoveLeft)
	v.Bind("down", ActionMoveDown)
	v.Bind("up", ActionMoveUp)
	v.Bind("right", ActionMoveRight)
	v.Bind("w", ActionMoveWordNext)
	v.Bind("b", ActionMoveWordPrev)
	v.Bind("0", ActionGotoLineStart)
	v.Bind("$", ActionGotoLineEnd)
	v.Bind("gg", ActionGotoTop)
	v.Bind("G", ActionGotoBottom)

	// Operations on selection
	v.Bind("d", ActionOperatorDelete)
	v.Bind("y", ActionOperatorYank)
	v.Bind("c", ActionOperatorChange)
	v.Bind("x", ActionOperatorDelete)
}

func (km *KeyMap) loadOperatorDefaults() {
	o := km.Operator

	// Cancel operator
	o.Bind("esc", ActionNormalMode)
	o.Bind("ctrl+c", ActionNormalMode)

	// Motions that can follow an operator
	o.Bind("h", ActionMoveLeft)
	o.Bind("j", ActionMoveDown)
	o.Bind("k", ActionMoveUp)
	o.Bind("l", ActionMoveRight)
	o.Bind("w", ActionMoveWordNext)
	o.Bind("b", ActionMoveWordPrev)
	o.Bind("0", ActionGotoLineStart)
	o.Bind("$", ActionGotoLineEnd)
	o.Bind("gg", ActionGotoTop)
	o.Bind("G", ActionGotoBottom)

	// Double operator = line operation (dd, yy, cc)
	o.Bind("d", ActionDeleteLine)  // dd
	o.Bind("y", ActionYankLine)    // yy
	o.Bind("c", ActionDeleteLine)  // cc (delete line, enter insert)
}

// GetBindingMap returns the binding map for a mode.
func (km *KeyMap) GetBindingMap(mode Mode) *BindingMap {
	switch mode {
	case ModeNormal:
		return km.Normal
	case ModeInsert:
		return km.Insert
	case ModeVisual:
		return km.Visual
	case ModeOperatorPending:
		return km.Operator
	default:
		return km.Normal
	}
}

// Bind adds a key binding to a specific mode.
func (km *KeyMap) Bind(mode Mode, key string, action Action) {
	km.GetBindingMap(mode).Bind(key, action)
}

// Unbind removes a key binding from a specific mode.
func (km *KeyMap) Unbind(mode Mode, key string) {
	km.GetBindingMap(mode).Unbind(key)
}

// Lookup looks up a key in the current mode.
func (km *KeyMap) Lookup(key string) LookupResult {
	mode := km.CurrentMode
	if km.State.HasPendingOperator() {
		mode = ModeOperatorPending
	}
	return km.GetBindingMap(mode).Lookup(km.State.KeyBuffer + key)
}

// ProcessKey processes a key press and returns the resulting command.
func (km *KeyMap) ProcessKey(key string) (Command, bool) {
	key = NormalizeKey(key)

	// Handle digit for count (but not '0' at start which is line start)
	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		if km.State.AddDigit(rune(key[0])) {
			return Command{}, false // Consumed digit, no command yet
		}
	}

	// Add key to buffer
	km.State.AddKey(key)

	// Lookup the current key sequence
	result := km.lookupCurrentSequence()

	switch result.Status {
	case LookupFound:
		// Complete match - execute command
		cmd := km.buildCommand(result.Action)
		km.State.Reset()
		return cmd, true

	case LookupPending:
		// Could be complete or could continue
		// For now, wait for more input or timeout
		// The caller can decide to execute based on result.Action
		if result.Action != ActionNone {
			// Has a valid action but could continue
			// We'll execute immediately for simplicity
			cmd := km.buildCommand(result.Action)
			km.State.Reset()
			return cmd, true
		}
		return Command{}, false

	case LookupPartialMatch:
		// Waiting for more keys
		return Command{}, false

	case LookupNotFound:
		// No binding - reset state
		km.State.Reset()
		return Command{}, false
	}

	return Command{}, false
}

// lookupCurrentSequence looks up the current key buffer.
func (km *KeyMap) lookupCurrentSequence() LookupResult {
	mode := km.CurrentMode
	if km.State.HasPendingOperator() {
		mode = ModeOperatorPending
	}

	return km.GetBindingMap(mode).Lookup(km.State.KeyBuffer)
}

// buildCommand builds a command from the current state and action.
func (km *KeyMap) buildCommand(action Action) Command {
	count := km.State.GetCount()

	if action.IsOperator() && !km.State.HasPendingOperator() {
		// Starting an operator - enter operator-pending mode
		km.State.SetOperator(action)
		km.State.ClearKeyBuffer()
		return Command{} // No command yet
	}

	if km.State.HasPendingOperator() {
		// Completing an operator with a motion
		return NewOperatorCommand(km.State.PendingOperator, count, action, 1)
	}

	// Simple command
	return NewCommand(action, count)
}

// SetMode sets the current mode.
func (km *KeyMap) SetMode(mode Mode) {
	km.CurrentMode = mode
	km.State.Reset()
}

// GetMode returns the current mode.
func (km *KeyMap) GetMode() Mode {
	if km.State.HasPendingOperator() {
		return ModeOperatorPending
	}
	return km.CurrentMode
}

// Reset resets the keymap state.
func (km *KeyMap) Reset() {
	km.State.Reset()
}

// Clone creates a copy of the keymap.
func (km *KeyMap) Clone() *KeyMap {
	return &KeyMap{
		Normal:      km.Normal.Clone(),
		Insert:      km.Insert.Clone(),
		Visual:      km.Visual.Clone(),
		Operator:    km.Operator.Clone(),
		State:       NewMotionState(),
		CurrentMode: km.CurrentMode,
	}
}