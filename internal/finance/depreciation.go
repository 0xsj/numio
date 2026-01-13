// internal/finance/depreciation.go

package finance

// ════════════════════════════════════════════════════════════════
// STRAIGHT-LINE DEPRECIATION
// ════════════════════════════════════════════════════════════════

// StraightLine calculates annual straight-line depreciation.
// cost: initial cost of asset
// salvage: salvage value at end of life
// life: useful life in years
// Returns annual depreciation amount.
func StraightLine(cost, salvage float64, life int) float64 {
	if life <= 0 {
		return 0
	}
	// SLN = (Cost - Salvage) / Life
	return (cost - salvage) / float64(life)
}

// StraightLineRate returns the annual depreciation rate.
func StraightLineRate(life int) float64 {
	if life <= 0 {
		return 0
	}
	return 1 / float64(life)
}

// StraightLineSchedule returns full depreciation schedule.
func StraightLineSchedule(cost, salvage float64, life int) []DepreciationYear {
	if life <= 0 {
		return nil
	}

	annual := StraightLine(cost, salvage, life)
	schedule := make([]DepreciationYear, life)
	bookValue := cost

	for i := 0; i < life; i++ {
		depreciation := annual
		// Handle final year rounding
		if i == life-1 {
			depreciation = bookValue - salvage
		}

		bookValue -= depreciation
		schedule[i] = DepreciationYear{
			Year:            i + 1,
			Depreciation:    RoundCurrency(depreciation),
			AccumulatedDepr: RoundCurrency(cost - bookValue),
			BookValue:       RoundCurrency(bookValue),
		}
	}

	return schedule
}

// ════════════════════════════════════════════════════════════════
// DECLINING BALANCE DEPRECIATION
// ════════════════════════════════════════════════════════════════

// DecliningBalance calculates depreciation for a specific year.
// cost: initial cost
// salvage: salvage value
// life: useful life in years
// year: which year (1-indexed)
// factor: depreciation factor (2 for double declining)
func DecliningBalance(cost, salvage float64, life, year int, factor float64) float64 {
	if life <= 0 || year <= 0 || year > life {
		return 0
	}

	rate := factor / float64(life)
	bookValue := cost

	for i := 1; i <= year; i++ {
		depreciation := bookValue * rate

		// Don't depreciate below salvage value
		if bookValue-depreciation < salvage {
			depreciation = bookValue - salvage
		}

		if i == year {
			return RoundCurrency(depreciation)
		}

		bookValue -= depreciation
	}

	return 0
}

// DoubleDeclining calculates double declining balance depreciation.
// Convenience function with factor = 2.
func DoubleDeclining(cost, salvage float64, life, year int) float64 {
	return DecliningBalance(cost, salvage, life, year, 2)
}

// DecliningBalanceSchedule returns full declining balance schedule.
func DecliningBalanceSchedule(cost, salvage float64, life int, factor float64) []DepreciationYear {
	if life <= 0 {
		return nil
	}

	rate := factor / float64(life)
	schedule := make([]DepreciationYear, life)
	bookValue := cost
	accumulated := 0.0

	for i := 0; i < life; i++ {
		depreciation := bookValue * rate

		// Don't depreciate below salvage value
		if bookValue-depreciation < salvage {
			depreciation = bookValue - salvage
		}

		// No negative depreciation
		if depreciation < 0 {
			depreciation = 0
		}

		accumulated += depreciation
		bookValue -= depreciation

		schedule[i] = DepreciationYear{
			Year:            i + 1,
			Depreciation:    RoundCurrency(depreciation),
			AccumulatedDepr: RoundCurrency(accumulated),
			BookValue:       RoundCurrency(bookValue),
		}
	}

	return schedule
}

// DoubleDecliningSchedule returns double declining balance schedule.
func DoubleDecliningSchedule(cost, salvage float64, life int) []DepreciationYear {
	return DecliningBalanceSchedule(cost, salvage, life, 2)
}

// ════════════════════════════════════════════════════════════════
// SUM-OF-YEARS'-DIGITS DEPRECIATION
// ════════════════════════════════════════════════════════════════

// SumOfYearsDigits calculates SYD depreciation for a specific year.
// cost: initial cost
// salvage: salvage value
// life: useful life in years
// year: which year (1-indexed)
func SumOfYearsDigits(cost, salvage float64, life, year int) float64 {
	if life <= 0 || year <= 0 || year > life {
		return 0
	}

	// Sum of years = n(n+1)/2
	sumOfYears := float64(life * (life + 1) / 2)

	// Remaining life at start of year
	remainingLife := float64(life - year + 1)

	// SYD = (Cost - Salvage) * (Remaining Life / Sum of Years)
	depreciable := cost - salvage
	return RoundCurrency(depreciable * remainingLife / sumOfYears)
}

// SumOfYearsDigitsSchedule returns full SYD schedule.
func SumOfYearsDigitsSchedule(cost, salvage float64, life int) []DepreciationYear {
	if life <= 0 {
		return nil
	}

	schedule := make([]DepreciationYear, life)
	bookValue := cost
	accumulated := 0.0

	for i := 0; i < life; i++ {
		depreciation := SumOfYearsDigits(cost, salvage, life, i+1)

		accumulated += depreciation
		bookValue -= depreciation

		schedule[i] = DepreciationYear{
			Year:            i + 1,
			Depreciation:    RoundCurrency(depreciation),
			AccumulatedDepr: RoundCurrency(accumulated),
			BookValue:       RoundCurrency(bookValue),
		}
	}

	return schedule
}

