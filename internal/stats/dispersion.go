// internal/stats/dispersion.go

package stats

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// VARIANCE
// ════════════════════════════════════════════════════════════════

// Variance returns the population variance.
func Variance(values []float64) float64 {
	return varianceInternal(values, false)
}

// SampleVariance returns the sample variance (n-1 denominator).
func SampleVariance(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	return varianceInternal(values, true)
}

// varianceInternal computes variance with optional sample correction.
func varianceInternal(values []float64, sample bool) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	if sample && n < 2 {
		return 0
	}

	mean := Mean(values)

	var sumSq float64
	for _, v := range values {
		diff := v - mean
		sumSq += diff * diff
	}

	denom := float64(n)
	if sample {
		denom = float64(n - 1)
	}

	return sumSq / denom
}

// ════════════════════════════════════════════════════════════════
// STANDARD DEVIATION
// ════════════════════════════════════════════════════════════════

// StdDev returns the population standard deviation.
func StdDev(values []float64) float64 {
	return math.Sqrt(Variance(values))
}

// SampleStdDev returns the sample standard deviation (n-1 denominator).
func SampleStdDev(values []float64) float64 {
	return math.Sqrt(SampleVariance(values))
}

// ════════════════════════════════════════════════════════════════
// RANGE
// ════════════════════════════════════════════════════════════════

// Range returns max - min.
func Range(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Max(values) - Min(values)
}

// ════════════════════════════════════════════════════════════════
// MEAN ABSOLUTE DEVIATION
// ════════════════════════════════════════════════════════════════

// MAD returns the mean absolute deviation from the mean.
func MAD(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := Mean(values)

	var sumAbs float64
	for _, v := range values {
		sumAbs += math.Abs(v - mean)
	}

	return sumAbs / float64(len(values))
}

// MADMedian returns the mean absolute deviation from the median.
func MADMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	median := Median(values)

	var sumAbs float64
	for _, v := range values {
		sumAbs += math.Abs(v - median)
	}

	return sumAbs / float64(len(values))
}

// ════════════════════════════════════════════════════════════════
// ROOT MEAN SQUARE
// ════════════════════════════════════════════════════════════════

// RMS returns the root mean square.
func RMS(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sumSq float64
	for _, v := range values {
		sumSq += v * v
	}

	return math.Sqrt(sumSq / float64(len(values)))
}

// ════════════════════════════════════════════════════════════════
// COEFFICIENT OF VARIATION
// ════════════════════════════════════════════════════════════════

// CoefficientOfVariation returns stddev / mean (as a ratio).
// Returns 0 if mean is zero.
func CoefficientOfVariation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := Mean(values)
	if mean == 0 {
		return 0
	}

	return StdDev(values) / math.Abs(mean)
}

// CoefficientOfVariationPercent returns CV as a percentage.
func CoefficientOfVariationPercent(values []float64) float64 {
	return CoefficientOfVariation(values) * 100
}

// ════════════════════════════════════════════════════════════════
// SKEWNESS & KURTOSIS
// ════════════════════════════════════════════════════════════════

// Skewness returns the skewness (measure of asymmetry).
// Positive = right-skewed, Negative = left-skewed.
func Skewness(values []float64) float64 {
	n := len(values)
	if n < 3 {
		return 0
	}

	mean := Mean(values)
	stddev := StdDev(values)

	if stddev == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		z := (v - mean) / stddev
		sum += z * z * z
	}

	return sum / float64(n)
}

// Kurtosis returns the excess kurtosis.
// Positive = heavy tails, Negative = light tails.
// Normal distribution has kurtosis = 0.
func Kurtosis(values []float64) float64 {
	n := len(values)
	if n < 4 {
		return 0
	}

	mean := Mean(values)
	stddev := StdDev(values)

	if stddev == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		z := (v - mean) / stddev
		sum += z * z * z * z
	}

	// Excess kurtosis (subtract 3 for normal distribution)
	return sum/float64(n) - 3
}

// ════════════════════════════════════════════════════════════════
// STANDARD ERROR
// ════════════════════════════════════════════════════════════════

// StandardError returns the standard error of the mean.
func StandardError(values []float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}

	return SampleStdDev(values) / math.Sqrt(float64(n))
}

// ════════════════════════════════════════════════════════════════
// RELATIVE MEASURES
// ════════════════════════════════════════════════════════════════

// RelativeRange returns range / mean.
func RelativeRange(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := Mean(values)
	if mean == 0 {
		return 0
	}

	return Range(values) / math.Abs(mean)
}

// QuartileDeviation returns (Q3 - Q1) / 2.
func QuartileDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := Sort(values)
	q1 := PercentileOfSorted(sorted, 25)
	q3 := PercentileOfSorted(sorted, 75)

	return (q3 - q1) / 2
}
