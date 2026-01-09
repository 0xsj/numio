// internal/eval/functions_math.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// LOGARITHMS & EXPONENTIALS
// ════════════════════════════════════════════════════════════════

// FnLog returns the natural logarithm.
func FnLog(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("log requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x <= 0 {
		return types.Errorf("log: argument must be positive, got %v", x)
	}
	return types.Number(math.Log(x))
}

// FnLog10 returns the base-10 logarithm.
func FnLog10(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("log10 requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x <= 0 {
		return types.Errorf("log10: argument must be positive, got %v", x)
	}
	return types.Number(math.Log10(x))
}

// FnLog2 returns the base-2 logarithm.
func FnLog2(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("log2 requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x <= 0 {
		return types.Errorf("log2: argument must be positive, got %v", x)
	}
	return types.Number(math.Log2(x))
}

// FnExp returns e^x.
func FnExp(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("exp requires exactly one argument")
	}
	return types.Number(math.Exp(args[0].AsFloat()))
}

// ════════════════════════════════════════════════════════════════
// TRIGONOMETRIC (radians)
// ════════════════════════════════════════════════════════════════

// FnSin returns the sine.
func FnSin(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sin requires exactly one argument")
	}
	return types.Number(math.Sin(args[0].AsFloat()))
}

// FnCos returns the cosine.
func FnCos(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("cos requires exactly one argument")
	}
	return types.Number(math.Cos(args[0].AsFloat()))
}

// FnTan returns the tangent.
func FnTan(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tan requires exactly one argument")
	}
	return types.Number(math.Tan(args[0].AsFloat()))
}

// FnAsin returns the arc sine.
func FnAsin(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("asin requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x < -1 || x > 1 {
		return types.Errorf("asin: argument must be in [-1, 1], got %v", x)
	}
	return types.Number(math.Asin(x))
}

// FnAcos returns the arc cosine.
func FnAcos(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("acos requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x < -1 || x > 1 {
		return types.Errorf("acos: argument must be in [-1, 1], got %v", x)
	}
	return types.Number(math.Acos(x))
}

// FnAtan returns the arc tangent.
func FnAtan(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("atan requires exactly one argument")
	}
	return types.Number(math.Atan(args[0].AsFloat()))
}

// FnAtan2 returns the arc tangent of y/x.
func FnAtan2(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("atan2 requires exactly two arguments")
	}
	y := args[0].AsFloat()
	x := args[1].AsFloat()
	return types.Number(math.Atan2(y, x))
}

// ════════════════════════════════════════════════════════════════
// HYPERBOLIC
// ════════════════════════════════════════════════════════════════

// FnSinh returns the hyperbolic sine.
func FnSinh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sinh requires exactly one argument")
	}
	return types.Number(math.Sinh(args[0].AsFloat()))
}

// FnCosh returns the hyperbolic cosine.
func FnCosh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("cosh requires exactly one argument")
	}
	return types.Number(math.Cosh(args[0].AsFloat()))
}

// FnTanh returns the hyperbolic tangent.
func FnTanh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tanh requires exactly one argument")
	}
	return types.Number(math.Tanh(args[0].AsFloat()))
}

// FnAsinh returns the inverse hyperbolic sine.
func FnAsinh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("asinh requires exactly one argument")
	}
	return types.Number(math.Asinh(args[0].AsFloat()))
}

// FnAcosh returns the inverse hyperbolic cosine.
func FnAcosh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("acosh requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x < 1 {
		return types.Errorf("acosh: argument must be >= 1, got %v", x)
	}
	return types.Number(math.Acosh(x))
}

// FnAtanh returns the inverse hyperbolic tangent.
func FnAtanh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("atanh requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x <= -1 || x >= 1 {
		return types.Errorf("atanh: argument must be in (-1, 1), got %v", x)
	}
	return types.Number(math.Atanh(x))
}

// ════════════════════════════════════════════════════════════════
// ANGLE CONVERSION
// ════════════════════════════════════════════════════════════════

// FnDeg converts radians to degrees.
func FnDeg(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("deg requires exactly one argument")
	}
	radians := args[0].AsFloat()
	return types.Number(radians * 180.0 / math.Pi)
}

// FnRad converts degrees to radians.
func FnRad(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("rad requires exactly one argument")
	}
	degrees := args[0].AsFloat()
	return types.Number(degrees * math.Pi / 180.0)
}

// ════════════════════════════════════════════════════════════════
// COMBINATORICS
// ════════════════════════════════════════════════════════════════

