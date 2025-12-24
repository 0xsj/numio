// internal/highlight/highlight.go

package highlight

import (
	"strings"

	"github.com/0xsj/numio/internal/lexer"
	"github.com/0xsj/numio/internal/token"
	"github.com/0xsj/numio/pkg/types"
)

// Highlighter applies syntax highlighting to numio expressions.
type Highlighter struct {
	theme *Theme
}

// New creates a new Highlighter with the given theme.
func New(theme *Theme) *Highlighter {
	if theme == nil {
		theme = DefaultTheme()
	}
	return &Highlighter{theme: theme}
}

// NewWithThemeName creates a Highlighter with a theme by name.
func NewWithThemeName(name string) *Highlighter {
	return New(GetTheme(name))
}

// Default creates a Highlighter with the default theme.
func Default() *Highlighter {
	return New(DefaultTheme())
}

// Theme returns the current theme.
func (h *Highlighter) Theme() *Theme {
	return h.theme
}

// SetTheme changes the highlighter's theme.
func (h *Highlighter) SetTheme(theme *Theme) {
	if theme != nil {
		h.theme = theme
	}
}

// ════════════════════════════════════════════════════════════════
// HIGHLIGHTING
// ════════════════════════════════════════════════════════════════

// Highlight applies syntax highlighting to an input string.
// Returns the highlighted string with ANSI color codes.
func (h *Highlighter) Highlight(input string) string {
	if input == "" {
		return ""
	}

	// Check for comment-only lines first
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return h.theme.Render(ClassComment, input)
	}

	// Tokenize the input
	tokens := lexer.Tokenize(input)

	// Build highlighted output
	var result strings.Builder
	lastEnd := 0

	for _, tok := range tokens {
		if tok.Type == token.EOF {
			break
		}

		// Preserve whitespace/characters between tokens
		if tok.Pos > lastEnd {
			result.WriteString(input[lastEnd:tok.Pos])
		}

		// Get token class and apply highlighting
		class := h.classifyToken(tok)
		result.WriteString(h.theme.Render(class, tok.Literal))

		lastEnd = tok.Pos + len(tok.Literal)
	}

	// Append any remaining content
	if lastEnd < len(input) {
		result.WriteString(input[lastEnd:])
	}

	return result.String()
}

// HighlightLine highlights a line, handling the cursor position.
// Returns highlighted text before cursor, cursor char (unstyled), and after cursor.
func (h *Highlighter) HighlightLine(input string, cursorPos int) (before, cursor, after string) {
	if input == "" {
		return "", "", ""
	}

	// Clamp cursor position
	if cursorPos < 0 {
		cursorPos = 0
	}
	if cursorPos > len(input) {
		cursorPos = len(input)
	}

	// If cursor is at end, highlight entire line
	if cursorPos >= len(input) {
		return h.Highlight(input), "", ""
	}

	// Split at cursor and highlight each part
	beforeText := input[:cursorPos]
	cursorChar := string(input[cursorPos])
	afterText := ""
	if cursorPos+1 < len(input) {
		afterText = input[cursorPos+1:]
	}

	return h.Highlight(beforeText), cursorChar, h.Highlight(afterText)
}

// ════════════════════════════════════════════════════════════════
// TOKEN CLASSIFICATION
// ════════════════════════════════════════════════════════════════

// classifyToken determines the TokenClass for a given token.
func (h *Highlighter) classifyToken(tok token.Token) TokenClass {
	switch tok.Type {
	// Numbers and percentages
	case token.NUMBER:
		return ClassNumber
	case token.PERCENT:
		return ClassPercent

	// Operators
	case token.PLUS, token.MINUS, token.STAR, token.SLASH, token.CARET, token.POWER:
		return ClassOperator

	// Parentheses
	case token.LPAREN, token.RPAREN:
		return ClassParen

	// Assignment
	case token.EQUALS:
		return ClassAssign

	// Keywords
	case token.IN, token.OF:
		return ClassKeyword

	// Currency symbols
	case token.DOLLAR, token.EURO, token.POUND, token.YEN, token.BITCOIN, token.CURRENCY:
		return ClassCurrency

	// Comments
	case token.COMMENT:
		return ClassComment

	// Identifiers - need further classification
	case token.IDENTIFIER:
		return h.classifyIdentifier(tok.Literal)

	// Comma
	case token.COMMA:
		return ClassOperator

	default:
		return ClassNone
	}
}

