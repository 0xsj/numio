// internal/eval/continuation.go

package eval

import (
	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// CONTINUATION EXPRESSIONS
// ════════════════════════════════════════════════════════════════

// EvalContinuation handles "+ 10", "* 2" etc. continuing from previous.
func EvalContinuation(expr *ast.ContinuationExpr, evalExpr func(ast.Expr) types.Value, ctx *Context) types.Value {
	if !ctx.HasPrevious() {
		// No previous - evaluate expression alone
		return evalExpr(expr.Expr)
	}

	prev := ctx.Previous()
	right := evalExpr(expr.Expr)
	if right.IsError() {
		return right
	}

	return ApplyBinaryOp(expr.Op, prev, right, ctx)
}

// EvalConversionContinuation handles "in EUR", "to miles" continuing from previous.
func EvalConversionContinuation(expr *ast.ConversionContinuation, ctx *Context) types.Value {
	if !ctx.HasPrevious() {
		return types.Error("no previous value to convert")
	}

	prev := ctx.Previous()
	return ConvertValue(prev, expr.Target, ctx)
}

// ════════════════════════════════════════════════════════════════
// CONTINUATION DETECTION
// ════════════════════════════════════════════════════════════════

// IsContinuation checks if a statement is a continuation expression.
func IsContinuation(stmt ast.Stmt) bool {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}

	switch exprStmt.Expr.(type) {
	case *ast.ContinuationExpr:
		return true
	case *ast.ConversionContinuation:
		return true
	default:
		return false
	}
}

// GetContinuationType returns the type of continuation.
func GetContinuationType(stmt ast.Stmt) ContinuationType {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return ContinuationNone
	}

	switch exprStmt.Expr.(type) {
	case *ast.ContinuationExpr:
		return ContinuationOperator
	case *ast.ConversionContinuation:
		return ContinuationConversion
	default:
		return ContinuationNone
	}
}

// ContinuationType represents the type of continuation.
type ContinuationType int

const (
	ContinuationNone       ContinuationType = iota
	ContinuationOperator                    // "+ 10", "* 2"
	ContinuationConversion                  // "in EUR", "to miles"
)

// String returns a string representation of the continuation type.
func (c ContinuationType) String() string {
	switch c {
	case ContinuationOperator:
		return "operator"
	case ContinuationConversion:
		return "conversion"
	default:
		return "none"
	}
}
