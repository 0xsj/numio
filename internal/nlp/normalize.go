// internal/nlp/normalize.go

package nlp

import (
	"strings"
	"unicode"
)

// ════════════════════════════════════════════════════════════════
// TEXT NORMALIZATION
// ════════════════════════════════════════════════════════════════

// Normalize prepares text for pattern matching.
// - Lowercases
// - Normalizes whitespace
// - Removes unnecessary punctuation
func Normalize(input string) string {
	// Lowercase
	input = strings.ToLower(input)

	// Normalize whitespace
	input = normalizeWhitespace(input)

	// Remove question marks (for "what time is it?")
	input = strings.ReplaceAll(input, "?", "")

	return strings.TrimSpace(input)
}

// normalizeWhitespace collapses multiple spaces into one.
func normalizeWhitespace(s string) string {
	var result strings.Builder
	inSpace := false

	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inSpace {
				result.WriteRune(' ')
				inSpace = true
			}
		} else {
			result.WriteRune(r)
			inSpace = false
		}
	}

	return result.String()
}

// ════════════════════════════════════════════════════════════════
// TOKENIZATION
// ════════════════════════════════════════════════════════════════

// Token represents a word or symbol in the input.
type Token struct {
	Text  string
	Start int
	End   int
}

// Tokenize splits input into tokens.
func Tokenize(input string) []Token {
	var tokens []Token
	var current strings.Builder
	start := 0

	for i, r := range input {
		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				tokens = append(tokens, Token{
					Text:  current.String(),
					Start: start,
					End:   i,
				})
				current.Reset()
			}
			start = i + 1
		} else if isPunctuation(r) {
			// Emit current token if any
			if current.Len() > 0 {
				tokens = append(tokens, Token{
					Text:  current.String(),
					Start: start,
					End:   i,
				})
				current.Reset()
			}
			// Emit punctuation as its own token
			tokens = append(tokens, Token{
				Text:  string(r),
				Start: i,
				End:   i + 1,
			})
			start = i + 1
		} else {
			if current.Len() == 0 {
				start = i
			}
			current.WriteRune(r)
		}
	}

	// Emit final token
	if current.Len() > 0 {
		tokens = append(tokens, Token{
			Text:  current.String(),
			Start: start,
			End:   len(input),
		})
	}

	return tokens
}

// isPunctuation checks if a rune is punctuation (but not currency symbols).
func isPunctuation(r rune) bool {
	// Keep currency symbols as part of tokens
	if r == '$' || r == '€' || r == '£' || r == '¥' || r == '₿' {
		return false
	}
	// Keep % with numbers
	if r == '%' {
		return false
	}
	return unicode.IsPunct(r)
}

// ════════════════════════════════════════════════════════════════
// NUMBER EXTRACTION
// ════════════════════════════════════════════════════════════════

// ExtractNumber tries to parse a number from a string.
// Handles formats like: 100, 1000, 1,000, 1.5, $100, 100k, 250K
func ExtractNumber(s string) (float64, bool) {
	// Remove currency symbols
	s = strings.TrimLeft(s, "$€£¥₿")

	// Remove commas
	s = strings.ReplaceAll(s, ",", "")

	// Handle k/K suffix (thousands)
	multiplier := 1.0
	if strings.HasSuffix(s, "k") || strings.HasSuffix(s, "K") {
		multiplier = 1000
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "m") || strings.HasSuffix(s, "M") {
		multiplier = 1000000
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "b") || strings.HasSuffix(s, "B") {
		multiplier = 1000000000
		s = s[:len(s)-1]
	}

	// Parse the number
	num, err := parseFloat(s)
	if err != nil {
		return 0, false
	}

	return num * multiplier, true
}

// parseFloat parses a float without the fmt package.
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, errInvalidNumber
	}

	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}

	var result float64
	var decimalPlace float64
	inDecimal := false

	for _, r := range s {
		if r == '.' {
			if inDecimal {
				return 0, errInvalidNumber
			}
			inDecimal = true
			decimalPlace = 0.1
			continue
		}

		if r < '0' || r > '9' {
			return 0, errInvalidNumber
		}

		digit := float64(r - '0')
		if inDecimal {
			result += digit * decimalPlace
			decimalPlace *= 0.1
		} else {
			result = result*10 + digit
		}
	}

	if negative {
		result = -result
	}

	return result, nil
}

var errInvalidNumber = &parseError{"invalid number"}

type parseError struct {
	msg string
}

func (e *parseError) Error() string {
	return e.msg
}

// ════════════════════════════════════════════════════════════════
// WORD LISTS
// ════════════════════════════════════════════════════════════════

