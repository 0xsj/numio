// internal/finance/tvm.go

package finance

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// PRESENT VALUE (PV)
// ════════════════════════════════════════════════════════════════

// PV calculates the present value of a future amount.
// fv: future value
// rate: periodic interest rate (as decimal)
// nper: number of periods
// Returns present value.
func PV(fv, rate float64, nper int) float64 {
	if nper == 0 {
		return fv
	}
	// PV = FV / (1 + r)^n
	return fv / math.Pow(1+rate, float64(nper))
}

// PVAnnuity calculates present value of an ordinary annuity.
// pmt: periodic payment
// rate: periodic interest rate (as decimal)
// nper: number of periods
// Returns present value.
func PVAnnuity(pmt, rate float64, nper int) float64 {
	if rate == 0 {
		return pmt * float64(nper)
	}
	// PVA = PMT * [(1 - (1 + r)^-n) / r]
	return pmt * (1 - math.Pow(1+rate, float64(-nper))) / rate
}

// PVAnnuityDue calculates present value of an annuity due (payments at start).
func PVAnnuityDue(pmt, rate float64, nper int) float64 {
	if rate == 0 {
		return pmt * float64(nper)
	}
	// PVADue = PVA * (1 + r)
	return PVAnnuity(pmt, rate, nper) * (1 + rate)
}

// PVPerpetuity calculates present value of a perpetuity.
// pmt: periodic payment (forever)
// rate: periodic interest rate
func PVPerpetuity(pmt, rate float64) (float64, error) {
	if rate <= 0 {
		return 0, ErrInvalidRate
	}
	// PV = PMT / r
	return pmt / rate, nil
}

// PVGrowingPerpetuity calculates PV of perpetuity with growth.
// pmt: first payment
// rate: discount rate
// growth: growth rate of payments
func PVGrowingPerpetuity(pmt, rate, growth float64) (float64, error) {
	if rate <= growth {
		return 0, ErrInvalidRate
	}
	// PV = PMT / (r - g)
	return pmt / (rate - growth), nil
}

// ════════════════════════════════════════════════════════════════
// FUTURE VALUE (FV)
// ════════════════════════════════════════════════════════════════

// FV calculates the future value of a present amount.
// pv: present value
// rate: periodic interest rate (as decimal)
// nper: number of periods
func FV(pv, rate float64, nper int) float64 {
	// FV = PV * (1 + r)^n
	return pv * math.Pow(1+rate, float64(nper))
}

// FVAnnuity calculates future value of an ordinary annuity.
// pmt: periodic payment
// rate: periodic interest rate
// nper: number of periods
func FVAnnuity(pmt, rate float64, nper int) float64 {
	if rate == 0 {
		return pmt * float64(nper)
	}
	// FVA = PMT * [((1 + r)^n - 1) / r]
	return pmt * (math.Pow(1+rate, float64(nper)) - 1) / rate
}

// FVAnnuityDue calculates future value of an annuity due.
func FVAnnuityDue(pmt, rate float64, nper int) float64 {
	if rate == 0 {
		return pmt * float64(nper)
	}
	// FVADue = FVA * (1 + r)
	return FVAnnuity(pmt, rate, nper) * (1 + rate)
}

// FVMixed calculates future value with initial amount plus regular payments.
// pv: initial present value
// pmt: periodic payment
// rate: periodic interest rate
// nper: number of periods
func FVMixed(pv, pmt, rate float64, nper int) float64 {
	return FV(pv, rate, nper) + FVAnnuity(pmt, rate, nper)
}

// ════════════════════════════════════════════════════════════════
// PAYMENT (PMT)
// ════════════════════════════════════════════════════════════════

// PMT calculates the periodic payment for a loan or annuity.
// pv: present value (loan amount)
// rate: periodic interest rate
// nper: number of periods
// Returns payment amount (positive for loan payments).
func PMT(pv, rate float64, nper int) float64 {
	if nper == 0 {
		return 0
	}
	if rate == 0 {
		return pv / float64(nper)
	}
	// PMT = PV * [r(1+r)^n] / [(1+r)^n - 1]
	factor := math.Pow(1+rate, float64(nper))
	return pv * (rate * factor) / (factor - 1)
}

