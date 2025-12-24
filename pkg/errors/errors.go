// internal/errors/errors.go

// Package errors defines custom error types for numio.
package errors

import "fmt"

// Kind represents the category of error.
type Kind int

const (
	KindUnknown    Kind = iota // Unknown or uncategorized error
	KindParse                  // Parsing/lexing error
	KindEval                   // Evaluation error
	KindConversion             // Unit/currency conversion error
	KindDivision               // Division by zero
	KindVariable               // Undefined variable
	KindFunction               // Unknown function or bad arguments
	KindType                   // Type mismatch
)

func (k Kind) String() string {
	switch k {
	case KindParse:
		return "parse error"
	case KindEval:
		return "evaluation error"
	case KindConversion:
		return "conversion error"
	case KindDivision:
		return "division by zero"
	case KindVariable:
		return "undefined variable"
	case KindFunction:
		return "function error"
	case KindType:
		return "type error"
	default:
		return "unknown error"
	}
}

// Error represents a numio error with kind, message, and optional position.
type Error struct {
	Kind    Kind
	Message string
	Pos     int // Character position in input, -1 if not applicable
	Line    int // Line number, -1 if not applicable
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Line >= 0 && e.Pos >= 0 {
		return fmt.Sprintf("%s: %s (line %d, col %d)", e.Kind, e.Message, e.Line, e.Pos)
	}
	if e.Pos >= 0 {
		return fmt.Sprintf("%s: %s (col %d)", e.Kind, e.Message, e.Pos)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

// New creates a new Error with the given kind and message.
func New(kind Kind, message string) *Error {
	return &Error{
		Kind:    kind,
		Message: message,
		Pos:     -1,
		Line:    -1,
	}
}

// Newf creates a new Error with a formatted message.
func Newf(kind Kind, format string, args ...any) *Error {
	return &Error{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
		Pos:     -1,
		Line:    -1,
	}
}

// WithPos returns a copy of the error with position information.
func (e *Error) WithPos(pos int) *Error {
	return &Error{
		Kind:    e.Kind,
		Message: e.Message,
		Pos:     pos,
		Line:    e.Line,
	}
}

// WithLine returns a copy of the error with line information.
func (e *Error) WithLine(line int) *Error {
	return &Error{
		Kind:    e.Kind,
		Message: e.Message,
		Pos:     e.Pos,
		Line:    line,
	}
}

// Is checks if the error is of a specific kind.
func Is(err error, kind Kind) bool {
	if e, ok := err.(*Error); ok {
		return e.Kind == kind
	}
	return false
}

// --- Convenience constructors ---

// ParseError creates a parse error.
func ParseError(message string) *Error {
	return New(KindParse, message)
}

// ParseErrorf creates a parse error with formatting.
func ParseErrorf(format string, args ...any) *Error {
	return Newf(KindParse, format, args...)
}

// EvalError creates an evaluation error.
func EvalError(message string) *Error {
	return New(KindEval, message)
}

// EvalErrorf creates an evaluation error with formatting.
func EvalErrorf(format string, args ...any) *Error {
	return Newf(KindEval, format, args...)
}

// ConversionError creates a conversion error.
func ConversionError(message string) *Error {
	return New(KindConversion, message)
}

// ConversionErrorf creates a conversion error with formatting.
func ConversionErrorf(format string, args ...any) *Error {
	return Newf(KindConversion, format, args...)
}

// DivisionByZero creates a division by zero error.
func DivisionByZero() *Error {
	return New(KindDivision, "division by zero")
}

// UndefinedVariable creates an undefined variable error.
func UndefinedVariable(name string) *Error {
	return Newf(KindVariable, "undefined variable: %s", name)
}

// UnknownFunction creates an unknown function error.
func UnknownFunction(name string) *Error {
	return Newf(KindFunction, "unknown function: %s", name)
}

// TypeError creates a type mismatch error.
func TypeError(message string) *Error {
	return New(KindType, message)
}

// TypeErrorf creates a type mismatch error with formatting.
func TypeErrorf(format string, args ...any) *Error {
	return Newf(KindType, format, args...)
}
