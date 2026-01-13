// internal/nlp/datetime.go

package nlp

import (
	"regexp"
	"strings"
)

// ════════════════════════════════════════════════════════════════
// DATE/TIME PATTERNS
// ════════════════════════════════════════════════════════════════

// RegisterDateTimePatterns registers all date/time patterns.
func RegisterDateTimePatterns(r *PatternRegistry) {
	// Register in priority order (most specific first)
	r.RegisterAll([]*Pattern{
		// High priority - specific patterns
		patternNowInTimezone(),
		patternTimeInTimezone(),
		patternWhatTimeIn(),

		// Date arithmetic
		patternTodayPlusDays(),
		patternNowPlusTime(),
		patternDateFromNow(),
		patternDateAgo(),

		// Relative days
		patternNextWeekday(),
		patternLastWeekday(),
		patternThisWeekday(),

		// Date references
		patternDaysSince(),
		patternDaysUntil(),
		patternDaysBetween(),

		// Specific dates
		patternMonthDayYear(),
		patternMonthDay(),
		patternNamedHoliday(),

		// Time parsing
		patternTimeAMPM(),
		patternTime24(),
	})
}

// ════════════════════════════════════════════════════════════════
// TIMEZONE PATTERNS
// ════════════════════════════════════════════════════════════════

// patternNowInTimezone: "now in PST", "current time in Tokyo"
func patternNowInTimezone() *Pattern {
	return NewPattern("now_in_timezone").
		Regex(`(?:now|current\s*time)\s+in\s+(\w+(?:/\w+)?)`).
		Keywords("in").
		Priority(100).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}
			tz := strings.ToUpper(matches[1])
			return `nowin("` + tz + `")`, true
		}).
		Build()
}

// patternTimeInTimezone: "3pm PST", "15:00 UTC", "3:30 pm EST"
func patternTimeInTimezone() *Pattern {
	return NewPattern("time_in_timezone").
		Regex(`(\d{1,2})(?::(\d{2}))?\s*(am|pm|a\.m\.|p\.m\.)?\s+(\w{2,5})`).
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			hour, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			minute := 0.0
			if matches[2] != "" {
				minute, _ = ExtractNumber(matches[2])
			}

			// Handle AM/PM
			ampm := strings.ToLower(strings.ReplaceAll(matches[3], ".", ""))
			h := int(hour)
			if ampm == "pm" && h < 12 {
				h += 12
			} else if ampm == "am" && h == 12 {
				h = 0
			}

			tz := strings.ToUpper(matches[4])

			return `intimezone(time(` + FormatInt(h) + `, ` + FormatInt(int(minute)) + `), "` + tz + `")`, true
		}).
		Build()
}

// patternWhatTimeIn: "what time is it in Tokyo", "time in London"
func patternWhatTimeIn() *Pattern {
	return NewPattern("what_time_in").
		Regex(`(?:what\s+time\s+(?:is\s+it\s+)?in|time\s+in)\s+(\w+(?:/\w+)?)`).
		Keywords("in").
		Priority(95).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}
			location := matches[1]
			// Map city names to timezone
			tz := cityToTimezone(location)
			return `nowin("` + tz + `")`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// DATE ARITHMETIC PATTERNS
// ════════════════════════════════════════════════════════════════

// patternTodayPlusDays: "today + 30 days", "today - 1 week"
func patternTodayPlusDays() *Pattern {
	return NewPattern("today_plus_days").
		Regex(`(today|tomorrow|yesterday)\s*([+-])\s*(\d+)\s*(` + rxTimeUnitNoCap + `)`).
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			base := matches[1]
			op := matches[2]
			num, ok := ExtractNumber(matches[3])
			if !ok {
				return "", false
			}
			unit, ok := ExtractTimeUnit(matches[4])
			if !ok {
				return "", false
			}

			// Adjust for operation
			if op == "-" {
				num = -num
			}

			// Map base to function
			baseFunc := "today()"
			if base == "tomorrow" {
				baseFunc = "adddays(today(), 1)"
			} else if base == "yesterday" {
				baseFunc = "adddays(today(), -1)"
			}

			// Map unit to function
			fn := unitToAddFunction(unit)

			return fn + "(" + baseFunc + ", " + FormatFloat(num) + ")", true
		}).
		Build()
}

