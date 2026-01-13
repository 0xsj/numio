// internal/date/utils.go

package date

import (
	"time"
)

// ════════════════════════════════════════════════════════════════
// LEAP YEAR
// ════════════════════════════════════════════════════════════════

// IsLeapYear checks if a year is a leap year.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsLeap checks if the date's year is a leap year.
func (d Date) IsLeap() bool {
	return IsLeapYear(d.Year())
}

// ════════════════════════════════════════════════════════════════
// DAYS IN MONTH
// ════════════════════════════════════════════════════════════════

// DaysInMonthOf returns the number of days in a given month.
func DaysInMonthOf(year, month int) int {
	// Handle month overflow/underflow
	if month < 1 {
		month = 1
	}
	if month > 12 {
		month = 12
	}

	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 30
	}
}

// DaysInMonth returns the number of days in the given month of current year.
func DaysInMonth(month int) int {
	return DaysInMonthOf(time.Now().Year(), month)
}

// ════════════════════════════════════════════════════════════════
// VALIDATION
// ════════════════════════════════════════════════════════════════

// IsValidDate checks if year, month, day form a valid date.
func IsValidDate(year, month, day int) bool {
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > DaysInMonthOf(year, month) {
		return false
	}
	return true
}

// IsValidTime checks if hour, minute, second form a valid time.
func IsValidTime(hour, minute, second int) bool {
	if hour < 0 || hour > 23 {
		return false
	}
	if minute < 0 || minute > 59 {
		return false
	}
	if second < 0 || second > 59 {
		return false
	}
	return true
}

// ════════════════════════════════════════════════════════════════
// SPECIAL DATES
// ════════════════════════════════════════════════════════════════

// Yesterday returns yesterday's date.
func Yesterday() Date {
	return Today().AddDays(-1)
}

// Tomorrow returns tomorrow's date.
func Tomorrow() Date {
	return Today().AddDays(1)
}

// FirstDayOfMonth returns the first day of the given month.
func FirstDayOfMonth(year, month int) Date {
	return New(year, month, 1)
}

// LastDayOfMonth returns the last day of the given month.
func LastDayOfMonth(year, month int) Date {
	return New(year, month, DaysInMonthOf(year, month))
}

// FirstDayOfYear returns January 1 of the given year.
func FirstDayOfYear(year int) Date {
	return New(year, 1, 1)
}

// LastDayOfYear returns December 31 of the given year.
func LastDayOfYear(year int) Date {
	return New(year, 12, 31)
}

// ════════════════════════════════════════════════════════════════
// NTH WEEKDAY
// ════════════════════════════════════════════════════════════════

// NthWeekdayOfMonth returns the nth occurrence of a weekday in a month.
// For example, NthWeekdayOfMonth(2024, 11, time.Thursday, 4) returns the 4th Thursday of November 2024.
func NthWeekdayOfMonth(year, month int, weekday time.Weekday, n int) Date {
	if n < 1 || n > 5 {
		return Date{}
	}

	first := FirstDayOfMonth(year, month)
	firstWeekday := int(first.Weekday())
	targetWeekday := int(weekday)

	// Days until first occurrence of target weekday
	daysUntil := (targetWeekday - firstWeekday + 7) % 7

	// Day of the nth occurrence
	day := 1 + daysUntil + (n-1)*7

	if day > DaysInMonthOf(year, month) {
		return Date{} // Invalid, doesn't exist
	}

	return New(year, month, day)
}

// LastWeekdayOfMonth returns the last occurrence of a weekday in a month.
func LastWeekdayOfMonth(year, month int, weekday time.Weekday) Date {
	last := LastDayOfMonth(year, month)
	lastWeekday := int(last.Weekday())
	targetWeekday := int(weekday)

	// Days back to target weekday
	daysBack := (lastWeekday - targetWeekday + 7) % 7

	return last.AddDays(-daysBack)
}

// ════════════════════════════════════════════════════════════════
// RELATIVE DESCRIPTIONS
// ════════════════════════════════════════════════════════════════

// RelativeDays returns a human-readable relative description.
func (d Date) RelativeDays() string {
	today := Today()
	days := DaysBetween(today, d)

	switch days {
	case -1:
		return "yesterday"
	case 0:
		return "today"
	case 1:
		return "tomorrow"
	default:
		if days < 0 {
			return intToStr(-days) + " days ago"
		}
		return "in " + intToStr(days) + " days"
	}
}

// RelativeTime returns a human-readable relative time description.
func (d Date) RelativeTime() string {
	now := Now()
	seconds := SecondsBetween(now, d)

	if seconds < 0 {
		seconds = -seconds
		return formatDuration(seconds) + " ago"
	}
	if seconds == 0 {
		return "now"
	}
	return "in " + formatDuration(seconds)
}

func formatDuration(seconds int64) string {
	if seconds < 60 {
		return intToStr64(seconds) + " seconds"
	}
	if seconds < 3600 {
		return intToStr64(seconds/60) + " minutes"
	}
	if seconds < 86400 {
		return intToStr64(seconds/3600) + " hours"
	}
	if seconds < 604800 {
		return intToStr64(seconds/86400) + " days"
	}
	if seconds < 2592000 {
		return intToStr64(seconds/604800) + " weeks"
	}
	if seconds < 31536000 {
		return intToStr64(seconds/2592000) + " months"
	}
	return intToStr64(seconds/31536000) + " years"
}

// ════════════════════════════════════════════════════════════════
// EPOCH & JULIAN
// ════════════════════════════════════════════════════════════════

// UnixEpoch returns the Unix epoch (1970-01-01).
func UnixEpoch() Date {
	return New(1970, 1, 1)
}

// JulianDay returns the Julian Day Number.
func (d Date) JulianDay() float64 {
	y := float64(d.Year())
	m := float64(d.MonthInt())
	day := float64(d.Day())

	if m <= 2 {
		y--
		m += 12
	}

	a := float64(int(y / 100))
	b := 2 - a + float64(int(a/4))

	jd := float64(int(365.25*(y+4716))) + float64(int(30.6001*(m+1))) + day + b - 1524.5

	// Add time fraction
	jd += (float64(d.Hour()) + float64(d.Minute())/60 + float64(d.Second())/3600) / 24

	return jd
}

// FromJulianDay creates a Date from Julian Day Number.
func FromJulianDay(jd float64) Date {
	jd += 0.5
	z := int(jd)
	f := jd - float64(z)

	var a int
	if z < 2299161 {
		a = z
	} else {
		alpha := int((float64(z) - 1867216.25) / 36524.25)
		a = z + 1 + alpha - alpha/4
	}

	b := a + 1524
	c := int((float64(b) - 122.1) / 365.25)
	d := int(365.25 * float64(c))
	e := int(float64(b-d) / 30.6001)

	day := b - d - int(30.6001*float64(e))

	var month int
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}

	var year int
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	// Extract time from fraction
	f *= 24
	hour := int(f)
	f = (f - float64(hour)) * 60
	minute := int(f)
	second := int((f - float64(minute)) * 60)

	return NewDateTime(year, month, day, hour, minute, second)
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}

func intToStr64(n int64) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
