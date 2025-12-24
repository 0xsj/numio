// internal/ast/ast.go

// Package ast defines the Abstract Syntax Tree for numio expressions.
package ast

import (
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// Node is the interface implemented by all AST nodes.
type Node interface {
	node()
	String() string
}

// Expr is the interface implemented by all expression nodes.
type Expr interface {
	Node
	expr()
}

// Stmt is the interface implemented by all statement nodes.
type Stmt interface {
	Node
	stmt()
}

// ════════════════════════════════════════════════════════════════
// STATEMENTS
// ════════════════════════════════════════════════════════════════

// Line represents a single line of input.
// It can be empty, a comment, an assignment, or an expression.
type Line struct {
	Stmt    Stmt   // The statement (nil if empty)
	Comment string // Trailing comment (if any)
	Raw     string // Original raw input
}

func (l *Line) node() {}
func (l *Line) stmt() {}

func (l *Line) String() string {
	if l.Stmt == nil {
		if l.Comment != "" {
			return l.Comment
		}
		return ""
	}
	if l.Comment != "" {
		return l.Stmt.String() + " " + l.Comment
	}
	return l.Stmt.String()
}

// EmptyStmt represents an empty line.
type EmptyStmt struct{}

func (e *EmptyStmt) node() {}
func (e *EmptyStmt) stmt() {}

func (e *EmptyStmt) String() string {
	return ""
}

// CommentStmt represents a comment-only line.
type CommentStmt struct {
	Text string // Full comment text including # or //
}

func (c *CommentStmt) node() {}
func (c *CommentStmt) stmt() {}

func (c *CommentStmt) String() string {
	return c.Text
}

// ExprStmt represents an expression statement.
type ExprStmt struct {
	Expr Expr
}

func (e *ExprStmt) node() {}
func (e *ExprStmt) stmt() {}

func (e *ExprStmt) String() string {
	if e.Expr == nil {
		return ""
	}
	return e.Expr.String()
}

// AssignStmt represents a variable assignment.
type AssignStmt struct {
	Name string // Variable name
	Expr Expr   // Value expression
}

func (a *AssignStmt) node() {}
func (a *AssignStmt) stmt() {}

func (a *AssignStmt) String() string {
	return a.Name + " = " + a.Expr.String()
}

// ════════════════════════════════════════════════════════════════
// EXPRESSIONS - LITERALS
// ════════════════════════════════════════════════════════════════

// NumberLit represents a numeric literal.
type NumberLit struct {
	Value float64
	Raw   string // Original text (for display)
}

func (n *NumberLit) node() {}
func (n *NumberLit) expr() {}

func (n *NumberLit) String() string {
	if n.Raw != "" {
		return n.Raw
	}
	return formatFloat(n.Value)
}

// PercentLit represents a percentage literal (e.g., 20%).
type PercentLit struct {
	Value float64 // Stored as decimal (0.20 for 20%)
	Raw   string  // Original text
}

func (p *PercentLit) node() {}
func (p *PercentLit) expr() {}

func (p *PercentLit) String() string {
	if p.Raw != "" {
		return p.Raw
	}
	return formatFloat(p.Value*100) + "%"
}

// CurrencyLit represents a currency literal (e.g., $100, €50).
type CurrencyLit struct {
	Amount   float64
	Currency *types.Currency
	Raw      string
}

func (c *CurrencyLit) node() {}
func (c *CurrencyLit) expr() {}

func (c *CurrencyLit) String() string {
	if c.Raw != "" {
		return c.Raw
	}
	if c.Currency != nil {
		if c.Currency.SymbolAfter {
			return formatFloat(c.Amount) + c.Currency.Symbol
		}
		return c.Currency.Symbol + formatFloat(c.Amount)
	}
	return formatFloat(c.Amount)
}

// UnitLit represents a value with a unit (e.g., 5 km, 2 hours).
type UnitLit struct {
	Amount float64
	Unit   *types.Unit
	Raw    string
}

func (u *UnitLit) node() {}
func (u *UnitLit) expr() {}

func (u *UnitLit) String() string {
	if u.Raw != "" {
		return u.Raw
	}
	if u.Unit != nil {
		return formatFloat(u.Amount) + " " + u.Unit.Code
	}
	return formatFloat(u.Amount)
}

// MetalLit represents a precious metal literal (e.g., 1 oz gold).
type MetalLit struct {
	Amount float64
	Metal  *types.Metal
	Raw    string
}

func (m *MetalLit) node() {}
func (m *MetalLit) expr() {}

func (m *MetalLit) String() string {
	if m.Raw != "" {
		return m.Raw
	}
	if m.Metal != nil {
		return formatFloat(m.Amount) + " " + m.Metal.Code
	}
	return formatFloat(m.Amount)
}

// CryptoLit represents a cryptocurrency literal (e.g., 0.5 BTC).
type CryptoLit struct {
	Amount float64
	Crypto *types.Crypto
	Raw    string
}

func (c *CryptoLit) node() {}
func (c *CryptoLit) expr() {}

func (c *CryptoLit) String() string {
	if c.Raw != "" {
		return c.Raw
	}
	if c.Crypto != nil {
		if c.Crypto.HasSymbol() {
			return c.Crypto.Symbol + formatFloat(c.Amount)
		}
		return formatFloat(c.Amount) + " " + c.Crypto.Code
	}
	return formatFloat(c.Amount)
}

// ════════════════════════════════════════════════════════════════
// EXPRESSIONS - REFERENCES
// ════════════════════════════════════════════════════════════════

// Identifier represents a variable reference.
type Identifier struct {
	Name string
}

func (i *Identifier) node() {}
func (i *Identifier) expr() {}

func (i *Identifier) String() string {
	return i.Name
}

// ════════════════════════════════════════════════════════════════
// EXPRESSIONS - OPERATORS
// ════════════════════════════════════════════════════════════════

// BinaryOp represents the binary operator type.
type BinaryOp int

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpPow
	OpMod
)

