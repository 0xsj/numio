// internal/date/components.go

package date

import (
	"time"
)

// ════════════════════════════════════════════════════════════════
// DATE COMPONENTS
// ════════════════════════════════════════════════════════════════

// Year returns the year.
func (d Date) Year() int {
	return d.Time.Year()
}

// Month returns the month (1-12).
func (d Date) MonthInt() int {
	return int(d.Time.Month())
}

// Day returns the day of the month (1-31).
func (d Date) Day() int {
	return d.Time.Day()
}

// Hour returns the hour (0-23).
func (d Date) Hour() int {
	return d.Time.Hour()
}

// Minute returns the minute (0-59).
func (d Date) Minute() int {
	return d.Time.Minute()
}

// Second returns the second (0-59).
func (d Date) Second() int {
	return d.Time.Second()
}

// Millisecond returns the millisecond (0-999).
func (d Date) Millisecond() int {
	return d.Time.Nanosecond() / 1e6
}

// ════════════════════════════════════════════════════════════════
// WEEKDAY
// ════════════════════════════════════════════════════════════════

// Weekday returns the day of the week (0=Sunday, 6=Saturday).
func (d Date) WeekdayInt() int {
	return int(d.Time.Weekday())
}

// WeekdayISO returns the ISO day of the week (1=Monday, 7=Sunday).
func (d Date) WeekdayISO() int {
	w := int(d.Time.Weekday())
	if w == 0 {
		return 7
	}
	return w
}

// WeekdayName returns the full weekday name.
func (d Date) WeekdayName() string {
	return d.Time.Weekday().String()
}

// WeekdayShort returns the short weekday name (Mon, Tue, etc).
func (d Date) WeekdayShort() string {
	return d.Format("Mon")
}

// IsWeekend returns true if the date is Saturday or Sunday.
func (d Date) IsWeekend() bool {
	w := d.Time.Weekday()
	return w == time.Saturday || w == time.Sunday
}

// IsWeekday returns true if the date is Monday-Friday.
func (d Date) IsWeekday() bool {
	return !d.IsWeekend()
}

// ════════════════════════════════════════════════════════════════
// MONTH INFO
// ════════════════════════════════════════════════════════════════

// MonthName returns the full month name.
func (d Date) MonthName() string {
	return d.Time.Month().String()
}

// MonthShort returns the short month name (Jan, Feb, etc).
func (d Date) MonthShort() string {
	return d.Format("Jan")
}

// ════════════════════════════════════════════════════════════════
// YEAR POSITION
// ════════════════════════════════════════════════════════════════

// DayOfYear returns the day of the year (1-366).
func (d Date) DayOfYear() int {
	return d.Time.YearDay()
}

// WeekOfYear returns the ISO week number (1-53).
func (d Date) WeekOfYear() int {
	_, week := d.Time.ISOWeek()
	return week
}

// ISOWeek returns the ISO year and week number.
func (d Date) ISOWeek() (year, week int) {
	return d.Time.ISOWeek()
}

// Quarter returns the quarter (1-4).
func (d Date) Quarter() int {
	return (int(d.Time.Month())-1)/3 + 1
}

// DaysInMonth returns the number of days in the current month.
func (d Date) DaysInMonth() int {
	return DaysInMonthOf(d.Year(), d.MonthInt())
}

// DaysInYear returns the number of days in the current year.
func (d Date) DaysInYear() int {
	if IsLeapYear(d.Year()) {
		return 366
	}
	return 365
}

// DaysRemainingInMonth returns days remaining in the month.
func (d Date) DaysRemainingInMonth() int {
	return d.DaysInMonth() - d.Day()
}

// DaysRemainingInYear returns days remaining in the year.
func (d Date) DaysRemainingInYear() int {
	return d.DaysInYear() - d.DayOfYear()
}

// ════════════════════════════════════════════════════════════════
// TIMEZONE
// ════════════════════════════════════════════════════════════════

// TimezoneInfo holds timezone data for fallback.
type TimezoneInfo struct {
	IANAName string
	Offset   int // Offset in seconds from UTC
}

