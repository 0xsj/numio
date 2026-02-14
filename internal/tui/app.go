// internal/tui/app.go

package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0xsj/numio/internal/fetch"
	"github.com/0xsj/numio/internal/graph"
	"github.com/0xsj/numio/internal/highlight"
	"github.com/0xsj/numio/internal/tui/keymap"
	"github.com/0xsj/numio/pkg/engine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	lineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666"))
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

	// Rate status styles
	rateStatusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#666"))
	rateFetchingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffa657"))
	rateSuccessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7ee787"))
	rateErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f85149"))

	// Explain popup styles
	explainBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#7ee787")).Padding(1, 2)
	explainTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7ee787"))
	explainInputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#79c0ff"))
	explainStepStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
	explainResultStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7ee787"))

	// History popup styles
	historyBorderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#ffa657")).Padding(1, 3)
	historyTitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffa657"))
	historyLabelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
	historyValueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0"))
	historyChangeUp      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7ee787"))
	historyChangeDown    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f85149"))
	historySparkStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffa657"))
	historyHintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Italic(true)
	historyForecastStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#d2a8ff"))
	historyForecastDim   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7e57c2"))
)

// App is the main model
type App struct {
	lines  []string
	row    int
	col    int
	width  int
	height int
	engine *engine.Engine

	// Syntax highlighting
	highlighter *highlight.Highlighter

	// Keymap
	keymap   *keymap.KeyMap
	showHelp bool

	// Explain mode
	showExplain bool
	lastExplain *engine.ExplainResult

	// Yank buffer
	yankBuffer string

	// Undo/Redo
	undoStack []editorState
	redoStack []editorState

	// Rate fetch status
	rateStatus   RateStatusInfo
	spinnerFrame int

	// History popup
	showHistory    bool
	historyResult  *fetch.HistoryResult
	historyRange   fetch.HistoryRange
	historyLoading bool
	historyError   string
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

	eng := engine.New()
	hl := highlight.Default()

	// Wire up user function registry for syntax highlighting
	hl.SetUserFuncRegistry(eng.UserFunctionRegistry())

	return &App{
		lines:        []string{""},
		row:          0,
		col:          0,
		width:        80,
		height:       24,
		engine:       eng,
		highlighter:  hl,
		keymap:       km,
		showHelp:     false,
		showExplain:  false,
		lastExplain:  nil,
		yankBuffer:   "",
		undoStack:    nil,
		redoStack:    nil,
		rateStatus:   RateStatusInfo{Status: RateStatusIdle},
		spinnerFrame: 0,
	}
}

// NewAppWithTheme creates a new app with a specific theme
func NewAppWithTheme(themeName string) *App {
	km, _ := keymap.LoadOrCreate(keymap.DefaultConfigPath())

	eng := engine.New()
	hl := highlight.NewWithThemeName(themeName)

	// Wire up user function registry for syntax highlighting
	hl.SetUserFuncRegistry(eng.UserFunctionRegistry())

	return &App{
		lines:        []string{""},
		row:          0,
		col:          0,
		width:        80,
		height:       24,
		engine:       eng,
		highlighter:  hl,
		keymap:       km,
		showHelp:     false,
		showExplain:  false,
		lastExplain:  nil,
		yankBuffer:   "",
		undoStack:    nil,
		redoStack:    nil,
		rateStatus:   RateStatusInfo{Status: RateStatusIdle},
		spinnerFrame: 0,
	}
}

// SetTheme changes the syntax highlighting theme
func (a *App) SetTheme(themeName string) {
	a.highlighter.SetTheme(highlight.GetTheme(themeName))
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	// Start fetching rates on startup
	a.rateStatus = RateStatusInfo{Status: RateStatusFetching}
	return tea.Batch(
		SpinnerTick(),
		StartRateFetch(a.fetchRates),
	)
}

