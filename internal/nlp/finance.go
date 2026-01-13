// internal/nlp/finance.go

package nlp

import (
	"strings"
)

// ════════════════════════════════════════════════════════════════
// FINANCE PATTERNS
// ════════════════════════════════════════════════════════════════

// RegisterFinancePatterns registers all finance patterns.
func RegisterFinancePatterns(r *PatternRegistry) {
	r.RegisterAll([]*Pattern{
		// Loan patterns
		patternLoanPayment(),
		patternMortgagePayment(),
		patternMonthlyPayment(),

		// Interest/compound patterns
		patternCompoundInterest(),
		patternSimpleInterest(),

		// Investment patterns
		patternFutureValue(),
		patternPresentValue(),
		patternROI(),
		patternCAGR(),

		// Percentage patterns
		patternPercentOf(),
		patternPercentChange(),
		patternWhatPercentOf(),
		patternPercentOff(),

		// Tip patterns
		patternTipCalculation(),
		patternTipSplit(),
	})
}

// ════════════════════════════════════════════════════════════════
// LOAN PATTERNS
// ════════════════════════════════════════════════════════════════

// patternLoanPayment: "loan payment on $250k at 6.5% for 30 years"
func patternLoanPayment() *Pattern {
	return NewPattern("loan_payment").
		Regex(`(?:loan|car|auto)\s+payment\s+(?:on|for)\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|months?|yrs?|mos?)`).
		Keywords("payment", "at", "for").
		Priority(100).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 6 {
				return "", false
			}

			principal, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			term, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			// Convert months to years if needed
			unit := strings.ToLower(matches[5])
			if strings.HasPrefix(unit, "mo") {
				term = term / 12
			}

			return `loan(` + FormatFloat(principal) + `, ` + FormatFloat(rate) + `, ` + FormatFloat(term) + `)`, true
		}).
		Build()
}

// patternMortgagePayment: "mortgage payment on $500k at 7% for 30 years"
func patternMortgagePayment() *Pattern {
	return NewPattern("mortgage_payment").
		Regex(`mortgage\s+(?:payment\s+)?(?:on|for)\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|yrs?)`).
		Keywords("mortgage", "at", "for").
		Priority(100).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 6 {
				return "", false
			}

			principal, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			term, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `loan(` + FormatFloat(principal) + `, ` + FormatFloat(rate) + `, ` + FormatFloat(term) + `)`, true
		}).
		Build()
}

// patternMonthlyPayment: "monthly payment on $25000 at 5% for 5 years"
func patternMonthlyPayment() *Pattern {
	return NewPattern("monthly_payment").
		Regex(`monthly\s+payment\s+(?:on|for)\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|months?|yrs?|mos?)`).
		Keywords("monthly", "payment", "at", "for").
		Priority(95).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 6 {
				return "", false
			}

			principal, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			term, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			// Convert months to years if needed
			unit := strings.ToLower(matches[5])
			if strings.HasPrefix(unit, "mo") {
				term = term / 12
			}

			return `loan(` + FormatFloat(principal) + `, ` + FormatFloat(rate) + `, ` + FormatFloat(term) + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// INTEREST PATTERNS
// ════════════════════════════════════════════════════════════════

// patternCompoundInterest: "$10000 at 5% for 10 years compounded monthly"
func patternCompoundInterest() *Pattern {
	return NewPattern("compound_interest").
		Regex(`([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|yrs?)\s*(?:compounded?\s*(daily|weekly|monthly|quarterly|yearly|annually)?)?`).
		Keywords("at", "for").
		Priority(70).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			// Skip if this looks like a loan pattern
			if strings.Contains(input, "payment") || strings.Contains(input, "loan") || strings.Contains(input, "mortgage") {
				return "", false
			}

			principal, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			years, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			// Default to monthly compounding
			n := 12
			if len(matches) > 6 && matches[6] != "" {
				n = compoundingFrequency(matches[6])
			}

			return `compound(` + FormatFloat(principal) + `, ` + FormatFloat(rate) + `, ` + FormatFloat(years) + `, ` + FormatInt(n) + `)`, true
		}).
		Build()
}

// patternSimpleInterest: "simple interest on $5000 at 3% for 2 years"
func patternSimpleInterest() *Pattern {
	return NewPattern("simple_interest").
		Regex(`simple\s+interest\s+(?:on\s+)?([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|months?|yrs?|mos?)`).
		Keywords("simple", "interest", "at", "for").
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 6 {
				return "", false
			}

			principal, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			time, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			// Convert months to years if needed
			unit := strings.ToLower(matches[5])
			if strings.HasPrefix(unit, "mo") {
				time = time / 12
			}

			// Simple interest: P * r * t
			return `(` + FormatFloat(principal) + ` * ` + FormatFloat(rate/100) + ` * ` + FormatFloat(time) + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// INVESTMENT PATTERNS
// ════════════════════════════════════════════════════════════════

// patternFutureValue: "future value of $10000 at 7% for 20 years"
func patternFutureValue() *Pattern {
	return NewPattern("future_value").
		Regex(`(?:future\s+value|fv)\s+(?:of\s+)?([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|yrs?)`).
		Keywords("at", "for").
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 6 {
				return "", false
			}

			pv, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			years, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `fv(` + FormatFloat(pv) + `, ` + FormatFloat(rate) + `, ` + FormatFloat(years) + `)`, true
		}).
		Build()
}

// patternPresentValue: "present value of $100000 at 5% for 10 years"
func patternPresentValue() *Pattern {
	return NewPattern("present_value").
		Regex(`(?:present\s+value|pv)\s+(?:of\s+)?([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+at\s+(\d+(?:\.\d+)?)\s*%\s+for\s+(\d+)\s*(years?|yrs?)`).
		Keywords("at", "for").
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 6 {
				return "", false
			}

			fv, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			rate, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			years, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `pv(` + FormatFloat(fv) + `, ` + FormatFloat(rate) + `, ` + FormatFloat(years) + `)`, true
		}).
		Build()
}