// PMTFromFV calculates payment needed to reach a future value.
// fv: target future value
// rate: periodic interest rate
// nper: number of periods
func PMTFromFV(fv, rate float64, nper int) float64 {
	if nper == 0 {
		return 0
	}
	if rate == 0 {
		return fv / float64(nper)
	}
	// PMT = FV * r / [(1+r)^n - 1]
	return fv * rate / (math.Pow(1+rate, float64(nper)) - 1)
}

// PMTDue calculates payment for annuity due (payments at beginning).
func PMTDue(pv, rate float64, nper int) float64 {
	if rate == 0 {
		return pv / float64(nper)
	}
	return PMT(pv, rate, nper) / (1 + rate)
}

// ════════════════════════════════════════════════════════════════
// NUMBER OF PERIODS (NPER)
// ════════════════════════════════════════════════════════════════

// NPER calculates the number of periods for a loan/investment.
// pmt: periodic payment
// pv: present value
// fv: future value (usually 0 for loans)
// rate: periodic interest rate
func NPER(pmt, pv, fv, rate float64) (float64, error) {
	if rate == 0 {
		if pmt == 0 {
			return 0, ErrDivisionByZero
		}
		return -(pv + fv) / pmt, nil
	}

	// n = ln[(PMT - FV*r) / (PMT + PV*r)] / ln(1 + r)
	numerator := pmt - fv*rate
	denominator := pmt + pv*rate

	if denominator == 0 || numerator/denominator <= 0 {
		return 0, ErrNoConvergence
	}

	return math.Log(numerator/denominator) / math.Log(1+rate), nil
}

// NPERSimple calculates periods to pay off a loan.
// principal: loan amount
// pmt: payment per period
// rate: periodic interest rate
func NPERSimple(principal, pmt, rate float64) (float64, error) {
	return NPER(pmt, -principal, 0, rate)
}

// ════════════════════════════════════════════════════════════════
// RATE
// ════════════════════════════════════════════════════════════════

// Rate calculates the interest rate per period.
// Uses Newton-Raphson iteration.
// nper: number of periods
// pmt: payment per period
// pv: present value
// fv: future value (default 0)
func Rate(nper int, pmt, pv, fv float64) (float64, error) {
	if nper <= 0 {
		return 0, ErrInvalidPeriods
	}

	// Initial guess
	rate := 0.1

	// Newton-Raphson iteration
	maxIterations := 100
	tolerance := 1e-10

	for i := 0; i < maxIterations; i++ {
		// f(r) = pv*(1+r)^n + pmt*((1+r)^n - 1)/r + fv = 0
		factor := math.Pow(1+rate, float64(nper))

		var f, df float64
		if rate == 0 {
			f = pv + pmt*float64(nper) + fv
			df = pmt * float64(nper)
		} else {
			f = pv*factor + pmt*(factor-1)/rate + fv
			// Derivative
			dfactor := float64(nper) * math.Pow(1+rate, float64(nper-1))
			df = pv*dfactor + pmt*(dfactor*rate-(factor-1))/(rate*rate)
		}

		if math.Abs(df) < tolerance {
			return 0, ErrNoConvergence
		}

		newRate := rate - f/df

		if math.Abs(newRate-rate) < tolerance {
			return newRate, nil
		}

		rate = newRate
	}

	return 0, ErrNoConvergence
}

// ════════════════════════════════════════════════════════════════
// NET PRESENT VALUE (NPV)
// ════════════════════════════════════════════════════════════════

// NPV calculates the net present value of cash flows.
// rate: discount rate per period
// cashFlows: series of cash flows (first is typically negative = investment)
func NPV(rate float64, cashFlows []float64) float64 {
	var npv float64
	for i, cf := range cashFlows {
		npv += cf / math.Pow(1+rate, float64(i))
	}
	return npv
}