// Common timezone mappings with offsets as fallback
var timezoneData = map[string]TimezoneInfo{
	// North America
	"EST":  {"America/New_York", -5 * 3600},
	"EDT":  {"America/New_York", -4 * 3600},
	"CST":  {"America/Chicago", -6 * 3600},
	"CDT":  {"America/Chicago", -5 * 3600},
	"MST":  {"America/Denver", -7 * 3600},
	"MDT":  {"America/Denver", -6 * 3600},
	"PST":  {"America/Los_Angeles", -8 * 3600},
	"PDT":  {"America/Los_Angeles", -7 * 3600},
	"AKST": {"America/Anchorage", -9 * 3600},
	"AKDT": {"America/Anchorage", -8 * 3600},
	"HST":  {"Pacific/Honolulu", -10 * 3600},

	// Europe
	"GMT":  {"Europe/London", 0},
	"UTC":  {"UTC", 0},
	"BST":  {"Europe/London", 1 * 3600},
	"CET":  {"Europe/Paris", 1 * 3600},
	"CEST": {"Europe/Paris", 2 * 3600},
	"EET":  {"Europe/Helsinki", 2 * 3600},
	"EEST": {"Europe/Helsinki", 3 * 3600},
	"WET":  {"Europe/Lisbon", 0},
	"WEST": {"Europe/Lisbon", 1 * 3600},

	// Asia
	"IST":       {"Asia/Kolkata", 5*3600 + 30*60}, // +5:30
	"JST":       {"Asia/Tokyo", 9 * 3600},
	"KST":       {"Asia/Seoul", 9 * 3600},
	"CST_CHINA": {"Asia/Shanghai", 8 * 3600},
	"HKT":       {"Asia/Hong_Kong", 8 * 3600},
	"SGT":       {"Asia/Singapore", 8 * 3600},
	"PHT":       {"Asia/Manila", 8 * 3600},
	"ICT":       {"Asia/Bangkok", 7 * 3600},
	"WIB":       {"Asia/Jakarta", 7 * 3600},

	// Middle East
	"TRT":  {"Europe/Istanbul", 3 * 3600},   // Turkey
	"AST":  {"Asia/Riyadh", 3 * 3600},       // Arabia
	"IRST": {"Asia/Tehran", 3*3600 + 30*60}, // Iran +3:30

	// Australia
	"AEST": {"Australia/Sydney", 10 * 3600},
	"AEDT": {"Australia/Sydney", 11 * 3600},
	"ACST": {"Australia/Adelaide", 9*3600 + 30*60},
	"ACDT": {"Australia/Adelaide", 10*3600 + 30*60},
	"AWST": {"Australia/Perth", 8 * 3600},

	// Others
	"NZST": {"Pacific/Auckland", 12 * 3600},
	"NZDT": {"Pacific/Auckland", 13 * 3600},
	"BRT":  {"America/Sao_Paulo", -3 * 3600},
	"ART":  {"America/Argentina/Buenos_Aires", -3 * 3600},
}

// GetTimezone returns a timezone location from abbreviation or name.
// Falls back to fixed offset if IANA database unavailable.
func GetTimezone(name string) (*time.Location, error) {
	// Handle UTC specially - always available
	if name == "UTC" || name == "utc" {
		return time.UTC, nil
	}

	// Check our timezone data first
	if info, ok := timezoneData[name]; ok {
		// Try IANA name first
		if loc, err := time.LoadLocation(info.IANAName); err == nil {
			return loc, nil
		}
		// Fall back to fixed offset
		return time.FixedZone(name, info.Offset), nil
	}

	// Try as IANA name directly
	if loc, err := time.LoadLocation(name); err == nil {
		return loc, nil
	}

	// Try parsing as offset string like "+05:30" or "-08:00"
	if loc, err := parseOffsetString(name); err == nil {
		return loc, nil
	}

	return nil, ErrInvalidTimezone
}

// parseOffsetString parses timezone offset strings like "+05:30", "-08:00", "+8", "-5"
func parseOffsetString(s string) (*time.Location, error) {
	if len(s) == 0 {
		return nil, ErrInvalidTimezone
	}

	// Check for sign
	sign := 1
	offset := s
	if s[0] == '+' {
		offset = s[1:]
	} else if s[0] == '-' {
		sign = -1
		offset = s[1:]
	}

	var hours, minutes int

	// Try parsing "HH:MM" format
	if idx := indexByte(offset, ':'); idx >= 0 {
		h, errH := parseSimpleInt(offset[:idx])
		m, errM := parseSimpleInt(offset[idx+1:])
		if errH != nil || errM != nil {
			return nil, ErrInvalidTimezone
		}
		hours, minutes = h, m
	} else {
		// Try parsing as just hours
		h, err := parseSimpleInt(offset)
		if err != nil {
			return nil, ErrInvalidTimezone
		}
		hours = h
	}

	// Validate
	if hours < 0 || hours > 14 || minutes < 0 || minutes > 59 {
		return nil, ErrInvalidTimezone
	}

	totalSeconds := sign * (hours*3600 + minutes*60)
	name := formatOffset(sign, hours, minutes)

	return time.FixedZone(name, totalSeconds), nil
}

// formatOffset creates a timezone name from offset components.
func formatOffset(sign, hours, minutes int) string {
	s := "+"
	if sign < 0 {
		s = "-"
	}
	if minutes == 0 {
		return s + intToStr(hours)
	}
	h := intToStr(hours)
	if hours < 10 {
		h = "0" + h
	}
	m := intToStr(minutes)
	if minutes < 10 {
		m = "0" + m
	}
	return s + h + ":" + m
}

// indexByte finds the index of a byte in a string.
func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// parseSimpleInt parses a simple integer from a string.
func parseSimpleInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, ErrInvalidTimezone
	}
	result := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, ErrInvalidTimezone
		}
		result = result*10 + int(c-'0')
	}
	return result, nil
}

// InTimezone returns the date converted to a timezone.
func (d Date) InTimezone(tz string) (Date, error) {
	loc, err := GetTimezone(tz)
	if err != nil {
		return d, err
	}
	return Date{d.Time.In(loc)}, nil
}

