// internal/highlight/color.go

// Package highlight provides syntax highlighting for numio expressions.
package highlight

import "github.com/charmbracelet/lipgloss"

// Color represents a color that can be applied to text.
// It wraps lipgloss.Color for consistency with the TUI.
type Color struct {
	value lipgloss.Color
}

// NewColor creates a Color from a hex string (e.g., "#ff0000") or ANSI code (e.g., "196").
func NewColor(value string) Color {
	return Color{value: lipgloss.Color(value)}
}

// Lipgloss returns the underlying lipgloss.Color.
func (c Color) Lipgloss() lipgloss.Color {
	return c.value
}

// Style returns a lipgloss.Style with this color as the foreground.
func (c Color) Style() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c.value)
}

// String returns the color value as a string.
func (c Color) String() string {
	return string(c.value)
}

// ════════════════════════════════════════════════════════════════
// SEMANTIC COLORS (what things mean)
// ════════════════════════════════════════════════════════════════

// TokenClass represents a semantic category for syntax highlighting.
type TokenClass int

const (
	ClassNone       TokenClass = iota // No special highlighting
	ClassNumber                       // Numeric literals: 42, 3.14
	ClassPercent                      // Percentages: 20%
	ClassOperator                     // Operators: +, -, *, /, ^
	ClassParen                        // Parentheses: (, )
	ClassIdentifier                   // Variable names
	ClassKeyword                      // Keywords: in, to, of
	ClassFunction                     // Function names: sum, avg
	ClassCurrency                     // Currency symbols and codes: $, €, USD
	ClassUnit                         // Unit codes: km, lb, hours
	ClassCrypto                       // Crypto codes: BTC, ETH
	ClassMetal                        // Metal codes: XAU, XAG
	ClassComment                      // Comments: # or //
	ClassError                        // Errors
	ClassAssign                       // Assignment: =
)

// String returns the token class name.
func (c TokenClass) String() string {
	switch c {
	case ClassNone:
		return "none"
	case ClassNumber:
		return "number"
	case ClassPercent:
		return "percent"
	case ClassOperator:
		return "operator"
	case ClassParen:
		return "paren"
	case ClassIdentifier:
		return "identifier"
	case ClassKeyword:
		return "keyword"
	case ClassFunction:
		return "function"
	case ClassCurrency:
		return "currency"
	case ClassUnit:
		return "unit"
	case ClassCrypto:
		return "crypto"
	case ClassMetal:
		return "metal"
	case ClassComment:
		return "comment"
	case ClassError:
		return "error"
	case ClassAssign:
		return "assign"
	default:
		return "unknown"
	}
}

// ════════════════════════════════════════════════════════════════
// COLOR PALETTE (common colors for themes)
// ════════════════════════════════════════════════════════════════

// Palette provides a set of commonly used colors.
var Palette = struct {
	// Grayscale
	White   Color
	Gray100 Color
	Gray200 Color
	Gray300 Color
	Gray400 Color
	Gray500 Color
	Gray600 Color
	Gray700 Color
	Gray800 Color
	Black   Color

	// Colors
	Red     Color
	Orange  Color
	Yellow  Color
	Green   Color
	Cyan    Color
	Blue    Color
	Purple  Color
	Pink    Color
	Magenta Color

	// Semantic
	Success Color
	Warning Color
	Error   Color
	Info    Color
}{
	// Grayscale
	White:   NewColor("#ffffff"),
	Gray100: NewColor("#f5f5f5"),
	Gray200: NewColor("#eeeeee"),
	Gray300: NewColor("#e0e0e0"),
	Gray400: NewColor("#bdbdbd"),
	Gray500: NewColor("#9e9e9e"),
	Gray600: NewColor("#757575"),
	Gray700: NewColor("#616161"),
	Gray800: NewColor("#424242"),
	Black:   NewColor("#000000"),

	// Colors
	Red:     NewColor("#f85149"),
	Orange:  NewColor("#ffa657"),
	Yellow:  NewColor("#e3b341"),
	Green:   NewColor("#7ee787"),
	Cyan:    NewColor("#56d4dd"),
	Blue:    NewColor("#79c0ff"),
	Purple:  NewColor("#d2a8ff"),
	Pink:    NewColor("#ff7eb6"),
	Magenta: NewColor("#bc8cff"),

	// Semantic
	Success: NewColor("#7ee787"),
	Warning: NewColor("#ffa657"),
	Error:   NewColor("#f85149"),
	Info:    NewColor("#79c0ff"),
}

// ════════════════════════════════════════════════════════════════
// STYLE HELPERS
// ════════════════════════════════════════════════════════════════

// Render applies a color to text and returns the styled string.
func (c Color) Render(text string) string {
	return c.Style().Render(text)
}

// Bold returns a style with this color and bold enabled.
func (c Color) Bold() lipgloss.Style {
	return c.Style().Bold(true)
}

// Italic returns a style with this color and italic enabled.
func (c Color) Italic() lipgloss.Style {
	return c.Style().Italic(true)
}

// Underline returns a style with this color and underline enabled.
func (c Color) Underline() lipgloss.Style {
	return c.Style().Underline(true)
}
