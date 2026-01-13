// internal/eval/eval.go

package eval

import (
	"strings"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/pkg/types"
)

// Evaluator evaluates AST nodes and produces values.
type Evaluator struct {
	ctx *Context
}

// New creates a new Evaluator with a fresh context.
func New() *Evaluator {
	return &Evaluator{
		ctx: NewContext(),
	}
}

// NewWithContext creates an Evaluator with an existing context.
func NewWithContext(ctx *Context) *Evaluator {
	return &Evaluator{
		ctx: ctx,
	}
}

// Context returns the evaluation context.
func (e *Evaluator) Context() *Context {
	return e.ctx
}

// ════════════════════════════════════════════════════════════════
// PUBLIC EVALUATION METHODS
// ════════════════════════════════════════════════════════════════

// EvalLine evaluates a parsed line and returns the result.
func (e *Evaluator) EvalLine(line *ast.Line) types.Value {
	if line == nil || line.Stmt == nil {
		return types.Empty()
	}

	result := e.evalStmt(line.Stmt)

	// Track result
	lr := LineResult{
		Input: line.Raw,
		Value: result,
	}

	// Check if this was a continuation
	if stmt, ok := line.Stmt.(*ast.ExprStmt); ok {
		if _, isCont := stmt.Expr.(*ast.ContinuationExpr); isCont {
			lr.IsContinuation = true
			e.ctx.MarkLastConsumed()
		}
		if _, isConvCont := stmt.Expr.(*ast.ConversionContinuation); isConvCont {
			lr.IsContinuation = true
			e.ctx.MarkLastConsumed()
		}
	}

	// Check if this was an assignment
	if assign, ok := line.Stmt.(*ast.AssignStmt); ok {
		lr.AssignedVar = assign.Name
	}

	// Update context
	e.ctx.AddLineResult(lr)
	e.ctx.SetPrevious(result)

	return result
}

// EvalExpr evaluates an expression and returns the result.
func (e *Evaluator) EvalExpr(expr ast.Expr) types.Value {
	return e.evalExpr(expr)
}

// ════════════════════════════════════════════════════════════════
// STATEMENT EVALUATION
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalStmt(stmt ast.Stmt) types.Value {
	switch s := stmt.(type) {
	case *ast.EmptyStmt:
		return types.Empty()

	case *ast.CommentStmt:
		return types.Empty()

	case *ast.ExprStmt:
		return e.evalExpr(s.Expr)

	case *ast.AssignStmt:
		return e.evalAssign(s)

	default:
		return types.Error("unknown statement type")
	}
}

func (e *Evaluator) evalAssign(stmt *ast.AssignStmt) types.Value {
	value := e.evalExpr(stmt.Expr)

	if !value.IsError() {
		e.ctx.SetVariable(stmt.Name, value)
	}

	return value
}

// ════════════════════════════════════════════════════════════════
// EXPRESSION EVALUATION
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalExpr(expr ast.Expr) types.Value {
	if expr == nil {
		return types.Empty()
	}

	switch ex := expr.(type) {
	// Literals
	case *ast.NumberLit:
		return types.Number(ex.Value)

	case *ast.PercentLit:
		return types.Percentage(ex.Value)

	case *ast.CurrencyLit:
		return types.CurrencyValue(ex.Amount, ex.Currency)

	case *ast.UnitLit:
		return types.UnitValue(ex.Amount, ex.Unit)

	case *ast.MetalLit:
		return types.MetalValue(ex.Amount, ex.Metal)

	case *ast.CryptoLit:
		return types.CryptoValue(ex.Amount, ex.Crypto)

	case *ast.StringLit:
		return types.StringValue(ex.Value)

	// References
	case *ast.Identifier:
		return e.evalIdentifier(ex)

	// Operators
	case *ast.BinaryExpr:
		return e.evalBinary(ex)

	case *ast.UnaryExpr:
		return e.evalUnary(ex)

	// Special forms
	case *ast.PercentOfExpr:
		return e.evalPercentOf(ex)

	case *ast.ConversionExpr:
		return e.evalConversion(ex)

	case *ast.CallExpr:
		return e.evalCall(ex)

	case *ast.GroupExpr:
		return e.evalExpr(ex.Expr)

	// Continuations
	case *ast.ContinuationExpr:
		return EvalContinuation(ex, e.evalExpr, e.ctx)

	case *ast.ConversionContinuation:
		return EvalConversionContinuation(ex, e.ctx)

	default:
		return types.Error("unknown expression type")
	}
}

// ════════════════════════════════════════════════════════════════
// IDENTIFIER EVALUATION
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalIdentifier(id *ast.Identifier) types.Value {
	// Check for math constants first
	if val, ok := GetMathConstant(id.Name); ok {
		return types.Number(val)
	}

	// Check for variables
	value, ok := e.ctx.GetVariable(id.Name)
	if !ok {
		if e.ctx.IsStrict() {
			return types.Errorf("undefined variable: %s", id.Name)
		}
		// In non-strict mode, treat as zero
		return types.Number(0)
	}
	return value
}

// ════════════════════════════════════════════════════════════════
// BINARY OPERATIONS
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalBinary(expr *ast.BinaryExpr) types.Value {
	left := e.evalExpr(expr.Left)
	if left.IsError() {
		return left
	}

	right := e.evalExpr(expr.Right)
	if right.IsError() {
		return right
	}

	return ApplyBinaryOp(expr.Op, left, right, e.ctx)
}

// ════════════════════════════════════════════════════════════════
// UNARY OPERATIONS
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalUnary(expr *ast.UnaryExpr) types.Value {
	value := e.evalExpr(expr.Expr)
	if value.IsError() {
		return value
	}

	return ApplyUnaryOp(expr.Op, value)
}

// ════════════════════════════════════════════════════════════════
// SPECIAL EXPRESSIONS
// ════════════════════════════════════════════════════════════════

// evalPercentOf handles "X% of Y" expressions.
func (e *Evaluator) evalPercentOf(expr *ast.PercentOfExpr) types.Value {
	percent := e.evalExpr(expr.Percent)
	if percent.IsError() {
		return percent
	}

	value := e.evalExpr(expr.Value)
	if value.IsError() {
		return value
	}

	// Get percentage as decimal
	var pct float64
	if percent.IsPercentage() {
		pct = percent.Num // Already decimal (0.20 for 20%)
	} else {
		pct = percent.AsFloat() / 100.0
	}

	result := value.AsFloat() * pct

	// Preserve value's type
	return value.WithAmount(result)
}

// evalConversion handles "value in target" expressions.
func (e *Evaluator) evalConversion(expr *ast.ConversionExpr) types.Value {
	value := e.evalExpr(expr.Value)
	if value.IsError() {
		return value
	}

	return ConvertValue(value, expr.Target, e.ctx)
}

// ════════════════════════════════════════════════════════════════
// FUNCTION CALLS
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalCall(expr *ast.CallExpr) types.Value {
	// Evaluate arguments
	args := make([]types.Value, len(expr.Args))
	for i, arg := range expr.Args {
		val := e.evalExpr(arg)
		if val.IsError() {
			return val
		}
		args[i] = val
	}

	// Look up and call function
	name := strings.ToLower(expr.Name)
	return CallFunction(name, args)
}
