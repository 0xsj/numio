// internal/eval/functions_finance.go

package eval

import (
	"github.com/0xsj/numio/internal/finance"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// INTEREST FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSimpleInterest calculates simple interest.
// Args: principal, rate (%), time (years)
func FnSimpleInterest(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("simple requires 3 arguments: principal, rate%, years")
	}

	principal := args[0].AsFloat()
	rate := args[1].AsFloat() / 100 // Convert percentage to decimal
	time := args[2].AsFloat()

	result := finance.SimpleInterest(principal, rate, time)
	return types.Number(finance.RoundCurrency(result))
}

// FnCompoundInterest calculates compound interest (final amount).
// Args: principal, rate (%), years [, periods per year]
func FnCompoundInterest(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("compound requires 3-4 arguments: principal, rate%, years [, periods]")
	}

	principal := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	years := args[2].AsFloat()
	periods := 12 // Default monthly

	if len(args) == 4 {
		periods = int(args[3].AsFloat())
	}

	result := finance.CompoundInterest(principal, rate, periods, years)
	return types.Number(finance.RoundCurrency(result))
}

// FnAPY converts APR to APY.
// Args: apr (%), periods per year
func FnAPY(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("apy requires 2 arguments: apr%, periods")
	}

	apr := args[0].AsFloat() / 100
	periods := int(args[1].AsFloat())

	result := finance.APRtoAPY(apr, periods) * 100
	return types.Number(finance.RoundRate(result))
}

// FnAPR converts APY to APR.
// Args: apy (%), periods per year
func FnAPR(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("apr requires 2 arguments: apy%, periods")
	}

	apy := args[0].AsFloat() / 100
	periods := int(args[1].AsFloat())

	result := finance.APYtoAPR(apy, periods) * 100
	return types.Number(finance.RoundRate(result))
}

// FnRule72 estimates years to double investment.
// Args: rate (%)
func FnRule72(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("rule72 requires 1 argument: rate%")
	}

	rate := args[0].AsFloat()
	result := finance.RuleOf72(rate)
	return types.Number(finance.RoundRate(result))
}

// FnEffectiveRate calculates effective annual rate.
// Args: nominal rate (%), periods per year
func FnEffectiveRate(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("effectiverate requires 2 arguments: rate%, periods")
	}

	rate := args[0].AsFloat() / 100
	periods := int(args[1].AsFloat())

	result := finance.EffectiveAnnualRate(rate, periods) * 100
	return types.Number(finance.RoundRate(result))
}

// FnRealRate calculates inflation-adjusted rate.
// Args: nominal rate (%), inflation rate (%)
func FnRealRate(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("realrate requires 2 arguments: nominal%, inflation%")
	}

	nominal := args[0].AsFloat() / 100
	inflation := args[1].AsFloat() / 100

	result := finance.RealRate(nominal, inflation) * 100
	return types.Number(finance.RoundRate(result))
}

// ════════════════════════════════════════════════════════════════
// TIME VALUE OF MONEY
// ════════════════════════════════════════════════════════════════

// FnPV calculates present value.
// Args: future value, rate (%), periods
func FnPV(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("pv requires 3 arguments: fv, rate%, periods")
	}

	fv := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	nper := int(args[2].AsFloat())

	result := finance.PV(fv, rate, nper)
	return types.Number(finance.RoundCurrency(result))
}

// FnFV calculates future value.
// Args: present value, rate (%), periods
func FnFV(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("fv requires 3 arguments: pv, rate%, periods")
	}

	pv := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	nper := int(args[2].AsFloat())

	result := finance.FV(pv, rate, nper)
	return types.Number(finance.RoundCurrency(result))
}

// FnPMT calculates periodic payment.
// Args: principal, rate (%), periods
func FnPMT(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("pmt requires 3 arguments: principal, rate%, periods")
	}

	pv := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	nper := int(args[2].AsFloat())

	result := finance.PMT(pv, rate, nper)
	return types.Number(finance.RoundCurrency(result))
}

