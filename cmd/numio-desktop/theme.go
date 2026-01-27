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
		return ColorMuted
	case theme.ColorNameSelection:
		return ColorSelection
	case theme.ColorNameShadow:
		return color.RGBA{R: 0, G: 0, B: 0, A: 60}
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
		return 6
	case theme.SizeNameInlineIcon:
		return 16
	case theme.SizeNameScrollBar:
		return 8
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// ════════════════════════════════════════════════════════════════
// COLORS - TRANSLUCENT DARK THEME
// ════════════════════════════════════════════════════════════════

var (
	// Base colors (with transparency for glass effect)
	ColorBackground      = color.RGBA{R: 30, G: 30, B: 35, A: 230} // Semi-transparent dark
	ColorBackgroundSolid = color.RGBA{R: 30, G: 30, B: 35, A: 255} // Solid for popups
	ColorForeground      = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	ColorMuted           = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	ColorPrimary         = color.RGBA{R: 130, G: 190, B: 240, A: 255}
	ColorSelection       = color.RGBA{R: 60, G: 85, B: 130, A: 200}

	// Syntax colors - softer, Numi-inspired
	ColorNumber   = color.RGBA{R: 180, G: 150, B: 220, A: 255} // Soft purple
	ColorOperator = color.RGBA{R: 160, G: 180, B: 210, A: 255} // Muted blue
	ColorPercent  = color.RGBA{R: 220, G: 180, B: 140, A: 255} // Warm orange
	ColorCurrency = color.RGBA{R: 140, G: 200, B: 140, A: 255} // Soft green
	ColorUnit     = color.RGBA{R: 140, G: 190, B: 210, A: 255} // Soft cyan
	ColorMetal    = color.RGBA{R: 220, G: 200, B: 120, A: 255} // Gold
	ColorCrypto   = color.RGBA{R: 240, G: 180, B: 100, A: 255} // Bitcoin orange
	ColorFunction = color.RGBA{R: 130, G: 180, B: 230, A: 255} // Function blue
	ColorVariable = color.RGBA{R: 220, G: 140, B: 140, A: 255} // Soft red
	ColorComment  = color.RGBA{R: 110, G: 110, B: 110, A: 255} // Gray
	ColorString   = color.RGBA{R: 220, G: 180, B: 140, A: 255} // Warm
	ColorKeyword  = color.RGBA{R: 220, G: 130, B: 180, A: 255} // Pink
	ColorResult   = color.RGBA{R: 140, G: 210, B: 150, A: 255} // Soft green
	ColorError    = color.RGBA{R: 240, G: 100, B: 100, A: 255} // Red
	ColorPending  = color.RGBA{R: 240, G: 180, B: 100, A: 255} // Orange

	// UI colors
	ColorTilde       = color.RGBA{R: 70, G: 70, B: 80, A: 255} // Vim tilde color
	ColorStatusBar   = color.RGBA{R: 40, G: 40, B: 48, A: 240} // Slightly lighter, semi-transparent
	ColorStatusText  = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	ColorCursor      = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColorCursorBlock = color.RGBA{R: 255, G: 255, B: 255, A: 160}

	// Mode colors
	ColorModeNormal = color.RGBA{R: 100, G: 150, B: 220, A: 255} // Blue
	ColorModeInsert = color.RGBA{R: 100, G: 190, B: 100, A: 255} // Green
	ColorModeVisual = color.RGBA{R: 220, G: 170, B: 80, A: 255}  // Yellow/Orange

	// Popup colors
	ColorPopupBackground = color.RGBA{R: 35, G: 35, B: 42, A: 250}
	ColorPopupBorder     = color.RGBA{R: 80, G: 120, B: 160, A: 255}
	ColorPopupTitle      = color.RGBA{R: 130, G: 190, B: 240, A: 255}
	ColorPopupHeading    = color.RGBA{R: 220, G: 170, B: 100, A: 255}
	ColorPopupKey        = color.RGBA{R: 180, G: 150, B: 220, A: 255}
	ColorPopupDesc       = color.RGBA{R: 180, G: 180, B: 180, A: 255}
	ColorPopupHint       = color.RGBA{R: 120, G: 120, B: 120, A: 255}
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
		return ColorMuted
	case "selection":
		return ColorForeground
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

// ModeName returns a short mode name for status bar.
func ModeName(mode string) string {
	switch mode {
	case "INSERT":
		return "INSERT"
	case "VISUAL":
		return "VISUAL"
	default:
		return "NORMAL"
	}
}
