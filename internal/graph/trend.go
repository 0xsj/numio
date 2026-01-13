// internal/graph/trend.go

package graph

import (
	"math"
	"strings"
)

// TrendStyle represents different trend indicator styles.
type TrendStyle int

const (
	TrendStyleArrow   TrendStyle = iota // ↑ ↓ → ↗ ↘
	TrendStyleLine                      // ╱ ╲ ─
	TrendStyleWord                      // up, down, flat
	TrendStyleEmoji                     // 📈 📉 ➡
	TrendStyleDelta                     // +5% -3% 0%
	TrendStyleCompact                   // ▲ ▼ ▶
)

// ════════════════════════════════════════════════════════════════
// TREND INDICATORS
// ════════════════════════════════════════════════════════════════

// TrendIndicator generates a trend indicator from values.
func TrendIndicator(values []float64) string {
	return TrendIndicatorStyled(values, TrendStyleArrow)
}

// TrendIndicatorStyled generates a trend indicator with specific style.
func TrendIndicatorStyled(values []float64, style TrendStyle) string {
	if len(values) < 2 {
		return trendSymbol(TrendFlat, style)
	}

	trend := detectTrend(values)
	return trendSymbol(trend, style)
}

// detectTrend analyzes values and returns the trend direction.
func detectTrend(values []float64) Trend {
	if len(values) < 2 {
		return TrendFlat
	}

	first := values[0]
	last := values[len(values)-1]

	// Calculate min/max for volatility check
	min, max := first, first
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	dataRange := max - min
	change := last - first

	// Threshold for significance (10% of range or 1% of value)
	threshold := dataRange * 0.1
	if threshold < math.Abs(first)*0.01 {
		threshold = math.Abs(first) * 0.01
	}
	if threshold == 0 {
		threshold = 0.001
	}

	// Check for volatility
	volatility := dataRange / math.Max(math.Abs(first), 1)
	if volatility > 0.5 && math.Abs(change) < threshold {
		return TrendVolatile
	}

	if change > threshold {
		return TrendUp
	} else if change < -threshold {
		return TrendDown
	}

	return TrendFlat
}

// trendSymbol returns the symbol for a trend in the given style.
func trendSymbol(trend Trend, style TrendStyle) string {
	switch style {
	case TrendStyleLine:
		switch trend {
		case TrendUp:
			return "╱"
		case TrendDown:
			return "╲"
		case TrendVolatile:
			return "∿"
		default:
			return "─"
		}

	case TrendStyleWord:
		switch trend {
		case TrendUp:
			return "up"
		case TrendDown:
			return "down"
		case TrendVolatile:
			return "volatile"
		default:
			return "flat"
		}

	case TrendStyleEmoji:
		switch trend {
		case TrendUp:
			return "📈"
		case TrendDown:
			return "📉"
		case TrendVolatile:
			return "📊"
		default:
			return "➡️"
		}

	case TrendStyleCompact:
		switch trend {
		case TrendUp:
			return "▲"
		case TrendDown:
			return "▼"
		case TrendVolatile:
			return "◆"
		default:
			return "▶"
		}

	default: // TrendStyleArrow
		switch trend {
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
}

// ════════════════════════════════════════════════════════════════
// CHANGE INDICATORS
// ════════════════════════════════════════════════════════════════

// ChangeIndicator shows the change between first and last value.
func ChangeIndicator(values []float64) string {
	if len(values) < 2 {
		return "0%"
	}

	first := values[0]
	last := values[len(values)-1]

	if first == 0 {
		if last > 0 {
			return "+∞%"
		} else if last < 0 {
			return "-∞%"
		}
		return "0%"
	}

	change := ((last - first) / math.Abs(first)) * 100

	return formatChange(change)
}

// formatChange formats a percentage change with sign.
func formatChange(change float64) string {
	var sb strings.Builder

	if change > 0 {
		sb.WriteString("+")
	}

	// Format with 1 decimal place
	if change == float64(int(change)) {
		sb.WriteString(intToString(int64(change)))
	} else {
		negative := change < 0
		if negative {
			change = -change
		}

		intPart := int64(change)
		fracPart := int64((change - float64(intPart)) * 10)

		if negative {
			sb.WriteString("-")
		}
		sb.WriteString(intToString(intPart))
		sb.WriteString(".")
		sb.WriteString(intToString(fracPart))
	}

	sb.WriteString("%")
	return sb.String()
}

// ChangeWithArrow shows change percentage with trend arrow.
func ChangeWithArrow(values []float64) string {
	if len(values) < 2 {
		return "→ 0%"
	}

	trend := detectTrend(values)
	change := ChangeIndicator(values)

	return trendSymbol(trend, TrendStyleArrow) + " " + change
}

// ════════════════════════════════════════════════════════════════
// DELTA COMPARISONS
// ════════════════════════════════════════════════════════════════

// Delta shows the absolute difference between two values.
func Delta(a, b float64) string {
	diff := b - a

	var sb strings.Builder

	if diff > 0 {
		sb.WriteString("+")
	}

	sb.WriteString(formatFloat(diff))

	return sb.String()
}

// DeltaPercent shows the percentage difference between two values.
func DeltaPercent(a, b float64) string {
	if a == 0 {
		if b > 0 {
			return "+∞%"
		} else if b < 0 {
			return "-∞%"
		}
		return "0%"
	}

	change := ((b - a) / math.Abs(a)) * 100
	return formatChange(change)
}

// CompareValues generates a comparison indicator.
func CompareValues(a, b float64) string {
	if a < b {
		return formatFloat(a) + " < " + formatFloat(b) + " (+" + formatFloat(b-a) + ")"
	} else if a > b {
		return formatFloat(a) + " > " + formatFloat(b) + " (" + formatFloat(b-a) + ")"
	}
	return formatFloat(a) + " = " + formatFloat(b)
}

// ════════════════════════════════════════════════════════════════
// MINI TREND LINE
// ════════════════════════════════════════════════════════════════

// MiniTrend generates a mini trend line showing direction.
func MiniTrend(values []float64) string {
	if len(values) < 2 {
		return "─"
	}

	// Simplify to a few points
	points := simplifyToPoints(values, 5)

	var sb strings.Builder
	for i := 1; i < len(points); i++ {
		if points[i] > points[i-1] {
			sb.WriteString("╱")
		} else if points[i] < points[i-1] {
			sb.WriteString("╲")
		} else {
			sb.WriteString("─")
		}
	}

	return sb.String()
}

// simplifyToPoints reduces a slice to n representative points.
func simplifyToPoints(values []float64, n int) []float64 {
	if len(values) <= n {
		return values
	}

	result := make([]float64, n)
	step := float64(len(values)-1) / float64(n-1)

	for i := 0; i < n; i++ {
		idx := int(float64(i) * step)
		if idx >= len(values) {
			idx = len(values) - 1
		}
		result[i] = values[idx]
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// SLOPE INDICATOR
// ════════════════════════════════════════════════════════════════

// Slope calculates and displays the slope of the data.
func Slope(values []float64) string {
	if len(values) < 2 {
		return "slope: 0"
	}

	// Simple linear regression slope
	n := float64(len(values))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	var sb strings.Builder
	sb.WriteString("slope: ")

	if slope > 0 {
		sb.WriteString("+")
	}

	sb.WriteString(formatFloat(slope))

	return sb.String()
}
