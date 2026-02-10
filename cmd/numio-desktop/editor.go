package main

import (
	"context"
	"fmt"
	"image/color"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/0xsj/numio/internal/editor"
	"github.com/0xsj/numio/internal/rpc"
	"github.com/0xsj/numio/internal/tui/keymap"
)

// ════════════════════════════════════════════════════════════════
// DESKTOP-SPECIFIC ACTIONS
// ════════════════════════════════════════════════════════════════

// Actions not in the shared keymap defaults but needed for desktop vim.
const (
	actionInsertLineStart keymap.Action = "insert_line_start"
	actionInsertLineEnd   keymap.Action = "insert_line_end"
	actionVisualLineMode  keymap.Action = "visual_line_mode"
)

// ════════════════════════════════════════════════════════════════
// POPUP TYPES
// ════════════════════════════════════════════════════════════════

type PopupType int

const (
	PopupNone PopupType = iota
	PopupHelp
	PopupExplain
)

// ════════════════════════════════════════════════════════════════
// EDITOR WIDGET
// ════════════════════════════════════════════════════════════════

// EditorWidget is a custom Fyne widget that wraps the numio editor.
type EditorWidget struct {
	widget.BaseWidget

	editor *editor.Editor
	state  *rpc.RenderState
	mu     sync.RWMutex

	// Layout constants
	lineHeight      float32
	charWidth       float32
	leftPadding     float32
	rightPadding    float32
	topPadding      float32
	statusBarHeight float32
	lineSpacing     float32
	lineNumWidth    float32

	// Cursor blink state
	cursorVisible bool

	// Modifier key state
	ctrlPressed  bool
	altPressed   bool
	shiftPressed bool
	superPressed bool

	// Popup state
	activePopup   PopupType
	explainResult string

	// Vim mode toggle (default: true)
	vimMode bool

	// Keymap for vim bindings
	km *keymap.KeyMap
}

// NewEditorWidget creates a new editor widget.
func NewEditorWidget() *EditorWidget {
	km := keymap.Default()

	// Add desktop-specific bindings not in shared defaults
	km.Normal.Bind("I", actionInsertLineStart)
	km.Normal.Bind("A", actionInsertLineEnd)
	km.Normal.Bind("V", actionVisualLineMode)

	w := &EditorWidget{
		editor:          editor.NewEditor(),
		lineHeight:      24,
		charWidth:       8.4,
		leftPadding:     20,
		rightPadding:    20,
		topPadding:      16,
		statusBarHeight: 28,
		lineSpacing:     4,
		lineNumWidth:    8.4 * 4, // 4 characters wide
		cursorVisible:   true,
		activePopup:     PopupNone,
		vimMode:         true, // Default: vim mode ON
		km:              km,
	}

	// Editor starts in ModeNormal by default — synced with keymap

	w.ExtendBaseWidget(w)
	w.updateState()

	// Fetch rates in background on startup
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		w.editor.Engine().RefreshRates(ctx)
	}()

	return w
}

// ════════════════════════════════════════════════════════════════
// VIM MODE TOGGLE
// ════════════════════════════════════════════════════════════════

// SetVimMode enables or disables vim keybindings.
func (w *EditorWidget) SetVimMode(enabled bool) {
	w.vimMode = enabled
	if enabled {
		w.editor.EnterNormalMode()
		w.km.SetMode(keymap.ModeNormal)
	} else {
		w.editor.EnterInsertMode()
		w.km.SetMode(keymap.ModeInsert)
	}
	w.updateState()
	w.Refresh()
}

// ToggleVimMode toggles vim keybindings on/off.
func (w *EditorWidget) ToggleVimMode() {
	w.SetVimMode(!w.vimMode)
}

// VimMode returns whether vim mode is enabled.
func (w *EditorWidget) VimMode() bool {
	return w.vimMode
}

// ════════════════════════════════════════════════════════════════
// WIDGET INTERFACE
// ════════════════════════════════════════════════════════════════

// CreateRenderer implements fyne.Widget.
func (w *EditorWidget) CreateRenderer() fyne.WidgetRenderer {
	return newEditorRenderer(w)
}

// MinSize returns the minimum size of the widget.
func (w *EditorWidget) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

