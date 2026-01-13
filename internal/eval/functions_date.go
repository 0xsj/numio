// internal/eval/functions_date.go

package eval

import (
	"strings"
	"time"

	"github.com/0xsj/numio/internal/date"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// CURRENT DATE/TIME
// ════════════════════════════════════════════════════════════════

// FnToday returns today's date.
func FnToday(args []types.Value) types.Value {
	return types.DateValue(date.Today().Time)
}

// FnNow returns the current date and time.
func FnNow(args []types.Value) types.Value {
	return types.DateValue(date.Now().Time)
}

// FnYesterday returns yesterday's date.
func FnYesterday(args []types.Value) types.Value {
	return types.DateValue(date.Yesterday().Time)
}

// FnTomorrow returns tomorrow's date.
func FnTomorrow(args []types.Value) types.Value {
	return types.DateValue(date.Tomorrow().Time)
}

// ════════════════════════════════════════════════════════════════
// DATE CREATION
// ════════════════════════════════════════════════════════════════

// FnDate creates a date from year, month, day.
// Args: year, month, day
func FnDate(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("date requires 3 arguments: year, month, day")
	}

	year := int(args[0].AsFloat())
	month := int(args[1].AsFloat())
	day := int(args[2].AsFloat())

	if !date.IsValidDate(year, month, day) {
		return types.Error("invalid date")
	}

	return types.DateValueFromComponents(year, month, day)
}

// FnDateTime creates a date with time from components.
// Args: year, month, day, hour, minute [, second]
func FnDateTime(args []types.Value) types.Value {
	if len(args) < 5 || len(args) > 6 {
		return types.Error("datetime requires 5-6 arguments: year, month, day, hour, minute [, second]")
	}

	year := int(args[0].AsFloat())
	month := int(args[1].AsFloat())
	day := int(args[2].AsFloat())
	hour := int(args[3].AsFloat())
	minute := int(args[4].AsFloat())
	second := 0
	if len(args) == 6 {
		second = int(args[5].AsFloat())
	}

	if !date.IsValidDate(year, month, day) {
		return types.Error("invalid date")
	}
	if !date.IsValidTime(hour, minute, second) {
		return types.Error("invalid time")
	}

	return types.DateTimeValue(year, month, day, hour, minute, second)
}

// FnTime creates a time on today's date.
// Args: hour, minute [, second]
func FnTime(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("time requires 2-3 arguments: hour, minute [, second]")
	}

	hour := int(args[0].AsFloat())
	minute := int(args[1].AsFloat())
	second := 0
	if len(args) == 3 {
		second = int(args[2].AsFloat())
	}

	if !date.IsValidTime(hour, minute, second) {
		return types.Error("invalid time")
	}

	today := date.Today()
	return types.DateTimeValue(today.Year(), today.MonthInt(), today.Day(), hour, minute, second)
}

// FnUnixToDate converts Unix timestamp to date.
// Args: timestamp (seconds)
func FnUnixToDate(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("unixtodate requires 1 argument: timestamp")
	}

	timestamp := int64(args[0].AsFloat())
	d := date.FromUnix(timestamp)
	return types.DateValue(d.Time)
}

// FnDateToUnix converts date to Unix timestamp.
// Args: date value
func FnDateToUnix(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("datetounix requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("datetounix requires a date argument")
	}

	t := args[0].AsTime()
	return types.Number(float64(t.Unix()))
}

// ════════════════════════════════════════════════════════════════
// DATE ARITHMETIC
// ════════════════════════════════════════════════════════════════

