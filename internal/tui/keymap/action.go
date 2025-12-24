// internal/tui/keymap/action.go

package keymap

// Action represents an editor action that can be triggered by a key binding.
type Action string

const (
	// No action
	ActionNone Action = ""

	// Mode switching
	ActionNormalMode Action = "normal_mode"
	ActionInsertMode Action = "insert_mode"
	ActionAppendMode Action = "append_mode"
	ActionVisualMode Action = "visual_mode"

	// Cursor movement
	ActionMoveUp        Action = "move_up"
	ActionMoveDown      Action = "move_down"
	ActionMoveLeft      Action = "move_left"
	ActionMoveRight     Action = "move_right"
	ActionMoveWordNext  Action = "move_word_next"
	ActionMoveWordPrev  Action = "move_word_prev"
	ActionGotoLineStart Action = "goto_line_start"
	ActionGotoLineEnd   Action = "goto_line_end"
	ActionGotoTop       Action = "goto_top"
	ActionGotoBottom    Action = "goto_bottom"
	ActionPageUp        Action = "page_up"
	ActionPageDown      Action = "page_down"

	// Editing - Normal mode
	ActionDeleteChar     Action = "delete_char"
	ActionDeleteCharBack Action = "delete_char_back"
	ActionDeleteLine     Action = "delete_line"
	ActionDeleteToEnd    Action = "delete_to_end"
	ActionYankLine       Action = "yank_line"
	ActionYank           Action = "yank"
	ActionPaste          Action = "paste"
	ActionPasteAbove     Action = "paste_above"
	ActionUndo           Action = "undo"
	ActionRedo           Action = "redo"
	ActionJoinLines      Action = "join_lines"

	// Editing - Insert mode
	ActionInsertChar    Action = "insert_char"
	ActionInsertNewline Action = "insert_newline"
	ActionBackspace     Action = "backspace"
	ActionDelete        Action = "delete"
	ActionInsertTab     Action = "insert_tab"

	// Line operations
	ActionOpenBelow Action = "open_below"
	ActionOpenAbove Action = "open_above"

	// Operators (take a motion)
	ActionOperatorDelete Action = "operator_delete"
	ActionOperatorYank   Action = "operator_yank"
	ActionOperatorChange Action = "operator_change"

	// General
	ActionQuit        Action = "quit"
	ActionForceQuit   Action = "force_quit"
	ActionSave        Action = "save"
	ActionSaveQuit    Action = "save_quit"
	ActionToggleHelp  Action = "toggle_help"
	ActionRefreshRate Action = "refresh_rate"

	// UI toggles
	ActionToggleLineNumbers Action = "toggle_line_numbers"
	ActionToggleWrap        Action = "toggle_wrap"
)

// ActionMetadata contains information about an action.
type ActionMetadata struct {
	Name        string
	Description string
	IsMotion    bool // Can be used with a count (e.g., 5j)
	IsOperator  bool // Takes a motion (e.g., d, y, c)
	Repeatable  bool // Can be repeated with .
}

