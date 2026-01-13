// internal/graph/histogram.go

package graph

import (
	"math"
	"strings"
)

// HistogramOptions configures histogram rendering.
type HistogramOptions struct {
	// Bins is the number of buckets (default: auto-detect)
	Bins int

	// Style determines the character set
	Style BlockStyle

	// Min/Max bounds (NaN = auto-detect)
	Min float64
	Max float64

	// ShowCounts appends count numbers
	ShowCounts bool
}

// DefaultHistogramOptions returns default histogram options.
func DefaultHistogramOptions() HistogramOptions {
	return HistogramOptions{
		Bins:       0, // Auto-detect
		Style:      BlockStyleDefault,
		Min:        math.NaN(),
		Max:        math.NaN(),
		ShowCounts: false,
	}
}

// Histogram generates a histogram from values.
func Histogram(values []float64) string {
	return HistogramWithOptions(values, DefaultHistogramOptions())
}

// HistogramWithOptions generates a histogram with custom options.
func HistogramWithOptions(values []float64, opts HistogramOptions) string {
	if len(values) == 0 {
		return ""
	}

	// Determine bounds
	min, max := opts.Min, opts.Max
	if math.IsNaN(min) || math.IsNaN(max) {
		detectedMin, detectedMax := detectBounds(values)
		if math.IsNaN(min) {
			min = detectedMin
		}
		if math.IsNaN(max) {
			max = detectedMax
		}
	}

	// Handle flat data
	if min >= max {
		max = min + 1
	}

	// Determine number of bins
	bins := opts.Bins
	if bins <= 0 {
		// Sturges' formula: k = ceil(log2(n) + 1)
		bins = int(math.Ceil(math.Log2(float64(len(values))) + 1))
		if bins < 5 {
			bins = 5
		}
		if bins > 20 {
			bins = 20
		}
	}

	// Count values in each bin
	counts := make([]int, bins)
	binWidth := (max - min) / float64(bins)

	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}

		// Determine bin index
		idx := int((v - min) / binWidth)
		if idx < 0 {
			idx = 0
		}
		if idx >= bins {
			idx = bins - 1
		}
		counts[idx]++
	}

	// Find max count for scaling
	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	if maxCount == 0 {
		return strings.Repeat(string(BlockChars(opts.Style)[0]), bins)
	}

	// Build histogram
	blocks := BlockChars(opts.Style)
	numBlocks := len(blocks)

	var sb strings.Builder
	sb.Grow(bins * 4)

	for _, c := range counts {
		if c == 0 {
			sb.WriteRune(blocks[0])
		} else {
			// Scale count to block index
			normalized := float64(c) / float64(maxCount)
			idx := int(normalized * float64(numBlocks-1))
			if idx >= numBlocks {
				idx = numBlocks - 1
			}
			sb.WriteRune(blocks[idx])
		}
	}

	return sb.String()
}

// HistogramWithBins generates a histogram with specified bin count.
func HistogramWithBins(values []float64, bins int) string {
	opts := DefaultHistogramOptions()
	opts.Bins = bins
	return HistogramWithOptions(values, opts)
}

// HistogramResult holds histogram with statistics.
type HistogramResult struct {
	Histogram string
	Bins      int
	Counts    []int
	Min       float64
	Max       float64
	BinWidth  float64
}

// HistogramWithStats generates histogram and returns statistics.
func HistogramWithStats(values []float64, bins int) HistogramResult {
	if len(values) == 0 {
		return HistogramResult{}
	}

	min, max := detectBounds(values)
	if min >= max {
		max = min + 1
	}

	if bins <= 0 {
		bins = int(math.Ceil(math.Log2(float64(len(values))) + 1))
		if bins < 5 {
			bins = 5
		}
		if bins > 20 {
			bins = 20
		}
	}

	counts := make([]int, bins)
	binWidth := (max - min) / float64(bins)

	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}
		idx := int((v - min) / binWidth)
		if idx < 0 {
			idx = 0
		}
		if idx >= bins {
			idx = bins - 1
		}
		counts[idx]++
	}

	opts := DefaultHistogramOptions()
	opts.Bins = bins

	return HistogramResult{
		Histogram: HistogramWithOptions(values, opts),
		Bins:      bins,
		Counts:    counts,
		Min:       min,
		Max:       max,
		BinWidth:  binWidth,
	}
}
