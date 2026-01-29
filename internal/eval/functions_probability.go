// internal/eval/functions_probability.go

package eval

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/0xsj/numio/pkg/types"
)

// init seeds the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// ════════════════════════════════════════════════════════════════
// BINOMIAL DISTRIBUTION
// ════════════════════════════════════════════════════════════════

// FnBinomialPDF calculates the binomial probability mass function.
// P(X = k) for n trials with probability p.
// Args: n (trials), p (probability), k (successes)
func FnBinomialPDF(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("binomialpdf requires 3 arguments: n, p, k")
	}

	n := int(args[0].AsFloat())
	p := args[1].AsFloat()
	k := int(args[2].AsFloat())

	if n < 0 {
		return types.Error("binomialpdf: n must be non-negative")
	}
	if p < 0 || p > 1 {
		return types.Error("binomialpdf: p must be between 0 and 1")
	}
	if k < 0 || k > n {
		return types.Number(0)
	}

	// P(X = k) = C(n,k) * p^k * (1-p)^(n-k)
	coeff := binomialCoeff(n, k)
	prob := coeff * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))

	return types.Number(prob)
}

// FnBinomialCDF calculates the binomial cumulative distribution function.
// P(X <= k) for n trials with probability p.
// Args: n (trials), p (probability), k (successes)
func FnBinomialCDF(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("binomialcdf requires 3 arguments: n, p, k")
	}

	n := int(args[0].AsFloat())
	p := args[1].AsFloat()
	k := int(args[2].AsFloat())

	if n < 0 {
		return types.Error("binomialcdf: n must be non-negative")
	}
	if p < 0 || p > 1 {
		return types.Error("binomialcdf: p must be between 0 and 1")
	}
	if k < 0 {
		return types.Number(0)
	}
	if k >= n {
		return types.Number(1)
	}

	// Sum P(X = i) for i = 0 to k
	var cdf float64
	for i := 0; i <= k; i++ {
		coeff := binomialCoeff(n, i)
		cdf += coeff * math.Pow(p, float64(i)) * math.Pow(1-p, float64(n-i))
	}

	return types.Number(cdf)
}

// FnBinomialMean returns the mean of a binomial distribution.
// Args: n (trials), p (probability)
func FnBinomialMean(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("binomialmean requires 2 arguments: n, p")
	}

	n := args[0].AsFloat()
	p := args[1].AsFloat()

	if n < 0 {
		return types.Error("binomialmean: n must be non-negative")
	}
	if p < 0 || p > 1 {
		return types.Error("binomialmean: p must be between 0 and 1")
	}

	return types.Number(n * p)
}

// FnBinomialVar returns the variance of a binomial distribution.
// Args: n (trials), p (probability)
func FnBinomialVar(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("binomialvar requires 2 arguments: n, p")
	}

	n := args[0].AsFloat()
	p := args[1].AsFloat()

	if n < 0 {
		return types.Error("binomialvar: n must be non-negative")
	}
	if p < 0 || p > 1 {
		return types.Error("binomialvar: p must be between 0 and 1")
	}

	return types.Number(n * p * (1 - p))
}

// ════════════════════════════════════════════════════════════════
// POISSON DISTRIBUTION
// ════════════════════════════════════════════════════════════════

// FnPoissonPDF calculates the Poisson probability mass function.
// P(X = k) for rate λ.
// Args: lambda (rate), k (occurrences)
func FnPoissonPDF(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("poissonpdf requires 2 arguments: lambda, k")
	}

	lambda := args[0].AsFloat()
	k := int(args[1].AsFloat())

	if lambda < 0 {
		return types.Error("poissonpdf: lambda must be non-negative")
	}
	if k < 0 {
		return types.Number(0)
	}

	// P(X = k) = (λ^k * e^(-λ)) / k!
	prob := math.Pow(lambda, float64(k)) * math.Exp(-lambda) / factorial64(k)

	return types.Number(prob)
}

// FnPoissonCDF calculates the Poisson cumulative distribution function.
// P(X <= k) for rate λ.
// Args: lambda (rate), k (occurrences)
func FnPoissonCDF(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("poissoncdf requires 2 arguments: lambda, k")
	}

	lambda := args[0].AsFloat()
	k := int(args[1].AsFloat())

	if lambda < 0 {
		return types.Error("poissoncdf: lambda must be non-negative")
	}
	if k < 0 {
		return types.Number(0)
	}

	// Sum P(X = i) for i = 0 to k
	var cdf float64
	for i := 0; i <= k; i++ {
		cdf += math.Pow(lambda, float64(i)) * math.Exp(-lambda) / factorial64(i)
	}

	return types.Number(cdf)
}

