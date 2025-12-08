// internal/eval/eval.go

package eval

import (
	"math"
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
		return e.evalContinuation(ex)

	case *ast.ConversionContinuation:
		return e.evalConversionContinuation(ex)

	default:
		return types.Error("unknown expression type")
	}
}

// ════════════════════════════════════════════════════════════════
// IDENTIFIER EVALUATION
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalIdentifier(id *ast.Identifier) types.Value {
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

	return e.applyBinaryOp(expr.Op, left, right)
}

func (e *Evaluator) applyBinaryOp(op ast.BinaryOp, left, right types.Value) types.Value {
	// Handle percentage operations specially
	if right.IsPercentage() && (op == ast.OpAdd || op == ast.OpSub) {
		return e.applyPercentageOp(op, left, right)
	}

	// Get numeric values
	leftNum := left.AsFloat()
	rightNum := right.AsFloat()

	// For percentages used in multiplication/division, use decimal value
	if right.IsPercentage() && (op == ast.OpMul || op == ast.OpDiv) {
		rightNum = right.Num
	}
	if left.IsPercentage() && (op == ast.OpMul || op == ast.OpDiv) {
		leftNum = left.Num
	}

	var result float64

	switch op {
	case ast.OpAdd:
		result = leftNum + rightNum
	case ast.OpSub:
		result = leftNum - rightNum
	case ast.OpMul:
		result = leftNum * rightNum
	case ast.OpDiv:
		if rightNum == 0 {
			return types.Error("division by zero")
		}
		result = leftNum / rightNum
	case ast.OpPow:
		result = math.Pow(leftNum, rightNum)
	case ast.OpMod:
		if rightNum == 0 {
			return types.Error("modulo by zero")
		}
		result = math.Mod(leftNum, rightNum)
	default:
		return types.Error("unknown operator")
	}

	// Determine result type based on operands
	return e.coerceResult(result, left, right, op)
}

// applyPercentageOp handles "value + percentage" and "value - percentage"
// e.g., 100 + 15% = 115, $50 - 10% = $45
func (e *Evaluator) applyPercentageOp(op ast.BinaryOp, left, right types.Value) types.Value {
	baseValue := left.AsFloat()
	percentage := right.Num // Already in decimal form (0.15 for 15%)

	var result float64
	if op == ast.OpAdd {
		result = baseValue * (1 + percentage)
	} else { // OpSub
		result = baseValue * (1 - percentage)
	}

	// Preserve the left operand's type
	return left.WithAmount(result)
}

// coerceResult determines the result type based on operands.
func (e *Evaluator) coerceResult(result float64, left, right types.Value, op ast.BinaryOp) types.Value {
	// For multiplication/division, special handling
	if op == ast.OpMul || op == ast.OpDiv {
		// If one is a plain number, inherit the other's type
		if left.IsNumber() && !right.IsNumber() {
			return right.WithAmount(result)
		}
		if right.IsNumber() && !left.IsNumber() {
			return left.WithAmount(result)
		}
		// Both typed - return plain number (or could be unit algebra in future)
		if !left.IsNumber() && !right.IsNumber() {
			return types.Number(result)
		}
	}

	// For addition/subtraction, types must be compatible
	if op == ast.OpAdd || op == ast.OpSub {
		// Same type - preserve it
		if left.Kind == right.Kind {
			return left.WithAmount(result)
		}

		// One is a plain number - inherit the typed one
		if left.IsNumber() {
			return right.WithAmount(result)
		}
		if right.IsNumber() {
			return left.WithAmount(result)
		}

		// Different typed values - need conversion
		// For currencies, convert right to left's currency
		if left.IsCurrency() && right.IsCurrency() {
			if left.Curr != nil && right.Curr != nil {
				converted, ok := e.ctx.Convert(right.Num, right.Curr.Code, left.Curr.Code)
				if ok {
					if op == ast.OpAdd {
						return left.WithAmount(left.Num + converted)
					}
					return left.WithAmount(left.Num - converted)
				}
			}
		}

		// For units, convert right to left's unit
		if left.IsUnit() && right.IsUnit() {
			if left.Unit != nil && right.Unit != nil {
				converted, ok := right.Unit.ConvertTo(right.Num, left.Unit)
				if ok {
					if op == ast.OpAdd {
						return left.WithAmount(left.Num + converted)
					}
					return left.WithAmount(left.Num - converted)
				}
			}
			return types.Error("incompatible units")
		}
	}

	return types.Number(result)
}

// ════════════════════════════════════════════════════════════════
// UNARY OPERATIONS
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) evalUnary(expr *ast.UnaryExpr) types.Value {
	value := e.evalExpr(expr.Expr)
	if value.IsError() {
		return value
	}

	switch expr.Op {
	case ast.OpNeg:
		return value.Negate()
	case ast.OpPos:
		return value
	default:
		return types.Error("unknown unary operator")
	}
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

	return e.convertValue(value, expr.Target)
}

