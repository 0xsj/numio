// internal/date/date.go

package date

import (
	"strings"
	"time"
)

// ════════════════════════════════════════════════════════════════
// DATE TYPE
// ════════════════════════════════════════════════════════════════

// Date represents a date/time value.
type Date struct {
	time.Time
}

// New creates a Date from year, month, day.
func New(year, month, day int) Date {
	return Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)}
}

// NewDateTime creates a Date from year, month, day, hour, minute, second.
func NewDateTime(year, month, day, hour, minute, second int) Date {
	return Date{time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)}
}

// FromTime creates a Date from time.Time.
func FromTime(t time.Time) Date {
	return Date{t}
}

// Today returns today's date at midnight.
func Today() Date {
	now := time.Now()
	return Date{time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)}
}

// Now returns the current date and time.
func Now() Date {
	return Date{time.Now()}
}

// ════════════════════════════════════════════════════════════════
// PARSING
// ════════════════════════════════════════════════════════════════

// Common date formats for parsing
var dateFormats = []string{
	"2006-01-02",
	"2006/01/02",
	"01/02/2006",
	"02/01/2006",
	"Jan 2, 2006",
	"January 2, 2006",
	"2 Jan 2006",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	time.RFC3339,
}

// Parse attempts to parse a date string.
func Parse(s string) (Date, error) {
	s = strings.TrimSpace(s)

	for _, format := range dateFormats {
		if t, err := time.Parse(format, s); err == nil {
			return Date{t}, nil
		}
	}

	return Date{}, ErrInvalidDateFormat
}

// ParseOrZero parses a date or returns zero date.
func ParseOrZero(s string) Date {
	d, err := Parse(s)
	if err != nil {
		return Date{}
	}
	return d
}

// FromUnix creates a Date from Unix timestamp (seconds).
func FromUnix(timestamp int64) Date {
	return Date{time.Unix(timestamp, 0)}
}

// FromUnixMilli creates a Date from Unix timestamp (milliseconds).
func FromUnixMilli(timestamp int64) Date {
	return Date{time.UnixMilli(timestamp)}
}

// ════════════════════════════════════════════════════════════════
// FORMATTING
// ════════════════════════════════════════════════════════════════

// String returns the date as a string.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	// If time is midnight, return date only
	if d.Hour() == 0 && d.Minute() == 0 && d.Second() == 0 {
		return d.Format("2006-01-02")
	}
	return d.Format("2006-01-02 15:04:05")
}

// DateOnly returns date in YYYY-MM-DD format.
func (d Date) DateOnly() string {
	return d.Format("2006-01-02")
}

// TimeOnly returns time in HH:MM:SS format.
func (d Date) TimeOnly() string {
	return d.Format("15:04:05")
}

// ISO returns the date in ISO 8601 format.
func (d Date) ISO() string {
	return d.Format(time.RFC3339)
}

// Unix returns the Unix timestamp (seconds).
func (d Date) Unix() int64 {
	return d.Time.Unix()
}

// UnixMilli returns the Unix timestamp (milliseconds).
func (d Date) UnixMilli() int64 {
	return d.Time.UnixMilli()
}

// FormatCustom formats the date with a custom layout.
func (d Date) FormatCustom(layout string) string {
	return d.Format(layout)
}

// ════════════════════════════════════════════════════════════════
// COMPARISON
// ════════════════════════════════════════════════════════════════

// Equal checks if two dates are equal.
func (d Date) Equal(other Date) bool {
	return d.Time.Equal(other.Time)
}

// Before checks if this date is before another.
func (d Date) Before(other Date) bool {
	return d.Time.Before(other.Time)
}

// After checks if this date is after another.
func (d Date) After(other Date) bool {
	return d.Time.After(other.Time)
}

// Between checks if this date is between two dates (inclusive).
func (d Date) Between(start, end Date) bool {
	return !d.Before(start) && !d.After(end)
}

// ════════════════════════════════════════════════════════════════
// ERRORS
// ════════════════════════════════════════════════════════════════

// DateError represents a date-related error.
type DateError string

func (e DateError) Error() string {
	return string(e)
}

const (
	ErrInvalidDateFormat DateError = "invalid date format"
	ErrInvalidDate       DateError = "invalid date"
	ErrInvalidTime       DateError = "invalid time"
	ErrInvalidTimezone   DateError = "invalid timezone"
)