// FnPoissonMean returns the mean of a Poisson distribution (equals λ).
func FnPoissonMean(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("poissonmean requires 1 argument: lambda")
	}

	lambda := args[0].AsFloat()
	if lambda < 0 {
		return types.Error("poissonmean: lambda must be non-negative")
	}

	return types.Number(lambda)
}

// FnPoissonVar returns the variance of a Poisson distribution (equals λ).
func FnPoissonVar(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("poissonvar requires 1 argument: lambda")
	}

	lambda := args[0].AsFloat()
	if lambda < 0 {
		return types.Error("poissonvar: lambda must be non-negative")
	}

	return types.Number(lambda)
}

// ════════════════════════════════════════════════════════════════
// NORMAL DISTRIBUTION
// ════════════════════════════════════════════════════════════════

// FnNormalPDF calculates the normal probability density function.
// Args: x, [mean], [stddev] (defaults: mean=0, stddev=1)
func FnNormalPDF(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 3 {
		return types.Error("normalpdf requires 1-3 arguments: x, [mean], [stddev]")
	}

	x := args[0].AsFloat()
	mean := 0.0
	stddev := 1.0

	if len(args) >= 2 {
		mean = args[1].AsFloat()
	}
	if len(args) >= 3 {
		stddev = args[2].AsFloat()
		if stddev <= 0 {
			return types.Error("normalpdf: stddev must be positive")
		}
	}

	// f(x) = (1 / (σ√(2π))) * e^(-(x-μ)²/(2σ²))
	z := (x - mean) / stddev
	pdf := math.Exp(-0.5*z*z) / (stddev * math.Sqrt(2*math.Pi))

	return types.Number(pdf)
}

// FnNormalCDF calculates the normal cumulative distribution function.
// Args: x, [mean], [stddev] (defaults: mean=0, stddev=1)
func FnNormalCDF(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 3 {
		return types.Error("normalcdf requires 1-3 arguments: x, [mean], [stddev]")
	}

	x := args[0].AsFloat()
	mean := 0.0
	stddev := 1.0

	if len(args) >= 2 {
		mean = args[1].AsFloat()
	}
	if len(args) >= 3 {
		stddev = args[2].AsFloat()
		if stddev <= 0 {
			return types.Error("normalcdf: stddev must be positive")
		}
	}

	// Use error function: Φ(z) = 0.5 * (1 + erf(z/√2))
	z := (x - mean) / stddev
	cdf := 0.5 * (1 + math.Erf(z/math.Sqrt2))

	return types.Number(cdf)
}

// FnNormalInv calculates the inverse normal (quantile function).
// Args: p (probability), [mean], [stddev]
func FnNormalInv(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 3 {
		return types.Error("norminv requires 1-3 arguments: p, [mean], [stddev]")
	}

	p := args[0].AsFloat()
	mean := 0.0
	stddev := 1.0

	if p <= 0 || p >= 1 {
		return types.Error("norminv: p must be between 0 and 1 (exclusive)")
	}

	if len(args) >= 2 {
		mean = args[1].AsFloat()
	}
	if len(args) >= 3 {
		stddev = args[2].AsFloat()
		if stddev <= 0 {
			return types.Error("norminv: stddev must be positive")
		}
	}

	// Use rational approximation for inverse error function
	z := inverseNormalCDF(p)
	x := mean + z*stddev

	return types.Number(x)
}

// FnZScore calculates the z-score.
// Args: x, mean, stddev
func FnZScoreCalc(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("zscore requires 3 arguments: x, mean, stddev")
	}

	x := args[0].AsFloat()
	mean := args[1].AsFloat()
	stddev := args[2].AsFloat()

	if stddev == 0 {
		return types.Error("zscore: stddev cannot be zero")
	}

	z := (x - mean) / stddev
	return types.Number(z)
}

// FnZScoreToValue converts a z-score to a value.
// Args: z, mean, stddev
func FnZScoreToValue(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("zvalue requires 3 arguments: z, mean, stddev")
	}

	z := args[0].AsFloat()
	mean := args[1].AsFloat()
	stddev := args[2].AsFloat()

	x := mean + z*stddev
	return types.Number(x)
}

// ════════════════════════════════════════════════════════════════
// EXPONENTIAL DISTRIBUTION
// ════════════════════════════════════════════════════════════════