// Resize handles widget resize.
func (w *EditorWidget) Resize(size fyne.Size) {
	w.BaseWidget.Resize(size)

	// Calculate viewport size in characters (subtract line number gutter in vim mode)
	gutterWidth := float32(0)
	if w.vimMode {
		gutterWidth = w.lineNumWidth
	}
	cols := int((size.Width - w.leftPadding - gutterWidth - w.rightPadding) / w.charWidth)
	rows := int((size.Height - w.topPadding*2 - w.statusBarHeight) / (w.lineHeight + w.lineSpacing))

	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	w.mu.Lock()
	w.editor.Resize(cols, rows)
	w.mu.Unlock()

	w.updateState()
}

// ════════════════════════════════════════════════════════════════
// FOCUSABLE INTERFACE
// ════════════════════════════════════════════════════════════════

// FocusGained is called when the widget gains focus.
func (w *EditorWidget) FocusGained() {
	w.cursorVisible = true
	w.Refresh()
}

// FocusLost is called when the widget loses focus.
func (w *EditorWidget) FocusLost() {
	w.Refresh()
}

// TypedRune handles character input.
func (w *EditorWidget) TypedRune(r rune) {
	// Close popup on any key
	if w.activePopup != PopupNone {
		w.activePopup = PopupNone
		w.Refresh()
		return
	}

	w.mu.Lock()

	if w.vimMode {
		if w.editor.Mode() == editor.ModeInsert {
			// Insert mode: type characters directly
			w.editor.InsertChar(r)
		} else {
			// Normal/Visual: feed to keymap
			w.processVimKeymapKey(string(r))
		}
	} else {
		// Simple editor mode: always insert
		w.editor.InsertChar(r)
	}

	w.mu.Unlock()
	w.updateState()
	w.Refresh()
}

// TypedKey handles special key input.
func (w *EditorWidget) TypedKey(ev *fyne.KeyEvent) {
	// F-key fallbacks (always available regardless of modifier state)
	switch ev.Name {
	case fyne.KeyF1:
		w.ToggleHelp()
		return
	case fyne.KeyF2:
		w.ShowExplain()
		return
	case fyne.KeyF3:
		w.ToggleVimMode()
		return
	}

	// Close popup on Escape or any key
	if w.activePopup != PopupNone {
		w.activePopup = PopupNone
		w.Refresh()
		return
	}

	w.mu.Lock()

	if w.vimMode {
		// Build keymap-compatible key string
		keymapKey := translateKeyForKeymap(ev.Name)

		// Handle Ctrl + letter (Fyne sends TypedKey, not TypedRune, when Ctrl held)
		if keymapKey == "" && w.ctrlPressed {
			letter := strings.ToLower(string(ev.Name))
			if len(letter) == 1 && letter[0] >= 'a' && letter[0] <= 'z' {
				keymapKey = "ctrl+" + letter
			}
		}

		if keymapKey == "" {
			// Regular character — handled by TypedRune
			w.mu.Unlock()
			return
		}

		// Add ctrl prefix for special keys when ctrl is held
		if w.ctrlPressed && !strings.HasPrefix(keymapKey, "ctrl+") {
			keymapKey = "ctrl+" + keymapKey
		}

		w.processVimKeymapKey(keymapKey)
	} else {
		key := translateKey(ev.Name)

		// Skip regular characters (handled by TypedRune)
		if len(key) == 1 && key[0] >= 32 && key[0] < 127 {
			w.mu.Unlock()
			return
		}

		w.processNormalEditorKey(key)
	}

	w.mu.Unlock()
	w.updateState()
	w.Refresh()
}

// ════════════════════════════════════════════════════════════════
// DESKTOP INTERFACE (for modifier keys)
// ════════════════════════════════════════════════════════════════

var _ desktop.Keyable = (*EditorWidget)(nil)

// KeyDown handles key press with modifiers.
func (w *EditorWidget) KeyDown(ev *fyne.KeyEvent) {
	switch ev.Name {
	case desktop.KeyControlLeft, desktop.KeyControlRight:
		w.ctrlPressed = true
	case desktop.KeyAltLeft, desktop.KeyAltRight:
		w.altPressed = true
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
		w.shiftPressed = true
	case desktop.KeySuperLeft, desktop.KeySuperRight:
		w.superPressed = true
	}
}

