// internal/eval/functions_calculus.go

package eval

import (
	"math"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// NUMERICAL DERIVATIVE
// ════════════════════════════════════════════════════════════════

// FnNDerivative computes the numerical derivative at a point.
// Uses central difference: f'(x) ≈ (f(x+h) - f(x-h)) / (2h)
// Args: coefficient (a for f(x)=ax²), x, [h]
func FnNDerivative(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("nderivative requires at least 2 arguments: a, x")
	}

	a := args[0].AsFloat()
	x := args[1].AsFloat()

	h := 1e-8
	if len(args) >= 3 {
		h = args[2].AsFloat()
		if h == 0 {
			h = 1e-8
		}
	}

	// For f(x) = a * x^2, derivative is 2ax
	// Using numerical approximation for generality
	fPlus := a * (x + h) * (x + h)
	fMinus := a * (x - h) * (x - h)
	derivative := (fPlus - fMinus) / (2 * h)

	return types.Number(derivative)
}

// FnNDerivative2 computes the second derivative at a point.
// Uses central difference: f”(x) ≈ (f(x+h) - 2f(x) + f(x-h)) / h²
func FnNDerivative2(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("nderivative2 requires at least 2 arguments: a, x")
	}

	a := args[0].AsFloat()
	x := args[1].AsFloat()

	h := 1e-5
	if len(args) >= 3 {
		h = args[2].AsFloat()
		if h == 0 {
			h = 1e-5
		}
	}

	// For f(x) = a*x^2, f''(x) = 2a
	fPlus := a * (x + h) * (x + h)
	fCenter := a * x * x
	fMinus := a * (x - h) * (x - h)

	secondDerivative := (fPlus - 2*fCenter + fMinus) / (h * h)

	return types.Number(secondDerivative)
}

// ════════════════════════════════════════════════════════════════
// NUMERICAL INTEGRATION
// ════════════════════════════════════════════════════════════════

// FnNIntegral computes the numerical integral using Simpson's rule.
// Args: coefficient (a for f(x)=ax²), lower bound, upper bound, [subdivisions]
func FnNIntegral(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("nintegral requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	n := 1000
	if len(args) >= 4 {
		n = int(args[3].AsFloat())
		if n < 2 {
			n = 2
		}
		if n%2 != 0 {
			n++
		}
	}

	h := (b - a) / float64(n)

	sum := evalQuadratic(coef, a) + evalQuadratic(coef, b)

	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		if i%2 == 0 {
			sum += 2 * evalQuadratic(coef, x)
		} else {
			sum += 4 * evalQuadratic(coef, x)
		}
	}

	integral := (h / 3) * sum

	return types.Number(integral)
}

// evalQuadratic evaluates f(x) = coef * x^2
func evalQuadratic(coef, x float64) float64 {
	return coef * x * x
}

// ════════════════════════════════════════════════════════════════
// LIMIT APPROXIMATION
// ════════════════════════════════════════════════════════════════

// FnNLimit approximates the limit as x approaches a value.
// Args: coefficient, approach point, [direction: 1=right, -1=left, 0=both]
func FnNLimit(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("nlimit requires at least 2 arguments: a, x")
	}

	coef := args[0].AsFloat()
	x := args[1].AsFloat()

	dir := 0.0
	if len(args) >= 3 {
		dir = args[2].AsFloat()
	}

	h := 1e-10

	var limitValue float64
	if dir >= 0 {
		limitValue = evalQuadratic(coef, x+h)
	}
	if dir <= 0 {
		leftLimit := evalQuadratic(coef, x-h)
		if dir == 0 {
			limitValue = (limitValue + leftLimit) / 2
		} else {
			limitValue = leftLimit
		}
	}

	return types.Number(limitValue)
}

// ════════════════════════════════════════════════════════════════
// ROOT FINDING
// ════════════════════════════════════════════════════════════════

