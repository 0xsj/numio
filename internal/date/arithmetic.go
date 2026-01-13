// internal/date/arithmetic.go

package date

import (
	"time"
)

// ════════════════════════════════════════════════════════════════
// ADD OPERATIONS
// ════════════════════════════════════════════════════════════════

// AddDays adds days to the date.
func (d Date) AddDays(days int) Date {
	return Date{d.Time.AddDate(0, 0, days)}
}

// AddWeeks adds weeks to the date.
func (d Date) AddWeeks(weeks int) Date {
	return Date{d.Time.AddDate(0, 0, weeks*7)}
}

// AddMonths adds months to the date.
func (d Date) AddMonths(months int) Date {
	return Date{d.Time.AddDate(0, months, 0)}
}

// AddYears adds years to the date.
func (d Date) AddYears(years int) Date {
	return Date{d.Time.AddDate(years, 0, 0)}
}

// AddHours adds hours to the date.
func (d Date) AddHours(hours int) Date {
	return Date{d.Time.Add(time.Duration(hours) * time.Hour)}
}

// AddMinutes adds minutes to the date.
func (d Date) AddMinutes(minutes int) Date {
	return Date{d.Time.Add(time.Duration(minutes) * time.Minute)}
}

// AddSeconds adds seconds to the date.
func (d Date) AddSeconds(seconds int) Date {
	return Date{d.Time.Add(time.Duration(seconds) * time.Second)}
}

// AddDuration adds a time.Duration to the date.
func (d Date) AddDuration(duration time.Duration) Date {
	return Date{d.Time.Add(duration)}
}

// ════════════════════════════════════════════════════════════════
// SUBTRACT OPERATIONS (convenience wrappers)
// ════════════════════════════════════════════════════════════════

// SubDays subtracts days from the date.
func (d Date) SubDays(days int) Date {
	return d.AddDays(-days)
}

// SubWeeks subtracts weeks from the date.
func (d Date) SubWeeks(weeks int) Date {
	return d.AddWeeks(-weeks)
}

// SubMonths subtracts months from the date.
func (d Date) SubMonths(months int) Date {
	return d.AddMonths(-months)
}

// SubYears subtracts years from the date.
func (d Date) SubYears(years int) Date {
	return d.AddYears(-years)
}

// SubHours subtracts hours from the date.
func (d Date) SubHours(hours int) Date {
	return d.AddHours(-hours)
}

// SubMinutes subtracts minutes from the date.
func (d Date) SubMinutes(minutes int) Date {
	return d.AddMinutes(-minutes)
}

// SubSeconds subtracts seconds from the date.
func (d Date) SubSeconds(seconds int) Date {
	return d.AddSeconds(-seconds)
}

// ════════════════════════════════════════════════════════════════
// DIFFERENCE CALCULATIONS
// ════════════════════════════════════════════════════════════════

// DaysBetween returns the number of days between two dates.
func DaysBetween(d1, d2 Date) int {
	// Normalize to midnight for accurate day count
	t1 := time.Date(d1.Year(), d1.Month(), d1.Day(), 0, 0, 0, 0, time.UTC)
	t2 := time.Date(d2.Year(), d2.Month(), d2.Day(), 0, 0, 0, 0, time.UTC)

	duration := t2.Sub(t1)
	return int(duration.Hours() / 24)
}

// WeeksBetween returns the number of complete weeks between two dates.
func WeeksBetween(d1, d2 Date) int {
	return DaysBetween(d1, d2) / 7
}

// MonthsBetween returns the number of months between two dates.
func MonthsBetween(d1, d2 Date) int {
	years := d2.Year() - d1.Year()
	months := int(d2.Month()) - int(d1.Month())

	total := years*12 + months

	// Adjust if day of month hasn't been reached
	if d2.Day() < d1.Day() {
		total--
	}

	return total
}

// YearsBetween returns the number of complete years between two dates.
func YearsBetween(d1, d2 Date) int {
	years := d2.Year() - d1.Year()

	// Adjust if anniversary hasn't been reached
	if d2.Month() < d1.Month() || (d2.Month() == d1.Month() && d2.Day() < d1.Day()) {
		years--
	}

	return years
}

// HoursBetween returns the number of hours between two dates.
func HoursBetween(d1, d2 Date) int {
	duration := d2.Time.Sub(d1.Time)
	return int(duration.Hours())
}

