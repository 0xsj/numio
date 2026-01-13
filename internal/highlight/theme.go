// internal/highlight/theme.go

package highlight

import "github.com/charmbracelet/lipgloss"

// Theme defines colors for each token class.
type Theme struct {
	Name   string
	Colors map[TokenClass]Color
}

// Get returns the color for a token class, or a default if not found.
func (t *Theme) Get(class TokenClass) Color {
	if c, ok := t.Colors[class]; ok {
		return c
	}
	return Palette.White
}

// Style returns a lipgloss.Style for a token class.
func (t *Theme) Style(class TokenClass) lipgloss.Style {
	return t.Get(class).Style()
}

// Render applies the theme color to text for a given token class.
func (t *Theme) Render(class TokenClass, text string) string {
	return t.Get(class).Render(text)
}

// ════════════════════════════════════════════════════════════════
// BUILT-IN THEMES
// ════════════════════════════════════════════════════════════════

// DefaultTheme returns the default dark theme.
func DefaultTheme() *Theme {
	return &Theme{
		Name: "default",
		Colors: map[TokenClass]Color{
			ClassNone:       Palette.White,
			ClassNumber:     Palette.Purple,
			ClassString:     Palette.Green, // ADD THIS
			ClassPercent:    Palette.Magenta,
			ClassOperator:   Palette.Cyan,
			ClassParen:      Palette.Gray400,
			ClassIdentifier: Palette.White,
			ClassKeyword:    Palette.Orange,
			ClassFunction:   Palette.Blue,
			ClassCurrency:   Palette.Green,
			ClassUnit:       Palette.Yellow,
			ClassCrypto:     Palette.Orange,
			ClassMetal:      Palette.Yellow,
			ClassComment:    Palette.Gray600,
			ClassError:      Palette.Error,
			ClassAssign:     Palette.Cyan,
		},
	}
}

// DraculaTheme returns a Dracula-inspired theme.
func DraculaTheme() *Theme {
	return &Theme{
		Name: "dracula",
		Colors: map[TokenClass]Color{
			ClassNone:       NewColor("#f8f8f2"),
			ClassNumber:     NewColor("#bd93f9"),
			ClassString:     NewColor("#f1fa8c"), // ADD THIS - Yellow
			ClassPercent:    NewColor("#ff79c6"),
			ClassOperator:   NewColor("#ff79c6"),
			ClassParen:      NewColor("#f8f8f2"),
			ClassIdentifier: NewColor("#f8f8f2"),
			ClassKeyword:    NewColor("#ff79c6"),
			ClassFunction:   NewColor("#50fa7b"),
			ClassCurrency:   NewColor("#50fa7b"),
			ClassUnit:       NewColor("#f1fa8c"),
			ClassCrypto:     NewColor("#ffb86c"),
			ClassMetal:      NewColor("#f1fa8c"),
			ClassComment:    NewColor("#6272a4"),
			ClassError:      NewColor("#ff5555"),
			ClassAssign:     NewColor("#ff79c6"),
		},
	}
}

// MonokaiTheme returns a Monokai-inspired theme.
func MonokaiTheme() *Theme {
	return &Theme{
		Name: "monokai",
		Colors: map[TokenClass]Color{
			ClassNone:       NewColor("#f8f8f2"),
			ClassNumber:     NewColor("#ae81ff"),
			ClassString:     NewColor("#e6db74"), // ADD THIS - Yellow
			ClassPercent:    NewColor("#ae81ff"),
			ClassOperator:   NewColor("#f92672"),
			ClassParen:      NewColor("#f8f8f2"),
			ClassIdentifier: NewColor("#f8f8f2"),
			ClassKeyword:    NewColor("#f92672"),
			ClassFunction:   NewColor("#66d9ef"),
			ClassCurrency:   NewColor("#a6e22e"),
			ClassUnit:       NewColor("#e6db74"),
			ClassCrypto:     NewColor("#fd971f"),
			ClassMetal:      NewColor("#e6db74"),
			ClassComment:    NewColor("#75715e"),
			ClassError:      NewColor("#f92672"),
			ClassAssign:     NewColor("#f92672"),
		},
	}
}

// GruvboxTheme returns a Gruvbox dark theme.
func GruvboxTheme() *Theme {
	return &Theme{
		Name: "gruvbox",
		Colors: map[TokenClass]Color{
			ClassNone:       NewColor("#ebdbb2"),
			ClassNumber:     NewColor("#d3869b"),
			ClassString:     NewColor("#b8bb26"), // ADD THIS - Green
			ClassPercent:    NewColor("#d3869b"),
			ClassOperator:   NewColor("#8ec07c"),
			ClassParen:      NewColor("#a89984"),
			ClassIdentifier: NewColor("#ebdbb2"),
			ClassKeyword:    NewColor("#fb4934"),
			ClassFunction:   NewColor("#83a598"),
			ClassCurrency:   NewColor("#b8bb26"),
			ClassUnit:       NewColor("#fabd2f"),
			ClassCrypto:     NewColor("#fe8019"),
			ClassMetal:      NewColor("#fabd2f"),
			ClassComment:    NewColor("#928374"),
			ClassError:      NewColor("#fb4934"),
			ClassAssign:     NewColor("#8ec07c"),
		},
	}
}

