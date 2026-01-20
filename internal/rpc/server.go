package rpc

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
)

// ════════════════════════════════════════════════════════════════
// SERVER STATE
// ════════════════════════════════════════════════════════════════

// ServerState represents the server lifecycle state.
type ServerState int32

const (
	StateUninitialized ServerState = iota
	StateStarting
	StateRunning
	StateShuttingDown
	StateStopped
)

// String returns the state name.
func (s ServerState) String() string {
	switch s {
	case StateUninitialized:
		return "uninitialized"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateShuttingDown:
		return "shutting_down"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// ════════════════════════════════════════════════════════════════
// SERVER
// ════════════════════════════════════════════════════════════════

// Server is a JSON-RPC 2.0 server.
type Server struct {
	codec    *Codec
	registry *Registry

	state  atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Pending responses for async handling
	pendingMu sync.RWMutex
	pending   map[int64]chan *Response

	// Error handler for internal errors
	onError func(error)

	// Hooks
	onStart    func()
	onShutdown func()
}

// ServerOption configures the server.
type ServerOption func(*Server)

// WithErrorHandler sets the error handler for internal errors.
func WithErrorHandler(fn func(error)) ServerOption {
	return func(s *Server) {
		s.onError = fn
	}
}

// WithOnStart sets a callback invoked when the server starts.
func WithOnStart(fn func()) ServerOption {
	return func(s *Server) {
		s.onStart = fn
	}
}

// WithOnShutdown sets a callback invoked when the server shuts down.
func WithOnShutdown(fn func()) ServerOption {
	return func(s *Server) {
		s.onShutdown = fn
	}
}

// WithMiddleware adds middleware to the server's registry.
func WithMiddleware(mw ...Middleware) ServerOption {
	return func(s *Server) {
		s.registry.Use(mw...)
	}
}

// NewServer creates a new RPC server.
func NewServer(transport Transport, opts ...ServerOption) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		codec:    NewCodec(transport),
		registry: NewRegistry(),
		ctx:      ctx,
		cancel:   cancel,
		pending:  make(map[int64]chan *Response),
		onError:  func(error) {}, // Default no-op
	}

	s.state.Store(int32(StateUninitialized))

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Registry returns the server's handler registry.
func (s *Server) Registry() *Registry {
	return s.registry
}

// State returns the current server state.
func (s *Server) State() ServerState {
	return ServerState(s.state.Load())
}

// ════════════════════════════════════════════════════════════════
// HANDLER REGISTRATION
// ════════════════════════════════════════════════════════════════

// Handle registers a handler for a method.
func (s *Server) Handle(method string, handler Handler) {
	s.registry.Register(method, handler)
}

// HandleFunc registers a function as a handler for a method.
func (s *Server) HandleFunc(method string, fn HandlerFunc) {
	s.registry.RegisterFunc(method, fn)
}

// HandleNotification registers a notification handler.
func (s *Server) HandleNotification(method string, handler NotificationHandler) {
	s.registry.RegisterNotification(method, handler)
}

// HandleNotificationFunc registers a function as a notification handler.
func (s *Server) HandleNotificationFunc(method string, fn NotificationHandlerFunc) {
	s.registry.RegisterNotificationFunc(method, fn)
}

// ════════════════════════════════════════════════════════════════
// SERVER LIFECYCLE
// ════════════════════════════════════════════════════════════════

// Start starts the server and begins processing messages.
// This method blocks until the server is stopped.
func (s *Server) Start() error {
	if !s.state.CompareAndSwap(int32(StateUninitialized), int32(StateStarting)) {
		return NewError(ErrCodeInternal, "server already started")
	}

	if s.onStart != nil {
		s.onStart()
	}

	s.state.Store(int32(StateRunning))

	return s.messageLoop()
}

// StartAsync starts the server in a goroutine.
// Returns immediately. Use Wait() to block until stopped.
func (s *Server) StartAsync() error {
	if !s.state.CompareAndSwap(int32(StateUninitialized), int32(StateStarting)) {
		return NewError(ErrCodeInternal, "server already started")
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		if s.onStart != nil {
			s.onStart()
		}

		s.state.Store(int32(StateRunning))

		if err := s.messageLoop(); err != nil && err != io.EOF {
			s.onError(err)
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if !s.state.CompareAndSwap(int32(StateRunning), int32(StateShuttingDown)) {
		// Already shutting down or stopped
		return nil
	}

	if s.onShutdown != nil {
		s.onShutdown()
	}

	// Cancel the server context
	s.cancel()

	// Wait for message loop to finish or context to expire
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.state.Store(int32(StateStopped))
		return nil
	case <-ctx.Done():
		s.state.Store(int32(StateStopped))
		return ctx.Err()
	}
}

// Stop immediately stops the server.
func (s *Server) Stop() {
	s.state.Store(int32(StateShuttingDown))
	s.cancel()
	s.codec.Close()
	s.state.Store(int32(StateStopped))
}

// Wait blocks until the server stops.
func (s *Server) Wait() {
	s.wg.Wait()
}

// ════════════════════════════════════════════════════════════════
// MESSAGE LOOP
// ════════════════════════════════════════════════════════════════

func (s *Server) messageLoop() error {
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
		}

		msg, msgType, err := s.codec.ReadAny()
		if err != nil {
			if err == io.EOF {
				return err
			}
			// Parse error - send error response
			s.sendParseError(err)
			continue
		}

		switch msgType {
		case MessageRequest:
			s.handleRequest(msg.(*Request))

		case MessageNotification:
			s.handleNotification(msg.(*Notification))

		case MessageResponse:
			s.handleResponse(msg.(*Response))

		default:
			s.sendInvalidRequest("unknown message type")
		}
	}
}

