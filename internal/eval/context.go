// internal/eval/context.go

// Package eval provides expression evaluation for numio.
package eval

import (
	"strings"
	"sync"

	"github.com/0xsj/numio/internal/cache"
	"github.com/0xsj/numio/internal/types"
)

// Context holds the evaluation state including variables and rate cache.
type Context struct {
	mu sync.RWMutex

	// Variables map
	variables map[string]types.Value

	// Rate cache for currency/crypto conversions
	rateCache *cache.RateCache

	// Previous result (for _ and ANS)
	previous types.Value

	// Line results (for continuation tracking)
	lines []LineResult

	// Settings
	precision int  // Decimal precision for display
	strict    bool // Strict mode (error on undefined variables)
}

// LineResult stores the result of evaluating a single line.
type LineResult struct {
	Input          string      // Original input
	Value          types.Value // Computed value
	IsConsumed     bool        // True if consumed by continuation
	IsContinuation bool        // True if this was a continuation
	AssignedVar    string      // Variable name if assignment
}

// NewContext creates a new evaluation context.
func NewContext() *Context {
	return &Context{
		variables: make(map[string]types.Value),
		rateCache: cache.New(),
		previous:  types.Empty(),
		lines:     nil,
		precision: 2,
		strict:    false,
	}
}

// NewContextWithCache creates a context with an existing rate cache.
func NewContextWithCache(rc *cache.RateCache) *Context {
	ctx := NewContext()
	if rc != nil {
		ctx.rateCache = rc
	}
	return ctx
}

// ════════════════════════════════════════════════════════════════
// VARIABLE OPERATIONS
// ════════════════════════════════════════════════════════════════

// GetVariable retrieves a variable value.
func (c *Context) GetVariable(name string) (types.Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Handle special variables
	lower := strings.ToLower(name)
	if lower == "_" || lower == "ans" {
		if !c.previous.IsEmpty() {
			return c.previous, true
		}
		return types.Empty(), false
	}

	if lower == "total" {
		return c.calculateTotal(), true
	}

	// Regular variable lookup
	v, ok := c.variables[name]
	if !ok {
		// Try case-insensitive
		v, ok = c.variables[lower]
	}
	return v, ok
}

// SetVariable sets a variable value.
func (c *Context) SetVariable(name string, value types.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Don't store special variables
	lower := strings.ToLower(name)
	if lower == "_" || lower == "ans" || lower == "total" {
		return
	}

	c.variables[name] = value
}

// DeleteVariable removes a variable.
func (c *Context) DeleteVariable(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.variables, name)
}

// HasVariable checks if a variable exists.
func (c *Context) HasVariable(name string) bool {
	_, ok := c.GetVariable(name)
	return ok
}

// Variables returns a copy of all user-defined variables.
func (c *Context) Variables() map[string]types.Value {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]types.Value, len(c.variables))
	for k, v := range c.variables {
		result[k] = v
	}
	return result
}

// VariableNames returns all variable names.
func (c *Context) VariableNames() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.variables))
	for k := range c.variables {
		names = append(names, k)
	}
	return names
}

// ClearVariables removes all user-defined variables.
func (c *Context) ClearVariables() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.variables = make(map[string]types.Value)
}

// ════════════════════════════════════════════════════════════════
// PREVIOUS RESULT (_, ANS)
// ════════════════════════════════════════════════════════════════

// Previous returns the previous non-empty result.
func (c *Context) Previous() types.Value {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.previous
}

// SetPrevious sets the previous result.
func (c *Context) SetPrevious(v types.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Only store non-empty, non-error values
	if !v.IsEmpty() && !v.IsError() {
		c.previous = v
	}
}

// HasPrevious returns true if there's a valid previous result.
func (c *Context) HasPrevious() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return !c.previous.IsEmpty() && !c.previous.IsError()
}

// ════════════════════════════════════════════════════════════════
// LINE TRACKING
// ════════════════════════════════════════════════════════════════

// AddLineResult adds a line result to the history.
func (c *Context) AddLineResult(result LineResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lines = append(c.lines, result)
}

// Lines returns all line results.
func (c *Context) Lines() []LineResult {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]LineResult, len(c.lines))
	copy(result, c.lines)
	return result
}

