// internal/stats/central.go

package stats

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// MEDIAN
// ════════════════════════════════════════════════════════════════

// Median returns the median (middle value) of the data.
func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := Sort(values)
	return MedianOfSorted(sorted)
}

// MedianOfSorted returns the median of pre-sorted data.
func MedianOfSorted(sorted []float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}

	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

// ════════════════════════════════════════════════════════════════
// MODE
// ════════════════════════════════════════════════════════════════

// Mode returns the most frequent value.
// If multiple modes exist, returns the smallest.
func Mode(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	counts := make(map[float64]int)
	for _, v := range values {
		counts[v]++
	}

	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}

	mode := math.Inf(1)
	for val, count := range counts {
		if count == maxCount && val < mode {
			mode = val
		}
	}

	return mode
}

// Modes returns all modes (values with highest frequency).
func Modes(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	counts := make(map[float64]int)
	for _, v := range values {
		counts[v]++
	}

	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}

	var modes []float64
	for val, count := range counts {
		if count == maxCount {
			modes = append(modes, val)
		}
	}

	return Sort(modes)
}

// IsMultimodal checks if the data has multiple modes.
func IsMultimodal(values []float64) bool {
	return len(Modes(values)) > 1
}

// ════════════════════════════════════════════════════════════════
// GEOMETRIC MEAN
// ════════════════════════════════════════════════════════════════

// GeometricMean returns the geometric mean.
// All values must be positive.
func GeometricMean(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	if !AllPositive(values) {
		return 0, ErrNonPositive
	}

	// Use log to avoid overflow: exp(mean(log(x)))
	var sumLog float64
	for _, v := range values {
		sumLog += math.Log(v)
	}

	return math.Exp(sumLog / float64(len(values))), nil
}

// GeometricMeanOrZero returns geometric mean or 0 on error.
func GeometricMeanOrZero(values []float64) float64 {
	result, err := GeometricMean(values)
	if err != nil {
		return 0
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// HARMONIC MEAN
// ════════════════════════════════════════════════════════════════

// HarmonicMean returns the harmonic mean.
// All values must be non-zero.
func HarmonicMean(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	if !AllNonZero(values) {
		return 0, ErrZeroValue
	}

	var sumRecip float64
	for _, v := range values {
		sumRecip += 1 / v
	}

	return float64(len(values)) / sumRecip, nil
}

// HarmonicMeanOrZero returns harmonic mean or 0 on error.
func HarmonicMeanOrZero(values []float64) float64 {
	result, err := HarmonicMean(values)
	if err != nil {
		return 0
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// WEIGHTED MEAN
// ════════════════════════════════════════════════════════════════

// WeightedMean returns the weighted arithmetic mean.
// values and weights must have the same length.
func WeightedMean(values, weights []float64) (float64, error) {
	if len(values) != len(weights) {
		return 0, ErrLengthMismatch
	}

	if len(values) == 0 {
		return 0, nil
	}

	var sumWeighted, sumWeights float64
	for i, v := range values {
		w := weights[i]
		if w < 0 {
			return 0, ErrNegativeWeight
		}
		sumWeighted += v * w
		sumWeights += w
	}

	if sumWeights == 0 {
		return 0, ErrZeroWeight
	}

	return sumWeighted / sumWeights, nil
}

// ════════════════════════════════════════════════════════════════
// TRIMMED MEAN
// ════════════════════════════════════════════════════════════════

// TrimmedMean returns the mean after trimming a percentage from each end.
// trimPercent is the percentage to trim from EACH end (0-50).
func TrimmedMean(values []float64, trimPercent float64) float64 {
	if len(values) == 0 {
		return 0
	}

	if trimPercent <= 0 {
		return Mean(values)
	}

	if trimPercent >= 50 {
		return Median(values)
	}

	sorted := Sort(values)
	n := len(sorted)

	// Calculate how many to trim from each end
	trimCount := int(float64(n) * trimPercent / 100)
	if trimCount == 0 {
		trimCount = 1
	}

	// Ensure we don't trim everything
	if trimCount*2 >= n {
		return Median(values)
	}

	// Calculate mean of remaining values
	trimmed := sorted[trimCount : n-trimCount]
	return Mean(trimmed)
}

// ════════════════════════════════════════════════════════════════
// MIDRANGE
// ════════════════════════════════════════════════════════════════

// Midrange returns (max + min) / 2.
func Midrange(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return (Min(values) + Max(values)) / 2
}

// ════════════════════════════════════════════════════════════════
// ERRORS
// ════════════════════════════════════════════════════════════════

// StatsError represents a statistics error.
type StatsError string

func (e StatsError) Error() string {
	return string(e)
}

const (
	ErrNonPositive    StatsError = "all values must be positive"
	ErrZeroValue      StatsError = "cannot include zero values"
	ErrLengthMismatch StatsError = "values and weights must have same length"
	ErrNegativeWeight StatsError = "weights cannot be negative"
	ErrZeroWeight     StatsError = "sum of weights cannot be zero"
)
