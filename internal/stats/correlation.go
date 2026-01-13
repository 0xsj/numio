// internal/stats/correlation.go

package stats

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// Z-SCORE
// ════════════════════════════════════════════════════════════════

// ZScore returns the z-score of a value given the dataset.
func ZScore(x float64, values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	mean := Mean(values)
	stddev := StdDev(values)

	if stddev == 0 {
		return 0, ErrZeroStdDev
	}

	return (x - mean) / stddev, nil
}

// ZScoreOrZero returns z-score or 0 on error.
func ZScoreOrZero(x float64, values []float64) float64 {
	result, err := ZScore(x, values)
	if err != nil {
		return 0
	}
	return result
}

// ZScores returns z-scores for all values.
func ZScores(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}

	mean := Mean(values)
	stddev := StdDev(values)

	if stddev == 0 {
		result := make([]float64, len(values))
		return result
	}

	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = (v - mean) / stddev
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// COVARIANCE
// ════════════════════════════════════════════════════════════════

// Covariance returns the population covariance of two datasets.
func Covariance(xs, ys []float64) (float64, error) {
	if len(xs) != len(ys) {
		return 0, ErrLengthMismatch
	}

	if len(xs) == 0 {
		return 0, nil
	}

	meanX := Mean(xs)
	meanY := Mean(ys)

	var cov float64
	for i := range xs {
		cov += (xs[i] - meanX) * (ys[i] - meanY)
	}

	return cov / float64(len(xs)), nil
}

// SampleCovariance returns the sample covariance (n-1 denominator).
func SampleCovariance(xs, ys []float64) (float64, error) {
	if len(xs) != len(ys) {
		return 0, ErrLengthMismatch
	}

	if len(xs) < 2 {
		return 0, nil
	}

	meanX := Mean(xs)
	meanY := Mean(ys)

	var cov float64
	for i := range xs {
		cov += (xs[i] - meanX) * (ys[i] - meanY)
	}

	return cov / float64(len(xs)-1), nil
}

// CovarianceFromPairs computes covariance from interleaved x,y pairs.
func CovarianceFromPairs(pairs []float64) (float64, error) {
	if len(pairs)%2 != 0 {
		return 0, ErrOddPairCount
	}

	n := len(pairs) / 2
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i := 0; i < n; i++ {
		xs[i] = pairs[i*2]
		ys[i] = pairs[i*2+1]
	}

	return Covariance(xs, ys)
}

// ════════════════════════════════════════════════════════════════
// CORRELATION
// ════════════════════════════════════════════════════════════════

// Correlation returns the Pearson correlation coefficient.
func Correlation(xs, ys []float64) (float64, error) {
	if len(xs) != len(ys) {
		return 0, ErrLengthMismatch
	}

	if len(xs) == 0 {
		return 0, nil
	}

	meanX := Mean(xs)
	meanY := Mean(ys)

	var cov, varX, varY float64
	for i := range xs {
		dx := xs[i] - meanX
		dy := ys[i] - meanY
		cov += dx * dy
		varX += dx * dx
		varY += dy * dy
	}

	if varX == 0 || varY == 0 {
		return 0, ErrZeroVariance
	}

	return cov / math.Sqrt(varX*varY), nil
}

// CorrelationFromPairs computes correlation from interleaved x,y pairs.
func CorrelationFromPairs(pairs []float64) (float64, error) {
	if len(pairs)%2 != 0 {
		return 0, ErrOddPairCount
	}

	n := len(pairs) / 2
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i := 0; i < n; i++ {
		xs[i] = pairs[i*2]
		ys[i] = pairs[i*2+1]
	}

	return Correlation(xs, ys)
}

// ════════════════════════════════════════════════════════════════
// SPEARMAN RANK CORRELATION
// ════════════════════════════════════════════════════════════════

// SpearmanCorrelation returns Spearman's rank correlation coefficient.
func SpearmanCorrelation(xs, ys []float64) (float64, error) {
	if len(xs) != len(ys) {
		return 0, ErrLengthMismatch
	}

	if len(xs) < 2 {
		return 0, nil
	}

	// Convert to ranks
	ranksX := calculateRanks(xs)
	ranksY := calculateRanks(ys)

	// Calculate Pearson correlation on ranks
	return Correlation(ranksX, ranksY)
}