// ════════════════════════════════════════════════════════════════
// REQUEST HANDLING
// ════════════════════════════════════════════════════════════════

func (s *Server) handleRequest(req *Request) {
	// Check if shutting down
	if ServerState(s.state.Load()) == StateShuttingDown {
		s.sendError(req.ID, ShuttingDownError())
		return
	}

	// Create request context
	ctx := s.ctx
	ctx = WithRequestID(ctx, req.ID)
	ctx = WithMethod(ctx, req.Method)

	// Dispatch to handler
	result, rpcErr := s.registry.Dispatch(ctx, req)

	// Send response
	if rpcErr != nil {
		s.sendError(req.ID, rpcErr)
	} else {
		s.sendResult(req.ID, result)
	}
}

func (s *Server) handleNotification(notif *Notification) {
	// Create context
	ctx := s.ctx
	ctx = WithMethod(ctx, notif.Method)

	// Dispatch (fire and forget)
	s.registry.DispatchNotification(ctx, notif)
}

func (s *Server) handleResponse(resp *Response) {
	// Look up pending request
	s.pendingMu.RLock()
	ch, ok := s.pending[resp.ID]
	s.pendingMu.RUnlock()

	if ok {
		// Send response to waiting goroutine
		select {
		case ch <- resp:
		default:
			// Channel full or closed, ignore
		}
	}
}

// ════════════════════════════════════════════════════════════════
// SENDING MESSAGES
// ════════════════════════════════════════════════════════════════

func (s *Server) sendResult(id int64, result any) {
	resp, err := NewResponse(id, result)
	if err != nil {
		s.sendError(id, InternalError(err.Error()))
		return
	}

	if err := s.codec.WriteResponse(resp); err != nil {
		s.onError(err)
	}
}

func (s *Server) sendError(id int64, rpcErr *Error) {
	resp := NewErrorResponse(id, rpcErr)
	if err := s.codec.WriteResponse(resp); err != nil {
		s.onError(err)
	}
}

func (s *Server) sendParseError(err error) {
	resp := NewErrorResponse(0, ParseError(err.Error()))
	if err := s.codec.WriteResponse(resp); err != nil {
		s.onError(err)
	}
}

func (s *Server) sendInvalidRequest(detail string) {
	resp := NewErrorResponse(0, InvalidRequestError(detail))
	if err := s.codec.WriteResponse(resp); err != nil {
		s.onError(err)
	}
}

// ════════════════════════════════════════════════════════════════
// NOTIFICATIONS (SERVER → CLIENT)
// ════════════════════════════════════════════════════════════════

// Notify sends a notification to the client.
func (s *Server) Notify(method string, params any) error {
	if ServerState(s.state.Load()) != StateRunning {
		return NewError(ErrCodeNotReady, "server not running")
	}

	notif, err := NewNotification(method, params)
	if err != nil {
		return err
	}

	return s.codec.WriteNotification(notif)
}

// NotifyRender sends a render state update to the client.
func (s *Server) NotifyRender(state *RenderState) error {
	return s.Notify(MethodRender, state)
}

// NotifyRateStatus sends a rate status update to the client.
func (s *Server) NotifyRateStatus(status *RateStatus) error {
	return s.Notify(MethodRateStatus, RateStatusParams{Status: *status})
}

// NotifyLog sends a log message to the client.
func (s *Server) NotifyLog(level, message string) error {
	return s.Notify(MethodLog, map[string]string{
		"level":   level,
		"message": message,
	})
}

// ════════════════════════════════════════════════════════════════
// REQUESTS (SERVER → CLIENT)
// ════════════════════════════════════════════════════════════════

// Request sends a request to the client and waits for a response.
// This is used for reverse requests (server asking client).
func (s *Server) Request(ctx context.Context, method string, params any) (*Response, error) {
	if ServerState(s.state.Load()) != StateRunning {
		return nil, NewError(ErrCodeNotReady, "server not running")
	}

	req, err := NewRequest(method, params)
	if err != nil {
		return nil, err
	}

	// Create response channel
	ch := make(chan *Response, 1)

	s.pendingMu.Lock()
	s.pending[req.ID] = ch
	s.pendingMu.Unlock()

	defer func() {
		s.pendingMu.Lock()
		delete(s.pending, req.ID)
		s.pendingMu.Unlock()
		close(ch)
	}()

	// Send request
	if err := s.codec.WriteRequest(req); err != nil {
		return nil, err
	}

	// Wait for response
	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	}
}

// ════════════════════════════════════════════════════════════════
// UTILITIES
// ════════════════════════════════════════════════════════════════

// IsRunning returns true if the server is running.
func (s *Server) IsRunning() bool {
	return ServerState(s.state.Load()) == StateRunning
}

// IsShuttingDown returns true if the server is shutting down.
func (s *Server) IsShuttingDown() bool {
	return ServerState(s.state.Load()) == StateShuttingDown
}

// IsStopped returns true if the server is stopped.
func (s *Server) IsStopped() bool {
	state := ServerState(s.state.Load())
	return state == StateStopped || state == StateUninitialized
}