// InTimezoneOrSame returns the date in timezone, or same date if error.
func (d Date) InTimezoneOrSame(tz string) Date {
	result, err := d.InTimezone(tz)
	if err != nil {
		return d
	}
	return result
}

// ToUTC converts the date to UTC.
func (d Date) ToUTC() Date {
	return Date{d.Time.UTC()}
}

// ToLocal converts the date to local timezone.
func (d Date) ToLocal() Date {
	return Date{d.Time.Local()}
}

// TimezoneOffset returns the offset from UTC in hours.
func (d Date) TimezoneOffset() float64 {
	_, offset := d.Time.Zone()
	return float64(offset) / 3600
}

// TimezoneName returns the timezone abbreviation.
func (d Date) TimezoneName() string {
	name, _ := d.Time.Zone()
	return name
}

// ListTimezones returns all supported timezone abbreviations.
func ListTimezones() []string {
	zones := make([]string, 0, len(timezoneData))
	for name := range timezoneData {
		zones = append(zones, name)
	}
	return zones
}

// ════════════════════════════════════════════════════════════════
// TIME OF DAY
// ════════════════════════════════════════════════════════════════

// IsMorning returns true if time is between 5:00 and 11:59.
func (d Date) IsMorning() bool {
	h := d.Hour()
	return h >= 5 && h < 12
}

// IsAfternoon returns true if time is between 12:00 and 16:59.
func (d Date) IsAfternoon() bool {
	h := d.Hour()
	return h >= 12 && h < 17
}

// IsEvening returns true if time is between 17:00 and 20:59.
func (d Date) IsEvening() bool {
	h := d.Hour()
	return h >= 17 && h < 21
}

// IsNight returns true if time is between 21:00 and 4:59.
func (d Date) IsNight() bool {
	h := d.Hour()
	return h >= 21 || h < 5
}

// Period returns the period of day as a string.
func (d Date) Period() string {
	if d.IsMorning() {
		return "morning"
	}
	if d.IsAfternoon() {
		return "afternoon"
	}
	if d.IsEvening() {
		return "evening"
	}
	return "night"
}

// Hour12 returns the hour in 12-hour format (1-12).
func (d Date) Hour12() int {
	h := d.Hour() % 12
	if h == 0 {
		h = 12
	}
	return h
}

// IsPM returns true if time is PM.
func (d Date) IsPM() bool {
	return d.Hour() >= 12
}

// IsAM returns true if time is AM.
func (d Date) IsAM() bool {
	return d.Hour() < 12
}

// AMPM returns "AM" or "PM".
func (d Date) AMPM() string {
	if d.IsPM() {
		return "PM"
	}
	return "AM"
}

// ════════════════════════════════════════════════════════════════
// TIME STRING PARSING (e.g., "1230pm", "3:45 PM")
// ════════════════════════════════════════════════════════════════

// TimeFormats for parsing time strings
var timeFormats = []string{
	"3:04pm",
	"3:04 pm",
	"3:04PM",
	"3:04 PM",
	"15:04",
	"15:04:05",
	"3pm",
	"3 pm",
	"3PM",
	"3 PM",
}

// ParseTime parses a time string and returns hour, minute, second.
func ParseTime(s string) (hour, minute, second int, err error) {
	// Handle formats like "1230pm" -> "12:30pm"
	s = normalizeTimeString(s)

	for _, format := range timeFormats {
		t, parseErr := time.Parse(format, s)
		if parseErr == nil {
			return t.Hour(), t.Minute(), t.Second(), nil
		}
	}

	return 0, 0, 0, ErrInvalidTime
}

// normalizeTimeString converts "1230pm" to "12:30pm" etc.
func normalizeTimeString(s string) string {
	// Remove spaces
	result := ""
	for _, c := range s {
		if c != ' ' {
			result += string(c)
		}
	}

	// Check for patterns like "1230pm" (4 digits + am/pm)
	if len(result) >= 6 {
		lower := toLower(result)
		if (endsWith(lower, "am") || endsWith(lower, "pm")) && isDigits(result[:4]) {
			// "1230pm" -> "12:30pm"
			return result[:2] + ":" + result[2:4] + result[4:]
		}
	}

	// Check for patterns like "130pm" (3 digits + am/pm)
	if len(result) >= 5 {
		lower := toLower(result)
		if (endsWith(lower, "am") || endsWith(lower, "pm")) && isDigits(result[:3]) {
			// "130pm" -> "1:30pm"
			return result[:1] + ":" + result[1:3] + result[3:]
		}
	}

	return result
}

// SetTime sets the time components of the date.
func (d Date) SetTime(hour, minute, second int) Date {
	return NewDateTime(d.Year(), d.MonthInt(), d.Day(), hour, minute, second)
}

// WithTimeString sets the time from a string like "1230pm".
func (d Date) WithTimeString(timeStr string) (Date, error) {
	hour, minute, second, err := ParseTime(timeStr)
	if err != nil {
		return d, err
	}
	return d.SetTime(hour, minute, second), nil
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func toLower(s string) string {
	result := ""
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			result += string(c + 32)
		} else {
			result += string(c)
		}
	}
	return result
}

func endsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

func isDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
