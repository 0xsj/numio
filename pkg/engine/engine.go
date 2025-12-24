// pkg/engine/engine.go

// Package engine provides the public API for numio.
package engine

import (
	"strings"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/internal/eval"
	"github.com/0xsj/numio/internal/parser"
	"github.com/0xsj/numio/pkg/cache"
	"github.com/0xsj/numio/pkg/errors"
	"github.com/0xsj/numio/pkg/types"
)

// Engine is the main entry point for numio calculations.
type Engine struct {
	evaluator *eval.Evaluator
	rateCache *cache.RateCache
}

// New creates a new Engine with default settings.
func New() *Engine {
	rc := cache.New()
	ctx := eval.NewContext()
	ctx.SetRateCacheAdapter(&rateCacheAdapter{rc: rc})

	return &Engine{
		evaluator: eval.NewWithContext(ctx),
		rateCache: rc,
	}
}

// NewWithCache creates an Engine with an existing rate cache.
func NewWithCache(rc *cache.RateCache) *Engine {
	if rc == nil {
		rc = cache.New()
	}
	ctx := eval.NewContext()
	ctx.SetRateCacheAdapter(&rateCacheAdapter{rc: rc})

	return &Engine{
		evaluator: eval.NewWithContext(ctx),
		rateCache: rc,
	}
}

// rateCacheAdapter adapts pkg/cache.RateCache to the interface expected by eval.
type rateCacheAdapter struct {
	rc *cache.RateCache
}

func (a *rateCacheAdapter) GetRate(from, to string) (float64, bool) {
	return a.rc.GetRate(from, to)
}

func (a *rateCacheAdapter) Convert(amount float64, from, to string) (float64, bool) {
	return a.rc.Convert(amount, from, to)
}

func (a *rateCacheAdapter) ConvertValue(v types.Value, target string) (types.Value, bool) {
	return a.rc.ConvertValue(v, target)
}

// ════════════════════════════════════════════════════════════════
// CORE EVALUATION
// ════════════════════════════════════════════════════════════════

// Eval evaluates a single line of input and returns the result.
func (e *Engine) Eval(input string) types.Value {
	// Skip empty lines
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return types.Empty()
	}

	// Skip comment-only lines
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return types.Empty()
	}

	// Parse and evaluate
	line, errs := parser.ParseLine(input)
	if len(errs) > 0 {
		return types.Error(errs[0].Message)
	}

	line.Raw = input
	return e.evaluator.EvalLine(line)
}

// EvalMultiple evaluates multiple lines and returns all results.
func (e *Engine) EvalMultiple(lines []string) []types.Value {
	results := make([]types.Value, len(lines))
	for i, line := range lines {
		results[i] = e.Eval(line)
	}
	return results
}

// EvalFile evaluates a multi-line string (like a file contents).
func (e *Engine) EvalFile(content string) []types.Value {
	lines := strings.Split(content, "\n")
	return e.EvalMultiple(lines)
}

// EvalPreview evaluates an expression without affecting state.
// Useful for live preview while typing.
func (e *Engine) EvalPreview(input string) types.Value {
	// Clone context for preview
	ctx := e.evaluator.Context().Clone()
	ctx.SetRateCacheAdapter(&rateCacheAdapter{rc: e.rateCache})
	tempEval := eval.NewWithContext(ctx)

	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return types.Empty()
	}

	line, errs := parser.ParseLine(input)
	if len(errs) > 0 {
		return types.Error(errs[0].Message)
	}

	line.Raw = input
	return tempEval.EvalLine(line)
}

// ════════════════════════════════════════════════════════════════
// VARIABLES
// ════════════════════════════════════════════════════════════════

// GetVariable retrieves a variable value.
func (e *Engine) GetVariable(name string) (types.Value, bool) {
	return e.evaluator.Context().GetVariable(name)
}

// SetVariable sets a variable value.
func (e *Engine) SetVariable(name string, value types.Value) {
	e.evaluator.Context().SetVariable(name, value)
}

// SetVariableFloat sets a variable to a number value.
func (e *Engine) SetVariableFloat(name string, value float64) {
	e.evaluator.Context().SetVariable(name, types.Number(value))
}

// DeleteVariable removes a variable.
func (e *Engine) DeleteVariable(name string) {
	e.evaluator.Context().DeleteVariable(name)
}

// Variables returns all user-defined variables.
func (e *Engine) Variables() map[string]types.Value {
	return e.evaluator.Context().Variables()
}