// fetchRates fetches rates with a timeout context.
func (a *App) fetchRates() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return a.engine.RefreshRates(ctx)
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		return a.handleKey(msg)

	// Rate fetch messages
	case SpinnerTickMsg:
		if a.rateStatus.Status == RateStatusFetching || a.historyLoading {
			a.spinnerFrame++
			return a, SpinnerTick()
		}

	case RateFetchDoneMsg:
		if msg.Err != nil {
			a.rateStatus = RateStatusInfo{
				Status:    RateStatusError,
				Error:     msg.Err,
				UpdatedAt: time.Now(),
			}
		} else {
			a.rateStatus = RateStatusInfo{
				Status:    RateStatusSuccess,
				Message:   formatRateCount(msg.Count),
				UpdatedAt: time.Now(),
			}
		}
		// Clear success/error message after 3 seconds
		return a, ClearStatusAfter(3 * time.Second)

	case RateStatusClearMsg:
		// Only clear if we're in success/error state
		if a.rateStatus.Status == RateStatusSuccess || a.rateStatus.Status == RateStatusError {
			a.rateStatus = RateStatusInfo{
				Status:    RateStatusIdle,
				UpdatedAt: a.rateStatus.UpdatedAt,
			}
		}

	case HistoryFetchDoneMsg:
		a.historyLoading = false
		if msg.Err != nil {
			a.historyError = msg.Err.Error()
			a.historyResult = nil
		} else {
			a.historyResult = msg.Result
			a.historyError = ""
		}
		a.showHistory = true
	}

	return a, nil
}

// formatRateCount formats the rate count for display.
func formatRateCount(count int) string {
	if count == 0 {
		return "Rates up to date"
	}
	if count == 1 {
		return "1 rate updated"
	}
	return intToStr(count) + " rates updated"
}

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Always handle Ctrl+C as force quit
	if key == "ctrl+c" {
		return a, tea.Quit
	}

	// Close explain popup with any key
	if a.showExplain {
		a.showExplain = false
		return a, nil
	}

	// History popup: consume ALL keys while showing (never leak to editor)
	if a.showHistory {
		if !a.historyLoading {
			switch key {
			case "ctrl+p", "right", "l":
				a.historyRange = a.historyRange.Next()
				return a.triggerHistory()
			case "left", "h":
				a.historyRange = a.historyRange.Prev()
				return a.triggerHistory()
			case "1", "2", "3", "4":
				if r, ok := fetch.RangeFromKey(key); ok {
					a.historyRange = r
					return a.triggerHistory()
				}
				// invalid key — close
				a.showHistory = false
				a.historyResult = nil
				a.historyError = ""
			default:
				// Any other key closes
				a.showHistory = false
				a.historyResult = nil
				a.historyError = ""
			}
		}
		// Always consume the key — don't pass to editor
		return a, nil
	}

	// Handle Ctrl+E to show explanation
	if key == "ctrl+e" {
		a.triggerExplain()
		return a, nil
	}

	// Handle Ctrl+P to show price history
	if key == "ctrl+p" {
		return a.handleHistoryToggle()
	}

	// Handle Ctrl+A to select all / clear
	if key == "ctrl+a" {
		a.selectAll()
		return a, nil
	}

	// Handle Ctrl+R to refresh rates
	if key == "ctrl+r" {
		if a.rateStatus.Status != RateStatusFetching {
			a.rateStatus = RateStatusInfo{Status: RateStatusFetching}
			a.spinnerFrame = 0
			return a, tea.Batch(
				SpinnerTick(),
				StartRateFetch(a.fetchRates),
			)
		}
		return a, nil
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

// triggerExplain shows explanation for the current or previous line.
func (a *App) triggerExplain() {
	// Ensure row exists
	a.ensureRowExists()

	// Get current line
	line := strings.TrimSpace(a.lines[a.row])

	// If current line is empty or "explain", use previous non-empty line
	if line == "" || strings.ToLower(line) == "explain" {
		for i := a.row - 1; i >= 0; i-- {
			prevLine := strings.TrimSpace(a.lines[i])
			if prevLine != "" && !strings.HasPrefix(prevLine, "#") && !strings.HasPrefix(prevLine, "//") {
				line = prevLine
				break
			}
		}
	}

	if line == "" {
		return
	}

	// Get explanation
	a.lastExplain = a.engine.Explain(line)
	a.showExplain = true
}

// handleHistoryToggle opens the history popup or cycles range if already open.
func (a *App) handleHistoryToggle() (tea.Model, tea.Cmd) {
	if a.showHistory && !a.historyLoading {
		// Cycle to next range
		a.historyRange = a.historyRange.Next()
		return a.triggerHistory()
	}

	// Fresh open — detect asset from current/previous line
	a.historyRange = fetch.HistoryRange7d
	return a.triggerHistory()
}

// triggerHistory starts fetching price history for the current line's asset.
func (a *App) triggerHistory() (tea.Model, tea.Cmd) {
	a.ensureRowExists()

	// Find a non-empty line to evaluate
	line := strings.TrimSpace(a.lines[a.row])
	if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
		for i := a.row - 1; i >= 0; i-- {
			candidate := strings.TrimSpace(a.lines[i])
			if candidate != "" && !strings.HasPrefix(candidate, "#") && !strings.HasPrefix(candidate, "//") {
				line = candidate
				break
			}
		}
	}

	if line == "" {
		a.showHistory = true
		a.historyError = "No chartable asset on this line"
		a.historyResult = nil
		return a, nil
	}

	// Evaluate the line to detect the asset (with text scanning fallback)
	result := a.engine.Eval(line)
	code, kind, ok := fetch.DetectAssetFromLine(result, line)
	if !ok {
		a.showHistory = true
		a.historyError = "No chartable asset on this line"
		a.historyResult = nil
		return a, nil
	}

	// Edge cases
	if kind == fetch.AssetKindMetal {
		a.showHistory = true
		a.historyError = "Historical charts not yet available for metals"
		a.historyResult = nil
		return a, nil
	}
	if kind == fetch.AssetKindFiat && strings.ToUpper(code) == "USD" {
		a.showHistory = true
		a.historyError = "USD is the base currency"
		a.historyResult = nil
		return a, nil
	}

	a.historyLoading = true
	a.historyError = ""
	a.showHistory = true

	hr := a.historyRange
	return a, tea.Batch(
		SpinnerTick(),
		StartHistoryFetch(func() (*fetch.HistoryResult, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			return fetch.FetchHistory(ctx, code, kind, hr)
		}),
	)
}

