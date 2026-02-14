// internal/nlp/nlp.go

// Package nlp provides natural language processing for numio expressions.
// It transforms natural language queries into structured expressions that
// the existing parser and evaluator can handle.
package nlp

import (
	"strings"
)

// ════════════════════════════════════════════════════════════════
// PROCESSOR
// ════════════════════════════════════════════════════════════════

// Processor handles natural language processing of input.
type Processor struct {
	registry *PatternRegistry
	enabled  bool
}

// New creates a new NLP processor with all patterns registered.
func New() *Processor {
	p := &Processor{
		registry: NewRegistry(),
		enabled:  true,
	}

	// Register all pattern categories
	RegisterDateTimePatterns(p.registry)
	RegisterCurrencyPatterns(p.registry)
	RegisterFinancePatterns(p.registry)
	RegisterWeatherPatterns(p.registry)

	return p
}

// NewWithPatterns creates a processor with custom patterns only.
func NewWithPatterns(patterns []*Pattern) *Processor {
	p := &Processor{
		registry: NewRegistry(),
		enabled:  true,
	}
	p.registry.RegisterAll(patterns)
	return p
}

// Enable turns on NLP processing.
func (p *Processor) Enable() {
	p.enabled = true
}

// Disable turns off NLP processing.
func (p *Processor) Disable() {
	p.enabled = false
}

// IsEnabled returns whether NLP processing is enabled.
func (p *Processor) IsEnabled() bool {
	return p.enabled
}

// ════════════════════════════════════════════════════════════════
// PROCESSING
// ════════════════════════════════════════════════════════════════

// Process attempts to transform natural language input into a structured expression.
// Returns the transformed expression and true if successful, or the original input
// and false if no transformation was applied.
func (p *Processor) Process(input string) (string, bool) {
	if !p.enabled {
		return input, false
	}

	// Skip empty input
	input = strings.TrimSpace(input)
	if input == "" {
		return input, false
	}

	// Handle assignments: try NLP on the right-hand side
	if eqIdx := strings.Index(input, "="); eqIdx > 0 && !strings.Contains(input, "==") {
		lhs := input[:eqIdx]
		rhs := strings.TrimSpace(input[eqIdx+1:])
		if rhs != "" && !looksLikeExpression(rhs) {
			match := p.registry.Match(rhs)
			if match != nil {
				return lhs + "= " + match.Transformed, true
			}
		}
		return input, false
	}

	// Skip if it looks like it's already a structured expression
	if looksLikeExpression(input) {
		return input, false
	}

	// Try to match against patterns
	match := p.registry.Match(input)
	if match != nil {
		return match.Transformed, true
	}

	return input, false
}

// ProcessWithInfo returns detailed information about the processing.
func (p *Processor) ProcessWithInfo(input string) *ProcessResult {
	result := &ProcessResult{
		Original: input,
		Output:   input,
		Matched:  false,
	}

	if !p.enabled {
		return result
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return result
	}

	// Handle assignments: try NLP on the right-hand side
	if eqIdx := strings.Index(input, "="); eqIdx > 0 && !strings.Contains(input, "==") {
		lhs := input[:eqIdx]
		rhs := strings.TrimSpace(input[eqIdx+1:])
		if rhs != "" && !looksLikeExpression(rhs) {
			match := p.registry.Match(rhs)
			if match != nil {
				result.Output = lhs + "= " + match.Transformed
				result.Matched = true
				result.PatternName = match.Pattern.Name
				result.Captures = match.Captures
				return result
			}
		}
		result.SkipReason = "assignment"
		return result
	}

	if looksLikeExpression(input) {
		result.SkipReason = "already structured"
		return result
	}

	match := p.registry.Match(input)
	if match != nil {
		result.Output = match.Transformed
		result.Matched = true
		result.PatternName = match.Pattern.Name
		result.Captures = match.Captures
	}

	return result
}

// ProcessResult contains detailed information about NLP processing.
type ProcessResult struct {
	Original    string   // Original input
	Output      string   // Transformed output (or original if no match)
	Matched     bool     // Whether a pattern matched
	PatternName string   // Name of matched pattern (if any)
	Captures    []string // Regex captures (if any)
	SkipReason  string   // Reason processing was skipped (if any)
}

// ════════════════════════════════════════════════════════════════
// DETECTION HELPERS
// ════════════════════════════════════════════════════════════════

// looksLikeExpression checks if input appears to already be a structured expression.
func looksLikeExpression(input string) bool {
	// Contains function call syntax
	if strings.Contains(input, "(") && strings.Contains(input, ")") {
		// But not natural language parentheses like "(optional)"
		// Check if it looks like a function call
		for i := 0; i < len(input)-1; i++ {
			if input[i] >= 'a' && input[i] <= 'z' || input[i] >= 'A' && input[i] <= 'Z' {
				// Found letter, check if followed by (
				for j := i + 1; j < len(input); j++ {
					if input[j] == '(' {
						return true
					}
					if input[j] == ' ' {
						break
					}
				}
			}
		}
	}

	// Contains assignment
	if strings.Contains(input, "=") && !strings.Contains(input, "==") {
		return true
	}

	// Starts with a number and has operators (basic math expression)
	if len(input) > 0 && isDigit(input[0]) {
		if containsOperator(input) && !containsNaturalLanguage(input) {
			return true
		}
	}

	// Starts with currency symbol followed by number
	if len(input) > 1 && isCurrencySymbol(rune(input[0])) && isDigit(input[1]) {
		if !containsNaturalLanguage(input) {
			return true
		}
	}

	return false
}

