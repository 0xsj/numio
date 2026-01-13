// internal/finance/finance.go

package finance

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// CONSTANTS
// ════════════════════════════════════════════════════════════════

const (
	// PeriodsPerYear for different compounding frequencies
	Annual     = 1
	SemiAnnual = 2
	Quarterly  = 4
	Monthly    = 12
	Weekly     = 52
	Daily      = 365
)

// ════════════════════════════════════════════════════════════════
// ROUNDING
// ════════════════════════════════════════════════════════════════

// RoundCurrency rounds to 2 decimal places (cents).
func RoundCurrency(amount float64) float64 {
	return math.Round(amount*100) / 100
}

// RoundRate rounds to 4 decimal places.
func RoundRate(rate float64) float64 {
	return math.Round(rate*10000) / 10000
}

// RoundPercent rounds to 2 decimal places for display.
func RoundPercent(rate float64) float64 {
	return math.Round(rate*10000) / 100
}

// ════════════════════════════════════════════════════════════════
// RATE CONVERSION
// ════════════════════════════════════════════════════════════════

// ToDecimal converts a percentage to decimal (5% -> 0.05).
func ToDecimal(percent float64) float64 {
	return percent / 100
}

// ToPercent converts decimal to percentage (0.05 -> 5%).
func ToPercent(decimal float64) float64 {
	return decimal * 100
}

// PeriodicRate converts annual rate to periodic rate.
func PeriodicRate(annualRate float64, periodsPerYear int) float64 {
	return annualRate / float64(periodsPerYear)
}

// AnnualRate converts periodic rate to annual rate.
func AnnualRate(periodicRate float64, periodsPerYear int) float64 {
	return periodicRate * float64(periodsPerYear)
}

// ════════════════════════════════════════════════════════════════
// VALIDATION
// ════════════════════════════════════════════════════════════════

// ValidateRate checks if rate is reasonable (0-100% typically).
func ValidateRate(rate float64) bool {
	return rate >= -1 && rate <= 10 // -100% to 1000%
}

// ValidatePeriods checks if periods is positive.
func ValidatePeriods(n int) bool {
	return n > 0
}

// ValidatePositive checks if value is positive.
func ValidatePositive(v float64) bool {
	return v > 0
}

// ValidateNonNegative checks if value is non-negative.
func ValidateNonNegative(v float64) bool {
	return v >= 0
}

// ════════════════════════════════════════════════════════════════
// ERRORS
// ════════════════════════════════════════════════════════════════

// FinanceError represents a finance calculation error.
type FinanceError string

func (e FinanceError) Error() string {
	return string(e)
}

const (
	ErrInvalidRate      FinanceError = "invalid rate"
	ErrInvalidPeriods   FinanceError = "periods must be positive"
	ErrInvalidPrincipal FinanceError = "principal must be positive"
	ErrInvalidPayment   FinanceError = "payment must be positive"
	ErrDivisionByZero   FinanceError = "division by zero"
	ErrNoConvergence    FinanceError = "calculation did not converge"
	ErrInvalidCashFlow  FinanceError = "invalid cash flow"
)

// ════════════════════════════════════════════════════════════════
// CASH FLOW HELPERS
// ════════════════════════════════════════════════════════════════

// CashFlow represents a series of cash flows.
type CashFlow struct {
	Amount float64
	Period int
}

// SumCashFlows returns the total of all cash flows.
func SumCashFlows(flows []float64) float64 {
	var sum float64
	for _, f := range flows {
		sum += f
	}
	return sum
}

// HasSignChange checks if cash flows change sign (required for IRR).
func HasSignChange(flows []float64) bool {
	if len(flows) < 2 {
		return false
	}

	positive := false
	negative := false

	for _, f := range flows {
		if f > 0 {
			positive = true
		} else if f < 0 {
			negative = true
		}
		if positive && negative {
			return true
		}
	}

	return false
}

// ════════════════════════════════════════════════════════════════
// MATH HELPERS
// ════════════════════════════════════════════════════════════════

// Pow is a helper for exponentiation.
func Pow(base float64, exp int) float64 {
	return math.Pow(base, float64(exp))
}

// Abs returns the absolute value.
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Sign returns -1, 0, or 1.
func Sign(x float64) float64 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}
