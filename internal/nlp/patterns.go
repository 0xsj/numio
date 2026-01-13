// internal/nlp/patterns.go

package nlp

import (
	"regexp"
	"strings"
)

// ════════════════════════════════════════════════════════════════
// PATTERN TYPES
// ════════════════════════════════════════════════════════════════

// Pattern represents a matching rule that transforms natural language.
type Pattern struct {
	Name     string         // Pattern name for debugging
	Regex    *regexp.Regexp // Regex to match
	Keywords []string       // Required keywords (all must be present)
	Priority int            // Higher priority patterns match first
	Handler  PatternHandler // Function to transform the match
}

// PatternHandler transforms a matched input into a normalized expression.
// Returns the transformed string and true if successful.
type PatternHandler func(input string, matches []string) (string, bool)

// Match represents a successful pattern match.
type Match struct {
	Pattern     *Pattern
	Input       string
	Captures    []string
	Transformed string
}

// ════════════════════════════════════════════════════════════════
// PATTERN REGISTRY
// ════════════════════════════════════════════════════════════════

// PatternRegistry holds all registered patterns.
type PatternRegistry struct {
	patterns []*Pattern
}

// NewRegistry creates a new pattern registry.
func NewRegistry() *PatternRegistry {
	return &PatternRegistry{
		patterns: make([]*Pattern, 0),
	}
}

// Register adds a pattern to the registry.
func (r *PatternRegistry) Register(p *Pattern) {
	r.patterns = append(r.patterns, p)
}

// RegisterAll adds multiple patterns.
func (r *PatternRegistry) RegisterAll(patterns []*Pattern) {
	for _, p := range patterns {
		r.Register(p)
	}
}

// Match tries to match input against all patterns.
// Returns the first successful match (by priority).
func (r *PatternRegistry) Match(input string) *Match {
	normalized := Normalize(input)

	// Sort by priority (higher first) - simple bubble for small list
	sorted := make([]*Pattern, len(r.patterns))
	copy(sorted, r.patterns)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Priority > sorted[i].Priority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	for _, p := range sorted {
		// Check keywords first (fast path)
		if !p.hasKeywords(normalized) {
			continue
		}

		// Try regex match
		if p.Regex != nil {
			matches := p.Regex.FindStringSubmatch(normalized)
			if matches != nil {
				if transformed, ok := p.Handler(normalized, matches); ok {
					return &Match{
						Pattern:     p,
						Input:       input,
						Captures:    matches,
						Transformed: transformed,
					}
				}
			}
		} else {
			// No regex, just keyword match
			if transformed, ok := p.Handler(normalized, nil); ok {
				return &Match{
					Pattern:     p,
					Input:       input,
					Captures:    nil,
					Transformed: transformed,
				}
			}
		}
	}

	return nil
}

// hasKeywords checks if input contains all required keywords.
func (p *Pattern) hasKeywords(input string) bool {
	for _, kw := range p.Keywords {
		if !strings.Contains(input, kw) {
			return false
		}
	}
	return true
}

// ════════════════════════════════════════════════════════════════
// PATTERN BUILDER
// ════════════════════════════════════════════════════════════════

// PatternBuilder provides a fluent API for building patterns.
type PatternBuilder struct {
	pattern *Pattern
}

// NewPattern starts building a new pattern.
func NewPattern(name string) *PatternBuilder {
	return &PatternBuilder{
		pattern: &Pattern{
			Name:     name,
			Priority: 0,
		},
	}
}

// Regex sets the pattern's regex.
func (b *PatternBuilder) Regex(pattern string) *PatternBuilder {
	b.pattern.Regex = regexp.MustCompile(pattern)
	return b
}

// Keywords sets required keywords.
func (b *PatternBuilder) Keywords(keywords ...string) *PatternBuilder {
	b.pattern.Keywords = keywords
	return b
}

// Priority sets the pattern priority.
func (b *PatternBuilder) Priority(p int) *PatternBuilder {
	b.pattern.Priority = p
	return b
}

// Handler sets the pattern handler.
func (b *PatternBuilder) Handler(h PatternHandler) *PatternBuilder {
	b.pattern.Handler = h
	return b
}

// Build returns the constructed pattern.
func (b *PatternBuilder) Build() *Pattern {
	return b.pattern
}

// ════════════════════════════════════════════════════════════════
// COMMON REGEX PATTERNS
// ════════════════════════════════════════════════════════════════