func (e *Evaluator) convertValue(value types.Value, target string) types.Value {
	// Try unit conversion first
	if value.IsUnit() && value.Unit != nil {
		targetUnit := types.ParseUnit(target)
		if targetUnit != nil {
			converted, ok := value.Unit.ConvertTo(value.Num, targetUnit)
			if ok {
				return types.UnitValue(converted, targetUnit)
			}
			return types.Errorf("cannot convert %s to %s", value.Unit.Code, target)
		}
	}

	// Try currency/crypto conversion
	converted, ok := e.ctx.ConvertValue(value, target)
	if ok {
		return converted
	}

	// Check if target is valid but conversion unavailable
	if types.ParseCurrency(target) != nil || types.ParseCrypto(target) != nil {
		return types.Errorf("no rate available for conversion to %s", target)
	}
	if types.ParseUnit(target) != nil {
		return types.Errorf("cannot convert to %s (incompatible types)", target)
	}

	return types.Errorf("unknown target: %s", target)
}

// ════════════════════════════════════════════════════════════════
// CONTINUATIONS
// ════════════════════════════════════════════════════════════════

// evalContinuation handles "+ 10", "* 2" etc. continuing from previous.
func (e *Evaluator) evalContinuation(expr *ast.ContinuationExpr) types.Value {
	if !e.ctx.HasPrevious() {
		// No previous - evaluate expression alone
		return e.evalExpr(expr.Expr)
	}

	prev := e.ctx.Previous()
	right := e.evalExpr(expr.Expr)
	if right.IsError() {
		return right
	}

	return e.applyBinaryOp(expr.Op, prev, right)
}

// evalConversionContinuation handles "in EUR", "to miles" continuing from previous.
func (e *Evaluator) evalConversionContinuation(expr *ast.ConversionContinuation) types.Value {
	if !e.ctx.HasPrevious() {
		return types.Error("no previous value to convert")
	}

	prev := e.ctx.Previous()
	return e.convertValue(prev, expr.Target)
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
	return e.callFunction(name, args)
}

func (e *Evaluator) callFunction(name string, args []types.Value) types.Value {
	switch name {
	// Aggregation functions
	case "sum":
		return e.fnSum(args)
	case "avg", "average", "mean":
		return e.fnAvg(args)
	case "min":
		return e.fnMin(args)
	case "max":
		return e.fnMax(args)
	case "count":
		return types.Number(float64(len(args)))

	// Math functions
	case "abs":
		return e.fnUnary(args, math.Abs)
	case "sqrt":
		return e.fnUnary(args, math.Sqrt)
	case "round":
		return e.fnUnary(args, math.Round)
	case "floor":
		return e.fnUnary(args, math.Floor)
	case "ceil":
		return e.fnUnary(args, math.Ceil)
	case "log", "log10":
		return e.fnUnary(args, math.Log10)
	case "ln":
		return e.fnUnary(args, math.Log)
	case "exp":
		return e.fnUnary(args, math.Exp)
	case "sin":
		return e.fnUnary(args, math.Sin)
	case "cos":
		return e.fnUnary(args, math.Cos)
	case "tan":
		return e.fnUnary(args, math.Tan)
	case "asin":
		return e.fnUnary(args, math.Asin)
	case "acos":
		return e.fnUnary(args, math.Acos)
	case "atan":
		return e.fnUnary(args, math.Atan)

	// Power function (2 args)
	case "pow":
		return e.fnPow(args)

	default:
		return types.Errorf("unknown function: %s", name)
	}
}

// ════════════════════════════════════════════════════════════════
// BUILT-IN FUNCTIONS
// ════════════════════════════════════════════════════════════════

func (e *Evaluator) fnSum(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Number(0)
	}

	var total float64
	var resultType types.Value = args[0]

	for _, arg := range args {
		if arg.IsError() {
			return arg
		}
		total += arg.AsFloat()
	}

	return resultType.WithAmount(total)
}

func (e *Evaluator) fnAvg(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Number(0)
	}

	sum := e.fnSum(args)
	if sum.IsError() {
		return sum
	}

	return sum.WithAmount(sum.AsFloat() / float64(len(args)))
}

func (e *Evaluator) fnMin(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("min requires at least one argument")
	}

	minVal := args[0]
	minNum := minVal.AsFloat()

	for _, arg := range args[1:] {
		if arg.IsError() {
			return arg
		}
		if arg.AsFloat() < minNum {
			minNum = arg.AsFloat()
			minVal = arg
		}
	}

	return minVal.WithAmount(minNum)
}

func (e *Evaluator) fnMax(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("max requires at least one argument")
	}

	maxVal := args[0]
	maxNum := maxVal.AsFloat()

	for _, arg := range args[1:] {
		if arg.IsError() {
			return arg
		}
		if arg.AsFloat() > maxNum {
			maxNum = arg.AsFloat()
			maxVal = arg
		}
	}

	return maxVal.WithAmount(maxNum)
}

func (e *Evaluator) fnUnary(args []types.Value, fn func(float64) float64) types.Value {
	if len(args) != 1 {
		return types.Error("function requires exactly one argument")
	}

	arg := args[0]
	if arg.IsError() {
		return arg
	}

	result := fn(arg.AsFloat())

	// Check for NaN/Inf
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return types.Error("invalid result")
	}

	return types.Number(result)
}

func (e *Evaluator) fnPow(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("pow requires exactly two arguments")
	}

	base := args[0].AsFloat()
	exp := args[1].AsFloat()

	result := math.Pow(base, exp)

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return types.Error("invalid result")
	}

	return types.Number(result)
}