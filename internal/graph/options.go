// internal/graph/options.go

package graph

import (
	"math"
)

// Options configures graph rendering.
type Options struct {
	// Style determines the character set used.
	Style BlockStyle

	// Min/Max bounds for scaling (NaN = auto-detect)
	Min float64
	Max float64

	// Width constrains output to fixed number of characters.
	// 0 = no constraint (one char per value)
	Width int

	// ShowTrend appends a trend indicator (↑↓→↕)
	ShowTrend bool

	// ShowStats appends statistics summary
	ShowStats bool

	// NaNChar is the character to use for NaN/missing values.
	NaNChar rune

	// GapChar is the character to use for gaps in data.
	GapChar rune

	// Prefix is prepended to the output.
	Prefix string

	// Suffix is appended to the output.
	Suffix string
}

// DefaultOptions returns the default rendering options.
func DefaultOptions() Options {
	return Options{
		Style:     BlockStyleDefault,
		Min:       math.NaN(),
		Max:       math.NaN(),
		Width:     0,
		ShowTrend: false,
		ShowStats: false,
		NaNChar:   '·',
		GapChar:   ' ',
		Prefix:    "",
		Suffix:    "",
	}
}

// OptionFunc is a function that modifies Options.
type OptionFunc func(*Options)

// Apply applies option functions to create final Options.
func (o Options) Apply(opts ...OptionFunc) Options {
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

// ════════════════════════════════════════════════════════════════
// OPTION FUNCTIONS
// ════════════════════════════════════════════════════════════════

// WithStyle sets the block style.
func WithStyle(style BlockStyle) OptionFunc {
	return func(o *Options) {
		o.Style = style
	}
}

// WithStyleName sets the block style by name.
func WithStyleName(name string) OptionFunc {
	return func(o *Options) {
		o.Style = ParseBlockStyle(name)
	}
}

// WithBounds sets explicit min/max bounds.
func WithBounds(min, max float64) OptionFunc {
	return func(o *Options) {
		o.Min = min
		o.Max = max
	}
}

// WithMin sets explicit minimum bound.
func WithMin(min float64) OptionFunc {
	return func(o *Options) {
		o.Min = min
	}
}

// WithMax sets explicit maximum bound.
func WithMax(max float64) OptionFunc {
	return func(o *Options) {
		o.Max = max
	}
}

// WithWidth sets a fixed output width.
func WithWidth(width int) OptionFunc {
	return func(o *Options) {
		o.Width = width
	}
}

// WithTrend enables the trend indicator.
func WithTrend(show bool) OptionFunc {
	return func(o *Options) {
		o.ShowTrend = show
	}
}

// WithStats enables the statistics summary.
func WithStats(show bool) OptionFunc {
	return func(o *Options) {
		o.ShowStats = show
	}
}

// WithNaNChar sets the character for NaN values.
func WithNaNChar(c rune) OptionFunc {
	return func(o *Options) {
		o.NaNChar = c
	}
}

// WithGapChar sets the character for gaps.
func WithGapChar(c rune) OptionFunc {
	return func(o *Options) {
		o.GapChar = c
	}
}

// WithPrefix sets a prefix string.
func WithPrefix(prefix string) OptionFunc {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

// WithSuffix sets a suffix string.
func WithSuffix(suffix string) OptionFunc {
	return func(o *Options) {
		o.Suffix = suffix
	}
}

// ════════════════════════════════════════════════════════════════
// BOUNDS DETECTION
// ════════════════════════════════════════════════════════════════

// ResolveBounds determines the final min/max bounds.
// Uses explicit bounds if set, otherwise auto-detects from data.
func (o Options) ResolveBounds(values []float64) (min, max float64) {
	// Start with explicit bounds if set
	min = o.Min
	max = o.Max

	// Auto-detect if needed
	autoMin := math.IsNaN(min)
	autoMax := math.IsNaN(max)

	if autoMin || autoMax {
		detectedMin, detectedMax := detectBounds(values)
		if autoMin {
			min = detectedMin
		}
		if autoMax {
			max = detectedMax
		}
	}

	// Handle edge case where min >= max
	if min >= max {
		if min == 0 {
			max = 1
		} else {
			// Add 10% padding
			padding := math.Abs(min) * 0.1
			min -= padding
			max += padding
		}
	}

	return min, max
}

// detectBounds finds min/max from values, ignoring NaN.
func detectBounds(values []float64) (min, max float64) {
	min = math.Inf(1)
	max = math.Inf(-1)

	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// If no valid values found
	if math.IsInf(min, 1) {
		min = 0
		max = 1
	}

	return min, max
}

// ════════════════════════════════════════════════════════════════
// DATA RESAMPLING
// ════════════════════════════════════════════════════════════════

// Resample reduces or expands data to fit a target width.
func Resample(values []float64, targetWidth int) []float64 {
	if targetWidth <= 0 || len(values) == 0 {
		return values
	}

	if len(values) == targetWidth {
		return values
	}

	result := make([]float64, targetWidth)

	if len(values) < targetWidth {
		// Expand: interpolate
		for i := 0; i < targetWidth; i++ {
			pos := float64(i) * float64(len(values)-1) / float64(targetWidth-1)
			lower := int(math.Floor(pos))
			upper := int(math.Ceil(pos))

			if lower == upper || upper >= len(values) {
				result[i] = values[lower]
			} else {
				// Linear interpolation
				frac := pos - float64(lower)
				result[i] = values[lower]*(1-frac) + values[upper]*frac
			}
		}
	} else {
		// Reduce: average buckets
		bucketSize := float64(len(values)) / float64(targetWidth)

		for i := 0; i < targetWidth; i++ {
			start := int(float64(i) * bucketSize)
			end := int(float64(i+1) * bucketSize)
			if end > len(values) {
				end = len(values)
			}

			sum := 0.0
			count := 0
			for j := start; j < end; j++ {
				if !math.IsNaN(values[j]) {
					sum += values[j]
					count++
				}
			}

			if count > 0 {
				result[i] = sum / float64(count)
			} else {
				result[i] = math.NaN()
			}
		}
	}

	return result
}