// String returns the operator symbol.
func (op BinaryOp) String() string {
	switch op {
	case OpAdd:
		return "+"
	case OpSub:
		return "-"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpPow:
		return "^"
	case OpMod:
		return "%"
	default:
		return "?"
	}
}

// Precedence returns the operator precedence (higher = tighter binding).
func (op BinaryOp) Precedence() int {
	switch op {
	case OpAdd, OpSub:
		return 1
	case OpMul, OpDiv, OpMod:
		return 2
	case OpPow:
		return 3
	default:
		return 0
	}
}

// BinaryExpr represents a binary operation (e.g., a + b).
type BinaryExpr struct {
	Left  Expr
	Op    BinaryOp
	Right Expr
}

func (b *BinaryExpr) node() {}
func (b *BinaryExpr) expr() {}

func (b *BinaryExpr) String() string {
	return "(" + b.Left.String() + " " + b.Op.String() + " " + b.Right.String() + ")"
}

// UnaryOp represents the unary operator type.
type UnaryOp int

const (
	OpNeg UnaryOp = iota // -x
	OpPos                // +x (no-op, but valid syntax)
)

// String returns the operator symbol.
func (op UnaryOp) String() string {
	switch op {
	case OpNeg:
		return "-"
	case OpPos:
		return "+"
	default:
		return "?"
	}
}

// UnaryExpr represents a unary operation (e.g., -x).
type UnaryExpr struct {
	Op   UnaryOp
	Expr Expr
}

func (u *UnaryExpr) node() {}
func (u *UnaryExpr) expr() {}

func (u *UnaryExpr) String() string {
	return "(" + u.Op.String() + u.Expr.String() + ")"
}

// ════════════════════════════════════════════════════════════════
// EXPRESSIONS - SPECIAL FORMS
// ════════════════════════════════════════════════════════════════

