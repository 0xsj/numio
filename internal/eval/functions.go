// internal/eval/functions.go

package eval

import (
	"strings"

	"github.com/0xsj/numio/internal/fuzzy"
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

// functionNames is a cached list of function names for fuzzy matching.
var functionNames []string

// fuzzyMatcher is the matcher used for function name suggestions.
var fuzzyMatcher *fuzzy.Fuzzy

func init() {
	// Register all functions
	registerCoreFunctions()
	registerMathFunctions()
	registerGraphFunctions()
	registerStatsFunctions()
	registerFinanceFunctions()
	registerDateFunctions()
	registerCalculusFunctions()

	// Build function names list for fuzzy matching
	buildFunctionNamesList()

	// Initialize fuzzy matcher
	fuzzyMatcher = fuzzy.ForFunctionNames()
}

// buildFunctionNamesList builds the cached list of function names.
func buildFunctionNamesList() {
	functionNames = make([]string, 0, len(FunctionRegistry))
	for name := range FunctionRegistry {
		functionNames = append(functionNames, name)
	}
}

// ════════════════════════════════════════════════════════════════
// FUNCTION DISPATCH
// ════════════════════════════════════════════════════════════════

// CallFunction looks up and calls a built-in function by name.
// Note: This only handles built-in functions. User-defined functions
// are handled by the Evaluator which checks them first.
func CallFunction(name string, args []types.Value) types.Value {
	name = strings.ToLower(name)

	fn, ok := FunctionRegistry[name]
	if !ok {
		// Try autocorrect for close matches
		corrected := autocorrectFunctionName(name)
		if corrected != "" && corrected != name {
			fn, ok = FunctionRegistry[corrected]
		}

		if !ok {
			return types.Errorf("unknown function: %s", name)
		}
	}

	// Validate argument count
	if err := validateArgs(fn, args); err.IsError() {
		return err
	}

	return fn.Handler(args)
}

// autocorrectFunctionName attempts to correct a misspelled function name.
// Returns the corrected name if confident, empty string otherwise.
func autocorrectFunctionName(name string) string {
	// Find the closest match
	result := fuzzyMatcher.FindClosest(name, functionNames)

	// Only autocorrect if very confident:
	// - Distance of 1 (single typo), OR
	// - Distance of 2 with high similarity (0.8+), OR
	// - Similarity >= 0.85
	if result.Distance <= 1 {
		return result.Text
	}

	if result.Distance == 2 && result.Similarity >= 0.8 {
		return result.Text
	}

	if result.Similarity >= 0.85 {
		return result.Text
	}

	return ""
}

// HasFunction checks if a built-in function exists.
func HasFunction(name string) bool {
	_, ok := FunctionRegistry[strings.ToLower(name)]
	return ok
}

// GetFunction returns a built-in function definition.
func GetFunction(name string) (FunctionDef, bool) {
	fn, ok := FunctionRegistry[strings.ToLower(name)]
	return fn, ok
}

// ListFunctions returns all registered built-in function names.
func ListFunctions() []string {
	names := make([]string, 0, len(FunctionRegistry))
	for name := range FunctionRegistry {
		names = append(names, name)
	}
	return names
}

// BuiltinFunctionCount returns the number of built-in functions.
func BuiltinFunctionCount() int {
	return len(FunctionRegistry)
}

// SuggestFunction returns function name suggestions for a misspelled name.
func SuggestFunction(name string) []string {
	return fuzzyMatcher.Suggest(strings.ToLower(name), functionNames)
}

// AutocorrectFunction returns the corrected function name if close enough.
func AutocorrectFunction(name string) string {
	name = strings.ToLower(name)

	// If exact match exists, return it
	if HasFunction(name) {
		return name
	}

	// Try autocorrect with distance 1 (very confident)
	return fuzzyMatcher.AutocorrectWithThreshold(name, functionNames, 1)
}

// ════════════════════════════════════════════════════════════════
// RESERVED NAME CHECKING
// ════════════════════════════════════════════════════════════════

// IsReservedFunctionName checks if a name is reserved (built-in function or keyword).
// This is used to prevent user-defined functions from shadowing built-ins.
func IsReservedFunctionName(name string) bool {
	lower := strings.ToLower(name)

	// Check built-in functions
	if HasFunction(lower) {
		return true
	}

	// Check math constants
	if _, ok := GetMathConstant(lower); ok {
		return true
	}

	// Check reserved keywords
	reserved := map[string]bool{
		"def": true, "in": true, "to": true, "of": true,
		"true": true, "false": true, "nil": true, "null": true,
		"if": true, "then": true, "else": true, // Reserved for future
		"and": true, "or": true, "not": true, // Reserved for future
		"for": true, "while": true, "return": true, // Reserved for future
	}

	return reserved[lower]
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

// ════════════════════════════════════════════════════════════════
// FINANCE FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerFinanceFunctions() {
	// Interest
	register("simple", 3, 3, false, FnSimpleInterest)
	register("compound", 3, 4, false, FnCompoundInterest)
	register("apy", 2, 2, false, FnAPY)
	register("apr", 2, 2, false, FnAPR)
	register("rule72", 1, 1, false, FnRule72)
	register("effectiverate", 2, 2, false, FnEffectiveRate)
	register("realrate", 2, 2, false, FnRealRate)

	// Time value of money
	register("pv", 3, 3, false, FnPV)
	register("fv", 3, 3, false, FnFV)
	register("pmt", 3, 3, false, FnPMT)
	register("nper", 3, 3, false, FnNPER)
	register("npv", 2, -1, true, FnNPV)
	register("irr", 2, -1, true, FnIRR)
	register("mirr", 4, -1, true, FnMIRR)

	// Loan
	register("loan", 3, 3, false, FnLoan)
	register("mortgage", 3, 3, false, FnMortgage)
	register("loantotal", 3, 3, false, FnLoanTotal)
	register("loaninterest", 3, 3, false, FnLoanInterest)
	register("loanbalance", 4, 4, false, FnLoanBalance)
	register("maxloan", 3, 3, false, FnMaxLoan)

	// Investment
	register("roi", 2, 2, false, FnROI)
	register("cagr", 3, 3, false, FnCAGR)
	register("payback", 2, 2, false, FnPayback)
	register("grossmargin", 2, 2, false, FnGrossMargin)
	register("netmargin", 2, 2, false, FnNetMargin)
	register("markup", 2, 2, false, FnMarkup)
	register("breakeven", 3, 3, false, FnBreakeven)

	// Depreciation
	register("sln", 3, 3, false, FnSLN)
	register("ddb", 4, 4, false, FnDDB)
	register("syd", 4, 4, false, FnSYD)
	register("macrs", 3, 3, false, FnMACRS)

	// Dividend & valuation
	register("divyield", 2, 2, false, FnDividendYield)
	register("pe", 2, 2, false, FnPE)
	register("eps", 2, 2, false, FnEPS)

	// Tips
	register("tip", 2, 2, false, FnTip)
	register("tipamount", 2, 2, false, FnTipAmount)
	register("splittip", 3, 3, false, FnSplitTip)
	register("percentchange", 2, 2, false, FnPercentChange)
}

// ════════════════════════════════════════════════════════════════
// DATE FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerDateFunctions() {
	// Current date/time
	register("today", 0, 0, false, FnToday)
	register("now", 0, 0, false, FnNow)
	register("yesterday", 0, 0, false, FnYesterday)
	register("tomorrow", 0, 0, false, FnTomorrow)

	// Date creation
	register("date", 3, 3, false, FnDate)
	register("datetime", 5, 6, false, FnDateTime)
	register("time", 2, 3, false, FnTime)
	register("unixtodate", 1, 1, false, FnUnixToDate)
	register("datetounix", 1, 1, false, FnDateToUnix)
	register("epoch", 0, 0, false, FnEpoch)

	// Date arithmetic
	register("adddays", 2, 2, false, FnAddDays)
	register("addweeks", 2, 2, false, FnAddWeeks)
	register("addmonths", 2, 2, false, FnAddMonths)
	register("addyears", 2, 2, false, FnAddYears)
	register("addhours", 2, 2, false, FnAddHours)
	register("addminutes", 2, 2, false, FnAddMinutes)
	register("addseconds", 2, 2, false, FnAddSeconds)
	register("addworkdays", 2, 2, false, FnAddWorkdays)

	// Date difference
	register("daysbetween", 2, 2, false, FnDaysBetween)
	register("weeksbetween", 2, 2, false, FnWeeksBetween)
	register("monthsbetween", 2, 2, false, FnMonthsBetween)
	register("yearsbetween", 2, 2, false, FnYearsBetween)
	register("hoursbetween", 2, 2, false, FnHoursBetween)
	register("minutesbetween", 2, 2, false, FnMinutesBetween)
	register("secondsbetween", 2, 2, false, FnSecondsBetween)
	register("workdaysbetween", 2, 2, false, FnWorkdaysBetween)

	// Date components
	register("year", 1, 1, false, FnYear)
	register("month", 1, 1, false, FnMonth)
	register("day", 1, 1, false, FnDay)
	register("hour", 1, 1, false, FnHour)
	register("minute", 1, 1, false, FnMinute)
	register("second", 1, 1, false, FnSecond)
	register("weekday", 1, 1, false, FnWeekday)
	register("weekdayiso", 1, 1, false, FnWeekdayISO)
	register("weekdayname", 1, 1, false, FnWeekdayName)
	register("monthname", 1, 1, false, FnMonthName)
	register("dayofyear", 1, 1, false, FnDayOfYear)
	register("weekofyear", 1, 1, false, FnWeekOfYear)
	register("quarter", 1, 1, false, FnQuarter)

	// Date utilities
	register("isleapyear", 1, 1, false, FnIsLeapYear)
	register("daysinmonth", 1, 2, false, FnDaysInMonth)
	register("isweekend", 1, 1, false, FnIsWeekend)
	register("isweekday", 1, 1, false, FnIsWeekday)
	register("age", 1, 1, false, FnAge)

	// Date rounding
	register("startofday", 1, 1, false, FnStartOfDay)
	register("endofday", 1, 1, false, FnEndOfDay)
	register("startofmonth", 1, 1, false, FnStartOfMonth)
	register("endofmonth", 1, 1, false, FnEndOfMonth)
	register("startofyear", 1, 1, false, FnStartOfYear)
	register("endofyear", 1, 1, false, FnEndOfYear)
	register("startofweek", 1, 1, false, FnStartOfWeek)
	register("endofweek", 1, 1, false, FnEndOfWeek)

	// Timezone
	register("intimezone", 2, 2, false, FnInTimezone)
	register("toutc", 1, 1, false, FnToUTC)
	register("tolocal", 1, 1, false, FnToLocal)
	register("timezoneoffset", 1, 1, false, FnTimezoneOffset)
	register("nowin", 1, 1, false, FnNowInTimezone)

	// Formatting
	register("formatdate", 2, 2, false, FnFormatDate)
	register("dateiso", 1, 1, false, FnDateISO)

	// Comparison
	register("isbefore", 2, 2, false, FnIsBefore)
	register("isafter", 2, 2, false, FnIsAfter)
	register("issameday", 2, 2, false, FnIsSameDay)

	// Weekday navigation
	register("nextweekday", 2, 2, false, FnNextWeekday)
	register("lastweekday", 2, 2, false, FnLastWeekday)
	register("thisweekday", 2, 2, false, FnThisWeekday)
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

// ════════════════════════════════════════════════════════════════
// CALCULUS FUNCTION REGISTRATION
// ════════════════════════════════════════════════════════════════

func registerCalculusFunctions() {
	// Numerical derivatives
	register("nderivative", 2, 3, false, FnNDerivative)
	register("nderiv", 2, 3, false, FnNDerivative)
	register("nderivative2", 2, 3, false, FnNDerivative2)
	register("nderiv2", 2, 3, false, FnNDerivative2)

	// Numerical integration
	register("nintegral", 3, 4, false, FnNIntegral)
	register("nintegrate", 3, 4, false, FnNIntegral)

	// Limit
	register("nlimit", 2, 3, false, FnNLimit)

	// Root finding
	register("findroot", 3, 4, false, FnFindRoot)
	register("bisect", 3, 4, false, FnFindRoot)

	// Extrema
	register("findmin", 3, 4, false, FnFindMin)
	register("findmax", 3, 4, false, FnFindMax)

	// Summation and product
	register("summation", 3, 3, false, FnSummation)
	register("sigma", 3, 3, false, FnSummation)
	register("seriesproduct", 3, 3, false, FnSeriesProduct)
	register("cappi", 3, 3, false, FnSeriesProduct)

	// Taylor series
	register("taylorsin", 1, 2, false, FnTaylorSin)
	register("taylorcos", 1, 2, false, FnTaylorCos)
	register("taylorexp", 1, 2, false, FnTaylorExp)
	register("taylorln", 1, 2, false, FnTaylorLn)

	// Riemann sums
	register("riemannleft", 3, 4, false, FnRiemannLeft)
	register("riemannright", 3, 4, false, FnRiemannRight)
	register("riemannmid", 3, 4, false, FnRiemannMid)
	register("trapezoidal", 3, 4, false, FnTrapezoidal)
	register("trapezoid", 3, 4, false, FnTrapezoidal)
}
