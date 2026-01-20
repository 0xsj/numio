package editor

import (
	"unicode/utf8"

	"github.com/0xsj/numio/internal/rpc"
)

// ════════════════════════════════════════════════════════════════
// RENDERER
// ════════════════════════════════════════════════════════════════

// Renderer produces styled output for the UI.
type Renderer struct {
	// StyleFunc maps token types to RPC span styles.
	// If nil, uses default mapping.
	StyleFunc func(tokenType string) rpc.SpanStyle

	// ShowLineNumbers enables line number rendering.
	ShowLineNumbers bool

	// LineNumberWidth is the width reserved for line numbers.
	LineNumberWidth int

	// ResultSeparator is the string between input and result.
	ResultSeparator string
}

// NewRenderer creates a renderer with default settings.
func NewRenderer() *Renderer {
	return &Renderer{
		StyleFunc:       defaultStyleFunc,
		ShowLineNumbers: true,
		LineNumberWidth: 4,
		ResultSeparator: " = ",
	}
}

// ════════════════════════════════════════════════════════════════
// RENDER STATE
// ════════════════════════════════════════════════════════════════

// RenderParams contains parameters for rendering.
type RenderParams struct {
	Buffer    *Buffer
	Cursor    *Cursor
	Selection *Selection
	Viewport  *Viewport
	Mode      Mode

	// Results maps line indices to their evaluation results.
	// Key is 0-based line index, value is the styled result.
	Results map[int][]rpc.Span

	// Errors maps line indices to error messages.
	Errors map[int]string
}

// Render produces a complete RenderState for the UI.
func (r *Renderer) Render(params RenderParams) *rpc.RenderState {
	state := &rpc.RenderState{
		Lines:  r.renderLines(params),
		Cursor: r.renderCursor(params),
		Mode:   modeToRPC(params.Mode),
		Dirty:  false, // TODO: track dirty state
	}

	state.StatusBar = r.renderStatusBar(params)

	return state
}

// ════════════════════════════════════════════════════════════════
// LINE RENDERING
// ════════════════════════════════════════════════════════════════

func (r *Renderer) renderLines(params RenderParams) []rpc.RenderLine {
	if params.Buffer == nil || params.Viewport == nil {
		return nil
	}

	startRow, endRow := params.Viewport.VisibleLines(params.Buffer.LineCount())
	lines := make([]rpc.RenderLine, 0, endRow-startRow)

	for row := startRow; row < endRow; row++ {
		line := r.renderLine(row, params)
		lines = append(lines, line)
	}

	return lines
}

func (r *Renderer) renderLine(row int, params RenderParams) rpc.RenderLine {
	content := params.Buffer.Line(row)
	isCurrentLine := params.Cursor != nil && params.Cursor.Row() == row

	renderLine := rpc.RenderLine{
		Number:        row + 1, // 1-based for display
		IsCurrentLine: isCurrentLine,
	}

	// Render input with syntax highlighting and selection
	renderLine.Input = r.renderLineContent(content, row, params)

	// Add result if available
	if params.Results != nil {
		if result, ok := params.Results[row]; ok && len(result) > 0 {
			renderLine.Result = result
		}
	}

	// Add error styling if there's an error
	if params.Errors != nil {
		if errMsg, ok := params.Errors[row]; ok && errMsg != "" {
			renderLine.Result = []rpc.Span{
				{Text: errMsg, Style: rpc.StyleError},
			}
		}
	}

	return renderLine
}

func (r *Renderer) renderLineContent(content string, row int, params RenderParams) []rpc.Span {
	if content == "" {
		return []rpc.Span{}
	}

	// Start with basic tokenization/highlighting
	spans := r.tokenizeLine(content)

	// Apply selection highlighting if applicable
	if params.Selection != nil && params.Selection.ContainsRow(row) {
		spans = r.applySelection(spans, row, params.Selection, len(content))
	}

	return spans
}

// tokenizeLine performs basic syntax highlighting.
// This is a simplified version - the full implementation would use internal/highlight.
func (r *Renderer) tokenizeLine(content string) []rpc.Span {
	if content == "" {
		return []rpc.Span{}
	}

	// Check for comments
	if len(content) >= 2 && (content[:2] == "//" || content[0] == '#') {
		return []rpc.Span{
			{Text: content, Style: rpc.StyleComment},
		}
	}

	// Basic tokenization using runes for proper Unicode support
	spans := make([]rpc.Span, 0)
	current := ""
	currentStyle := rpc.StyleDefault

	for i := 0; i < len(content); {
		r, size := utf8.DecodeRuneInString(content[i:])
		newStyle := classifyRune(r)

		if newStyle != currentStyle && current != "" {
			spans = append(spans, rpc.Span{Text: current, Style: currentStyle})
			current = ""
		}

		current += string(r)
		currentStyle = newStyle
		i += size
	}

	if current != "" {
		spans = append(spans, rpc.Span{Text: current, Style: currentStyle})
	}

	return spans
}

func classifyRune(r rune) rpc.SpanStyle {
	switch {
	case isDigitRune(r) || r == '.':
		return rpc.StyleNumber

	case isOperatorRune(r):
		return rpc.StyleOperator

	case r == '%':
		return rpc.StylePercent

	case isCurrencyRune(r):
		return rpc.StyleCurrency

	default:
		return rpc.StyleDefault
	}
}

func isDigitRune(r rune) bool {
	return r >= '0' && r <= '9'
}

func isOperatorRune(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' ||
		r == '^' || r == '=' || r == '<' || r == '>'
}

