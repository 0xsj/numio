// internal/tui/editor.go

package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Editor is a vim-style text editor component.
type Editor struct {
	lines       []string
	cursor      Cursor
	mode        Mode
	keymap      KeyMap
	styles      Styles
	highlighter *Highlighter

	// Viewport
	width      int
	height     int
	scrollY    int
	showLines  bool

	// Editing state
	yankBuffer string
	lastKey    string
	modified   bool

	// Undo/Redo
	undoStack []EditorState
	redoStack []EditorState
}

// Cursor represents the cursor position.
type Cursor struct {
	Line int
	Col  int
}

// EditorState represents a snapshot for undo/redo.
type EditorState struct {
	Lines  []string
	Cursor Cursor
}

// NewEditor creates a new editor instance.
func NewEditor(styles Styles, keymap KeyMap) *Editor {
	return &Editor{
		lines:       []string{""},
		cursor:      Cursor{Line: 0, Col: 0},
		mode:        ModeNormal,
		keymap:      keymap,
		styles:      styles,
		highlighter: NewHighlighter(styles),
		width:       80,
		height:      24,
		scrollY:     0,
		showLines:   true,
		yankBuffer:  "",
		lastKey:     "",
		modified:    false,
		undoStack:   nil,
		redoStack:   nil,
	}
}

// SetSize updates the editor dimensions.
func (e *Editor) SetSize(width, height int) {
	e.width = width
	e.height = height
	e.ensureCursorVisible()
}

// SetContent sets the editor content.
func (e *Editor) SetContent(content string) {
	e.saveUndo()
	if content == "" {
		e.lines = []string{""}
	} else {
		e.lines = strings.Split(content, "\n")
	}
	e.cursor = Cursor{Line: 0, Col: 0}
	e.modified = false
	e.scrollY = 0
}

// Content returns the current editor content.
func (e *Editor) Content() string {
	return strings.Join(e.lines, "\n")
}

// Lines returns all lines.
func (e *Editor) Lines() []string {
	return e.lines
}

// Mode returns the current editor mode.
func (e *Editor) Mode() Mode {
	return e.mode
}

// Cursor returns the current cursor position.
func (e *Editor) CursorPos() Cursor {
	return e.cursor
}

// Modified returns whether the content has been modified.
func (e *Editor) Modified() bool {
	return e.modified
}

// LineCount returns the number of lines.
func (e *Editor) LineCount() int {
	return len(e.lines)
}

// ToggleLineNumbers toggles line number display.
func (e *Editor) ToggleLineNumbers() {
	e.showLines = !e.showLines
}

// ════════════════════════════════════════════════════════════════
// UPDATE (Key handling)
// ════════════════════════════════════════════════════════════════

// Update handles key events and returns any command.
func (e *Editor) Update(msg tea.KeyMsg) tea.Cmd {
	e.lastKey = msg.String()

	switch e.mode {
	case ModeNormal:
		return e.updateNormal(msg)
	case ModeInsert:
		return e.updateInsert(msg)
	case ModeVisual:
		return e.updateVisual(msg)
	}
	return nil
}

