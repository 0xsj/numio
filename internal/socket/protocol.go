// internal/socket/protocol.go

// Package socket provides Unix domain socket IPC for numio.
package socket

import (
	"encoding/json"
	"fmt"
)

// ════════════════════════════════════════════════════════════════
// REQUEST TYPES
// ════════════════════════════════════════════════════════════════

// Request represents an incoming IPC request.
type Request struct {
	// ID is a client-provided identifier for matching responses.
	// Optional but recommended for async clients.
	ID any `json:"id,omitempty"`

	// Method is the operation to perform.
	Method Method `json:"method"`

	// Params contains method-specific parameters.
	// Use json.RawMessage to defer parsing until we know the method.
	Params json.RawMessage `json:"params,omitempty"`
}

// Method represents an API method.
type Method string

const (
	MethodEval      Method = "eval"
	MethodEvalLines Method = "eval_lines"
	MethodGetVars   Method = "get_vars"
	MethodSetVar    Method = "set_var"
	MethodClear     Method = "clear"
	MethodPing      Method = "ping"
	MethodShutdown  Method = "shutdown"
)

// EvalParams contains parameters for the eval method.
type EvalParams struct {
	Expression string `json:"expression"`
}

// EvalLinesParams contains parameters for the eval_lines method.
type EvalLinesParams struct {
	Lines []string `json:"lines"`
}

// SetVarParams contains parameters for the set_var method.
type SetVarParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ════════════════════════════════════════════════════════════════
// RESPONSE TYPES
// ════════════════════════════════════════════════════════════════

// Response represents an outgoing IPC response.
type Response struct {
	// ID echoes the request ID for matching.
	ID any `json:"id,omitempty"`

	// Result contains the successful result (if no error).
	Result any `json:"result,omitempty"`

	// Error contains error information (if failed).
	Error *Error `json:"error,omitempty"`
}

// Error represents an error response.
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// ErrorCode represents error categories.
type ErrorCode int

const (
	ErrParse          ErrorCode = -32700 // Invalid JSON
	ErrInvalidReq     ErrorCode = -32600 // Invalid request structure
	ErrMethodNotFound ErrorCode = -32601 // Unknown method
	ErrInvalidParams  ErrorCode = -32602 // Invalid method parameters
	ErrInternal       ErrorCode = -32603 // Internal server error
	ErrEval           ErrorCode = -32000 // Evaluation error
)

// ════════════════════════════════════════════════════════════════
// RESULT TYPES
// ════════════════════════════════════════════════════════════════

// EvalResult contains the result of an eval operation.
type EvalResult struct {
	// Value is the formatted result string.
	Value string `json:"value"`

	// Type is the result type (number, currency, date, etc.).
	Type string `json:"type"`

	// Raw is the numeric value if applicable.
	Raw *float64 `json:"raw,omitempty"`

	// Unit is the unit/currency code if applicable.
	Unit string `json:"unit,omitempty"`
}

// EvalLinesResult contains results for multiple evaluations.
type EvalLinesResult struct {
	Results []EvalResult `json:"results"`
}

// VarsResult contains the current variables.
type VarsResult struct {
	Vars map[string]any `json:"vars"`
}

// OkResult indicates success for void operations.
type OkResult struct {
	Ok bool `json:"ok"`
}

// PingResult is the response to a ping.
type PingResult struct {
	Pong bool `json:"pong"`
}

// ════════════════════════════════════════════════════════════════
// CONSTRUCTORS
// ════════════════════════════════════════════════════════════════

// NewResponse creates a successful response.
func NewResponse(id any, result any) *Response {
	return &Response{
		ID:     id,
		Result: result,
	}
}

// NewErrorResponse creates an error response.
func NewErrorResponse(id any, code ErrorCode, message string) *Response {
	return &Response{
		ID: id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}

// NewParseError creates a parse error response.
func NewParseError(id any, message string) *Response {
	return NewErrorResponse(id, ErrParse, message)
}

// NewInvalidRequestError creates an invalid request error response.
func NewInvalidRequestError(id any, message string) *Response {
	return NewErrorResponse(id, ErrInvalidReq, message)
}

// NewMethodNotFoundError creates a method not found error response.
func NewMethodNotFoundError(id any, method string) *Response {
	return NewErrorResponse(id, ErrMethodNotFound, fmt.Sprintf("method not found: %s", method))
}

// NewInvalidParamsError creates an invalid params error response.
func NewInvalidParamsError(id any, message string) *Response {
	return NewErrorResponse(id, ErrInvalidParams, message)
}

// NewInternalError creates an internal error response.
func NewInternalError(id any, message string) *Response {
	return NewErrorResponse(id, ErrInternal, message)
}

// NewEvalError creates an evaluation error response.
func NewEvalError(id any, message string) *Response {
	return NewErrorResponse(id, ErrEval, message)
}

// ════════════════════════════════════════════════════════════════
// ENCODING / DECODING
// ════════════════════════════════════════════════════════════════

// DecodeRequest parses a JSON request.
func DecodeRequest(data []byte) (*Request, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return &req, nil
}

// EncodeResponse serializes a response to JSON.
func EncodeResponse(resp *Response) ([]byte, error) {
	return json.Marshal(resp)
}

// ParseParams parses request params into the given type.
func ParseParams[T any](req *Request) (*T, error) {
	if req.Params == nil {
		return nil, fmt.Errorf("missing params")
	}

	var params T
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	return &params, nil
}

// ════════════════════════════════════════════════════════════════
// SOCKET PATHS
// ════════════════════════════════════════════════════════════════

const (
	// DefaultSocketPath is the default Unix socket path.
	DefaultSocketPath = "/tmp/numio.sock"

	// DefaultWindowsPipe is the default Windows named pipe.
	DefaultWindowsPipe = `\\.\pipe\numio`
)

// GetSocketPath returns the appropriate socket path for the OS.
func GetSocketPath() string {
	// Could check runtime.GOOS here for Windows
	return DefaultSocketPath
}