// FnFindRoot finds a root using bisection method.
// Args: coefficient, lower bound, upper bound, [tolerance]
func FnFindRoot(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("findroot requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	tol := 1e-10
	if len(args) >= 4 {
		tol = args[3].AsFloat()
		if tol <= 0 {
			tol = 1e-10
		}
	}

	maxIter := 100

	fa := evalQuadratic(coef, a)
	fb := evalQuadratic(coef, b)

	if fa*fb > 0 {
		return types.Number((a + b) / 2)
	}

	for i := 0; i < maxIter; i++ {
		mid := (a + b) / 2
		fmid := evalQuadratic(coef, mid)

		if math.Abs(fmid) < tol || (b-a)/2 < tol {
			return types.Number(mid)
		}

		if fa*fmid < 0 {
			b = mid
			fb = fmid
		} else {
			a = mid
			fa = fmid
		}
	}

	return types.Number((a + b) / 2)
}

// ════════════════════════════════════════════════════════════════
// EXTREMA FINDING
// ════════════════════════════════════════════════════════════════

// FnFindMin finds the minimum using golden section search.
// Args: coefficient, lower bound, upper bound, [tolerance]
func FnFindMin(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("findmin requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	tol := 1e-10
	if len(args) >= 4 {
		tol = args[3].AsFloat()
		if tol <= 0 {
			tol = 1e-10
		}
	}

	phi := (1 + math.Sqrt(5)) / 2
	resphi := 2 - phi

	x1 := a + resphi*(b-a)
	x2 := b - resphi*(b-a)
	f1 := evalQuadratic(coef, x1)
	f2 := evalQuadratic(coef, x2)

	maxIter := 100
	for i := 0; i < maxIter && (b-a) > tol; i++ {
		if f1 < f2 {
			b = x2
			x2 = x1
			f2 = f1
			x1 = a + resphi*(b-a)
			f1 = evalQuadratic(coef, x1)
		} else {
			a = x1
			x1 = x2
			f1 = f2
			x2 = b - resphi*(b-a)
			f2 = evalQuadratic(coef, x2)
		}
	}

	return types.Number((a + b) / 2)
}

// FnFindMax finds the maximum using golden section search.
// Args: coefficient, lower bound, upper bound, [tolerance]
func FnFindMax(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("findmax requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	tol := 1e-10
	if len(args) >= 4 {
		tol = args[3].AsFloat()
		if tol <= 0 {
			tol = 1e-10
		}
	}

	phi := (1 + math.Sqrt(5)) / 2
	resphi := 2 - phi

	x1 := a + resphi*(b-a)
	x2 := b - resphi*(b-a)
	f1 := -evalQuadratic(coef, x1)
	f2 := -evalQuadratic(coef, x2)

	maxIter := 100
	for i := 0; i < maxIter && (b-a) > tol; i++ {
		if f1 < f2 {
			b = x2
			x2 = x1
			f2 = f1
			x1 = a + resphi*(b-a)
			f1 = -evalQuadratic(coef, x1)
		} else {
			a = x1
			x1 = x2
			f1 = f2
			x2 = b - resphi*(b-a)
			f2 = -evalQuadratic(coef, x2)
		}
	}

	return types.Number((a + b) / 2)
}

// ════════════════════════════════════════════════════════════════
// SUMMATION AND PRODUCT
// ════════════════════════════════════════════════════════════════

// FnSummation computes Σ f(i) for i from start to end.
// Args: coefficient, start, end
func FnSummation(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("summation requires 3 arguments: a, start, end")
	}

	coef := args[0].AsFloat()
	start := int(args[1].AsFloat())
	end := int(args[2].AsFloat())

	if end < start {
		return types.Number(0)
	}

	if end-start > 1000000 {
		return types.Error("summation range too large (max 1,000,000)")
	}

	var sum float64
	for i := start; i <= end; i++ {
		sum += evalQuadratic(coef, float64(i))
	}

	return types.Number(sum)
}

// FnSeriesProduct computes Π f(i) for i from start to end.
// Args: coefficient, start, end
func FnSeriesProduct(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("seriesproduct requires 3 arguments: a, start, end")
	}

	coef := args[0].AsFloat()
	start := int(args[1].AsFloat())
	end := int(args[2].AsFloat())

	if end < start {
		return types.Number(1)
	}

	if end-start > 1000 {
		return types.Error("seriesproduct range too large (max 1,000)")
	}

	product := 1.0
	for i := start; i <= end; i++ {
		product *= evalQuadratic(coef, float64(i))
	}

	return types.Number(product)
}

// ════════════════════════════════════════════════════════════════
// TAYLOR SERIES
// ════════════════════════════════════════════════════════════════

// FnTaylorSin computes sin(x) using Taylor series.
// Args: x, [terms]
func FnTaylorSin(args []types.Value) types.Value {
	if len(args) < 1 {
		return types.Error("taylorsin requires at least 1 argument: x")
	}

	x := args[0].AsFloat()
	terms := 10
	if len(args) >= 2 {
		terms = int(args[1].AsFloat())
		if terms < 1 {
			terms = 1
		}
		if terms > 50 {
			terms = 50
		}
	}

	result := 0.0
	for n := 0; n < terms; n++ {
		sign := math.Pow(-1, float64(n))
		power := float64(2*n + 1)
		term := sign * math.Pow(x, power) / calcFactorial(int(power))
		result += term
	}

	return types.Number(result)
}

// FnTaylorCos computes cos(x) using Taylor series.
// Args: x, [terms]
func FnTaylorCos(args []types.Value) types.Value {
	if len(args) < 1 {
		return types.Error("taylorcos requires at least 1 argument: x")
	}

	x := args[0].AsFloat()
	terms := 10
	if len(args) >= 2 {
		terms = int(args[1].AsFloat())
		if terms < 1 {
			terms = 1
		}
		if terms > 50 {
			terms = 50
		}
	}

	result := 0.0
	for n := 0; n < terms; n++ {
		sign := math.Pow(-1, float64(n))
		power := float64(2 * n)
		term := sign * math.Pow(x, power) / calcFactorial(int(power))
		result += term
	}

	return types.Number(result)
}

// FnTaylorExp computes e^x using Taylor series.
// Args: x, [terms]
func FnTaylorExp(args []types.Value) types.Value {
	if len(args) < 1 {
		return types.Error("taylorexp requires at least 1 argument: x")
	}

	x := args[0].AsFloat()
	terms := 20
	if len(args) >= 2 {
		terms = int(args[1].AsFloat())
		if terms < 1 {
			terms = 1
		}
		if terms > 50 {
			terms = 50
		}
	}

	result := 0.0
	for n := 0; n < terms; n++ {
		term := math.Pow(x, float64(n)) / calcFactorial(n)
		result += term
	}

	return types.Number(result)
}

// FnTaylorLn computes ln(1+x) using Taylor series (valid for -1 < x <= 1).
// Args: x, [terms]
func FnTaylorLn(args []types.Value) types.Value {
	if len(args) < 1 {
		return types.Error("taylorln requires at least 1 argument: x")
	}

	x := args[0].AsFloat()
	if x <= -1 || x > 1 {
		return types.Error("taylorln: x must be in range (-1, 1]")
	}

	terms := 20
	if len(args) >= 2 {
		terms = int(args[1].AsFloat())
		if terms < 1 {
			terms = 1
		}
		if terms > 100 {
			terms = 100
		}
	}

	// ln(1+x) = x - x²/2 + x³/3 - x⁴/4 + ...
	result := 0.0
	for n := 1; n <= terms; n++ {
		sign := math.Pow(-1, float64(n+1))
		term := sign * math.Pow(x, float64(n)) / float64(n)
		result += term
	}

	return types.Number(result)
}

// calcFactorial computes n! for Taylor series.
func calcFactorial(n int) float64 {
	if n <= 1 {
		return 1
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// RIEMANN SUM
// ════════════════════════════════════════════════════════════════

// FnRiemannLeft computes left Riemann sum.
// Args: coefficient, lower, upper, [rectangles]
func FnRiemannLeft(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("riemannleft requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	n := 100
	if len(args) >= 4 {
		n = int(args[3].AsFloat())
		if n < 1 {
			n = 1
		}
	}

	dx := (b - a) / float64(n)
	sum := 0.0

	for i := 0; i < n; i++ {
		x := a + float64(i)*dx
		sum += evalQuadratic(coef, x) * dx
	}

	return types.Number(sum)
}

// FnRiemannRight computes right Riemann sum.
// Args: coefficient, lower, upper, [rectangles]
func FnRiemannRight(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("riemannright requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	n := 100
	if len(args) >= 4 {
		n = int(args[3].AsFloat())
		if n < 1 {
			n = 1
		}
	}

	dx := (b - a) / float64(n)
	sum := 0.0

	for i := 1; i <= n; i++ {
		x := a + float64(i)*dx
		sum += evalQuadratic(coef, x) * dx
	}

	return types.Number(sum)
}

// FnRiemannMid computes midpoint Riemann sum.
// Args: coefficient, lower, upper, [rectangles]
func FnRiemannMid(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("riemannmid requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	n := 100
	if len(args) >= 4 {
		n = int(args[3].AsFloat())
		if n < 1 {
			n = 1
		}
	}

	dx := (b - a) / float64(n)
	sum := 0.0

	for i := 0; i < n; i++ {
		x := a + (float64(i)+0.5)*dx
		sum += evalQuadratic(coef, x) * dx
	}

	return types.Number(sum)
}

// FnTrapezoidal computes integral using trapezoidal rule.
// Args: coefficient, lower, upper, [subdivisions]
func FnTrapezoidal(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("trapezoidal requires at least 3 arguments: a, lower, upper")
	}

	coef := args[0].AsFloat()
	a := args[1].AsFloat()
	b := args[2].AsFloat()

	n := 100
	if len(args) >= 4 {
		n = int(args[3].AsFloat())
		if n < 1 {
			n = 1
		}
	}

	dx := (b - a) / float64(n)
	sum := (evalQuadratic(coef, a) + evalQuadratic(coef, b)) / 2

	for i := 1; i < n; i++ {
		x := a + float64(i)*dx
		sum += evalQuadratic(coef, x)
	}

	return types.Number(sum * dx)
}