func (e *Editor) updateNormal(msg tea.KeyMsg) tea.Cmd {
	switch {
	// Mode switching
	case key.Matches(msg, e.keymap.InsertMode):
		e.mode = ModeInsert
	case key.Matches(msg, e.keymap.InsertModeAppend):
		e.mode = ModeInsert
		e.cursorRight()
	case key.Matches(msg, e.keymap.InsertModeStart):
		e.mode = ModeInsert
		e.cursor.Col = 0
	case key.Matches(msg, e.keymap.InsertModeEnd):
		e.mode = ModeInsert
		e.cursor.Col = len(e.currentLine())
	case key.Matches(msg, e.keymap.VisualMode):
		e.mode = ModeVisual

	// Navigation
	case key.Matches(msg, e.keymap.Up):
		e.cursorUp()
	case key.Matches(msg, e.keymap.Down):
		e.cursorDown()
	case key.Matches(msg, e.keymap.Left):
		e.cursorLeft()
	case key.Matches(msg, e.keymap.Right):
		e.cursorRight()
	case key.Matches(msg, e.keymap.Home):
		e.cursor.Col = 0
	case key.Matches(msg, e.keymap.End):
		e.cursor.Col = len(e.currentLine())
	case key.Matches(msg, e.keymap.PageUp):
		e.pageUp()
	case key.Matches(msg, e.keymap.PageDown):
		e.pageDown()
	case key.Matches(msg, e.keymap.GotoBottom):
		e.cursor.Line = len(e.lines) - 1
		e.cursor.Col = 0
		e.ensureCursorVisible()
	case key.Matches(msg, e.keymap.WordNext):
		e.wordNext()
	case key.Matches(msg, e.keymap.WordPrev):
		e.wordPrev()

	// Editing
	case key.Matches(msg, e.keymap.NewLineBelow):
		e.saveUndo()
		e.insertLineBelow()
		e.mode = ModeInsert
	case key.Matches(msg, e.keymap.NewLineAbove):
		e.saveUndo()
		e.insertLineAbove()
		e.mode = ModeInsert
	case key.Matches(msg, e.keymap.DeleteChar):
		e.saveUndo()
		e.deleteChar()
	case key.Matches(msg, e.keymap.DeleteLine):
		// Need to wait for second 'd' for dd
		if e.lastKey == "d" {
			e.saveUndo()
			e.deleteLine()
		}
	case key.Matches(msg, e.keymap.Yank):
		// Need to wait for second 'y' for yy
		if e.lastKey == "y" {
			e.yankLine()
		}
	case key.Matches(msg, e.keymap.Paste):
		e.saveUndo()
		e.paste()
	case key.Matches(msg, e.keymap.Undo):
		e.undo()
	case key.Matches(msg, e.keymap.Redo):
		e.redo()

	// Handle 'gg' for go to top
	case msg.String() == "g":
		if e.lastKey == "g" {
			e.cursor.Line = 0
			e.cursor.Col = 0
			e.ensureCursorVisible()
		}
	}

	return nil
}

func (e *Editor) updateInsert(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, e.keymap.NormalMode):
		e.mode = ModeNormal
		// Adjust cursor when leaving insert mode
		if e.cursor.Col > 0 {
			e.cursor.Col--
		}

	case key.Matches(msg, e.keymap.Backspace):
		e.saveUndo()
		e.backspace()

	case key.Matches(msg, e.keymap.Delete):
		e.saveUndo()
		e.deleteChar()

	case key.Matches(msg, e.keymap.Enter):
		e.saveUndo()
		e.insertNewline()

	case key.Matches(msg, e.keymap.Tab):
		e.saveUndo()
		e.insertText("  ") // 2 spaces for tab

	// Navigation in insert mode
	case key.Matches(msg, e.keymap.Up):
		e.cursorUp()
	case key.Matches(msg, e.keymap.Down):
		e.cursorDown()
	case key.Matches(msg, e.keymap.Left):
		e.cursorLeft()
	case key.Matches(msg, e.keymap.Right):
		e.cursorRight()
	case key.Matches(msg, e.keymap.Home):
		e.cursor.Col = 0
	case key.Matches(msg, e.keymap.End):
		e.cursor.Col = len(e.currentLine())

	default:
		// Insert regular characters
		if len(msg.Runes) > 0 {
			e.saveUndo()
			e.insertText(string(msg.Runes))
		}
	}

	return nil
}

func (e *Editor) updateVisual(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, e.keymap.NormalMode):
		e.mode = ModeNormal

	// Navigation
	case key.Matches(msg, e.keymap.Up):
		e.cursorUp()
	case key.Matches(msg, e.keymap.Down):
		e.cursorDown()
	case key.Matches(msg, e.keymap.Left):
		e.cursorLeft()
	case key.Matches(msg, e.keymap.Right):
		e.cursorRight()
	}

	return nil
}

