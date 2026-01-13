// internal/eval/functions.go

package eval

import (
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// FunctionHandler is the signature for function implementations.
type FunctionHandler func(args []types.Value) types.Value

// FunctionDef defines a function's metadata and handler.
type FunctionDef struct {
	Name     string
	MinArgs  int
	MaxArgs  int // -1 for variadic
	Variadic bool
	Handler  FunctionHandler
}

// FunctionRegistry holds all registered functions.
var FunctionRegistry = map[string]FunctionDef{}

func init() {
	// Register all functions
	registerCoreFunctions()
	registerMathFunctions()
	registerGraphFunctions()
	registerStatsFunctions()
}

// ════════════════════════════════════════════════════════════════
// FUNCTION DISPATCH
// ════════════════════════════════════════════════════════════════

// CallFunction looks up and calls a function by name.
func CallFunction(name string, args []types.Value) types.Value {
	name = strings.ToLower(name)

	fn, ok := FunctionRegistry[name]
	if !ok {
		return types.Errorf("unknown function: %s", name)
	}

	// Validate argument count
	if err := validateArgs(fn, args); err.IsError() {
		return err
	}

	return fn.Handler(args)
}

// HasFunction checks if a function exists.
func HasFunction(name string) bool {
	_, ok := FunctionRegistry[strings.ToLower(name)]
	return ok
}

// GetFunction returns a function definition.
func GetFunction(name string) (FunctionDef, bool) {
	fn, ok := FunctionRegistry[strings.ToLower(name)]
	return fn, ok
}

// ListFunctions returns all registered function names.
func ListFunctions() []string {
	names := make([]string, 0, len(FunctionRegistry))
	for name := range FunctionRegistry {
		names = append(names, name)
	}
	return names
}

// ════════════════════════════════════════════════════════════════
// ARGUMENT VALIDATION
// ════════════════════════════════════════════════════════════════

func validateArgs(fn FunctionDef, args []types.Value) types.Value {
	argc := len(args)

	if argc < fn.MinArgs {
		if fn.MinArgs == fn.MaxArgs {
			return types.Errorf("%s requires exactly %d argument(s), got %d", fn.Name, fn.MinArgs, argc)
		}
		return types.Errorf("%s requires at least %d argument(s), got %d", fn.Name, fn.MinArgs, argc)
	}

	if !fn.Variadic && fn.MaxArgs >= 0 && argc > fn.MaxArgs {
		if fn.MinArgs == fn.MaxArgs {
			return types.Errorf("%s requires exactly %d argument(s), got %d", fn.Name, fn.MaxArgs, argc)
		}
		return types.Errorf("%s requires at most %d argument(s), got %d", fn.Name, fn.MaxArgs, argc)
	}

	// Return empty value (not error) on success
	return types.Empty()
}

// ════════════════════════════════════════════════════════════════
// CORE FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerCoreFunctions() {
	// Aggregation
	register("sum", 0, -1, true, FnSum)
	register("avg", 1, -1, true, FnAvg)
	register("average", 1, -1, true, FnAvg)
	register("mean", 1, -1, true, FnAvg)
	register("min", 1, -1, true, FnMin)
	register("max", 1, -1, true, FnMax)
	register("count", 0, -1, true, FnCount)

	// Basic math
	register("abs", 1, 1, false, FnAbs)
	register("sqrt", 1, 1, false, FnSqrt)
	register("cbrt", 1, 1, false, FnCbrt)
	register("round", 1, 2, false, FnRound)
	register("floor", 1, 1, false, FnFloor)
	register("ceil", 1, 1, false, FnCeil)
	register("pow", 2, 2, false, FnPow)
}

// ════════════════════════════════════════════════════════════════
// MATH FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerMathFunctions() {
	// Logarithms & exponentials
	register("log", 1, 1, false, FnLog)
	register("ln", 1, 1, false, FnLog)
	register("log10", 1, 1, false, FnLog10)
	register("log2", 1, 1, false, FnLog2)
	register("exp", 1, 1, false, FnExp)

	// Trigonometric
	register("sin", 1, 1, false, FnSin)
	register("cos", 1, 1, false, FnCos)
	register("tan", 1, 1, false, FnTan)
	register("asin", 1, 1, false, FnAsin)
	register("acos", 1, 1, false, FnAcos)
	register("atan", 1, 1, false, FnAtan)
	register("atan2", 2, 2, false, FnAtan2)

	// Hyperbolic
	register("sinh", 1, 1, false, FnSinh)
	register("cosh", 1, 1, false, FnCosh)
	register("tanh", 1, 1, false, FnTanh)
	register("asinh", 1, 1, false, FnAsinh)
	register("acosh", 1, 1, false, FnAcosh)
	register("atanh", 1, 1, false, FnAtanh)

	// Angle conversion
	register("deg", 1, 1, false, FnDeg)
	register("rad", 1, 1, false, FnRad)

	// Combinatorics
	register("factorial", 1, 1, false, FnFactorial)
	register("npr", 2, 2, false, FnPermutations)
	register("ncr", 2, 2, false, FnCombinations)
	register("perm", 2, 2, false, FnPermutations)
	register("comb", 2, 2, false, FnCombinations)

	// Number theory
	register("gcd", 2, -1, true, FnGCD)
	register("lcm", 2, -1, true, FnLCM)
	register("mod", 2, 2, false, FnMod)

	// Rounding & sign
	register("sign", 1, 1, false, FnSign)
	register("trunc", 1, 1, false, FnTrunc)
	register("frac", 1, 1, false, FnFrac)

	// Special
	register("hypot", 2, 2, false, FnHypot)
}

