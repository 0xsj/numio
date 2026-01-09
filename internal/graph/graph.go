// internal/graph/graph.go

package graph

// BlockStyle represents different character sets for rendering.
type BlockStyle int

const (
	BlockStyleDefault BlockStyle = iota // ▁▂▃▄▅▆▇█
	BlockStyleBraille                   // Braille patterns
	BlockStyleBars                      // ░▒▓█
	BlockStyleDots                      // ⋅∘○●
	BlockStyleASCII                     // _.-=+*#@
)

// String returns the style name.
func (s BlockStyle) String() string {
	switch s {
	case BlockStyleDefault:
		return "blocks"
	case BlockStyleBraille:
		return "braille"
	case BlockStyleBars:
		return "bars"
	case BlockStyleDots:
		return "dots"
	case BlockStyleASCII:
		return "ascii"
	default:
		return "unknown"
	}
}

// ParseBlockStyle parses a style name to BlockStyle.
func ParseBlockStyle(name string) BlockStyle {
	switch name {
	case "braille":
		return BlockStyleBraille
	case "bars":
		return BlockStyleBars
	case "dots":
		return BlockStyleDots
	case "ascii":
		return BlockStyleASCII
	default:
		return BlockStyleDefault
	}
}

// BlockChars returns the character set for a given style.
func BlockChars(style BlockStyle) []rune {
	switch style {
	case BlockStyleBraille:
		return []rune{'⣀', '⣄', '⣤', '⣦', '⣶', '⣷', '⣿'}
	case BlockStyleBars:
		return []rune{'░', '▒', '▓', '█'}
	case BlockStyleDots:
		return []rune{'⋅', '∘', '○', '●'}
	case BlockStyleASCII:
		return []rune{'_', '.', '-', '=', '+', '*', '#', '@'}
	default:
		return []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	}
}

// Trend represents the direction of data.
type Trend int

const (
	TrendFlat Trend = iota
	TrendUp
	TrendDown
	TrendVolatile
)

// String returns the trend as a symbol.
func (t Trend) String() string {
	switch t {
	case TrendUp:
		return "↑"
	case TrendDown:
		return "↓"
	case TrendVolatile:
		return "↕"
	default:
		return "→"
	}
}

// Symbol returns just the arrow symbol.
func (t Trend) Symbol() string {
	return t.String()
}

// Statistics holds computed statistics for a dataset.
type Statistics struct {
	Count  int
	Min    float64
	Max    float64
	Sum    float64
	Avg    float64
	First  float64
	Last   float64
	Trend  Trend
	Change float64 // Percentage change from first to last
}

// ComputeStats calculates statistics for a slice of values.
func ComputeStats(values []float64) Statistics {
	if len(values) == 0 {
		return Statistics{}
	}

	stats := Statistics{
		Count: len(values),
		Min:   values[0],
		Max:   values[0],
		First: values[0],
		Last:  values[len(values)-1],
	}

	for _, v := range values {
		stats.Sum += v
		if v < stats.Min {
			stats.Min = v
		}
		if v > stats.Max {
			stats.Max = v
		}
	}

	stats.Avg = stats.Sum / float64(stats.Count)

	// Calculate trend
	if stats.First != 0 {
		stats.Change = ((stats.Last - stats.First) / stats.First) * 100
	}

	threshold := (stats.Max - stats.Min) * 0.1 // 10% of range
	if stats.Last > stats.First+threshold {
		stats.Trend = TrendUp
	} else if stats.Last < stats.First-threshold {
		stats.Trend = TrendDown
	} else if stats.Max-stats.Min > stats.Avg*0.5 {
		stats.Trend = TrendVolatile
	} else {
		stats.Trend = TrendFlat
	}

	return stats
}
