// internal/explain/explain.go

// Package explain provides step-by-step explanation of calculations.
package explain

import (
	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/pkg/types"
)

// Step represents a single evaluation step.
type Step struct {
	// Expr is the AST node being evaluated (optional, for context)
	Expr ast.Expr

	// ExprString is the string representation of the expression
	ExprString string

	// Result is the computed value of this step
	Result types.Value

	// SubSteps are child steps (for nested expressions)
	SubSteps []*Step

	// Op is the operation being performed (for binary/unary ops)
	Op string

	// Left and Right are operand representations (for binary ops)
	Left  string
	Right string
}

// Trace holds the complete evaluation trace for an expression.
type Trace struct {
	// Input is the original input string
	Input string

	// Root is the top-level step
	Root *Step

	// FinalResult is the final computed value
	FinalResult types.Value

	// Steps is a flat list of all steps in evaluation order
	Steps []*Step
}

// NewTrace creates a new empty trace.
func NewTrace(input string) *Trace {
	return &Trace{
		Input: input,
		Steps: make([]*Step, 0),
	}
}

// NewStep creates a new evaluation step.
func NewStep(exprString string, result types.Value) *Step {
	return &Step{
		ExprString: exprString,
		Result:     result,
		SubSteps:   make([]*Step, 0),
	}
}

// NewBinaryStep creates a step for a binary operation.
func NewBinaryStep(left, op, right string, result types.Value) *Step {
	return &Step{
		ExprString: left + " " + op + " " + right,
		Result:     result,
		Op:         op,
		Left:       left,
		Right:      right,
		SubSteps:   make([]*Step, 0),
	}
}

// NewUnaryStep creates a step for a unary operation.
func NewUnaryStep(op, operand string, result types.Value) *Step {
	return &Step{
		ExprString: op + operand,
		Result:     result,
		Op:         op,
		SubSteps:   make([]*Step, 0),
	}
}

// AddSubStep adds a child step.
func (s *Step) AddSubStep(sub *Step) {
	s.SubSteps = append(s.SubSteps, sub)
}

// HasSubSteps returns true if this step has child steps.
func (s *Step) HasSubSteps() bool {
	return len(s.SubSteps) > 0
}

// AddStep adds a step to the trace.
func (t *Trace) AddStep(step *Step) {
	t.Steps = append(t.Steps, step)
}

// SetRoot sets the root step and final result.
func (t *Trace) SetRoot(step *Step) {
	t.Root = step
	t.FinalResult = step.Result
}

// IsEmpty returns true if the trace has no steps.
func (t *Trace) IsEmpty() bool {
	return len(t.Steps) == 0 && t.Root == nil
}

// Depth returns the maximum nesting depth of the trace.
func (t *Trace) Depth() int {
	if t.Root == nil {
		return 0
	}
	return stepDepth(t.Root)
}

// stepDepth recursively calculates step depth.
func stepDepth(s *Step) int {
	if s == nil || len(s.SubSteps) == 0 {
		return 1
	}
	maxChild := 0
	for _, sub := range s.SubSteps {
		d := stepDepth(sub)
		if d > maxChild {
			maxChild = d
		}
	}
	return maxChild + 1
}

// StepKind identifies the type of step for formatting purposes.
type StepKind int

const (
	StepLiteral    StepKind = iota // A literal value (number, currency, etc.)
	StepVariable                   // A variable reference
	StepBinaryOp                   // A binary operation
	StepUnaryOp                    // A unary operation
	StepFuncCall                   // A function call
	StepConversion                 // A unit/currency conversion
	StepPercentOf                  // A "percent of" expression
	StepGroup                      // A parenthesized group
)

// Kind determines the kind of step based on its properties.
func (s *Step) Kind() StepKind {
	if s.Op != "" {
		if s.Left != "" && s.Right != "" {
			return StepBinaryOp
		}
		return StepUnaryOp
	}
	if len(s.SubSteps) == 0 {
		return StepLiteral
	}
	return StepGroup
}

// OperatorSymbol returns a display-friendly operator symbol.
func OperatorSymbol(op string) string {
	switch op {
	case "+":
		return "+"
	case "-":
		return "−" // minus sign (U+2212)
	case "*", "x", "X":
		return "×" // multiplication sign
	case "/":
		return "÷" // division sign
	case "^", "**":
		return "^"
	case "%":
		return "%"
	default:
		return op
	}
}
