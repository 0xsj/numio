// internal/eval/functions_stats.go

package eval

import (
	"github.com/0xsj/numio/internal/stats"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// CENTRAL TENDENCY
// ════════════════════════════════════════════════════════════════

// FnMedian returns the median of the arguments.
func FnMedian(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("median requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Median(values))
}

// FnMode returns the most frequent value.
func FnMode(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("mode requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Mode(values))
}

// FnGeomean returns the geometric mean.
func FnGeomean(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("geomean requires at least one argument")
	}
	values := valuesToFloats(args)
	result, err := stats.GeometricMean(values)
	if err != nil {
		return types.Error("geomean: all values must be positive")
	}
	return types.Number(result)
}

// FnHarmean returns the harmonic mean.
func FnHarmean(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("harmean requires at least one argument")
	}
	values := valuesToFloats(args)
	result, err := stats.HarmonicMean(values)
	if err != nil {
		return types.Error("harmean: cannot include zero values")
	}
	return types.Number(result)
}

// FnTrimmedMean returns the trimmed mean.
func FnTrimmedMean(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("trimmean requires at least 2 arguments: trim% and values")
	}
	trimPercent := args[0].AsFloat()
	values := valuesToFloats(args[1:])
	return types.Number(stats.TrimmedMean(values, trimPercent))
}

// FnMidrange returns (max + min) / 2.
func FnMidrange(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("midrange requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Midrange(values))
}

// ════════════════════════════════════════════════════════════════
// DISPERSION
// ════════════════════════════════════════════════════════════════

// FnVariance returns the population variance.
func FnVariance(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("variance requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Variance(values))
}

// FnSampleVariance returns the sample variance.
func FnSampleVariance(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("svariance requires at least two arguments")
	}
	values := valuesToFloats(args)
	return types.Number(stats.SampleVariance(values))
}

// FnStdDev returns the population standard deviation.
func FnStdDev(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("stddev requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.StdDev(values))
}

// FnSampleStdDev returns the sample standard deviation.
func FnSampleStdDev(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("sstddev requires at least two arguments")
	}
	values := valuesToFloats(args)
	return types.Number(stats.SampleStdDev(values))
}

// FnStatRange returns max - min.
func FnStatRange(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("range requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Range(values))
}

// FnMAD returns the mean absolute deviation.
func FnMAD(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("mad requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.MAD(values))
}

// FnRMS returns the root mean square.
func FnRMS(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("rms requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.RMS(values))
}

// FnCV returns the coefficient of variation (as percentage).
func FnCV(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("cv requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.CoefficientOfVariationPercent(values))
}

// FnSkewness returns the skewness.
func FnSkewness(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("skewness requires at least 3 arguments")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Skewness(values))
}

// FnKurtosis returns the excess kurtosis.
func FnKurtosis(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("kurtosis requires at least 4 arguments")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Kurtosis(values))
}

// FnStdErr returns the standard error of the mean.
func FnStdErr(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("stderr requires at least 2 arguments")
	}
	values := valuesToFloats(args)
	return types.Number(stats.StandardError(values))
}

// ════════════════════════════════════════════════════════════════
// PERCENTILES & QUARTILES
// ════════════════════════════════════════════════════════════════

// FnPercentile returns the pth percentile.
func FnPercentile(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("percentile requires at least 2 arguments: p and values")
	}
	p := args[0].AsFloat()
	if p < 0 || p > 100 {
		return types.Error("percentile: p must be between 0 and 100")
	}
	values := valuesToFloats(args[1:])
	return types.Number(stats.Percentile(values, p))
}

// FnQuartile returns Q1, Q2, or Q3.
func FnQuartile(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("quartile requires at least 2 arguments: q and values")
	}
	q := int(args[0].AsFloat())
	if q < 1 || q > 3 {
		return types.Error("quartile: q must be 1, 2, or 3")
	}
	values := valuesToFloats(args[1:])
	return types.Number(stats.Quartile(values, q))
}

