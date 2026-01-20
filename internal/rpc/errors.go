package rpc

import "fmt"

// ════════════════════════════════════════════════════════════════
// PREDEFINED ERRORS
// ════════════════════════════════════════════════════════════════

// Standard JSON-RPC 2.0 errors.
var (
	ErrParseError    = NewError(ErrCodeParse, "Parse error")
	ErrInvalidReq    = NewError(ErrCodeInvalidRequest, "Invalid request")
	ErrNotFound      = NewError(ErrCodeMethodNotFound, "Method not found")
	ErrInvalidP      = NewError(ErrCodeInvalidParams, "Invalid params")
	ErrInternalError = NewError(ErrCodeInternal, "Internal error")
)

// Application-specific errors.
var (
	ErrNotInitialized = NewError(ErrCodeNotReady, "Core not initialized")
	ErrShutdown       = NewError(ErrCodeShuttingDown, "Core is shutting down")
)

// ════════════════════════════════════════════════════════════════
// ERROR BUILDERS
// ════════════════════════════════════════════════════════════════

// ParseError creates a parse error with detail.
func ParseError(detail string) *Error {
	return &Error{
		Code:    ErrCodeParse,
		Message: "Parse error",
		Data:    detail,
	}
}

// InvalidRequestError creates an invalid request error with detail.
func InvalidRequestError(detail string) *Error {
	return &Error{
		Code:    ErrCodeInvalidRequest,
		Message: "Invalid request",
		Data:    detail,
	}
}

// MethodNotFoundError creates a method not found error.
func MethodNotFoundError(method string) *Error {
	return &Error{
		Code:    ErrCodeMethodNotFound,
		Message: fmt.Sprintf("Method not found: %s", method),
		Data:    method,
	}
}

// InvalidParamsError creates an invalid params error with detail.
func InvalidParamsError(detail string) *Error {
	return &Error{
		Code:    ErrCodeInvalidParams,
		Message: "Invalid params",
		Data:    detail,
	}
}

// InternalError creates an internal error with detail.
func InternalError(detail string) *Error {
	return &Error{
		Code:    ErrCodeInternal,
		Message: "Internal error",
		Data:    detail,
	}
}

// EvaluationError creates an evaluation error.
func EvaluationError(detail string) *Error {
	return &Error{
		Code:    ErrCodeEvaluation,
		Message: "Evaluation error",
		Data:    detail,
	}
}

// NotReadyError creates a not ready error with optional detail.
func NotReadyError(detail string) *Error {
	msg := "Core not initialized"
	if detail != "" {
		msg = fmt.Sprintf("Core not initialized: %s", detail)
	}
	return &Error{
		Code:    ErrCodeNotReady,
		Message: msg,
		Data:    detail,
	}
}

// ShuttingDownError creates a shutting down error.
func ShuttingDownError() *Error {
	return &Error{
		Code:    ErrCodeShuttingDown,
		Message: "Core is shutting down",
	}
}

// ════════════════════════════════════════════════════════════════
// ERROR WRAPPING
// ════════════════════════════════════════════════════════════════

// WrapError wraps a Go error as an RPC internal error.
func WrapError(err error) *Error {
	if err == nil {
		return nil
	}

	// If already an RPC error, return as-is
	if rpcErr, ok := err.(*Error); ok {
		return rpcErr
	}

	return &Error{
		Code:    ErrCodeInternal,
		Message: "Internal error",
		Data:    err.Error(),
	}
}

// WrapErrorWithCode wraps a Go error with a specific code.
func WrapErrorWithCode(code int, err error) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}

// ════════════════════════════════════════════════════════════════
// ERROR CHECKING
// ════════════════════════════════════════════════════════════════

// IsParseError checks if an error is a parse error.
func IsParseError(err *Error) bool {
	return err != nil && err.Code == ErrCodeParse
}

// IsInvalidRequest checks if an error is an invalid request error.
func IsInvalidRequest(err *Error) bool {
	return err != nil && err.Code == ErrCodeInvalidRequest
}

// IsMethodNotFound checks if an error is a method not found error.
func IsMethodNotFound(err *Error) bool {
	return err != nil && err.Code == ErrCodeMethodNotFound
}

// IsInvalidParams checks if an error is an invalid params error.
func IsInvalidParams(err *Error) bool {
	return err != nil && err.Code == ErrCodeInvalidParams
}

// IsInternalError checks if an error is an internal error.
func IsInternalError(err *Error) bool {
	return err != nil && err.Code == ErrCodeInternal
}

// IsEvaluationError checks if an error is an evaluation error.
func IsEvaluationError(err *Error) bool {
	return err != nil && err.Code == ErrCodeEvaluation
}

// IsNotReady checks if an error is a not ready error.
func IsNotReady(err *Error) bool {
	return err != nil && err.Code == ErrCodeNotReady
}

// IsShuttingDown checks if an error is a shutting down error.
func IsShuttingDown(err *Error) bool {
	return err != nil && err.Code == ErrCodeShuttingDown
}

// ════════════════════════════════════════════════════════════════
// ERROR CATEGORIES
// ════════════════════════════════════════════════════════════════

// IsClientError returns true if the error is a client-side error.
// Client errors are in the range -32600 to -32699.
func IsClientError(err *Error) bool {
	if err == nil {
		return false
	}
	return err.Code >= -32699 && err.Code <= -32600
}

// IsServerError returns true if the error is a server-side error.
// Server errors are in the range -32000 to -32099.
func IsServerError(err *Error) bool {
	if err == nil {
		return false
	}
	return err.Code >= -32099 && err.Code <= -32000
}

// IsApplicationError returns true if the error is an application error.
// Application errors have positive codes.
func IsApplicationError(err *Error) bool {
	if err == nil {
		return false
	}
	return err.Code > 0
}
