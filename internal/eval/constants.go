// internal/eval/constants.go

package eval

import "math"

// Mathematical constants available in expressions.
var MathConstants = map[string]float64{
	// Circle constants
	"pi":  math.Pi,
	"PI":  math.Pi,
	"tau": 2 * math.Pi,
	"TAU": 2 * math.Pi,

	// Euler's number
	"e": math.E,
	"E": math.E,

	// Golden ratio
	"phi": math.Phi,
	"PHI": math.Phi,

	// Other useful constants
	"sqrt2":  math.Sqrt2,
	"sqrt3":  1.7320508075688772,
	"ln2":    math.Ln2,
	"ln10":   math.Ln10,
	"inf":    math.Inf(1),
	"neginf": math.Inf(-1),
}

// IsMathConstant checks if a name is a mathematical constant.
func IsMathConstant(name string) bool {
	_, ok := MathConstants[name]
	return ok
}

// GetMathConstant returns a mathematical constant value.
func GetMathConstant(name string) (float64, bool) {
	v, ok := MathConstants[name]
	return v, ok
}