// FnQ1 returns the first quartile.
func FnQ1(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("q1 requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Q1(values))
}

// FnQ2 returns the second quartile (median).
func FnQ2(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("q2 requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Q2(values))
}

// FnQ3 returns the third quartile.
func FnQ3(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("q3 requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.Q3(values))
}

// FnIQR returns the interquartile range.
func FnIQR(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("iqr requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(stats.IQR(values))
}

// FnDecile returns the nth decile.
func FnDecile(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("decile requires at least 2 arguments: n and values")
	}
	n := int(args[0].AsFloat())
	if n < 1 || n > 9 {
		return types.Error("decile: n must be between 1 and 9")
	}
	values := valuesToFloats(args[1:])
	return types.Number(stats.Decile(values, n))
}

// FnOutlierCount returns the number of outliers.
func FnOutlierCount(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("outliers requires at least one argument")
	}
	values := valuesToFloats(args)
	return types.Number(float64(stats.OutlierCount(values)))
}

// ════════════════════════════════════════════════════════════════
// CORRELATION & REGRESSION
// ════════════════════════════════════════════════════════════════

// FnZScore returns the z-score of a value in a dataset.
func FnZScore(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("zscore requires at least 2 arguments: x and dataset")
	}
	x := args[0].AsFloat()
	values := valuesToFloats(args[1:])
	result, err := stats.ZScore(x, values)
	if err != nil {
		return types.Error("zscore: standard deviation is zero")
	}
	return types.Number(result)
}

// FnCovariance returns the covariance of x,y pairs.
func FnCovariance(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("covariance requires at least 4 arguments (x,y pairs)")
	}
	if len(args)%2 != 0 {
		return types.Error("covariance requires even number of arguments (x,y pairs)")
	}
	pairs := valuesToFloats(args)
	result, err := stats.CovarianceFromPairs(pairs)
	if err != nil {
		return types.Errorf("covariance: %s", err.Error())
	}
	return types.Number(result)
}

// FnCorrelation returns the Pearson correlation coefficient.
func FnCorrelation(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("correlation requires at least 4 arguments (x,y pairs)")
	}
	if len(args)%2 != 0 {
		return types.Error("correlation requires even number of arguments (x,y pairs)")
	}
	pairs := valuesToFloats(args)
	result, err := stats.CorrelationFromPairs(pairs)
	if err != nil {
		return types.Errorf("correlation: %s", err.Error())
	}
	return types.Number(result)
}

// FnRSquared returns the coefficient of determination.
func FnRSquared(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("rsquared requires at least 4 arguments (x,y pairs)")
	}
	if len(args)%2 != 0 {
		return types.Error("rsquared requires even number of arguments (x,y pairs)")
	}

	n := len(args) / 2
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := 0; i < n; i++ {
		xs[i] = args[i*2].AsFloat()
		ys[i] = args[i*2+1].AsFloat()
	}

	result, err := stats.RSquared(xs, ys)
	if err != nil {
		return types.Errorf("rsquared: %s", err.Error())
	}
	return types.Number(result)
}

// FnLinReg returns the slope of a linear regression.
func FnLinReg(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("linreg requires at least 4 arguments (x,y pairs)")
	}
	if len(args)%2 != 0 {
		return types.Error("linreg requires even number of arguments (x,y pairs)")
	}

	pairs := valuesToFloats(args)
	result, err := stats.LinearRegressionFromPairs(pairs)
	if err != nil {
		return types.Errorf("linreg: %s", err.Error())
	}

	// Return slope (most commonly needed)
	return types.Number(result.Slope)
}

// FnLinRegIntercept returns the intercept of a linear regression.
func FnLinRegIntercept(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("linregb requires at least 4 arguments (x,y pairs)")
	}
	if len(args)%2 != 0 {
		return types.Error("linregb requires even number of arguments (x,y pairs)")
	}

	pairs := valuesToFloats(args)
	result, err := stats.LinearRegressionFromPairs(pairs)
	if err != nil {
		return types.Errorf("linregb: %s", err.Error())
	}

	return types.Number(result.Intercept)
}

// FnSpearman returns Spearman's rank correlation.
func FnSpearman(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("spearman requires at least 4 arguments (x,y pairs)")
	}
	if len(args)%2 != 0 {
		return types.Error("spearman requires even number of arguments (x,y pairs)")
	}

	n := len(args) / 2
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := 0; i < n; i++ {
		xs[i] = args[i*2].AsFloat()
		ys[i] = args[i*2+1].AsFloat()
	}

	result, err := stats.SpearmanCorrelation(xs, ys)
	if err != nil {
		return types.Errorf("spearman: %s", err.Error())
	}
	return types.Number(result)
}

// FnAutocorr returns the autocorrelation at a given lag.
func FnAutocorr(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("autocorr requires at least 2 arguments: lag and values")
	}
	lag := int(args[0].AsFloat())
	values := valuesToFloats(args[1:])
	return types.Number(stats.Autocorrelation(values, lag))
}

// ════════════════════════════════════════════════════════════════
// OTHER
// ════════════════════════════════════════════════════════════════

// FnProduct returns the product of all values.
func FnProduct(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Number(1)
	}
	values := valuesToFloats(args)
	return types.Number(stats.Product(values))
}
