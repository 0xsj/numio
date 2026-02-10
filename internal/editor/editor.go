package editor

import (
	"github.com/0xsj/numio/internal/rpc"
	"github.com/0xsj/numio/pkg/engine"
)

// ════════════════════════════════════════════════════════════════
// EDITOR
// ════════════════════════════════════════════════════════════════

// Editor is a headless editor that manages text, cursor, and evaluation.
// It can be used by any UI frontend (TUI, native, web).
type Editor struct {
	buffer    *Buffer
	cursor    *Cursor
	selection *Selection
	viewport  *Viewport
	history   *History
	renderer  *Renderer

	mode   Mode
	engine *engine.Engine

	// Results cache: line index -> styled result spans
	results map[int][]rpc.Span

	// Errors cache: line index -> error message
	errors map[int]string

	// Yank buffer for copy/paste
	yankBuffer string
	yankIsLine bool

	// Dirty flag
	dirty bool
}

// NewEditor creates a new editor with default settings.
func NewEditor() *Editor {
	return &Editor{
		buffer:     NewBuffer(),
		cursor:     NewCursor(),
		selection:  nil,
		viewport:   NewViewport(80, 24),
		history:    NewHistory(),
		renderer:   NewRenderer(),
		mode:       ModeNormal,
		engine:     engine.New(),
		results:    make(map[int][]rpc.Span),
		errors:     make(map[int]string),
		yankBuffer: "",
		yankIsLine: false,
		dirty:      false,
	}
}