// ════════════════════════════════════════════════════════════════
// CURSOR MOVEMENT
// ════════════════════════════════════════════════════════════════

func (e *Editor) cursorUp() {
	if e.cursor.Line > 0 {
		e.cursor.Line--
		e.clampCursorCol()
		e.ensureCursorVisible()
	}
}

func (e *Editor) cursorDown() {
	if e.cursor.Line < len(e.lines)-1 {
		e.cursor.Line++
		e.clampCursorCol()
		e.ensureCursorVisible()
	}
}

func (e *Editor) cursorLeft() {
	if e.cursor.Col > 0 {
		e.cursor.Col--
	}
}

func (e *Editor) cursorRight() {
	line := e.currentLine()
	maxCol := len(line)
	if e.mode == ModeNormal && maxCol > 0 {
		maxCol--
	}
	if e.cursor.Col < maxCol {
		e.cursor.Col++
	}
}

func (e *Editor) pageUp() {
	e.cursor.Line -= e.height - 2
	if e.cursor.Line < 0 {
		e.cursor.Line = 0
	}
	e.clampCursorCol()
	e.ensureCursorVisible()
}

func (e *Editor) pageDown() {
	e.cursor.Line += e.height - 2
	if e.cursor.Line >= len(e.lines) {
		e.cursor.Line = len(e.lines) - 1
	}
	e.clampCursorCol()
	e.ensureCursorVisible()
}

func (e *Editor) wordNext() {
	line := e.currentLine()
	col := e.cursor.Col

	// Skip current word
	for col < len(line) && !isWordChar(line[col]) {
		col++
	}
	for col < len(line) && isWordChar(line[col]) {
		col++
	}
	// Skip whitespace
	for col < len(line) && !isWordChar(line[col]) {
		col++
	}

	if col >= len(line) && e.cursor.Line < len(e.lines)-1 {
		// Move to next line
		e.cursor.Line++
		e.cursor.Col = 0
	} else {
		e.cursor.Col = col
	}
}

func (e *Editor) wordPrev() {
	line := e.currentLine()
	col := e.cursor.Col

	if col == 0 && e.cursor.Line > 0 {
		// Move to end of previous line
		e.cursor.Line--
		e.cursor.Col = len(e.lines[e.cursor.Line])
		return
	}

	// Skip whitespace backwards
	for col > 0 && !isWordChar(line[col-1]) {
		col--
	}
	// Skip word backwards
	for col > 0 && isWordChar(line[col-1]) {
		col--
	}

	e.cursor.Col = col
}

func (e *Editor) clampCursorCol() {
	line := e.currentLine()
	maxCol := len(line)
	if e.mode == ModeNormal && maxCol > 0 {
		maxCol--
	}
	if e.cursor.Col > maxCol {
		e.cursor.Col = maxCol
	}
	if e.cursor.Col < 0 {
		e.cursor.Col = 0
	}
}

func (e *Editor) ensureCursorVisible() {
	// Scroll up if cursor is above viewport
	if e.cursor.Line < e.scrollY {
		e.scrollY = e.cursor.Line
	}
	// Scroll down if cursor is below viewport
	visibleLines := e.height - 1
	if e.cursor.Line >= e.scrollY+visibleLines {
		e.scrollY = e.cursor.Line - visibleLines + 1
	}
}

// ════════════════════════════════════════════════════════════════
// TEXT EDITING
// ════════════════════════════════════════════════════════════════

func (e *Editor) insertText(text string) {
	line := e.currentLine()
	newLine := line[:e.cursor.Col] + text + line[e.cursor.Col:]
	e.lines[e.cursor.Line] = newLine
	e.cursor.Col += len(text)
	e.modified = true
}