// patternNowPlusTime: "now + 2 hours", "now - 30 minutes"
func patternNowPlusTime() *Pattern {
	return NewPattern("now_plus_time").
		Regex(`now\s*([+-])\s*(\d+)\s*(` + rxTimeUnitNoCap + `)`).
		Keywords("now").
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			op := matches[1]
			num, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}
			unit, ok := ExtractTimeUnit(matches[3])
			if !ok {
				return "", false
			}

			if op == "-" {
				num = -num
			}

			fn := unitToAddFunction(unit)

			return fn + "(now(), " + FormatFloat(num) + ")", true
		}).
		Build()
}

// patternDateFromNow: "30 days from now", "2 weeks from today"
func patternDateFromNow() *Pattern {
	return NewPattern("date_from_now").
		Regex(`(\d+)\s*(` + rxTimeUnitNoCap + `)\s+from\s+(now|today)`).
		Keywords("from").
		Priority(75).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			num, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}
			unit, ok := ExtractTimeUnit(matches[2])
			if !ok {
				return "", false
			}

			base := "today()"
			if matches[3] == "now" {
				base = "now()"
			}

			fn := unitToAddFunction(unit)

			return fn + "(" + base + ", " + FormatFloat(num) + ")", true
		}).
		Build()
}

// patternDateAgo: "30 days ago", "2 weeks ago"
func patternDateAgo() *Pattern {
	return NewPattern("date_ago").
		Regex(`(\d+)\s*(` + rxTimeUnitNoCap + `)\s+ago`).
		Keywords("ago").
		Priority(75).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 3 {
				return "", false
			}

			num, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}
			unit, ok := ExtractTimeUnit(matches[2])
			if !ok {
				return "", false
			}

			fn := unitToAddFunction(unit)

			return fn + "(today(), " + FormatFloat(-num) + ")", true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// RELATIVE WEEKDAY PATTERNS
// ════════════════════════════════════════════════════════════════

// patternNextWeekday: "next monday", "next friday"
func patternNextWeekday() *Pattern {
	return NewPattern("next_weekday").
		Regex(`next\s+(` + rxWeekdayNoCap + `)`).
		Keywords("next").
		Priority(70).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}

			day, ok := ExtractWeekday(matches[1])
			if !ok {
				return "", false
			}

			return `nextweekday(today(), ` + FormatInt(day) + `)`, true
		}).
		Build()
}

// patternLastWeekday: "last monday", "last friday"
func patternLastWeekday() *Pattern {
	return NewPattern("last_weekday").
		Regex(`last\s+(` + rxWeekdayNoCap + `)`).
		Keywords("last").
		Priority(70).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}

			day, ok := ExtractWeekday(matches[1])
			if !ok {
				return "", false
			}

			return `lastweekday(today(), ` + FormatInt(day) + `)`, true
		}).
		Build()
}

// patternThisWeekday: "this monday", "this friday"
func patternThisWeekday() *Pattern {
	return NewPattern("this_weekday").
		Regex(`this\s+(` + rxWeekdayNoCap + `)`).
		Keywords("this").
		Priority(70).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}

			day, ok := ExtractWeekday(matches[1])
			if !ok {
				return "", false
			}

			// "this weekday" means the occurrence in the current week
			return `thisweekday(today(), ` + FormatInt(day) + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// DATE RANGE PATTERNS
// ════════════════════════════════════════════════════════════════

// patternDaysSince: "days since jan 1 2025", "days since christmas"
func patternDaysSince() *Pattern {
	return NewPattern("days_since").
		Regex(`days\s+since\s+(.+)`).
		Keywords("days", "since").
		Priority(65).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}

			dateExpr := parseNaturalDate(matches[1])
			if dateExpr == "" {
				return "", false
			}

			return `daysbetween(` + dateExpr + `, today())`, true
		}).
		Build()
}

