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

	// Local bindings for function evaluation (parameter values)
	localBindings map[string]types.Value
}

// New creates a new Evaluator with a fresh context.
func New() *Evaluator {
	return &Evaluator{
		ctx:           NewContext(),
		trace:         nil,
		localBindings: nil,
	}
}

// NewWithContext creates an Evaluator with an existing context.
func NewWithContext(ctx *Context) *Evaluator {
	return &Evaluator{
		ctx:           ctx,
		trace:         nil,
		localBindings: nil,
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

	// Don't track function definitions in line results
	if _, isFuncDef := line.Stmt.(*ast.FuncDefStmt); isFuncDef {
		return result
	}

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

	// Don't track function definitions in line results
	if _, isFuncDef := line.Stmt.(*ast.FuncDefStmt); isFuncDef {
		return result, trace
	}

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

	case *ast.FuncDefStmt:
		return e.evalFuncDef(s)

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

// evalFuncDef evaluates a function definition statement.
func (e *Evaluator) evalFuncDef(stmt *ast.FuncDefStmt) types.Value {
	// Create user function from AST
	fn := NewUserFunction(stmt)

	// Register in context
	if err := e.ctx.DefineFunc(fn); err != "" {
		return types.Error(err)
	}

	// Return a confirmation message as a string value
	return types.StringValue("defined " + fn.Signature())
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

	// Rate expressions
	case *ast.RateExpr:
		return e.evalRateExpr(ex)

	case *ast.PeriodExpr:
		return e.evalPeriodExpr(ex)

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
// RATE AND PERIOD EVALUATION
// ════════════════════════════════════════════════════════════════

// evalRateExpr evaluates a rate expression (e.g., 45000 TL/month).
func (e *Evaluator) evalRateExpr(expr *ast.RateExpr) types.Value {
	// Evaluate the base value
	value := e.evalExpr(expr.Value)
	if value.IsError() {
		return value
	}

	// Add the period to create a rate value
	result := value.WithPeriod(expr.Period)

	if e.isTracing() {
		e.trace.RecordLiteral(expr, result)
	}

	return result
}

// evalPeriodExpr evaluates a period expression (e.g., year, 2 months).
// Returns a special "period value" that can be used in multiplication with rates.
func (e *Evaluator) evalPeriodExpr(expr *ast.PeriodExpr) types.Value {
	// Create a period value
	// We store it as a number with a special marker so multiplication can detect it
	result := types.Number(expr.Count).WithPeriod(expr.Period)

	if e.isTracing() {
		e.trace.RecordLiteral(expr, result)
	}

	return result
}

// ════════════════════════════════════════════════════════════════
// IDENTIFIER EVALUATION
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalIdentifier(id *ast.Identifier) types.Value {
	// Check local bindings first (function parameters)
	if e.localBindings != nil {
		if val, ok := e.localBindings[id.Name]; ok {
			if e.isTracing() {
				e.trace.RecordVariable(id.Name, val)
			}
			return val
		}
		// Also check case-insensitive
		if val, ok := e.localBindings[strings.ToLower(id.Name)]; ok {
			if e.isTracing() {
				e.trace.RecordVariable(id.Name, val)
			}
			return val
		}
	}

	// Check for math constants
	if val, ok := GetMathConstant(id.Name); ok {
		result := types.Number(val)
		if e.isTracing() {
			e.trace.RecordVariable(id.Name, result)
		}
		return result
	}

	// Check for time period constants (fallback when used as plain identifier)
	if period := types.ParsePeriod(id.Name); period != types.PeriodNone {
		// Return the period's default multiplier as a number
		// This handles cases like "rent * year" when year is not parsed as PeriodExpr
		result := types.Number(period.DefaultMultiplier())
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

	// Special handling for rate * period multiplication
	if expr.Op == ast.OpMul {
		if result, handled := e.tryRatePeriodMultiply(left, right); handled {
			if e.isTracing() {
				e.trace.RecordBinaryOp(left, expr.Op, right, result)
			}
			return result
		}
	}

	// Special handling for rate / period division (rate conversion)
	if expr.Op == ast.OpDiv {
		if result, handled := e.tryRatePeriodDivide(left, right); handled {
			if e.isTracing() {
				e.trace.RecordBinaryOp(left, expr.Op, right, result)
			}
			return result
		}
	}

	result := ApplyBinaryOp(expr.Op, left, right, e.ctx)

	if e.isTracing() {
		e.trace.RecordBinaryOp(left, expr.Op, right, result)
	}

	return result
}

// tryRatePeriodMultiply handles multiplication between rates and periods.
// Returns (result, true) if this was a rate*period operation, (empty, false) otherwise.
func (e *Evaluator) tryRatePeriodMultiply(left, right types.Value) (types.Value, bool) {
	// Case 1: rate * period (e.g., 45000TL/month * year)
	if left.IsRate() && right.IsRate() && right.Kind == types.ValueNumber {
		// right is a period value (number with period set)
		return left.MultiplyByPeriod(right.Period, right.Num), true
	}

	// Case 2: period * rate (e.g., year * 45000TL/month)
	if right.IsRate() && left.IsRate() && left.Kind == types.ValueNumber {
		// left is a period value
		return right.MultiplyByPeriod(left.Period, left.Num), true
	}

	// Case 3: value * period (no rate info, use period's default multiplier)
	// e.g., rent * year where rent is just 45000TL (not a rate)
	if !left.IsRate() && right.IsRate() && right.Kind == types.ValueNumber {
		// right is a period value, left is a regular value
		// Multiply by the period count times the default multiplier relative to month
		multiplier := right.Num * right.Period.DefaultMultiplier()
		return left.WithAmount(left.Num * multiplier), true
	}

	// Case 4: period * value
	if left.IsRate() && left.Kind == types.ValueNumber && !right.IsRate() {
		multiplier := left.Num * left.Period.DefaultMultiplier()
		return right.WithAmount(right.Num * multiplier), true
	}

	return types.Empty(), false
}

// tryRatePeriodDivide handles division of rates by periods (rate conversion).
// e.g., 540000TL/year / 12 = 45000TL/month (if properly structured)
func (e *Evaluator) tryRatePeriodDivide(left, right types.Value) (types.Value, bool) {
	// For now, just handle rate conversion
	// e.g., salary/year -> salary/month would need explicit conversion syntax
	// This is a placeholder for future enhancement

	return types.Empty(), false
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

	name := strings.ToLower(expr.Name)

	// Check for user-defined function first
	if userFn, ok := e.ctx.GetFunc(name); ok {
		result := e.callUserFunc(userFn, args)
		if e.isTracing() {
			e.trace.RecordFuncCall(expr.Name, args, result)
		}
		return result
	}

	// Fall back to built-in function
	result := CallFunction(name, args)

	if e.isTracing() {
		e.trace.RecordFuncCall(expr.Name, args, result)
	}

	return result
}

// callUserFunc invokes a user-defined function with the given arguments.
func (e *Evaluator) callUserFunc(fn *UserFunction, args []types.Value) types.Value {
	// Validate argument count
	if len(args) != fn.Arity() {
		if fn.Arity() == 1 {
			return types.Errorf("%s requires exactly 1 argument, got %d", fn.Name, len(args))
		}
		return types.Errorf("%s requires exactly %d arguments, got %d", fn.Name, fn.Arity(), len(args))
	}

	// Build parameter bindings
	bindings := make(map[string]types.Value, len(fn.Params))
	for i, param := range fn.Params {
		bindings[param] = args[i]
		// Also store lowercase for case-insensitive lookup
		bindings[strings.ToLower(param)] = args[i]
	}

	// Evaluate body with bindings
	return e.evalWithBindings(fn.Body, bindings)
}

// evalWithBindings evaluates an expression with local variable bindings.
// Used for function parameter passing.
func (e *Evaluator) evalWithBindings(expr ast.Expr, bindings map[string]types.Value) types.Value {
	// Save current bindings
	prevBindings := e.localBindings

	// Set new bindings
	e.localBindings = bindings

	// Evaluate
	result := e.evalExpr(expr)

	// Restore previous bindings
	e.localBindings = prevBindings

	return result
}

// ════════════════════════════════════════════════════════════════
// EXPRESSION EVALUATION WITH BINDINGS (for external use)
// ════════════════════════════════════════════════════════════════

// EvalExprWithBindings evaluates an expression with the given variable bindings.
// This is useful for evaluating function bodies with parameter values.
func (e *Evaluator) EvalExprWithBindings(expr ast.Expr, bindings map[string]types.Value) types.Value {
	return e.evalWithBindings(expr, bindings)
}