// FnFactorial returns n!
func FnFactorial(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("factorial requires exactly one argument")
	}

	n := args[0].AsFloat()

	if n < 0 {
		return types.Errorf("factorial: argument must be non-negative, got %v", n)
	}
	if n != math.Floor(n) {
		return types.Errorf("factorial: argument must be an integer, got %v", n)
	}
	if n > 170 {
		return types.Errorf("factorial: argument too large (max 170), got %v", n)
	}

	result := 1.0
	for i := 2.0; i <= n; i++ {
		result *= i
	}
	return types.Number(result)
}

// FnPermutations returns nPr (permutations).
func FnPermutations(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("nPr requires exactly two arguments")
	}

	n := args[0].AsFloat()
	r := args[1].AsFloat()

	if n < 0 || r < 0 {
		return types.Error("nPr: arguments must be non-negative")
	}
	if n != math.Floor(n) || r != math.Floor(r) {
		return types.Error("nPr: arguments must be integers")
	}
	if r > n {
		return types.Number(0)
	}

	// nPr = n! / (n-r)!
	result := 1.0
	for i := n; i > n-r; i-- {
		result *= i
	}
	return types.Number(result)
}

// FnCombinations returns nCr (combinations).
func FnCombinations(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("nCr requires exactly two arguments")
	}

	n := args[0].AsFloat()
	r := args[1].AsFloat()

	if n < 0 || r < 0 {
		return types.Error("nCr: arguments must be non-negative")
	}
	if n != math.Floor(n) || r != math.Floor(r) {
		return types.Error("nCr: arguments must be integers")
	}
	if r > n {
		return types.Number(0)
	}

	// Use smaller r for efficiency: C(n,r) = C(n, n-r)
	if r > n-r {
		r = n - r
	}

	// nCr = n! / (r! * (n-r)!)
	result := 1.0
	for i := 0.0; i < r; i++ {
		result = result * (n - i) / (i + 1)
	}
	return types.Number(math.Round(result))
}

// ════════════════════════════════════════════════════════════════
// NUMBER THEORY
// ════════════════════════════════════════════════════════════════

// FnGCD returns the greatest common divisor.
func FnGCD(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("gcd requires at least two arguments")
	}

	result := int64(math.Abs(args[0].AsFloat()))
	for i := 1; i < len(args); i++ {
		b := int64(math.Abs(args[i].AsFloat()))
		result = gcd(result, b)
	}
	return types.Number(float64(result))
}

// FnLCM returns the least common multiple.
func FnLCM(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("lcm requires at least two arguments")
	}

	result := int64(math.Abs(args[0].AsFloat()))
	for i := 1; i < len(args); i++ {
		b := int64(math.Abs(args[i].AsFloat()))
		result = lcm(result, b)
	}
	return types.Number(float64(result))
}

// FnMod returns the modulo.
func FnMod(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("mod requires exactly two arguments")
	}
	a := args[0].AsFloat()
	b := args[1].AsFloat()
	if b == 0 {
		return types.Error("mod: division by zero")
	}
	return types.Number(math.Mod(a, b))
}

// gcd computes greatest common divisor using Euclidean algorithm.
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm computes least common multiple.
func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return (a / gcd(a, b)) * b
}

// ════════════════════════════════════════════════════════════════
// ROUNDING & SIGN
// ════════════════════════════════════════════════════════════════

// FnSign returns the sign of a number (-1, 0, or 1).
func FnSign(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sign requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x > 0 {
		return types.Number(1)
	}
	if x < 0 {
		return types.Number(-1)
	}
	return types.Number(0)
}

// FnTrunc truncates toward zero.
func FnTrunc(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("trunc requires exactly one argument")
	}
	return types.Number(math.Trunc(args[0].AsFloat()))
}

// FnFrac returns the fractional part.
func FnFrac(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("frac requires exactly one argument")
	}
	x := args[0].AsFloat()
	_, frac := math.Modf(x)
	return types.Number(frac)
}

// ════════════════════════════════════════════════════════════════
// SPECIAL FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnHypot returns sqrt(x^2 + y^2).
func FnHypot(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("hypot requires exactly two arguments")
	}
	x := args[0].AsFloat()
	y := args[1].AsFloat()
	return types.Number(math.Hypot(x, y))
}

// FnCbrt returns the cube root.
func FnCbrt(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("cbrt requires exactly one argument")
	}
	return types.Number(math.Cbrt(args[0].AsFloat()))
}