// FnExponentialPDF calculates the exponential probability density function.
// Args: x, lambda (rate)
func FnExponentialPDF(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("exppdf requires 2 arguments: x, lambda")
	}

	x := args[0].AsFloat()
	lambda := args[1].AsFloat()

	if lambda <= 0 {
		return types.Error("exppdf: lambda must be positive")
	}
	if x < 0 {
		return types.Number(0)
	}

	// f(x) = λ * e^(-λx)
	pdf := lambda * math.Exp(-lambda*x)

	return types.Number(pdf)
}

// FnExponentialCDF calculates the exponential cumulative distribution function.
// Args: x, lambda (rate)
func FnExponentialCDF(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("expcdf requires 2 arguments: x, lambda")
	}

	x := args[0].AsFloat()
	lambda := args[1].AsFloat()

	if lambda <= 0 {
		return types.Error("expcdf: lambda must be positive")
	}
	if x < 0 {
		return types.Number(0)
	}

	// F(x) = 1 - e^(-λx)
	cdf := 1 - math.Exp(-lambda*x)

	return types.Number(cdf)
}

// FnExponentialInv calculates the inverse exponential (quantile function).
// Args: p, lambda
func FnExponentialInv(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("expinv requires 2 arguments: p, lambda")
	}

	p := args[0].AsFloat()
	lambda := args[1].AsFloat()

	if p < 0 || p >= 1 {
		return types.Error("expinv: p must be in [0, 1)")
	}
	if lambda <= 0 {
		return types.Error("expinv: lambda must be positive")
	}

	// x = -ln(1-p) / λ
	x := -math.Log(1-p) / lambda

	return types.Number(x)
}

// ════════════════════════════════════════════════════════════════
// UNIFORM DISTRIBUTION
// ════════════════════════════════════════════════════════════════

// FnUniformPDF calculates the uniform probability density function.
// Args: x, a (min), b (max)
func FnUniformPDF(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("uniformpdf requires 3 arguments: x, a, b")
	}

	x := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	if a >= b {
		return types.Error("uniformpdf: a must be less than b")
	}

	if x < a || x > b {
		return types.Number(0)
	}

	pdf := 1.0 / (b - a)
	return types.Number(pdf)
}

// FnUniformCDF calculates the uniform cumulative distribution function.
// Args: x, a (min), b (max)
func FnUniformCDF(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("uniformcdf requires 3 arguments: x, a, b")
	}

	x := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	if a >= b {
		return types.Error("uniformcdf: a must be less than b")
	}

	if x < a {
		return types.Number(0)
	}
	if x > b {
		return types.Number(1)
	}

	cdf := (x - a) / (b - a)
	return types.Number(cdf)
}

// ════════════════════════════════════════════════════════════════
// CONFIDENCE INTERVALS
// ════════════════════════════════════════════════════════════════

// FnConfidence calculates the margin of error for a confidence interval.
// Args: confidence level (e.g., 0.95), stddev, sample size
func FnConfidence(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("confidence requires 3 arguments: level, stddev, n")
	}

	level := args[0].AsFloat()
	stddev := args[1].AsFloat()
	n := args[2].AsFloat()

	if level <= 0 || level >= 1 {
		return types.Error("confidence: level must be between 0 and 1")
	}
	if stddev <= 0 {
		return types.Error("confidence: stddev must be positive")
	}
	if n <= 0 {
		return types.Error("confidence: n must be positive")
	}

	// z-score for the confidence level
	alpha := 1 - level
	z := inverseNormalCDF(1 - alpha/2)

	// Margin of error = z * (σ / √n)
	margin := z * (stddev / math.Sqrt(n))

	return types.Number(margin)
}

// FnConfidenceInterval returns both bounds of a confidence interval.
// Args: mean, confidence level, stddev, sample size
func FnConfidenceInterval(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("confint requires 4 arguments: mean, level, stddev, n")
	}

	mean := args[0].AsFloat()
	level := args[1].AsFloat()
	stddev := args[2].AsFloat()
	n := args[3].AsFloat()

	if level <= 0 || level >= 1 {
		return types.Error("confint: level must be between 0 and 1")
	}
	if stddev <= 0 {
		return types.Error("confint: stddev must be positive")
	}
	if n <= 0 {
		return types.Error("confint: n must be positive")
	}

	alpha := 1 - level
	z := inverseNormalCDF(1 - alpha/2)
	margin := z * (stddev / math.Sqrt(n))

	lower := mean - margin
	upper := mean + margin

	return types.StringValue(fmt.Sprintf("[%.4f, %.4f]", lower, upper))
}