// KeyUp handles key release.
func (w *EditorWidget) KeyUp(ev *fyne.KeyEvent) {
	switch ev.Name {
	case desktop.KeyControlLeft, desktop.KeyControlRight:
		w.ctrlPressed = false
	case desktop.KeyAltLeft, desktop.KeyAltRight:
		w.altPressed = false
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
		w.shiftPressed = false
	case desktop.KeySuperLeft, desktop.KeySuperRight:
		w.superPressed = false
	}
}

// ════════════════════════════════════════════════════════════════
// HELP TOGGLE
// ════════════════════════════════════════════════════════════════

// ToggleHelp toggles the help popup on/off.
func (w *EditorWidget) ToggleHelp() {
	if w.activePopup == PopupHelp {
		w.activePopup = PopupNone
	} else {
		w.activePopup = PopupHelp
	}
	w.Refresh()
}

// ════════════════════════════════════════════════════════════════
// SIMPLE EDITOR KEY HANDLING (non-vim)
// ════════════════════════════════════════════════════════════════

func (w *EditorWidget) processNormalEditorKey(key string) {
	switch key {
	case "Backspace":
		w.editor.DeleteCharBack()
	case "Delete":
		w.editor.DeleteChar()
	case "Enter":
		w.editor.InsertNewline()
	case "Tab":
		w.editor.InsertString("  ")
	case "ArrowLeft":
		w.editor.MoveLeft(1)
	case "ArrowRight":
		w.editor.MoveRight(1)
	case "ArrowUp":
		w.editor.MoveUp(1)
	case "ArrowDown":
		w.editor.MoveDown(1)
	case "Home":
		w.editor.MoveToLineStart()
	case "End":
		w.editor.MoveToLineEnd()
	}
}

// ════════════════════════════════════════════════════════════════
// VIM KEYMAP PROCESSING
// ════════════════════════════════════════════════════════════════

// processVimKeymapKey feeds a key through the keymap and dispatches the result.
// This replaces the old hand-rolled processVimKey / processVimNormalModeKey / etc.
func (w *EditorWidget) processVimKeymapKey(key string) {
	km := w.km
	key = keymap.NormalizeKey(key)

	// Handle digit for count (but not '0' at start which is line start)
	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		if km.State.AddDigit(rune(key[0])) {
			return
		}
	}

	// Add key to buffer
	km.State.AddKey(key)

	// Determine which binding map to use
	mode := km.CurrentMode
	if km.State.HasPendingOperator() {
		mode = keymap.ModeOperatorPending
	}

	result := km.GetBindingMap(mode).Lookup(km.State.KeyBuffer)

	switch result.Status {
	case keymap.LookupFound:
		w.handleKeymapAction(result.Action)

	case keymap.LookupPending:
		if result.Action != keymap.ActionNone {
			w.handleKeymapAction(result.Action)
		}
		// else: partial sequence, wait for more keys

	case keymap.LookupPartialMatch:
		// waiting for more keys

	case keymap.LookupNotFound:
		km.State.Reset()
	}
}

// handleKeymapAction processes a resolved action, managing operator-pending state.
func (w *EditorWidget) handleKeymapAction(action keymap.Action) {
	km := w.km
	count := km.State.GetCount()
	if count < 1 {
		count = 1
	}

	// In visual mode, operators act on the selection immediately
	if action.IsOperator() && km.CurrentMode == keymap.ModeVisual {
		cmd := keymap.NewCommand(action, count)
		km.State.Reset()
		w.executeCommand(cmd)
		return
	}

	// Starting a new operator (d, y, c) — enter operator-pending mode
	if action.IsOperator() && !km.State.HasPendingOperator() {
		km.State.SetOperator(action)
		km.State.ClearKeyBuffer()
		return // wait for motion
	}

	// Completing an operator with a motion
	if km.State.HasPendingOperator() {
		cmd := keymap.NewOperatorCommand(km.State.PendingOperator, count, action, 1)
		km.State.Reset()
		w.executeCommand(cmd)
		return
	}

	// Simple command
	cmd := keymap.NewCommand(action, count)
	km.State.Reset()
	w.executeCommand(cmd)
}

// ════════════════════════════════════════════════════════════════
// COMMAND DISPATCH
// ════════════════════════════════════════════════════════════════

