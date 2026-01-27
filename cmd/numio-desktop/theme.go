package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ════════════════════════════════════════════════════════════════
// NUMIO THEME
// ════════════════════════════════════════════════════════════════

// NumioTheme implements fyne.Theme with numio's color scheme.
type NumioTheme struct{}

var _ fyne.Theme = (*NumioTheme)(nil)

// Color returns theme colors.
func (t *NumioTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return ColorBackground
	case theme.ColorNameForeground:
		return ColorForeground
	case theme.ColorNamePrimary:
		return ColorPrimary
	case theme.ColorNameDisabled:
		return ColorComment
	case theme.ColorNameSelection:
		return ColorSelection
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font returns the theme font.
func (t *NumioTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns theme icons.
func (t *NumioTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns theme sizes.
func (t *NumioTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameInlineIcon:
		return 16
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// ════════════════════════════════════════════════════════════════
// COLORS
// ════════════════════════════════════════════════════════════════

var (
	// Base colors
	ColorBackground = color.RGBA{R: 26, G: 26, B: 30, A: 255}    // #1a1a1e
	ColorForeground = color.RGBA{R: 217, G: 217, B: 217, A: 255} // #d9d9d9
	ColorPrimary    = color.RGBA{R: 125, G: 185, B: 235, A: 255} // #7db9eb
	ColorSelection  = color.RGBA{R: 64, G: 89, B: 140, A: 255}   // #40598c

	// Syntax colors
	ColorNumber   = color.RGBA{R: 186, G: 150, B: 227, A: 255} // #ba96e3 - purple
	ColorOperator = color.RGBA{R: 171, G: 194, B: 227, A: 255} // #abc2e3 - light blue
	ColorPercent  = color.RGBA{R: 227, G: 186, B: 150, A: 255} // #e3ba96 - orange
	ColorCurrency = color.RGBA{R: 150, G: 227, B: 161, A: 255} // #96e3a1 - green
	ColorUnit     = color.RGBA{R: 150, G: 209, B: 227, A: 255} // #96d1e3 - cyan
	ColorMetal    = color.RGBA{R: 227, G: 209, B: 125, A: 255} // #e3d17d - gold
	ColorCrypto   = color.RGBA{R: 227, G: 166, B: 84, A: 255}  // #e3a654 - bitcoin orange
	ColorFunction = color.RGBA{R: 125, G: 186, B: 235, A: 255} // #7dbae3 - blue
	ColorVariable = color.RGBA{R: 235, G: 125, B: 125, A: 255} // #eb7d7d - red
	ColorComment  = color.RGBA{R: 128, G: 128, B: 128, A: 255} // #808080 - gray
	ColorString   = color.RGBA{R: 227, G: 186, B: 150, A: 255} // #e3ba96 - orange
	ColorKeyword  = color.RGBA{R: 227, G: 125, B: 181, A: 255} // #e37db5 - pink
	ColorResult   = color.RGBA{R: 125, G: 227, B: 135, A: 255} // #7de387 - bright green
	ColorError    = color.RGBA{R: 247, G: 82, B: 74, A: 255}   // #f7524a - red
	ColorPending  = color.RGBA{R: 255, G: 166, B: 87, A: 255}  // #ffa657 - orange

	// UI colors
	ColorLineNumber   = color.RGBA{R: 102, G: 102, B: 102, A: 255} // #666666
	ColorStatusBar    = color.RGBA{R: 38, G: 38, B: 43, A: 255}    // #26262b
	ColorCursor       = color.RGBA{R: 255, G: 255, B: 255, A: 255} // #ffffff
	ColorCursorInsert = color.RGBA{R: 255, G: 255, B: 255, A: 200} // white with alpha

	// Mode colors
	ColorModeNormal = color.RGBA{R: 102, G: 153, B: 230, A: 255} // #6699e6 - blue
	ColorModeInsert = color.RGBA{R: 102, G: 204, B: 102, A: 255} // #66cc66 - green
	ColorModeVisual = color.RGBA{R: 204, G: 153, B: 51, A: 255}  // #cc9933 - yellow
)

// ════════════════════════════════════════════════════════════════
// STYLE HELPERS
// ════════════════════════════════════════════════════════════════

// StyleColor returns the color for an RPC span style.
func StyleColor(style string) color.Color {
	switch style {
	case "number":
		return ColorNumber
	case "operator":
		return ColorOperator
	case "percent":
		return ColorPercent
	case "currency":
		return ColorCurrency
	case "unit":
		return ColorUnit
	case "metal":
		return ColorMetal
	case "crypto":
		return ColorCrypto
	case "function":
		return ColorFunction
	case "variable":
		return ColorVariable
	case "comment":
		return ColorComment
	case "string":
		return ColorString
	case "keyword":
		return ColorKeyword
	case "result":
		return ColorResult
	case "error":
		return ColorError
	case "pending":
		return ColorPending
	case "lineNumber":
		return ColorLineNumber
	case "selection":
		return ColorForeground // Selection uses background highlight
	default:
		return ColorForeground
	}
}

// ModeColor returns the color for an editor mode.
func ModeColor(mode string) color.Color {
	switch mode {
	case "INSERT":
		return ColorModeInsert
	case "VISUAL":
		return ColorModeVisual
	default:
		return ColorModeNormal
	}
}
