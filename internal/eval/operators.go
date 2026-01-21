// internal/eval/operators.go

package eval

import (
	"math"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// BINARY OPERATIONS
// ════════════════════════════════════════════════════════════════

// ApplyBinaryOp applies a binary operator to two values.
func ApplyBinaryOp(op ast.BinaryOp, left, right types.Value, ctx *Context) types.Value {
	// Handle percentage operations specially
	if right.IsPercentage() && (op == ast.OpAdd || op == ast.OpSub) {
		return applyPercentageOp(op, left, right)
	}

	// Get numeric values
	leftNum := left.AsFloat()
	rightNum := right.AsFloat()

	// For percentages used in multiplication/division, use decimal value
	if right.IsPercentage() && (op == ast.OpMul || op == ast.OpDiv) {
		rightNum = right.Num
	}
	if left.IsPercentage() && (op == ast.OpMul || op == ast.OpDiv) {
		leftNum = left.Num
	}

	var result float64

	switch op {
	case ast.OpAdd:
		result = leftNum + rightNum
	case ast.OpSub:
		result = leftNum - rightNum
	case ast.OpMul:
		result = leftNum * rightNum
	case ast.OpDiv:
		if rightNum == 0 {
			return types.Error("division by zero")
		}
		result = leftNum / rightNum
	case ast.OpPow:
		result = math.Pow(leftNum, rightNum)
	case ast.OpMod:
		if rightNum == 0 {
			return types.Error("modulo by zero")
		}
		result = math.Mod(leftNum, rightNum)
	default:
		return types.Error("unknown operator")
	}

	// Determine result type based on operands
	return coerceResult(result, left, right, op, ctx)
}

// applyPercentageOp handles "value + percentage" and "value - percentage"
// e.g., 100 + 15% = 115, $50 - 10% = $45
func applyPercentageOp(op ast.BinaryOp, left, right types.Value) types.Value {
	baseValue := left.AsFloat()
	percentage := right.Num // Already in decimal form (0.15 for 15%)

	var result float64
	if op == ast.OpAdd {
		result = baseValue * (1 + percentage)
	} else { // OpSub
		result = baseValue * (1 - percentage)
	}

	// Preserve the left operand's type (including period if it's a rate)
	return left.WithAmount(result)
}

// coerceResult determines the result type based on operands.
func coerceResult(result float64, left, right types.Value, op ast.BinaryOp, ctx *Context) types.Value {
	// For multiplication/division, special handling
	if op == ast.OpMul || op == ast.OpDiv {
		return coerceMultiplyDivide(result, left, right, op)
	}

	// For addition/subtraction, types must be compatible
	if op == ast.OpAdd || op == ast.OpSub {
		return coerceAddSubtract(result, left, right, op, ctx)
	}

	// For power/mod, return plain number
	return types.Number(result)
}

// coerceMultiplyDivide handles type coercion for multiplication and division.
func coerceMultiplyDivide(result float64, left, right types.Value, op ast.BinaryOp) types.Value {
	// If one is a plain number (not a rate), inherit the other's type
	if left.IsNumber() && !left.IsRate() && !right.IsNumber() {
		return right.WithAmount(result)
	}
	if right.IsNumber() && !right.IsRate() && !left.IsNumber() {
		return left.WithAmount(result)
	}

	// If one is a plain number with no period, inherit the rate type
	if left.Kind == types.ValueNumber && !left.IsRate() {
		if right.IsRate() {
			// number * rate = rate (scaled)
			return right.WithAmount(result)
		}
		return right.WithAmount(result)
	}
	if right.Kind == types.ValueNumber && !right.IsRate() {
		if left.IsRate() {
			// rate * number = rate (scaled)
			return left.WithAmount(result)
		}
		return left.WithAmount(result)
	}

	// Both typed - return plain number (or could be unit algebra in future)
	// But if one has a period and is the dominant type, preserve it
	if left.IsRate() && (right.Kind == types.ValueNumber && !right.IsRate()) {
		return left.WithAmount(result)
	}
	if right.IsRate() && (left.Kind == types.ValueNumber && !left.IsRate()) {
		return right.WithAmount(result)
	}

	// Division of rate by number preserves the rate
	if op == ast.OpDiv && left.IsRate() && right.Kind == types.ValueNumber {
		return left.WithAmount(result)
	}

	return types.Number(result)
}

// coerceAddSubtract handles type coercion for addition and subtraction.
func coerceAddSubtract(result float64, left, right types.Value, op ast.BinaryOp, ctx *Context) types.Value {
	// Same type - preserve it (including period)
	if left.Kind == right.Kind {
		// If both are rates with same period, preserve it
		if left.IsRate() && right.IsRate() && left.Period == right.Period {
			return left.WithAmount(result)
		}
		// If both are rates with different periods, this is an error or needs conversion
		if left.IsRate() && right.IsRate() && left.Period != right.Period {
			// Convert right to left's period and recalculate
			rightConverted := right.ConvertRateTo(left.Period)
			if op == ast.OpAdd {
				return left.WithAmount(left.Num + rightConverted.Num)
			}
			return left.WithAmount(left.Num - rightConverted.Num)
		}
		return left.WithAmount(result)
	}

	// One is a plain number - inherit the typed one
	if left.IsNumber() && !left.IsRate() {
		return right.WithAmount(result)
	}
	if right.IsNumber() && !right.IsRate() {
		return left.WithAmount(result)
	}

	// Different typed values - need conversion
	// For currencies, convert right to left's currency
	if left.IsCurrency() && right.IsCurrency() {
		if left.Curr != nil && right.Curr != nil && ctx != nil {
			converted, ok := ctx.Convert(right.Num, right.Curr.Code, left.Curr.Code)
			if ok {
				var newResult float64
				if op == ast.OpAdd {
					newResult = left.Num + converted
				} else {
					newResult = left.Num - converted
				}
				// Preserve period if left is a rate
				res := left.WithAmount(newResult)
				return res
			}
		}
	}

	// For units, convert right to left's unit
	if left.IsUnit() && right.IsUnit() {
		if left.Unit != nil && right.Unit != nil {
			converted, ok := right.Unit.ConvertTo(right.Num, left.Unit)
			if ok {
				var newResult float64
				if op == ast.OpAdd {
					newResult = left.Num + converted
				} else {
					newResult = left.Num - converted
				}
				return left.WithAmount(newResult)
			}
		}
		return types.Error("incompatible units")
	}

	return types.Number(result)
}

// ════════════════════════════════════════════════════════════════
// UNARY OPERATIONS
// ════════════════════════════════════════════════════════════════

// ApplyUnaryOp applies a unary operator to a value.
func ApplyUnaryOp(op ast.UnaryOp, value types.Value) types.Value {
	switch op {
	case ast.OpNeg:
		return value.Negate()
	case ast.OpPos:
		return value
	default:
		return types.Error("unknown unary operator")
	}
}