func (a *App) handleInsertKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Check for bound keys in insert mode
	result := a.keymap.Insert.Lookup(key)
	if result.Status == keymap.LookupFound {
		cmd := keymap.NewCommand(result.Action, 1)
		return a.executeCommand(cmd)
	}

	// Handle Enter key - check for "explain" command
	if key == "enter" {
		a.ensureRowExists()
		line := strings.TrimSpace(a.lines[a.row])
		if strings.ToLower(line) == "explain" {
			// Find the previous non-empty line to explain
			for i := a.row - 1; i >= 0; i-- {
				prevLine := strings.TrimSpace(a.lines[i])
				if prevLine != "" && !strings.HasPrefix(prevLine, "#") && !strings.HasPrefix(prevLine, "//") && strings.ToLower(prevLine) != "explain" {
					a.lastExplain = a.engine.Explain(prevLine)
					a.showExplain = true
					return a, nil
				}
			}
		}
		// Normal newline
		a.saveUndo()
		a.newLine()
		return a, nil
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
		a.ensureRowExists()
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
		a.ensureRowExists()
		a.col = len(a.lines[a.row])

	case keymap.ActionGotoTop:
		a.row = 0
		a.col = 0

	case keymap.ActionGotoBottom:
		// Find last non-empty line
		lastContent := 0
		for i := len(a.lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(a.lines[i]) != "" {
				lastContent = i
				break
			}
		}
		a.row = lastContent
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

	case keymap.ActionSelectAll:
		a.selectAll()

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

// ensureRowExists makes sure the current row exists in the lines slice.
func (a *App) ensureRowExists() {
	for a.row >= len(a.lines) {
		a.lines = append(a.lines, "")
	}
}

func (a *App) cursorUp() {
	if a.row > 0 {
		a.row--
		a.clampCol()
	}
}

func (a *App) cursorDown() {
	a.row++
	// Expand lines slice if needed
	a.ensureRowExists()
	a.clampCol()
}

func (a *App) cursorLeft() {
	if a.col > 0 {
		a.col--
	}
}

func (a *App) cursorRight() {
	a.ensureRowExists()

	maxCol := len(a.lines[a.row])
	if a.keymap.CurrentMode == keymap.ModeNormal && maxCol > 0 {
		maxCol--
	}
	if a.col < maxCol {
		a.col++
	}
}

func (a *App) clampCol() {
	a.ensureRowExists()

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
	a.ensureRowExists()

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

	if col >= len(line) {
		// Move to next line
		a.row++
		a.ensureRowExists()
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

	a.ensureRowExists()

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
	a.ensureRowExists()

	line := a.lines[a.row]
	if a.col > len(line) {
		a.col = len(line)
	}
	a.lines[a.row] = line[:a.col] + string(r) + line[a.col:]
	a.col++
}

func (a *App) newLine() {
	a.ensureRowExists()

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
	a.ensureRowExists()

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
	a.ensureRowExists()

	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row]...)
	newLines = append(newLines, "")
	newLines = append(newLines, a.lines[a.row:]...)
	a.lines = newLines
	a.col = 0
}

