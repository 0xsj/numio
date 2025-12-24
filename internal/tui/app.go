// internal/tui/app.go

package tui

import (
	"fmt"
	"strings"

	"github.com/0xsj/numio/internal/tui/keymap"
	"github.com/0xsj/numio/pkg/engine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	lineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666"))
	commentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Italic(true)
	resultStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7ee787"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f85149"))
	cursorStyle  = lipgloss.NewStyle().Reverse(true)
	tildeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#444"))
	pendingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffa657"))

	// Help styles
	helpBorderStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#79c0ff")).Padding(1, 2)
	helpTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#79c0ff"))
	helpSectionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffa657")).MarginTop(1)
	helpKeyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#79c0ff")).Width(14)
	helpDescStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
	helpFooterStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Italic(true).MarginTop(1)
)

// App is the main model
type App struct {
	lines  []string
	row    int
	col    int
	width  int
	height int
	engine *engine.Engine

	// Keymap
	keymap   *keymap.KeyMap
	showHelp bool

	// Yank buffer
	yankBuffer string

	// Undo/Redo
	undoStack []editorState
	redoStack []editorState
}

// editorState for undo/redo
type editorState struct {
	lines []string
	row   int
	col   int
}

// NewApp creates a new app
func NewApp() *App {
	// Load keymap (with user config if exists)
	km, _ := keymap.LoadOrCreate(keymap.DefaultConfigPath())

	return &App{
		lines:      []string{""},
		row:        0,
		col:        0,
		width:      80,
		height:     24,
		engine:     engine.New(),
		keymap:     km,
		showHelp:   false,
		yankBuffer: "",
		undoStack:  nil,
		redoStack:  nil,
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		return a.handleKey(msg)
	}

	return a, nil
}

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Always handle Ctrl+C as force quit
	if key == "ctrl+c" {
		return a, tea.Quit
	}

	// In insert mode, handle text input specially
	if a.keymap.CurrentMode == keymap.ModeInsert {
		return a.handleInsertKey(msg)
	}

	// Close help with any key
	if a.showHelp {
		a.showHelp = false
		return a, nil
	}

	// Process key through keymap
	cmd, ok := a.keymap.ProcessKey(key)
	if !ok {
		// No command yet (could be pending sequence or count)
		return a, nil
	}

	// Execute the command
	return a.executeCommand(cmd)
}

func (a *App) handleInsertKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Check for bound keys in insert mode
	result := a.keymap.Insert.Lookup(key)
	if result.Status == keymap.LookupFound {
		cmd := keymap.NewCommand(result.Action, 1)
		return a.executeCommand(cmd)
	}

	// Handle regular character input
	if len(msg.Runes) > 0 {
		a.saveUndo()
		for _, r := range msg.Runes {
			a.insertChar(r)
		}
	}

	return a, nil
}