func (e *Editor) insertNewline() {
	line := e.currentLine()
	before := line[:e.cursor.Col]
	after := line[e.cursor.Col:]

	e.lines[e.cursor.Line] = before

	// Insert new line after current
	newLines := make([]string, len(e.lines)+1)
	copy(newLines[:e.cursor.Line+1], e.lines[:e.cursor.Line+1])
	newLines[e.cursor.Line+1] = after
	copy(newLines[e.cursor.Line+2:], e.lines[e.cursor.Line+1:])
	e.lines = newLines

	e.cursor.Line++
	e.cursor.Col = 0
	e.modified = true
	e.ensureCursorVisible()
}

func (e *Editor) insertLineBelow() {
	newLines := make([]string, len(e.lines)+1)
	copy(newLines[:e.cursor.Line+1], e.lines[:e.cursor.Line+1])
	newLines[e.cursor.Line+1] = ""
	copy(newLines[e.cursor.Line+2:], e.lines[e.cursor.Line+1:])
	e.lines = newLines

	e.cursor.Line++
	e.cursor.Col = 0
	e.modified = true
	e.ensureCursorVisible()
}

func (e *Editor) insertLineAbove() {
	newLines := make([]string, len(e.lines)+1)
	copy(newLines[:e.cursor.Line], e.lines[:e.cursor.Line])
	newLines[e.cursor.Line] = ""
	copy(newLines[e.cursor.Line+1:], e.lines[e.cursor.Line:])
	e.lines = newLines

	e.cursor.Col = 0
	e.modified = true
	e.ensureCursorVisible()
}

func (e *Editor) backspace() {
	if e.cursor.Col > 0 {
		line := e.currentLine()
		e.lines[e.cursor.Line] = line[:e.cursor.Col-1] + line[e.cursor.Col:]
		e.cursor.Col--
		e.modified = true
	} else if e.cursor.Line > 0 {
		// Join with previous line
		prevLine := e.lines[e.cursor.Line-1]
		currLine := e.currentLine()
		e.lines[e.cursor.Line-1] = prevLine + currLine

		// Remove current line
		e.lines = append(e.lines[:e.cursor.Line], e.lines[e.cursor.Line+1:]...)

		e.cursor.Line--
		e.cursor.Col = len(prevLine)
		e.modified = true
		e.ensureCursorVisible()
	}
}

func (e *Editor) deleteChar() {
	line := e.currentLine()
	if e.cursor.Col < len(line) {
		e.lines[e.cursor.Line] = line[:e.cursor.Col] + line[e.cursor.Col+1:]
		e.modified = true
	} else if e.cursor.Line < len(e.lines)-1 {
		// Join with next line
		nextLine := e.lines[e.cursor.Line+1]
		e.lines[e.cursor.Line] = line + nextLine

		// Remove next line
		e.lines = append(e.lines[:e.cursor.Line+1], e.lines[e.cursor.Line+2:]...)
		e.modified = true
	}
}

func (e *Editor) deleteLine() {
	e.yankBuffer = e.currentLine()

	if len(e.lines) == 1 {
		e.lines[0] = ""
		e.cursor.Col = 0
	} else {
		e.lines = append(e.lines[:e.cursor.Line], e.lines[e.cursor.Line+1:]...)
		if e.cursor.Line >= len(e.lines) {
			e.cursor.Line = len(e.lines) - 1
		}
		e.clampCursorCol()
	}
	e.modified = true
	e.ensureCursorVisible()
}

func (e *Editor) yankLine() {
	e.yankBuffer = e.currentLine()
}

func (e *Editor) paste() {
	if e.yankBuffer == "" {
		return
	}

	// Paste below current line
	newLines := make([]string, len(e.lines)+1)
	copy(newLines[:e.cursor.Line+1], e.lines[:e.cursor.Line+1])
	newLines[e.cursor.Line+1] = e.yankBuffer
	copy(newLines[e.cursor.Line+2:], e.lines[e.cursor.Line+1:])
	e.lines = newLines

	e.cursor.Line++
	e.cursor.Col = 0
	e.modified = true
	e.ensureCursorVisible()
}

