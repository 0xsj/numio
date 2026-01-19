// internal/explain/builder.go

// Package explain provides step-by-step explanation of calculations.
package explain

import (
	"strings"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/pkg/types"
)

// Builder constructs an evaluation trace.
type Builder struct {
	trace   *Trace
	stack   []*Step // Stack for nested expressions
	enabled bool
}

// NewBuilder creates a new trace builder.
func NewBuilder(input string) *Builder {
	return &Builder{
		trace:   NewTrace(input),
		stack:   make([]*Step, 0),
		enabled: true,
	}
}

// NewDisabledBuilder creates a builder that doesn't record (for performance).
func NewDisabledBuilder() *Builder {
	return &Builder{
		enabled: false,
	}
}

// IsEnabled returns true if tracing is active.
func (b *Builder) IsEnabled() bool {
	return b.enabled
}

// Trace returns the completed trace.
func (b *Builder) Trace() *Trace {
	if !b.enabled {
		return nil
	}
	return b.trace
}

// ════════════════════════════════════════════════════════════════
// STEP RECORDING
// ════════════════════════════════════════════════════════════════

// RecordLiteral records a literal value evaluation.
func (b *Builder) RecordLiteral(expr ast.Expr, value types.Value) *Step {
	if !b.enabled {
		return nil
	}

	step := NewStep(exprToString(expr), value)
	b.addStep(step)
	return step
}

// RecordVariable records a variable lookup.
func (b *Builder) RecordVariable(name string, value types.Value) *Step {
	if !b.enabled {
		return nil
	}

	step := NewStep(name, value)
	b.addStep(step)
	return step
}

// RecordBinaryOp records a binary operation.
func (b *Builder) RecordBinaryOp(left types.Value, op ast.BinaryOp, right types.Value, result types.Value) *Step {
	if !b.enabled {
		return nil
	}

	opStr := binaryOpToString(op)
	step := NewBinaryStep(
		valueToString(left),
		opStr,
		valueToString(right),
		result,
	)
	b.addStep(step)
	return step
}

// RecordUnaryOp records a unary operation.
func (b *Builder) RecordUnaryOp(op ast.UnaryOp, operand types.Value, result types.Value) *Step {
	if !b.enabled {
		return nil
	}

	opStr := unaryOpToString(op)
	step := NewUnaryStep(opStr, valueToString(operand), result)
	b.addStep(step)
	return step
}

// RecordFuncCall records a function call.
func (b *Builder) RecordFuncCall(name string, args []types.Value, result types.Value) *Step {
	if !b.enabled {
		return nil
	}

	// Build argument string
	argStrs := make([]string, len(args))
	for i, arg := range args {
		argStrs[i] = valueToString(arg)
	}
	exprStr := name + "(" + strings.Join(argStrs, ", ") + ")"

	step := NewStep(exprStr, result)
	step.Op = name
	b.addStep(step)
	return step
}

// RecordConversion records a unit/currency conversion.
func (b *Builder) RecordConversion(from types.Value, target string, result types.Value) *Step {
	if !b.enabled {
		return nil
	}

	exprStr := valueToString(from) + " → " + target
	step := NewStep(exprStr, result)
	step.Op = "convert"
	b.addStep(step)
	return step
}

// RecordPercentOf records a "percent of" expression.
// RecordPercentOf records a "percent of" expression.
func (b *Builder) RecordPercentOf(percent types.Value, value types.Value, result types.Value) *Step {
	if !b.enabled {
		return nil
	}

	step := NewBinaryStep(
		valueToString(percent),
		"of",
		valueToString(value),
		result,
	)
	b.addStep(step)
	return step
}

// RecordGroup records a grouped (parenthesized) expression.
func (b *Builder) RecordGroup(inner types.Value, result types.Value) *Step {
	if !b.enabled {
		return nil
	}

	step := NewStep("("+valueToString(inner)+")", result)
	b.addStep(step)
	return step
}

// ════════════════════════════════════════════════════════════════
// NESTING SUPPORT
// ════════════════════════════════════════════════════════════════

// PushContext starts a nested evaluation context.
// Returns the parent step that children will be added to.
func (b *Builder) PushContext(exprStr string) {
	if !b.enabled {
		return
	}

	parent := &Step{
		ExprString: exprStr,
		SubSteps:   make([]*Step, 0),
	}
	b.stack = append(b.stack, parent)
}

// PopContext ends a nested evaluation context and records the result.
func (b *Builder) PopContext(result types.Value) *Step {
	if !b.enabled || len(b.stack) == 0 {
		return nil
	}

	// Pop the parent step
	parent := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	parent.Result = result

	// Add to parent's parent or to trace root
	b.addStep(parent)
	return parent
}

// CurrentDepth returns the current nesting depth.
func (b *Builder) CurrentDepth() int {
	return len(b.stack)
}

// ════════════════════════════════════════════════════════════════
// FINALIZATION
// ════════════════════════════════════════════════════════════════

// Finalize completes the trace with the final result.
func (b *Builder) Finalize(result types.Value) *Trace {
	if !b.enabled {
		return nil
	}

	b.trace.FinalResult = result

	// If we have steps but no root, use the last step as root
	if b.trace.Root == nil && len(b.trace.Steps) > 0 {
		b.trace.Root = b.trace.Steps[len(b.trace.Steps)-1]
	}

	return b.trace
}

// ════════════════════════════════════════════════════════════════
// INTERNAL HELPERS
// ════════════════════════════════════════════════════════════════

// addStep adds a step to the current context or trace.
func (b *Builder) addStep(step *Step) {
	if len(b.stack) > 0 {
		// Add to current parent
		parent := b.stack[len(b.stack)-1]
		parent.AddSubStep(step)
	} else {
		// Add to trace directly
		b.trace.AddStep(step)
		b.trace.Root = step // Latest top-level step becomes root
	}
}

// ════════════════════════════════════════════════════════════════
// STRING CONVERSION HELPERS
// ════════════════════════════════════════════════════════════════

// exprToString converts an AST expression to a string representation.
func exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	return expr.String()
}

// valueToString converts a Value to a display string.
func valueToString(v types.Value) string {
	return v.String()
}

// binaryOpToString converts a BinaryOp to its string symbol.
func binaryOpToString(op ast.BinaryOp) string {
	switch op {
	case ast.OpAdd:
		return "+"
	case ast.OpSub:
		return "-"
	case ast.OpMul:
		return "*"
	case ast.OpDiv:
		return "/"
	case ast.OpPow:
		return "^"
	case ast.OpMod:
		return "%"
	default:
		return "?"
	}
}

// unaryOpToString converts a UnaryOp to its string symbol.
func unaryOpToString(op ast.UnaryOp) string {
	switch op {
	case ast.OpNeg:
		return "-"
	case ast.OpPos:
		return "+"
	default:
		return "?"
	}
}
