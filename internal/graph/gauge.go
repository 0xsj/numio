// internal/graph/gauge.go

package graph

import (
	"strings"
)

// GaugeStyle represents different gauge rendering styles.
type GaugeStyle int

const (
	GaugeStyleBlock   GaugeStyle = iota // [████████░░]
	GaugeStyleDot                       // [●●●●●●●○○○]
	GaugeStyleArrow                     // [========> ]
	GaugeStylePipe                      // [||||||    ]
	GaugeStyleBraille                   // [⣿⣿⣿⣿⣀⣀⣀]
)

// GaugeChars returns the filled and empty characters for a style.
func GaugeChars(style GaugeStyle) (filled, empty rune) {
	switch style {
	case GaugeStyleDot:
		return '●', '○'
	case GaugeStyleArrow:
		return '=', ' '
	case GaugeStylePipe:
		return '|', ' '
	case GaugeStyleBraille:
		return '⣿', '⣀'
	default:
		return '█', '░'
	}
}

// GaugeOptions configures gauge rendering.
type GaugeOptions struct {
	// Width is the total width of the gauge bar (excluding brackets)
	Width int

	// Style determines the character set
	Style GaugeStyle

	// Min/Max bounds
	Min float64
	Max float64

	// ShowPercent appends percentage
	ShowPercent bool

	// ShowValue appends the raw value
	ShowValue bool

	// Brackets around the gauge
	LeftBracket  string
	RightBracket string

	// Pointer character for arrow style
	Pointer rune
}

// DefaultGaugeOptions returns default gauge options.
func DefaultGaugeOptions() GaugeOptions {
	return GaugeOptions{
		Width:        20,
		Style:        GaugeStyleBlock,
		Min:          0,
		Max:          100,
		ShowPercent:  false,
		ShowValue:    false,
		LeftBracket:  "[",
		RightBracket: "]",
		Pointer:      '>',
	}
}

// Gauge generates a gauge visualization.
func Gauge(value float64) string {
	opts := DefaultGaugeOptions()
	opts.ShowPercent = true
	return GaugeWithOptions(value, opts)
}

// GaugeWithOptions generates a gauge with custom options.
func GaugeWithOptions(value float64, opts GaugeOptions) string {
	// Calculate ratio
	ratio := (value - opts.Min) / (opts.Max - opts.Min)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	// Calculate filled width
	filledWidth := int(ratio * float64(opts.Width))
	emptyWidth := opts.Width - filledWidth

	// Get characters
	filled, empty := GaugeChars(opts.Style)

	// Build gauge
	var sb strings.Builder
	sb.Grow(opts.Width + 20)

	sb.WriteString(opts.LeftBracket)

	// Handle arrow style specially
	if opts.Style == GaugeStyleArrow {
		for i := 0; i < filledWidth; i++ {
			if i == filledWidth-1 && filledWidth > 0 {
				sb.WriteRune(opts.Pointer)
			} else {
				sb.WriteRune(filled)
			}
		}
		for i := 0; i < emptyWidth; i++ {
			sb.WriteRune(empty)
		}
	} else {
		for i := 0; i < filledWidth; i++ {
			sb.WriteRune(filled)
		}
		for i := 0; i < emptyWidth; i++ {
			sb.WriteRune(empty)
		}
	}

	sb.WriteString(opts.RightBracket)

	// Add percentage
	if opts.ShowPercent {
		percent := int(ratio * 100)
		sb.WriteString(" ")
		sb.WriteString(intToString(int64(percent)))
		sb.WriteString("%")
	}

	// Add value
	if opts.ShowValue {
		sb.WriteString(" (")
		sb.WriteString(formatFloat(value))
		sb.WriteString(")")
	}

	return sb.String()
}

// GaugeWithMax generates a gauge with custom max value.
func GaugeWithMax(value, max float64) string {
	opts := DefaultGaugeOptions()
	opts.Max = max
	opts.ShowPercent = true
	return GaugeWithOptions(value, opts)
}

// GaugeWithRange generates a gauge with custom min/max.
func GaugeWithRange(value, min, max float64) string {
	opts := DefaultGaugeOptions()
	opts.Min = min
	opts.Max = max
	opts.ShowPercent = true
	return GaugeWithOptions(value, opts)
}

// GaugeStyled generates a gauge with a specific style.
func GaugeStyled(value float64, style GaugeStyle) string {
	opts := DefaultGaugeOptions()
	opts.Style = style
	opts.ShowPercent = true
	return GaugeWithOptions(value, opts)
}

// ════════════════════════════════════════════════════════════════
// SPECIALIZED GAUGES
// ════════════════════════════════════════════════════════════════

// Battery generates a battery-style gauge.
func Battery(percent float64) string {
	opts := GaugeOptions{
		Width:        10,
		Style:        GaugeStyleBlock,
		Min:          0,
		Max:          100,
		ShowPercent:  true,
		LeftBracket:  "[",
		RightBracket: "]▌",
	}
	return GaugeWithOptions(percent, opts)
}

// Meter generates a simple meter visualization.
func Meter(value, max float64) string {
	opts := DefaultGaugeOptions()
	opts.Width = 15
	opts.Max = max
	opts.ShowValue = true
	return GaugeWithOptions(value, opts)
}

// Signal generates a signal strength indicator.
func Signal(strength float64) string {
	// Clamp to 0-100
	if strength < 0 {
		strength = 0
	}
	if strength > 100 {
		strength = 100
	}

	// Use increasing height blocks
	bars := []rune{'▁', '▂', '▄', '▆', '█'}
	numBars := 5

	// Calculate how many bars to show filled
	filledBars := int((strength / 100.0) * float64(numBars))
	if strength > 0 && filledBars == 0 {
		filledBars = 1
	}

	var sb strings.Builder
	for i := 0; i < numBars; i++ {
		if i < filledBars {
			sb.WriteRune(bars[i])
		} else {
			sb.WriteRune('▁')
		}
	}

	return sb.String()
}

// Stars generates a star rating.
func Stars(rating, maxRating float64) string {
	if maxRating <= 0 {
		maxRating = 5
	}

	ratio := rating / maxRating
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	totalStars := int(maxRating)
	fullStars := int(rating)
	hasHalf := (rating - float64(fullStars)) >= 0.5

	var sb strings.Builder
	for i := 0; i < totalStars; i++ {
		if i < fullStars {
			sb.WriteString("★")
		} else if i == fullStars && hasHalf {
			sb.WriteString("⯨") // Half star
		} else {
			sb.WriteString("☆")
		}
	}

	return sb.String()
}

// Hearts generates a heart-based life indicator.
func Hearts(current, max float64) string {
	if max <= 0 {
		max = 5
	}

	totalHearts := int(max)
	fullHearts := int(current)
	hasHalf := (current - float64(fullHearts)) >= 0.5

	var sb strings.Builder
	for i := 0; i < totalHearts; i++ {
		if i < fullHearts {
			sb.WriteString("❤")
		} else if i == fullHearts && hasHalf {
			sb.WriteString("💔")
		} else {
			sb.WriteString("♡")
		}
	}

	return sb.String()
}