// Common regex building blocks
const (
	// Numbers
	rxNumber  = `(\d+(?:\.\d+)?)`
	rxNumberK = `(\d+(?:\.\d+)?[kmb]?)`
	rxInt     = `(\d+)`
	rxOrdinal = `(\d+(?:st|nd|rd|th)|\w+)`

	// Currency
	rxCurrency = `([$€£¥₿]?\d+(?:,\d{3})*(?:\.\d+)?[kmb]?)`
	rxCurrCode = `([a-zA-Z]{3})`

	// Time
	rxTime12   = `(\d{1,2})(?::(\d{2}))?(?::(\d{2}))?\s*(am|pm|a\.m\.|p\.m\.)?`
	rxTime24   = `(\d{1,2}):(\d{2})(?::(\d{2}))?`
	rxTimezone = `([a-zA-Z]{2,5}|[a-zA-Z]+/[a-zA-Z_]+)`

	// Date components
	rxMonth   = `(jan(?:uary)?|feb(?:ruary)?|mar(?:ch)?|apr(?:il)?|may|june?|july?|aug(?:ust)?|sep(?:t(?:ember)?)?|oct(?:ober)?|nov(?:ember)?|dec(?:ember)?)`
	rxWeekday = `(sun(?:day)?|mon(?:day)?|tue(?:s(?:day)?)?|wed(?:nesday)?|thu(?:r(?:s(?:day)?)?)?|fri(?:day)?|sat(?:urday)?)`
	rxYear    = `((?:19|20)\d{2}|\d{2})`

	// Time units
	rxTimeUnit = `(seconds?|secs?|minutes?|mins?|hours?|hrs?|days?|weeks?|wks?|months?|mos?|years?|yrs?)`

	// Relative
	rxRelative  = `(yesterday|today|tomorrow|now)`
	rxDirection = `(last|next|this|from|ago|before|after|since|until|til|till)`
)

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// JoinOr creates a regex alternation from strings.
func JoinOr(items ...string) string {
	return "(" + strings.Join(items, "|") + ")"
}

// Optional wraps a pattern to make it optional.
func Optional(pattern string) string {
	return "(?:" + pattern + ")?"
}

// Group wraps a pattern in a non-capturing group.
func Group(pattern string) string {
	return "(?:" + pattern + ")"
}

// Capture wraps a pattern in a capturing group.
func Capture(pattern string) string {
	return "(" + pattern + ")"
}

// Word matches a complete word.
func Word(word string) string {
	return `\b` + word + `\b`
}

// AnyOf creates a pattern matching any of the given words.
func AnyOf(words ...string) string {
	escaped := make([]string, len(words))
	for i, w := range words {
		escaped[i] = regexp.QuoteMeta(w)
	}
	return `\b` + JoinOr(escaped...) + `\b`
}

// ════════════════════════════════════════════════════════════════
// EXTRACTION HELPERS
// ════════════════════════════════════════════════════════════════

// ExtractTimeUnit normalizes a time unit string.
func ExtractTimeUnit(s string) (string, bool) {
	s = strings.ToLower(s)
	if unit, ok := TimeUnits[s]; ok {
		return unit, true
	}
	return "", false
}

// ExtractMonth extracts month number from name.
func ExtractMonth(s string) (int, bool) {
	s = strings.ToLower(s)
	if m, ok := Months[s]; ok {
		return m, true
	}
	return 0, false
}

// ExtractWeekday extracts weekday number from name.
func ExtractWeekday(s string) (int, bool) {
	s = strings.ToLower(s)
	if d, ok := Weekdays[s]; ok {
		return d, true
	}
	return 0, false
}

// ExtractOrdinal extracts a number from an ordinal.
func ExtractOrdinal(s string) (int, bool) {
	s = strings.ToLower(s)

	// Check word ordinals
	if n, ok := Ordinals[s]; ok {
		return n, true
	}

	// Try parsing numeric ordinal (1st, 2nd, etc.)
	s = strings.TrimSuffix(s, "st")
	s = strings.TrimSuffix(s, "nd")
	s = strings.TrimSuffix(s, "rd")
	s = strings.TrimSuffix(s, "th")

	if num, ok := ExtractNumber(s); ok {
		return int(num), true
	}

	return 0, false
}

// ExtractNumberWord extracts a number from a word.
func ExtractNumberWord(s string) (int, bool) {
	s = strings.ToLower(s)

	// Check number words
	if n, ok := NumberWords[s]; ok {
		return n, true
	}

	// Try parsing as number
	if num, ok := ExtractNumber(s); ok {
		return int(num), true
	}

	return 0, false
}

// FormatInt converts an int to a string without fmt.
func FormatInt(n int) string {
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

// FormatFloat converts a float to a string without fmt.
func FormatFloat(f float64) string {
	if f == float64(int(f)) {
		return FormatInt(int(f))
	}

	// Handle negative
	if f < 0 {
		return "-" + FormatFloat(-f)
	}

	intPart := int(f)
	fracPart := f - float64(intPart)

	// Get up to 6 decimal places
	fracStr := ""
	for i := 0; i < 6 && fracPart > 0.0000001; i++ {
		fracPart *= 10
		digit := int(fracPart)
		fracStr += string(byte('0' + digit))
		fracPart -= float64(digit)
	}

	// Trim trailing zeros
	fracStr = strings.TrimRight(fracStr, "0")

	if fracStr == "" {
		return FormatInt(intPart)
	}

	return FormatInt(intPart) + "." + fracStr
}