// patternROI: "roi on $1000 investment now worth $1500"
func patternROI() *Pattern {
	return NewPattern("roi").
		Regex(`roi\s+(?:on\s+)?([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+(?:investment\s+)?(?:now\s+)?worth\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)`).
		Keywords("roi", "worth").
		Priority(75).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			initial, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			final, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `roi(` + FormatFloat(initial) + `, ` + FormatFloat(final) + `)`, true
		}).
		Build()
}

// patternCAGR: "cagr from $10000 to $25000 over 5 years"
func patternCAGR() *Pattern {
	return NewPattern("cagr").
		Regex(`cagr\s+(?:from\s+)?([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+to\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+(?:over\s+)?(\d+)\s*(years?|yrs?)`).
		Keywords("cagr", "to").
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 7 {
				return "", false
			}

			initial, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			final, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			years, ok := ExtractNumber(matches[5])
			if !ok {
				return "", false
			}

			return `cagr(` + FormatFloat(initial) + `, ` + FormatFloat(final) + `, ` + FormatFloat(years) + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// PERCENTAGE PATTERNS
// ════════════════════════════════════════════════════════════════

// patternPercentOf: "20% of 150", "what is 15% of $200"
func patternPercentOf() *Pattern {
	return NewPattern("percent_of").
		Regex(`(?:what\s+is\s+)?(\d+(?:\.\d+)?)\s*%\s+of\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)`).
		Keywords("of").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			percent, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			value, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			symbol := matches[2]
			if symbol != "" {
				return FormatFloat(percent) + `% of ` + symbol + FormatFloat(value), true
			}

			return FormatFloat(percent) + `% of ` + FormatFloat(value), true
		}).
		Build()
}

// patternPercentChange: "percent change from 50 to 75"
func patternPercentChange() *Pattern {
	return NewPattern("percent_change").
		Regex(`(?:percent(?:age)?\s+)?change\s+from\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+to\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)`).
		Keywords("change", "from", "to").
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			from, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			to, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `percentchange(` + FormatFloat(from) + `, ` + FormatFloat(to) + `)`, true
		}).
		Build()
}

// patternWhatPercentOf: "what percent of 200 is 50", "50 is what percent of 200"
func patternWhatPercentOf() *Pattern {
	return NewPattern("what_percent_of").
		Regex(`(?:what\s+percent(?:age)?\s+(?:of|is)\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+is\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)|([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+is\s+what\s+percent(?:age)?\s+of\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?))`).
		Keywords("what", "percent").
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			var whole, part float64
			var ok bool

			if matches[2] != "" {
				// "what percent of 200 is 50"
				whole, ok = ExtractNumber(matches[2])
				if !ok {
					return "", false
				}
				part, ok = ExtractNumber(matches[4])
				if !ok {
					return "", false
				}
			} else {
				// "50 is what percent of 200"
				part, ok = ExtractNumber(matches[6])
				if !ok {
					return "", false
				}
				whole, ok = ExtractNumber(matches[8])
				if !ok {
					return "", false
				}
			}

			return `(` + FormatFloat(part) + ` / ` + FormatFloat(whole) + ` * 100)`, true
		}).
		Build()
}

// patternPercentOff: "20% off $150", "$50 off $200"
func patternPercentOff() *Pattern {
	return NewPattern("percent_off").
		Regex(`(\d+(?:\.\d+)?)\s*%\s+off\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)`).
		Keywords("off").
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			percent, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			price, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			symbol := matches[2]
			discount := price * (percent / 100)
			final := price - discount

			if symbol != "" {
				return symbol + FormatFloat(final), true
			}

			return FormatFloat(final), true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// TIP PATTERNS
// ════════════════════════════════════════════════════════════════

// patternTipCalculation: "15% tip on $85", "tip 20% on $50"
func patternTipCalculation() *Pattern {
	return NewPattern("tip_calculation").
		Regex(`(?:(\d+(?:\.\d+)?)\s*%\s+tip\s+on|tip\s+(\d+(?:\.\d+)?)\s*%\s+on)\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)`).
		Keywords("tip", "on").
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			percentStr := matches[1]
			if percentStr == "" {
				percentStr = matches[2]
			}

			percent, ok := ExtractNumber(percentStr)
			if !ok {
				return "", false
			}

			bill, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `tip(` + FormatFloat(bill) + `, ` + FormatFloat(percent) + `)`, true
		}).
		Build()
}

// patternTipSplit: "split $100 bill with 20% tip 4 ways"
func patternTipSplit() *Pattern {
	return NewPattern("tip_split").
		Regex(`split\s+([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s+(?:bill\s+)?(?:with\s+)?(\d+(?:\.\d+)?)\s*%\s+tip\s+(\d+)\s+ways`).
		Keywords("split", "tip", "ways").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			bill, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			tipPercent, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}

			people, ok := ExtractNumber(matches[4])
			if !ok {
				return "", false
			}

			return `splittip(` + FormatFloat(bill) + `, ` + FormatFloat(tipPercent) + `, ` + FormatInt(int(people)) + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// compoundingFrequency returns the number of compounding periods per year.
func compoundingFrequency(freq string) int {
	freq = strings.ToLower(freq)

	switch freq {
	case "daily":
		return 365
	case "weekly":
		return 52
	case "monthly":
		return 12
	case "quarterly":
		return 4
	case "yearly", "annually", "annual":
		return 1
	default:
		return 12 // Default to monthly
	}
}