func (a *App) executeCommand(cmd keymap.Command) (tea.Model, tea.Cmd) {
	count := cmd.TotalCount()

	switch cmd.Action {
	// Mode switching
	case keymap.ActionNormalMode:
		a.keymap.SetMode(keymap.ModeNormal)
		if a.col > 0 {
			a.col--
		}

	case keymap.ActionInsertMode:
		a.keymap.SetMode(keymap.ModeInsert)

	case keymap.ActionAppendMode:
		a.keymap.SetMode(keymap.ModeInsert)
		if a.col < len(a.lines[a.row]) {
			a.col++
		}

	case keymap.ActionVisualMode:
		a.keymap.SetMode(keymap.ModeVisual)

	// Cursor movement
	case keymap.ActionMoveUp:
		for i := 0; i < count; i++ {
			a.cursorUp()
		}

	case keymap.ActionMoveDown:
		for i := 0; i < count; i++ {
			a.cursorDown()
		}

	case keymap.ActionMoveLeft:
		for i := 0; i < count; i++ {
			a.cursorLeft()
		}

	case keymap.ActionMoveRight:
		for i := 0; i < count; i++ {
			a.cursorRight()
		}

	case keymap.ActionMoveWordNext:
		for i := 0; i < count; i++ {
			a.wordNext()
		}

	case keymap.ActionMoveWordPrev:
		for i := 0; i < count; i++ {
			a.wordPrev()
		}

	case keymap.ActionGotoLineStart:
		a.col = 0

	case keymap.ActionGotoLineEnd:
		a.col = len(a.lines[a.row])

	case keymap.ActionGotoTop:
		a.row = 0
		a.col = 0

	case keymap.ActionGotoBottom:
		a.row = len(a.lines) - 1
		a.clampCol()

	case keymap.ActionPageUp:
		pageSize := a.height - 4
		if pageSize < 1 {
			pageSize = 1
		}
		for i := 0; i < pageSize*count; i++ {
			a.cursorUp()
		}

	case keymap.ActionPageDown:
		pageSize := a.height - 4
		if pageSize < 1 {
			pageSize = 1
		}
		for i := 0; i < pageSize*count; i++ {
			a.cursorDown()
		}

	// Editing
	case keymap.ActionDeleteChar:
		a.saveUndo()
		for i := 0; i < count; i++ {
			a.deleteChar()
		}

	case keymap.ActionDeleteCharBack:
		a.saveUndo()
		for i := 0; i < count; i++ {
			a.backspace()
		}

	case keymap.ActionDeleteLine:
		a.saveUndo()
		for i := 0; i < count; i++ {
			a.deleteLine()
		}

	case keymap.ActionDeleteToEnd:
		a.saveUndo()
		a.deleteToEnd()

	case keymap.ActionYankLine:
		a.yankLine()

	case keymap.ActionPaste:
		a.saveUndo()
		a.paste()

	case keymap.ActionPasteAbove:
		a.saveUndo()
		a.pasteAbove()

	case keymap.ActionUndo:
		a.undo()

	case keymap.ActionRedo:
		a.redo()

	case keymap.ActionJoinLines:
		a.saveUndo()
		a.joinLines()

	case keymap.ActionOpenBelow:
		a.saveUndo()
		a.newLineBelow()
		a.keymap.SetMode(keymap.ModeInsert)

	case keymap.ActionOpenAbove:
		a.saveUndo()
		a.newLineAbove()
		a.keymap.SetMode(keymap.ModeInsert)

	// Insert mode actions
	case keymap.ActionBackspace:
		a.saveUndo()
		a.backspace()

	case keymap.ActionDelete:
		a.saveUndo()
		a.deleteChar()

	case keymap.ActionInsertNewline:
		a.saveUndo()
		a.newLine()

	case keymap.ActionInsertTab:
		a.saveUndo()
		a.insertChar(' ')
		a.insertChar(' ')

	// Operators with motions
	case keymap.ActionOperatorDelete:
		if cmd.Motion != keymap.ActionNone {
			a.saveUndo()
			a.deleteWithMotion(cmd.Motion, count)
		}

	case keymap.ActionOperatorYank:
		if cmd.Motion != keymap.ActionNone {
			a.yankWithMotion(cmd.Motion, count)
		}

	case keymap.ActionOperatorChange:
		if cmd.Motion != keymap.ActionNone {
			a.saveUndo()
			a.deleteWithMotion(cmd.Motion, count)
			a.keymap.SetMode(keymap.ModeInsert)
		}

	// General
	case keymap.ActionQuit:
		return a, tea.Quit

	case keymap.ActionForceQuit:
		return a, tea.Quit

	case keymap.ActionSave:
		// TODO: Implement save

	case keymap.ActionSaveQuit:
		// TODO: Implement save
		return a, tea.Quit

	case keymap.ActionToggleHelp:
		a.showHelp = !a.showHelp

	case keymap.ActionToggleLineNumbers:
		// TODO: Implement

	case keymap.ActionToggleWrap:
		// TODO: Implement
	}

	return a, nil
}