// executeCommand dispatches a keymap Command to the editor.
func (w *EditorWidget) executeCommand(cmd keymap.Command) {
	if cmd.Action == keymap.ActionNone {
		return
	}

	ed := w.editor
	count := cmd.Count
	if count < 1 {
		count = 1
	}

	// Operator + motion commands
	if cmd.Motion != keymap.ActionNone {
		w.executeOperatorMotion(cmd)
		w.syncModeToKeymap()
		return
	}

	switch cmd.Action {
	// ── Mode switching ──────────────────────────────────────────
	case keymap.ActionInsertMode:
		ed.EnterInsertMode()
	case keymap.ActionAppendMode:
		ed.EnterInsertModeAppend()
	case actionInsertLineStart:
		ed.EnterInsertModeLineStart()
	case actionInsertLineEnd:
		ed.EnterInsertModeLineEnd()
	case keymap.ActionNormalMode:
		ed.EnterNormalMode()
	case keymap.ActionVisualMode:
		ed.EnterVisualMode()
	case actionVisualLineMode:
		ed.EnterVisualLineMode()

	// ── Movement ────────────────────────────────────────────────
	case keymap.ActionMoveUp:
		ed.MoveUp(count)
	case keymap.ActionMoveDown:
		ed.MoveDown(count)
	case keymap.ActionMoveLeft:
		ed.MoveLeft(count)
	case keymap.ActionMoveRight:
		ed.MoveRight(count)
	case keymap.ActionMoveWordNext:
		ed.MoveWordForward(count)
	case keymap.ActionMoveWordPrev:
		ed.MoveWordBackward(count)
	case keymap.ActionGotoLineStart:
		ed.MoveToLineStart()
	case keymap.ActionGotoLineEnd:
		ed.MoveToLineEnd()
	case keymap.ActionGotoTop:
		ed.MoveToTop()
	case keymap.ActionGotoBottom:
		ed.MoveToBottom()
	case keymap.ActionPageUp:
		ed.MoveUp(ed.Viewport().Height())
	case keymap.ActionPageDown:
		ed.MoveDown(ed.Viewport().Height())

	// ── Editing ─────────────────────────────────────────────────
	case keymap.ActionDeleteChar:
		for i := 0; i < count; i++ {
			ed.DeleteChar()
		}
	case keymap.ActionDeleteCharBack:
		for i := 0; i < count; i++ {
			ed.DeleteCharBack()
		}
	case keymap.ActionDeleteLine:
		for i := 0; i < count; i++ {
			ed.DeleteLine()
		}
	case keymap.ActionDeleteToEnd:
		ed.DeleteToLineEnd()
	case keymap.ActionYankLine:
		ed.YankLine()
	case keymap.ActionPaste:
		for i := 0; i < count; i++ {
			ed.Paste()
		}
	case keymap.ActionPasteAbove:
		for i := 0; i < count; i++ {
			ed.PasteBefore()
		}
	case keymap.ActionUndo:
		for i := 0; i < count; i++ {
			ed.Undo()
		}
	case keymap.ActionRedo:
		for i := 0; i < count; i++ {
			ed.Redo()
		}
	case keymap.ActionJoinLines:
		for i := 0; i < count; i++ {
			ed.JoinLines()
		}

	// ── Line operations ─────────────────────────────────────────
	case keymap.ActionOpenBelow:
		ed.OpenLineBelow()
	case keymap.ActionOpenAbove:
		ed.OpenLineAbove()

	// ── Insert mode editing ─────────────────────────────────────
	case keymap.ActionBackspace:
		ed.DeleteCharBack()
	case keymap.ActionDelete:
		ed.DeleteChar()
	case keymap.ActionInsertNewline:
		ed.InsertNewline()
	case keymap.ActionInsertTab:
		ed.InsertString("  ")

	// ── Visual mode operators ───────────────────────────────────
	case keymap.ActionOperatorDelete:
		ed.YankSelection()
		ed.EnterNormalMode()
	case keymap.ActionOperatorYank:
		ed.YankSelection()
		ed.EnterNormalMode()
	case keymap.ActionOperatorChange:
		ed.YankSelection()
		ed.EnterNormalMode()

	// ── UI ──────────────────────────────────────────────────────
	case keymap.ActionToggleHelp:
		w.ToggleHelp()
	}

	w.syncModeToKeymap()
}

