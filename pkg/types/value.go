// internal/types/value.go

package types

import (
	"strings"
)

// ValueKind represents the type of a Value.
type ValueKind int

const (
	ValueEmpty      ValueKind = iota // No value (empty line, comment)
	ValueNumber                      // Plain number: 42, 3.14
	ValuePercentage                  // Percentage: 20% (stored as 0.20)
	ValueCurrency                    // Currency: $100, €50
	ValueWithUnit                    // Value with unit: 5 km, 2 hours
	ValueMetal                       // Precious metal: 1 oz gold
	ValueCrypto                      // Cryptocurrency: 0.5 BTC
	ValueError                       // Error during evaluation
)

// String returns the kind name.
func (k ValueKind) String() string {
	switch k {
	case ValueEmpty:
		return "empty"
	case ValueNumber:
		return "number"
	case ValuePercentage:
		return "percentage"
	case ValueCurrency:
		return "currency"
	case ValueWithUnit:
		return "unit"
	case ValueMetal:
		return "metal"
	case ValueCrypto:
		return "crypto"
	case ValueError:
		return "error"
	default:
		return "unknown"
	}
}

// Value represents a computed value in numio.
// It's a sum type that can hold different kinds of values.
type Value struct {
	Kind ValueKind

	// Numeric value (used by all numeric kinds)
	Num float64

	// Type-specific data
	Curr   *Currency // For ValueCurrency
	Unit   *Unit     // For ValueWithUnit
	Metal  *Metal    // For ValueMetal
	Crypto *Crypto   // For ValueCrypto

	// Error message (for ValueError)
	Err string
}

// ════════════════════════════════════════════════════════════════
// CONSTRUCTORS
// ════════════════════════════════════════════════════════════════

// Empty returns an empty value.
func Empty() Value {
	return Value{Kind: ValueEmpty}
}

// Number creates a plain number value.
func Number(n float64) Value {
	return Value{
		Kind: ValueNumber,
		Num:  n,
	}
}

// Percentage creates a percentage value.
// The input should be the decimal form (e.g., 0.20 for 20%).
func Percentage(p float64) Value {
	return Value{
		Kind: ValuePercentage,
		Num:  p,
	}
}

// PercentageFromDisplay creates a percentage from display form (e.g., 20 for 20%).
func PercentageFromDisplay(p float64) Value {
	return Value{
		Kind: ValuePercentage,
		Num:  p / 100.0,
	}
}

// CurrencyValue creates a currency value.
func CurrencyValue(amount float64, curr *Currency) Value {
	return Value{
		Kind: ValueCurrency,
		Num:  amount,
		Curr: curr,
	}
}

// UnitValue creates a value with a unit.
func UnitValue(amount float64, unit *Unit) Value {
	return Value{
		Kind: ValueWithUnit,
		Num:  amount,
		Unit: unit,
	}
}

// MetalValue creates a precious metal value.
func MetalValue(amount float64, metal *Metal) Value {
	return Value{
		Kind:  ValueMetal,
		Num:   amount,
		Metal: metal,
	}
}

// CryptoValue creates a cryptocurrency value.
func CryptoValue(amount float64, crypto *Crypto) Value {
	return Value{
		Kind:   ValueCrypto,
		Num:    amount,
		Crypto: crypto,
	}
}

// Error creates an error value.
func Error(message string) Value {
	return Value{
		Kind: ValueError,
		Err:  message,
	}
}

// Errorf creates an error value with formatted message.
func Errorf(format string, args ...any) Value {
	// Simple formatting without fmt package
	msg := format
	for _, arg := range args {
		if idx := strings.Index(msg, "%"); idx >= 0 {
			var replacement string
			switch v := arg.(type) {
			case string:
				replacement = v
			case int:
				replacement = itoa(int64(v))
			case int64:
				replacement = itoa(v)
			case float64:
				replacement = formatFloat(v, 2)
			default:
				replacement = "?"
			}
			// Skip format specifier (e.g., %s, %d, %v)
			end := idx + 2
			if end > len(msg) {
				end = len(msg)
			}
			msg = msg[:idx] + replacement + msg[end:]
		}
	}
	return Value{
		Kind: ValueError,
		Err:  msg,
	}
}

// ════════════════════════════════════════════════════════════════
// PREDICATES
// ════════════════════════════════════════════════════════════════

// IsEmpty returns true if the value is empty.
func (v Value) IsEmpty() bool {
	return v.Kind == ValueEmpty
}

