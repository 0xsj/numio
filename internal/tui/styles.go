// internal/tui/styles.go

package tui

import "github.com/charmbracelet/lipgloss"

// Color palette (matching numr's dark theme)
var (
	// Base colors
	ColorBg        = lipgloss.Color("#1a1a2e")
	ColorFg        = lipgloss.Color("#eaeaea")
	ColorFgDim     = lipgloss.Color("#6b6b6b")
	ColorBorder    = lipgloss.Color("#3a3a5a")
	ColorHighlight = lipgloss.Color("#4a4a7a")

	// Syntax colors
	ColorNumber     = lipgloss.Color("#79c0ff") // Cyan/blue for numbers
	ColorCurrency   = lipgloss.Color("#7ee787") // Green for currency symbols
	ColorUnit       = lipgloss.Color("#a5d6ff") // Light blue for units
	ColorOperator   = lipgloss.Color("#ffffff") // White for operators
	ColorComment    = lipgloss.Color("#6b6b6b") // Gray for comments
	ColorVariable   = lipgloss.Color("#ffa657") // Orange for variables
	ColorKeyword    = lipgloss.Color("#ff7b72") // Red for keywords
	ColorPercentage = lipgloss.Color("#d2a8ff") // Purple for percentages
	ColorString     = lipgloss.Color("#a5d6ff") // Light blue for strings
	ColorError      = lipgloss.Color("#f85149") // Red for errors

	// Result colors
	ColorResultPositive = lipgloss.Color("#7ee787") // Green for positive
	ColorResultNegative = lipgloss.Color("#f85149") // Red for negative
	ColorResultNeutral  = lipgloss.Color("#79c0ff") // Blue for neutral

	// UI colors
	ColorStatusBg     = lipgloss.Color("#2d2d4a")
	ColorStatusFg     = lipgloss.Color("#eaeaea")
	ColorModeNormal   = lipgloss.Color("#79c0ff") // Blue for NORMAL
	ColorModeInsert   = lipgloss.Color("#7ee787") // Green for INSERT
	ColorModeVisual   = lipgloss.Color("#d2a8ff") // Purple for VISUAL
	ColorHelpBg       = lipgloss.Color("#2d2d4a")
	ColorHelpBorder   = lipgloss.Color("#4a4a7a")
	ColorHelpHeader   = lipgloss.Color("#79c0ff")
	ColorCursor       = lipgloss.Color("#79c0ff")
	ColorLineNumber   = lipgloss.Color("#4a4a6a")
	ColorCurrentLine  = lipgloss.Color("#2a2a4a")
)

// Styles defines all the lipgloss styles for the TUI.
type Styles struct {
	// Layout
	App        lipgloss.Style
	EditorPane lipgloss.Style
	ResultPane lipgloss.Style
	StatusBar  lipgloss.Style
	HelpModal  lipgloss.Style

	// Editor
	LineNumber    lipgloss.Style
	CurrentLine   lipgloss.Style
	Cursor        lipgloss.Style
	EditorContent lipgloss.Style

	// Results
	ResultLine     lipgloss.Style
	ResultPositive lipgloss.Style
	ResultNegative lipgloss.Style
	ResultNeutral  lipgloss.Style
	ResultError    lipgloss.Style

	// Syntax highlighting
	Number     lipgloss.Style
	Currency   lipgloss.Style
	Unit       lipgloss.Style
	Operator   lipgloss.Style
	Comment    lipgloss.Style
	Variable   lipgloss.Style
	Keyword    lipgloss.Style
	Percentage lipgloss.Style

	// Status bar
	ModeNormal lipgloss.Style
	ModeInsert lipgloss.Style
	ModeVisual lipgloss.Style
	StatusText lipgloss.Style
	StatusHint lipgloss.Style
	Total      lipgloss.Style

	// Help
	HelpTitle   lipgloss.Style
	HelpKey     lipgloss.Style
	HelpDesc    lipgloss.Style
	HelpSection lipgloss.Style
}

