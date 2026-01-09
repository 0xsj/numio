// internal/graph/sparkline.go

package graph

import (
	"math"
	"strings"
)

// Sparkline generates a sparkline string from values with default options.
func Sparkline(values []float64) string {
	return SparklineWithOptions(values, DefaultOptions())
}

// SparklineWithOptions generates a sparkline with custom options.
func SparklineWithOptions(values []float64, opts Options) string {
	if len(values) == 0 {
		return ""
	}

	// Resample if width is specified
	data := values
	if opts.Width > 0 && len(values) != opts.Width {
		data = Resample(values, opts.Width)
	}

	// Resolve bounds
	min, max := opts.ResolveBounds(data)

	// Get character set
	blocks := BlockChars(opts.Style)
	numBlocks := len(blocks)

	// Build sparkline
	var sb strings.Builder
	sb.Grow(len(data)*4 + len(opts.Prefix) + len(opts.Suffix) + 10)

	// Add prefix
	if opts.Prefix != "" {
		sb.WriteString(opts.Prefix)
	}

	// Render each value
	dataRange := max - min
	for _, v := range data {
		// Handle special values
		if math.IsNaN(v) {
			sb.WriteRune(opts.NaNChar)
			continue
		}
		if math.IsInf(v, 0) {
			if v > 0 {
				sb.WriteRune(blocks[numBlocks-1])
			} else {
				sb.WriteRune(blocks[0])
			}
			continue
		}

		// Normalize and map to block index
		index := valueToIndex(v, min, dataRange, numBlocks)
		sb.WriteRune(blocks[index])
	}

	// Add trend indicator
	if opts.ShowTrend && len(values) >= 2 {
		stats := ComputeStats(values)
		sb.WriteString(" ")
		sb.WriteString(stats.Trend.Symbol())
	}

	// Add suffix
	if opts.Suffix != "" {
		sb.WriteString(opts.Suffix)
	}

	return sb.String()
}

// valueToIndex maps a value to a block index.
func valueToIndex(v, min, dataRange float64, numBlocks int) int {
	if dataRange == 0 {
		return numBlocks / 2 // Middle block for flat data
	}

	// Normalize to 0-1 range
	normalized := (v - min) / dataRange

	// Clamp to [0, 1]
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	// Map to block index
	index := int(normalized * float64(numBlocks-1))

	// Clamp to valid range
	if index < 0 {
		index = 0
	}
	if index >= numBlocks {
		index = numBlocks - 1
	}

	return index
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE CONSTRUCTORS
// ════════════════════════════════════════════════════════════════

// SparklineStyled generates a sparkline with a specific style.
func SparklineStyled(values []float64, style BlockStyle) string {
	opts := DefaultOptions()
	opts.Style = style
	return SparklineWithOptions(values, opts)
}

// SparklineBounded generates a sparkline with explicit bounds.
func SparklineBounded(values []float64, min, max float64) string {
	opts := DefaultOptions()
	opts.Min = min
	opts.Max = max
	return SparklineWithOptions(values, opts)
}

// SparklineFixed generates a sparkline with fixed width.
func SparklineFixed(values []float64, width int) string {
	opts := DefaultOptions()
	opts.Width = width
	return SparklineWithOptions(values, opts)
}

// SparklineWithTrend generates a sparkline with trend indicator.
func SparklineWithTrend(values []float64) string {
	opts := DefaultOptions()
	opts.ShowTrend = true
	return SparklineWithOptions(values, opts)
}

// ════════════════════════════════════════════════════════════════
// SPARKLINE RESULT
// ════════════════════════════════════════════════════════════════

// SparklineResult holds a sparkline with its statistics.
type SparklineResult struct {
	Sparkline string
	Stats     Statistics
	Options   Options
}

// SparklineWithStats generates a sparkline and computes statistics.
func SparklineWithStats(values []float64, opts ...OptionFunc) SparklineResult {
	options := DefaultOptions().Apply(opts...)

	return SparklineResult{
		Sparkline: SparklineWithOptions(values, options),
		Stats:     ComputeStats(values),
		Options:   options,
	}
}

// String returns the sparkline string.
func (r SparklineResult) String() string {
	return r.Sparkline
}

// Full returns sparkline with stats summary.
func (r SparklineResult) Full() string {
	var sb strings.Builder
	sb.WriteString(r.Sparkline)
	sb.WriteString(" (")
	sb.WriteString("min:")
	sb.WriteString(formatFloat(r.Stats.Min))
	sb.WriteString(" max:")
	sb.WriteString(formatFloat(r.Stats.Max))
	sb.WriteString(" avg:")
	sb.WriteString(formatFloat(r.Stats.Avg))
	sb.WriteString(")")
	return sb.String()
}

// WithTrendSymbol returns sparkline with trend arrow.
func (r SparklineResult) WithTrendSymbol() string {
	return r.Sparkline + " " + r.Stats.Trend.Symbol()
}

// ════════════════════════════════════════════════════════════════
// INTEGER HELPERS
// ════════════════════════════════════════════════════════════════

// SparklineInt generates a sparkline from integers.
func SparklineInt(values []int) string {
	return Sparkline(intsToFloats(values))
}

// SparklineInt64 generates a sparkline from int64s.
func SparklineInt64(values []int64) string {
	floats := make([]float64, len(values))
	for i, v := range values {
		floats[i] = float64(v)
	}
	return Sparkline(floats)
}

// intsToFloats converts []int to []float64.
func intsToFloats(values []int) []float64 {
	floats := make([]float64, len(values))
	for i, v := range values {
		floats[i] = float64(v)
	}
	return floats
}

// ════════════════════════════════════════════════════════════════
// FORMATTING HELPERS
// ════════════════════════════════════════════════════════════════

// formatFloat formats a float for display.
func formatFloat(f float64) string {
	if math.IsNaN(f) {
		return "NaN"
	}
	if math.IsInf(f, 1) {
		return "+Inf"
	}
	if math.IsInf(f, -1) {
		return "-Inf"
	}

	// Simple formatting without fmt package
	if f == float64(int64(f)) {
		return intToString(int64(f))
	}

	// Format with 2 decimal places
	negative := f < 0
	if negative {
		f = -f
	}

	intPart := int64(f)
	fracPart := int64((f - float64(intPart)) * 100)

	result := intToString(intPart) + "." + padLeft(intToString(fracPart), 2, '0')

	// Trim trailing zeros
	result = strings.TrimRight(result, "0")
	result = strings.TrimRight(result, ".")

	if negative {
		result = "-" + result
	}

	return result
}

// intToString converts an int64 to string without fmt.
func intToString(n int64) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}

// padLeft pads a string on the left to reach a minimum length.
func padLeft(s string, length int, pad rune) string {
	for len(s) < length {
		s = string(pad) + s
	}
	return s
}
