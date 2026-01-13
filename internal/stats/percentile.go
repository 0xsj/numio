// internal/stats/percentile.go

package stats

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// PERCENTILE
// ════════════════════════════════════════════════════════════════

// Percentile returns the pth percentile of the data.
// p should be between 0 and 100.
func Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Clamp p to valid range
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}

	sorted := Sort(values)
	return PercentileOfSorted(sorted, p)
}

// PercentileRank returns the percentile rank of a value in the dataset.
// Returns the percentage of values that are less than or equal to x.
func PercentileRank(values []float64, x float64) float64 {
	if len(values) == 0 {
		return 0
	}

	countBelow := 0
	countEqual := 0

	for _, v := range values {
		if v < x {
			countBelow++
		} else if v == x {
			countEqual++
		}
	}

	// Use midpoint method
	rank := float64(countBelow) + float64(countEqual)/2
	return (rank / float64(len(values))) * 100
}

// ════════════════════════════════════════════════════════════════
// QUARTILES
// ════════════════════════════════════════════════════════════════

// Q1 returns the first quartile (25th percentile).
func Q1(values []float64) float64 {
	return Percentile(values, 25)
}

// Q2 returns the second quartile (50th percentile, median).
func Q2(values []float64) float64 {
	return Percentile(values, 50)
}

// Q3 returns the third quartile (75th percentile).
func Q3(values []float64) float64 {
	return Percentile(values, 75)
}

// Quartile returns the qth quartile (q = 1, 2, or 3).
func Quartile(values []float64, q int) float64 {
	switch q {
	case 1:
		return Q1(values)
	case 2:
		return Q2(values)
	case 3:
		return Q3(values)
	default:
		return 0
	}
}

// Quartiles returns all three quartiles.
func Quartiles(values []float64) (q1, q2, q3 float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	sorted := Sort(values)
	q1 = PercentileOfSorted(sorted, 25)
	q2 = PercentileOfSorted(sorted, 50)
	q3 = PercentileOfSorted(sorted, 75)
	return
}

// ════════════════════════════════════════════════════════════════
// INTERQUARTILE RANGE
// ════════════════════════════════════════════════════════════════

// IQR returns the interquartile range (Q3 - Q1).
func IQR(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := Sort(values)
	q1 := PercentileOfSorted(sorted, 25)
	q3 := PercentileOfSorted(sorted, 75)

	return q3 - q1
}

// ════════════════════════════════════════════════════════════════
// QUINTILES & DECILES
// ════════════════════════════════════════════════════════════════

// Quintile returns the nth quintile (n = 1 to 4).
// Quintiles divide data into 5 equal parts.
func Quintile(values []float64, n int) float64 {
	if n < 1 || n > 4 {
		return 0
	}
	return Percentile(values, float64(n)*20)
}

// Quintiles returns all four quintile boundaries.
func Quintiles(values []float64) [4]float64 {
	var result [4]float64
	if len(values) == 0 {
		return result
	}

	sorted := Sort(values)
	for i := 1; i <= 4; i++ {
		result[i-1] = PercentileOfSorted(sorted, float64(i)*20)
	}
	return result
}

// Decile returns the nth decile (n = 1 to 9).
// Deciles divide data into 10 equal parts.
func Decile(values []float64, n int) float64 {
	if n < 1 || n > 9 {
		return 0
	}
	return Percentile(values, float64(n)*10)
}