// executeOperatorMotion handles operator + motion commands (dw, d$, yy, etc.).
func (w *EditorWidget) executeOperatorMotion(cmd keymap.Command) {
	count := cmd.Count
	if count < 1 {
		count = 1
	}

	switch cmd.Action {
	case keymap.ActionOperatorDelete:
		w.deleteWithMotion(cmd.Motion, count)
	case keymap.ActionOperatorYank:
		w.yankWithMotion(cmd.Motion, count)
	case keymap.ActionOperatorChange:
		w.changeWithMotion(cmd.Motion, count)
	}
}

func (w *EditorWidget) deleteWithMotion(motion keymap.Action, count int) {
	ed := w.editor
	switch motion {
	case keymap.ActionMoveWordNext: // dw
		for i := 0; i < count; i++ {
			ed.DeleteWord()
		}
	case keymap.ActionGotoLineEnd: // d$
		ed.DeleteToLineEnd()
	case keymap.ActionDeleteLine: // dd
		for i := 0; i < count; i++ {
			ed.DeleteLine()
		}
	case keymap.ActionMoveDown: // dj
		for i := 0; i < count+1; i++ {
			ed.DeleteLine()
		}
	case keymap.ActionMoveUp: // dk
		if ed.Cursor().Row() > 0 {
			ed.MoveUp(1)
		}
		for i := 0; i < count+1; i++ {
			ed.DeleteLine()
		}
	case keymap.ActionGotoTop: // dgg
		row := ed.Cursor().Row()
		ed.MoveToTop()
		for i := 0; i <= row; i++ {
			ed.DeleteLine()
		}
	case keymap.ActionGotoBottom: // dG
		row := ed.Cursor().Row()
		total := ed.Buffer().LineCount()
		for i := 0; i < total-row; i++ {
			ed.DeleteLine()
		}
	default:
		for i := 0; i < count; i++ {
			ed.DeleteLine()
		}
	}
}

func (w *EditorWidget) yankWithMotion(motion keymap.Action, count int) {
	ed := w.editor
	switch motion {
	case keymap.ActionYankLine: // yy
		ed.YankLine()
	default:
		ed.YankLine()
	}
}

func (w *EditorWidget) changeWithMotion(motion keymap.Action, count int) {
	ed := w.editor
	switch motion {
	case keymap.ActionMoveWordNext: // cw
		for i := 0; i < count; i++ {
			ed.DeleteWord()
		}
		ed.EnterInsertMode()
	case keymap.ActionGotoLineEnd: // c$
		ed.DeleteToLineEnd()
		ed.EnterInsertMode()
	case keymap.ActionDeleteLine: // cc
		ed.MoveToLineStart()
		ed.DeleteToLineEnd()
		ed.EnterInsertMode()
	default:
		ed.MoveToLineStart()
		ed.DeleteToLineEnd()
		ed.EnterInsertMode()
	}
}

// syncModeToKeymap keeps the keymap mode in sync with the editor mode.
func (w *EditorWidget) syncModeToKeymap() {
	switch w.editor.Mode() {
	case editor.ModeNormal:
		w.km.SetMode(keymap.ModeNormal)
	case editor.ModeInsert:
		w.km.SetMode(keymap.ModeInsert)
	case editor.ModeVisual:
		w.km.SetMode(keymap.ModeVisual)
	}
}

// ════════════════════════════════════════════════════════════════
// EXPLAIN POPUP
// ════════════════════════════════════════════════════════════════

// ShowExplain triggers the explain popup for the current line.
func (w *EditorWidget) ShowExplain() {
	// Get current line
	cursor := w.editor.Cursor()
	line := w.editor.Buffer().Line(cursor.Row())

	// If current line is empty, find the last non-empty line
	if line == "" {
		for i := cursor.Row() - 1; i >= 0; i-- {
			candidate := w.editor.Buffer().Line(i)
			if candidate != "" {
				line = candidate
				break
			}
		}
	}

	if line == "" {
		return
	}

	// Get explanation from engine
	result := w.editor.Engine().Explain(line)
	if result == nil {
		return
	}

	// Build explanation text
	w.explainResult = "Input: " + line + "\n\n"

	if len(result.Steps) > 0 {
		w.explainResult += "Steps:\n"
		for _, step := range result.Steps {
			w.explainResult += "  " + step + "\n"
		}
		w.explainResult += "\n"
	}

	w.explainResult += "Result: " + result.Value.String()

	w.activePopup = PopupExplain
	w.Refresh()
}