// actionRegistry maps actions to their metadata.
var actionRegistry = map[Action]ActionMetadata{
	// Mode switching
	ActionNormalMode: {"Normal Mode", "Switch to normal mode", false, false, false},
	ActionInsertMode: {"Insert Mode", "Switch to insert mode", false, false, false},
	ActionAppendMode: {"Append Mode", "Insert after cursor", false, false, false},
	ActionVisualMode: {"Visual Mode", "Switch to visual mode", false, false, false},

	// Cursor movement (all are motions)
	ActionMoveUp:        {"Move Up", "Move cursor up", true, false, false},
	ActionMoveDown:      {"Move Down", "Move cursor down", true, false, false},
	ActionMoveLeft:      {"Move Left", "Move cursor left", true, false, false},
	ActionMoveRight:     {"Move Right", "Move cursor right", true, false, false},
	ActionMoveWordNext:  {"Word Next", "Move to next word", true, false, false},
	ActionMoveWordPrev:  {"Word Prev", "Move to previous word", true, false, false},
	ActionGotoLineStart: {"Line Start", "Go to start of line", true, false, false},
	ActionGotoLineEnd:   {"Line End", "Go to end of line", true, false, false},
	ActionGotoTop:       {"Go to Top", "Go to first line", true, false, false},
	ActionGotoBottom:    {"Go to Bottom", "Go to last line", true, false, false},
	ActionPageUp:        {"Page Up", "Move page up", true, false, false},
	ActionPageDown:      {"Page Down", "Move page down", true, false, false},

	// Editing - Normal mode
	ActionDeleteChar:     {"Delete Char", "Delete character under cursor", false, false, true},
	ActionDeleteCharBack: {"Delete Back", "Delete character before cursor", false, false, true},
	ActionDeleteLine:     {"Delete Line", "Delete current line", false, false, true},
	ActionDeleteToEnd:    {"Delete to End", "Delete to end of line", false, false, true},
	ActionYankLine:       {"Yank Line", "Copy current line", false, false, false},
	ActionYank:           {"Yank", "Copy selection", false, false, false},
	ActionPaste:          {"Paste", "Paste after cursor", false, false, true},
	ActionPasteAbove:     {"Paste Above", "Paste before cursor", false, false, true},
	ActionUndo:           {"Undo", "Undo last change", false, false, false},
	ActionRedo:           {"Redo", "Redo last undone change", false, false, false},
	ActionJoinLines:      {"Join Lines", "Join current and next line", false, false, true},

	// Editing - Insert mode
	ActionInsertChar:    {"Insert Char", "Insert character", false, false, false},
	ActionInsertNewline: {"Insert Newline", "Insert new line", false, false, false},
	ActionBackspace:     {"Backspace", "Delete character before cursor", false, false, false},
	ActionDelete:        {"Delete", "Delete character under cursor", false, false, false},
	ActionInsertTab:     {"Insert Tab", "Insert tab/spaces", false, false, false},

	// Line operations
	ActionOpenBelow: {"Open Below", "Insert line below", false, false, true},
	ActionOpenAbove: {"Open Above", "Insert line above", false, false, true},

	// Operators
	ActionOperatorDelete: {"Delete Operator", "Delete with motion", false, true, true},
	ActionOperatorYank:   {"Yank Operator", "Yank with motion", false, true, false},
	ActionOperatorChange: {"Change Operator", "Change with motion", false, true, true},

	// General
	ActionQuit:        {"Quit", "Quit editor", false, false, false},
	ActionForceQuit:   {"Force Quit", "Quit without saving", false, false, false},
	ActionSave:        {"Save", "Save file", false, false, false},
	ActionSaveQuit:    {"Save & Quit", "Save and quit", false, false, false},
	ActionToggleHelp:  {"Toggle Help", "Show/hide help", false, false, false},
	ActionRefreshRate: {"Refresh Rates", "Refresh currency rates", false, false, false},

	// UI toggles
	ActionToggleLineNumbers: {"Toggle Line Numbers", "Show/hide line numbers", false, false, false},
	ActionToggleWrap:        {"Toggle Wrap", "Toggle line wrapping", false, false, false},
}

// Metadata returns the metadata for an action.
func (a Action) Metadata() ActionMetadata {
	if meta, ok := actionRegistry[a]; ok {
		return meta
	}
	return ActionMetadata{Name: string(a), Description: "Unknown action"}
}

// IsMotion returns true if this action can be used with a count.
func (a Action) IsMotion() bool {
	return a.Metadata().IsMotion
}

// IsOperator returns true if this action takes a motion.
func (a Action) IsOperator() bool {
	return a.Metadata().IsOperator
}

// IsRepeatable returns true if this action can be repeated with dot.
func (a Action) IsRepeatable() bool {
	return a.Metadata().Repeatable
}

// String returns the action as a string.
func (a Action) String() string {
	return string(a)
}

// ParseAction converts a string to an Action.
func ParseAction(s string) Action {
	action := Action(s)
	if _, ok := actionRegistry[action]; ok {
		return action
	}
	return ActionNone
}

// AllActions returns all registered actions.
func AllActions() []Action {
	actions := make([]Action, 0, len(actionRegistry))
	for action := range actionRegistry {
		actions = append(actions, action)
	}
	return actions
}

// MotionActions returns all actions that are motions.
func MotionActions() []Action {
	var actions []Action
	for action, meta := range actionRegistry {
		if meta.IsMotion {
			actions = append(actions, action)
		}
	}
	return actions
}

// OperatorActions returns all actions that are operators.
func OperatorActions() []Action {
	var actions []Action
	for action, meta := range actionRegistry {
		if meta.IsOperator {
			actions = append(actions, action)
		}
	}
	return actions
}