// ════════════════════════════════════════════════════════════════
// CURSOR MOVEMENT
// ════════════════════════════════════════════════════════════════

func (a *App) cursorUp() {
	if a.row > 0 {
		a.row--
		a.clampCol()
	}
}

func (a *App) cursorDown() {
	if a.row < len(a.lines)-1 {
		a.row++
		a.clampCol()
	}
}

func (a *App) cursorLeft() {
	if a.col > 0 {
		a.col--
	}
}

func (a *App) cursorRight() {
	maxCol := len(a.lines[a.row])
	if a.keymap.CurrentMode == keymap.ModeNormal && maxCol > 0 {
		maxCol--
	}
	if a.col < maxCol {
		a.col++
	}
}

func (a *App) clampCol() {
	maxCol := len(a.lines[a.row])
	if a.keymap.CurrentMode == keymap.ModeNormal && maxCol > 0 {
		maxCol--
	}
	if a.col > maxCol {
		a.col = maxCol
	}
	if a.col < 0 {
		a.col = 0
	}
}

func (a *App) wordNext() {
	line := a.lines[a.row]
	col := a.col

	// Skip current word
	for col < len(line) && isWordChar(line[col]) {
		col++
	}
	// Skip whitespace
	for col < len(line) && !isWordChar(line[col]) {
		col++
	}

	if col >= len(line) && a.row < len(a.lines)-1 {
		a.row++
		a.col = 0
	} else {
		a.col = col
	}
}

func (a *App) wordPrev() {
	if a.col == 0 && a.row > 0 {
		a.row--
		a.col = len(a.lines[a.row])
		return
	}

	line := a.lines[a.row]
	col := a.col

	// Move back one if at word
	if col > 0 {
		col--
	}

	// Skip whitespace backwards
	for col > 0 && !isWordChar(line[col]) {
		col--
	}
	// Skip word backwards
	for col > 0 && isWordChar(line[col-1]) {
		col--
	}

	a.col = col
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

// ════════════════════════════════════════════════════════════════
// TEXT EDITING
// ════════════════════════════════════════════════════════════════

func (a *App) insertChar(r rune) {
	line := a.lines[a.row]
	if a.col > len(line) {
		a.col = len(line)
	}
	a.lines[a.row] = line[:a.col] + string(r) + line[a.col:]
	a.col++
}

func (a *App) newLine() {
	line := a.lines[a.row]
	if a.col > len(line) {
		a.col = len(line)
	}
	before := line[:a.col]
	after := line[a.col:]
	a.lines[a.row] = before

	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row+1]...)
	newLines = append(newLines, after)
	if a.row+1 < len(a.lines) {
		newLines = append(newLines, a.lines[a.row+1:]...)
	}
	a.lines = newLines

	a.row++
	a.col = 0
}

func (a *App) newLineBelow() {
	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row+1]...)
	newLines = append(newLines, "")
	if a.row+1 < len(a.lines) {
		newLines = append(newLines, a.lines[a.row+1:]...)
	}
	a.lines = newLines
	a.row++
	a.col = 0
}

func (a *App) newLineAbove() {
	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row]...)
	newLines = append(newLines, "")
	newLines = append(newLines, a.lines[a.row:]...)
	a.lines = newLines
	a.col = 0
}

func (a *App) backspace() {
	if a.col > 0 {
		line := a.lines[a.row]
		a.lines[a.row] = line[:a.col-1] + line[a.col:]
		a.col--
	} else if a.row > 0 {
		prevLen := len(a.lines[a.row-1])
		a.lines[a.row-1] += a.lines[a.row]
		a.lines = append(a.lines[:a.row], a.lines[a.row+1:]...)
		a.row--
		a.col = prevLen
	}
}