// ════════════════════════════════════════════════════════════════
// STATE MANAGEMENT
// ════════════════════════════════════════════════════════════════

func (w *EditorWidget) updateState() {
	w.mu.Lock()
	w.state = w.editor.Render()
	w.mu.Unlock()
}

// State returns the current render state.
func (w *EditorWidget) State() *rpc.RenderState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state
}

// Editor returns the underlying editor.
func (w *EditorWidget) Editor() *editor.Editor {
	return w.editor
}

// ════════════════════════════════════════════════════════════════
// CURSOR BLINK
// ════════════════════════════════════════════════════════════════

// ToggleCursor toggles cursor visibility for blinking.
func (w *EditorWidget) ToggleCursor() {
	w.cursorVisible = !w.cursorVisible
	w.Refresh()
}

// ResetCursor makes cursor visible and resets blink.
func (w *EditorWidget) ResetCursor() {
	w.cursorVisible = true
}

// ════════════════════════════════════════════════════════════════
// EDITOR RENDERER
// ════════════════════════════════════════════════════════════════

type editorRenderer struct {
	widget     *EditorWidget
	background *canvas.Rectangle
	objects    []fyne.CanvasObject
}

func newEditorRenderer(w *EditorWidget) *editorRenderer {
	r := &editorRenderer{
		widget:     w,
		background: canvas.NewRectangle(ColorBackground),
	}
	r.objects = []fyne.CanvasObject{r.background}
	return r
}

// Layout positions all objects.
func (r *editorRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
}

// MinSize returns minimum size.
func (r *editorRenderer) MinSize() fyne.Size {
	return r.widget.MinSize()
}

// Refresh redraws the widget.
func (r *editorRenderer) Refresh() {
	r.background.FillColor = ColorBackground
	r.background.Refresh()

	// Rebuild all objects
	r.objects = []fyne.CanvasObject{r.background}

	state := r.widget.State()
	if state == nil {
		return
	}

	size := r.widget.Size()
	w := r.widget

	// Calculate visible lines
	visibleLines := int((size.Height - w.topPadding*2 - w.statusBarHeight) / (w.lineHeight + w.lineSpacing))

	// Render lines
	for i := 0; i < visibleLines; i++ {
		y := w.topPadding + float32(i)*(w.lineHeight+w.lineSpacing)

		if i < len(state.Lines) {
			r.renderLine(state.Lines[i], y, size.Width)
		} else if w.vimMode {
			// Only show tildes in vim mode
			r.renderTilde(y)
		}
	}

	// Render cursor
	if state.Cursor.Visible && r.widget.cursorVisible {
		r.renderCursor(state.Cursor)
	}

	// Render status bar
	r.renderStatusBar(state, size)

	// Render popup if active
	if w.activePopup != PopupNone {
		r.renderPopup(size)
	}
}

func (r *editorRenderer) renderLine(line rpc.RenderLine, y float32, width float32) {
	w := r.widget

	// Relative line numbers (vim mode only)
	gutterWidth := float32(0)
	if w.vimMode {
		gutterWidth = w.lineNumWidth
		numStr := fmt.Sprintf("%3d", line.RelativeNumber)
		numColor := ColorMuted
		if line.IsCurrentLine {
			numColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}
		}
		numText := canvas.NewText(numStr, numColor)
		numText.TextSize = 14
		numText.TextStyle = fyne.TextStyle{Monospace: true}
		numText.Move(fyne.NewPos(w.leftPadding, y))
		r.objects = append(r.objects, numText)
	}

	// Input spans (left-aligned, after gutter)
	x := w.leftPadding + gutterWidth
	for _, span := range line.Input {
		text := canvas.NewText(span.Text, StyleColor(string(span.Style)))
		text.TextSize = 14
		text.TextStyle = fyne.TextStyle{Monospace: true}
		text.Move(fyne.NewPos(x, y))
		r.objects = append(r.objects, text)
		x += text.MinSize().Width
	}

	// Result spans (right-aligned)
	if len(line.Result) > 0 {
		// Calculate result width
		resultWidth := float32(0)
		for _, span := range line.Result {
			text := canvas.NewText(span.Text, StyleColor(string(span.Style)))
			text.TextSize = 14
			text.TextStyle = fyne.TextStyle{Monospace: true}
			resultWidth += text.MinSize().Width
		}

		// Position from right
		resultX := width - w.rightPadding - resultWidth

		for _, span := range line.Result {
			text := canvas.NewText(span.Text, StyleColor(string(span.Style)))
			text.TextSize = 14
			text.TextStyle = fyne.TextStyle{Monospace: true}
			text.Move(fyne.NewPos(resultX, y))
			r.objects = append(r.objects, text)
			resultX += text.MinSize().Width
		}
	}
}