// LastLine returns the last non-empty, non-error line result.
func (c *Context) LastLine() (LineResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := len(c.lines) - 1; i >= 0; i-- {
		lr := c.lines[i]
		if !lr.Value.IsEmpty() && !lr.Value.IsError() {
			return lr, true
		}
	}
	return LineResult{}, false
}

// MarkLastConsumed marks the last valid line as consumed by continuation.
func (c *Context) MarkLastConsumed() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := len(c.lines) - 1; i >= 0; i-- {
		if !c.lines[i].Value.IsEmpty() && !c.lines[i].Value.IsError() {
			c.lines[i].IsConsumed = true
			return
		}
	}
}

// ClearLines removes all line history.
func (c *Context) ClearLines() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lines = nil
}

// ════════════════════════════════════════════════════════════════
// TOTALS
// ════════════════════════════════════════════════════════════════

// calculateTotal calculates the sum of all non-consumed line values.
func (c *Context) calculateTotal() types.Value {
	var total float64

	for _, lr := range c.lines {
		if lr.IsConsumed {
			continue
		}
		if lr.Value.IsNumeric() {
			total += lr.Value.AsFloat()
		}
	}

	return types.Number(total)
}

// Total returns the running total of all results.
func (c *Context) Total() types.Value {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.calculateTotal()
}

// GroupedTotals returns totals grouped by type (currency, unit type, etc).
func (c *Context) GroupedTotals() []types.Value {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Track totals by type
	currencyTotals := make(map[string]float64)     // currency code -> amount
	unitTotals := make(map[types.UnitType]float64) // unit type -> base amount
	var lastCurrency *types.Currency
	lastUnits := make(map[types.UnitType]*types.Unit)
	var plainTotal float64

	for _, lr := range c.lines {
		if lr.IsConsumed || lr.Value.IsEmpty() || lr.Value.IsError() {
			continue
		}

		switch lr.Value.Kind {
		case types.ValueCurrency:
			if lr.Value.Curr != nil {
				code := lr.Value.Curr.Code
				// Convert to USD for summing
				if usdAmount, ok := c.rateCache.Convert(lr.Value.Num, code, "USD"); ok {
					currencyTotals["USD"] += usdAmount
				} else {
					currencyTotals[code] += lr.Value.Num
				}
				lastCurrency = lr.Value.Curr
			}

		case types.ValueCrypto:
			if lr.Value.Crypto != nil {
				code := lr.Value.Crypto.Code
				// Convert to USD for summing
				if usdAmount, ok := c.rateCache.Convert(lr.Value.Num, code, "USD"); ok {
					currencyTotals["USD"] += usdAmount
				}
			}

		case types.ValueWithUnit:
			if lr.Value.Unit != nil {
				unitType := lr.Value.Unit.Type
				// Convert to base unit for summing
				baseAmount := lr.Value.Num * lr.Value.Unit.ToBase
				unitTotals[unitType] += baseAmount
				lastUnits[unitType] = lr.Value.Unit
			}

		case types.ValueNumber:
			plainTotal += lr.Value.Num

		case types.ValuePercentage:
			// Percentages usually don't sum meaningfully
		}
	}

	var results []types.Value

	// Add currency total (converted back to last used currency)
	if len(currencyTotals) > 0 {
		usdTotal := currencyTotals["USD"]
		if lastCurrency != nil && lastCurrency.Code != "USD" {
			if converted, ok := c.rateCache.Convert(usdTotal, "USD", lastCurrency.Code); ok {
				results = append(results, types.CurrencyValue(converted, lastCurrency))
			} else {
				usdCurr := types.ParseCurrency("USD")
				results = append(results, types.CurrencyValue(usdTotal, usdCurr))
			}
		} else {
			usdCurr := types.ParseCurrency("USD")
			results = append(results, types.CurrencyValue(usdTotal, usdCurr))
		}
	}

	// Add unit totals (converted back to last used unit of each type)
	for unitType, baseTotal := range unitTotals {
		if lastUnit, ok := lastUnits[unitType]; ok {
			// Convert from base back to last used unit
			amount := baseTotal / lastUnit.ToBase
			results = append(results, types.UnitValue(amount, lastUnit))
		}
	}

	// Add plain number total if any
	if plainTotal != 0 {
		results = append(results, types.Number(plainTotal))
	}

	return results
}