// LightTheme returns a light theme suitable for light terminals.
func LightTheme() *Theme {
	return &Theme{
		Name: "light",
		Colors: map[TokenClass]Color{
			ClassNone:       NewColor("#24292e"),
			ClassNumber:     NewColor("#6f42c1"),
			ClassString:     NewColor("#032f62"), // ADD THIS - Dark blue
			ClassPercent:    NewColor("#6f42c1"),
			ClassOperator:   NewColor("#d73a49"),
			ClassParen:      NewColor("#24292e"),
			ClassIdentifier: NewColor("#24292e"),
			ClassKeyword:    NewColor("#d73a49"),
			ClassFunction:   NewColor("#005cc5"),
			ClassCurrency:   NewColor("#22863a"),
			ClassUnit:       NewColor("#b08800"),
			ClassCrypto:     NewColor("#e36209"),
			ClassMetal:      NewColor("#b08800"),
			ClassComment:    NewColor("#6a737d"),
			ClassError:      NewColor("#cb2431"),
			ClassAssign:     NewColor("#d73a49"),
		},
	}
}

// ════════════════════════════════════════════════════════════════
// THEME REGISTRY
// ════════════════════════════════════════════════════════════════

// builtinThemes holds all registered themes.
var builtinThemes = map[string]func() *Theme{
	"default": DefaultTheme,
	"dracula": DraculaTheme,
	"monokai": MonokaiTheme,
	"gruvbox": GruvboxTheme,
	"light":   LightTheme,
}

// GetTheme returns a theme by name, or the default theme if not found.
func GetTheme(name string) *Theme {
	if fn, ok := builtinThemes[name]; ok {
		return fn()
	}
	return DefaultTheme()
}

// ThemeNames returns a list of all available theme names.
func ThemeNames() []string {
	names := make([]string, 0, len(builtinThemes))
	for name := range builtinThemes {
		names = append(names, name)
	}
	return names
}

// ════════════════════════════════════════════════════════════════
// CUSTOM THEME BUILDER
// ════════════════════════════════════════════════════════════════

// ThemeBuilder provides a fluent API for building custom themes.
type ThemeBuilder struct {
	theme *Theme
}

// NewThemeBuilder creates a new theme builder starting from the default theme.
func NewThemeBuilder(name string) *ThemeBuilder {
	base := DefaultTheme()
	return &ThemeBuilder{
		theme: &Theme{
			Name:   name,
			Colors: base.Colors,
		},
	}
}

// From creates a theme builder starting from an existing theme.
func (b *ThemeBuilder) From(base *Theme) *ThemeBuilder {
	// Copy colors from base
	colors := make(map[TokenClass]Color, len(base.Colors))
	for k, v := range base.Colors {
		colors[k] = v
	}
	b.theme.Colors = colors
	return b
}

// Set sets the color for a token class.
func (b *ThemeBuilder) Set(class TokenClass, color Color) *ThemeBuilder {
	b.theme.Colors[class] = color
	return b
}

// SetHex sets the color for a token class using a hex string.
func (b *ThemeBuilder) SetHex(class TokenClass, hex string) *ThemeBuilder {
	b.theme.Colors[class] = NewColor(hex)
	return b
}

// Number sets the number color.
func (b *ThemeBuilder) Number(hex string) *ThemeBuilder {
	return b.SetHex(ClassNumber, hex)
}

// Operator sets the operator color.
func (b *ThemeBuilder) Operator(hex string) *ThemeBuilder {
	return b.SetHex(ClassOperator, hex)
}

// Keyword sets the keyword color.
func (b *ThemeBuilder) Keyword(hex string) *ThemeBuilder {
	return b.SetHex(ClassKeyword, hex)
}

// Function sets the function color.
func (b *ThemeBuilder) Function(hex string) *ThemeBuilder {
	return b.SetHex(ClassFunction, hex)
}

// Currency sets the currency color.
func (b *ThemeBuilder) Currency(hex string) *ThemeBuilder {
	return b.SetHex(ClassCurrency, hex)
}

// Unit sets the unit color.
func (b *ThemeBuilder) Unit(hex string) *ThemeBuilder {
	return b.SetHex(ClassUnit, hex)
}

// Comment sets the comment color.
func (b *ThemeBuilder) Comment(hex string) *ThemeBuilder {
	return b.SetHex(ClassComment, hex)
}

// Build returns the constructed theme.
func (b *ThemeBuilder) Build() *Theme {
	return b.theme
}