// TimeUnits maps time unit words to singular form.
var TimeUnits = map[string]string{
	"second":  "second",
	"seconds": "second",
	"sec":     "second",
	"secs":    "second",
	"minute":  "minute",
	"minutes": "minute",
	"min":     "minute",
	"mins":    "minute",
	"hour":    "hour",
	"hours":   "hour",
	"hr":      "hour",
	"hrs":     "hour",
	"day":     "day",
	"days":    "day",
	"week":    "week",
	"weeks":   "week",
	"wk":      "week",
	"wks":     "week",
	"month":   "month",
	"months":  "month",
	"mo":      "month",
	"mos":     "month",
	"year":    "year",
	"years":   "year",
	"yr":      "year",
	"yrs":     "year",
}

// RelativeWords for date expressions.
var RelativeWords = map[string]int{
	"yesterday": -1,
	"today":     0,
	"tomorrow":  1,
	"now":       0,
}

// Weekdays maps day names to weekday number (0=Sunday).
var Weekdays = map[string]int{
	"sunday":    0,
	"sun":       0,
	"monday":    1,
	"mon":       1,
	"tuesday":   2,
	"tue":       2,
	"tues":      2,
	"wednesday": 3,
	"wed":       3,
	"thursday":  4,
	"thu":       4,
	"thur":      4,
	"thurs":     4,
	"friday":    5,
	"fri":       5,
	"saturday":  6,
	"sat":       6,
}

// Months maps month names to month number (1-12).
var Months = map[string]int{
	"january":   1,
	"jan":       1,
	"february":  2,
	"feb":       2,
	"march":     3,
	"mar":       3,
	"april":     4,
	"apr":       4,
	"may":       5,
	"june":      6,
	"jun":       6,
	"july":      7,
	"jul":       7,
	"august":    8,
	"aug":       8,
	"september": 9,
	"sep":       9,
	"sept":      9,
	"october":   10,
	"oct":       10,
	"november":  11,
	"nov":       11,
	"december":  12,
	"dec":       12,
}

// Ordinals maps ordinal words to numbers.
var Ordinals = map[string]int{
	"first":          1,
	"second":         2,
	"third":          3,
	"fourth":         4,
	"fifth":          5,
	"sixth":          6,
	"seventh":        7,
	"eighth":         8,
	"ninth":          9,
	"tenth":          10,
	"eleventh":       11,
	"twelfth":        12,
	"thirteenth":     13,
	"fourteenth":     14,
	"fifteenth":      15,
	"sixteenth":      16,
	"seventeenth":    17,
	"eighteenth":     18,
	"nineteenth":     19,
	"twentieth":      20,
	"twenty-first":   21,
	"twenty-second":  22,
	"twenty-third":   23,
	"twenty-fourth":  24,
	"twenty-fifth":   25,
	"twenty-sixth":   26,
	"twenty-seventh": 27,
	"twenty-eighth":  28,
	"twenty-ninth":   29,
	"thirtieth":      30,
	"thirty-first":   31,
	"1st":            1,
	"2nd":            2,
	"3rd":            3,
	"4th":            4,
	"5th":            5,
	"6th":            6,
	"7th":            7,
	"8th":            8,
	"9th":            9,
	"10th":           10,
	"11th":           11,
	"12th":           12,
	"13th":           13,
	"14th":           14,
	"15th":           15,
	"16th":           16,
	"17th":           17,
	"18th":           18,
	"19th":           19,
	"20th":           20,
	"21st":           21,
	"22nd":           22,
	"23rd":           23,
	"24th":           24,
	"25th":           25,
	"26th":           26,
	"27th":           27,
	"28th":           28,
	"29th":           29,
	"30th":           30,
	"31st":           31,
}

// NumberWords maps number words to values.
var NumberWords = map[string]int{
	"zero":      0,
	"one":       1,
	"two":       2,
	"three":     3,
	"four":      4,
	"five":      5,
	"six":       6,
	"seven":     7,
	"eight":     8,
	"nine":      9,
	"ten":       10,
	"eleven":    11,
	"twelve":    12,
	"thirteen":  13,
	"fourteen":  14,
	"fifteen":   15,
	"sixteen":   16,
	"seventeen": 17,
	"eighteen":  18,
	"nineteen":  19,
	"twenty":    20,
	"thirty":    30,
	"forty":     40,
	"fifty":     50,
	"sixty":     60,
	"seventy":   70,
	"eighty":    80,
	"ninety":    90,
	"hundred":   100,
	"thousand":  1000,
	"million":   1000000,
	"billion":   1000000000,
	"a":         1,
	"an":        1,
}
