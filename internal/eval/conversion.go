// internal/eval/conversion.go

package eval

import (
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// VALUE CONVERSION
// ════════════════════════════════════════════════════════════════

// ConvertValue converts a value to a target currency/unit.
func ConvertValue(value types.Value, target string, ctx *Context) types.Value {
	// Try unit conversion first
	if value.IsUnit() && value.Unit != nil {
		return convertUnit(value, target)
	}

	// Try currency/crypto/metal conversion via context
	if ctx != nil {
		converted, ok := ctx.ConvertValue(value, target)
		if ok {
			return converted
		}
	}

	// Check if target is valid but conversion unavailable
	if types.ParseCurrency(target) != nil || types.ParseCrypto(target) != nil {
		return types.Errorf("no rate available for conversion to %s", target)
	}
	if types.ParseMetal(target) != nil {
		return types.Errorf("no rate available for conversion to %s", target)
	}
	if types.ParseUnit(target) != nil {
		return types.Errorf("cannot convert to %s (incompatible types)", target)
	}

	return types.Errorf("unknown target: %s", target)
}

// convertUnit handles unit-to-unit conversion.
func convertUnit(value types.Value, target string) types.Value {
	targetUnit := types.ParseUnit(target)
	if targetUnit == nil {
		return types.Errorf("unknown unit: %s", target)
	}

	converted, ok := value.Unit.ConvertTo(value.Num, targetUnit)
	if ok {
		return types.UnitValue(converted, targetUnit)
	}

	return types.Errorf("cannot convert %s to %s", value.Unit.Code, target)
}

// ════════════════════════════════════════════════════════════════
// TYPE CHECKING HELPERS
// ════════════════════════════════════════════════════════════════

// IsConvertibleTo checks if a value can be converted to the target.
func IsConvertibleTo(value types.Value, target string) bool {
	// Units
	if value.IsUnit() && value.Unit != nil {
		targetUnit := types.ParseUnit(target)
		if targetUnit != nil {
			return value.Unit.Type == targetUnit.Type
		}
		return false
	}

	// Currencies
	if value.IsCurrency() {
		return types.ParseCurrency(target) != nil ||
			types.ParseCrypto(target) != nil ||
			types.ParseMetal(target) != nil
	}

	// Crypto
	if value.IsCrypto() {
		return types.ParseCurrency(target) != nil ||
			types.ParseCrypto(target) != nil
	}

	// Metal
	if value.IsMetal() {
		return types.ParseCurrency(target) != nil ||
			types.ParseMetal(target) != nil
	}

	return false
}

// GetTargetType determines what type a target string represents.
func GetTargetType(target string) TargetType {
	if types.ParseUnit(target) != nil {
		return TargetUnit
	}
	if types.ParseCurrency(target) != nil {
		return TargetCurrency
	}
	if types.ParseCrypto(target) != nil {
		return TargetCrypto
	}
	if types.ParseMetal(target) != nil {
		return TargetMetal
	}
	return TargetUnknown
}

// TargetType represents the type of a conversion target.
type TargetType int

const (
	TargetUnknown TargetType = iota
	TargetUnit
	TargetCurrency
	TargetCrypto
	TargetMetal
)

// String returns a string representation of the target type.
func (t TargetType) String() string {
	switch t {
	case TargetUnit:
		return "unit"
	case TargetCurrency:
		return "currency"
	case TargetCrypto:
		return "crypto"
	case TargetMetal:
		return "metal"
	default:
		return "unknown"
	}
}