// NewEditorWithEngine creates an editor with a specific engine.
func NewEditorWithEngine(eng *engine.Engine) *Editor {
	e := NewEditor()
	e.engine = eng
	return e
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// Buffer returns the text buffer.
func (e *Editor) Buffer() *Buffer {
	return e.buffer
}

// Cursor returns the cursor.
func (e *Editor) Cursor() *Cursor {
	return e.cursor
}

// Selection returns the current selection, or nil if none.
func (e *Editor) Selection() *Selection {
	return e.selection
}

// Viewport returns the viewport.
func (e *Editor) Viewport() *Viewport {
	return e.viewport
}

// Mode returns the current editor mode.
func (e *Editor) Mode() Mode {
	return e.mode
}

// Engine returns the evaluation engine.
func (e *Editor) Engine() *engine.Engine {
	return e.engine
}

// IsDirty returns true if there are unsaved changes.
func (e *Editor) IsDirty() bool {
	return e.dirty
}

// ════════════════════════════════════════════════════════════════
// MODE MANAGEMENT
// ════════════════════════════════════════════════════════════════

// SetMode changes the editor mode.
func (e *Editor) SetMode(mode Mode) {
	prevMode := e.mode
	e.mode = mode

	switch mode {
	case ModeNormal:
		// Clear selection when entering normal mode
		e.selection = nil
		// Adjust cursor if past end of line
		e.clampCursor()

	case ModeInsert:
		// Clear selection
		e.selection = nil

	case ModeVisual:
		// Start selection at cursor
		if prevMode != ModeVisual {
			row, col := e.cursor.Position()
			e.selection = NewSelection(row, col)
		}
	}
}

// EnterInsertMode enters insert mode.
func (e *Editor) EnterInsertMode() {
	e.SetMode(ModeInsert)
}

// EnterInsertModeAppend enters insert mode after the cursor.
func (e *Editor) EnterInsertModeAppend() {
	e.cursor.MoveRight(1)
	e.clampCursor()
	e.SetMode(ModeInsert)
}

// EnterInsertModeLineEnd enters insert mode at end of line.
func (e *Editor) EnterInsertModeLineEnd() {
	lineLen := e.buffer.LineLength(e.cursor.Row())
	e.cursor.SetCol(lineLen)
	e.SetMode(ModeInsert)
}

// EnterInsertModeLineStart enters insert mode at start of line.
func (e *Editor) EnterInsertModeLineStart() {
	e.cursor.SetCol(0)
	e.SetMode(ModeInsert)
}

// EnterNormalMode enters normal mode.
func (e *Editor) EnterNormalMode() {
	e.SetMode(ModeNormal)
	// Move cursor back one if past end
	if e.cursor.Col() > 0 {
		lineLen := e.buffer.LineLength(e.cursor.Row())
		if lineLen > 0 && e.cursor.Col() >= lineLen {
			e.cursor.SetCol(lineLen - 1)
		}
	}
}

// EnterVisualMode enters visual (character) selection mode.
func (e *Editor) EnterVisualMode() {
	e.SetMode(ModeVisual)
}

// EnterVisualLineMode enters visual line selection mode.
func (e *Editor) EnterVisualLineMode() {
	row, col := e.cursor.Position()
	e.selection = NewSelectionWithMode(row, col, SelectionLine)
	e.mode = ModeVisual
}

// ════════════════════════════════════════════════════════════════
// CURSOR MOVEMENT
// ════════════════════════════════════════════════════════════════

// MoveUp moves the cursor up.
func (e *Editor) MoveUp(n int) {
	e.cursor.MoveUp(n)
	lineLen := e.buffer.LineLength(e.cursor.Row())
	e.cursor.ApplyPreferredCol(lineLen, e.mode == ModeInsert)
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveDown moves the cursor down, expanding the buffer if needed.
func (e *Editor) MoveDown(n int) {
	e.cursor.MoveDown(n)
	e.buffer.EnsureLines(e.cursor.Row())
	lineLen := e.buffer.LineLength(e.cursor.Row())
	e.cursor.ApplyPreferredCol(lineLen, e.mode == ModeInsert)
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveLeft moves the cursor left.
func (e *Editor) MoveLeft(n int) {
	e.cursor.MoveLeft(n)
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveRight moves the cursor right.
func (e *Editor) MoveRight(n int) {
	e.cursor.MoveRight(n)
	e.clampCursor()
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveToLineStart moves cursor to start of line.
func (e *Editor) MoveToLineStart() {
	e.cursor.MoveToLineStart()
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveToLineEnd moves cursor to end of line.
func (e *Editor) MoveToLineEnd() {
	lineLen := e.buffer.LineLength(e.cursor.Row())
	e.cursor.MoveToLineEnd(lineLen)
	e.clampCursor()
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveToFirstNonBlank moves cursor to first non-whitespace character.
func (e *Editor) MoveToFirstNonBlank() {
	line := e.buffer.Line(e.cursor.Row())
	e.cursor.MoveToFirstNonBlank(line)
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveWordForward moves cursor to next word.
func (e *Editor) MoveWordForward(n int) {
	for i := 0; i < n; i++ {
		line := e.buffer.Line(e.cursor.Row())
		newCol := e.cursor.MoveWordForward(line)

		// If at end of line, move to next line
		if newCol >= len(line) && e.cursor.Row() < e.buffer.LastRow() {
			e.cursor.SetRow(e.cursor.Row() + 1)
			e.cursor.SetCol(0)
		}
	}
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveWordBackward moves cursor to previous word.
func (e *Editor) MoveWordBackward(n int) {
	for i := 0; i < n; i++ {
		if e.cursor.Col() == 0 && e.cursor.Row() > 0 {
			e.cursor.SetRow(e.cursor.Row() - 1)
			lineLen := e.buffer.LineLength(e.cursor.Row())
			e.cursor.SetCol(lineLen)
		}
		line := e.buffer.Line(e.cursor.Row())
		e.cursor.MoveWordBackward(line)
	}
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveToTop moves cursor to first line.
func (e *Editor) MoveToTop() {
	e.cursor.SetPosition(0, 0)
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveToBottom moves cursor to last line.
func (e *Editor) MoveToBottom() {
	lastRow := e.buffer.LastRow()
	e.cursor.SetRow(lastRow)
	e.cursor.SetCol(0)
	e.updateSelection()
	e.ensureCursorVisible()
}

// MoveToLine moves cursor to a specific line (1-indexed).
func (e *Editor) MoveToLine(lineNum int) {
	row := lineNum - 1
	if row < 0 {
		row = 0
	}
	if row > e.buffer.LastRow() {
		row = e.buffer.LastRow()
	}
	e.cursor.SetRow(row)
	e.cursor.SetCol(0)
	e.updateSelection()
	e.ensureCursorVisible()
}

// ════════════════════════════════════════════════════════════════
// TEXT EDITING
// ════════════════════════════════════════════════════════════════

// InsertChar inserts a character at the cursor.
func (e *Editor) InsertChar(ch rune) {
	e.saveUndo()
	row, col := e.cursor.Position()
	e.buffer.EnsureLines(row)
	e.buffer.InsertChar(row, col, ch)
	e.cursor.MoveRight(1)
	e.dirty = true
	e.evaluateLine(row)
	e.ensureCursorVisible()
}

// InsertString inserts a string at the cursor.
func (e *Editor) InsertString(s string) {
	if s == "" {
		return
	}
	e.saveUndo()
	row, col := e.cursor.Position()
	e.buffer.EnsureLines(row)
	e.buffer.InsertString(row, col, s)
	e.cursor.SetCol(col + len(s))
	e.dirty = true
	e.evaluateLine(row)
	e.ensureCursorVisible()
}

// InsertNewline inserts a newline, splitting the current line.
func (e *Editor) InsertNewline() {
	e.saveUndo()
	row, col := e.cursor.Position()
	e.buffer.EnsureLines(row)
	e.buffer.SplitLine(row, col)
	e.cursor.SetPosition(row+1, 0)
	e.dirty = true
	e.shiftResults(row + 1)
	e.evaluateLine(row)
	e.evaluateLine(row + 1)
	e.ensureCursorVisible()
}

// DeleteChar deletes the character at the cursor (like 'x' in vim).
func (e *Editor) DeleteChar() {
	e.saveUndo()
	row, col := e.cursor.Position()
	e.buffer.DeleteChar(row, col)
	e.clampCursor()
	e.dirty = true
	e.evaluateLine(row)
}

// DeleteCharBack deletes the character before the cursor (backspace).
func (e *Editor) DeleteCharBack() {
	row, col := e.cursor.Position()

	if col > 0 {
		e.saveUndo()
		e.buffer.DeleteCharBack(row, col)
		e.cursor.MoveLeft(1)
		e.dirty = true
		e.evaluateLine(row)
	} else if row > 0 {
		// Join with previous line
		e.saveUndo()
		prevLineLen := e.buffer.LineLength(row - 1)
		e.buffer.JoinLines(row - 1)
		e.cursor.SetPosition(row-1, prevLineLen)
		e.dirty = true
		e.shiftResultsUp(row)
		e.evaluateLine(row - 1)
	}
	e.ensureCursorVisible()
}

// DeleteLine deletes the current line (like 'dd' in vim).
func (e *Editor) DeleteLine() {
	e.saveUndo()
	row := e.cursor.Row()

	// Yank the line first
	e.yankBuffer = e.buffer.Line(row)
	e.yankIsLine = true

	e.buffer.DeleteLine(row)
	e.clampCursorRow()
	e.clampCursor()
	e.dirty = true

	// Clear result for deleted line and shift others
	delete(e.results, row)
	delete(e.errors, row)
	e.shiftResultsUp(row + 1)

	e.evaluateLine(e.cursor.Row())
}

// DeleteToLineEnd deletes from cursor to end of line (like 'D' in vim).
func (e *Editor) DeleteToLineEnd() {
	e.saveUndo()
	row, col := e.cursor.Position()
	deleted := e.buffer.DeleteToLineEnd(row, col)
	e.yankBuffer = deleted
	e.yankIsLine = false
	e.dirty = true
	e.evaluateLine(row)
}

// DeleteWord deletes from cursor to end of word (like 'dw' in vim).
func (e *Editor) DeleteWord() {
	e.saveUndo()
	row, col := e.cursor.Position()
	deleted, _ := e.buffer.DeleteWord(row, col)
	e.yankBuffer = deleted
	e.yankIsLine = false
	e.dirty = true
	e.evaluateLine(row)
}

// ════════════════════════════════════════════════════════════════
// LINE OPERATIONS
// ════════════════════════════════════════════════════════════════

// OpenLineBelow inserts a new line below and enters insert mode (like 'o').
func (e *Editor) OpenLineBelow() {
	e.saveUndo()
	row := e.cursor.Row()
	e.buffer.InsertLine(row+1, "")
	e.cursor.SetPosition(row+1, 0)
	e.SetMode(ModeInsert)
	e.dirty = true
	e.shiftResults(row + 1)
	e.ensureCursorVisible()
}

// OpenLineAbove inserts a new line above and enters insert mode (like 'O').
func (e *Editor) OpenLineAbove() {
	e.saveUndo()
	row := e.cursor.Row()
	e.buffer.InsertLine(row, "")
	e.cursor.SetPosition(row, 0)
	e.SetMode(ModeInsert)
	e.dirty = true
	e.shiftResults(row)
	e.ensureCursorVisible()
}

// JoinLines joins the current line with the next (like 'J' in vim).
func (e *Editor) JoinLines() {
	row := e.cursor.Row()
	if row >= e.buffer.LastRow() {
		return
	}

	e.saveUndo()
	joinCol := e.buffer.JoinLines(row)
	e.cursor.SetCol(joinCol)
	e.dirty = true
	e.shiftResultsUp(row + 1)
	e.evaluateLine(row)
}

// ════════════════════════════════════════════════════════════════
// YANK / PASTE
// ════════════════════════════════════════════════════════════════

// YankLine yanks (copies) the current line.
func (e *Editor) YankLine() {
	e.yankBuffer = e.buffer.Line(e.cursor.Row())
	e.yankIsLine = true
}

// YankSelection yanks the current selection.
func (e *Editor) YankSelection() {
	if e.selection == nil {
		return
	}

	startRow, startCol, endRow, endCol := e.selection.Bounds()
	e.yankBuffer = e.buffer.GetRange(startRow, startCol, endRow, endCol)
	e.yankIsLine = e.selection.Mode() == SelectionLine
}

// Paste pastes after the cursor (like 'p' in vim).
func (e *Editor) Paste() {
	if e.yankBuffer == "" {
		return
	}

	e.saveUndo()

	if e.yankIsLine {
		// Paste line below
		row := e.cursor.Row()
		e.buffer.InsertLine(row+1, e.yankBuffer)
		e.cursor.SetPosition(row+1, 0)
		e.shiftResults(row + 1)
	} else {
		// Paste inline
		row, col := e.cursor.Position()
		e.buffer.InsertString(row, col+1, e.yankBuffer)
		e.cursor.SetCol(col + len(e.yankBuffer))
	}

	e.dirty = true
	e.evaluateLine(e.cursor.Row())
	e.ensureCursorVisible()
}

// PasteBefore pastes before the cursor (like 'P' in vim).
func (e *Editor) PasteBefore() {
	if e.yankBuffer == "" {
		return
	}

	e.saveUndo()

	if e.yankIsLine {
		// Paste line above
		row := e.cursor.Row()
		e.buffer.InsertLine(row, e.yankBuffer)
		e.cursor.SetPosition(row, 0)
		e.shiftResults(row)
	} else {
		// Paste inline before cursor
		row, col := e.cursor.Position()
		e.buffer.InsertString(row, col, e.yankBuffer)
		e.cursor.SetCol(col + len(e.yankBuffer) - 1)
	}

	e.dirty = true
	e.evaluateLine(e.cursor.Row())
	e.ensureCursorVisible()
}

// ════════════════════════════════════════════════════════════════
// UNDO / REDO
// ════════════════════════════════════════════════════════════════

func (e *Editor) saveUndo() {
	row, col := e.cursor.Position()
	e.history.Save(e.buffer.Lines(), row, col)
}

// Undo undoes the last change.
func (e *Editor) Undo() bool {
	row, col := e.cursor.Position()
	snapshot := e.history.Undo(e.buffer.Lines(), row, col)
	if snapshot == nil {
		return false
	}

	e.buffer.ReplaceLines(snapshot.Lines)
	e.cursor.SetPosition(snapshot.CursorRow, snapshot.CursorCol)
	e.clampCursor()
	e.dirty = true

	// Re-evaluate all lines
	e.evaluateAll()
	e.ensureCursorVisible()

	return true
}

// Redo redoes the last undone change.
func (e *Editor) Redo() bool {
	row, col := e.cursor.Position()
	snapshot := e.history.Redo(e.buffer.Lines(), row, col)
	if snapshot == nil {
		return false
	}

	e.buffer.ReplaceLines(snapshot.Lines)
	e.cursor.SetPosition(snapshot.CursorRow, snapshot.CursorCol)
	e.clampCursor()
	e.dirty = true

	// Re-evaluate all lines
	e.evaluateAll()
	e.ensureCursorVisible()

	return true
}

// ════════════════════════════════════════════════════════════════
// EVALUATION
// ════════════════════════════════════════════════════════════════

// evaluateLine evaluates a single line and caches the result.
func (e *Editor) evaluateLine(row int) {
	line := e.buffer.Line(row)

	// Clear previous result/error
	delete(e.results, row)
	delete(e.errors, row)

	if line == "" {
		return
	}

	// Skip comments
	if len(line) >= 2 && (line[:2] == "//" || line[0] == '#') {
		return
	}

	// Evaluate
	result := e.engine.Eval(line)

	if result.IsError() {
		e.errors[row] = result.ErrorMessage()
	} else if !result.IsEmpty() {
		e.results[row] = []rpc.Span{
			{Text: result.String(), Style: rpc.StyleResult},
		}
	}
}

// evaluateAll re-evaluates all lines.
func (e *Editor) evaluateAll() {
	e.results = make(map[int][]rpc.Span)
	e.errors = make(map[int]string)

	// Clear engine state for fresh evaluation
	e.engine.Clear()

	for row := 0; row < e.buffer.LineCount(); row++ {
		e.evaluateLine(row)
	}
}

// shiftResults shifts results down after inserting a line.
func (e *Editor) shiftResults(fromRow int) {
	newResults := make(map[int][]rpc.Span)
	newErrors := make(map[int]string)

	for row, result := range e.results {
		if row >= fromRow {
			newResults[row+1] = result
		} else {
			newResults[row] = result
		}
	}

	for row, err := range e.errors {
		if row >= fromRow {
			newErrors[row+1] = err
		} else {
			newErrors[row] = err
		}
	}

	e.results = newResults
	e.errors = newErrors
}

// shiftResultsUp shifts results up after deleting a line.
func (e *Editor) shiftResultsUp(fromRow int) {
	newResults := make(map[int][]rpc.Span)
	newErrors := make(map[int]string)

	for row, result := range e.results {
		if row >= fromRow {
			newResults[row-1] = result
		} else {
			newResults[row] = result
		}
	}

	for row, err := range e.errors {
		if row >= fromRow {
			newErrors[row-1] = err
		} else {
			newErrors[row] = err
		}
	}

	e.results = newResults
	e.errors = newErrors
}

// ════════════════════════════════════════════════════════════════
// VIEWPORT
// ════════════════════════════════════════════════════════════════

// Resize updates the viewport size.
func (e *Editor) Resize(width, height int) {
	e.viewport.SetSize(width, height)
	e.ensureCursorVisible()
}

func (e *Editor) ensureCursorVisible() {
	row, col := e.cursor.Position()
	e.viewport.EnsureCursorVisible(row, col)
}

// ════════════════════════════════════════════════════════════════
// SELECTION HELPERS
// ════════════════════════════════════════════════════════════════

func (e *Editor) updateSelection() {
	if e.mode == ModeVisual && e.selection != nil {
		row, col := e.cursor.Position()
		e.selection.SetActive(row, col)
	}
}

// ════════════════════════════════════════════════════════════════
// CURSOR HELPERS
// ════════════════════════════════════════════════════════════════

func (e *Editor) clampCursor() {
	lineLen := e.buffer.LineLength(e.cursor.Row())
	e.cursor.ClampCol(lineLen, e.mode == ModeInsert)
}

func (e *Editor) clampCursorRow() {
	maxRow := e.buffer.LastRow()
	if e.cursor.Row() > maxRow {
		e.cursor.SetRow(maxRow)
	}
}

// ════════════════════════════════════════════════════════════════
// RENDERING
// ════════════════════════════════════════════════════════════════

// Render produces the current render state for the UI.
func (e *Editor) Render() *rpc.RenderState {
	params := RenderParams{
		Buffer:    e.buffer,
		Cursor:    e.cursor,
		Selection: e.selection,
		Viewport:  e.viewport,
		Mode:      e.mode,
		Results:   e.results,
		Errors:    e.errors,
	}

	state := e.renderer.Render(params)
	state.Dirty = e.dirty

	return state
}

// ════════════════════════════════════════════════════════════════
// CONTENT ACCESS
// ════════════════════════════════════════════════════════════════

// Content returns the entire buffer content as a string.
func (e *Editor) Content() string {
	return e.buffer.String()
}

// Lines returns all buffer lines.
func (e *Editor) Lines() []string {
	return e.buffer.Lines()
}

// SetContent replaces the buffer content.
func (e *Editor) SetContent(content string) {
	e.saveUndo()
	e.buffer.Replace(content)
	e.cursor.SetPosition(0, 0)
	e.selection = nil
	e.dirty = true
	e.evaluateAll()
}

// Clear resets the editor to empty state.
func (e *Editor) Clear() {
	e.buffer.Clear()
	e.cursor.SetPosition(0, 0)
	e.selection = nil
	e.history.Clear()
	e.results = make(map[int][]rpc.Span)
	e.errors = make(map[int]string)
	e.dirty = false
	e.engine.Clear()
}

// ════════════════════════════════════════════════════════════════
// STATE MANAGEMENT
// ════════════════════════════════════════════════════════════════

// MarkClean marks the editor as clean (no unsaved changes).
func (e *Editor) MarkClean() {
	e.dirty = false
}