func (a *App) deleteChar() {
	line := a.lines[a.row]
	if a.col < len(line) {
		a.lines[a.row] = line[:a.col] + line[a.col+1:]
	} else if a.row < len(a.lines)-1 {
		a.lines[a.row] = line + a.lines[a.row+1]
		a.lines = append(a.lines[:a.row+1], a.lines[a.row+2:]...)
	}
}

func (a *App) deleteLine() {
	a.yankBuffer = a.lines[a.row] + "\n"

	if len(a.lines) == 1 {
		a.lines[0] = ""
		a.col = 0
	} else {
		a.lines = append(a.lines[:a.row], a.lines[a.row+1:]...)
		if a.row >= len(a.lines) {
			a.row = len(a.lines) - 1
		}
		a.clampCol()
	}
}

func (a *App) deleteToEnd() {
	line := a.lines[a.row]
	a.yankBuffer = line[a.col:]
	a.lines[a.row] = line[:a.col]
	a.clampCol()
}

func (a *App) joinLines() {
	if a.row < len(a.lines)-1 {
		a.lines[a.row] = a.lines[a.row] + " " + strings.TrimLeft(a.lines[a.row+1], " \t")
		a.lines = append(a.lines[:a.row+1], a.lines[a.row+2:]...)
	}
}

func (a *App) yankLine() {
	a.yankBuffer = a.lines[a.row] + "\n"
}

func (a *App) paste() {
	if a.yankBuffer == "" {
		return
	}

	if strings.HasSuffix(a.yankBuffer, "\n") {
		// Paste line below
		content := strings.TrimSuffix(a.yankBuffer, "\n")
		newLines := make([]string, 0, len(a.lines)+1)
		newLines = append(newLines, a.lines[:a.row+1]...)
		newLines = append(newLines, content)
		if a.row+1 < len(a.lines) {
			newLines = append(newLines, a.lines[a.row+1:]...)
		}
		a.lines = newLines
		a.row++
		a.col = 0
	} else {
		// Paste inline
		line := a.lines[a.row]
		a.lines[a.row] = line[:a.col+1] + a.yankBuffer + line[a.col+1:]
		a.col += len(a.yankBuffer)
	}
}

func (a *App) pasteAbove() {
	if a.yankBuffer == "" {
		return
	}

	if strings.HasSuffix(a.yankBuffer, "\n") {
		// Paste line above
		content := strings.TrimSuffix(a.yankBuffer, "\n")
		newLines := make([]string, 0, len(a.lines)+1)
		newLines = append(newLines, a.lines[:a.row]...)
		newLines = append(newLines, content)
		newLines = append(newLines, a.lines[a.row:]...)
		a.lines = newLines
		a.col = 0
	} else {
		// Paste inline before cursor
		line := a.lines[a.row]
		a.lines[a.row] = line[:a.col] + a.yankBuffer + line[a.col:]
	}
}

// ════════════════════════════════════════════════════════════════
// OPERATOR + MOTION
// ════════════════════════════════════════════════════════════════

