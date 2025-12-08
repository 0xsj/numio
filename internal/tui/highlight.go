// internal/tui/highlight.go

package tui

import (
	"strings"

	"github.com/0xsj/numio/internal/lexer"
	"github.com/0xsj/numio/internal/token"
	"github.com/charmbracelet/lipgloss"
)

// Highlighter provides syntax highlighting for numio expressions.
type Highlighter struct {
	styles Styles
}

// NewHighlighter creates a new syntax highlighter.
func NewHighlighter(styles Styles) *Highlighter {
	return &Highlighter{styles: styles}
}

// Highlight returns a syntax-highlighted version of the input line.
func (h *Highlighter) Highlight(line string) string {
	if line == "" {
		return ""
	}

	// Check for comment-only line
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return h.styles.Comment.Render(line)
	}

	// Tokenize the line
	tokens := lexer.Tokenize(line)

	// Build highlighted output
	var result strings.Builder
	lastPos := 0

	for _, tok := range tokens {
		// Add any whitespace/characters between tokens
		if tok.Pos > lastPos {
			result.WriteString(line[lastPos:tok.Pos])
		}

		// Highlight token based on type
		highlighted := h.highlightToken(tok)
		result.WriteString(highlighted)

		// Update position
		lastPos = tok.Pos + len(tok.Literal)
	}

	// Add any remaining characters
	if lastPos < len(line) {
		result.WriteString(line[lastPos:])
	}

	return result.String()
}

// highlightToken returns the highlighted version of a single token.
func (h *Highlighter) highlightToken(tok token.Token) string {
	switch tok.Type {
	// Numbers
	case token.NUMBER:
		return h.styles.Number.Render(tok.Literal)

	// Percentages
	case token.PERCENT:
		return h.styles.Percentage.Render(tok.Literal)

	// Currency symbols
	case token.DOLLAR, token.EURO, token.POUND, token.YEN, token.BITCOIN, token.CURRENCY:
		return h.styles.Currency.Render(tok.Literal)

	// Operators
	case token.PLUS, token.MINUS, token.STAR, token.SLASH, token.CARET, token.POWER, token.EQUALS:
		return h.styles.Operator.Render(tok.Literal)

	// Parentheses and punctuation
	case token.LPAREN, token.RPAREN, token.COMMA:
		return h.styles.Operator.Render(tok.Literal)

	// Keywords
	case token.IN, token.OF:
		return h.styles.Keyword.Render(tok.Literal)

	// Comments
	case token.COMMENT:
		return h.styles.Comment.Render(tok.Literal)

	// Identifiers (variables, units, currencies)
	case token.IDENTIFIER:
		return h.highlightIdentifier(tok.Literal)

	// Default
	default:
		return tok.Literal
	}
}

// highlightIdentifier determines the appropriate style for an identifier.
func (h *Highlighter) highlightIdentifier(literal string) string {
	lower := strings.ToLower(literal)

	// Check for known units
	if isKnownUnit(lower) {
		return h.styles.Unit.Render(literal)
	}

	// Check for known currencies
	if isKnownCurrency(lower) {
		return h.styles.Currency.Render(literal)
	}

	// Check for known functions
	if isKnownFunction(lower) {
		return h.styles.Keyword.Render(literal)
	}

	// Default to variable style
	return h.styles.Variable.Render(literal)
}