// ════════════════════════════════════════════════════════════════
// UNDO / REDO
// ════════════════════════════════════════════════════════════════

func (e *Editor) saveUndo() {
	state := EditorState{
		Lines:  make([]string, len(e.lines)),
		Cursor: e.cursor,
	}
	copy(state.Lines, e.lines)
	e.undoStack = append(e.undoStack, state)

	// Limit undo stack size
	if len(e.undoStack) > 100 {
		e.undoStack = e.undoStack[1:]
	}

	// Clear redo stack on new change
	e.redoStack = nil
}

func (e *Editor) undo() {
	if len(e.undoStack) == 0 {
		return
	}

	// Save current state to redo stack
	redoState := EditorState{
		Lines:  make([]string, len(e.lines)),
		Cursor: e.cursor,
	}
	copy(redoState.Lines, e.lines)
	e.redoStack = append(e.redoStack, redoState)

	// Restore from undo stack
	state := e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]

	e.lines = state.Lines
	e.cursor = state.Cursor
	e.ensureCursorVisible()
}

func (e *Editor) redo() {
	if len(e.redoStack) == 0 {
		return
	}

	// Save current state to undo stack
	undoState := EditorState{
		Lines:  make([]string, len(e.lines)),
		Cursor: e.cursor,
	}
	copy(undoState.Lines, e.lines)
	e.undoStack = append(e.undoStack, undoState)

	// Restore from redo stack
	state := e.redoStack[len(e.redoStack)-1]
	e.redoStack = e.redoStack[:len(e.redoStack)-1]

	e.lines = state.Lines
	e.cursor = state.Cursor
	e.ensureCursorVisible()
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func (e *Editor) currentLine() string {
	if e.cursor.Line >= 0 && e.cursor.Line < len(e.lines) {
		return e.lines[e.cursor.Line]
	}
	return ""
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

// ════════════════════════════════════════════════════════════════
// RENDER
// ════════════════════════════════════════════════════════════════

// Render returns the rendered editor content.
func (e *Editor) Render() string {
	var result strings.Builder

	visibleLines := e.height - 1
	lineNumWidth := 0
	_ = lineNumWidth
	if e.showLines {
		lineNumWidth = 5 // "  1 │"
	}

	for i := 0; i < visibleLines; i++ {
		lineIdx := e.scrollY + i

		if lineIdx < len(e.lines) {
			// Line number
			if e.showLines {
				isCurrent := lineIdx == e.cursor.Line
				lineNum := e.highlighter.HighlightLineNumber(lineIdx+1, isCurrent)
				result.WriteString(lineNum)
				result.WriteString(" │")
			}

			// Line content with syntax highlighting
			line := e.lines[lineIdx]
			highlighted := e.highlighter.Highlight(line)

			// Add cursor if on this line
			if lineIdx == e.cursor.Line {
				highlighted = e.renderCursor(line, highlighted)
			}

			result.WriteString(highlighted)
		} else {
			// Empty line (show ~ like vim)
			if e.showLines {
				result.WriteString("     │")
			}
			tildeStyle := e.styles.Comment
			result.WriteString(tildeStyle.Render("~"))
		}

		if i < visibleLines-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// renderCursor adds cursor highlighting to the line.
func (e *Editor) renderCursor(original, highlighted string) string {
	if e.mode != ModeInsert && e.mode != ModeNormal {
		return highlighted
	}

	col := e.cursor.Col
	if col > len(original) {
		col = len(original)
	}

	// For simplicity, render cursor as a block character
	// In a full implementation, you'd track positions through highlighting
	if col < len(original) {
		char := string(original[col])
		cursorChar := e.styles.Cursor.Render(char)

		// This is a simplified approach - full implementation would need
		// to track ANSI escape code positions
		before := original[:col]
		after := original[col+1:]
		return e.highlighter.Highlight(before) + cursorChar + e.highlighter.Highlight(after)
	}

	// Cursor at end of line
	return highlighted + e.styles.Cursor.Render(" ")
}