func (a *App) deleteWithMotion(motion keymap.Action, count int) {
	startRow, startCol := a.row, a.col

	// Execute motion to find end position
	for i := 0; i < count; i++ {
		a.executeMotion(motion)
	}

	endRow, endCol := a.row, a.col

	// Ensure start <= end
	if endRow < startRow || (endRow == startRow && endCol < startCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	// Delete the range
	if startRow == endRow {
		// Same line
		line := a.lines[startRow]
		if endCol > len(line) {
			endCol = len(line)
		}
		a.yankBuffer = line[startCol:endCol]
		a.lines[startRow] = line[:startCol] + line[endCol:]
		a.row = startRow
		a.col = startCol
	} else {
		// Multiple lines
		var yanked strings.Builder
		yanked.WriteString(a.lines[startRow][startCol:])
		yanked.WriteString("\n")
		for i := startRow + 1; i < endRow; i++ {
			yanked.WriteString(a.lines[i])
			yanked.WriteString("\n")
		}
		if endRow < len(a.lines) {
			yanked.WriteString(a.lines[endRow][:endCol])
		}
		a.yankBuffer = yanked.String()

		// Join lines
		newLine := a.lines[startRow][:startCol]
		if endRow < len(a.lines) {
			newLine += a.lines[endRow][endCol:]
		}
		a.lines[startRow] = newLine

		// Remove middle lines
		a.lines = append(a.lines[:startRow+1], a.lines[endRow+1:]...)

		a.row = startRow
		a.col = startCol
	}

	a.clampCol()
}

func (a *App) yankWithMotion(motion keymap.Action, count int) {
	startRow, startCol := a.row, a.col

	// Execute motion to find end position
	for i := 0; i < count; i++ {
		a.executeMotion(motion)
	}

	endRow, endCol := a.row, a.col

	// Restore position
	a.row, a.col = startRow, startCol

	// Ensure start <= end
	if endRow < startRow || (endRow == startRow && endCol < startCol) {
		startRow, endRow = endRow, startRow
		startCol, endCol = endCol, startCol
	}

	// Yank the range
	if startRow == endRow {
		line := a.lines[startRow]
		if endCol > len(line) {
			endCol = len(line)
		}
		a.yankBuffer = line[startCol:endCol]
	} else {
		var yanked strings.Builder
		yanked.WriteString(a.lines[startRow][startCol:])
		yanked.WriteString("\n")
		for i := startRow + 1; i < endRow; i++ {
			yanked.WriteString(a.lines[i])
			yanked.WriteString("\n")
		}
		if endRow < len(a.lines) {
			yanked.WriteString(a.lines[endRow][:endCol])
		}
		a.yankBuffer = yanked.String()
	}
}

func (a *App) executeMotion(motion keymap.Action) {
	switch motion {
	case keymap.ActionMoveUp:
		a.cursorUp()
	case keymap.ActionMoveDown:
		a.cursorDown()
	case keymap.ActionMoveLeft:
		a.cursorLeft()
	case keymap.ActionMoveRight:
		a.cursorRight()
	case keymap.ActionMoveWordNext:
		a.wordNext()
	case keymap.ActionMoveWordPrev:
		a.wordPrev()
	case keymap.ActionGotoLineStart:
		a.col = 0
	case keymap.ActionGotoLineEnd:
		a.col = len(a.lines[a.row])
	case keymap.ActionGotoTop:
		a.row = 0
		a.col = 0
	case keymap.ActionGotoBottom:
		a.row = len(a.lines) - 1
		a.clampCol()
	}
}

// ════════════════════════════════════════════════════════════════
// UNDO / REDO
// ════════════════════════════════════════════════════════════════

func (a *App) saveUndo() {
	state := editorState{
		lines: make([]string, len(a.lines)),
		row:   a.row,
		col:   a.col,
	}
	copy(state.lines, a.lines)
	a.undoStack = append(a.undoStack, state)

	// Limit undo stack
	if len(a.undoStack) > 100 {
		a.undoStack = a.undoStack[1:]
	}

	// Clear redo stack
	a.redoStack = nil
}

func (a *App) undo() {
	if len(a.undoStack) == 0 {
		return
	}

	// Save current state to redo
	redoState := editorState{
		lines: make([]string, len(a.lines)),
		row:   a.row,
		col:   a.col,
	}
	copy(redoState.lines, a.lines)
	a.redoStack = append(a.redoStack, redoState)

	// Restore from undo
	state := a.undoStack[len(a.undoStack)-1]
	a.undoStack = a.undoStack[:len(a.undoStack)-1]

	a.lines = state.lines
	a.row = state.row
	a.col = state.col
}

func (a *App) redo() {
	if len(a.redoStack) == 0 {
		return
	}

	// Save current state to undo
	undoState := editorState{
		lines: make([]string, len(a.lines)),
		row:   a.row,
		col:   a.col,
	}
	copy(undoState.lines, a.lines)
	a.undoStack = append(a.undoStack, undoState)

	// Restore from redo
	state := a.redoStack[len(a.redoStack)-1]
	a.redoStack = a.redoStack[:len(a.redoStack)-1]

	a.lines = state.lines
	a.row = state.row
	a.col = state.col
}

// ════════════════════════════════════════════════════════════════
// VIEW
// ════════════════════════════════════════════════════════════════

func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	if a.showHelp {
		return a.renderHelp()
	}

	var b strings.Builder

	contentHeight := a.height - 2
	if contentHeight < 1 {
		contentHeight = 20
	}

	lineNumWidth := 5
	resultWidth := 20
	editorWidth := a.width - lineNumWidth - resultWidth - 4

	if editorWidth < 20 {
		editorWidth = 20
	}

	a.engine.Clear()

	for i := 0; i < contentHeight; i++ {
		if i < len(a.lines) {
			b.WriteString(lineNumStyle.Render(fmt.Sprintf("%3d ", i+1)))
		} else {
			b.WriteString(lineNumStyle.Render("    "))
		}

		b.WriteString("│")

		var editorContent string
		var resultContent string

		if i < len(a.lines) {
			line := a.lines[i]

			if i == a.row {
				editorContent = a.renderLineWithCursor(line)
			} else {
				editorContent = a.highlightLine(line)
			}

			resultContent = a.evaluateLine(line)
		} else {
			editorContent = tildeStyle.Render("~")
			resultContent = ""
		}

		editorLen := lipgloss.Width(editorContent)
		if editorLen < editorWidth {
			editorContent += strings.Repeat(" ", editorWidth-editorLen)
		} else if editorLen > editorWidth {
			editorContent = editorContent[:editorWidth]
		}

		resultContent = fmt.Sprintf("%*s", resultWidth, resultContent)

		b.WriteString(editorContent)
		b.WriteString("│")
		b.WriteString(resultContent)
		b.WriteString("\n")
	}

	b.WriteString(a.renderStatusBar())

	return b.String()
}