// containsOperator checks if input contains math operators.
func containsOperator(input string) bool {
	operators := []string{"+", "-", "*", "/", "^", "%"}
	for _, op := range operators {
		if strings.Contains(input, op) {
			return true
		}
	}
	return false
}

// containsNaturalLanguage checks if input contains natural language keywords.
func containsNaturalLanguage(input string) bool {
	lower := strings.ToLower(input)
	keywords := []string{
		// Question words
		"what", "how", "when", "where",
		// Conversion
		"convert", "from", "to",
		// Time references
		"today", "tomorrow", "yesterday", "now",
		"ago", "since", "until", "til", "till",
		"next", "last", "this",
		// Time units
		"days", "day", "weeks", "week", "months", "month", "years", "year",
		"hours", "hour", "minutes", "minute", "seconds", "second",
		// Finance
		"payment", "mortgage", "loan", "tip", "split",
		"compound", "interest", "percent",
		// Weather
		"weather", "temperature", "humidity",
		// Ranges
		"between", "and",
	}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// isDigit checks if a byte is a digit.
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// isCurrencySymbol checks if a rune is a currency symbol.
func isCurrencySymbol(r rune) bool {
	return r == '$' || r == '€' || r == '£' || r == '¥' || r == '₿'
}

// ════════════════════════════════════════════════════════════════
// REGISTRY ACCESS
// ════════════════════════════════════════════════════════════════

// RegisterPattern adds a custom pattern to the processor.
func (p *Processor) RegisterPattern(pattern *Pattern) {
	p.registry.Register(pattern)
}

// RegisterPatterns adds multiple custom patterns.
func (p *Processor) RegisterPatterns(patterns []*Pattern) {
	p.registry.RegisterAll(patterns)
}

// PatternCount returns the number of registered patterns.
func (p *Processor) PatternCount() int {
	return len(p.registry.patterns)
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// defaultProcessor is the shared default processor.
var defaultProcessor *Processor

// init initializes the default processor.
func init() {
	defaultProcessor = New()
}

// Process uses the default processor to transform input.
func Process(input string) (string, bool) {
	return defaultProcessor.Process(input)
}

// ProcessWithInfo uses the default processor and returns detailed info.
func ProcessWithInfo(input string) *ProcessResult {
	return defaultProcessor.ProcessWithInfo(input)
}

// Enable enables the default processor.
func Enable() {
	defaultProcessor.Enable()
}

// Disable disables the default processor.
func Disable() {
	defaultProcessor.Disable()
}

// IsEnabled returns whether the default processor is enabled.
func IsEnabled() bool {
	return defaultProcessor.IsEnabled()
}

// ════════════════════════════════════════════════════════════════
// DEBUG / TESTING
// ════════════════════════════════════════════════════════════════

// TestPattern tests a single pattern against input.
func TestPattern(pattern *Pattern, input string) *Match {
	r := NewRegistry()
	r.Register(pattern)
	return r.Match(input)
}

// ListPatterns returns all pattern names in the default processor.
func ListPatterns() []string {
	names := make([]string, 0, len(defaultProcessor.registry.patterns))
	for _, p := range defaultProcessor.registry.patterns {
		names = append(names, p.Name)
	}
	return names
}

// ════════════════════════════════════════════════════════════════
// EXAMPLE EXPRESSIONS
// ════════════════════════════════════════════════════════════════

// Examples returns a list of example natural language expressions
// and their expected transformations.
func Examples() []Example {
	return []Example{
		// Date/Time
		{"now in PST", `nowin("PST")`},
		{"what time is it in Tokyo", `nowin("JST")`},
		{"3pm PST", `intimezone(time(15, 0), "PST")`},
		{"today + 30 days", `adddays(today(), 30)`},
		{"2 weeks from now", `addweeks(today(), 2)`},
		{"30 days ago", `adddays(today(), -30)`},
		{"next monday", `nextweekday(today(), 1)`},
		{"last friday", `lastweekday(today(), 5)`},
		{"days until christmas", `daysbetween(today(), date(year(today()), 12, 25))`},
		{"days since jan 1 2025", `daysbetween(date(2025, 1, 1), today())`},

		// Currency
		{"$100 to EUR", `100 USD in EUR`},
		{"100 USD to GBP", `100 USD in GBP`},
		{"convert 50 euros to dollars", `50 EUR in USD`},
		{"how much is $100 in pounds", `100 USD in GBP`},
		{"1 BTC to USD", `1 BTC in USD`},
		{"1 oz gold to USD", `1 XAU in USD`},

		// Finance
		{"loan payment on $250k at 6.5% for 30 years", `loan(250000, 6.5, 30)`},
		{"mortgage on $500k at 7% for 30 years", `loan(500000, 7, 30)`},
		{"$10000 at 5% for 10 years compounded monthly", `compound(10000, 5, 10, 12)`},
		{"future value of $10000 at 7% for 20 years", `fv(10000, 7, 20)`},
		{"roi on $1000 investment now worth $1500", `roi(1000, 1500)`},
		{"cagr from $10000 to $25000 over 5 years", `cagr(10000, 25000, 5)`},

		// Percentages
		{"20% of 150", `20% of 150`},
		{"what is 15% of $200", `15% of $200`},
		{"percent change from 50 to 75", `percentchange(50, 75)`},
		{"20% off $150", `$120`},

		// Tips
		{"15% tip on $85", `tip(85, 15)`},
		{"split $100 bill with 20% tip 4 ways", `splittip(100, 20, 4)`},
	}
}

// Example represents a natural language example and its expected output.
type Example struct {
	Input    string
	Expected string
}
