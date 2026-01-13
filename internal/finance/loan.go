// internal/finance/loan.go

package finance

import (
	"math"
)

// ════════════════════════════════════════════════════════════════
// LOAN PAYMENT
// ════════════════════════════════════════════════════════════════

// LoanPayment calculates the monthly payment for a loan.
// principal: loan amount
// annualRate: annual interest rate (as decimal, e.g., 0.06 for 6%)
// years: loan term in years
// Returns monthly payment amount.
func LoanPayment(principal, annualRate float64, years int) float64 {
	if years <= 0 {
		return 0
	}

	monthlyRate := annualRate / 12
	nper := years * 12

	return PMT(principal, monthlyRate, nper)
}

// LoanPaymentMonths calculates payment given months instead of years.
func LoanPaymentMonths(principal, annualRate float64, months int) float64 {
	if months <= 0 {
		return 0
	}

	monthlyRate := annualRate / 12
	return PMT(principal, monthlyRate, months)
}

// LoanPaymentCustom calculates payment with custom compounding.
// principal: loan amount
// annualRate: annual rate (decimal)
// totalPeriods: total number of payment periods
// periodsPerYear: compounding frequency
func LoanPaymentCustom(principal, annualRate float64, totalPeriods, periodsPerYear int) float64 {
	if totalPeriods <= 0 || periodsPerYear <= 0 {
		return 0
	}

	periodicRate := annualRate / float64(periodsPerYear)
	return PMT(principal, periodicRate, totalPeriods)
}

// ════════════════════════════════════════════════════════════════
// LOAN TOTALS
// ════════════════════════════════════════════════════════════════

// LoanTotalPaid returns the total amount paid over the loan term.
func LoanTotalPaid(principal, annualRate float64, years int) float64 {
	payment := LoanPayment(principal, annualRate, years)
	return payment * float64(years*12)
}

// LoanTotalInterest returns the total interest paid over the loan term.
func LoanTotalInterest(principal, annualRate float64, years int) float64 {
	return LoanTotalPaid(principal, annualRate, years) - principal
}

// LoanInterestRatio returns interest as a ratio of principal.
func LoanInterestRatio(principal, annualRate float64, years int) float64 {
	if principal == 0 {
		return 0
	}
	return LoanTotalInterest(principal, annualRate, years) / principal
}

// ════════════════════════════════════════════════════════════════
// LOAN BALANCE
// ════════════════════════════════════════════════════════════════

// LoanBalance calculates remaining balance after n payments.
// principal: original loan amount
// annualRate: annual interest rate (decimal)
// years: total loan term in years
// paymentsMade: number of payments already made
func LoanBalance(principal, annualRate float64, years, paymentsMade int) float64 {
	if paymentsMade <= 0 {
		return principal
	}

	monthlyRate := annualRate / 12
	totalPayments := years * 12

	if paymentsMade >= totalPayments {
		return 0
	}

	payment := PMT(principal, monthlyRate, totalPayments)

	// Balance = P(1+r)^n - PMT*[(1+r)^n - 1]/r
	factor := math.Pow(1+monthlyRate, float64(paymentsMade))
	if monthlyRate == 0 {
		return principal - payment*float64(paymentsMade)
	}
	return principal*factor - payment*(factor-1)/monthlyRate
}

// LoanBalancePercent returns the percentage of loan remaining.
func LoanBalancePercent(principal, annualRate float64, years, paymentsMade int) float64 {
	if principal == 0 {
		return 0
	}
	return LoanBalance(principal, annualRate, years, paymentsMade) / principal * 100
}

// ════════════════════════════════════════════════════════════════
// PAYMENT BREAKDOWN
// ════════════════════════════════════════════════════════════════

// PaymentBreakdown holds the principal and interest portions of a payment.
type PaymentBreakdown struct {
	Payment          float64
	Principal        float64
	Interest         float64
	RemainingBalance float64
}

// LoanPaymentBreakdown calculates breakdown for a specific payment.
// paymentNumber: which payment (1-indexed)
func LoanPaymentBreakdown(principal, annualRate float64, years, paymentNumber int) PaymentBreakdown {
	if paymentNumber <= 0 {
		return PaymentBreakdown{RemainingBalance: principal}
	}

	monthlyRate := annualRate / 12
	totalPayments := years * 12

	if paymentNumber > totalPayments {
		return PaymentBreakdown{}
	}

	payment := PMT(principal, monthlyRate, totalPayments)

	// Get balance before this payment
	balanceBefore := principal
	if paymentNumber > 1 {
		balanceBefore = LoanBalance(principal, annualRate, years, paymentNumber-1)
	}

	// Interest portion = balance * monthly rate
	interestPortion := balanceBefore * monthlyRate

	// Principal portion = payment - interest
	principalPortion := payment - interestPortion

	// Handle final payment rounding
	balanceAfter := balanceBefore - principalPortion
	if paymentNumber == totalPayments {
		balanceAfter = 0
		principalPortion = balanceBefore
		payment = principalPortion + interestPortion
	}

	return PaymentBreakdown{
		Payment:          RoundCurrency(payment),
		Principal:        RoundCurrency(principalPortion),
		Interest:         RoundCurrency(interestPortion),
		RemainingBalance: RoundCurrency(balanceAfter),
	}
}

// ════════════════════════════════════════════════════════════════
// AMORTIZATION
// ════════════════════════════════════════════════════════════════

// AmortizationSchedule generates full amortization schedule.
func AmortizationSchedule(principal, annualRate float64, years int) []PaymentBreakdown {
	totalPayments := years * 12
	schedule := make([]PaymentBreakdown, totalPayments)

	for i := 1; i <= totalPayments; i++ {
		schedule[i-1] = LoanPaymentBreakdown(principal, annualRate, years, i)
	}

	return schedule
}