// MinutesBetween returns the number of minutes between two dates.
func MinutesBetween(d1, d2 Date) int {
	duration := d2.Time.Sub(d1.Time)
	return int(duration.Minutes())
}

// SecondsBetween returns the number of seconds between two dates.
func SecondsBetween(d1, d2 Date) int64 {
	duration := d2.Time.Sub(d1.Time)
	return int64(duration.Seconds())
}

// Duration returns the time.Duration between two dates.
func Duration(d1, d2 Date) time.Duration {
	return d2.Time.Sub(d1.Time)
}

// ════════════════════════════════════════════════════════════════
// AGE CALCULATION
// ════════════════════════════════════════════════════════════════

// Age calculates the age in years from a birth date to today.
func Age(birthDate Date) int {
	return YearsBetween(birthDate, Today())
}

// AgeAt calculates the age in years from a birth date to a specific date.
func AgeAt(birthDate, atDate Date) int {
	return YearsBetween(birthDate, atDate)
}

// AgeDetailed returns years, months, and days.
type AgeDetail struct {
	Years  int
	Months int
	Days   int
}

// AgeDetailedAt calculates detailed age from birth date to a specific date.
func AgeDetailedAt(birthDate, atDate Date) AgeDetail {
	years := YearsBetween(birthDate, atDate)

	// Calculate remaining months
	afterYears := birthDate.AddYears(years)
	months := MonthsBetween(afterYears, atDate)

	// Calculate remaining days
	afterMonths := afterYears.AddMonths(months)
	days := DaysBetween(afterMonths, atDate)

	return AgeDetail{
		Years:  years,
		Months: months,
		Days:   days,
	}
}

// ════════════════════════════════════════════════════════════════
// WORKING DAYS
// ════════════════════════════════════════════════════════════════

// AddWorkdays adds working days (excluding weekends).
func (d Date) AddWorkdays(days int) Date {
	result := d
	direction := 1
	if days < 0 {
		direction = -1
		days = -days
	}

	for days > 0 {
		result = result.AddDays(direction)
		if result.Weekday() != time.Saturday && result.Weekday() != time.Sunday {
			days--
		}
	}

	return result
}

// WorkdaysBetween returns the number of working days between two dates.
func WorkdaysBetween(d1, d2 Date) int {
	if d1.After(d2) {
		d1, d2 = d2, d1
	}

	count := 0
	current := d1

	for !current.After(d2) {
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			count++
		}
		current = current.AddDays(1)
	}

	return count
}

// ════════════════════════════════════════════════════════════════
// ROUNDING
// ════════════════════════════════════════════════════════════════

// StartOfDay returns the date at midnight.
func (d Date) StartOfDay() Date {
	return New(d.Year(), int(d.Month()), d.Day())
}

// EndOfDay returns the date at 23:59:59.
func (d Date) EndOfDay() Date {
	return NewDateTime(d.Year(), int(d.Month()), d.Day(), 23, 59, 59)
}

// StartOfMonth returns the first day of the month.
func (d Date) StartOfMonth() Date {
	return New(d.Year(), int(d.Month()), 1)
}

// EndOfMonth returns the last day of the month.
func (d Date) EndOfMonth() Date {
	return New(d.Year(), int(d.Month())+1, 0)
}

// StartOfYear returns January 1 of the year.
func (d Date) StartOfYear() Date {
	return New(d.Year(), 1, 1)
}

// EndOfYear returns December 31 of the year.
func (d Date) EndOfYear() Date {
	return New(d.Year(), 12, 31)
}

// StartOfWeek returns the Monday of the current week.
func (d Date) StartOfWeek() Date {
	weekday := int(d.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	return d.AddDays(-(weekday - 1))
}

// EndOfWeek returns the Sunday of the current week.
func (d Date) EndOfWeek() Date {
	return d.StartOfWeek().AddDays(6)
}

// StartOfQuarter returns the first day of the quarter.
func (d Date) StartOfQuarter() Date {
	quarter := (int(d.Month()) - 1) / 3
	return New(d.Year(), quarter*3+1, 1)
}

// EndOfQuarter returns the last day of the quarter.
func (d Date) EndOfQuarter() Date {
	return d.StartOfQuarter().AddMonths(3).AddDays(-1)
}
