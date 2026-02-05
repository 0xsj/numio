// pkg/types/period.go

package types

import "strings"

// Period represents a time period for rate calculations.
type Period int

const (
	PeriodNone Period = iota
	PeriodSecond
	PeriodMinute
	PeriodHour
	PeriodDay
	PeriodWeek
	PeriodMonth
	PeriodQuarter
	PeriodYear
)

// String returns the period name.
func (p Period) String() string {
	switch p {
	case PeriodSecond:
		return "second"
	case PeriodMinute:
		return "minute"
	case PeriodHour:
		return "hour"
	case PeriodDay:
		return "day"
	case PeriodWeek:
		return "week"
	case PeriodMonth:
		return "month"
	case PeriodQuarter:
		return "quarter"
	case PeriodYear:
		return "year"
	default:
		return ""
	}
}

// Symbol returns the short symbol for display (e.g., "/month").
func (p Period) Symbol() string {
	if p == PeriodNone {
		return ""
	}
	return "/" + p.String()
}

// Plural returns the plural form of the period name.
func (p Period) Plural() string {
	switch p {
	case PeriodSecond:
		return "seconds"
	case PeriodMinute:
		return "minutes"
	case PeriodHour:
		return "hours"
	case PeriodDay:
		return "days"
	case PeriodWeek:
		return "weeks"
	case PeriodMonth:
		return "months"
	case PeriodQuarter:
		return "quarters"
	case PeriodYear:
		return "years"
	default:
		return ""
	}
}

// ════════════════════════════════════════════════════════════════
// PERIOD CONVERSION
// ════════════════════════════════════════════════════════════════

// ToBase returns the number of base units (seconds) in this period.
// Used for converting between periods.
func (p Period) ToBase() float64 {
	switch p {
	case PeriodSecond:
		return 1
	case PeriodMinute:
		return 60
	case PeriodHour:
		return 3600
	case PeriodDay:
		return 86400
	case PeriodWeek:
		return 604800
	case PeriodMonth:
		return 2629746 // Average month (30.44 days)
	case PeriodQuarter:
		return 7889238 // 3 months
	case PeriodYear:
		return 31556952 // Average year (365.25 days)
	default:
		return 1
	}
}

// ConvertTo converts a count of this period to another period.
// Example: 1 year → 12 months, 1 week → 7 days
func (p Period) ConvertTo(count float64, target Period) float64 {
	if p == target || p == PeriodNone || target == PeriodNone {
		return count
	}
	// Convert via base units (seconds)
	baseUnits := count * p.ToBase()
	return baseUnits / target.ToBase()
}

// InMonths returns how many months this period represents.
// Useful for financial calculations where month is the common base.
func (p Period) InMonths() float64 {
	switch p {
	case PeriodSecond:
		return 1.0 / 2629746
	case PeriodMinute:
		return 1.0 / 43829.1
	case PeriodHour:
		return 1.0 / 730.485
	case PeriodDay:
		return 1.0 / 30.4375
	case PeriodWeek:
		return 7.0 / 30.4375
	case PeriodMonth:
		return 1
	case PeriodQuarter:
		return 3
	case PeriodYear:
		return 12
	default:
		return 1
	}
}

// InDays returns how many days this period represents.
func (p Period) InDays() float64 {
	switch p {
	case PeriodSecond:
		return 1.0 / 86400
	case PeriodMinute:
		return 1.0 / 1440
	case PeriodHour:
		return 1.0 / 24
	case PeriodDay:
		return 1
	case PeriodWeek:
		return 7
	case PeriodMonth:
		return 30.4375
	case PeriodQuarter:
		return 91.3125
	case PeriodYear:
		return 365.25
	default:
		return 1
	}
}

// ════════════════════════════════════════════════════════════════
// PERIOD PARSING
// ════════════════════════════════════════════════════════════════

// ParsePeriod parses a string into a Period.
// Accepts singular, plural, and abbreviations.
func ParsePeriod(s string) Period {
	switch strings.ToLower(s) {
	case "s", "sec", "second", "seconds":
		return PeriodSecond
	case "m", "min", "minute", "minutes":
		return PeriodMinute
	case "h", "hr", "hour", "hours":
		return PeriodHour
	case "d", "day", "days":
		return PeriodDay
	case "w", "wk", "week", "weeks":
		return PeriodWeek
	case "mo", "month", "months":
		return PeriodMonth
	case "q", "qtr", "quarter", "quarters":
		return PeriodQuarter
	case "y", "yr", "year", "years":
		return PeriodYear
	default:
		return PeriodNone
	}
}

// IsPeriod checks if a string is a valid period name.
func IsPeriod(s string) bool {
	return ParsePeriod(s) != PeriodNone
}

// ════════════════════════════════════════════════════════════════
// PERIOD CONSTANTS (for use as multipliers)
// ════════════════════════════════════════════════════════════════

// PeriodMultiplier returns the default multiplier for a period.
// When used as a bare constant (e.g., `rent * year`), these are the values.
// Base unit is assumed to be "per month" for financial contexts.
func (p Period) DefaultMultiplier() float64 {
	switch p {
	case PeriodSecond:
		return 1.0 / 2629746 // seconds in a month
	case PeriodMinute:
		return 1.0 / 43829.1 // minutes in a month
	case PeriodHour:
		return 1.0 / 730.485 // hours in a month
	case PeriodDay:
		return 1.0 / 30.4375 // days in a month
	case PeriodWeek:
		return 7.0 / 30.4375 // ~0.23 months
	case PeriodMonth:
		return 1
	case PeriodQuarter:
		return 3
	case PeriodYear:
		return 12
	default:
		return 1
	}
}

// ════════════════════════════════════════════════════════════════
// ALL PERIODS
// ════════════════════════════════════════════════════════════════

// AllPeriods returns all valid periods.
func AllPeriods() []Period {
	return []Period{
		PeriodSecond,
		PeriodMinute,
		PeriodHour,
		PeriodDay,
		PeriodWeek,
		PeriodMonth,
		PeriodQuarter,
		PeriodYear,
	}
}

// PeriodNames returns all period names (for autocompletion, etc.).
func PeriodNames() []string {
	return []string{
		"second", "seconds",
		"minute", "minutes",
		"hour", "hours",
		"day", "days",
		"week", "weeks",
		"month", "months",
		"quarter", "quarters",
		"year", "years",
	}
}

// ParsePeriodStrict parses a string into a Period, but excludes single-letter
// abbreviations (s, m, h, d, w, q, y) that commonly collide with variable names.
// Use this in contexts where the input could be either a variable or a period
// (e.g., standalone identifiers in expressions).
// For unambiguous contexts like rate suffixes (/y, per m), use ParsePeriod instead.
func ParsePeriodStrict(s string) Period {
	switch strings.ToLower(s) {
	case "sec", "second", "seconds":
		return PeriodSecond
	case "min", "minute", "minutes":
		return PeriodMinute
	case "hr", "hour", "hours":
		return PeriodHour
	case "day", "days":
		return PeriodDay
	case "wk", "week", "weeks":
		return PeriodWeek
	case "mo", "month", "months":
		return PeriodMonth
	case "qtr", "quarter", "quarters":
		return PeriodQuarter
	case "yr", "year", "years":
		return PeriodYear
	default:
		return PeriodNone
	}
}