// ════════════════════════════════════════════════════════════════
// GRAPH FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerGraphFunctions() {
	// Sparklines
	register("spark", 1, -1, true, FnSparkline)
	register("sparkline", 1, -1, true, FnSparkline)
	register("sparkstats", 1, -1, true, FnSparklineStats)
	register("sparktrend", 1, -1, true, FnSparklineTrend)
	register("sparkstyle", 2, -1, true, FnSparklineStyled)
	register("sparkbound", 3, -1, true, FnSparklineBounded)
	register("sparkwidth", 2, -1, true, FnSparklineFixed)

	// Histograms
	register("hist", 1, -1, true, FnHistogram)
	register("histogram", 1, -1, true, FnHistogram)
	register("histbins", 2, -1, true, FnHistogramBins)

	// Gauges
	register("gauge", 1, 2, false, FnGauge)
	register("gaugerange", 3, 3, false, FnGaugeRange)
	register("battery", 1, 1, false, FnBattery)
	register("meter", 1, 2, false, FnMeter)
	register("signal", 1, 1, false, FnSignal)
	register("stars", 1, 2, false, FnStars)
	register("hearts", 1, 2, false, FnHearts)

	// Bars
	register("bar", 1, 2, false, FnBar)
	register("progress", 1, 2, false, FnProgress)

	// Trends
	register("trend", 2, -1, true, FnTrend)
	register("trendstyle", 3, -1, true, FnTrendStyled)
	register("change", 2, -1, true, FnChange)
	register("minitrend", 2, -1, true, FnMiniTrend)
	register("slope", 2, -1, true, FnSlope)
	register("delta", 2, 2, false, FnDelta)
	register("deltapct", 2, 2, false, FnDeltaPercent)

	// Dot plots
	register("dots", 1, -1, true, FnDotPlot)
	register("dotplot", 1, -1, true, FnDotPlot)
	register("distdots", 1, -1, true, FnDistDots)
	register("numline", 3, 3, false, FnNumberLine)
	register("scatter", 4, -1, true, FnScatter)

	// Comparison
	register("compare", 2, 2, false, FnCompare)
}

// ════════════════════════════════════════════════════════════════
// STATS FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerStatsFunctions() {
	// Central tendency
	register("median", 1, -1, true, FnMedian)
	register("mode", 1, -1, true, FnMode)
	register("geomean", 1, -1, true, FnGeomean)
	register("harmean", 1, -1, true, FnHarmean)
	register("trimmean", 2, -1, true, FnTrimmedMean)
	register("midrange", 1, -1, true, FnMidrange)

	// Dispersion
	register("variance", 1, -1, true, FnVariance)
	register("var", 1, -1, true, FnVariance)
	register("svariance", 2, -1, true, FnSampleVariance)
	register("stddev", 1, -1, true, FnStdDev)
	register("stdev", 1, -1, true, FnStdDev)
	register("sstddev", 2, -1, true, FnSampleStdDev)
	register("sstdev", 2, -1, true, FnSampleStdDev)
	register("range", 1, -1, true, FnStatRange)
	register("mad", 1, -1, true, FnMAD)
	register("rms", 1, -1, true, FnRMS)
	register("cv", 1, -1, true, FnCV)
	register("skewness", 3, -1, true, FnSkewness)
	register("kurtosis", 4, -1, true, FnKurtosis)
	register("stderr", 2, -1, true, FnStdErr)

	// Percentiles & quartiles
	register("percentile", 2, -1, true, FnPercentile)
	register("pctl", 2, -1, true, FnPercentile)
	register("quartile", 2, -1, true, FnQuartile)
	register("q1", 1, -1, true, FnQ1)
	register("q2", 1, -1, true, FnQ2)
	register("q3", 1, -1, true, FnQ3)
	register("iqr", 1, -1, true, FnIQR)
	register("decile", 2, -1, true, FnDecile)
	register("outliers", 1, -1, true, FnOutlierCount)

	// Correlation & regression
	register("zscore", 2, -1, true, FnZScore)
	register("covariance", 4, -1, true, FnCovariance)
	register("cov", 4, -1, true, FnCovariance)
	register("correlation", 4, -1, true, FnCorrelation)
	register("corr", 4, -1, true, FnCorrelation)
	register("rsquared", 4, -1, true, FnRSquared)
	register("r2", 4, -1, true, FnRSquared)
	register("linreg", 4, -1, true, FnLinReg)
	register("linregb", 4, -1, true, FnLinRegIntercept)
	register("spearman", 4, -1, true, FnSpearman)
	register("autocorr", 2, -1, true, FnAutocorr)

	// Other
	register("product", 1, -1, true, FnProduct)
	register("prod", 1, -1, true, FnProduct)
}

// register is a helper to add functions to the registry.
func register(name string, minArgs, maxArgs int, variadic bool, handler FunctionHandler) {
	FunctionRegistry[name] = FunctionDef{
		Name:     name,
		MinArgs:  minArgs,
		MaxArgs:  maxArgs,
		Variadic: variadic,
		Handler:  handler,
	}
}