// IsError returns true if the value is an error.
func (v Value) IsError() bool {
	return v.Kind == ValueError
}

// IsNumeric returns true if the value has a numeric component.
func (v Value) IsNumeric() bool {
	switch v.Kind {
	case ValueNumber, ValuePercentage, ValueCurrency, ValueWithUnit, ValueMetal, ValueCrypto:
		return true
	default:
		return false
	}
}

// IsNumber returns true if the value is a plain number.
func (v Value) IsNumber() bool {
	return v.Kind == ValueNumber
}

// IsPercentage returns true if the value is a percentage.
func (v Value) IsPercentage() bool {
	return v.Kind == ValuePercentage
}

// IsCurrency returns true if the value is a currency.
func (v Value) IsCurrency() bool {
	return v.Kind == ValueCurrency
}

// IsUnit returns true if the value has a unit.
func (v Value) IsUnit() bool {
	return v.Kind == ValueWithUnit
}

// IsMetal returns true if the value is a precious metal.
func (v Value) IsMetal() bool {
	return v.Kind == ValueMetal
}

// IsCrypto returns true if the value is a cryptocurrency.
func (v Value) IsCrypto() bool {
	return v.Kind == ValueCrypto
}

// ════════════════════════════════════════════════════════════════
// ACCESSORS
// ════════════════════════════════════════════════════════════════

// AsFloat returns the numeric value as float64.
// Returns 0 for non-numeric values.
func (v Value) AsFloat() float64 {
	return v.Num
}

// AsPercentageDisplay returns the percentage in display form (e.g., 20 for 20%).
func (v Value) AsPercentageDisplay() float64 {
	if v.Kind == ValuePercentage {
		return v.Num * 100.0
	}
	return v.Num
}

// ErrorMessage returns the error message, or empty string if not an error.
func (v Value) ErrorMessage() string {
	return v.Err
}

// UnitType returns the unit type if the value has a unit.
func (v Value) UnitType() (UnitType, bool) {
	if v.Kind == ValueWithUnit && v.Unit != nil {
		return v.Unit.Type, true
	}
	return 0, false
}

// ════════════════════════════════════════════════════════════════
// OPERATIONS
// ════════════════════════════════════════════════════════════════

// WithAmount returns a new value with a different numeric amount.
// Preserves the kind and type information.
func (v Value) WithAmount(amount float64) Value {
	result := v
	result.Num = amount
	return result
}

// Negate returns the negated value.
func (v Value) Negate() Value {
	if v.IsError() || v.IsEmpty() {
		return v
	}
	return v.WithAmount(-v.Num)
}

// ════════════════════════════════════════════════════════════════
// FORMATTING
// ════════════════════════════════════════════════════════════════

// String returns a human-readable representation of the value.
func (v Value) String() string {
	switch v.Kind {
	case ValueEmpty:
		return ""

	case ValueNumber:
		return formatNumber(v.Num)

	case ValuePercentage:
		return formatNumber(v.Num*100) + "%"

	case ValueCurrency:
		if v.Curr != nil {
			return formatCurrency(v.Num, v.Curr)
		}
		return formatNumber(v.Num)

	case ValueWithUnit:
		if v.Unit != nil {
			return formatNumber(v.Num) + " " + v.Unit.Code
		}
		return formatNumber(v.Num)

	case ValueMetal:
		if v.Metal != nil {
			return formatNumber(v.Num) + " " + v.Metal.Code
		}
		return formatNumber(v.Num)

	case ValueCrypto:
		if v.Crypto != nil {
			return formatCrypto(v.Num, v.Crypto)
		}
		return formatNumber(v.Num)

	case ValueError:
		return "Error: " + v.Err

	default:
		return "?"
	}
}

// formatNumber formats a number with appropriate precision.
func formatNumber(n float64) string {
	// Handle negative
	if n < 0 {
		return "-" + formatNumber(-n)
	}

	// Determine precision based on magnitude
	var decimals int
	if n == float64(int64(n)) {
		decimals = 0
	} else if n >= 100 {
		decimals = 2
	} else if n >= 1 {
		decimals = 2
	} else if n >= 0.01 {
		decimals = 4
	} else {
		decimals = 6
	}

	return formatFloatTrimmed(n, decimals)
}

