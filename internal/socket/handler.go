// internal/socket/handler.go

package socket

import (
	"github.com/0xsj/numio/pkg/engine"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// HANDLER
// ════════════════════════════════════════════════════════════════

// Handler processes incoming requests.
type Handler struct {
	engine *engine.Engine

	// onShutdown is called when a shutdown request is received.
	onShutdown func()
}

// NewHandler creates a new request handler.
func NewHandler(eng *engine.Engine) *Handler {
	return &Handler{
		engine: eng,
	}
}

// SetShutdownCallback sets the function to call on shutdown request.
func (h *Handler) SetShutdownCallback(fn func()) {
	h.onShutdown = fn
}

// ════════════════════════════════════════════════════════════════
// DISPATCH
// ════════════════════════════════════════════════════════════════

// Handle processes a request and returns a response.
func (h *Handler) Handle(req *Request) *Response {
	switch req.Method {
	case MethodEval:
		return h.handleEval(req)
	case MethodEvalLines:
		return h.handleEvalLines(req)
	case MethodGetVars:
		return h.handleGetVars(req)
	case MethodSetVar:
		return h.handleSetVar(req)
	case MethodClear:
		return h.handleClear(req)
	case MethodPing:
		return h.handlePing(req)
	case MethodShutdown:
		return h.handleShutdown(req)
	default:
		return NewMethodNotFoundError(req.ID, string(req.Method))
	}
}

// ════════════════════════════════════════════════════════════════
// METHOD HANDLERS
// ════════════════════════════════════════════════════════════════

func (h *Handler) handleEval(req *Request) *Response {
	params, err := ParseParams[EvalParams](req)
	if err != nil {
		return NewInvalidParamsError(req.ID, err.Error())
	}

	if params.Expression == "" {
		return NewInvalidParamsError(req.ID, "expression is required")
	}

	result := h.engine.Eval(params.Expression)

	return NewResponse(req.ID, valueToEvalResult(result))
}

func (h *Handler) handleEvalLines(req *Request) *Response {
	params, err := ParseParams[EvalLinesParams](req)
	if err != nil {
		return NewInvalidParamsError(req.ID, err.Error())
	}

	if len(params.Lines) == 0 {
		return NewInvalidParamsError(req.ID, "lines is required")
	}

	results := make([]EvalResult, len(params.Lines))
	for i, line := range params.Lines {
		result := h.engine.Eval(line)
		results[i] = valueToEvalResult(result)
	}

	return NewResponse(req.ID, EvalLinesResult{Results: results})
}

func (h *Handler) handleGetVars(req *Request) *Response {
	vars := h.engine.Variables()

	// Convert types.Value to basic types for JSON
	result := make(map[string]any, len(vars))
	for name, val := range vars {
		result[name] = valueToAny(val)
	}

	return NewResponse(req.ID, VarsResult{Vars: result})
}

func (h *Handler) handleSetVar(req *Request) *Response {
	params, err := ParseParams[SetVarParams](req)
	if err != nil {
		return NewInvalidParamsError(req.ID, err.Error())
	}

	if params.Name == "" {
		return NewInvalidParamsError(req.ID, "name is required")
	}

	// Evaluate the value expression
	result := h.engine.Eval(params.Value)
	if result.IsError() {
		return NewEvalError(req.ID, result.String())
	}

	// Set the variable
	h.engine.SetVariable(params.Name, result)

	return NewResponse(req.ID, OkResult{Ok: true})
}

func (h *Handler) handleClear(req *Request) *Response {
	h.engine.Clear()
	return NewResponse(req.ID, OkResult{Ok: true})
}

func (h *Handler) handlePing(req *Request) *Response {
	return NewResponse(req.ID, PingResult{Pong: true})
}

func (h *Handler) handleShutdown(req *Request) *Response {
	if h.onShutdown != nil {
		// Call shutdown in goroutine so we can send response first
		go h.onShutdown()
	}
	return NewResponse(req.ID, OkResult{Ok: true})
}

// ════════════════════════════════════════════════════════════════
// VALUE CONVERSION
// ════════════════════════════════════════════════════════════════

// valueToEvalResult converts a types.Value to an EvalResult.
func valueToEvalResult(val types.Value) EvalResult {
	result := EvalResult{
		Value: val.String(),
		Type:  val.Kind.String(),
	}

	// Add raw numeric value if applicable
	if val.IsNumeric() {
		raw := val.AsFloat()
		result.Raw = &raw
	}

	// Add unit/currency code if applicable
	if val.IsCurrency() && val.Curr != nil {
		result.Unit = val.Curr.Code
	} else if val.IsUnit() && val.Unit != nil {
		result.Unit = val.Unit.Code
	} else if val.IsCrypto() && val.Crypto != nil {
		result.Unit = val.Crypto.Code
	} else if val.IsMetal() && val.Metal != nil {
		result.Unit = val.Metal.Code
	}

	return result
}

// valueToAny converts a types.Value to a JSON-friendly type.
func valueToAny(val types.Value) any {
	switch {
	case val.IsError():
		return map[string]any{
			"error": val.ErrorMessage(),
		}
	case val.IsEmpty():
		return nil
	case val.IsCurrency() && val.Curr != nil:
		return map[string]any{
			"value":    val.AsFloat(),
			"currency": val.Curr.Code,
			"display":  val.String(),
		}
	case val.IsUnit() && val.Unit != nil:
		return map[string]any{
			"value":   val.AsFloat(),
			"unit":    val.Unit.Code,
			"display": val.String(),
		}
	case val.IsCrypto() && val.Crypto != nil:
		return map[string]any{
			"value":   val.AsFloat(),
			"crypto":  val.Crypto.Code,
			"display": val.String(),
		}
	case val.IsMetal() && val.Metal != nil:
		return map[string]any{
			"value":   val.AsFloat(),
			"metal":   val.Metal.Code,
			"display": val.String(),
		}
	case val.IsPercentage():
		return map[string]any{
			"value":   val.AsFloat(),
			"display": val.String(),
		}
	case val.IsNumber():
		return val.AsFloat()
	case val.IsDate():
		return map[string]any{
			"value":   val.String(),
			"display": val.String(),
		}
	case val.IsString():
		return val.AsString()
	default:
		return val.String()
	}
}