// VariableNames returns all variable names.
func (e *Engine) VariableNames() []string {
	return e.evaluator.Context().VariableNames()
}

// HasVariable checks if a variable exists.
func (e *Engine) HasVariable(name string) bool {
	return e.evaluator.Context().HasVariable(name)
}

// ════════════════════════════════════════════════════════════════
// PREVIOUS RESULT
// ════════════════════════════════════════════════════════════════

// Previous returns the previous result (accessible as _ or ANS).
func (e *Engine) Previous() types.Value {
	return e.evaluator.Context().Previous()
}

// HasPrevious returns true if there's a valid previous result.
func (e *Engine) HasPrevious() bool {
	return e.evaluator.Context().HasPrevious()
}

// ════════════════════════════════════════════════════════════════
// TOTALS
// ════════════════════════════════════════════════════════════════

// Total returns the sum of all non-consumed results.
func (e *Engine) Total() types.Value {
	return e.evaluator.Context().Total()
}

// GroupedTotals returns totals grouped by type (currency, unit, etc).
func (e *Engine) GroupedTotals() []types.Value {
	return e.evaluator.Context().GroupedTotals()
}

// ════════════════════════════════════════════════════════════════
// LINE HISTORY
// ════════════════════════════════════════════════════════════════

// LineResult represents the result of evaluating a line.
type LineResult = eval.LineResult

// Lines returns all evaluated line results.
func (e *Engine) Lines() []LineResult {
	return e.evaluator.Context().Lines()
}

// LineCount returns the number of evaluated lines.
func (e *Engine) LineCount() int {
	return len(e.evaluator.Context().Lines())
}

// ════════════════════════════════════════════════════════════════
// EXCHANGE RATES
// ════════════════════════════════════════════════════════════════

// RateCache returns the rate cache.
func (e *Engine) RateCache() *cache.RateCache {
	return e.rateCache
}

// SetRate sets an exchange rate.
func (e *Engine) SetRate(from, to string, rate float64) {
	e.rateCache.SetRate(from, to, rate)
}

// GetRate gets an exchange rate.
func (e *Engine) GetRate(from, to string) (float64, bool) {
	return e.rateCache.GetRate(from, to)
}

// ApplyRawRates applies rates from an API response.
func (e *Engine) ApplyRawRates(rates map[string]float64) {
	e.rateCache.ApplyRawRates(rates)
}

// SaveRatesToFile saves rates to the file cache.
func (e *Engine) SaveRatesToFile() error {
	return e.rateCache.SaveToFile()
}

// LoadRatesFromFile loads rates from the file cache.
func (e *Engine) LoadRatesFromFile() bool {
	return e.rateCache.LoadFromFile()
}

// IsRateCacheValid returns true if the rate cache is not expired.
func (e *Engine) IsRateCacheValid() bool {
	return e.rateCache.IsValid()
}

// Convert converts an amount between currencies/assets.
func (e *Engine) Convert(amount float64, from, to string) (float64, bool) {
	return e.rateCache.Convert(amount, from, to)
}

// ════════════════════════════════════════════════════════════════
// SETTINGS
// ════════════════════════════════════════════════════════════════

// Precision returns the display precision.
func (e *Engine) Precision() int {
	return e.evaluator.Context().Precision()
}

// SetPrecision sets the display precision.
func (e *Engine) SetPrecision(p int) {
	e.evaluator.Context().SetPrecision(p)
}

// IsStrict returns whether strict mode is enabled.
func (e *Engine) IsStrict() bool {
	return e.evaluator.Context().IsStrict()
}

// SetStrict enables or disables strict mode.
// In strict mode, undefined variables cause errors.
func (e *Engine) SetStrict(strict bool) {
	e.evaluator.Context().SetStrict(strict)
}

// ════════════════════════════════════════════════════════════════
// STATE MANAGEMENT
// ════════════════════════════════════════════════════════════════

// Clear resets the engine to initial state.
// Clears variables, line history, and previous result.
// Does not clear the rate cache.
func (e *Engine) Clear() {
	e.evaluator.Context().Clear()
}

// ClearVariables removes all user-defined variables.
func (e *Engine) ClearVariables() {
	e.evaluator.Context().ClearVariables()
}

// ClearLines removes all line history.
func (e *Engine) ClearLines() {
	e.evaluator.Context().ClearLines()
}

// Reset is an alias for Clear.
func (e *Engine) Reset() {
	e.Clear()
}