// formatFloatTrimmed formats a float and trims trailing zeros.
func formatFloatTrimmed(n float64, maxDecimals int) string {
	str := formatFloat(n, maxDecimals)

	// Trim trailing zeros after decimal point
	if strings.Contains(str, ".") {
		str = strings.TrimRight(str, "0")
		str = strings.TrimRight(str, ".")
	}

	return str
}

// formatCurrency formats a currency value.
func formatCurrency(amount float64, curr *Currency) string {
	// Format with 2 decimal places for currency
	numStr := formatFloat(absFloat(amount), 2)

	var result string
	if curr.SymbolAfter {
		result = numStr + curr.Symbol
	} else {
		result = curr.Symbol + numStr
	}

	if amount < 0 {
		result = "-" + result
	}

	return result
}

// formatCrypto formats a cryptocurrency value.
func formatCrypto(amount float64, crypto *Crypto) string {
	// Use crypto's preferred decimal places
	decimals := crypto.Decimals
	if decimals == 0 {
		decimals = 4
	}

	numStr := formatFloatTrimmed(absFloat(amount), decimals)

	// Use symbol if available, otherwise code
	symbol := crypto.Code
	if crypto.HasSymbol() {
		symbol = crypto.Symbol
	}

	var result string
	// Crypto symbols typically come before the amount
	result = symbol + numStr

	if amount < 0 {
		result = "-" + result
	}

	return result
}

// absFloat returns the absolute value of a float.
func absFloat(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}

// ════════════════════════════════════════════════════════════════
// TYPE COMPATIBILITY
// ════════════════════════════════════════════════════════════════

// CanCombineWith checks if this value can be combined with another
// through arithmetic operations.
func (v Value) CanCombineWith(other Value) bool {
	// Errors and empty can't combine
	if v.IsError() || v.IsEmpty() || other.IsError() || other.IsEmpty() {
		return false
	}

	// Percentages can combine with anything numeric
	if v.IsPercentage() || other.IsPercentage() {
		return true
	}

	// Plain numbers can combine with anything
	if v.IsNumber() || other.IsNumber() {
		return true
	}

	// Same kind can combine
	if v.Kind == other.Kind {
		// For units, must be same unit type
		if v.Kind == ValueWithUnit {
			return v.Unit != nil && other.Unit != nil && v.Unit.Type == other.Unit.Type
		}
		// For currency, any currencies can combine (conversion needed)
		return true
	}

	return false
}

// ResultKind determines the resulting kind when combining two values.
// Returns the "stronger" type (currency > unit > number).
func ResultKind(a, b Value) ValueKind {
	// Priority: Currency > Crypto > Metal > Unit > Percentage > Number
	priority := map[ValueKind]int{
		ValueNumber:     1,
		ValuePercentage: 2,
		ValueWithUnit:   3,
		ValueMetal:      4,
		ValueCrypto:     5,
		ValueCurrency:   6,
	}

	pa, oka := priority[a.Kind]
	pb, okb := priority[b.Kind]

	if !oka {
		return b.Kind
	}
	if !okb {
		return a.Kind
	}

	if pa >= pb {
		return a.Kind
	}
	return b.Kind
}

// ════════════════════════════════════════════════════════════════
// JSON-LIKE REPRESENTATION (for API/export)
// ════════════════════════════════════════════════════════════════

// ToMap returns a map representation of the value.
func (v Value) ToMap() map[string]any {
	m := map[string]any{
		"kind": v.Kind.String(),
	}

	switch v.Kind {
	case ValueEmpty:
		// Nothing extra

	case ValueNumber:
		m["value"] = v.Num

	case ValuePercentage:
		m["value"] = v.Num
		m["display"] = v.Num * 100

	case ValueCurrency:
		m["amount"] = v.Num
		if v.Curr != nil {
			m["currency"] = v.Curr.Code
			m["symbol"] = v.Curr.Symbol
		}

	case ValueWithUnit:
		m["amount"] = v.Num
		if v.Unit != nil {
			m["unit"] = v.Unit.Code
			m["unitType"] = v.Unit.Type.String()
		}

	case ValueMetal:
		m["amount"] = v.Num
		if v.Metal != nil {
			m["metal"] = v.Metal.Code
			m["name"] = v.Metal.Name
		}

	case ValueCrypto:
		m["amount"] = v.Num
		if v.Crypto != nil {
			m["crypto"] = v.Crypto.Code
			m["name"] = v.Crypto.Name
		}

	case ValueError:
		m["error"] = v.Err
	}

	m["display"] = v.String()

	return m
}