func (r *editorRenderer) renderTilde(y float32) {
	w := r.widget
	gutterWidth := float32(0)
	if w.vimMode {
		gutterWidth = w.lineNumWidth
	}
	tilde := canvas.NewText("~", ColorTilde)
	tilde.TextSize = 14
	tilde.TextStyle = fyne.TextStyle{Monospace: true}
	tilde.Move(fyne.NewPos(w.leftPadding+gutterWidth, y))
	r.objects = append(r.objects, tilde)
}

func (r *editorRenderer) renderCursor(cursor rpc.CursorState) {
	w := r.widget

	gutterWidth := float32(0)
	if w.vimMode {
		gutterWidth = w.lineNumWidth
	}
	x := w.leftPadding + gutterWidth + float32(cursor.Col)*w.charWidth
	y := w.topPadding + float32(cursor.Row)*(w.lineHeight+w.lineSpacing)

	// Cursor should match text size
	cursorHeight := float32(16)
	cursorY := y + (w.lineHeight-cursorHeight)/2

	var cursorRect *canvas.Rectangle

	switch cursor.Style {
	case rpc.CursorBlock:
		cursorRect = canvas.NewRectangle(ColorCursorBlock)
		cursorRect.Resize(fyne.NewSize(w.charWidth, cursorHeight))
		cursorRect.Move(fyne.NewPos(x, cursorY))

	case rpc.CursorLine:
		cursorRect = canvas.NewRectangle(ColorCursor)
		cursorRect.Resize(fyne.NewSize(2, cursorHeight))
		cursorRect.Move(fyne.NewPos(x, cursorY))

	case rpc.CursorUnderline:
		cursorRect = canvas.NewRectangle(ColorCursor)
		cursorRect.Resize(fyne.NewSize(w.charWidth, 2))
		cursorRect.Move(fyne.NewPos(x, y+w.lineHeight-4))
	}

	if cursorRect != nil {
		r.objects = append(r.objects, cursorRect)
	}
}

func (r *editorRenderer) renderStatusBar(state *rpc.RenderState, size fyne.Size) {
	w := r.widget

	// Background
	barY := size.Height - w.statusBarHeight
	barBg := canvas.NewRectangle(ColorStatusBar)
	barBg.Resize(fyne.NewSize(size.Width, w.statusBarHeight))
	barBg.Move(fyne.NewPos(0, barY))
	r.objects = append(r.objects, barBg)

	textY := barY + 6

	// Mode indicator (left) — show in vim mode, also show pending operator
	if w.vimMode && state.StatusBar != nil && state.StatusBar.Mode != "" {
		modeLabel := state.StatusBar.Mode
		modeColor := ModeColor(state.StatusBar.Mode)

		// Show pending operator/count
		if pending := w.km.State.PendingDisplay(); pending != "" {
			modeLabel += " " + pending
			modeColor = ColorPending
		}

		modeText := canvas.NewText(modeLabel, modeColor)
		modeText.TextSize = 12
		modeText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		modeText.Move(fyne.NewPos(w.leftPadding, textY))
		r.objects = append(r.objects, modeText)
	}

	// Help hints (center)
	var hints string
	if w.vimMode {
		hints = "⌘/ help   ⌘E explain   ⌘K simple mode"
	} else {
		hints = "⌘/ help   ⌘E explain   ⌘K vim mode"
	}
	hintsText := canvas.NewText(hints, ColorStatusText)
	hintsText.TextSize = 11
	hintsText.TextStyle = fyne.TextStyle{Monospace: true}
	hintsX := (size.Width - hintsText.MinSize().Width) / 2
	hintsText.Move(fyne.NewPos(hintsX, textY))
	r.objects = append(r.objects, hintsText)

	// Position (right)
	if state.StatusBar != nil && state.StatusBar.Position != "" {
		posText := canvas.NewText(state.StatusBar.Position, ColorStatusText)
		posText.TextSize = 12
		posText.TextStyle = fyne.TextStyle{Monospace: true}
		posX := size.Width - w.rightPadding - posText.MinSize().Width
		posText.Move(fyne.NewPos(posX, textY))
		r.objects = append(r.objects, posText)
	}
}

