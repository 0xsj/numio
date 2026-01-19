// internal/eval/eval.go

package eval

import (
	"strings"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/internal/explain"
	"github.com/0xsj/numio/pkg/types"
)

// Evaluator evaluates AST nodes and produces values.
type Evaluator struct {
	ctx   *Context
	trace *explain.Builder
}

// New creates a new Evaluator with a fresh context.
func New() *Evaluator {
	return &Evaluator{
		ctx:   NewContext(),
		trace: nil,
	}
}

// NewWithContext creates an Evaluator with an existing context.
func NewWithContext(ctx *Context) *Evaluator {
	return &Evaluator{
		ctx:   ctx,
		trace: nil,
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
		if _, isMultiConvCont := stmt.Expr.(*ast.MultiConversionContinuation); isMultiConvCont {
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

// EvalLineWithTrace evaluates a line and returns both result and trace.
func (e *Evaluator) EvalLineWithTrace(line *ast.Line, input string) (types.Value, *explain.Trace) {
	if line == nil || line.Stmt == nil {
		return types.Empty(), nil
	}

	// Enable tracing
	e.trace = explain.NewBuilder(input)

	// Evaluate
	result := e.evalStmt(line.Stmt)

	// Finalize trace
	trace := e.trace.Finalize(result)

	// Disable tracing
	e.trace = nil

	// Track result (same as EvalLine)
	lr := LineResult{
		Input: line.Raw,
		Value: result,
	}

	if stmt, ok := line.Stmt.(*ast.ExprStmt); ok {
		if _, isCont := stmt.Expr.(*ast.ContinuationExpr); isCont {
			lr.IsContinuation = true
			e.ctx.MarkLastConsumed()
		}
		if _, isConvCont := stmt.Expr.(*ast.ConversionContinuation); isConvCont {
			lr.IsContinuation = true
			e.ctx.MarkLastConsumed()
		}
		if _, isMultiConvCont := stmt.Expr.(*ast.MultiConversionContinuation); isMultiConvCont {
			lr.IsContinuation = true
			e.ctx.MarkLastConsumed()
		}
	}

	if assign, ok := line.Stmt.(*ast.AssignStmt); ok {
		lr.AssignedVar = assign.Name
	}

	e.ctx.AddLineResult(lr)
	e.ctx.SetPrevious(result)

	return result, trace
}

// EvalExprWithTrace evaluates an expression and returns both result and trace.
func (e *Evaluator) EvalExprWithTrace(expr ast.Expr, input string) (types.Value, *explain.Trace) {
	// Enable tracing
	e.trace = explain.NewBuilder(input)

	// Evaluate
	result := e.evalExpr(expr)

	// Finalize trace
	trace := e.trace.Finalize(result)

	// Disable tracing
	e.trace = nil

	return result, trace
}

// isTracing returns true if trace recording is active.
func (e *Evaluator) isTracing() bool {
	return e.trace != nil && e.trace.IsEnabled()
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
		result := types.Number(ex.Value)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

	case *ast.PercentLit:
		result := types.Percentage(ex.Value)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

	case *ast.CurrencyLit:
		result := types.CurrencyValue(ex.Amount, ex.Currency)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

	case *ast.UnitLit:
		result := types.UnitValue(ex.Amount, ex.Unit)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

	case *ast.MetalLit:
		result := types.MetalValue(ex.Amount, ex.Metal)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

	case *ast.CryptoLit:
		result := types.CryptoValue(ex.Amount, ex.Crypto)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

	case *ast.StringLit:
		result := types.StringValue(ex.Value)
		if e.isTracing() {
			e.trace.RecordLiteral(ex, result)
		}
		return result

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

	case *ast.MultiConversionExpr:
		return e.evalMultiConversion(ex)

	case *ast.CallExpr:
		return e.evalCall(ex)

	case *ast.GroupExpr:
		result := e.evalExpr(ex.Expr)
		if e.isTracing() {
			e.trace.RecordGroup(result, result)
		}
		return result

	// Continuations
	case *ast.ContinuationExpr:
		return EvalContinuation(ex, e.evalExpr, e.ctx)

	case *ast.ConversionContinuation:
		return EvalConversionContinuation(ex, e.ctx)

	case *ast.MultiConversionContinuation:
		return e.evalMultiConversionContinuation(ex)

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
		result := types.Number(val)
		if e.isTracing() {
			e.trace.RecordVariable(id.Name, result)
		}
		return result
	}

	// Check for variables
	value, ok := e.ctx.GetVariable(id.Name)
	if !ok {
		if e.ctx.IsStrict() {
			return types.Errorf("undefined variable: %s", id.Name)
		}
		// In non-strict mode, treat as zero
		value = types.Number(0)
	}

	if e.isTracing() {
		e.trace.RecordVariable(id.Name, value)
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

	result := ApplyBinaryOp(expr.Op, left, right, e.ctx)

	if e.isTracing() {
		e.trace.RecordBinaryOp(left, expr.Op, right, result)
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// UNARY OPERATIONS
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalUnary(expr *ast.UnaryExpr) types.Value {
	value := e.evalExpr(expr.Expr)
	if value.IsError() {
		return value
	}

	result := ApplyUnaryOp(expr.Op, value)

	if e.isTracing() {
		e.trace.RecordUnaryOp(expr.Op, value, result)
	}

	return result
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

	resultNum := value.AsFloat() * pct

	// Preserve value's type
	result := value.WithAmount(resultNum)

	if e.isTracing() {
		e.trace.RecordPercentOf(percent, value, result)
	}

	return result
}

// evalConversion handles "value in target" expressions.
func (e *Evaluator) evalConversion(expr *ast.ConversionExpr) types.Value {
	value := e.evalExpr(expr.Value)
	if value.IsError() {
		return value
	}

	result := ConvertValue(value, expr.Target, e.ctx)

	if e.isTracing() {
		e.trace.RecordConversion(value, expr.Target, result)
	}

	return result
}

// evalMultiConversion handles "value in target1, target2, ..." expressions.
func (e *Evaluator) evalMultiConversion(expr *ast.MultiConversionExpr) types.Value {
	value := e.evalExpr(expr.Value)
	if value.IsError() {
		return value
	}

	return e.convertToMultiple(value, expr.Targets)
}

// evalMultiConversionContinuation handles "in target1, target2, ..." continuation.
func (e *Evaluator) evalMultiConversionContinuation(expr *ast.MultiConversionContinuation) types.Value {
	// Get previous result
	prev := e.ctx.Previous()
	if prev.IsEmpty() {
		return types.Error("no previous result for conversion")
	}
	if prev.IsError() {
		return prev
	}

	return e.convertToMultiple(prev, expr.Targets)
}

// convertToMultiple converts a value to multiple targets and returns a MultiValue.
func (e *Evaluator) convertToMultiple(value types.Value, targets []string) types.Value {
	results := make([]types.Value, len(targets))

	for i, target := range targets {
		converted := ConvertValue(value, target, e.ctx)
		results[i] = converted

		if e.isTracing() {
			e.trace.RecordConversion(value, target, converted)
		}
	}

	return types.MultiValue(results...)
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
	result := CallFunction(name, args)

	if e.isTracing() {
		e.trace.RecordFuncCall(expr.Name, args, result)
	}

	return result
}