// FnMarginOfError calculates the margin of error.
// Args: z-score or confidence level, stddev, sample size
func FnMarginOfError(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("marginoferror requires 3 arguments: z/level, stddev, n")
	}

	zOrLevel := args[0].AsFloat()
	stddev := args[1].AsFloat()
	n := args[2].AsFloat()

	if stddev <= 0 {
		return types.Error("marginoferror: stddev must be positive")
	}
	if n <= 0 {
		return types.Error("marginoferror: n must be positive")
	}

	var z float64
	if zOrLevel > 0 && zOrLevel < 1 {
		// Treat as confidence level
		alpha := 1 - zOrLevel
		z = inverseNormalCDF(1 - alpha/2)
	} else {
		// Treat as z-score
		z = zOrLevel
	}

	margin := z * (stddev / math.Sqrt(n))
	return types.Number(margin)
}

// FnSampleSize calculates required sample size for a given margin of error.
// Args: confidence level, stddev, margin of error
func FnSampleSize(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("samplesize requires 3 arguments: level, stddev, margin")
	}

	level := args[0].AsFloat()
	stddev := args[1].AsFloat()
	margin := args[2].AsFloat()

	if level <= 0 || level >= 1 {
		return types.Error("samplesize: level must be between 0 and 1")
	}
	if stddev <= 0 {
		return types.Error("samplesize: stddev must be positive")
	}
	if margin <= 0 {
		return types.Error("samplesize: margin must be positive")
	}

	alpha := 1 - level
	z := inverseNormalCDF(1 - alpha/2)

	// n = (z * σ / E)²
	n := math.Pow(z*stddev/margin, 2)

	return types.Number(math.Ceil(n))
}

// ════════════════════════════════════════════════════════════════
// SIMULATION & RANDOM
// ════════════════════════════════════════════════════════════════

// FnDice simulates rolling dice.
// Args: [count], [sides] (defaults: 1d6)
func FnDice(args []types.Value) types.Value {
	count := 1
	sides := 6

	if len(args) >= 1 {
		count = int(args[0].AsFloat())
		if count < 1 {
			count = 1
		}
		if count > 100 {
			count = 100
		}
	}
	if len(args) >= 2 {
		sides = int(args[1].AsFloat())
		if sides < 2 {
			sides = 2
		}
		if sides > 1000 {
			sides = 1000
		}
	}

	total := 0
	for i := 0; i < count; i++ {
		total += rand.Intn(sides) + 1
	}

	return types.Number(float64(total))
}

// FnCoin simulates flipping a coin.
// Args: [count] (defaults: 1)
// Returns: number of heads
func FnCoin(args []types.Value) types.Value {
	count := 1

	if len(args) >= 1 {
		count = int(args[0].AsFloat())
		if count < 1 {
			count = 1
		}
		if count > 10000 {
			count = 10000
		}
	}

	heads := 0
	for i := 0; i < count; i++ {
		if rand.Float64() < 0.5 {
			heads++
		}
	}

	return types.Number(float64(heads))
}

// FnBernoulli simulates a Bernoulli trial.
// Args: p (probability of success)
// Returns: 1 (success) or 0 (failure)
func FnBernoulli(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bernoulli requires 1 argument: p")
	}

	p := args[0].AsFloat()
	if p < 0 || p > 1 {
		return types.Error("bernoulli: p must be between 0 and 1")
	}

	if rand.Float64() < p {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnRandNormal generates a random number from normal distribution.
// Args: [mean], [stddev] (defaults: 0, 1)
func FnRandNormal(args []types.Value) types.Value {
	mean := 0.0
	stddev := 1.0

	if len(args) >= 1 {
		mean = args[0].AsFloat()
	}
	if len(args) >= 2 {
		stddev = args[1].AsFloat()
		if stddev <= 0 {
			return types.Error("randnormal: stddev must be positive")
		}
	}

	// Box-Muller transform
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	return types.Number(mean + z*stddev)
}

// FnRandUniform generates a random number from uniform distribution.
// Args: [min], [max] (defaults: 0, 1)
func FnRandUniform(args []types.Value) types.Value {
	min := 0.0
	max := 1.0

	if len(args) >= 1 {
		min = args[0].AsFloat()
	}
	if len(args) >= 2 {
		max = args[1].AsFloat()
	}

	if min >= max {
		return types.Error("randuniform: min must be less than max")
	}

	return types.Number(min + rand.Float64()*(max-min))
}

// FnRandInt generates a random integer in a range.
// Args: min, max (inclusive)
func FnRandInt(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("randint requires 2 arguments: min, max")
	}

	min := int(args[0].AsFloat())
	max := int(args[1].AsFloat())

	if min > max {
		min, max = max, min
	}

	return types.Number(float64(min + rand.Intn(max-min+1)))
}