func (r *editorRenderer) renderPopup(size fyne.Size) {
	w := r.widget

	var title string
	var content []string

	switch w.activePopup {
	case PopupHelp:
		title = "Help"
		if w.vimMode {
			content = []string{
				"",
				"Navigation",
				"  h/j/k/l      Move cursor",
				"  w/b          Next/prev word",
				"  0/$          Start/end of line",
				"  gg/G         Top/bottom of file",
				"  5j           Move 5 lines down",
				"  Ctrl+U/D     Page up/down",
				"",
				"Editing",
				"  i/a          Insert/append mode",
				"  I/A          Insert start/end of line",
				"  o/O          Open line below/above",
				"  x            Delete character",
				"  dd           Delete line",
				"  dw           Delete word",
				"  d$           Delete to end of line",
				"  yy           Yank line",
				"  p/P          Paste after/before",
				"  u            Undo",
				"  Ctrl+R       Redo",
				"  J            Join lines",
				"",
				"General",
				"  Esc          Normal mode",
				"  v/V          Visual / visual-line mode",
				"  ?            Toggle this help",
				"  ⌘/           Toggle help",
				"  ⌘E           Explain calculation",
				"  ⌘K           Switch to simple mode",
				"",
				"Press any key to close",
			}
		} else {
			content = []string{
				"",
				"Numio Calculator",
				"",
				"  Type expressions and see results",
				"  Examples:",
				"    1 + 1",
				"    500 usd to eur",
				"    100 km to miles",
				"    time in tokyo",
				"",
				"Shortcuts",
				"  ⌘/            Help",
				"  ⌘E            Explain calculation",
				"  ⌘K            Switch to vim mode",
				"",
				"Press any key to close",
			}
		}

	case PopupExplain:
		title = "Explanation"
		content = []string{""}
		line := ""
		for _, ch := range w.explainResult {
			if ch == '\n' {
				content = append(content, "  "+line)
				line = ""
			} else {
				line += string(ch)
			}
		}
		if line != "" {
			content = append(content, "  "+line)
		}
		content = append(content, "", "Press any key to close")
	}

	// Calculate popup dimensions
	popupWidth := float32(360)
	popupHeight := float32(len(content)*18 + 50)
	popupX := (size.Width - popupWidth) / 2
	popupY := (size.Height - popupHeight) / 2

	// Background
	bg := canvas.NewRectangle(ColorPopupBackground)
	bg.Resize(fyne.NewSize(popupWidth, popupHeight))
	bg.Move(fyne.NewPos(popupX, popupY))
	bg.StrokeColor = ColorPopupBorder
	bg.StrokeWidth = 1
	r.objects = append(r.objects, bg)

	// Title
	titleText := canvas.NewText(title, ColorPopupTitle)
	titleText.TextSize = 14
	titleText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	titleText.Move(fyne.NewPos(popupX+16, popupY+12))
	r.objects = append(r.objects, titleText)

	// Content
	y := popupY + 36
	for _, line := range content {
		var textColor color.Color = ColorPopupDesc

		// Check if it's a heading (no leading spaces, not empty)
		if len(line) > 0 && line[0] != ' ' {
			textColor = ColorPopupHeading
		}

		// Check if it's the hint line
		if line == "Press any key to close" {
			textColor = ColorPopupHint
		}

		text := canvas.NewText(line, textColor)
		text.TextSize = 12
		text.TextStyle = fyne.TextStyle{Monospace: true}
		text.Move(fyne.NewPos(popupX+16, y))
		r.objects = append(r.objects, text)
		y += 18
	}
}

// Objects returns all canvas objects.
func (r *editorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy cleans up resources.
func (r *editorRenderer) Destroy() {}