// DefaultStyles returns the default style configuration.
func DefaultStyles() Styles {
	return Styles{
		// Layout styles
		App: lipgloss.NewStyle().
			Background(ColorBg).
			Foreground(ColorFg),

		EditorPane: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1),

		ResultPane: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1).
			Align(lipgloss.Right),

		StatusBar: lipgloss.NewStyle().
			Background(ColorStatusBg).
			Foreground(ColorStatusFg).
			Padding(0, 1),

		HelpModal: lipgloss.NewStyle().
			Background(ColorHelpBg).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorHelpBorder).
			Padding(1, 2),

		// Editor styles
		LineNumber: lipgloss.NewStyle().
			Foreground(ColorLineNumber).
			Width(4).
			Align(lipgloss.Right).
			MarginRight(1),

		CurrentLine: lipgloss.NewStyle().
			Background(ColorCurrentLine),

		Cursor: lipgloss.NewStyle().
			Background(ColorCursor).
			Foreground(ColorBg),

		EditorContent: lipgloss.NewStyle().
			Foreground(ColorFg),

		// Result styles
		ResultLine: lipgloss.NewStyle().
			Foreground(ColorFg).
			Align(lipgloss.Right),

		ResultPositive: lipgloss.NewStyle().
			Foreground(ColorResultPositive).
			Align(lipgloss.Right),

		ResultNegative: lipgloss.NewStyle().
			Foreground(ColorResultNegative).
			Align(lipgloss.Right),

		ResultNeutral: lipgloss.NewStyle().
			Foreground(ColorResultNeutral).
			Align(lipgloss.Right),

		ResultError: lipgloss.NewStyle().
			Foreground(ColorError).
			Align(lipgloss.Right),

		// Syntax styles
		Number: lipgloss.NewStyle().
			Foreground(ColorNumber),

		Currency: lipgloss.NewStyle().
			Foreground(ColorCurrency),

		Unit: lipgloss.NewStyle().
			Foreground(ColorUnit),

		Operator: lipgloss.NewStyle().
			Foreground(ColorOperator),

		Comment: lipgloss.NewStyle().
			Foreground(ColorComment).
			Italic(true),

		Variable: lipgloss.NewStyle().
			Foreground(ColorVariable),

		Keyword: lipgloss.NewStyle().
			Foreground(ColorKeyword),

		Percentage: lipgloss.NewStyle().
			Foreground(ColorPercentage),

		// Status bar styles
		ModeNormal: lipgloss.NewStyle().
			Background(ColorModeNormal).
			Foreground(ColorBg).
			Bold(true).
			Padding(0, 1),

		ModeInsert: lipgloss.NewStyle().
			Background(ColorModeInsert).
			Foreground(ColorBg).
			Bold(true).
			Padding(0, 1),

		ModeVisual: lipgloss.NewStyle().
			Background(ColorModeVisual).
			Foreground(ColorBg).
			Bold(true).
			Padding(0, 1),

		StatusText: lipgloss.NewStyle().
			Foreground(ColorStatusFg),

		StatusHint: lipgloss.NewStyle().
			Foreground(ColorFgDim),

		Total: lipgloss.NewStyle().
			Foreground(ColorResultPositive).
			Bold(true),

		// Help styles
		HelpTitle: lipgloss.NewStyle().
			Foreground(ColorHelpHeader).
			Bold(true).
			MarginBottom(1),

		HelpKey: lipgloss.NewStyle().
			Foreground(ColorNumber).
			Width(16),

		HelpDesc: lipgloss.NewStyle().
			Foreground(ColorFgDim),

		HelpSection: lipgloss.NewStyle().
			Foreground(ColorVariable).
			Bold(true).
			MarginTop(1).
			MarginBottom(0),
	}
}

// WithDimensions returns styles adjusted for the given dimensions.
func (s Styles) WithDimensions(width, height int) Styles {
	editorWidth := width * 2 / 3
	resultWidth := width - editorWidth - 4 // Account for borders

	s.EditorPane = s.EditorPane.Width(editorWidth).Height(height - 3)
	s.ResultPane = s.ResultPane.Width(resultWidth).Height(height - 3)
	s.StatusBar = s.StatusBar.Width(width)

	return s
}