// PercentOfExpr represents "X% of Y" (e.g., 20% of 150).
type PercentOfExpr struct {
	Percent Expr // The percentage (should evaluate to percentage)
	Value   Expr // The value to take percentage of
}

func (p *PercentOfExpr) node() {}
func (p *PercentOfExpr) expr() {}

func (p *PercentOfExpr) String() string {
	return p.Percent.String() + " of " + p.Value.String()
}

// ConversionExpr represents a unit/currency conversion (e.g., $100 in EUR, 5 km to miles).
type ConversionExpr struct {
	Value  Expr   // The value to convert
	Target string // Target unit/currency (raw string, resolved at eval time)
}

func (c *ConversionExpr) node() {}
func (c *ConversionExpr) expr() {}

func (c *ConversionExpr) String() string {
	return c.Value.String() + " in " + c.Target
}

// CallExpr represents a function call (e.g., sum(1, 2, 3), sqrt(16)).
type CallExpr struct {
	Name string
	Args []Expr
}

func (c *CallExpr) node() {}
func (c *CallExpr) expr() {}

func (c *CallExpr) String() string {
	var args []string
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return c.Name + "(" + strings.Join(args, ", ") + ")"
}

// GroupExpr represents a parenthesized expression.
type GroupExpr struct {
	Expr Expr
}

func (g *GroupExpr) node() {}
func (g *GroupExpr) expr() {}

func (g *GroupExpr) String() string {
	return "(" + g.Expr.String() + ")"
}

// ════════════════════════════════════════════════════════════════
// EXPRESSIONS - CONTINUATION (for line chaining)
// ════════════════════════════════════════════════════════════════

// ContinuationExpr represents a continuation from previous result.
// e.g., "+ 10" continues from previous line's result.
type ContinuationExpr struct {
	Op   BinaryOp // The operator
	Expr Expr     // The expression to apply
}

func (c *ContinuationExpr) node() {}
func (c *ContinuationExpr) expr() {}

func (c *ContinuationExpr) String() string {
	return c.Op.String() + " " + c.Expr.String()
}

// ConversionContinuation represents "in X" or "to X" continuing from previous.
// e.g., "in EUR" converts previous result to EUR.
type ConversionContinuation struct {
	Target string
}

func (c *ConversionContinuation) node() {}
func (c *ConversionContinuation) expr() {}

func (c *ConversionContinuation) String() string {
	return "in " + c.Target
}

// ════════════════════════════════════════════════════════════════
// FUTURE: CONDITIONALS (placeholder for later)
// ════════════════════════════════════════════════════════════════

// CondExpr represents a conditional expression.
// e.g., "if x > 100 then 100 else x" or "x > 0 ? x : -x"
type CondExpr struct {
	Cond Expr // Condition
	Then Expr // True branch
	Else Expr // False branch
}

func (c *CondExpr) node() {}
func (c *CondExpr) expr() {}

func (c *CondExpr) String() string {
	return "if " + c.Cond.String() + " then " + c.Then.String() + " else " + c.Else.String()
}

// ════════════════════════════════════════════════════════════════
// FUTURE: RANGES (placeholder for later)
// ════════════════════════════════════════════════════════════════

// RangeExpr represents a range expression (e.g., 1..10).
type RangeExpr struct {
	Start Expr
	End   Expr
	Step  Expr // Optional step (nil = 1)
}

func (r *RangeExpr) node() {}
func (r *RangeExpr) expr() {}