// patternDaysUntil: "days until christmas", "days until jan 1 2026"
func patternDaysUntil() *Pattern {
	return NewPattern("days_until").
		Regex(`days\s+(?:until|til|till|to)\s+(.+)`).
		Keywords("days").
		Priority(65).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}

			dateExpr := parseNaturalDate(matches[1])
			if dateExpr == "" {
				return "", false
			}

			return `daysbetween(today(), ` + dateExpr + `)`, true
		}).
		Build()
}

// patternDaysBetween: "days between jan 1 and dec 31"
func patternDaysBetween() *Pattern {
	return NewPattern("days_between").
		Regex(`days\s+between\s+(.+?)\s+and\s+(.+)`).
		Keywords("days", "between", "and").
		Priority(60).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 3 {
				return "", false
			}

			date1 := parseNaturalDate(matches[1])
			date2 := parseNaturalDate(matches[2])

			if date1 == "" || date2 == "" {
				return "", false
			}

			return `daysbetween(` + date1 + `, ` + date2 + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// SPECIFIC DATE PATTERNS
// ════════════════════════════════════════════════════════════════

// patternMonthDayYear: "jan 15 2025", "january 15, 2025", "15 jan 2025"
func patternMonthDayYear() *Pattern {
	return NewPattern("month_day_year").
		Regex(`(` + rxMonthNoCap + `)\s+(\d{1,2})(?:st|nd|rd|th)?,?\s*(\d{4})|(\d{1,2})(?:st|nd|rd|th)?\s+(` + rxMonthNoCap + `)\s+(\d{4})`).
		Priority(55).
		Handler(func(input string, matches []string) (string, bool) {
			var month, day, year int
			var ok bool

			if matches[1] != "" {
				// "jan 15 2025" format
				month, ok = ExtractMonth(matches[1])
				if !ok {
					return "", false
				}
				day, ok = ExtractOrdinal(matches[2])
				if !ok {
					return "", false
				}
				y, ok := ExtractNumber(matches[3])
				if !ok {
					return "", false
				}
				year = int(y)
			} else {
				// "15 jan 2025" format
				day, ok = ExtractOrdinal(matches[4])
				if !ok {
					return "", false
				}
				month, ok = ExtractMonth(matches[5])
				if !ok {
					return "", false
				}
				y, ok := ExtractNumber(matches[6])
				if !ok {
					return "", false
				}
				year = int(y)
			}

			return `date(` + FormatInt(year) + `, ` + FormatInt(month) + `, ` + FormatInt(day) + `)`, true
		}).
		Build()
}

// patternMonthDay: "jan 15", "january 15th" (assumes current year)
func patternMonthDay() *Pattern {
	return NewPattern("month_day").
		Regex(`(` + rxMonthNoCap + `)\s+(\d{1,2})(?:st|nd|rd|th)?(?:\s|$)`).
		Priority(50).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 3 {
				return "", false
			}

			month, ok := ExtractMonth(matches[1])
			if !ok {
				return "", false
			}
			day, ok := ExtractOrdinal(matches[2])
			if !ok {
				return "", false
			}

			// Use year() to get current year
			return `date(year(today()), ` + FormatInt(month) + `, ` + FormatInt(day) + `)`, true
		}).
		Build()
}