// FnNPER calculates number of periods.
// Args: payment, principal, rate (%)
func FnNPER(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("nper requires 3 arguments: payment, principal, rate%")
	}

	pmt := args[0].AsFloat()
	principal := args[1].AsFloat()
	rate := args[2].AsFloat() / 100

	result, err := finance.NPERSimple(principal, pmt, rate)
	if err != nil {
		return types.Error("nper: calculation error")
	}
	return types.Number(result)
}

// FnNPV calculates net present value.
// Args: rate (%), cash flows...
func FnNPV(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("npv requires at least 2 arguments: rate% and cash flows")
	}

	rate := args[0].AsFloat() / 100
	cashFlows := make([]float64, len(args)-1)
	for i, arg := range args[1:] {
		cashFlows[i] = arg.AsFloat()
	}

	result := finance.NPV(rate, cashFlows)
	return types.Number(finance.RoundCurrency(result))
}

// FnIRR calculates internal rate of return.
// Args: cash flows...
func FnIRR(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("irr requires at least 2 cash flows")
	}

	cashFlows := valuesToFloats(args)

	result, err := finance.IRR(cashFlows)
	if err != nil {
		return types.Error("irr: " + err.Error())
	}
	return types.Number(finance.RoundRate(result * 100))
}

// FnMIRR calculates modified internal rate of return.
// Args: finance rate (%), reinvest rate (%), cash flows...
func FnMIRR(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("mirr requires at least 4 arguments: financeRate%, reinvestRate%, cash flows...")
	}

	financeRate := args[0].AsFloat() / 100
	reinvestRate := args[1].AsFloat() / 100
	cashFlows := valuesToFloats(args[2:])

	result, err := finance.MIRR(cashFlows, financeRate, reinvestRate)
	if err != nil {
		return types.Error("mirr: " + err.Error())
	}
	return types.Number(finance.RoundRate(result * 100))
}

// ════════════════════════════════════════════════════════════════
// LOAN FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnLoan calculates monthly loan payment.
// Args: principal, annual rate (%), years
func FnLoan(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("loan requires 3 arguments: principal, rate%, years")
	}

	principal := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	years := int(args[2].AsFloat())

	result := finance.LoanPayment(principal, rate, years)
	return types.Number(finance.RoundCurrency(result))
}

// FnMortgage is an alias for loan.
func FnMortgage(args []types.Value) types.Value {
	return FnLoan(args)
}

// FnLoanTotal calculates total amount paid over loan term.
// Args: principal, annual rate (%), years
func FnLoanTotal(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("loantotal requires 3 arguments: principal, rate%, years")
	}

	principal := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	years := int(args[2].AsFloat())

	result := finance.LoanTotalPaid(principal, rate, years)
	return types.Number(finance.RoundCurrency(result))
}

// FnLoanInterest calculates total interest paid.
// Args: principal, annual rate (%), years
func FnLoanInterest(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("loaninterest requires 3 arguments: principal, rate%, years")
	}

	principal := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	years := int(args[2].AsFloat())

	result := finance.LoanTotalInterest(principal, rate, years)
	return types.Number(finance.RoundCurrency(result))
}

// FnLoanBalance calculates remaining balance after n payments.
// Args: principal, annual rate (%), years, payments made
func FnLoanBalance(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("loanbalance requires 4 arguments: principal, rate%, years, payments")
	}

	principal := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	years := int(args[2].AsFloat())
	payments := int(args[3].AsFloat())

	result := finance.LoanBalance(principal, rate, years, payments)
	return types.Number(finance.RoundCurrency(result))
}