func (r *RangeExpr) String() string {
	s := r.Start.String() + ".." + r.End.String()
	if r.Step != nil {
		s += " step " + r.Step.String()
	}
	return s
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// formatFloat formats a float64 for display without fmt package.
func formatFloat(v float64) string {
	// Handle negative
	if v < 0 {
		return "-" + formatFloat(-v)
	}

	// Handle zero
	if v == 0 {
		return "0"
	}

	// Check if it's a whole number
	if v == float64(int64(v)) {
		return itoa(int64(v))
	}

	// Format with up to 6 decimal places
	intPart := int64(v)
	fracPart := v - float64(intPart)

	// Multiply to get decimal digits
	const decimals = 6
	multiplier := 1000000.0
	frac := int64(fracPart*multiplier + 0.5)

	// Build fractional string
	fracStr := itoa(frac)
	for len(fracStr) < decimals {
		fracStr = "0" + fracStr
	}

	// Trim trailing zeros
	fracStr = strings.TrimRight(fracStr, "0")

	if fracStr == "" {
		return itoa(intPart)
	}

	return itoa(intPart) + "." + fracStr
}

// itoa converts int64 to string.
func itoa(n int64) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var buf [20]byte
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}

	if negative {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}

// ════════════════════════════════════════════════════════════════
// AST VISITOR (for traversal)
// ════════════════════════════════════════════════════════════════

// Visitor defines the interface for AST visitors.
type Visitor interface {
	Visit(node Node) Visitor
}

// Walk traverses an AST in depth-first order.
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Line:
		if n.Stmt != nil {
			Walk(v, n.Stmt)
		}

	case *ExprStmt:
		if n.Expr != nil {
			Walk(v, n.Expr)
		}

	case *AssignStmt:
		Walk(v, n.Expr)

	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *UnaryExpr:
		Walk(v, n.Expr)

	case *PercentOfExpr:
		Walk(v, n.Percent)
		Walk(v, n.Value)

	case *ConversionExpr:
		Walk(v, n.Value)

	case *CallExpr:
		for _, arg := range n.Args {
			Walk(v, arg)
		}

	case *GroupExpr:
		Walk(v, n.Expr)

	case *ContinuationExpr:
		Walk(v, n.Expr)

	case *CondExpr:
		Walk(v, n.Cond)
		Walk(v, n.Then)
		Walk(v, n.Else)

	case *RangeExpr:
		Walk(v, n.Start)
		Walk(v, n.End)
		if n.Step != nil {
			Walk(v, n.Step)
		}
	}
}

// ════════════════════════════════════════════════════════════════
// AST INSPECTION HELPERS
// ════════════════════════════════════════════════════════════════

// IsLiteral returns true if the expression is a literal value.
func IsLiteral(e Expr) bool {
	switch e.(type) {
	case *NumberLit, *PercentLit, *CurrencyLit, *UnitLit, *MetalLit, *CryptoLit:
		return true
	default:
		return false
	}
}

// IsContinuation returns true if the expression is a continuation.
func IsContinuation(e Expr) bool {
	switch e.(type) {
	case *ContinuationExpr, *ConversionContinuation:
		return true
	default:
		return false
	}
}

// GetIdentifiers returns all identifiers referenced in an expression.
func GetIdentifiers(e Expr) []string {
	var ids []string
	var collect func(Expr)

	collect = func(expr Expr) {
		switch n := expr.(type) {
		case *Identifier:
			ids = append(ids, n.Name)
		case *BinaryExpr:
			collect(n.Left)
			collect(n.Right)
		case *UnaryExpr:
			collect(n.Expr)
		case *PercentOfExpr:
			collect(n.Percent)
			collect(n.Value)
		case *ConversionExpr:
			collect(n.Value)
		case *CallExpr:
			for _, arg := range n.Args {
				collect(arg)
			}
		case *GroupExpr:
			collect(n.Expr)
		case *ContinuationExpr:
			collect(n.Expr)
		case *CondExpr:
			collect(n.Cond)
			collect(n.Then)
			collect(n.Else)
		case *RangeExpr:
			collect(n.Start)
			collect(n.End)
			if n.Step != nil {
				collect(n.Step)
			}
		}
	}

	collect(e)
	return ids
}