// classifyIdentifier determines if an identifier is a function, currency, unit, etc.
func (h *Highlighter) classifyIdentifier(name string) TokenClass {
	lower := strings.ToLower(name)

	// Check if it's a known function
	if isFunction(lower) {
		return ClassFunction
	}

	// Check if it's a currency code or name
	if types.ParseCurrency(name) != nil {
		return ClassCurrency
	}

	// Check if it's a crypto code or name
	if types.ParseCrypto(name) != nil {
		return ClassCrypto
	}

	// Check if it's a metal code or name
	if types.ParseMetal(name) != nil {
		return ClassMetal
	}

	// Check if it's a unit code or name
	if types.ParseUnit(name) != nil {
		return ClassUnit
	}

	// Default to identifier (variable)
	return ClassIdentifier
}

// isFunction checks if a name is a known built-in function.
func isFunction(name string) bool {
	functions := map[string]bool{
		// Aggregation
		"sum":     true,
		"avg":     true,
		"average": true,
		"mean":    true,
		"min":     true,
		"max":     true,
		"count":   true,

		// Math
		"abs":   true,
		"sqrt":  true,
		"round": true,
		"floor": true,
		"ceil":  true,
		"log":   true,
		"log10": true,
		"ln":    true,
		"exp":   true,
		"pow":   true,

		// Trigonometry
		"sin":  true,
		"cos":  true,
		"tan":  true,
		"asin": true,
		"acos": true,
		"atan": true,
	}

	return functions[name]
}

// ════════════════════════════════════════════════════════════════
// SPAN-BASED HIGHLIGHTING (for more control)
// ════════════════════════════════════════════════════════════════

// Span represents a highlighted segment of text.
type Span struct {
	Start int        // Start position in original string
	End   int        // End position (exclusive)
	Text  string     // The text content
	Class TokenClass // The token class for coloring
}

// HighlightSpans returns highlighting information as spans.
// Useful for custom rendering or editors that need position info.
func (h *Highlighter) HighlightSpans(input string) []Span {
	if input == "" {
		return nil
	}

	// Check for comment-only lines
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return []Span{{
			Start: 0,
			End:   len(input),
			Text:  input,
			Class: ClassComment,
		}}
	}

	tokens := lexer.Tokenize(input)
	spans := make([]Span, 0, len(tokens))
	lastEnd := 0

	for _, tok := range tokens {
		if tok.Type == token.EOF {
			break
		}

		// Add span for whitespace/gaps (as ClassNone)
		if tok.Pos > lastEnd {
			spans = append(spans, Span{
				Start: lastEnd,
				End:   tok.Pos,
				Text:  input[lastEnd:tok.Pos],
				Class: ClassNone,
			})
		}

		// Add span for this token
		class := h.classifyToken(tok)
		end := tok.Pos + len(tok.Literal)
		spans = append(spans, Span{
			Start: tok.Pos,
			End:   end,
			Text:  tok.Literal,
			Class: class,
		})

		lastEnd = end
	}

	// Add any remaining content
	if lastEnd < len(input) {
		spans = append(spans, Span{
			Start: lastEnd,
			End:   len(input),
			Text:  input[lastEnd:],
			Class: ClassNone,
		})
	}

	return spans
}

// RenderSpans renders spans with the current theme.
func (h *Highlighter) RenderSpans(spans []Span) string {
	var result strings.Builder
	for _, span := range spans {
		result.WriteString(h.theme.Render(span.Class, span.Text))
	}
	return result.String()
}

// ════════════════════════════════════════════════════════════════
// UTILITIES
// ════════════════════════════════════════════════════════════════

// StripHighlighting removes ANSI color codes from a string.
func StripHighlighting(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}

// VisibleLength returns the visible length of a string (excluding ANSI codes).
func VisibleLength(s string) int {
	return len(StripHighlighting(s))
}
