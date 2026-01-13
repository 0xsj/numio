// internal/finance/investment.go

package finance

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// RETURN ON INVESTMENT (ROI)
// ════════════════════════════════════════════════════════════════

// ROI calculates return on investment as a percentage.
// gain: total gain or final value
// cost: initial investment
// Returns ROI as percentage (e.g., 25 for 25%).
func ROI(gain, cost float64) float64 {
	if cost == 0 {
		return 0
	}
	return ((gain - cost) / cost) * 100
}

// ROIFromValues calculates ROI from initial and final values.
func ROIFromValues(initialValue, finalValue float64) float64 {
	return ROI(finalValue, initialValue)
}

// ROIAnnualized calculates annualized ROI.
// totalReturn: total return as decimal (e.g., 0.50 for 50%)
// years: holding period in years
func ROIAnnualized(totalReturn, years float64) float64 {
	if years <= 0 {
		return 0
	}
	// (1 + total return)^(1/years) - 1
	return (math.Pow(1+totalReturn, 1/years) - 1) * 100
}

// ════════════════════════════════════════════════════════════════
// COMPOUND ANNUAL GROWTH RATE (CAGR)
// ════════════════════════════════════════════════════════════════

// CAGR calculates compound annual growth rate.
// beginValue: starting value
// endValue: ending value
// years: number of years
// Returns CAGR as decimal (e.g., 0.10 for 10%).
func CAGR(beginValue, endValue, years float64) float64 {
	if beginValue <= 0 || years <= 0 {
		return 0
	}
	// CAGR = (EndValue / BeginValue)^(1/years) - 1
	return math.Pow(endValue/beginValue, 1/years) - 1
}

// CAGRPercent returns CAGR as percentage.
func CAGRPercent(beginValue, endValue, years float64) float64 {
	return CAGR(beginValue, endValue, years) * 100
}

// CAGRToFinalValue calculates final value given CAGR.
func CAGRToFinalValue(beginValue, cagr, years float64) float64 {
	return beginValue * math.Pow(1+cagr, years)
}

// ════════════════════════════════════════════════════════════════
// PAYBACK PERIOD
// ════════════════════════════════════════════════════════════════

// PaybackPeriod calculates simple payback period.
// initialInvestment: upfront cost
// annualCashFlow: expected annual cash flow
// Returns years to recover investment.
func PaybackPeriod(initialInvestment, annualCashFlow float64) float64 {
	if annualCashFlow <= 0 {
		return math.Inf(1)
	}
	return initialInvestment / annualCashFlow
}

// PaybackPeriodUneven calculates payback for uneven cash flows.
// initialInvestment: upfront cost
// cashFlows: yearly cash flows
// Returns years to recover (fractional).
func PaybackPeriodUneven(initialInvestment float64, cashFlows []float64) float64 {
	if len(cashFlows) == 0 {
		return math.Inf(1)
	}

	remaining := initialInvestment
	for i, cf := range cashFlows {
		if cf >= remaining {
			// Payback occurs this year
			fraction := remaining / cf
			return float64(i) + fraction
		}
		remaining -= cf
	}

	// Not paid back within cash flow period
	return math.Inf(1)
}

// DiscountedPayback calculates payback period with discounting.
func DiscountedPayback(initialInvestment float64, cashFlows []float64, discountRate float64) float64 {
	if len(cashFlows) == 0 {
		return math.Inf(1)
	}

	remaining := initialInvestment
	for i, cf := range cashFlows {
		// Discount the cash flow
		discountedCF := cf / math.Pow(1+discountRate, float64(i+1))

		if discountedCF >= remaining {
			fraction := remaining / discountedCF
			return float64(i) + fraction
		}
		remaining -= discountedCF
	}

	return math.Inf(1)
}

// ════════════════════════════════════════════════════════════════
// PROFIT MARGINS
// ════════════════════════════════════════════════════════════════

// GrossMargin calculates gross profit margin.
// revenue: total revenue
// cogs: cost of goods sold
// Returns margin as percentage.
func GrossMargin(revenue, cogs float64) float64 {
	if revenue == 0 {
		return 0
	}
	return ((revenue - cogs) / revenue) * 100
}

// NetMargin calculates net profit margin.
// revenue: total revenue
// netIncome: net income after all expenses
func NetMargin(revenue, netIncome float64) float64 {
	if revenue == 0 {
		return 0
	}
	return (netIncome / revenue) * 100
}

// OperatingMargin calculates operating profit margin.
func OperatingMargin(revenue, operatingIncome float64) float64 {
	if revenue == 0 {
		return 0
	}
	return (operatingIncome / revenue) * 100
}

// ════════════════════════════════════════════════════════════════
// MARKUP & MARGIN CONVERSION
// ════════════════════════════════════════════════════════════════

// Markup calculates markup percentage.
// cost: product cost
// price: selling price
// Returns markup as percentage.
func Markup(cost, price float64) float64 {
	if cost == 0 {
		return 0
	}
	return ((price - cost) / cost) * 100
}

// MarkupToPrice calculates selling price from cost and markup.
// cost: product cost
// markupPercent: desired markup percentage
func MarkupToPrice(cost, markupPercent float64) float64 {
	return cost * (1 + markupPercent/100)
}