// Deciles returns all nine decile boundaries.
func Deciles(values []float64) [9]float64 {
	var result [9]float64
	if len(values) == 0 {
		return result
	}

	sorted := Sort(values)
	for i := 1; i <= 9; i++ {
		result[i-1] = PercentileOfSorted(sorted, float64(i)*10)
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// OUTLIERS
// ════════════════════════════════════════════════════════════════

// OutlierBounds returns the lower and upper bounds for outlier detection.
// Uses the 1.5 * IQR rule.
func OutlierBounds(values []float64) (lower, upper float64) {
	if len(values) == 0 {
		return 0, 0
	}

	sorted := Sort(values)
	q1 := PercentileOfSorted(sorted, 25)
	q3 := PercentileOfSorted(sorted, 75)
	iqr := q3 - q1

	lower = q1 - 1.5*iqr
	upper = q3 + 1.5*iqr
	return
}

// OutlierBoundsExtreme returns bounds using 3 * IQR rule.
func OutlierBoundsExtreme(values []float64) (lower, upper float64) {
	if len(values) == 0 {
		return 0, 0
	}

	sorted := Sort(values)
	q1 := PercentileOfSorted(sorted, 25)
	q3 := PercentileOfSorted(sorted, 75)
	iqr := q3 - q1

	lower = q1 - 3*iqr
	upper = q3 + 3*iqr
	return
}

// IsOutlier checks if a value is an outlier using 1.5 * IQR rule.
func IsOutlier(values []float64, x float64) bool {
	lower, upper := OutlierBounds(values)
	return x < lower || x > upper
}

// Outliers returns all outlier values in the dataset.
func Outliers(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	lower, upper := OutlierBounds(values)

	var outliers []float64
	for _, v := range values {
		if v < lower || v > upper {
			outliers = append(outliers, v)
		}
	}

	return outliers
}

// OutlierCount returns the number of outliers.
func OutlierCount(values []float64) int {
	return len(Outliers(values))
}

// RemoveOutliers returns values with outliers removed.
func RemoveOutliers(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	lower, upper := OutlierBounds(values)

	var result []float64
	for _, v := range values {
		if v >= lower && v <= upper {
			result = append(result, v)
		}
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// FIVE NUMBER SUMMARY
// ════════════════════════════════════════════════════════════════

// FiveNumberSummary holds the five-number summary.
type FiveNumberSummary struct {
	Min    float64
	Q1     float64
	Median float64
	Q3     float64
	Max    float64
}

// FiveNumber returns the five-number summary.
func FiveNumber(values []float64) FiveNumberSummary {
	if len(values) == 0 {
		return FiveNumberSummary{}
	}

	sorted := Sort(values)
	n := len(sorted)

	return FiveNumberSummary{
		Min:    sorted[0],
		Q1:     PercentileOfSorted(sorted, 25),
		Median: PercentileOfSorted(sorted, 50),
		Q3:     PercentileOfSorted(sorted, 75),
		Max:    sorted[n-1],
	}
}

// ════════════════════════════════════════════════════════════════
// BOX PLOT DATA
// ════════════════════════════════════════════════════════════════

// BoxPlotData holds data for rendering a box plot.
type BoxPlotData struct {
	Min          float64
	Q1           float64
	Median       float64
	Q3           float64
	Max          float64
	IQR          float64
	LowerWhisker float64
	UpperWhisker float64
	Outliers     []float64
	OutlierCount int
}

// BoxPlot computes box plot data.
func BoxPlot(values []float64) BoxPlotData {
	if len(values) == 0 {
		return BoxPlotData{}
	}

	sorted := Sort(values)
	n := len(sorted)

	q1 := PercentileOfSorted(sorted, 25)
	q3 := PercentileOfSorted(sorted, 75)
	iqr := q3 - q1

	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr

	// Find whisker endpoints (furthest non-outlier values)
	lowerWhisker := sorted[0]
	upperWhisker := sorted[n-1]

	var outliers []float64
	for _, v := range sorted {
		if v < lowerBound || v > upperBound {
			outliers = append(outliers, v)
		} else {
			if v < lowerWhisker || lowerWhisker < lowerBound {
				lowerWhisker = v
			}
		}
	}

	// Find upper whisker
	for i := n - 1; i >= 0; i-- {
		if sorted[i] <= upperBound {
			upperWhisker = sorted[i]
			break
		}
	}

	// Find lower whisker
	for i := 0; i < n; i++ {
		if sorted[i] >= lowerBound {
			lowerWhisker = sorted[i]
			break
		}
	}

	return BoxPlotData{
		Min:          sorted[0],
		Q1:           q1,
		Median:       PercentileOfSorted(sorted, 50),
		Q3:           q3,
		Max:          sorted[n-1],
		IQR:          iqr,
		LowerWhisker: lowerWhisker,
		UpperWhisker: upperWhisker,
		Outliers:     outliers,
		OutlierCount: len(outliers),
	}
}

// ════════════════════════════════════════════════════════════════
// WINSORIZE
// ════════════════════════════════════════════════════════════════

// Winsorize limits extreme values to the pth and (100-p)th percentiles.
func Winsorize(values []float64, p float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	if p < 0 {
		p = 0
	}
	if p > 50 {
		p = 50
	}

	sorted := Sort(values)
	lower := PercentileOfSorted(sorted, p)
	upper := PercentileOfSorted(sorted, 100-p)

	result := make([]float64, len(values))
	for i, v := range values {
		if v < lower {
			result[i] = lower
		} else if v > upper {
			result[i] = upper
		} else {
			result[i] = v
		}
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// CLAMP
// ════════════════════════════════════════════════════════════════

// Clamp limits values to a min/max range.
func Clamp(values []float64, min, max float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	if min > max {
		min, max = max, min
	}

	result := make([]float64, len(values))
	for i, v := range values {
		if v < min {
			result[i] = min
		} else if v > max {
			result[i] = max
		} else {
			result[i] = v
		}
	}

	return result
}

// ClampToPercentiles limits values to percentile bounds.
func ClampToPercentiles(values []float64, lowerP, upperP float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	sorted := Sort(values)
	lower := PercentileOfSorted(sorted, lowerP)
	upper := PercentileOfSorted(sorted, upperP)

	return Clamp(values, lower, upper)
}

// ════════════════════════════════════════════════════════════════
// NORMALIZE TO PERCENTILE
// ════════════════════════════════════════════════════════════════

// NormalizeToPercentile converts values to their percentile ranks.
func NormalizeToPercentile(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = PercentileRank(values, v)
	}

	return result
}

// Rank returns the rank of each value (1-based).
func Rank(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	n := len(values)

	// Create index slice
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	// Sort indices by values
	sortedValues := Sort(values)

	result := make([]float64, n)
	for i, v := range values {
		// Find rank using binary search
		rank := 1
		for j, sv := range sortedValues {
			if math.Abs(sv-v) < 1e-10 {
				rank = j + 1
				break
			}
		}
		result[i] = float64(rank)
	}

	return result
}
