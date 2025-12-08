// internal/parser/parser.go

// Package parser provides a recursive descent parser for numio expressions.
package parser

import (
	"strconv"
	"strings"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/internal/errors"
	"github.com/0xsj/numio/internal/lexer"
	"github.com/0xsj/numio/internal/token"
	"github.com/0xsj/numio/internal/types"
)

// Parser parses tokens into an AST.
type Parser struct {
	lexer  *lexer.Lexer
	tokens []token.Token
	pos    int
	errors []*errors.Error
}

// New creates a new Parser for the given input.
func New(input string) *Parser {
	l := lexer.New(input)
	return &Parser{
		lexer:  l,
		tokens: l.Tokenize(),
		pos:    0,
		errors: nil,
	}
}

// NewFromTokens creates a parser from pre-tokenized input.
func NewFromTokens(tokens []token.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: nil,
	}
}

// Errors returns all parsing errors.
func (p *Parser) Errors() []*errors.Error {
	return p.errors
}

// HasErrors returns true if there were parsing errors.
func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

// ════════════════════════════════════════════════════════════════
// TOKEN NAVIGATION
// ════════════════════════════════════════════════════════════════

// current returns the current token.
func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.New(token.EOF, "", -1)
	}
	return p.tokens[p.pos]
}

// peek returns the next token without advancing.
func (p *Parser) peek() token.Token {
	if p.pos+1 >= len(p.tokens) {
		return token.New(token.EOF, "", -1)
	}
	return p.tokens[p.pos+1]
}

