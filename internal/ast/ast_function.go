// internal/ast/ast_function.go

package ast

import (
	"strings"
)

// ════════════════════════════════════════════════════════════════
// FUNCTION DEFINITION STATEMENT
// ════════════════════════════════════════════════════════════════

// FuncDefStmt represents a user-defined function definition.
// Syntax: def name(param1, param2, ...): body
//
// Examples:
//   - def withTax(amount): amount * 1.08
//   - def double(x): x * 2
//   - def toEUR(amount): amount in EUR
type FuncDefStmt struct {
	Name   string   // Function name
	Params []string // Parameter names
	Body   Expr     // Body expression
	Raw    string   // Original source text (for persistence/display)
}

func (f *FuncDefStmt) node() {}
func (f *FuncDefStmt) stmt() {}

func (f *FuncDefStmt) String() string {
	var sb strings.Builder
	sb.WriteString("def ")
	sb.WriteString(f.Name)
	sb.WriteString("(")
	sb.WriteString(strings.Join(f.Params, ", "))
	sb.WriteString("): ")
	if f.Body != nil {
		sb.WriteString(f.Body.String())
	}
	return sb.String()
}

// Arity returns the number of parameters.
func (f *FuncDefStmt) Arity() int {
	return len(f.Params)
}

// HasParams returns true if the function has parameters.
func (f *FuncDefStmt) HasParams() bool {
	return len(f.Params) > 0
}

// ════════════════════════════════════════════════════════════════
// FUNCTION DEFINITION HELPERS
// ════════════════════════════════════════════════════════════════

// ValidateFuncName checks if a function name is valid.
// Returns an error message if invalid, empty string if valid.
func ValidateFuncName(name string) string {
	if name == "" {
		return "function name cannot be empty"
	}

	// Check first character is letter or underscore
	first := rune(name[0])
	if !isValidIdentifierStart(first) {
		return "function name must start with a letter or underscore"
	}

	// Check remaining characters
	for i, ch := range name {
		if i == 0 {
			continue
		}
		if !isValidIdentifierChar(ch) {
			return "function name contains invalid character"
		}
	}

	// Check for reserved names
	reserved := map[string]bool{
		"def": true, "in": true, "to": true, "of": true,
		"true": true, "false": true, "nil": true, "null": true,
		"if": true, "then": true, "else": true, // Reserved for future
	}
	if reserved[strings.ToLower(name)] {
		return "function name is a reserved keyword"
	}

	return ""
}

// ValidateParamName checks if a parameter name is valid.
// Returns an error message if invalid, empty string if valid.
func ValidateParamName(name string) string {
	if name == "" {
		return "parameter name cannot be empty"
	}

	// Check first character
	first := rune(name[0])
	if !isValidIdentifierStart(first) {
		return "parameter name must start with a letter or underscore"
	}

	// Check remaining characters
	for i, ch := range name {
		if i == 0 {
			continue
		}
		if !isValidIdentifierChar(ch) {
			return "parameter name contains invalid character"
		}
	}

	return ""
}

// ValidateParams checks a list of parameter names.
// Returns an error message if invalid, empty string if valid.
func ValidateParams(params []string) string {
	seen := make(map[string]bool)

	for _, param := range params {
		// Validate individual param
		if err := ValidateParamName(param); err != "" {
			return err
		}

		// Check for duplicates
		lower := strings.ToLower(param)
		if seen[lower] {
			return "duplicate parameter name: " + param
		}
		seen[lower] = true
	}

	return ""
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// isValidIdentifierStart checks if a rune can start an identifier.
func isValidIdentifierStart(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

// isValidIdentifierChar checks if a rune can be part of an identifier.
func isValidIdentifierChar(ch rune) bool {
	return isValidIdentifierStart(ch) ||
		(ch >= '0' && ch <= '9')
}