// NPVWithInitial calculates NPV with explicit initial investment.
// initialInvestment: initial cost (positive number, will be negated)
// rate: discount rate
// cashFlows: future cash flows (all positive typically)
func NPVWithInitial(initialInvestment, rate float64, cashFlows []float64) float64 {
	allFlows := make([]float64, len(cashFlows)+1)
	allFlows[0] = -initialInvestment
	copy(allFlows[1:], cashFlows)
	return NPV(rate, allFlows)
}

// ════════════════════════════════════════════════════════════════
// INTERNAL RATE OF RETURN (IRR)
// ════════════════════════════════════════════════════════════════

// IRR calculates the internal rate of return.
// cashFlows: series of cash flows (must have at least one sign change)
// Uses Newton-Raphson iteration.
func IRR(cashFlows []float64) (float64, error) {
	if len(cashFlows) < 2 {
		return 0, ErrInvalidCashFlow
	}

	if !HasSignChange(cashFlows) {
		return 0, ErrInvalidCashFlow
	}

	// Initial guess based on simple return
	totalIn := 0.0
	totalOut := 0.0
	for _, cf := range cashFlows {
		if cf > 0 {
			totalIn += cf
		} else {
			totalOut -= cf
		}
	}

	rate := 0.1
	if totalOut > 0 {
		rate = (totalIn - totalOut) / totalOut / float64(len(cashFlows))
	}

	// Newton-Raphson iteration
	maxIterations := 100
	tolerance := 1e-10

	for i := 0; i < maxIterations; i++ {
		npv := 0.0
		dnpv := 0.0

		for j, cf := range cashFlows {
			factor := math.Pow(1+rate, float64(j))
			npv += cf / factor
			if j > 0 {
				dnpv -= float64(j) * cf / (factor * (1 + rate))
			}
		}

		if math.Abs(dnpv) < tolerance {
			return 0, ErrNoConvergence
		}

		newRate := rate - npv/dnpv

		if math.Abs(newRate-rate) < tolerance {
			return newRate, nil
		}

		rate = newRate

		// Prevent divergence
		if rate < -0.99 {
			rate = -0.99
		}
		if rate > 10 {
			rate = 10
		}
	}

	return 0, ErrNoConvergence
}

// MIRR calculates the modified internal rate of return.
// cashFlows: series of cash flows
// financeRate: rate for negative cash flows (borrowing cost)
// reinvestRate: rate for positive cash flows (reinvestment return)
func MIRR(cashFlows []float64, financeRate, reinvestRate float64) (float64, error) {
	if len(cashFlows) < 2 {
		return 0, ErrInvalidCashFlow
	}

	n := len(cashFlows)

	// PV of negative cash flows at finance rate
	pvNegative := 0.0
	// FV of positive cash flows at reinvest rate
	fvPositive := 0.0

	for i, cf := range cashFlows {
		if cf < 0 {
			pvNegative += cf / math.Pow(1+financeRate, float64(i))
		} else {
			fvPositive += cf * math.Pow(1+reinvestRate, float64(n-1-i))
		}
	}

	if pvNegative >= 0 {
		return 0, ErrInvalidCashFlow
	}

	// MIRR = (FV positive / |PV negative|)^(1/(n-1)) - 1
	return math.Pow(-fvPositive/pvNegative, 1/float64(n-1)) - 1, nil
}

// ════════════════════════════════════════════════════════════════
// XNPV / XIRR (Variable timing)
// ════════════════════════════════════════════════════════════════

// XNPV calculates NPV with irregular cash flow dates.
// rate: annual discount rate
// cashFlows: cash flow amounts
// days: days from start for each cash flow
func XNPV(rate float64, cashFlows []float64, days []int) (float64, error) {
	if len(cashFlows) != len(days) {
		return 0, ErrInvalidCashFlow
	}

	var npv float64
	for i, cf := range cashFlows {
		years := float64(days[i]) / 365.0
		npv += cf / math.Pow(1+rate, years)
	}

	return npv, nil
}