// FnMaxLoan calculates maximum loan amount for a payment.
// Args: payment, annual rate (%), years
func FnMaxLoan(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("maxloan requires 3 arguments: payment, rate%, years")
	}

	payment := args[0].AsFloat()
	rate := args[1].AsFloat() / 100
	years := int(args[2].AsFloat())

	result := finance.MaxLoanAmount(payment, rate, years)
	return types.Number(finance.RoundCurrency(result))
}

// ════════════════════════════════════════════════════════════════
// INVESTMENT FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnROI calculates return on investment.
// Args: gain (or final value), cost (or initial value)
func FnROI(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("roi requires 2 arguments: gain, cost")
	}

	gain := args[0].AsFloat()
	cost := args[1].AsFloat()

	result := finance.ROI(gain, cost)
	return types.Number(finance.RoundRate(result))
}

// FnCAGR calculates compound annual growth rate.
// Args: begin value, end value, years
func FnCAGR(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("cagr requires 3 arguments: beginValue, endValue, years")
	}

	beginValue := args[0].AsFloat()
	endValue := args[1].AsFloat()
	years := args[2].AsFloat()

	result := finance.CAGRPercent(beginValue, endValue, years)
	return types.Number(finance.RoundRate(result))
}

// FnPayback calculates simple payback period.
// Args: initial investment, annual cash flow
func FnPayback(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("payback requires 2 arguments: investment, annualCashFlow")
	}

	investment := args[0].AsFloat()
	cashFlow := args[1].AsFloat()

	result := finance.PaybackPeriod(investment, cashFlow)
	return types.Number(finance.RoundRate(result))
}

// FnGrossMargin calculates gross profit margin.
// Args: revenue, cost of goods sold
func FnGrossMargin(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("grossmargin requires 2 arguments: revenue, cogs")
	}

	revenue := args[0].AsFloat()
	cogs := args[1].AsFloat()

	result := finance.GrossMargin(revenue, cogs)
	return types.Number(finance.RoundRate(result))
}

// FnNetMargin calculates net profit margin.
// Args: revenue, net income
func FnNetMargin(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("netmargin requires 2 arguments: revenue, netIncome")
	}

	revenue := args[0].AsFloat()
	netIncome := args[1].AsFloat()

	result := finance.NetMargin(revenue, netIncome)
	return types.Number(finance.RoundRate(result))
}

// FnMarkup calculates markup percentage.
// Args: cost, price
func FnMarkup(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("markup requires 2 arguments: cost, price")
	}

	cost := args[0].AsFloat()
	price := args[1].AsFloat()

	result := finance.Markup(cost, price)
	return types.Number(finance.RoundRate(result))
}

// FnBreakeven calculates break-even units.
// Args: fixed costs, price per unit, variable cost per unit
func FnBreakeven(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("breakeven requires 3 arguments: fixedCosts, price, variableCost")
	}

	fixedCosts := args[0].AsFloat()
	price := args[1].AsFloat()
	variableCost := args[2].AsFloat()

	result := finance.BreakEvenUnits(fixedCosts, price, variableCost)
	return types.Number(result)
}

// ════════════════════════════════════════════════════════════════
// DEPRECIATION FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSLN calculates straight-line depreciation.
// Args: cost, salvage, life
func FnSLN(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("sln requires 3 arguments: cost, salvage, life")
	}

	cost := args[0].AsFloat()
	salvage := args[1].AsFloat()
	life := int(args[2].AsFloat())

	result := finance.StraightLine(cost, salvage, life)
	return types.Number(finance.RoundCurrency(result))
}

// FnDDB calculates double declining balance depreciation.
// Args: cost, salvage, life, year
func FnDDB(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("ddb requires 4 arguments: cost, salvage, life, year")
	}

	cost := args[0].AsFloat()
	salvage := args[1].AsFloat()
	life := int(args[2].AsFloat())
	year := int(args[3].AsFloat())

	result := finance.DoubleDeclining(cost, salvage, life, year)
	return types.Number(result)
}

