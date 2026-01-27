package main

import (
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/0xsj/numio/internal/editor"
	"github.com/0xsj/numio/internal/rpc"
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

	// Vim mode toggle
	vimMode bool
}

// NewEditorWidget creates a new editor widget.
func NewEditorWidget() *EditorWidget {
	w := &EditorWidget{
		editor:          editor.NewEditor(),
		lineHeight:      24,
		charWidth:       8.4,
		leftPadding:     20,
		rightPadding:    20,
		topPadding:      16,
		statusBarHeight: 28,
		lineSpacing:     4,
		cursorVisible:   true,
		activePopup:     PopupNone,
		vimMode:         false, // Default: normal editor mode
	}

	// Start in insert mode when vim mode is off
	w.editor.EnterInsertMode()

	w.ExtendBaseWidget(w)
	w.updateState()

	return w
}

// ════════════════════════════════════════════════════════════════
// VIM MODE TOGGLE
// ════════════════════════════════════════════════════════════════

// SetVimMode enables or disables vim keybindings.
func (w *EditorWidget) SetVimMode(enabled bool) {
	w.vimMode = enabled
	if !enabled {
		// Always stay in insert mode when vim mode is off
		w.editor.EnterInsertMode()
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

	// Calculate viewport size in characters
	cols := int((size.Width - w.leftPadding - w.rightPadding) / w.charWidth)
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
	// Handle Cmd+key combinations
	if w.superPressed {
		switch r {
		case '/', '÷': // ÷ is what macOS might send for Cmd+/
			w.toggleHelp()
			return
		case 'e', 'E':
			w.showExplain()
			w.Refresh()
			return
		case 'k', 'K':
			w.ToggleVimMode()
			return
		}
	}

	// Close popup on any key
	if w.activePopup != PopupNone {
		w.activePopup = PopupNone
		w.Refresh()
		return
	}

	w.mu.Lock()

	if w.vimMode {
		// Vim mode: handle based on current mode
		if w.editor.Mode() == editor.ModeInsert {
			w.editor.InsertChar(r)
		} else {
			w.processVimKey(KeyEvent{Key: string(r)})
		}
	} else {
		// Normal editor mode: always insert
		w.editor.InsertChar(r)
	}

	w.mu.Unlock()
	w.updateState()
	w.Refresh()
}

// TypedKey handles special key input.
func (w *EditorWidget) TypedKey(ev *fyne.KeyEvent) {
	// Handle Cmd+key via TypedKey as well (some keys come here)
	if w.superPressed {
		switch ev.Name {
		case fyne.KeySlash:
			w.toggleHelp()
			return
		case fyne.KeyE:
			w.showExplain()
			w.Refresh()
			return
		case fyne.KeyK:
			w.ToggleVimMode()
			return
		}
	}

	// Keep F-keys as backup
	switch ev.Name {
	case fyne.KeyF1:
		w.toggleHelp()
		return
	case fyne.KeyF2:
		w.showExplain()
		w.Refresh()
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

	key := translateKey(ev.Name)

	w.mu.Lock()

	// Skip if it's a regular character (handled by TypedRune)
	if len(key) == 1 && key[0] >= 32 && key[0] < 127 {
		w.mu.Unlock()
		return
	}

	if w.vimMode {
		w.processVimKey(KeyEvent{
			Key:       key,
			Modifiers: w.currentModifiers(),
		})
	} else {
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

func (w *EditorWidget) toggleHelp() {
	if w.activePopup == PopupHelp {
		w.activePopup = PopupNone
	} else {
		w.activePopup = PopupHelp
	}
	w.Refresh()
}

// ════════════════════════════════════════════════════════════════
// NORMAL EDITOR KEY HANDLING (non-vim)
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
// VIM KEY HANDLING
// ════════════════════════════════════════════════════════════════

func (w *EditorWidget) currentModifiers() []string {
	mods := make([]string, 0, 4)
	if w.ctrlPressed {
		mods = append(mods, "ctrl")
	}
	if w.altPressed {
		mods = append(mods, "alt")
	}
	if w.shiftPressed {
		mods = append(mods, "shift")
	}
	if w.superPressed {
		mods = append(mods, "cmd")
	}
	return mods
}

func (w *EditorWidget) processVimKey(ev KeyEvent) {
	mode := w.editor.Mode()

	switch mode {
	case editor.ModeNormal:
		w.processVimNormalModeKey(ev)
	case editor.ModeInsert:
		w.processVimInsertModeKey(ev)
	case editor.ModeVisual:
		w.processVimVisualModeKey(ev)
	}
}

func (w *EditorWidget) processVimNormalModeKey(ev KeyEvent) {
	key := ev.Key

	switch key {
	// Help & Explain (also available via Cmd+/)
	case "?":
		w.toggleHelp()
	case "e":
		w.showExplain()

	// Mode switching
	case "i":
		w.editor.EnterInsertMode()
	case "I":
		w.editor.EnterInsertModeLineStart()
	case "a":
		w.editor.EnterInsertModeAppend()
	case "A":
		w.editor.EnterInsertModeLineEnd()
	case "o":
		w.editor.OpenLineBelow()
	case "O":
		w.editor.OpenLineAbove()
	case "v":
		w.editor.EnterVisualMode()
	case "V":
		w.editor.EnterVisualLineMode()

	// Movement
	case "h", "ArrowLeft":
		w.editor.MoveLeft(1)
	case "j", "ArrowDown":
		w.editor.MoveDown(1)
	case "k", "ArrowUp":
		w.editor.MoveUp(1)
	case "l", "ArrowRight":
		w.editor.MoveRight(1)
	case "w":
		w.editor.MoveWordForward(1)
	case "b":
		w.editor.MoveWordBackward(1)
	case "0":
		w.editor.MoveToLineStart()
	case "$":
		w.editor.MoveToLineEnd()
	case "^":
		w.editor.MoveToFirstNonBlank()
	case "g":
		w.editor.MoveToTop()
	case "G":
		w.editor.MoveToBottom()

	// Editing
	case "x":
		w.editor.DeleteChar()
	case "X":
		w.editor.DeleteCharBack()
	case "d":
		w.editor.DeleteLine()
	case "D":
		w.editor.DeleteToLineEnd()
	case "y":
		w.editor.YankLine()
	case "p":
		w.editor.Paste()
	case "P":
		w.editor.PasteBefore()
	case "u":
		w.editor.Undo()
	case "J":
		w.editor.JoinLines()
	}
}

func (w *EditorWidget) processVimInsertModeKey(ev KeyEvent) {
	key := ev.Key

	// Escape exits insert mode
	if key == "Escape" {
		w.editor.EnterNormalMode()
		return
	}

	// Special keys
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
	}
}

func (w *EditorWidget) processVimVisualModeKey(ev KeyEvent) {
	key := ev.Key

	// Escape exits visual mode
	if key == "Escape" {
		w.editor.EnterNormalMode()
		return
	}

	// Movement extends selection
	switch key {
	case "h", "ArrowLeft":
		w.editor.MoveLeft(1)
	case "j", "ArrowDown":
		w.editor.MoveDown(1)
	case "k", "ArrowUp":
		w.editor.MoveUp(1)
	case "l", "ArrowRight":
		w.editor.MoveRight(1)
	case "w":
		w.editor.MoveWordForward(1)
	case "b":
		w.editor.MoveWordBackward(1)
	case "0":
		w.editor.MoveToLineStart()
	case "$":
		w.editor.MoveToLineEnd()
	case "G":
		w.editor.MoveToBottom()
	case "g":
		w.editor.MoveToTop()

	// Actions
	case "y":
		w.editor.YankSelection()
		w.editor.EnterNormalMode()
	case "d", "x":
		w.editor.YankSelection()
		w.editor.EnterNormalMode()
	case "v":
		w.editor.EnterNormalMode()
	}
}

// ════════════════════════════════════════════════════════════════
// EXPLAIN POPUP
// ════════════════════════════════════════════════════════════════

func (w *EditorWidget) showExplain() {
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

	// Input spans (left-aligned)
	x := w.leftPadding
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
	tilde := canvas.NewText("~", ColorTilde)
	tilde.TextSize = 14
	tilde.TextStyle = fyne.TextStyle{Monospace: true}
	tilde.Move(fyne.NewPos(w.leftPadding, y))
	r.objects = append(r.objects, tilde)
}

func (r *editorRenderer) renderCursor(cursor rpc.CursorState) {
	w := r.widget

	x := w.leftPadding + float32(cursor.Col)*w.charWidth
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

	// Mode indicator (left) - only show in vim mode
	if w.vimMode && state.StatusBar != nil && state.StatusBar.Mode != "" {
		modeText := canvas.NewText(state.StatusBar.Mode, ModeColor(state.StatusBar.Mode))
		modeText.TextSize = 12
		modeText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		modeText.Move(fyne.NewPos(w.leftPadding, textY))
		r.objects = append(r.objects, modeText)
	}

	// Help hints (center)
	var hints string
	if w.vimMode {
		hints = "⌘/ help   ⌘E explain   ⌘K normal mode"
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
				"",
				"Editing",
				"  i/a          Insert/append mode",
				"  o/O          Open line below/above",
				"  x            Delete character",
				"  dd           Delete line",
				"  u            Undo",
				"  p            Paste",
				"",
				"General",
				"  Esc          Normal mode",
				"  ⌘/           Toggle help",
				"  ⌘E           Explain calculation",
				"  ⌘K           Switch to normal mode",
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