// MarginToMarkup converts profit margin to markup.
// marginPercent: profit margin as percentage
func MarginToMarkup(marginPercent float64) float64 {
	if marginPercent >= 100 {
		return math.Inf(1)
	}
	return (marginPercent / (100 - marginPercent)) * 100
}

// MarkupToMargin converts markup to profit margin.
// markupPercent: markup as percentage
func MarkupToMargin(markupPercent float64) float64 {
	return (markupPercent / (100 + markupPercent)) * 100
}

// ════════════════════════════════════════════════════════════════
// BREAK-EVEN ANALYSIS
// ════════════════════════════════════════════════════════════════

// BreakEvenUnits calculates units needed to break even.
// fixedCosts: total fixed costs
// pricePerUnit: selling price per unit
// variableCostPerUnit: variable cost per unit
func BreakEvenUnits(fixedCosts, pricePerUnit, variableCostPerUnit float64) float64 {
	contribution := pricePerUnit - variableCostPerUnit
	if contribution <= 0 {
		return math.Inf(1)
	}
	return fixedCosts / contribution
}

// BreakEvenRevenue calculates revenue needed to break even.
func BreakEvenRevenue(fixedCosts, pricePerUnit, variableCostPerUnit float64) float64 {
	units := BreakEvenUnits(fixedCosts, pricePerUnit, variableCostPerUnit)
	return units * pricePerUnit
}

// ContributionMargin calculates contribution margin per unit.
func ContributionMargin(pricePerUnit, variableCostPerUnit float64) float64 {
	return pricePerUnit - variableCostPerUnit
}

// ContributionMarginRatio calculates contribution margin ratio.
func ContributionMarginRatio(pricePerUnit, variableCostPerUnit float64) float64 {
	if pricePerUnit == 0 {
		return 0
	}
	return (pricePerUnit - variableCostPerUnit) / pricePerUnit
}

// ════════════════════════════════════════════════════════════════
// PROFITABILITY INDEX
// ════════════════════════════════════════════════════════════════

// ProfitabilityIndex calculates PI (benefit-cost ratio).
// pvFutureCashFlows: present value of future cash flows
// initialInvestment: initial investment
func ProfitabilityIndex(pvFutureCashFlows, initialInvestment float64) float64 {
	if initialInvestment == 0 {
		return 0
	}
	return pvFutureCashFlows / initialInvestment
}

// ProfitabilityIndexFromNPV calculates PI from NPV.
func ProfitabilityIndexFromNPV(npv, initialInvestment float64) float64 {
	if initialInvestment == 0 {
		return 0
	}
	return (npv + initialInvestment) / initialInvestment
}

// ════════════════════════════════════════════════════════════════
// DIVIDEND CALCULATIONS
// ════════════════════════════════════════════════════════════════

// DividendYield calculates annual dividend yield.
// annualDividend: total annual dividend per share
// pricePerShare: current share price
func DividendYield(annualDividend, pricePerShare float64) float64 {
	if pricePerShare == 0 {
		return 0
	}
	return (annualDividend / pricePerShare) * 100
}

// DividendPayoutRatio calculates payout ratio.
// dividendPerShare: dividend per share
// earningsPerShare: earnings per share
func DividendPayoutRatio(dividendPerShare, earningsPerShare float64) float64 {
	if earningsPerShare == 0 {
		return 0
	}
	return (dividendPerShare / earningsPerShare) * 100
}

// DividendGrowthValue calculates stock value using dividend growth model.
// dividend: next year's expected dividend
// requiredReturn: required rate of return (decimal)
// growthRate: dividend growth rate (decimal)
func DividendGrowthValue(dividend, requiredReturn, growthRate float64) float64 {
	if requiredReturn <= growthRate {
		return math.Inf(1)
	}
	// Gordon Growth Model: P = D1 / (r - g)
	return dividend / (requiredReturn - growthRate)
}

// ════════════════════════════════════════════════════════════════
// EARNINGS & VALUATION
// ════════════════════════════════════════════════════════════════

// EarningsPerShare calculates EPS.
func EarningsPerShare(netIncome, sharesOutstanding float64) float64 {
	if sharesOutstanding == 0 {
		return 0
	}
	return netIncome / sharesOutstanding
}

// PriceToEarnings calculates P/E ratio.
func PriceToEarnings(pricePerShare, earningsPerShare float64) float64 {
	if earningsPerShare == 0 {
		return 0
	}
	return pricePerShare / earningsPerShare
}

// PriceToBook calculates P/B ratio.
func PriceToBook(pricePerShare, bookValuePerShare float64) float64 {
	if bookValuePerShare == 0 {
		return 0
	}
	return pricePerShare / bookValuePerShare
}

// PEGRatio calculates PEG ratio.
// peRatio: price to earnings ratio
// earningsGrowthRate: expected earnings growth rate (percentage)
func PEGRatio(peRatio, earningsGrowthRate float64) float64 {
	if earningsGrowthRate == 0 {
		return 0
	}
	return peRatio / earningsGrowthRate
}
