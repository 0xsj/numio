// internal/eval/functions_math.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// LOGARITHMS & EXPONENTIALS
// ════════════════════════════════════════════════════════════════

func fnLog2(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("log2 requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x <= 0 {
		return types.Errorf("log2: argument must be positive, got %v", x)
	}
	return types.Number(math.Log2(x))
}

// ════════════════════════════════════════════════════════════════
// TRIGONOMETRIC (additional)
// ════════════════════════════════════════════════════════════════

func fnAtan2(args []types.Value) types.Value {
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

func fnSinh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sinh requires exactly one argument")
	}
	return types.Number(math.Sinh(args[0].AsFloat()))
}

func fnCosh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("cosh requires exactly one argument")
	}
	return types.Number(math.Cosh(args[0].AsFloat()))
}

func fnTanh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tanh requires exactly one argument")
	}
	return types.Number(math.Tanh(args[0].AsFloat()))
}

func fnAsinh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("asinh requires exactly one argument")
	}
	return types.Number(math.Asinh(args[0].AsFloat()))
}

func fnAcosh(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("acosh requires exactly one argument")
	}
	x := args[0].AsFloat()
	if x < 1 {
		return types.Errorf("acosh: argument must be >= 1, got %v", x)
	}
	return types.Number(math.Acosh(x))
}

func fnAtanh(args []types.Value) types.Value {
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

func fnDeg(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("deg requires exactly one argument")
	}
	radians := args[0].AsFloat()
	return types.Number(radians * 180.0 / math.Pi)
}

func fnRad(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("rad requires exactly one argument")
	}
	degrees := args[0].AsFloat()
	return types.Number(degrees * math.Pi / 180.0)
}

// ════════════════════════════════════════════════════════════════
// COMBINATORICS
// ════════════════════════════════════════════════════════════════

func fnFactorial(args []types.Value) types.Value {
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

func fnPermutations(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("nPr requires exactly two arguments")
	}

	n := args[0].AsFloat()
	r := args[1].AsFloat()

	if n < 0 || r < 0 {
		return types.Errorf("nPr: arguments must be non-negative")
	}
	if n != math.Floor(n) || r != math.Floor(r) {
		return types.Errorf("nPr: arguments must be integers")
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

func fnCombinations(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("nCr requires exactly two arguments")
	}

	n := args[0].AsFloat()
	r := args[1].AsFloat()

	if n < 0 || r < 0 {
		return types.Errorf("nCr: arguments must be non-negative")
	}
	if n != math.Floor(n) || r != math.Floor(r) {
		return types.Errorf("nCr: arguments must be integers")
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

func fnGCD(args []types.Value) types.Value {
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

func fnLCM(args []types.Value) types.Value {
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

func fnMod(args []types.Value) types.Value {
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

func fnSign(args []types.Value) types.Value {
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

func fnTrunc(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("trunc requires exactly one argument")
	}
	return types.Number(math.Trunc(args[0].AsFloat()))
}

func fnFrac(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("frac requires exactly one argument")
	}
	x := args[0].AsFloat()
	_, frac := math.Modf(x)
	return types.Number(frac)
}

// ════════════════════════════════════════════════════════════════
// SPECIAL
// ════════════════════════════════════════════════════════════════

func fnHypot(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("hypot requires exactly two arguments")
	}
	x := args[0].AsFloat()
	y := args[1].AsFloat()
	return types.Number(math.Hypot(x, y))
}

func fnCbrt(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("cbrt requires exactly one argument")
	}
	return types.Number(math.Cbrt(args[0].AsFloat()))
}