// FnAddDays adds days to a date.
// Args: date, days
func FnAddDays(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("adddays requires 2 arguments: date, days")
	}

	if !args[0].IsDate() {
		return types.Error("adddays requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	days := int(args[1].AsFloat())
	result := d.AddDays(days)
	return types.DateValue(result.Time)
}

// FnAddWeeks adds weeks to a date.
// Args: date, weeks
func FnAddWeeks(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addweeks requires 2 arguments: date, weeks")
	}

	if !args[0].IsDate() {
		return types.Error("addweeks requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	weeks := int(args[1].AsFloat())
	result := d.AddWeeks(weeks)
	return types.DateValue(result.Time)
}

// FnAddMonths adds months to a date.
// Args: date, months
func FnAddMonths(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addmonths requires 2 arguments: date, months")
	}

	if !args[0].IsDate() {
		return types.Error("addmonths requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	months := int(args[1].AsFloat())
	result := d.AddMonths(months)
	return types.DateValue(result.Time)
}

// FnAddYears adds years to a date.
// Args: date, years
func FnAddYears(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addyears requires 2 arguments: date, years")
	}

	if !args[0].IsDate() {
		return types.Error("addyears requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	years := int(args[1].AsFloat())
	result := d.AddYears(years)
	return types.DateValue(result.Time)
}

// FnAddHours adds hours to a date.
// Args: date, hours
func FnAddHours(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addhours requires 2 arguments: date, hours")
	}

	if !args[0].IsDate() {
		return types.Error("addhours requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	hours := int(args[1].AsFloat())
	result := d.AddHours(hours)
	return types.DateValue(result.Time)
}

// FnAddMinutes adds minutes to a date.
// Args: date, minutes
func FnAddMinutes(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addminutes requires 2 arguments: date, minutes")
	}

	if !args[0].IsDate() {
		return types.Error("addminutes requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	minutes := int(args[1].AsFloat())
	result := d.AddMinutes(minutes)
	return types.DateValue(result.Time)
}

// FnAddSeconds adds seconds to a date.
// Args: date, seconds
func FnAddSeconds(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addseconds requires 2 arguments: date, seconds")
	}

	if !args[0].IsDate() {
		return types.Error("addseconds requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	seconds := int(args[1].AsFloat())
	result := d.AddSeconds(seconds)
	return types.DateValue(result.Time)
}

// FnAddWorkdays adds working days (excluding weekends).
// Args: date, workdays
func FnAddWorkdays(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("addworkdays requires 2 arguments: date, workdays")
	}

	if !args[0].IsDate() {
		return types.Error("addworkdays requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	days := int(args[1].AsFloat())
	result := d.AddWorkdays(days)
	return types.DateValue(result.Time)
}

// ════════════════════════════════════════════════════════════════
// DATE DIFFERENCE
// ════════════════════════════════════════════════════════════════

// FnDaysBetween returns days between two dates.
// Args: date1, date2
func FnDaysBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("daysbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("daysbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.DaysBetween(d1, d2)))
}

// FnWeeksBetween returns weeks between two dates.
// Args: date1, date2
func FnWeeksBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("weeksbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("weeksbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.WeeksBetween(d1, d2)))
}

// FnMonthsBetween returns months between two dates.
// Args: date1, date2
func FnMonthsBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("monthsbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("monthsbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.MonthsBetween(d1, d2)))
}

// FnYearsBetween returns years between two dates.
// Args: date1, date2
func FnYearsBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("yearsbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("yearsbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.YearsBetween(d1, d2)))
}

// FnHoursBetween returns hours between two dates.
// Args: date1, date2
func FnHoursBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("hoursbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("hoursbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.HoursBetween(d1, d2)))
}

// FnMinutesBetween returns minutes between two dates.
// Args: date1, date2
func FnMinutesBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("minutesbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("minutesbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.MinutesBetween(d1, d2)))
}

// FnSecondsBetween returns seconds between two dates.
// Args: date1, date2
func FnSecondsBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("secondsbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("secondsbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.SecondsBetween(d1, d2)))
}

// FnWorkdaysBetween returns working days between two dates.
// Args: date1, date2
func FnWorkdaysBetween(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("workdaysbetween requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("workdaysbetween requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())
	return types.Number(float64(date.WorkdaysBetween(d1, d2)))
}

// ════════════════════════════════════════════════════════════════
// DATE COMPONENTS
// ════════════════════════════════════════════════════════════════

// FnYear extracts the year from a date.
func FnYear(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("year requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("year requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.Year()))
}

// FnMonth extracts the month from a date (1-12).
func FnMonth(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("month requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("month requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.MonthInt()))
}

// FnDay extracts the day of month from a date (1-31).
func FnDay(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("day requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("day requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.Day()))
}

// FnHour extracts the hour from a date (0-23).
func FnHour(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hour requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("hour requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.Hour()))
}

// FnMinute extracts the minute from a date (0-59).
func FnMinute(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("minute requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("minute requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.Minute()))
}

// FnSecond extracts the second from a date (0-59).
func FnSecond(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("second requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("second requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.Second()))
}

// FnWeekday returns the day of week (0=Sunday, 6=Saturday).
func FnWeekday(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("weekday requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("weekday requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.WeekdayInt()))
}

