// internal/socket/socket.go

package socket

import (
	"bufio"
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/0xsj/numio/pkg/engine"
)

// ════════════════════════════════════════════════════════════════
// SERVER
// ════════════════════════════════════════════════════════════════

// Server is a Unix domain socket server for numio IPC.
type Server struct {
	socketPath string
	handler    *Handler
	listener   net.Listener

	// Configuration
	readTimeout  time.Duration
	writeTimeout time.Duration

	// State
	mu       sync.Mutex
	running  bool
	shutdown chan struct{}
	wg       sync.WaitGroup
}

// NewServer creates a new Unix socket server.
func NewServer(eng *engine.Engine) *Server {
	return &Server{
		socketPath:   DefaultSocketPath,
		handler:      NewHandler(eng),
		readTimeout:  30 * time.Second,
		writeTimeout: 10 * time.Second,
		shutdown:     make(chan struct{}),
	}
}

// NewServerWithPath creates a server with a custom socket path.
func NewServerWithPath(eng *engine.Engine, socketPath string) *Server {
	s := NewServer(eng)
	s.socketPath = socketPath
	return s
}

// ════════════════════════════════════════════════════════════════
// CONFIGURATION
// ════════════════════════════════════════════════════════════════

// SetReadTimeout sets the read timeout for client connections.
func (s *Server) SetReadTimeout(d time.Duration) {
	s.readTimeout = d
}

// SetWriteTimeout sets the write timeout for client connections.
func (s *Server) SetWriteTimeout(d time.Duration) {
	s.writeTimeout = d
}

// SocketPath returns the socket path.
func (s *Server) SocketPath() string {
	return s.socketPath
}

// ════════════════════════════════════════════════════════════════
// LIFECYCLE
// ════════════════════════════════════════════════════════════════

// Start begins listening for connections.
// This method blocks until the server is stopped.
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return ErrAlreadyRunning
	}

	// Clean up stale socket file
	if err := s.cleanup(); err != nil {
		s.mu.Unlock()
		return err
	}

	// Create listener
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		s.mu.Unlock()
		return err
	}

	// Set restrictive permissions (owner only)
	if err := os.Chmod(s.socketPath, 0600); err != nil {
		listener.Close()
		s.mu.Unlock()
		return err
	}

	s.listener = listener
	s.running = true
	s.shutdown = make(chan struct{})

	// Set shutdown callback
	s.handler.SetShutdownCallback(func() {
		s.Stop()
	})

	s.mu.Unlock()

	// Accept loop
	return s.acceptLoop()
}

// StartAsync starts the server in a background goroutine.
// Returns a channel that receives the error when the server stops.
func (s *Server) StartAsync() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start()
	}()
	return errCh
}

// Stop gracefully stops the server.
func (s *Server) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}

	s.running = false
	close(s.shutdown)

	if s.listener != nil {
		s.listener.Close()
	}
	s.mu.Unlock()

	// Wait for active connections to finish
	s.wg.Wait()

	// Clean up socket file
	return s.cleanup()
}

// IsRunning returns true if the server is running.
func (s *Server) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// cleanup removes stale socket file.
func (s *Server) cleanup() error {
	// Check if socket file exists
	if _, err := os.Stat(s.socketPath); os.IsNotExist(err) {
		return nil
	}

	// Try to connect to see if another server is running
	conn, err := net.DialTimeout("unix", s.socketPath, 100*time.Millisecond)
	if err == nil {
		// Another server is running
		conn.Close()
		return ErrSocketInUse
	}

	// Stale socket file, remove it
	return os.Remove(s.socketPath)
}

// ════════════════════════════════════════════════════════════════
// CONNECTION HANDLING
// ════════════════════════════════════════════════════════════════

// acceptLoop accepts incoming connections.
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				// Expected error from shutdown
				return nil
			default:
				// Unexpected error
				return err
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection handles a single client connection.
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		// Check for shutdown
		select {
		case <-s.shutdown:
			return
		default:
		}

		// Set read deadline
		if s.readTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		}

		// Read line (JSON request)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			// Connection closed or timeout
			return
		}

		// Process request
		response := s.processRequest(line)

		// Set write deadline
		if s.writeTimeout > 0 {
			conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
		}

		// Send response
		responseBytes, err := EncodeResponse(response)
		if err != nil {
			// Encoding error - send internal error
			errorResp := NewInternalError(nil, "failed to encode response")
			responseBytes, _ = EncodeResponse(errorResp)
		}

		// Append newline for JSON Lines protocol
		responseBytes = append(responseBytes, '\n')

		if _, err := conn.Write(responseBytes); err != nil {
			return
		}
	}
}

// processRequest parses and handles a request.
func (s *Server) processRequest(data []byte) *Response {
	req, err := DecodeRequest(data)
	if err != nil {
		return NewParseError(nil, err.Error())
	}

	if req.Method == "" {
		return NewInvalidRequestError(req.ID, "method is required")
	}

	return s.handler.Handle(req)
}

// ════════════════════════════════════════════════════════════════
// SIGNAL HANDLING
// ════════════════════════════════════════════════════════════════

// ListenForSignals starts the server and handles OS signals for graceful shutdown.
// This is a convenience method for typical server usage.
func (s *Server) ListenForSignals(ctx context.Context) error {
	// Start server in background
	errCh := s.StartAsync()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		// Server stopped with error
		return err
	case <-sigCh:
		// Received shutdown signal
		return s.Stop()
	case <-ctx.Done():
		// Context cancelled
		return s.Stop()
	}
}

// ════════════════════════════════════════════════════════════════
// ERRORS
// ════════════════════════════════════════════════════════════════

// Server errors.
var (
	ErrAlreadyRunning = &serverError{"server is already running"}
	ErrSocketInUse    = &serverError{"socket is in use by another process"}
	ErrNotRunning     = &serverError{"server is not running"}
)

type serverError struct {
	msg string
}

func (e *serverError) Error() string {
	return e.msg
}