// ════════════════════════════════════════════════════════════════
// RATE CACHE
// ════════════════════════════════════════════════════════════════

// RateCache returns the rate cache.
func (c *Context) RateCache() *cache.RateCache {
	return c.rateCache
}

// SetRateCache sets the rate cache.
func (c *Context) SetRateCache(rc *cache.RateCache) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rateCache = rc
}

// GetRate gets an exchange rate.
func (c *Context) GetRate(from, to string) (float64, bool) {
	if c.rateCache == nil {
		return 0, false
	}
	return c.rateCache.GetRate(from, to)
}

// SetRate sets an exchange rate.
func (c *Context) SetRate(from, to string, rate float64) {
	if c.rateCache == nil {
		c.rateCache = cache.New()
	}
	c.rateCache.SetRate(from, to, rate)
}

// Convert converts an amount between currencies.
func (c *Context) Convert(amount float64, from, to string) (float64, bool) {
	if c.rateCache == nil {
		return 0, false
	}
	return c.rateCache.Convert(amount, from, to)
}

// ConvertValue converts a value to a target currency/unit.
func (c *Context) ConvertValue(v types.Value, target string) (types.Value, bool) {
	// Handle unit conversion
	if v.Kind == types.ValueWithUnit && v.Unit != nil {
		targetUnit := types.ParseUnit(target)
		if targetUnit != nil {
			converted, ok := v.Unit.ConvertTo(v.Num, targetUnit)
			if ok {
				return types.UnitValue(converted, targetUnit), true
			}
		}
		return v, false
	}

	// Handle currency/crypto conversion
	if c.rateCache == nil {
		return v, false
	}
	return c.rateCache.ConvertValue(v, target)
}

// ════════════════════════════════════════════════════════════════
// SETTINGS
// ════════════════════════════════════════════════════════════════

// Precision returns the display precision.
func (c *Context) Precision() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.precision
}

// SetPrecision sets the display precision.
func (c *Context) SetPrecision(p int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if p >= 0 && p <= 15 {
		c.precision = p
	}
}

// IsStrict returns whether strict mode is enabled.
func (c *Context) IsStrict() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.strict
}

// SetStrict enables or disables strict mode.
func (c *Context) SetStrict(strict bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.strict = strict
}

// ════════════════════════════════════════════════════════════════
// RESET / CLEAR
// ════════════════════════════════════════════════════════════════

// Clear resets the context to initial state.
func (c *Context) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.variables = make(map[string]types.Value)
	c.previous = types.Empty()
	c.lines = nil
}

// Reset is an alias for Clear.
func (c *Context) Reset() {
	c.Clear()
}

// ════════════════════════════════════════════════════════════════
// CLONE / SNAPSHOT
// ════════════════════════════════════════════════════════════════

// Clone creates a copy of the context (shares rate cache).
func (c *Context) Clone() *Context {
	c.mu.RLock()
	defer c.mu.RUnlock()

	clone := &Context{
		variables: make(map[string]types.Value, len(c.variables)),
		rateCache: c.rateCache, // Shared
		previous:  c.previous,
		lines:     make([]LineResult, len(c.lines)),
		precision: c.precision,
		strict:    c.strict,
	}

	for k, v := range c.variables {
		clone.variables[k] = v
	}
	copy(clone.lines, c.lines)

	return clone
}

// Snapshot returns a read-only snapshot of the current state.
type Snapshot struct {
	Variables map[string]types.Value
	Previous  types.Value
	Lines     []LineResult
	Total     types.Value
}

// Snapshot creates a snapshot of the current state.
func (c *Context) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	vars := make(map[string]types.Value, len(c.variables))
	for k, v := range c.variables {
		vars[k] = v
	}

	lines := make([]LineResult, len(c.lines))
	copy(lines, c.lines)

	return Snapshot{
		Variables: vars,
		Previous:  c.previous,
		Lines:     lines,
		Total:     c.calculateTotal(),
	}
}