// AmortizationSummary provides summary statistics.
type AmortizationSummary struct {
	Principal      float64
	MonthlyPayment float64
	TotalPayments  int
	TotalPaid      float64
	TotalInterest  float64
	InterestRatio  float64
}

// LoanSummary calculates loan summary statistics.
func LoanSummary(principal, annualRate float64, years int) AmortizationSummary {
	payment := LoanPayment(principal, annualRate, years)
	totalPayments := years * 12
	totalPaid := payment * float64(totalPayments)
	totalInterest := totalPaid - principal

	return AmortizationSummary{
		Principal:      principal,
		MonthlyPayment: RoundCurrency(payment),
		TotalPayments:  totalPayments,
		TotalPaid:      RoundCurrency(totalPaid),
		TotalInterest:  RoundCurrency(totalInterest),
		InterestRatio:  totalInterest / principal,
	}
}

// ════════════════════════════════════════════════════════════════
// AFFORDABILITY
// ════════════════════════════════════════════════════════════════

// MaxLoanAmount calculates maximum loan for a given payment.
// payment: maximum affordable monthly payment
// annualRate: annual interest rate (decimal)
// years: loan term in years
func MaxLoanAmount(payment, annualRate float64, years int) float64 {
	if years <= 0 {
		return 0
	}

	monthlyRate := annualRate / 12
	nper := years * 12

	if monthlyRate == 0 {
		return payment * float64(nper)
	}

	// PV = PMT * [(1 - (1+r)^-n) / r]
	return payment * (1 - math.Pow(1+monthlyRate, float64(-nper))) / monthlyRate
}

// RequiredIncome estimates required income for a loan (28% rule).
// principal: loan amount
// annualRate: annual interest rate
// years: loan term
// Returns estimated annual income needed.
func RequiredIncome(principal, annualRate float64, years int) float64 {
	monthlyPayment := LoanPayment(principal, annualRate, years)
	// 28% of gross monthly income guideline
	return (monthlyPayment / 0.28) * 12
}

// ════════════════════════════════════════════════════════════════
// EXTRA PAYMENTS
// ════════════════════════════════════════════════════════════════

// LoanWithExtraPayment calculates new payoff time with extra payment.
// principal: loan amount
// annualRate: annual rate
// years: original term
// extraMonthly: additional monthly payment
// Returns new number of months to payoff.
func LoanWithExtraPayment(principal, annualRate float64, years int, extraMonthly float64) int {
	monthlyRate := annualRate / 12
	basePayment := LoanPayment(principal, annualRate, years)
	totalPayment := basePayment + extraMonthly

	if totalPayment <= principal*monthlyRate {
		// Payment doesn't cover interest
		return -1
	}

	// Calculate months to payoff
	// n = -ln(1 - P*r/PMT) / ln(1 + r)
	if monthlyRate == 0 {
		return int(math.Ceil(principal / totalPayment))
	}

	ratio := 1 - (principal*monthlyRate)/totalPayment
	if ratio <= 0 {
		return -1
	}

	months := -math.Log(ratio) / math.Log(1+monthlyRate)
	return int(math.Ceil(months))
}

// ExtraPaymentSavings calculates interest saved with extra payments.
func ExtraPaymentSavings(principal, annualRate float64, years int, extraMonthly float64) float64 {
	originalInterest := LoanTotalInterest(principal, annualRate, years)

	newMonths := LoanWithExtraPayment(principal, annualRate, years, extraMonthly)
	if newMonths <= 0 {
		return 0
	}

	basePayment := LoanPayment(principal, annualRate, years)
	totalPayment := basePayment + extraMonthly
	newTotalPaid := totalPayment * float64(newMonths)
	newInterest := newTotalPaid - principal

	return originalInterest - newInterest
}

// ════════════════════════════════════════════════════════════════
// REFINANCE
// ════════════════════════════════════════════════════════════════

// RefinanceAnalysis holds refinance comparison data.
type RefinanceAnalysis struct {
	CurrentPayment   float64
	NewPayment       float64
	MonthlySavings   float64
	NewTotalInterest float64
	InterestSavings  float64
	BreakEvenMonths  int
	ClosingCosts     float64
}

// AnalyzeRefinance compares current loan to refinance option.
func AnalyzeRefinance(currentBalance, currentRate float64, currentRemainingYears int,
	newRate float64, newYears int, closingCosts float64) RefinanceAnalysis {

	currentPayment := LoanPayment(currentBalance, currentRate, currentRemainingYears)
	currentTotalInterest := LoanTotalInterest(currentBalance, currentRate, currentRemainingYears)

	// New loan includes closing costs rolled in (or paid upfront)
	newPayment := LoanPayment(currentBalance, newRate, newYears)
	newTotalInterest := LoanTotalInterest(currentBalance, newRate, newYears) + closingCosts

	monthlySavings := currentPayment - newPayment
	interestSavings := currentTotalInterest - newTotalInterest

	breakEvenMonths := 0
	if monthlySavings > 0 {
		breakEvenMonths = int(math.Ceil(closingCosts / monthlySavings))
	}

	return RefinanceAnalysis{
		CurrentPayment:   RoundCurrency(currentPayment),
		NewPayment:       RoundCurrency(newPayment),
		MonthlySavings:   RoundCurrency(monthlySavings),
		NewTotalInterest: RoundCurrency(newTotalInterest),
		InterestSavings:  RoundCurrency(interestSavings),
		BreakEvenMonths:  breakEvenMonths,
		ClosingCosts:     closingCosts,
	}
}
