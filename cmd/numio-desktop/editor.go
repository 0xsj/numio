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
	lineNumWidth    float32
	statusBarHeight float32

	// Cursor blink state
	cursorVisible bool

	// Modifier key state
	ctrlPressed  bool
	altPressed   bool
	shiftPressed bool
	superPressed bool
}

// NewEditorWidget creates a new editor widget.
func NewEditorWidget() *EditorWidget {
	w := &EditorWidget{
		editor:          editor.NewEditor(),
		lineHeight:      20,
		charWidth:       8.4,
		leftPadding:     12,
		rightPadding:    12,
		topPadding:      8,
		lineNumWidth:    40,
		statusBarHeight: 24,
		cursorVisible:   true,
	}

	w.ExtendBaseWidget(w)
	w.updateState()

	return w
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
	cols := int((size.Width - w.leftPadding - w.rightPadding - w.lineNumWidth) / w.charWidth)
	rows := int((size.Height - w.topPadding*2 - w.statusBarHeight) / w.lineHeight)

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
	w.mu.Lock()
	w.handleRune(r)
	w.mu.Unlock()

	w.updateState()
	w.Refresh()
}

// TypedKey handles special key input.
func (w *EditorWidget) TypedKey(ev *fyne.KeyEvent) {
	w.mu.Lock()
	w.handleKey(ev)
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
	// Track modifier state
	switch ev.Name {
	case desktop.KeyControlLeft, desktop.KeyControlRight:
		w.ctrlPressed = true
		return
	case desktop.KeyAltLeft, desktop.KeyAltRight:
		w.altPressed = true
		return
	case desktop.KeyShiftLeft, desktop.KeyShiftRight:
		w.shiftPressed = true
		return
	case desktop.KeySuperLeft, desktop.KeySuperRight:
		w.superPressed = true
		return
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
// KEY HANDLING
// ════════════════════════════════════════════════════════════════

func (w *EditorWidget) handleKey(ev *fyne.KeyEvent) {
	keyEvent := TranslateKeyEventWithMods(ev, w.ctrlPressed, w.altPressed, w.shiftPressed, w.superPressed)
	w.processKey(keyEvent)
}

func (w *EditorWidget) handleRune(r rune) {
	keyEvent := TranslateRuneWithMods(r, w.ctrlPressed, w.altPressed, w.shiftPressed, w.superPressed)
	w.processKey(keyEvent)
}

func (w *EditorWidget) processKey(ev KeyEvent) {
	mode := w.editor.Mode()

	switch mode {
	case editor.ModeNormal:
		w.processNormalModeKey(ev)
	case editor.ModeInsert:
		w.processInsertModeKey(ev)
	case editor.ModeVisual:
		w.processVisualModeKey(ev)
	}
}

func (w *EditorWidget) processNormalModeKey(ev KeyEvent) {
	key := ev.Key

	// Handle Ctrl combinations
	if ev.IsCtrl() {
		switch key {
		case "r":
			w.editor.Redo()
		case "u":
			w.editor.MoveUp(w.editor.Viewport().Height() / 2)
		case "d":
			w.editor.MoveDown(w.editor.Viewport().Height() / 2)
		}
		return
	}

	// Normal mode keys
	switch key {
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

func (w *EditorWidget) processInsertModeKey(ev KeyEvent) {
	key := ev.Key

	// Escape exits insert mode
	if key == "Escape" {
		w.editor.EnterNormalMode()
		return
	}

	// Handle Ctrl combinations
	if ev.IsCtrl() {
		switch key {
		case "c":
			w.editor.EnterNormalMode()
		case "w":
			w.editor.DeleteWord()
		case "u":
			w.editor.DeleteToLineEnd()
		}
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
	default:
		// Insert regular character
		if len(key) == 1 {
			w.editor.InsertChar(rune(key[0]))
		}
	}
}

func (w *EditorWidget) processVisualModeKey(ev KeyEvent) {
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

	// Render lines
	for i, line := range state.Lines {
		y := r.widget.topPadding + float32(i)*r.widget.lineHeight
		r.renderLine(line, y, size.Width)
	}

	// Render cursor
	if state.Cursor.Visible && r.widget.cursorVisible {
		r.renderCursor(state.Cursor)
	}

	// Render status bar
	if state.StatusBar != nil {
		r.renderStatusBar(state.StatusBar, size)
	}
}

func (r *editorRenderer) renderLine(line rpc.RenderLine, y float32, width float32) {
	w := r.widget

	// Line number
	lineNumText := canvas.NewText(intToStr(line.Number), ColorLineNumber)
	lineNumText.TextSize = 13
	lineNumText.TextStyle = fyne.TextStyle{Monospace: true}
	lineNumX := w.leftPadding + w.lineNumWidth - lineNumText.MinSize().Width - 8
	lineNumText.Move(fyne.NewPos(lineNumX, y))
	r.objects = append(r.objects, lineNumText)

	// Input spans
	x := w.leftPadding + w.lineNumWidth
	for _, span := range line.Input {
		text := canvas.NewText(span.Text, StyleColor(string(span.Style)))
		text.TextSize = 13
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
			text.TextSize = 13
			text.TextStyle = fyne.TextStyle{Monospace: true}
			resultWidth += text.MinSize().Width
		}

		// Add "=" separator
		eqText := canvas.NewText(" = ", ColorComment)
		eqText.TextSize = 13
		eqText.TextStyle = fyne.TextStyle{Monospace: true}
		resultWidth += eqText.MinSize().Width

		// Position from right
		resultX := width - w.rightPadding - resultWidth

		eqText.Move(fyne.NewPos(resultX, y))
		r.objects = append(r.objects, eqText)
		resultX += eqText.MinSize().Width

		for _, span := range line.Result {
			text := canvas.NewText(span.Text, StyleColor(string(span.Style)))
			text.TextSize = 13
			text.TextStyle = fyne.TextStyle{Monospace: true}
			text.Move(fyne.NewPos(resultX, y))
			r.objects = append(r.objects, text)
			resultX += text.MinSize().Width
		}
	}
}

func (r *editorRenderer) renderCursor(cursor rpc.CursorState) {
	w := r.widget

	x := w.leftPadding + w.lineNumWidth + float32(cursor.Col)*w.charWidth
	y := w.topPadding + float32(cursor.Row)*w.lineHeight

	var cursorRect *canvas.Rectangle

	switch cursor.Style {
	case rpc.CursorBlock:
		cursorRect = canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 180})
		cursorRect.Resize(fyne.NewSize(w.charWidth, w.lineHeight-4))
		cursorRect.Move(fyne.NewPos(x, y+2))

	case rpc.CursorLine:
		cursorRect = canvas.NewRectangle(ColorCursor)
		cursorRect.Resize(fyne.NewSize(2, w.lineHeight-4))
		cursorRect.Move(fyne.NewPos(x, y+2))

	case rpc.CursorUnderline:
		cursorRect = canvas.NewRectangle(ColorCursor)
		cursorRect.Resize(fyne.NewSize(w.charWidth, 2))
		cursorRect.Move(fyne.NewPos(x, y+w.lineHeight-2))
	}

	if cursorRect != nil {
		r.objects = append(r.objects, cursorRect)
	}
}

func (r *editorRenderer) renderStatusBar(statusBar *rpc.StatusBar, size fyne.Size) {
	w := r.widget

	// Background
	barY := size.Height - w.statusBarHeight
	barBg := canvas.NewRectangle(ColorStatusBar)
	barBg.Resize(fyne.NewSize(size.Width, w.statusBarHeight))
	barBg.Move(fyne.NewPos(0, barY))
	r.objects = append(r.objects, barBg)

	// Mode indicator
	if statusBar.Mode != "" {
		modeText := canvas.NewText(" "+statusBar.Mode+" ", ModeColor(statusBar.Mode))
		modeText.TextSize = 11
		modeText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		modeText.Move(fyne.NewPos(w.leftPadding, barY+5))
		r.objects = append(r.objects, modeText)
	}

	// Position
	if statusBar.Position != "" {
		posText := canvas.NewText(statusBar.Position, ColorComment)
		posText.TextSize = 11
		posText.TextStyle = fyne.TextStyle{Monospace: true}
		posX := size.Width - w.rightPadding - posText.MinSize().Width
		posText.Move(fyne.NewPos(posX, barY+5))
		r.objects = append(r.objects, posText)
	}
}

// Objects returns all canvas objects.
func (r *editorRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy cleans up resources.
func (r *editorRenderer) Destroy() {}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var buf [20]byte
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}

	if negative {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}
