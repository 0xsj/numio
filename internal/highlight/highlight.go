// internal/highlight/highlight.go

package highlight

import (
	"strings"

	"github.com/0xsj/numio/internal/eval"
	"github.com/0xsj/numio/internal/lexer"
	"github.com/0xsj/numio/internal/token"
	"github.com/0xsj/numio/pkg/types"
)

// Highlighter applies syntax highlighting to numio expressions.
type Highlighter struct {
	theme *Theme

	// Optional: user function registry for highlighting user-defined functions
	userFuncs *eval.UserFuncRegistry
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

// SetUserFuncRegistry sets the user function registry for highlighting user-defined functions.
func (h *Highlighter) SetUserFuncRegistry(registry *eval.UserFuncRegistry) {
	h.userFuncs = registry
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

	// Track context for function definition highlighting
	inFuncDef := false
	expectFuncName := false
	expectParams := false
	parenDepth := 0

	for i, tok := range tokens {
		if tok.Type == token.EOF {
			break
		}

		// Preserve whitespace/characters between tokens
		if tok.Pos > lastEnd {
			result.WriteString(input[lastEnd:tok.Pos])
		}

		// Update context tracking
		if tok.Type == token.DEF {
			inFuncDef = true
			expectFuncName = true
		}

		// Get token class and apply highlighting
		class := h.classifyTokenWithContext(tok, tokens, i, inFuncDef, expectFuncName, expectParams, parenDepth)
		result.WriteString(h.theme.Render(class, tok.Literal))

		// Update context after processing
		if expectFuncName && tok.Type == token.IDENTIFIER {
			expectFuncName = false
			expectParams = true
		}
		if tok.Type == token.LPAREN && expectParams {
			parenDepth++
		}
		if tok.Type == token.RPAREN && expectParams {
			parenDepth--
			if parenDepth == 0 {
				expectParams = false
			}
		}
		if tok.Type == token.COLON && inFuncDef {
			inFuncDef = false
		}

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
	return h.classifyTokenWithContext(tok, nil, 0, false, false, false, 0)
}

// classifyTokenWithContext classifies a token with function definition context.
func (h *Highlighter) classifyTokenWithContext(tok token.Token, tokens []token.Token, idx int, inFuncDef, expectFuncName, expectParams bool, parenDepth int) TokenClass {
	switch tok.Type {
	// Numbers and percentages
	case token.NUMBER:
		return ClassNumber
	case token.PERCENT:
		return ClassPercent

	case token.STRING:
		return ClassString

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

	// Function definition keyword
	case token.DEF:
		return ClassKeyword

	// Colon (used in function definitions)
	case token.COLON:
		return ClassOperator

	// Currency symbols
	case token.DOLLAR, token.EURO, token.POUND, token.YEN, token.BITCOIN, token.CURRENCY:
		return ClassCurrency

	// Comments
	case token.COMMENT:
		return ClassComment

	// Identifiers - need further classification
	case token.IDENTIFIER:
		// In function definition context
		if expectFuncName {
			return ClassFuncName
		}
		if expectParams && parenDepth > 0 {
			return ClassParam
		}
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

	// Check if it's a known built-in function
	if isBuiltinFunction(lower) {
		return ClassFunction
	}

	// Check if it's a user-defined function
	if h.userFuncs != nil && h.userFuncs.Has(name) {
		return ClassUserFunc
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

// isBuiltinFunction checks if a name is a known built-in function.
func isBuiltinFunction(name string) bool {
	// Use the eval package's HasFunction for accuracy
	return eval.HasFunction(name)
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

	// Track context for function definition highlighting
	inFuncDef := false
	expectFuncName := false
	expectParams := false
	parenDepth := 0

	for i, tok := range tokens {
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

		// Update context tracking
		if tok.Type == token.DEF {
			inFuncDef = true
			expectFuncName = true
		}

		// Add span for this token
		class := h.classifyTokenWithContext(tok, tokens, i, inFuncDef, expectFuncName, expectParams, parenDepth)
		end := tok.Pos + len(tok.Literal)
		spans = append(spans, Span{
			Start: tok.Pos,
			End:   end,
			Text:  tok.Literal,
			Class: class,
		})

		// Update context after processing
		if expectFuncName && tok.Type == token.IDENTIFIER {
			expectFuncName = false
			expectParams = true
		}
		if tok.Type == token.LPAREN && expectParams {
			parenDepth++
		}
		if tok.Type == token.RPAREN && expectParams {
			parenDepth--
			if parenDepth == 0 {
				expectParams = false
			}
		}
		if tok.Type == token.COLON && inFuncDef {
			inFuncDef = false
		}

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