func (a *App) backspace() {
	a.ensureRowExists()

	if a.col > 0 {
		line := a.lines[a.row]
		if a.col <= len(line) {
			a.lines[a.row] = line[:a.col-1] + line[a.col:]
		}
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
	a.ensureRowExists()

	line := a.lines[a.row]
	if a.col < len(line) {
		a.lines[a.row] = line[:a.col] + line[a.col+1:]
	} else if a.row < len(a.lines)-1 {
		a.lines[a.row] = line + a.lines[a.row+1]
		a.lines = append(a.lines[:a.row+1], a.lines[a.row+2:]...)
	}
}

func (a *App) deleteLine() {
	a.ensureRowExists()

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
	a.ensureRowExists()

	line := a.lines[a.row]
	if a.col < len(line) {
		a.yankBuffer = line[a.col:]
		a.lines[a.row] = line[:a.col]
	}
	a.clampCol()
}

func (a *App) joinLines() {
	a.ensureRowExists()

	if a.row < len(a.lines)-1 {
		a.lines[a.row] = a.lines[a.row] + " " + strings.TrimLeft(a.lines[a.row+1], " \t")
		a.lines = append(a.lines[:a.row+1], a.lines[a.row+2:]...)
	}
}

func (a *App) selectAll() {
	a.saveUndo()
	a.lines = []string{""}
	a.row = 0
	a.col = 0
	a.keymap.SetMode(keymap.ModeInsert)
}

func (a *App) yankLine() {
	a.ensureRowExists()
	a.yankBuffer = a.lines[a.row] + "\n"
}

func (a *App) paste() {
	if a.yankBuffer == "" {
		return
	}

	a.ensureRowExists()

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
		insertPos := a.col + 1
		if insertPos > len(line) {
			insertPos = len(line)
		}
		a.lines[a.row] = line[:insertPos] + a.yankBuffer + line[insertPos:]
		a.col = insertPos + len(a.yankBuffer) - 1
	}
}

func (a *App) pasteAbove() {
	if a.yankBuffer == "" {
		return
	}

	a.ensureRowExists()

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
	a.ensureRowExists()

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
			yanked := a.lines[endRow]
			if endCol < len(yanked) {
				newLine += yanked[endCol:]
			}
		}
		a.lines[startRow] = newLine

		// Remove middle lines
		if endRow < len(a.lines) {
			a.lines = append(a.lines[:startRow+1], a.lines[endRow+1:]...)
		} else {
			a.lines = a.lines[:startRow+1]
		}

		a.row = startRow
		a.col = startCol
	}

	a.clampCol()
}