// calculateRanks computes average ranks for values.
func calculateRanks(values []float64) []float64 {
	n := len(values)
	if n == 0 {
		return nil
	}

	// Create indices sorted by values
	type indexedValue struct {
		index int
		value float64
	}

	indexed := make([]indexedValue, n)
	for i, v := range values {
		indexed[i] = indexedValue{i, v}
	}

	// Sort by value
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if indexed[i].value > indexed[j].value {
				indexed[i], indexed[j] = indexed[j], indexed[i]
			}
		}
	}

	// Assign ranks (handling ties with average rank)
	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i
		// Find all tied values
		for j < n && indexed[j].value == indexed[i].value {
			j++
		}
		// Average rank for tied values
		avgRank := float64(i+j+1) / 2
		for k := i; k < j; k++ {
			ranks[indexed[k].index] = avgRank
		}
		i = j
	}

	return ranks
}

// ════════════════════════════════════════════════════════════════
// R-SQUARED
// ════════════════════════════════════════════════════════════════

// RSquared returns the coefficient of determination (r²).
func RSquared(xs, ys []float64) (float64, error) {
	r, err := Correlation(xs, ys)
	if err != nil {
		return 0, err
	}
	return r * r, nil
}

// ════════════════════════════════════════════════════════════════
// LINEAR REGRESSION
// ════════════════════════════════════════════════════════════════

// LinearRegression returns slope and intercept for y = mx + b.
type RegressionResult struct {
	Slope     float64
	Intercept float64
	RSquared  float64
}

// LinearRegression computes simple linear regression.
func LinearRegression(xs, ys []float64) (RegressionResult, error) {
	if len(xs) != len(ys) {
		return RegressionResult{}, ErrLengthMismatch
	}

	if len(xs) < 2 {
		return RegressionResult{}, nil
	}

	meanX := Mean(xs)
	meanY := Mean(ys)

	var sumXY, sumXX float64
	for i := range xs {
		dx := xs[i] - meanX
		sumXY += dx * (ys[i] - meanY)
		sumXX += dx * dx
	}

	if sumXX == 0 {
		return RegressionResult{}, ErrZeroVariance
	}

	slope := sumXY / sumXX
	intercept := meanY - slope*meanX

	// Calculate R²
	var ssRes, ssTot float64
	for i := range xs {
		predicted := slope*xs[i] + intercept
		ssRes += (ys[i] - predicted) * (ys[i] - predicted)
		ssTot += (ys[i] - meanY) * (ys[i] - meanY)
	}

	var rSquared float64
	if ssTot != 0 {
		rSquared = 1 - ssRes/ssTot
	}

	return RegressionResult{
		Slope:     slope,
		Intercept: intercept,
		RSquared:  rSquared,
	}, nil
}

// LinearRegressionFromPairs computes regression from interleaved pairs.
func LinearRegressionFromPairs(pairs []float64) (RegressionResult, error) {
	if len(pairs)%2 != 0 {
		return RegressionResult{}, ErrOddPairCount
	}

	n := len(pairs) / 2
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i := 0; i < n; i++ {
		xs[i] = pairs[i*2]
		ys[i] = pairs[i*2+1]
	}

	return LinearRegression(xs, ys)
}

// Predict returns the predicted y value for a given x.
func (r RegressionResult) Predict(x float64) float64 {
	return r.Slope*x + r.Intercept
}

// ════════════════════════════════════════════════════════════════
// AUTOCORRELATION
// ════════════════════════════════════════════════════════════════

// Autocorrelation returns the autocorrelation at a given lag.
func Autocorrelation(values []float64, lag int) float64 {
	n := len(values)
	if n == 0 || lag >= n || lag < 0 {
		return 0
	}

	mean := Mean(values)

	var numerator, denominator float64
	for i := 0; i < n; i++ {
		denominator += (values[i] - mean) * (values[i] - mean)
	}

	if denominator == 0 {
		return 0
	}

	for i := 0; i < n-lag; i++ {
		numerator += (values[i] - mean) * (values[i+lag] - mean)
	}

	return numerator / denominator
}

// AutocorrelationSeries returns autocorrelation for lags 0 to maxLag.
func AutocorrelationSeries(values []float64, maxLag int) []float64 {
	if maxLag < 0 {
		maxLag = len(values) / 4
	}
	if maxLag >= len(values) {
		maxLag = len(values) - 1
	}

	result := make([]float64, maxLag+1)
	for lag := 0; lag <= maxLag; lag++ {
		result[lag] = Autocorrelation(values, lag)
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// ADDITIONAL ERRORS
// ════════════════════════════════════════════════════════════════

const (
	ErrZeroStdDev   StatsError = "standard deviation is zero"
	ErrZeroVariance StatsError = "variance is zero"
	ErrOddPairCount StatsError = "odd number of values (expected x,y pairs)"
)