// patternNamedHoliday: "christmas", "new years", "thanksgiving"
func patternNamedHoliday() *Pattern {
	return NewPattern("named_holiday").
		Regex(`\b(christmas|xmas|new\s*years?|thanksgiving|halloween|valentines?|easter|july\s*4th|independence\s*day|labor\s*day|memorial\s*day)\b`).
		Priority(45).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 2 {
				return "", false
			}

			return holidayToDate(matches[1]), true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// TIME PATTERNS
// ════════════════════════════════════════════════════════════════

// patternTimeAMPM: "3pm", "3:30pm", "3:30:00pm"
func patternTimeAMPM() *Pattern {
	return NewPattern("time_ampm").
		Regex(`^(\d{1,2})(?::(\d{2}))?(?::(\d{2}))?\s*(am|pm|a\.m\.|p\.m\.)$`).
		Priority(40).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			hour, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			minute := 0.0
			if matches[2] != "" {
				minute, _ = ExtractNumber(matches[2])
			}

			second := 0.0
			if matches[3] != "" {
				second, _ = ExtractNumber(matches[3])
			}

			ampm := strings.ToLower(strings.ReplaceAll(matches[4], ".", ""))
			h := int(hour)
			if ampm == "pm" && h < 12 {
				h += 12
			} else if ampm == "am" && h == 12 {
				h = 0
			}

			if second > 0 {
				return `time(` + FormatInt(h) + `, ` + FormatInt(int(minute)) + `, ` + FormatInt(int(second)) + `)`, true
			}
			return `time(` + FormatInt(h) + `, ` + FormatInt(int(minute)) + `)`, true
		}).
		Build()
}