// isKnownUnit checks if the identifier is a known unit.
func isKnownUnit(s string) bool {
	units := map[string]bool{
		// Length
		"m": true, "meter": true, "meters": true, "metre": true, "metres": true,
		"km": true, "kilometer": true, "kilometers": true, "kilometre": true, "kilometres": true,
		"cm": true, "centimeter": true, "centimeters": true,
		"mm": true, "millimeter": true, "millimeters": true,
		"mi": true, "mile": true, "miles": true,
		"yd": true, "yard": true, "yards": true,
		"ft": true, "foot": true, "feet": true,
		"in": true, "inch": true, "inches": true,
		// Weight
		"kg": true, "kilogram": true, "kilograms": true,
		"g": true, "gram": true, "grams": true,
		"mg": true, "milligram": true, "milligrams": true,
		"lb": true, "pound": true, "pounds": true,
		"oz": true, "ounce": true, "ounces": true,
		// Time
		"s": true, "sec": true, "second": true, "seconds": true,
		"min": true, "minute": true, "minutes": true,
		"h": true, "hr": true, "hour": true, "hours": true,
		"d": true, "day": true, "days": true,
		"wk": true, "week": true, "weeks": true,
		"mo": true, "month": true, "months": true,
		"yr": true, "year": true, "years": true,
		// Temperature
		"c": true, "celsius": true,
		"f": true, "fahrenheit": true,
		"k": true, "kelvin": true,
		// Data
		"b": true, "byte": true, "bytes": true,
		"kb": true, "kilobyte": true, "kilobytes": true,
		"mb": true, "megabyte": true, "megabytes": true,
		"gb": true, "gigabyte": true, "gigabytes": true,
		"tb": true, "terabyte": true, "terabytes": true,
	}
	return units[s]
}

// isKnownCurrency checks if the identifier is a known currency.
func isKnownCurrency(s string) bool {
	currencies := map[string]bool{
		// Fiat
		"usd": true, "dollar": true, "dollars": true,
		"eur": true, "euro": true, "euros": true,
		"gbp": true, "pound": true, "pounds": true, "sterling": true,
		"jpy": true, "yen": true,
		"cny": true, "yuan": true, "rmb": true,
		"inr": true, "rupee": true, "rupees": true,
		"rub": true, "ruble": true, "rubles": true,
		"krw": true, "won": true,
		"try": true, "lira": true,
		"chf": true, "franc": true, "francs": true,
		"cad": true, "aud": true, "nzd": true,
		"mxn": true, "peso": true, "pesos": true,
		"brl": true, "real": true,
		"zar": true, "rand": true,
		// Crypto
		"btc": true, "bitcoin": true,
		"eth": true, "ethereum": true, "ether": true,
		"sol": true, "solana": true,
		"bnb": true,
		"xrp": true, "ripple": true,
		"ada": true, "cardano": true,
		"doge": true, "dogecoin": true,
		"dot": true, "polkadot": true,
		"matic": true, "polygon": true,
		"ltc": true, "litecoin": true,
		"usdt": true, "tether": true,
		"usdc": true,
		// Metals
		"xau": true, "gold": true,
		"xag": true, "silver": true,
		"xpt": true, "platinum": true,
		"xpd": true, "palladium": true,
	}
	return currencies[s]
}

// isKnownFunction checks if the identifier is a known function.
func isKnownFunction(s string) bool {
	functions := map[string]bool{
		"sum": true, "avg": true, "average": true, "mean": true,
		"min": true, "max": true, "count": true,
		"abs": true, "sqrt": true, "round": true, "floor": true, "ceil": true,
		"log": true, "log10": true, "ln": true, "exp": true,
		"sin": true, "cos": true, "tan": true,
		"asin": true, "acos": true, "atan": true,
		"pow": true,
	}
	return functions[s]
}

// HighlightResult returns a styled result value.
func (h *Highlighter) HighlightResult(value string, isError bool, isNegative bool) string {
	if isError {
		return h.styles.ResultError.Render(value)
	}
	if isNegative {
		return h.styles.ResultNegative.Render(value)
	}
	return h.styles.ResultPositive.Render(value)
}

// HighlightLineNumber returns a styled line number.
func (h *Highlighter) HighlightLineNumber(num int, isCurrent bool) string {
	style := h.styles.LineNumber
	if isCurrent {
		style = style.Foreground(lipgloss.Color("#ffffff"))
	}
	return style.Render(formatLineNumber(num))
}

// formatLineNumber formats a line number with padding.
func formatLineNumber(num int) string {
	if num < 10 {
		return "  " + itoa(num)
	}
	if num < 100 {
		return " " + itoa(num)
	}
	if num < 1000 {
		return itoa(num)
	}
	return itoa(num)
}

// itoa converts int to string without fmt.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [10]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}