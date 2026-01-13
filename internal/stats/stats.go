// internal/stats/stats.go

package stats

import (
	"math"
	"sort"
)

// ════════════════════════════════════════════════════════════════
// BASIC OPERATIONS
// ════════════════════════════════════════════════════════════════

// Sum returns the sum of values.
func Sum(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum
}

// Product returns the product of all values.
func Product(values []float64) float64 {
	if len(values) == 0 {
		return 1
	}

	product := 1.0
	for _, v := range values {
		product *= v
	}
	return product
}

// Mean returns the arithmetic mean.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Sum(values) / float64(len(values))
}

// Min returns the minimum value.
func Min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum value.
func Max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// ════════════════════════════════════════════════════════════════
// SORTING
// ════════════════════════════════════════════════════════════════

// Sort returns a sorted copy of the values.
func Sort(values []float64) []float64 {
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	return sorted
}

// IsSorted checks if values are sorted.
func IsSorted(values []float64) bool {
	return sort.Float64sAreSorted(values)
}

// ════════════════════════════════════════════════════════════════
// VALIDATION
// ════════════════════════════════════════════════════════════════

// AllPositive checks if all values are positive.
func AllPositive(values []float64) bool {
	for _, v := range values {
		if v <= 0 {
			return false
		}
	}
	return true
}

// AllNonZero checks if all values are non-zero.
func AllNonZero(values []float64) bool {
	for _, v := range values {
		if v == 0 {
			return false
		}
	}
	return true
}

// HasNaN checks if any value is NaN.
func HasNaN(values []float64) bool {
	for _, v := range values {
		if math.IsNaN(v) {
			return true
		}
	}
	return false
}

// FilterNaN removes NaN values.
func FilterNaN(values []float64) []float64 {
	result := make([]float64, 0, len(values))
	for _, v := range values {
		if !math.IsNaN(v) {
			result = append(result, v)
		}
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// STATISTICS RESULT
// ════════════════════════════════════════════════════════════════

// Summary holds a statistical summary of data.
type Summary struct {
	Count    int
	Sum      float64
	Mean     float64
	Min      float64
	Max      float64
	Range    float64
	Variance float64
	StdDev   float64
	Median   float64
	Q1       float64
	Q3       float64
	IQR      float64
}

// Summarize computes a full statistical summary.
func Summarize(values []float64) Summary {
	if len(values) == 0 {
		return Summary{}
	}

	sorted := Sort(values)
	n := len(values)

	sum := Sum(values)
	mean := sum / float64(n)

	// Variance
	var sumSq float64
	for _, v := range values {
		diff := v - mean
		sumSq += diff * diff
	}
	variance := sumSq / float64(n)

	return Summary{
		Count:    n,
		Sum:      sum,
		Mean:     mean,
		Min:      sorted[0],
		Max:      sorted[n-1],
		Range:    sorted[n-1] - sorted[0],
		Variance: variance,
		StdDev:   math.Sqrt(variance),
		Median:   PercentileOfSorted(sorted, 50),
		Q1:       PercentileOfSorted(sorted, 25),
		Q3:       PercentileOfSorted(sorted, 75),
		IQR:      PercentileOfSorted(sorted, 75) - PercentileOfSorted(sorted, 25),
	}
}

// PercentileOfSorted calculates percentile of pre-sorted values.
func PercentileOfSorted(sorted []float64, p float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return sorted[0]
	}

	rank := (p / 100) * float64(n-1)
	lower := int(math.Floor(rank))
	upper := int(math.Ceil(rank))

	if lower == upper || upper >= n {
		return sorted[lower]
	}

	frac := rank - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}