// FnSYD calculates sum-of-years'-digits depreciation.
// Args: cost, salvage, life, year
func FnSYD(args []types.Value) types.Value {
	if len(args) != 4 {
		return types.Error("syd requires 4 arguments: cost, salvage, life, year")
	}

	cost := args[0].AsFloat()
	salvage := args[1].AsFloat()
	life := int(args[2].AsFloat())
	year := int(args[3].AsFloat())

	result := finance.SumOfYearsDigits(cost, salvage, life, year)
	return types.Number(result)
}

// FnMACRS calculates MACRS depreciation.
// Args: cost, property class (3,5,7,10,15), year
func FnMACRS(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("macrs requires 3 arguments: cost, propertyClass, year")
	}

	cost := args[0].AsFloat()
	propertyClass := int(args[1].AsFloat())
	year := int(args[2].AsFloat())

	result := finance.MACRS(cost, propertyClass, year)
	return types.Number(result)
}

// ════════════════════════════════════════════════════════════════
// DIVIDEND & VALUATION
// ════════════════════════════════════════════════════════════════

// FnDividendYield calculates dividend yield.
// Args: annual dividend, price per share
func FnDividendYield(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("divyield requires 2 arguments: dividend, price")
	}

	dividend := args[0].AsFloat()
	price := args[1].AsFloat()

	result := finance.DividendYield(dividend, price)
	return types.Number(finance.RoundRate(result))
}

// FnPE calculates price-to-earnings ratio.
// Args: price per share, earnings per share
func FnPE(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("pe requires 2 arguments: price, eps")
	}

	price := args[0].AsFloat()
	eps := args[1].AsFloat()

	result := finance.PriceToEarnings(price, eps)
	return types.Number(finance.RoundRate(result))
}

// FnEPS calculates earnings per share.
// Args: net income, shares outstanding
func FnEPS(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("eps requires 2 arguments: netIncome, shares")
	}

	netIncome := args[0].AsFloat()
	shares := args[1].AsFloat()

	result := finance.EarningsPerShare(netIncome, shares)
	return types.Number(finance.RoundCurrency(result))
}

// ════════════════════════════════════════════════════════════════
// TIP FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnTip calculates tip amount and total.
// Args: bill, tipPercent
// Returns: total (bill + tip)
func FnTip(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("tip requires 2 arguments: bill, tipPercent")
	}

	bill := args[0].AsFloat()
	tipPercent := args[1].AsFloat() / 100

	tip := bill * tipPercent
	total := bill + tip

	return types.Number(finance.RoundCurrency(total))
}

// FnTipAmount calculates just the tip amount.
// Args: bill, tipPercent
func FnTipAmount(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("tipamount requires 2 arguments: bill, tipPercent")
	}

	bill := args[0].AsFloat()
	tipPercent := args[1].AsFloat() / 100

	tip := bill * tipPercent

	return types.Number(finance.RoundCurrency(tip))
}

// FnSplitTip calculates per-person amount for a split bill with tip.
// Args: bill, tipPercent, people
func FnSplitTip(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("splittip requires 3 arguments: bill, tipPercent, people")
	}

	bill := args[0].AsFloat()
	tipPercent := args[1].AsFloat() / 100
	people := args[2].AsFloat()

	if people <= 0 {
		return types.Error("splittip: people must be greater than 0")
	}

	total := bill * (1 + tipPercent)
	perPerson := total / people

	return types.Number(finance.RoundCurrency(perPerson))
}

// FnPercentChange calculates percentage change between two values.
// Args: oldValue, newValue
// Returns: percentage change (e.g., 50 for 50% increase)
func FnPercentChange(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("percentchange requires 2 arguments: oldValue, newValue")
	}

	oldVal := args[0].AsFloat()
	newVal := args[1].AsFloat()

	if oldVal == 0 {
		return types.Error("percentchange: old value cannot be zero")
	}

	change := ((newVal - oldVal) / oldVal) * 100

	return types.Number(finance.RoundRate(change))
}
