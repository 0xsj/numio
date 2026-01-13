// internal/finance/interest.go

package finance

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// SIMPLE INTEREST
// ════════════════════════════════════════════════════════════════

// SimpleInterest calculates simple interest.
// principal: initial amount
// rate: annual interest rate (as decimal, e.g., 0.05 for 5%)
// time: time in years
// Returns the interest earned.
func SimpleInterest(principal, rate, time float64) float64 {
	return principal * rate * time
}

// SimpleInterestTotal returns principal + interest.
func SimpleInterestTotal(principal, rate, time float64) float64 {
	return principal + SimpleInterest(principal, rate, time)
}

// SimpleInterestRate calculates the rate given other values.
// Returns rate as decimal.
func SimpleInterestRate(principal, interest, time float64) (float64, error) {
	if principal == 0 || time == 0 {
		return 0, ErrDivisionByZero
	}
	return interest / (principal * time), nil
}

// SimpleInterestTime calculates time given other values.
func SimpleInterestTime(principal, interest, rate float64) (float64, error) {
	if principal == 0 || rate == 0 {
		return 0, ErrDivisionByZero
	}
	return interest / (principal * rate), nil
}

// ════════════════════════════════════════════════════════════════
// COMPOUND INTEREST
// ════════════════════════════════════════════════════════════════

// CompoundInterest calculates compound interest.
// principal: initial amount
// rate: annual interest rate (as decimal)
// n: compounding periods per year
// time: time in years
// Returns the final amount (principal + interest).
func CompoundInterest(principal, rate float64, n int, time float64) float64 {
	if n <= 0 {
		n = 1
	}
	// A = P(1 + r/n)^(nt)
	return principal * math.Pow(1+rate/float64(n), float64(n)*time)
}

// CompoundInterestOnly returns just the interest earned.
func CompoundInterestOnly(principal, rate float64, n int, time float64) float64 {
	return CompoundInterest(principal, rate, n, time) - principal
}

// ContinuousCompound calculates continuously compounded interest.
// A = Pe^(rt)
func ContinuousCompound(principal, rate, time float64) float64 {
	return principal * math.Exp(rate*time)
}

// ContinuousCompoundOnly returns just the interest from continuous compounding.
func ContinuousCompoundOnly(principal, rate, time float64) float64 {
	return ContinuousCompound(principal, rate, time) - principal
}

// ════════════════════════════════════════════════════════════════
// APR & APY (Annual Percentage Rate / Yield)
// ════════════════════════════════════════════════════════════════

// APRtoAPY converts Annual Percentage Rate to Annual Percentage Yield.
// apr: annual percentage rate (as decimal)
// n: compounding periods per year
// Returns APY as decimal.
func APRtoAPY(apr float64, n int) float64 {
	if n <= 0 {
		n = 1
	}
	// APY = (1 + APR/n)^n - 1
	return math.Pow(1+apr/float64(n), float64(n)) - 1
}

// APYtoAPR converts Annual Percentage Yield to Annual Percentage Rate.
// apy: annual percentage yield (as decimal)
// n: compounding periods per year
// Returns APR as decimal.
func APYtoAPR(apy float64, n int) float64 {
	if n <= 0 {
		n = 1
	}
	// APR = n * ((1 + APY)^(1/n) - 1)
	return float64(n) * (math.Pow(1+apy, 1/float64(n)) - 1)
}

// EffectiveAnnualRate calculates the effective annual rate.
// Same as APRtoAPY but more explicit naming.
func EffectiveAnnualRate(nominalRate float64, periodsPerYear int) float64 {
	return APRtoAPY(nominalRate, periodsPerYear)
}

// NominalRate calculates the nominal rate from effective rate.
// Same as APYtoAPR.
func NominalRate(effectiveRate float64, periodsPerYear int) float64 {
	return APYtoAPR(effectiveRate, periodsPerYear)
}

// ════════════════════════════════════════════════════════════════
// RATE DOUBLING
// ════════════════════════════════════════════════════════════════

// RuleOf72 estimates years to double at a given rate.
// rate: annual rate as percentage (e.g., 6 for 6%)
func RuleOf72(ratePercent float64) float64 {
	if ratePercent == 0 {
		return math.Inf(1)
	}
	return 72 / ratePercent
}

// ExactDoubleTime calculates exact time to double.
// rate: annual rate as decimal (e.g., 0.06 for 6%)
func ExactDoubleTime(rate float64) float64 {
	if rate <= 0 {
		return math.Inf(1)
	}
	// t = ln(2) / ln(1 + r)
	return math.Log(2) / math.Log(1+rate)
}

// DoubleTimeCompound calculates time to double with n compounding periods.
func DoubleTimeCompound(rate float64, n int) float64 {
	if rate <= 0 || n <= 0 {
		return math.Inf(1)
	}
	// t = ln(2) / (n * ln(1 + r/n))
	return math.Log(2) / (float64(n) * math.Log(1+rate/float64(n)))
}

// ════════════════════════════════════════════════════════════════
// RATE REQUIRED
// ════════════════════════════════════════════════════════════════

// RateToReach calculates rate needed to reach target in given time.
// principal: starting amount
// target: desired final amount
// time: years
// n: compounding periods per year
// Returns annual rate as decimal.
func RateToReach(principal, target, time float64, n int) (float64, error) {
	if principal <= 0 || time <= 0 || n <= 0 {
		return 0, ErrDivisionByZero
	}
	if target <= 0 {
		return 0, ErrInvalidPrincipal
	}

	// A = P(1 + r/n)^(nt)
	// (A/P)^(1/nt) = 1 + r/n
	// r = n * ((A/P)^(1/nt) - 1)
	ratio := target / principal
	exponent := 1 / (float64(n) * time)
	periodicRate := math.Pow(ratio, exponent) - 1
	annualRate := periodicRate * float64(n)

	return annualRate, nil
}

// TimeToReach calculates time needed to reach target at given rate.
func TimeToReach(principal, target, rate float64, n int) (float64, error) {
	if rate <= 0 || n <= 0 {
		return 0, ErrDivisionByZero
	}
	if principal <= 0 || target <= 0 {
		return 0, ErrInvalidPrincipal
	}

	// t = ln(A/P) / (n * ln(1 + r/n))
	return math.Log(target/principal) / (float64(n) * math.Log(1+rate/float64(n))), nil
}

// ════════════════════════════════════════════════════════════════
// INFLATION ADJUSTMENT
// ════════════════════════════════════════════════════════════════

// RealRate calculates real rate adjusted for inflation.
// nominalRate: nominal interest rate (decimal)
// inflationRate: inflation rate (decimal)
// Returns real rate as decimal.
func RealRate(nominalRate, inflationRate float64) float64 {
	// Fisher equation: (1 + r) = (1 + n) / (1 + i)
	// r = (1 + n) / (1 + i) - 1
	if inflationRate == -1 {
		return math.Inf(1)
	}
	return (1+nominalRate)/(1+inflationRate) - 1
}

// NominalFromReal calculates nominal rate from real rate and inflation.
func NominalFromReal(realRate, inflationRate float64) float64 {
	// n = (1 + r)(1 + i) - 1
	return (1+realRate)*(1+inflationRate) - 1
}

// InflationAdjustedValue calculates future purchasing power.
func InflationAdjustedValue(amount, inflationRate, years float64) float64 {
	if inflationRate <= -1 {
		return math.Inf(1)
	}
	return amount / math.Pow(1+inflationRate, years)
}
