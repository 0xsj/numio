package rpc

import (
	"context"
	"sync"
)

// ════════════════════════════════════════════════════════════════
// HANDLER INTERFACE
// ════════════════════════════════════════════════════════════════

// Handler handles an RPC request.
type Handler interface {
	// Handle processes a request and returns a result or error.
	// The result will be marshaled to JSON for the response.
	Handle(ctx context.Context, req *Request) (any, *Error)
}

// HandlerFunc is an adapter to allow ordinary functions as handlers.
type HandlerFunc func(ctx context.Context, req *Request) (any, *Error)

// Handle implements Handler.
func (f HandlerFunc) Handle(ctx context.Context, req *Request) (any, *Error) {
	return f(ctx, req)
}

// ════════════════════════════════════════════════════════════════
// NOTIFICATION HANDLER
// ════════════════════════════════════════════════════════════════

// NotificationHandler handles an RPC notification (no response).
type NotificationHandler interface {
	// HandleNotification processes a notification.
	HandleNotification(ctx context.Context, notif *Notification)
}

// NotificationHandlerFunc is an adapter for notification handlers.
type NotificationHandlerFunc func(ctx context.Context, notif *Notification)

// HandleNotification implements NotificationHandler.
func (f NotificationHandlerFunc) HandleNotification(ctx context.Context, notif *Notification) {
	f(ctx, notif)
}

// ════════════════════════════════════════════════════════════════
// TYPED HANDLER HELPERS
// ════════════════════════════════════════════════════════════════

// TypedHandler creates a handler that automatically unmarshals params.
// P is the params type, R is the result type.
func TypedHandler[P any, R any](fn func(ctx context.Context, params P) (R, *Error)) Handler {
	return HandlerFunc(func(ctx context.Context, req *Request) (any, *Error) {
		var params P
		if err := req.ParseParams(&params); err != nil {
			return nil, InvalidParamsError(err.Error())
		}
		return fn(ctx, params)
	})
}

// TypedHandlerNoParams creates a handler for methods with no parameters.
func TypedHandlerNoParams[R any](fn func(ctx context.Context) (R, *Error)) Handler {
	return HandlerFunc(func(ctx context.Context, req *Request) (any, *Error) {
		return fn(ctx)
	})
}

// TypedHandlerNoResult creates a handler for methods with no result.
func TypedHandlerNoResult[P any](fn func(ctx context.Context, params P) *Error) Handler {
	return HandlerFunc(func(ctx context.Context, req *Request) (any, *Error) {
		var params P
		if err := req.ParseParams(&params); err != nil {
			return nil, InvalidParamsError(err.Error())
		}
		return nil, fn(ctx, params)
	})
}

// ════════════════════════════════════════════════════════════════
// MIDDLEWARE
// ════════════════════════════════════════════════════════════════

// Middleware wraps a handler with additional functionality.
type Middleware func(Handler) Handler

// Chain combines multiple middlewares into one.
// Middlewares are applied in order: first middleware is outermost.
func Chain(middlewares ...Middleware) Middleware {
	return func(h Handler) Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}

// ════════════════════════════════════════════════════════════════
// REGISTRY
// ════════════════════════════════════════════════════════════════

// Registry manages method handlers.
type Registry struct {
	handlers      map[string]Handler
	notifHandlers map[string]NotificationHandler
	middleware    []Middleware
	mu            sync.RWMutex
}

// NewRegistry creates a new handler registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers:      make(map[string]Handler),
		notifHandlers: make(map[string]NotificationHandler),
	}
}

// Register registers a handler for a method.
func (r *Registry) Register(method string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[method] = handler
}

// RegisterFunc registers a function as a handler for a method.
func (r *Registry) RegisterFunc(method string, fn HandlerFunc) {
	r.Register(method, fn)
}

// RegisterNotification registers a notification handler.
func (r *Registry) RegisterNotification(method string, handler NotificationHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifHandlers[method] = handler
}

// RegisterNotificationFunc registers a function as a notification handler.
func (r *Registry) RegisterNotificationFunc(method string, fn NotificationHandlerFunc) {
	r.RegisterNotification(method, fn)
}