func isCurrencyRune(r rune) bool {
	switch r {
	case '$', '€', '£', '¥', '₺', '₹', '₽', '₿', '₩', '฿':
		return true
	default:
		return false
	}
}

// ════════════════════════════════════════════════════════════════
// SELECTION RENDERING
// ════════════════════════════════════════════════════════════════

func (r *Renderer) applySelection(spans []rpc.Span, row int, sel *Selection, lineLen int) []rpc.Span {
	startRow, startCol, endRow, endCol := sel.Bounds()

	// Determine selection range for this row
	var selStart, selEnd int

	switch sel.Mode() {
	case SelectionLine:
		selStart = 0
		selEnd = lineLen

	case SelectionChar:
		if row == startRow && row == endRow {
			selStart = startCol
			selEnd = endCol
		} else if row == startRow {
			selStart = startCol
			selEnd = lineLen
		} else if row == endRow {
			selStart = 0
			selEnd = endCol
		} else {
			selStart = 0
			selEnd = lineLen
		}

	case SelectionBlock:
		minCol, maxCol := sel.BlockBounds()
		selStart = minCol
		selEnd = maxCol
	}

	// Apply selection style to spans
	return r.splitSpansForSelection(spans, selStart, selEnd)
}

func (r *Renderer) splitSpansForSelection(spans []rpc.Span, selStart, selEnd int) []rpc.Span {
	if selStart >= selEnd {
		return spans
	}

	result := make([]rpc.Span, 0, len(spans)*2)
	pos := 0

	for _, span := range spans {
		spanLen := utf8.RuneCountInString(span.Text)
		spanStart := pos
		spanEnd := pos + spanLen

		// Before selection
		if spanEnd <= selStart || spanStart >= selEnd {
			result = append(result, span)
		} else {
			// Span overlaps with selection
			runes := []rune(span.Text)

			// Part before selection
			if spanStart < selStart {
				beforeLen := selStart - spanStart
				result = append(result, rpc.Span{
					Text:  string(runes[:beforeLen]),
					Style: span.Style,
				})
			}

			// Selected part
			selectStart := maxInt(0, selStart-spanStart)
			selectEnd := minInt(spanLen, selEnd-spanStart)
			if selectStart < selectEnd {
				result = append(result, rpc.Span{
					Text:  string(runes[selectStart:selectEnd]),
					Style: rpc.StyleSelection,
				})
			}

			// Part after selection
			if spanEnd > selEnd {
				afterStart := selEnd - spanStart
				result = append(result, rpc.Span{
					Text:  string(runes[afterStart:]),
					Style: span.Style,
				})
			}
		}

		pos = spanEnd
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// CURSOR RENDERING
// ════════════════════════════════════════════════════════════════

func (r *Renderer) renderCursor(params RenderParams) rpc.CursorState {
	if params.Cursor == nil {
		return rpc.CursorState{Visible: false}
	}

	cursorState := rpc.CursorState{
		Row:     params.Cursor.Row(),
		Col:     params.Cursor.Col(),
		Visible: true,
		Style:   cursorStyleForMode(params.Mode),
	}

	// Add selection info
	if params.Selection != nil {
		startRow, startCol, endRow, endCol := params.Selection.Bounds()
		cursorState.Selection = &rpc.Selection{
			StartRow: startRow,
			StartCol: startCol,
			EndRow:   endRow,
			EndCol:   endCol,
		}
	}

	return cursorState
}

func cursorStyleForMode(mode Mode) rpc.CursorStyle {
	switch mode {
	case ModeInsert:
		return rpc.CursorLine
	case ModeVisual:
		return rpc.CursorBlock
	default:
		return rpc.CursorBlock
	}
}

// ════════════════════════════════════════════════════════════════
// STATUS BAR RENDERING
// ════════════════════════════════════════════════════════════════

func (r *Renderer) renderStatusBar(params RenderParams) *rpc.StatusBar {
	if params.Cursor == nil {
		return nil
	}

	row, col := params.Cursor.Position()

	return &rpc.StatusBar{
		Mode:     modeString(params.Mode),
		Position: formatPosition(row+1, col+1), // 1-based for display
	}
}

// ════════════════════════════════════════════════════════════════
// MODE CONVERSION
// ════════════════════════════════════════════════════════════════

// Mode represents the editor mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeVisual
)

func modeToRPC(mode Mode) rpc.EditorMode {
	switch mode {
	case ModeInsert:
		return rpc.ModeInsert
	case ModeVisual:
		return rpc.ModeVisual
	default:
		return rpc.ModeNormal
	}
}

func modeString(mode Mode) string {
	switch mode {
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	default:
		return "NORMAL"
	}
}

// ════════════════════════════════════════════════════════════════
// DEFAULT STYLE MAPPING
// ════════════════════════════════════════════════════════════════

func defaultStyleFunc(tokenType string) rpc.SpanStyle {
	switch tokenType {
	case "number":
		return rpc.StyleNumber
	case "operator":
		return rpc.StyleOperator
	case "percent":
		return rpc.StylePercent
	case "currency":
		return rpc.StyleCurrency
	case "unit":
		return rpc.StyleUnit
	case "metal":
		return rpc.StyleMetal
	case "crypto":
		return rpc.StyleCrypto
	case "function":
		return rpc.StyleFunction
	case "variable":
		return rpc.StyleVariable
	case "comment":
		return rpc.StyleComment
	case "string":
		return rpc.StyleString
	case "keyword":
		return rpc.StyleKeyword
	case "error":
		return rpc.StyleError
	default:
		return rpc.StyleDefault
	}
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func formatPosition(row, col int) string {
	return intToString(row) + ":" + intToString(col)
}

func intToString(n int) string {
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