func (a *App) renderHelp() string {
	var content strings.Builder

	content.WriteString(helpTitleStyle.Render("Help"))
	content.WriteString("\n\n")

	content.WriteString(helpSectionStyle.Render("Navigation"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("[count]k/↑") + helpDescStyle.Render("Move up") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]j/↓") + helpDescStyle.Render("Move down") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]h/←") + helpDescStyle.Render("Move left") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]l/→") + helpDescStyle.Render("Move right") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]w") + helpDescStyle.Render("Next word") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]b") + helpDescStyle.Render("Previous word") + "\n")
	content.WriteString(helpKeyStyle.Render("0 / $") + helpDescStyle.Render("Start / End of line") + "\n")
	content.WriteString(helpKeyStyle.Render("gg / G") + helpDescStyle.Render("Top / Bottom of file") + "\n")

	content.WriteString(helpSectionStyle.Render("Editing"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("i / a") + helpDescStyle.Render("Insert / Append mode") + "\n")
	content.WriteString(helpKeyStyle.Render("o / O") + helpDescStyle.Render("Open line below/above") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]x") + helpDescStyle.Render("Delete character") + "\n")
	content.WriteString(helpKeyStyle.Render("[count]dd") + helpDescStyle.Render("Delete line") + "\n")
	content.WriteString(helpKeyStyle.Render("d{motion}") + helpDescStyle.Render("Delete with motion") + "\n")
	content.WriteString(helpKeyStyle.Render("yy / y{motion}") + helpDescStyle.Render("Yank line/motion") + "\n")
	content.WriteString(helpKeyStyle.Render("p / P") + helpDescStyle.Render("Paste after/before") + "\n")
	content.WriteString(helpKeyStyle.Render("u / Ctrl+r") + helpDescStyle.Render("Undo / Redo") + "\n")

	content.WriteString(helpSectionStyle.Render("General"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("Esc") + helpDescStyle.Render("Normal mode") + "\n")
	content.WriteString(helpKeyStyle.Render("?") + helpDescStyle.Render("Toggle help") + "\n")
	content.WriteString(helpKeyStyle.Render("q") + helpDescStyle.Render("Quit") + "\n")
	content.WriteString(helpKeyStyle.Render("Ctrl+C") + helpDescStyle.Render("Force quit") + "\n")

	content.WriteString(helpSectionStyle.Render("Examples"))
	content.WriteString("\n")
	content.WriteString(helpDescStyle.Render("5j      → Move down 5 lines") + "\n")
	content.WriteString(helpDescStyle.Render("3dd     → Delete 3 lines") + "\n")
	content.WriteString(helpDescStyle.Render("d3w     → Delete 3 words") + "\n")
	content.WriteString(helpDescStyle.Render("y$      → Yank to end of line") + "\n")

	content.WriteString(helpFooterStyle.Render("\nPress any key to close"))

	helpBox := helpBorderStyle.Render(content.String())

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, helpBox)
}

