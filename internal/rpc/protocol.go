// Package rpc provides JSON-RPC 2.0 communication for numio desktop.
package rpc

import (
	"encoding/json"
	"sync/atomic"
)

// ════════════════════════════════════════════════════════════════
// CONSTANTS
// ════════════════════════════════════════════════════════════════

// JSONRPCVersion is the JSON-RPC protocol version.
const JSONRPCVersion = "2.0"

// Method names for UI → Core requests.
const (
	MethodInitialize = "initialize" // Start session, get initial state
	MethodShutdown   = "shutdown"   // Graceful shutdown
	MethodKeypress   = "keypress"   // Send key event
	MethodResize     = "resize"     // Viewport resized
	MethodGetState   = "getState"   // Request current render state
	MethodSetOption  = "setOption"  // Set configuration option
)

// Method names for Core → UI notifications.
const (
	MethodRender     = "render"     // Push new render state
	MethodRateStatus = "rateStatus" // Rate fetch status update
	MethodLog        = "log"        // Debug/info logging
)

// ════════════════════════════════════════════════════════════════
// ID GENERATION
// ════════════════════════════════════════════════════════════════

// idCounter is used to generate unique request IDs.
var idCounter atomic.Int64

// nextID generates a unique request ID.
func nextID() int64 {
	return idCounter.Add(1)
}

// ════════════════════════════════════════════════════════════════
// REQUEST
// ════════════════════════════════════════════════════════════════

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// NewRequest creates a new request with an auto-generated ID.
func NewRequest(method string, params any) (*Request, error) {
	req := &Request{
		JSONRPC: JSONRPCVersion,
		ID:      nextID(),
		Method:  method,
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		req.Params = data
	}

	return req, nil
}

// NewRequestWithID creates a new request with a specific ID.
func NewRequestWithID(id int64, method string, params any) (*Request, error) {
	req := &Request{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  method,
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		req.Params = data
	}

	return req, nil
}

// ParseParams unmarshals the request params into the given target.
func (r *Request) ParseParams(target any) error {
	if r.Params == nil {
		return nil
	}
	return json.Unmarshal(r.Params, target)
}

// ════════════════════════════════════════════════════════════════
// RESPONSE
// ════════════════════════════════════════════════════════════════

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// NewResponse creates a successful response.
func NewResponse(id int64, result any) (*Response, error) {
	resp := &Response{
		JSONRPC: JSONRPCVersion,
		ID:      id,
	}

	if result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		resp.Result = data
	}

	return resp, nil
}

// NewErrorResponse creates an error response.
func NewErrorResponse(id int64, err *Error) *Response {
	return &Response{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error:   err,
	}
}

// ParseResult unmarshals the response result into the given target.
func (r *Response) ParseResult(target any) error {
	if r.Result == nil {
		return nil
	}
	return json.Unmarshal(r.Result, target)
}

// IsError returns true if this is an error response.
func (r *Response) IsError() bool {
	return r.Error != nil
}

// ════════════════════════════════════════════════════════════════
// NOTIFICATION
// ════════════════════════════════════════════════════════════════

// Notification represents a JSON-RPC 2.0 notification (no response expected).
type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// NewNotification creates a new notification.
func NewNotification(method string, params any) (*Notification, error) {
	notif := &Notification{
		JSONRPC: JSONRPCVersion,
		Method:  method,
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		notif.Params = data
	}

	return notif, nil
}

// ParseParams unmarshals the notification params into the given target.
func (n *Notification) ParseParams(target any) error {
	if n.Params == nil {
		return nil
	}
	return json.Unmarshal(n.Params, target)
}

// ════════════════════════════════════════════════════════════════
// ERROR
// ════════════════════════════════════════════════════════════════

// Standard JSON-RPC 2.0 error codes.
const (
	ErrCodeParse          = -32700 // Invalid JSON
	ErrCodeInvalidRequest = -32600 // Not a valid request object
	ErrCodeMethodNotFound = -32601 // Method does not exist
	ErrCodeInvalidParams  = -32602 // Invalid method parameters
	ErrCodeInternal       = -32603 // Internal JSON-RPC error
)

// Application-specific error codes (must be outside -32000 to -32099).
const (
	ErrCodeEvaluation   = 1 // Evaluation error
	ErrCodeNotReady     = 2 // Core not initialized
	ErrCodeShuttingDown = 3 // Core is shutting down
)

// Error represents a JSON-RPC 2.0 error object.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new RPC error.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithData creates a new RPC error with additional data.
func NewErrorWithData(code int, message string, data any) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Standard error constructors.

// ErrParse creates a parse error.
func ErrParse(detail string) *Error {
	return NewErrorWithData(ErrCodeParse, "Parse error", detail)
}

// ErrInvalidRequest creates an invalid request error.
func ErrInvalidRequest(detail string) *Error {
	return NewErrorWithData(ErrCodeInvalidRequest, "Invalid request", detail)
}

// ErrMethodNotFound creates a method not found error.
func ErrMethodNotFound(method string) *Error {
	return NewErrorWithData(ErrCodeMethodNotFound, "Method not found", method)
}

// ErrInvalidParams creates an invalid params error.
func ErrInvalidParams(detail string) *Error {
	return NewErrorWithData(ErrCodeInvalidParams, "Invalid params", detail)
}

// ErrInternal creates an internal error.
func ErrInternal(detail string) *Error {
	return NewErrorWithData(ErrCodeInternal, "Internal error", detail)
}

// ════════════════════════════════════════════════════════════════
// MESSAGE DETECTION
// ════════════════════════════════════════════════════════════════

// MessageType indicates the type of JSON-RPC message.
type MessageType int

const (
	MessageUnknown MessageType = iota
	MessageRequest
	MessageResponse
	MessageNotification
)

// messageProbe is used to determine message type without full parsing.
type messageProbe struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *int64          `json:"id"`
	Method  *string         `json:"method"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error"`
}

// DetectMessageType determines the type of a JSON-RPC message.
func DetectMessageType(data []byte) (MessageType, error) {
	var probe messageProbe
	if err := json.Unmarshal(data, &probe); err != nil {
		return MessageUnknown, err
	}

	// Response: has id and (result or error), no method
	if probe.ID != nil && probe.Method == nil {
		return MessageResponse, nil
	}

	// Request: has id and method
	if probe.ID != nil && probe.Method != nil {
		return MessageRequest, nil
	}

	// Notification: has method but no id
	if probe.ID == nil && probe.Method != nil {
		return MessageNotification, nil
	}

	return MessageUnknown, nil
}

// ParseMessage parses a JSON-RPC message into the appropriate type.
func ParseMessage(data []byte) (any, MessageType, error) {
	msgType, err := DetectMessageType(data)
	if err != nil {
		return nil, MessageUnknown, err
	}

	switch msgType {
	case MessageRequest:
		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, MessageRequest, err
		}
		return &req, MessageRequest, nil

	case MessageResponse:
		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, MessageResponse, err
		}
		return &resp, MessageResponse, nil

	case MessageNotification:
		var notif Notification
		if err := json.Unmarshal(data, &notif); err != nil {
			return nil, MessageNotification, err
		}
		return &notif, MessageNotification, nil

	default:
		return nil, MessageUnknown, nil
	}
}
