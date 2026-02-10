package session

import (
	"github.com/0xsj/numio/pkg/types"
)

// serializeValue breaks a types.Value into the columns stored in the
// variables table: kind (int), num (float64), str (string),
// type_code (e.g. "USD", "BTC", "km"), period (int).
func serializeValue(v types.Value) (kind int, num float64, str string, typeCode string, period int) {
	kind = int(v.Kind)
	num = v.Num
	str = v.Str
	period = int(v.Period)

	switch v.Kind {
	case types.ValueCurrency:
		if v.Curr != nil {
			typeCode = v.Curr.Code
		}
	case types.ValueCrypto:
		if v.Crypto != nil {
			typeCode = v.Crypto.Code
		}
	case types.ValueMetal:
		if v.Metal != nil {
			typeCode = v.Metal.Code
		}
	case types.ValueWithUnit:
		if v.Unit != nil {
			typeCode = v.Unit.Code
		}
	}
	return
}

// deserializeValue reconstructs a types.Value from DB columns using the
// existing Parse* lookup functions in pkg/types.
func deserializeValue(kind int, num float64, str string, typeCode string, period int) types.Value {
	k := types.ValueKind(kind)
	p := types.Period(period)

	switch k {
	case types.ValueNumber:
		v := types.Number(num)
		v.Period = p
		return v

	case types.ValuePercentage:
		return types.Percentage(num)

	case types.ValueCurrency:
		curr := types.ParseCurrency(typeCode)
		if curr == nil {
			// Fallback: create a minimal currency from code.
			curr = types.CurrencyFromCode(typeCode)
		}
		v := types.CurrencyValue(num, curr)
		v.Period = p
		return v

	case types.ValueCrypto:
		crypto := types.ParseCrypto(typeCode)
		if crypto == nil {
			return types.Number(num)
		}
		v := types.CryptoValue(num, crypto)
		v.Period = p
		return v

	case types.ValueMetal:
		metal := types.ParseMetal(typeCode)
		if metal == nil {
			return types.Number(num)
		}
		v := types.MetalValue(num, metal)
		v.Period = p
		return v

	case types.ValueWithUnit:
		unit := types.ParseUnit(typeCode)
		if unit == nil {
			return types.Number(num)
		}
		v := types.UnitValue(num, unit)
		v.Period = p
		return v

	case types.ValueString:
		return types.StringValue(str)

	default:
		return types.Number(num)
	}
}