func (a *App) renderLineWithCursor(line string) string {
	col := a.col
	if col > len(line) {
		col = len(line)
	}

	if col == len(line) {
		return a.highlightLine(line) + cursorStyle.Render(" ")
	}

	before := line[:col]
	cursorChar := string(line[col])
	after := line[col+1:]

	return a.highlightLine(before) + cursorStyle.Render(cursorChar) + a.highlightLine(after)
}

func (a *App) highlightLine(line string) string {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return commentStyle.Render(line)
	}

	return line
}

func (a *App) evaluateLine(line string) string {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return ""
	}

	result := a.engine.Eval(line)

	if result.IsEmpty() {
		return ""
	}

	if result.IsError() {
		return errorStyle.Render("err")
	}

	return resultStyle.Render(result.String())
}

func (a *App) renderStatusBar() string {
	mode := a.keymap.GetMode()

	var modeStyle lipgloss.Style
	switch mode {
	case keymap.ModeInsert:
		modeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#7ee787")).Padding(0, 1)
	case keymap.ModeVisual:
		modeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#d2a8ff")).Padding(0, 1)
	case keymap.ModeOperatorPending:
		modeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#ffa657")).Padding(0, 1)
	default:
		modeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000")).Background(lipgloss.Color("#79c0ff")).Padding(0, 1)
	}

	modeStr := modeStyle.Render(mode.String())

	// Show pending keys
	pending := a.keymap.State.PendingDisplay()
	if pending != "" {
		modeStr += " " + pendingStyle.Render(pending)
	}

	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Render("  ? help  ^s save")

	pos := fmt.Sprintf("%d:%d", a.row+1, a.col+1)

	total := a.engine.Total()
	totalStr := ""
	if !total.IsEmpty() && total.AsFloat() != 0 {
		totalStr = resultStyle.Render(fmt.Sprintf("total: %s", total.String())) + "  "
	}

	left := modeStr + hint
	right := totalStr + pos

	spaces := a.width - lipgloss.Width(left) - lipgloss.Width(right)
	if spaces < 0 {
		spaces = 1
	}

	statusBg := lipgloss.NewStyle().Background(lipgloss.Color("#1a1a2e"))
	return statusBg.Render(left + strings.Repeat(" ", spaces) + right)
}

// ════════════════════════════════════════════════════════════════
// RUN
// ════════════════════════════════════════════════════════════════

// Run starts the TUI
func Run() error {
	p := tea.NewProgram(NewApp(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// RunWithFile starts with file content
func RunWithFile(filename, content string) error {
	app := NewApp()
	if content != "" {
		app.lines = strings.Split(content, "\n")
	}
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
