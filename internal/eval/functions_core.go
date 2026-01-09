// internal/eval/functions_core.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// AGGREGATION FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSum returns the sum of all arguments.
func FnSum(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Number(0)
	}

	var total float64
	var resultType types.Value = args[0]

	for _, arg := range args {
		if arg.IsError() {
			return arg
		}
		total += arg.AsFloat()
	}

	return resultType.WithAmount(total)
}

// FnAvg returns the average of all arguments.
func FnAvg(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Number(0)
	}

	sum := FnSum(args)
	if sum.IsError() {
		return sum
	}

	return sum.WithAmount(sum.AsFloat() / float64(len(args)))
}

// FnMin returns the minimum value.
func FnMin(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("min requires at least one argument")
	}

	minVal := args[0]
	minNum := minVal.AsFloat()

	for _, arg := range args[1:] {
		if arg.IsError() {
			return arg
		}
		if arg.AsFloat() < minNum {
			minNum = arg.AsFloat()
			minVal = arg
		}
	}

	return minVal.WithAmount(minNum)
}

// FnMax returns the maximum value.
func FnMax(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("max requires at least one argument")
	}

	maxVal := args[0]
	maxNum := maxVal.AsFloat()

	for _, arg := range args[1:] {
		if arg.IsError() {
			return arg
		}
		if arg.AsFloat() > maxNum {
			maxNum = arg.AsFloat()
			maxVal = arg
		}
	}

	return maxVal.WithAmount(maxNum)
}

// FnCount returns the count of arguments.
func FnCount(args []types.Value) types.Value {
	return types.Number(float64(len(args)))
}

// ════════════════════════════════════════════════════════════════
// BASIC MATH FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnAbs returns the absolute value.
func FnAbs(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("abs requires exactly one argument")
	}
	arg := args[0]
	if arg.IsError() {
		return arg
	}
	return arg.WithAmount(math.Abs(arg.AsFloat()))
}

// FnSqrt returns the square root.
func FnSqrt(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sqrt requires exactly one argument")
	}
	arg := args[0]
	if arg.IsError() {
		return arg
	}
	x := arg.AsFloat()
	if x < 0 {
		return types.Errorf("sqrt: argument must be non-negative, got %v", x)
	}
	return types.Number(math.Sqrt(x))
}

// FnRound rounds to the nearest integer.
func FnRound(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("round requires one or two arguments")
	}
	arg := args[0]
	if arg.IsError() {
		return arg
	}

	// Optional decimal places argument
	places := 0
	if len(args) == 2 {
		places = int(args[1].AsFloat())
		if places < 0 {
			places = 0
		}
	}

	if places == 0 {
		return arg.WithAmount(math.Round(arg.AsFloat()))
	}

	// Round to specific decimal places
	multiplier := math.Pow(10, float64(places))
	rounded := math.Round(arg.AsFloat()*multiplier) / multiplier
	return arg.WithAmount(rounded)
}

// FnFloor rounds down to the nearest integer.
func FnFloor(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("floor requires exactly one argument")
	}
	arg := args[0]
	if arg.IsError() {
		return arg
	}
	return arg.WithAmount(math.Floor(arg.AsFloat()))
}

// FnCeil rounds up to the nearest integer.
func FnCeil(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ceil requires exactly one argument")
	}
	arg := args[0]
	if arg.IsError() {
		return arg
	}
	return arg.WithAmount(math.Ceil(arg.AsFloat()))
}

// FnPow returns base raised to the power of exp.
func FnPow(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("pow requires exactly two arguments")
	}

	base := args[0].AsFloat()
	exp := args[1].AsFloat()

	result := math.Pow(base, exp)

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return types.Error("pow: invalid result")
	}

	return types.Number(result)
}

// ════════════════════════════════════════════════════════════════
// HELPER FOR SIMPLE UNARY FUNCTIONS
// ════════════════════════════════════════════════════════════════

// UnaryMathFn applies a math function to a single argument.
func UnaryMathFn(args []types.Value, fn func(float64) float64, name string) types.Value {
	if len(args) != 1 {
		return types.Errorf("%s requires exactly one argument", name)
	}

	arg := args[0]
	if arg.IsError() {
		return arg
	}

	result := fn(arg.AsFloat())

	// Check for NaN/Inf
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return types.Errorf("%s: invalid result", name)
	}

	return types.Number(result)
}