// FnWeekdayISO returns the ISO day of week (1=Monday, 7=Sunday).
func FnWeekdayISO(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("weekdayiso requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("weekdayiso requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.WeekdayISO()))
}

// FnWeekdayName returns the weekday name.
func FnWeekdayName(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("weekdayname requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("weekdayname requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.StringValue(d.WeekdayName())
}

// FnMonthName returns the month name.
func FnMonthName(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("monthname requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("monthname requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.StringValue(d.MonthName())
}

// FnDayOfYear returns the day of year (1-366).
func FnDayOfYear(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("dayofyear requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("dayofyear requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.DayOfYear()))
}

// FnWeekOfYear returns the ISO week number (1-53).
func FnWeekOfYear(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("weekofyear requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("weekofyear requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.WeekOfYear()))
}

// FnQuarter returns the quarter (1-4).
func FnQuarter(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("quarter requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("quarter requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(float64(d.Quarter()))
}

// ════════════════════════════════════════════════════════════════
// DATE UTILITIES
// ════════════════════════════════════════════════════════════════

// FnIsLeapYear checks if the year is a leap year.
func FnIsLeapYear(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isleapyear requires 1 argument: year or date")
	}

	var year int
	if args[0].IsDate() {
		year = date.FromTime(args[0].AsTime()).Year()
	} else {
		year = int(args[0].AsFloat())
	}

	if date.IsLeapYear(year) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnDaysInMonth returns the number of days in the month.
func FnDaysInMonth(args []types.Value) types.Value {
	if len(args) == 1 {
		if args[0].IsDate() {
			d := date.FromTime(args[0].AsTime())
			return types.Number(float64(d.DaysInMonth()))
		}
		// Assume current year
		month := int(args[0].AsFloat())
		return types.Number(float64(date.DaysInMonth(month)))
	}

	if len(args) == 2 {
		year := int(args[0].AsFloat())
		month := int(args[1].AsFloat())
		return types.Number(float64(date.DaysInMonthOf(year, month)))
	}

	return types.Error("daysinmonth requires 1-2 arguments: date or [year,] month")
}

// FnIsWeekend checks if the date is a weekend.
func FnIsWeekend(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isweekend requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("isweekend requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	if d.IsWeekend() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsWeekday checks if the date is a weekday.
func FnIsWeekday(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isweekday requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("isweekday requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	if d.IsWeekday() {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnAge calculates age in years from a birth date.
func FnAge(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("age requires 1 argument: birthdate")
	}

	if !args[0].IsDate() {
		return types.Error("age requires a date argument")
	}

	birthDate := date.FromTime(args[0].AsTime())
	return types.Number(float64(date.Age(birthDate)))
}

// ════════════════════════════════════════════════════════════════
// DATE ROUNDING
// ════════════════════════════════════════════════════════════════

// FnStartOfDay returns the date at midnight.
func FnStartOfDay(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("startofday requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("startofday requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.StartOfDay().Time)
}

// FnEndOfDay returns the date at 23:59:59.
func FnEndOfDay(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("endofday requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("endofday requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.EndOfDay().Time)
}

// FnStartOfMonth returns the first day of the month.
func FnStartOfMonth(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("startofmonth requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("startofmonth requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.StartOfMonth().Time)
}

// FnEndOfMonth returns the last day of the month.
func FnEndOfMonth(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("endofmonth requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("endofmonth requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.EndOfMonth().Time)
}

// FnStartOfYear returns January 1 of the year.
func FnStartOfYear(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("startofyear requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("startofyear requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.StartOfYear().Time)
}

// FnEndOfYear returns December 31 of the year.
func FnEndOfYear(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("endofyear requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("endofyear requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.EndOfYear().Time)
}

// FnStartOfWeek returns Monday of the current week.
func FnStartOfWeek(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("startofweek requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("startofweek requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.StartOfWeek().Time)
}

// FnEndOfWeek returns Sunday of the current week.
func FnEndOfWeek(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("endofweek requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("endofweek requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.EndOfWeek().Time)
}

// ════════════════════════════════════════════════════════════════
// TIMEZONE
// ════════════════════════════════════════════════════════════════

// FnInTimezone converts a date to a different timezone.
// Args: date, timezone (e.g., "PST", "UTC", "America/New_York")
func FnInTimezone(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("intimezone requires 2 arguments: date, timezone")
	}

	if !args[0].IsDate() {
		return types.Errorf("intimezone requires a date as first argument, got %s", args[0].Kind.String())
	}

	// Debug: check what we're getting for the timezone arg
	tzArg := args[1]
	var tz string

	if tzArg.IsString() {
		tz = tzArg.Str
	} else {
		// Maybe it's coming through differently
		tz = tzArg.AsString()
	}

	// Strip quotes if present (parser might include them)
	tz = stripQuotes(tz)

	if tz == "" {
		return types.Error("intimezone requires a timezone string")
	}

	d := date.FromTime(args[0].AsTime())

	result, err := d.InTimezone(tz)
	if err != nil {
		return types.Errorf("invalid timezone: %s", tz)
	}

	return types.DateValue(result.Time)
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// FnToUTC converts a date to UTC.
func FnToUTC(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("toutc requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("toutc requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.ToUTC().Time)
}

// FnToLocal converts a date to local timezone.
func FnToLocal(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tolocal requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("tolocal requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.DateValue(d.ToLocal().Time)
}

// FnTimezoneOffset returns the timezone offset in hours.
func FnTimezoneOffset(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("timezoneoffset requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("timezoneoffset requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.Number(d.TimezoneOffset())
}

// ════════════════════════════════════════════════════════════════
// FORMATTING
// ════════════════════════════════════════════════════════════════

// FnFormatDate formats a date with a custom pattern.
// Args: date, format
// Format uses Go layout: 2006-01-02 15:04:05
func FnFormatDate(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("formatdate requires 2 arguments: date, format")
	}

	if !args[0].IsDate() {
		return types.Error("formatdate requires a date as first argument")
	}

	d := date.FromTime(args[0].AsTime())
	format := args[1].AsString()

	// Convert common format tokens to Go layout
	format = convertDateFormat(format)

	return types.StringValue(d.FormatCustom(format))
}

// FnDateISO returns the date in ISO 8601 format.
func FnDateISO(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("dateiso requires 1 argument: date")
	}

	if !args[0].IsDate() {
		return types.Error("dateiso requires a date argument")
	}

	d := date.FromTime(args[0].AsTime())
	return types.StringValue(d.ISO())
}

// convertDateFormat converts common format tokens to Go layout.
func convertDateFormat(format string) string {
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"hh":   "03",
		"mm":   "04",
		"ss":   "05",
		"SSS":  "000",
		"A":    "PM",
		"a":    "pm",
	}

	result := format
	for token, goLayout := range replacements {
		result = strings.ReplaceAll(result, token, goLayout)
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// COMPARISON
// ════════════════════════════════════════════════════════════════

// FnIsBefore checks if first date is before second date.
func FnIsBefore(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("isbefore requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("isbefore requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())

	if d1.Before(d2) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsAfter checks if first date is after second date.
func FnIsAfter(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("isafter requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("isafter requires two date arguments")
	}

	d1 := date.FromTime(args[0].AsTime())
	d2 := date.FromTime(args[1].AsTime())

	if d1.After(d2) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsSameDay checks if two dates are on the same day.
func FnIsSameDay(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("issameday requires 2 arguments: date1, date2")
	}

	if !args[0].IsDate() || !args[1].IsDate() {
		return types.Error("issameday requires two date arguments")
	}

	t1 := args[0].AsTime()
	t2 := args[1].AsTime()

	if t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day() {
		return types.Number(1)
	}
	return types.Number(0)
}

// ════════════════════════════════════════════════════════════════
// SPECIAL
// ════════════════════════════════════════════════════════════════

// FnEpoch returns the Unix epoch date (1970-01-01).
func FnEpoch(args []types.Value) types.Value {
	return types.DateValue(date.UnixEpoch().Time)
}

// FnNowInTimezone returns current time in a specific timezone.
// Args: timezone
func FnNowInTimezone(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("nowin requires 1 argument: timezone")
	}

	var tz string
	if args[0].IsString() {
		tz = args[0].Str
	} else {
		tz = args[0].AsString()
	}

	// Strip quotes if present
	tz = stripQuotes(tz)

	if tz == "" {
		return types.Error("nowin requires a timezone string")
	}

	loc, err := date.GetTimezone(tz)
	if err != nil {
		return types.Errorf("invalid timezone: %s", tz)
	}

	return types.DateValue(time.Now().In(loc))
}