// FnRandExp generates a random number from exponential distribution.
// Args: lambda (rate)
func FnRandExp(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("randexp requires 1 argument: lambda")
	}

	lambda := args[0].AsFloat()
	if lambda <= 0 {
		return types.Error("randexp: lambda must be positive")
	}

	// Inverse transform: x = -ln(U) / λ
	return types.Number(-math.Log(rand.Float64()) / lambda)
}

// FnRandBinomial generates a random number from binomial distribution.
// Args: n, p
func FnRandBinomial(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("randbinomial requires 2 arguments: n, p")
	}

	n := int(args[0].AsFloat())
	p := args[1].AsFloat()

	if n < 0 {
		return types.Error("randbinomial: n must be non-negative")
	}
	if p < 0 || p > 1 {
		return types.Error("randbinomial: p must be between 0 and 1")
	}

	successes := 0
	for i := 0; i < n; i++ {
		if rand.Float64() < p {
			successes++
		}
	}

	return types.Number(float64(successes))
}

// FnRandPoisson generates a random number from Poisson distribution.
// Args: lambda
func FnRandPoisson(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("randpoisson requires 1 argument: lambda")
	}

	lambda := args[0].AsFloat()
	if lambda < 0 {
		return types.Error("randpoisson: lambda must be non-negative")
	}

	// Knuth algorithm
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}

	return types.Number(float64(k - 1))
}

// ════════════════════════════════════════════════════════════════
// PROBABILITY UTILITIES
// ════════════════════════════════════════════════════════════════

// FnOdds converts probability to odds.
// Args: p (probability)
func FnOdds(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("odds requires 1 argument: p")
	}

	p := args[0].AsFloat()
	if p <= 0 || p >= 1 {
		return types.Error("odds: p must be between 0 and 1 (exclusive)")
	}

	odds := p / (1 - p)
	return types.Number(odds)
}

// FnOddsToProb converts odds to probability.
// Args: odds
func FnOddsToProb(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("oddstoprob requires 1 argument: odds")
	}

	odds := args[0].AsFloat()
	if odds < 0 {
		return types.Error("oddstoprob: odds must be non-negative")
	}

	p := odds / (1 + odds)
	return types.Number(p)
}

// FnExpectedValue calculates expected value.
// Args: values and probabilities alternating (v1, p1, v2, p2, ...)
func FnExpectedValue(args []types.Value) types.Value {
	if len(args) < 2 || len(args)%2 != 0 {
		return types.Error("expected requires pairs of (value, probability)")
	}

	var ev float64
	var totalProb float64

	for i := 0; i < len(args); i += 2 {
		value := args[i].AsFloat()
		prob := args[i+1].AsFloat()

		if prob < 0 || prob > 1 {
			return types.Errorf("expected: probability must be between 0 and 1, got %v", prob)
		}

		ev += value * prob
		totalProb += prob
	}

	// Warn if probabilities don't sum to 1
	if math.Abs(totalProb-1) > 0.001 {
		// Still return result but it may be scaled
	}

	return types.Number(ev)
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// binomialCoeff calculates C(n, k) = n! / (k! * (n-k)!)
func binomialCoeff(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	if k > n-k {
		k = n - k
	}

	result := 1.0
	for i := 0; i < k; i++ {
		result *= float64(n-i) / float64(i+1)
	}
	return result
}

// factorial64 calculates n! as float64
func factorial64(n int) float64 {
	if n <= 1 {
		return 1
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result
}

// inverseNormalCDF calculates the inverse of the standard normal CDF.
// Uses Abramowitz and Stegun approximation.
func inverseNormalCDF(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}

	// Rational approximation for lower region
	if p < 0.5 {
		return -rationalApproxForInverseNormal(math.Sqrt(-2 * math.Log(p)))
	}
	return rationalApproxForInverseNormal(math.Sqrt(-2 * math.Log(1-p)))
}

func rationalApproxForInverseNormal(t float64) float64 {
	// Coefficients for rational approximation
	c := []float64{2.515517, 0.802853, 0.010328}
	d := []float64{1.432788, 0.189269, 0.001308}

	return t - (c[0]+c[1]*t+c[2]*t*t)/(1+d[0]*t+d[1]*t*t+d[2]*t*t*t)
}