// ════════════════════════════════════════════════════════════════
// UNITS OF PRODUCTION DEPRECIATION
// ════════════════════════════════════════════════════════════════

// UnitsOfProduction calculates depreciation based on usage.
// cost: initial cost
// salvage: salvage value
// totalUnits: total estimated units of production
// unitsThisPeriod: units produced in this period
func UnitsOfProduction(cost, salvage float64, totalUnits, unitsThisPeriod float64) float64 {
	if totalUnits <= 0 {
		return 0
	}

	// Depreciation per unit
	depreciationPerUnit := (cost - salvage) / totalUnits

	return RoundCurrency(depreciationPerUnit * unitsThisPeriod)
}

// ════════════════════════════════════════════════════════════════
// MACRS DEPRECIATION (US Tax)
// ════════════════════════════════════════════════════════════════

// MACRS rates for common property classes
var MACRSRates = map[int][]float64{
	3:  {0.3333, 0.4445, 0.1481, 0.0741},
	5:  {0.2000, 0.3200, 0.1920, 0.1152, 0.1152, 0.0576},
	7:  {0.1429, 0.2449, 0.1749, 0.1249, 0.0893, 0.0892, 0.0893, 0.0446},
	10: {0.1000, 0.1800, 0.1440, 0.1152, 0.0922, 0.0737, 0.0655, 0.0655, 0.0656, 0.0655, 0.0328},
	15: {0.0500, 0.0950, 0.0855, 0.0770, 0.0693, 0.0623, 0.0590, 0.0590, 0.0591, 0.0590, 0.0591, 0.0590, 0.0591, 0.0590, 0.0591, 0.0295},
}

// MACRS calculates MACRS depreciation for a specific year.
// cost: initial cost (no salvage in MACRS)
// propertyClass: 3, 5, 7, 10, or 15 year property
// year: which year (1-indexed)
func MACRS(cost float64, propertyClass, year int) float64 {
	rates, ok := MACRSRates[propertyClass]
	if !ok || year <= 0 || year > len(rates) {
		return 0
	}

	return RoundCurrency(cost * rates[year-1])
}

// MACRSSchedule returns full MACRS schedule.
func MACRSSchedule(cost float64, propertyClass int) []DepreciationYear {
	rates, ok := MACRSRates[propertyClass]
	if !ok {
		return nil
	}

	schedule := make([]DepreciationYear, len(rates))
	bookValue := cost
	accumulated := 0.0

	for i, rate := range rates {
		depreciation := RoundCurrency(cost * rate)

		accumulated += depreciation
		bookValue -= depreciation

		schedule[i] = DepreciationYear{
			Year:            i + 1,
			Depreciation:    depreciation,
			AccumulatedDepr: RoundCurrency(accumulated),
			BookValue:       RoundCurrency(bookValue),
		}
	}

	return schedule
}

// ════════════════════════════════════════════════════════════════
// DEPRECIATION TYPES
// ════════════════════════════════════════════════════════════════

// DepreciationYear holds data for one year of depreciation.
type DepreciationYear struct {
	Year            int
	Depreciation    float64
	AccumulatedDepr float64
	BookValue       float64
}

// DepreciationMethod represents depreciation calculation method.
type DepreciationMethod int

const (
	MethodStraightLine DepreciationMethod = iota
	MethodDoubleDeclining
	MethodSumOfYears
	MethodUnitsOfProduction
	MethodMACRS
)

// ════════════════════════════════════════════════════════════════
// COMPARISON
// ════════════════════════════════════════════════════════════════

// CompareDepreciation returns schedules for multiple methods.
func CompareDepreciation(cost, salvage float64, life int) map[string][]DepreciationYear {
	return map[string][]DepreciationYear{
		"straight_line":    StraightLineSchedule(cost, salvage, life),
		"double_declining": DoubleDecliningSchedule(cost, salvage, life),
		"sum_of_years":     SumOfYearsDigitsSchedule(cost, salvage, life),
	}
}

// TotalDepreciation calculates total depreciation (should equal cost - salvage).
func TotalDepreciation(schedule []DepreciationYear) float64 {
	if len(schedule) == 0 {
		return 0
	}
	return schedule[len(schedule)-1].AccumulatedDepr
}

// ════════════════════════════════════════════════════════════════
// BOOK VALUE QUERIES
// ════════════════════════════════════════════════════════════════

// BookValueStraightLine calculates book value at end of year.
func BookValueStraightLine(cost, salvage float64, life, year int) float64 {
	if year <= 0 {
		return cost
	}
	if year >= life {
		return salvage
	}

	annual := StraightLine(cost, salvage, life)
	return cost - (annual * float64(year))
}

// BookValueDeclining calculates book value using declining balance.
func BookValueDeclining(cost, salvage float64, life, year int, factor float64) float64 {
	if year <= 0 {
		return cost
	}

	rate := factor / float64(life)
	bookValue := cost

	for i := 1; i <= year && i <= life; i++ {
		depreciation := bookValue * rate
		if bookValue-depreciation < salvage {
			depreciation = bookValue - salvage
		}
		bookValue -= depreciation
	}

	return RoundCurrency(bookValue)
}

// AccumulatedDepreciation calculates total depreciation through a year.
func AccumulatedDepreciation(cost, salvage float64, life, year int) float64 {
	if year <= 0 {
		return 0
	}

	annual := StraightLine(cost, salvage, life)
	accumulated := annual * float64(year)

	// Cap at depreciable amount
	maxDepr := cost - salvage
	if accumulated > maxDepr {
		accumulated = maxDepr
	}

	return RoundCurrency(accumulated)
}