// advance moves to the next token and returns the previous one.
func (p *Parser) advance() token.Token {
	tok := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

// check returns true if the current token is of the given type.
func (p *Parser) check(t token.Type) bool {
	return p.current().Type == t
}

// checkAny returns true if the current token is any of the given types.
func (p *Parser) checkAny(types ...token.Type) bool {
	for _, t := range types {
		if p.current().Type == t {
			return true
		}
	}
	return false
}

// match consumes the current token if it matches, returns true if matched.
func (p *Parser) match(t token.Type) bool {
	if p.check(t) {
		p.advance()
		return true
	}
	return false
}

// expect consumes the current token if it matches, otherwise adds an error.
func (p *Parser) expect(t token.Type, msg string) bool {
	if p.check(t) {
		p.advance()
		return true
	}
	p.addError(msg)
	return false
}

// addError adds a parsing error.
func (p *Parser) addError(msg string) {
	err := errors.ParseError(msg).WithPos(p.current().Pos)
	p.errors = append(p.errors, err)
}

// addErrorf adds a formatted parsing error.
func (p *Parser) addErrorf(format string, args ...any) {
	err := errors.ParseErrorf(format, args...).WithPos(p.current().Pos)
	p.errors = append(p.errors, err)
}

// ════════════════════════════════════════════════════════════════
// PUBLIC PARSING METHODS
// ════════════════════════════════════════════════════════════════

// ParseLine parses a single line of input.
func (p *Parser) ParseLine() *ast.Line {
	// Skip leading newlines
	for p.match(token.NEWLINE) {
	}

	// Check for EOF
	if p.check(token.EOF) {
		return &ast.Line{Stmt: &ast.EmptyStmt{}}
	}

	// Check for comment-only line
	if p.check(token.COMMENT) {
		comment := p.advance()
		return &ast.Line{
			Stmt:    &ast.CommentStmt{Text: comment.Literal},
			Comment: comment.Literal,
		}
	}

	// Try to parse a statement
	stmt := p.parseStatement()

	// Check for trailing comment
	var comment string
	if p.check(token.COMMENT) {
		comment = p.advance().Literal
	}

	return &ast.Line{
		Stmt:    stmt,
		Comment: comment,
	}
}

// Parse parses the entire input and returns the first line.
// For multi-line input, use ParseLines.
func (p *Parser) Parse() *ast.Line {
	return p.ParseLine()
}

// ParseLines parses multiple lines of input.
func (p *Parser) ParseLines() []*ast.Line {
	var lines []*ast.Line

	for !p.check(token.EOF) {
		line := p.ParseLine()
		lines = append(lines, line)

		// Consume newline if present
		p.match(token.NEWLINE)
	}

	return lines
}

// ════════════════════════════════════════════════════════════════
// STATEMENT PARSING
// ════════════════════════════════════════════════════════════════

// parseStatement parses a statement (assignment or expression).
func (p *Parser) parseStatement() ast.Stmt {
	// Check for assignment: identifier = expr
	if p.check(token.IDENTIFIER) && p.peek().Type == token.EQUALS {
		return p.parseAssignment()
	}

	// Check for continuation (line starting with operator)
	if p.checkAny(token.PLUS, token.MINUS, token.STAR, token.SLASH, token.CARET, token.POWER) {
		return p.parseContinuation()
	}

	// Check for conversion continuation (line starting with "in" or "to")
	if p.check(token.IN) {
		return p.parseConversionContinuation()
	}

	// Otherwise, parse as expression
	expr := p.parseExpression()
	if expr == nil {
		return &ast.EmptyStmt{}
	}

	return &ast.ExprStmt{Expr: expr}
}

// parseAssignment parses a variable assignment.
func (p *Parser) parseAssignment() *ast.AssignStmt {
	name := p.advance().Literal // identifier
	p.advance()                 // =

	expr := p.parseExpression()
	if expr == nil {
		p.addError("expected expression after '='")
		return &ast.AssignStmt{Name: name, Expr: &ast.NumberLit{Value: 0}}
	}

	return &ast.AssignStmt{Name: name, Expr: expr}
}

// parseContinuation parses a continuation expression (e.g., "+ 10").
func (p *Parser) parseContinuation() ast.Stmt {
	op := p.parseBinaryOp()
	expr := p.parseExpression()

	if expr == nil {
		p.addError("expected expression after operator")
		return &ast.EmptyStmt{}
	}

	return &ast.ExprStmt{
		Expr: &ast.ContinuationExpr{Op: op, Expr: expr},
	}
}

// parseConversionContinuation parses "in X" or "to X" continuation.
func (p *Parser) parseConversionContinuation() ast.Stmt {
	p.advance() // consume "in" or "to"

	if !p.check(token.IDENTIFIER) {
		p.addError("expected unit or currency after 'in'/'to'")
		return &ast.EmptyStmt{}
	}

	target := p.advance().Literal

	return &ast.ExprStmt{
		Expr: &ast.ConversionContinuation{Target: target},
	}
}

// ════════════════════════════════════════════════════════════════
// EXPRESSION PARSING (Pratt parser / precedence climbing)
// ════════════════════════════════════════════════════════════════

// parseExpression parses an expression.
func (p *Parser) parseExpression() ast.Expr {
	return p.parseBinaryExpr(0)
}

// parseBinaryExpr parses binary expressions with precedence climbing.
func (p *Parser) parseBinaryExpr(minPrec int) ast.Expr {
	left := p.parseUnaryExpr()
	if left == nil {
		return nil
	}

	for {
		// Check for binary operator
		if !p.isBinaryOp() {
			break
		}

		op := p.currentBinaryOp()
		prec := op.Precedence()

		// Stop if precedence is too low
		if prec < minPrec {
			break
		}

		p.advance() // consume operator

		// Parse right side with higher precedence for left-associativity
		// Use prec+1 for left-associative, prec for right-associative
		rightPrec := prec + 1
		if op == ast.OpPow {
			rightPrec = prec // Power is right-associative
		}

		right := p.parseBinaryExpr(rightPrec)
		if right == nil {
			p.addError("expected expression after operator")
			return left
		}

		left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
	}

	// Check for conversion suffix: "in EUR", "to miles"
	if p.check(token.IN) {
		p.advance()
		if p.check(token.IDENTIFIER) {
			target := p.advance().Literal
			left = &ast.ConversionExpr{Value: left, Target: target}
		}
	}

	return left
}

// parseUnaryExpr parses unary expressions.
func (p *Parser) parseUnaryExpr() ast.Expr {
	// Unary minus or plus
	if p.checkAny(token.MINUS, token.PLUS) {
		tok := p.advance()
		expr := p.parseUnaryExpr()
		if expr == nil {
			return nil
		}

		var op ast.UnaryOp
		if tok.Type == token.MINUS {
			op = ast.OpNeg
		} else {
			op = ast.OpPos
		}

		return &ast.UnaryExpr{Op: op, Expr: expr}
	}

	return p.parsePostfixExpr()
}

// parsePostfixExpr parses postfix expressions (function calls, etc).
func (p *Parser) parsePostfixExpr() ast.Expr {
	expr := p.parsePrimaryExpr()
	if expr == nil {
		return nil
	}

	// Check for "of" (percent of): 20% of 150
	if p.check(token.OF) {
		// Only valid if expr is a percentage
		if _, ok := expr.(*ast.PercentLit); ok {
			p.advance() // consume "of"
			value := p.parseUnaryExpr()
			if value == nil {
				p.addError("expected expression after 'of'")
				return expr
			}
			return &ast.PercentOfExpr{Percent: expr, Value: value}
		}
	}

	return expr
}

// parsePrimaryExpr parses primary expressions (literals, identifiers, groups).
func (p *Parser) parsePrimaryExpr() ast.Expr {
	tok := p.current()

	switch tok.Type {
	case token.NUMBER:
		return p.parseNumber()

	case token.PERCENT:
		return p.parsePercent()

	case token.DOLLAR, token.EURO, token.POUND, token.YEN, token.BITCOIN, token.CURRENCY:
		return p.parseCurrencyWithSymbol()

	case token.IDENTIFIER:
		return p.parseIdentifierOrValue()

	case token.LPAREN:
		return p.parseGroupExpr()

	case token.EOF, token.NEWLINE, token.COMMENT:
		return nil

	default:
		// Don't error on valid statement terminators
		if tok.Type != token.RPAREN && tok.Type != token.COMMA {
			p.addErrorf("unexpected token: %s", tok.Literal)
		}
		return nil
	}
}

// ════════════════════════════════════════════════════════════════
// LITERAL PARSING
// ════════════════════════════════════════════════════════════════

// parseNumber parses a numeric literal, possibly followed by a unit/currency.
func (p *Parser) parseNumber() ast.Expr {
	tok := p.advance()
	value, err := parseFloat(tok.Literal)
	if err != nil {
		p.addErrorf("invalid number: %s", tok.Literal)
		return &ast.NumberLit{Value: 0, Raw: tok.Literal}
	}

	// Check for unit or currency suffix
	if p.check(token.IDENTIFIER) {
		suffix := p.current().Literal

		// Try currency
		if curr := types.ParseCurrency(suffix); curr != nil {
			p.advance()
			return &ast.CurrencyLit{Amount: value, Currency: curr, Raw: tok.Literal + " " + suffix}
		}

		// Try crypto
		if crypto := types.ParseCrypto(suffix); crypto != nil {
			p.advance()
			return &ast.CryptoLit{Amount: value, Crypto: crypto, Raw: tok.Literal + " " + suffix}
		}

		// Try metal
		if metal := types.ParseMetal(suffix); metal != nil {
			p.advance()
			return &ast.MetalLit{Amount: value, Metal: metal, Raw: tok.Literal + " " + suffix}
		}

		// Try unit
		if unit := types.ParseUnit(suffix); unit != nil {
			p.advance()
			return &ast.UnitLit{Amount: value, Unit: unit, Raw: tok.Literal + " " + suffix}
		}
	}

	return &ast.NumberLit{Value: value, Raw: tok.Literal}
}

// parsePercent parses a percentage literal (e.g., "20%").
func (p *Parser) parsePercent() ast.Expr {
	tok := p.advance()
	// Remove % suffix for parsing
	numStr := strings.TrimSuffix(tok.Literal, "%")
	value, err := parseFloat(numStr)
	if err != nil {
		p.addErrorf("invalid percentage: %s", tok.Literal)
		return &ast.PercentLit{Value: 0, Raw: tok.Literal}
	}

	return &ast.PercentLit{Value: value / 100.0, Raw: tok.Literal}
}

// parseCurrencyWithSymbol parses currency with leading symbol (e.g., "$100").
func (p *Parser) parseCurrencyWithSymbol() ast.Expr {
	symbolTok := p.advance()
	symbol := symbolTok.Literal

	// Look up currency by symbol
	var curr *types.Currency
	var crypto *types.Crypto

	// Check fiat currency
	curr = types.LookupCurrencyBySymbol(symbol)

	// Check crypto
	if curr == nil {
		crypto = types.LookupCrypto(symbol)
	}

	// Expect a number to follow
	if !p.check(token.NUMBER) {
		p.addError("expected number after currency symbol")
		if curr != nil {
			return &ast.CurrencyLit{Amount: 0, Currency: curr, Raw: symbol}
		}
		if crypto != nil {
			return &ast.CryptoLit{Amount: 0, Crypto: crypto, Raw: symbol}
		}
		return &ast.NumberLit{Value: 0, Raw: symbol}
	}

	numTok := p.advance()
	amount, err := parseFloat(numTok.Literal)
	if err != nil {
		p.addErrorf("invalid number: %s", numTok.Literal)
		amount = 0
	}

	raw := symbol + numTok.Literal

	if curr != nil {
		return &ast.CurrencyLit{Amount: amount, Currency: curr, Raw: raw}
	}
	if crypto != nil {
		return &ast.CryptoLit{Amount: amount, Crypto: crypto, Raw: raw}
	}

	// Unknown symbol, treat as number
	return &ast.NumberLit{Value: amount, Raw: raw}
}

// parseIdentifierOrValue parses an identifier, which could be:
// - A variable reference
// - A function call
// - A currency/unit name (e.g., "dollars", "kilometers")
func (p *Parser) parseIdentifierOrValue() ast.Expr {
	tok := p.advance()
	name := tok.Literal

	// Check for function call: name(args)
	if p.check(token.LPAREN) {
		return p.parseFunctionCall(name)
	}

	// Check if it's a currency name (e.g., "dollars")
	if curr := types.ParseCurrency(name); curr != nil {
		// It's a currency name used alone - could be for conversion target
		// But as a primary expression, treat as identifier
		// The conversion will be handled at a higher level
	}

	// Check for special identifiers
	lower := strings.ToLower(name)
	if lower == "_" || lower == "ans" {
		return &ast.Identifier{Name: "_"} // Normalize to _
	}

	return &ast.Identifier{Name: name}
}

// parseFunctionCall parses a function call.
func (p *Parser) parseFunctionCall(name string) ast.Expr {
	p.advance() // consume (

	var args []ast.Expr

	// Parse arguments
	if !p.check(token.RPAREN) {
		for {
			arg := p.parseExpression()
			if arg != nil {
				args = append(args, arg)
			}

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	p.expect(token.RPAREN, "expected ')' after function arguments")

	return &ast.CallExpr{Name: name, Args: args}
}

// parseGroupExpr parses a parenthesized expression.
func (p *Parser) parseGroupExpr() ast.Expr {
	p.advance() // consume (

	expr := p.parseExpression()
	if expr == nil {
		p.addError("expected expression inside parentheses")
		expr = &ast.NumberLit{Value: 0}
	}

	p.expect(token.RPAREN, "expected ')' after expression")

	return &ast.GroupExpr{Expr: expr}
}

// ════════════════════════════════════════════════════════════════
// OPERATOR HELPERS
// ════════════════════════════════════════════════════════════════

// isBinaryOp returns true if current token is a binary operator.
func (p *Parser) isBinaryOp() bool {
	return p.checkAny(token.PLUS, token.MINUS, token.STAR, token.SLASH, token.CARET, token.POWER)
}

// currentBinaryOp returns the current token as a BinaryOp.
func (p *Parser) currentBinaryOp() ast.BinaryOp {
	switch p.current().Type {
	case token.PLUS:
		return ast.OpAdd
	case token.MINUS:
		return ast.OpSub
	case token.STAR:
		return ast.OpMul
	case token.SLASH:
		return ast.OpDiv
	case token.CARET, token.POWER:
		return ast.OpPow
	default:
		return ast.OpAdd
	}
}

// parseBinaryOp consumes and returns the current binary operator.
func (p *Parser) parseBinaryOp() ast.BinaryOp {
	op := p.currentBinaryOp()
	p.advance()
	return op
}

// ════════════════════════════════════════════════════════════════
// NUMBER PARSING HELPERS
// ════════════════════════════════════════════════════════════════

// parseFloat parses a float from string, handling thousands separators.
func parseFloat(s string) (float64, error) {
	// Remove thousands separators
	s = strings.ReplaceAll(s, ",", "")
	return strconv.ParseFloat(s, 64)
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// ParseLine parses a single line of input.
func ParseLine(input string) (*ast.Line, []*errors.Error) {
	p := New(input)
	line := p.ParseLine()
	return line, p.Errors()
}

// ParseExpr parses a single expression.
func ParseExpr(input string) (ast.Expr, []*errors.Error) {
	p := New(input)
	expr := p.parseExpression()
	return expr, p.Errors()
}

// MustParseLine parses a line, panicking on error (for tests).
func MustParseLine(input string) *ast.Line {
	line, errs := ParseLine(input)
	if len(errs) > 0 {
		panic("parse error: " + errs[0].Error())
	}
	return line
}

// MustParseExpr parses an expression, panicking on error (for tests).
func MustParseExpr(input string) ast.Expr {
	expr, errs := ParseExpr(input)
	if len(errs) > 0 {
		panic("parse error: " + errs[0].Error())
	}
	return expr
}