func (a *App) yankWithMotion(motion keymap.Action, count int) {
	a.ensureRowExists()

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
			line := a.lines[endRow]
			if endCol <= len(line) {
				yanked.WriteString(line[:endCol])
			} else {
				yanked.WriteString(line)
			}
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
		a.ensureRowExists()
		a.col = len(a.lines[a.row])
	case keymap.ActionGotoTop:
		a.row = 0
		a.col = 0
	case keymap.ActionGotoBottom:
		lastContent := 0
		for i := len(a.lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(a.lines[i]) != "" {
				lastContent = i
				break
			}
		}
		a.row = lastContent
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

	if a.showExplain {
		return a.renderExplain()
	}

	if a.showHistory {
		return a.renderHistory()
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

	// Calculate scroll offset to keep cursor visible
	scrollOffset := 0
	if a.row >= contentHeight {
		scrollOffset = a.row - contentHeight + 1
	}

	for i := 0; i < contentHeight; i++ {
		lineIdx := i + scrollOffset

		if lineIdx < len(a.lines) {
			b.WriteString(lineNumStyle.Render(fmt.Sprintf("%3d ", lineIdx+1)))
		} else {
			b.WriteString(lineNumStyle.Render("    "))
		}

		b.WriteString("│")

		var editorContent string
		var resultContent string

		if lineIdx < len(a.lines) {
			line := a.lines[lineIdx]

			if lineIdx == a.row {
				editorContent = a.renderLineWithCursor(line)
			} else {
				editorContent = a.highlighter.Highlight(line)
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

func (a *App) renderExplain() string {
	if a.lastExplain == nil {
		return a.View()
	}

	var content strings.Builder

	content.WriteString(explainTitleStyle.Render("Explanation"))
	content.WriteString("\n\n")

	// Input expression
	content.WriteString(explainInputStyle.Render("Input: "))
	if a.lastExplain.Trace != nil {
		content.WriteString(a.lastExplain.Trace.Input)
	}
	content.WriteString("\n\n")

	// Steps
	if len(a.lastExplain.Steps) > 0 {
		content.WriteString(explainStepStyle.Render("Steps:"))
		content.WriteString("\n")
		for _, step := range a.lastExplain.Steps {
			content.WriteString(explainStepStyle.Render("  " + step))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// One-line explanation
	if a.lastExplain.Explanation != "" {
		content.WriteString(explainStepStyle.Render("Calculation: "))
		content.WriteString(a.lastExplain.Explanation)
		content.WriteString("\n\n")
	}

	// Result
	content.WriteString(explainResultStyle.Render("Result: "))
	content.WriteString(explainResultStyle.Render(a.lastExplain.Value.String()))
	content.WriteString("\n")

	content.WriteString(helpFooterStyle.Render("\nPress any key to close"))

	explainBox := explainBorderStyle.Render(content.String())

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, explainBox)
}

func (a *App) renderHistory() string {
	var content strings.Builder

	if a.historyLoading {
		frame := SpinnerFrames[a.spinnerFrame%len(SpinnerFrames)]
		content.WriteString(historyTitleStyle.Render("Price Chart"))
		content.WriteString("\n\n")
		content.WriteString(frame + " Fetching price history...")
		content.WriteString("\n")

		historyBox := historyBorderStyle.Render(content.String())
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, historyBox)
	}

	if a.historyError != "" {
		content.WriteString(historyTitleStyle.Render("Price Chart"))
		content.WriteString("\n\n")
		content.WriteString(errorStyle.Render(a.historyError))
		content.WriteString("\n")
		content.WriteString(historyHintStyle.Render("\nPress any key to close"))

		historyBox := historyBorderStyle.Render(content.String())
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, historyBox)
	}

	if a.historyResult == nil {
		return a.View()
	}

	hr := a.historyResult
	prices := hr.Prices()
	stats := graph.ComputeStats(prices)

	// Title
	title := hr.Asset + "/" + hr.Base + " Price Chart"
	content.WriteString(historyTitleStyle.Render(title))
	content.WriteString("\n")

	// Range selector: [1:7d] [2:30d] [3:90d] [4:1y]
	rangeLabels := []struct {
		key   string
		label string
		r     fetch.HistoryRange
	}{
		{"1", "7d", fetch.HistoryRange7d},
		{"2", "30d", fetch.HistoryRange30d},
		{"3", "90d", fetch.HistoryRange90d},
		{"4", "1y", fetch.HistoryRange1y},
	}
	content.WriteString("  ")
	for _, rl := range rangeLabels {
		if rl.r == hr.Range {
			content.WriteString(historyTitleStyle.Render("[" + rl.key + ":" + rl.label + "]"))
		} else {
			content.WriteString(historyHintStyle.Render(" " + rl.key + ":" + rl.label + " "))
		}
	}
	content.WriteString("\n\n")

	// Sparkline with forecast extension
	forecast := fetch.LinearForecast(prices)
	sparkWidth := 44
	forecastWidth := 0
	if forecast != nil {
		forecastWidth = len(forecast.Points)
		if forecastWidth > 8 {
			forecastWidth = 8
		}
	}

	sparkline := graph.SparklineFixed(prices, sparkWidth)
	content.WriteString("  " + historySparkStyle.Render(sparkline))
	if forecast != nil && forecastWidth > 0 {
		// Render forecast sparkline with different color using combined bounds
		allPrices := append(prices, forecast.Points[:forecastWidth]...)
		fStats := graph.ComputeStats(allPrices)
		forecastSpark := graph.SparklineBounded(forecast.Points[:forecastWidth], fStats.Min, fStats.Max)
		content.WriteString(historyForecastDim.Render("╎"))
		content.WriteString(historyForecastStyle.Render(forecastSpark))
	}
	content.WriteString(" " + stats.Trend.Symbol())
	content.WriteString("\n")
	if forecast != nil {
		content.WriteString(historyHintStyle.Render("  " + padRight("historical", sparkWidth) + " forecast"))
		content.WriteString("\n")
	}
	content.WriteString("\n")

	// Stats
	open := formatPrice(stats.First)
	close_ := formatPrice(stats.Last)
	high := formatPrice(stats.Max)
	low := formatPrice(stats.Min)
	avg := formatPrice(stats.Avg)

	// Change with color
	changeStr := formatChange(stats.Change)
	var changeStyled string
	if stats.Change >= 0 {
		changeStyled = historyChangeUp.Render(changeStr)
	} else {
		changeStyled = historyChangeDown.Render(changeStr)
	}

	content.WriteString(historyLabelStyle.Render("  Open:   ") + historyValueStyle.Render(padRight(open, 12)) + historyLabelStyle.Render("High:  ") + historyValueStyle.Render(high))
	content.WriteString("\n")
	content.WriteString(historyLabelStyle.Render("  Close:  ") + historyValueStyle.Render(padRight(close_, 12)) + historyLabelStyle.Render("Low:   ") + historyValueStyle.Render(low))
	content.WriteString("\n")
	content.WriteString(historyLabelStyle.Render("  Change: ") + padRight(changeStyled, 12+10) + historyLabelStyle.Render("Avg:   ") + historyValueStyle.Render(avg))
	content.WriteString("\n")

	// Forecast section
	if forecast != nil {
		content.WriteString("\n")
		content.WriteString(historyForecastStyle.Render("  Forecast (" + forecast.Label + ")"))
		content.WriteString("\n")
		targetStr := formatPrice(forecast.Target)
		var fcChangeStyled string
		fcChangeStr := formatChange(forecast.ChangePct)
		if forecast.ChangePct >= 0 {
			fcChangeStyled = historyChangeUp.Render(fcChangeStr)
		} else {
			fcChangeStyled = historyChangeDown.Render(fcChangeStr)
		}
		content.WriteString(historyLabelStyle.Render("  Target: ") + historyValueStyle.Render(padRight(targetStr, 12)))
		content.WriteString(fcChangeStyled)
		content.WriteString(historyLabelStyle.Render("  R²: "))
		content.WriteString(historyValueStyle.Render(formatR2(forecast.RSquared)))
		content.WriteString(historyLabelStyle.Render(" (" + forecast.Confidence() + ")"))
		content.WriteString("\n")
	}

	content.WriteString(historyHintStyle.Render("\n  1-4 select · ←/→ cycle · any key to close"))

	historyBox := historyBorderStyle.Render(content.String())
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, historyBox)
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

	content.WriteString(helpSectionStyle.Render("Functions"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("def f(x):") + helpDescStyle.Render("Define function") + "\n")
	content.WriteString(helpDescStyle.Render("  Example: def double(x): x * 2") + "\n")

	content.WriteString(helpSectionStyle.Render("General"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("Esc") + helpDescStyle.Render("Normal mode") + "\n")
	content.WriteString(helpKeyStyle.Render("?") + helpDescStyle.Render("Toggle help") + "\n")
	content.WriteString(helpKeyStyle.Render("Ctrl+e") + helpDescStyle.Render("Explain calculation") + "\n")
	content.WriteString(helpKeyStyle.Render("Ctrl+p") + helpDescStyle.Render("Price chart") + "\n")
	content.WriteString(helpKeyStyle.Render("Ctrl+r") + helpDescStyle.Render("Refresh rates") + "\n")
	content.WriteString(helpKeyStyle.Render("q") + helpDescStyle.Render("Quit") + "\n")
	content.WriteString(helpKeyStyle.Render("Ctrl+C") + helpDescStyle.Render("Force quit") + "\n")

	content.WriteString(helpSectionStyle.Render("Examples"))
	content.WriteString("\n")
	content.WriteString(helpDescStyle.Render("5j      → Move down 5 lines") + "\n")
	content.WriteString(helpDescStyle.Render("3dd     → Delete 3 lines") + "\n")
	content.WriteString(helpDescStyle.Render("d3w     → Delete 3 words") + "\n")
	content.WriteString(helpDescStyle.Render("y$      → Yank to end of line") + "\n")
	content.WriteString(helpDescStyle.Render("explain → Show step-by-step") + "\n")

	content.WriteString(helpFooterStyle.Render("\nPress any key to close"))

	helpBox := helpBorderStyle.Render(content.String())

	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, helpBox)
}

func (a *App) renderLineWithCursor(line string) string {
	col := a.col
	if col > len(line) {
		col = len(line)
	}

	// Cursor at end of line
	if col == len(line) {
		return a.highlighter.Highlight(line) + cursorStyle.Render(" ")
	}

	// Get highlighted spans for precise cursor placement
	spans := a.highlighter.HighlightSpans(line)

	var result strings.Builder
	cursorRendered := false

	for _, span := range spans {
		// Check if cursor is within this span
		if !cursorRendered && col >= span.Start && col < span.End {
			// Cursor is in this span - split it
			relativeCol := col - span.Start
			beforeCursor := span.Text[:relativeCol]
			cursorChar := string(span.Text[relativeCol])
			afterCursor := span.Text[relativeCol+1:]

			// Render parts with the span's color
			style := a.highlighter.Theme().Style(span.Class)
			if beforeCursor != "" {
				result.WriteString(style.Render(beforeCursor))
			}
			result.WriteString(cursorStyle.Render(cursorChar))
			if afterCursor != "" {
				result.WriteString(style.Render(afterCursor))
			}
			cursorRendered = true
		} else {
			// Render entire span normally
			result.WriteString(a.highlighter.Theme().Render(span.Class, span.Text))
		}
	}

	return result.String()
}

func (a *App) evaluateLine(line string) string {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return ""
	}

	// Don't evaluate "explain" command
	if strings.ToLower(trimmed) == "explain" {
		return ""
	}

	result := a.engine.Eval(line)

	if result.IsEmpty() {
		return ""
	}

	if result.IsError() {
		errMsg := result.String()
		if len(errMsg) > 25 {
			errMsg = errMsg[:22] + "..."
		}
		return errorStyle.Render(errMsg)
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

	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Render("  ? help  ^e explain  ^p chart  ^r rates")

	pos := fmt.Sprintf("%d:%d", a.row+1, a.col+1)

	// Rate status display
	rateStatusStr := a.renderRateStatus()

	total := a.engine.Total()
	totalStr := ""
	if !total.IsEmpty() && total.AsFloat() != 0 {
		totalStr = resultStyle.Render(fmt.Sprintf("Σ %s", total.String())) + "  "
	}

	left := modeStr + hint
	right := rateStatusStr + totalStr + pos

	spaces := a.width - lipgloss.Width(left) - lipgloss.Width(right)
	if spaces < 0 {
		spaces = 1
	}

	statusBg := lipgloss.NewStyle().Background(lipgloss.Color("#1a1a2e"))
	return statusBg.Render(left + strings.Repeat(" ", spaces) + right)
}

// renderRateStatus renders the rate status with appropriate styling.
func (a *App) renderRateStatus() string {
	statusText := FormatRateStatus(a.rateStatus, a.spinnerFrame)
	if statusText == "" {
		return ""
	}

	var style lipgloss.Style
	switch a.rateStatus.Status {
	case RateStatusFetching:
		style = rateFetchingStyle
	case RateStatusSuccess:
		style = rateSuccessStyle
	case RateStatusError:
		style = rateErrorStyle
	default:
		style = rateStatusStyle
	}

	return style.Render(statusText) + "  "
}

// ════════════════════════════════════════════════════════════════
// HISTORY HELPERS
// ════════════════════════════════════════════════════════════════

// formatPrice formats a price for display in the history popup.
func formatPrice(n float64) string {
	if n == 0 {
		return "0"
	}

	// Handle negative
	prefix := ""
	if n < 0 {
		prefix = "-"
		n = -n
	}

	// Determine decimals based on magnitude
	var decimals int
	if n >= 1000 {
		decimals = 0
	} else if n >= 1 {
		decimals = 2
	} else if n >= 0.01 {
		decimals = 4
	} else {
		decimals = 6
	}

	// Format with comma separators for large numbers
	intPart := int64(n)
	fracPart := n - float64(intPart)

	intStr := formatIntWithCommas(intPart)

	if decimals == 0 {
		return prefix + intStr
	}

	// Build decimal part
	mul := 1.0
	for i := 0; i < decimals; i++ {
		mul *= 10
	}
	frac := int64(fracPart * mul)
	fracStr := intToStr(int(frac))
	for len(fracStr) < decimals {
		fracStr = "0" + fracStr
	}

	// Trim trailing zeros
	fracStr = strings.TrimRight(fracStr, "0")
	if fracStr == "" {
		return prefix + intStr
	}

	return prefix + intStr + "." + fracStr
}

// formatIntWithCommas formats an integer with comma separators.
func formatIntWithCommas(n int64) string {
	s := intToStr(int(n))
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	start := len(s) % 3
	if start > 0 {
		result.WriteString(s[:start])
	}
	for i := start; i < len(s); i += 3 {
		if result.Len() > 0 {
			result.WriteByte(',')
		}
		result.WriteString(s[i : i+3])
	}
	return result.String()
}

// formatChange formats a percentage change like "+12.0%" or "-5.3%".
func formatChange(pct float64) string {
	prefix := "+"
	neg := pct < 0
	if neg {
		prefix = "-"
		pct = -pct
	}

	intPart := int(pct)
	fracPart := int((pct - float64(intPart)) * 10)

	return prefix + intToStr(intPart) + "." + intToStr(fracPart) + "%"
}

// formatR2 formats the R² coefficient (0..1) with 2 decimal places.
func formatR2(r2 float64) string {
	intPart := int(r2)
	fracPart := int((r2 - float64(intPart)) * 100)
	frac := intToStr(fracPart)
	if len(frac) < 2 {
		frac = "0" + frac
	}
	return intToStr(intPart) + "." + frac
}

// padRight pads a string to a minimum width with spaces.
func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
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

// RunWithTheme starts the TUI with a specific theme
func RunWithTheme(themeName string) error {
	app := NewAppWithTheme(themeName)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