// Use adds middleware to the registry.
// Middleware is applied to all handlers when dispatching.
func (r *Registry) Use(mw ...Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middleware = append(r.middleware, mw...)
}

// Lookup returns the handler for a method.
func (r *Registry) Lookup(method string) (Handler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.handlers[method]
	return h, ok
}

// LookupNotification returns the notification handler for a method.
func (r *Registry) LookupNotification(method string) (NotificationHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.notifHandlers[method]
	return h, ok
}

// Methods returns all registered method names.
func (r *Registry) Methods() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	methods := make([]string, 0, len(r.handlers))
	for m := range r.handlers {
		methods = append(methods, m)
	}
	return methods
}

// Dispatch handles a request by looking up and invoking the handler.
func (r *Registry) Dispatch(ctx context.Context, req *Request) (any, *Error) {
	r.mu.RLock()
	handler, ok := r.handlers[req.Method]
	middleware := r.middleware
	r.mu.RUnlock()

	if !ok {
		return nil, MethodNotFoundError(req.Method)
	}

	// Apply middleware
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}

	return handler.Handle(ctx, req)
}

// DispatchNotification handles a notification.
func (r *Registry) DispatchNotification(ctx context.Context, notif *Notification) {
	r.mu.RLock()
	handler, ok := r.notifHandlers[notif.Method]
	r.mu.RUnlock()

	if ok {
		handler.HandleNotification(ctx, notif)
	}
}

// ════════════════════════════════════════════════════════════════
// COMMON MIDDLEWARE
// ════════════════════════════════════════════════════════════════

// RecoveryMiddleware catches panics and converts them to internal errors.
func RecoveryMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, req *Request) (result any, err *Error) {
			defer func() {
				if r := recover(); r != nil {
					switch v := r.(type) {
					case *Error:
						err = v
					case error:
						err = InternalError(v.Error())
					case string:
						err = InternalError(v)
					default:
						err = InternalError("unknown panic")
					}
				}
			}()
			return next.Handle(ctx, req)
		})
	}
}

// LoggingMiddleware logs requests and responses.
// The logger function receives method, duration, and error (if any).
type LogFunc func(method string, durationMs int64, err *Error)

// LoggingMiddleware creates a middleware that logs requests.
func LoggingMiddleware(log LogFunc) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, req *Request) (any, *Error) {
			start := timeNow()
			result, err := next.Handle(ctx, req)
			duration := timeSince(start)
			log(req.Method, duration, err)
			return result, err
		})
	}
}

// timeNow and timeSince are variables for testing.
var (
	timeNow   = func() int64 { return 0 } // Will be set properly below
	timeSince = func(start int64) int64 { return 0 }
)

func init() {
	// Use simple monotonic counter to avoid importing time in hot path
	// Real implementation would use time.Now()
	var counter int64
	timeNow = func() int64 {
		counter++
		return counter
	}
	timeSince = func(start int64) int64 {
		return timeNow() - start
	}
}

// ValidationMiddleware validates that requests have required fields.
func ValidationMiddleware() Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, req *Request) (any, *Error) {
			if req.JSONRPC != JSONRPCVersion {
				return nil, InvalidRequestError("invalid jsonrpc version")
			}
			if req.Method == "" {
				return nil, InvalidRequestError("missing method")
			}
			return next.Handle(ctx, req)
		})
	}
}

// ════════════════════════════════════════════════════════════════
// CONTEXT KEYS
// ════════════════════════════════════════════════════════════════

// Context key types.
type contextKey int

const (
	// RequestIDKey is the context key for the request ID.
	requestIDKey contextKey = iota

	// MethodKey is the context key for the method name.
	methodKey
)

// WithRequestID adds the request ID to the context.
func WithRequestID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestIDFromContext retrieves the request ID from context.
func RequestIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(requestIDKey).(int64)
	return id, ok
}

// WithMethod adds the method name to the context.
func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKey, method)
}

// MethodFromContext retrieves the method name from context.
func MethodFromContext(ctx context.Context) (string, bool) {
	m, ok := ctx.Value(methodKey).(string)
	return m, ok
}