// patternTime24: "15:00", "15:30:00"
func patternTime24() *Pattern {
	return NewPattern("time_24").
		Regex(`^(\d{1,2}):(\d{2})(?::(\d{2}))?$`).
		Priority(35).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 3 {
				return "", false
			}

			hour, ok := ExtractNumber(matches[1])
			if !ok || hour > 23 {
				return "", false
			}

			minute, ok := ExtractNumber(matches[2])
			if !ok || minute > 59 {
				return "", false
			}

			second := 0.0
			if len(matches) > 3 && matches[3] != "" {
				second, _ = ExtractNumber(matches[3])
			}

			if second > 0 {
				return `time(` + FormatInt(int(hour)) + `, ` + FormatInt(int(minute)) + `, ` + FormatInt(int(second)) + `)`, true
			}
			return `time(` + FormatInt(int(hour)) + `, ` + FormatInt(int(minute)) + `)`, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// HELPER PATTERNS (non-capturing for use in other patterns)
// ════════════════════════════════════════════════════════════════

const (
	rxTimeUnitNoCap = `(?:seconds?|secs?|minutes?|mins?|hours?|hrs?|days?|weeks?|wks?|months?|mos?|years?|yrs?)`
	rxWeekdayNoCap  = `(?:sun(?:day)?|mon(?:day)?|tue(?:s(?:day)?)?|wed(?:nesday)?|thu(?:r(?:s(?:day)?)?)?|fri(?:day)?|sat(?:urday)?)`
	rxMonthNoCap    = `(?:jan(?:uary)?|feb(?:ruary)?|mar(?:ch)?|apr(?:il)?|may|june?|july?|aug(?:ust)?|sep(?:t(?:ember)?)?|oct(?:ober)?|nov(?:ember)?|dec(?:ember)?)`
)

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// unitToAddFunction maps time units to add functions.
func unitToAddFunction(unit string) string {
	switch unit {
	case "second":
		return "addseconds"
	case "minute":
		return "addminutes"
	case "hour":
		return "addhours"
	case "day":
		return "adddays"
	case "week":
		return "addweeks"
	case "month":
		return "addmonths"
	case "year":
		return "addyears"
	default:
		return "adddays"
	}
}

// cityToTimezone maps city names to timezone identifiers.
func cityToTimezone(city string) string {
	city = strings.ToLower(city)

	cities := map[string]string{
		"tokyo":         "JST",
		"london":        "GMT",
		"paris":         "CET",
		"berlin":        "CET",
		"moscow":        "MSK",
		"beijing":       "CST",
		"shanghai":      "CST",
		"hong kong":     "HKT",
		"hongkong":      "HKT",
		"singapore":     "SGT",
		"sydney":        "AEST",
		"melbourne":     "AEST",
		"dubai":         "GST",
		"mumbai":        "IST",
		"delhi":         "IST",
		"new york":      "EST",
		"newyork":       "EST",
		"los angeles":   "PST",
		"losangeles":    "PST",
		"chicago":       "CST",
		"denver":        "MST",
		"seattle":       "PST",
		"san francisco": "PST",
		"sanfrancisco":  "PST",
		"toronto":       "EST",
		"vancouver":     "PST",
		"amsterdam":     "CET",
		"rome":          "CET",
		"madrid":        "CET",
		"istanbul":      "TRT",
		"cairo":         "EET",
		"johannesburg":  "SAST",
		"lagos":         "WAT",
		"nairobi":       "EAT",
		"seoul":         "KST",
		"bangkok":       "ICT",
		"jakarta":       "WIB",
		"manila":        "PHT",
		"auckland":      "NZST",
		"hawaii":        "HST",
		"alaska":        "AKST",
	}

	if tz, ok := cities[city]; ok {
		return tz
	}

	// Return as-is (might be a timezone code already)
	return strings.ToUpper(city)
}

// holidayToDate returns a date expression for a named holiday.
func holidayToDate(holiday string) string {
	holiday = strings.ToLower(holiday)
	holiday = strings.ReplaceAll(holiday, " ", "")

	// These return dates relative to current year
	switch holiday {
	case "christmas", "xmas":
		return "date(year(today()), 12, 25)"
	case "newyear", "newyears":
		return "date(year(today()) + 1, 1, 1)"
	case "halloween":
		return "date(year(today()), 10, 31)"
	case "valentine", "valentines":
		return "date(year(today()), 2, 14)"
	case "july4th", "independenceday":
		return "date(year(today()), 7, 4)"
	default:
		// For holidays that need calculation (Easter, Thanksgiving, etc.)
		// return a placeholder for now
		return "today()"
	}
}

// parseNaturalDate parses a natural language date string.
func parseNaturalDate(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))

	// Try relative words first
	if s == "today" {
		return "today()"
	}
	if s == "tomorrow" {
		return "adddays(today(), 1)"
	}
	if s == "yesterday" {
		return "adddays(today(), -1)"
	}
	if s == "now" {
		return "now()"
	}

	// Try holiday names
	if isHoliday(s) {
		return holidayToDate(s)
	}

	// Try "month day year" pattern
	mdyRegex := regexp.MustCompile(`(` + rxMonthNoCap + `)\s+(\d{1,2})(?:st|nd|rd|th)?,?\s*(\d{4})?`)
	if matches := mdyRegex.FindStringSubmatch(s); matches != nil {
		month, ok := ExtractMonth(matches[1])
		if !ok {
			return ""
		}
		day, ok := ExtractOrdinal(matches[2])
		if !ok {
			return ""
		}

		if matches[3] != "" {
			year, _ := ExtractNumber(matches[3])
			return `date(` + FormatInt(int(year)) + `, ` + FormatInt(month) + `, ` + FormatInt(day) + `)`
		}
		return `date(year(today()), ` + FormatInt(month) + `, ` + FormatInt(day) + `)`
	}

	// Try "day month year" pattern
	dmyRegex := regexp.MustCompile(`(\d{1,2})(?:st|nd|rd|th)?\s+(` + rxMonthNoCap + `)\s*(\d{4})?`)
	if matches := dmyRegex.FindStringSubmatch(s); matches != nil {
		day, ok := ExtractOrdinal(matches[1])
		if !ok {
			return ""
		}
		month, ok := ExtractMonth(matches[2])
		if !ok {
			return ""
		}

		if matches[3] != "" {
			year, _ := ExtractNumber(matches[3])
			return `date(` + FormatInt(int(year)) + `, ` + FormatInt(month) + `, ` + FormatInt(day) + `)`
		}
		return `date(year(today()), ` + FormatInt(month) + `, ` + FormatInt(day) + `)`
	}

	return ""
}

// isHoliday checks if a string is a known holiday name.
func isHoliday(s string) bool {
	holidays := []string{
		"christmas", "xmas", "newyear", "newyears", "new year", "new years",
		"halloween", "valentine", "valentines", "thanksgiving",
		"easter", "july4th", "july 4th", "independenceday", "independence day",
		"laborday", "labor day", "memorialday", "memorial day",
	}

	s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	for _, h := range holidays {
		if strings.ReplaceAll(h, " ", "") == s {
			return true
		}
	}
	return false
}