// Clone creates a copy of the engine (shares rate cache).
func (e *Engine) Clone() *Engine {
	ctx := e.evaluator.Context().Clone()
	ctx.SetRateCacheAdapter(&rateCacheAdapter{rc: e.rateCache})

	return &Engine{
		evaluator: eval.NewWithContext(ctx),
		rateCache: e.rateCache,
	}
}

// ════════════════════════════════════════════════════════════════
// PARSING UTILITIES
// ════════════════════════════════════════════════════════════════

// Parse parses input without evaluating.
func (e *Engine) Parse(input string) (*ast.Line, []*errors.Error) {
	return parser.ParseLine(input)
}

// ParseExpr parses an expression without evaluating.
func (e *Engine) ParseExpr(input string) (ast.Expr, []*errors.Error) {
	return parser.ParseExpr(input)
}

// IsValidExpression checks if an input is a valid expression.
func (e *Engine) IsValidExpression(input string) bool {
	_, errs := parser.ParseLine(input)
	return len(errs) == 0
}

// ════════════════════════════════════════════════════════════════
// TYPE UTILITIES
// ════════════════════════════════════════════════════════════════

// LookupCurrency finds a currency by code, symbol, or alias.
func LookupCurrency(s string) *types.Currency {
	return types.ParseCurrency(s)
}

// LookupUnit finds a unit by code or alias.
func LookupUnit(s string) *types.Unit {
	return types.ParseUnit(s)
}

// LookupCrypto finds a cryptocurrency by code, symbol, or alias.
func LookupCrypto(s string) *types.Crypto {
	return types.ParseCrypto(s)
}

// LookupMetal finds a metal by code or alias.
func LookupMetal(s string) *types.Metal {
	return types.ParseMetal(s)
}

// IsCurrency checks if a string refers to a currency.
func IsCurrency(s string) bool {
	return types.ParseCurrency(s) != nil
}

// IsUnit checks if a string refers to a unit.
func IsUnit(s string) bool {
	return types.ParseUnit(s) != nil
}

// IsCrypto checks if a string refers to a cryptocurrency.
func IsCrypto(s string) bool {
	return types.ParseCrypto(s) != nil
}

// IsMetal checks if a string refers to a metal.
func IsMetal(s string) bool {
	return types.ParseMetal(s) != nil
}

// ════════════════════════════════════════════════════════════════
// VALUE CONSTRUCTORS (convenience)
// ════════════════════════════════════════════════════════════════

// Number creates a number value.
func Number(n float64) types.Value {
	return types.Number(n)
}

// Percentage creates a percentage value (input as decimal, e.g., 0.20 for 20%).
func Percentage(p float64) types.Value {
	return types.Percentage(p)
}

// PercentageFromDisplay creates a percentage from display form (e.g., 20 for 20%).
func PercentageFromDisplay(p float64) types.Value {
	return types.PercentageFromDisplay(p)
}

// Currency creates a currency value.
func Currency(amount float64, currencyCode string) types.Value {
	curr := types.ParseCurrency(currencyCode)
	if curr == nil {
		curr = types.CurrencyFromCode(currencyCode)
	}
	return types.CurrencyValue(amount, curr)
}

// Unit creates a value with a unit.
func Unit(amount float64, unitCode string) types.Value {
	unit := types.ParseUnit(unitCode)
	if unit == nil {
		return types.Error("unknown unit: " + unitCode)
	}
	return types.UnitValue(amount, unit)
}

// ════════════════════════════════════════════════════════════════
// STATIC UTILITIES
// ════════════════════════════════════════════════════════════════

// IsCacheFileValid checks if the rate cache file is valid (not expired).
func IsCacheFileValid() bool {
	return cache.IsCacheFileValid()
}

// AllCurrencies returns all curated currencies.
func AllCurrencies() []types.Currency {
	return types.AllCurrencies()
}

// AllUnits returns all curated units.
func AllUnits() []types.Unit {
	return types.AllUnits()
}

// AllCryptos returns all curated cryptocurrencies.
func AllCryptos() []types.Crypto {
	return types.AllCryptos()
}

// AllMetals returns all curated metals.
func AllMetals() []types.Metal {
	return types.AllMetals()
}

// ════════════════════════════════════════════════════════════════
// QUICK EVAL (stateless convenience)
// ════════════════════════════════════════════════════════════════

// QuickEval evaluates an expression with a fresh engine.
// Useful for one-off calculations.
func QuickEval(input string) types.Value {
	return New().Eval(input)
}

// QuickEvalMultiple evaluates multiple expressions with a fresh engine.
func QuickEvalMultiple(inputs []string) []types.Value {
	return New().EvalMultiple(inputs)
}
