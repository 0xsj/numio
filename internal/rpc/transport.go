package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"
)

// ════════════════════════════════════════════════════════════════
// TRANSPORT INTERFACE
// ════════════════════════════════════════════════════════════════

// Transport defines the interface for RPC message transport.
type Transport interface {
	// ReadMessage reads the next message from the transport.
	// Returns io.EOF when the transport is closed.
	ReadMessage() ([]byte, error)

	// WriteMessage writes a message to the transport.
	WriteMessage(data []byte) error

	// Close closes the transport.
	Close() error
}

// ════════════════════════════════════════════════════════════════
// STDIO TRANSPORT
// ════════════════════════════════════════════════════════════════

// StdioTransport implements Transport over stdin/stdout.
// Uses newline-delimited JSON (NDJSON) framing.
type StdioTransport struct {
	reader *bufio.Reader
	writer io.Writer

	writeMu sync.Mutex // Protects concurrent writes
	closed  bool
	closeMu sync.RWMutex
}

// NewStdioTransport creates a new stdio transport.
func NewStdioTransport(stdin io.Reader, stdout io.Writer) *StdioTransport {
	return &StdioTransport{
		reader: bufio.NewReader(stdin),
		writer: stdout,
	}
}

// ReadMessage reads a newline-delimited JSON message.
func (t *StdioTransport) ReadMessage() ([]byte, error) {
	t.closeMu.RLock()
	if t.closed {
		t.closeMu.RUnlock()
		return nil, io.EOF
	}
	t.closeMu.RUnlock()

	line, err := t.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Trim the newline
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	// Also trim carriage return for Windows compatibility
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}

	return line, nil
}

// WriteMessage writes a newline-delimited JSON message.
func (t *StdioTransport) WriteMessage(data []byte) error {
	t.closeMu.RLock()
	if t.closed {
		t.closeMu.RUnlock()
		return io.ErrClosedPipe
	}
	t.closeMu.RUnlock()

	t.writeMu.Lock()
	defer t.writeMu.Unlock()

	// Write message followed by newline
	if _, err := t.writer.Write(data); err != nil {
		return err
	}
	if _, err := t.writer.Write([]byte{'\n'}); err != nil {
		return err
	}

	// Flush if writer supports it
	if f, ok := t.writer.(interface{ Flush() error }); ok {
		return f.Flush()
	}

	return nil
}

// Close closes the transport.
func (t *StdioTransport) Close() error {
	t.closeMu.Lock()
	defer t.closeMu.Unlock()

	t.closed = true

	// Close writer if it supports closing
	if c, ok := t.writer.(io.Closer); ok {
		return c.Close()
	}

	return nil
}

// ════════════════════════════════════════════════════════════════
// MESSAGE CODEC
// ════════════════════════════════════════════════════════════════

// Codec handles encoding and decoding of RPC messages.
type Codec struct {
	transport Transport
}

// NewCodec creates a new codec wrapping a transport.
func NewCodec(transport Transport) *Codec {
	return &Codec{transport: transport}
}

// ReadRequest reads and decodes a request from the transport.
func (c *Codec) ReadRequest() (*Request, error) {
	data, err := c.transport.ReadMessage()
	if err != nil {
		return nil, err
	}

	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

// ReadNotification reads and decodes a notification from the transport.
func (c *Codec) ReadNotification() (*Notification, error) {
	data, err := c.transport.ReadMessage()
	if err != nil {
		return nil, err
	}

	var notif Notification
	if err := json.Unmarshal(data, &notif); err != nil {
		return nil, err
	}

	return &notif, nil
}

// ReadResponse reads and decodes a response from the transport.
func (c *Codec) ReadResponse() (*Response, error) {
	data, err := c.transport.ReadMessage()
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ReadAny reads and decodes any message type from the transport.
func (c *Codec) ReadAny() (any, MessageType, error) {
	data, err := c.transport.ReadMessage()
	if err != nil {
		return nil, MessageUnknown, err
	}

	return ParseMessage(data)
}

// WriteRequest encodes and writes a request to the transport.
func (c *Codec) WriteRequest(req *Request) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return c.transport.WriteMessage(data)
}

// WriteResponse encodes and writes a response to the transport.
func (c *Codec) WriteResponse(resp *Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return c.transport.WriteMessage(data)
}

// WriteNotification encodes and writes a notification to the transport.
func (c *Codec) WriteNotification(notif *Notification) error {
	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	return c.transport.WriteMessage(data)
}

// WriteError writes an error response.
func (c *Codec) WriteError(id int64, err *Error) error {
	resp := NewErrorResponse(id, err)
	return c.WriteResponse(resp)
}

// Close closes the underlying transport.
func (c *Codec) Close() error {
	return c.transport.Close()
}

// ════════════════════════════════════════════════════════════════
// BUFFERED TRANSPORT
// ════════════════════════════════════════════════════════════════

// BufferedTransport wraps a transport with buffered writing.
type BufferedTransport struct {
	transport Transport
	writer    *bufio.Writer
	writeMu   sync.Mutex
}

// NewBufferedTransport creates a buffered transport wrapper.
func NewBufferedTransport(t Transport, bufferSize int) *BufferedTransport {
	// We need access to the underlying writer for buffering
	// This only works with StdioTransport for now
	return &BufferedTransport{
		transport: t,
	}
}

// ════════════════════════════════════════════════════════════════
// CHANNEL TRANSPORT (for testing)
// ════════════════════════════════════════════════════════════════

// ChannelTransport implements Transport using channels.
// Useful for testing without actual I/O.
type ChannelTransport struct {
	in     chan []byte
	out    chan []byte
	closed bool
	mu     sync.RWMutex
}

// NewChannelTransport creates a new channel-based transport.
func NewChannelTransport(bufferSize int) *ChannelTransport {
	return &ChannelTransport{
		in:  make(chan []byte, bufferSize),
		out: make(chan []byte, bufferSize),
	}
}

// NewChannelTransportPair creates two connected transports.
// Messages written to one appear as reads on the other.
func NewChannelTransportPair(bufferSize int) (*ChannelTransport, *ChannelTransport) {
	a := &ChannelTransport{
		in:  make(chan []byte, bufferSize),
		out: make(chan []byte, bufferSize),
	}
	b := &ChannelTransport{
		in:  a.out, // b reads from a's output
		out: a.in,  // b writes to a's input
	}
	return a, b
}

// ReadMessage reads from the input channel.
func (t *ChannelTransport) ReadMessage() ([]byte, error) {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return nil, io.EOF
	}
	t.mu.RUnlock()

	data, ok := <-t.in
	if !ok {
		return nil, io.EOF
	}

	return data, nil
}

// WriteMessage writes to the output channel.
func (t *ChannelTransport) WriteMessage(data []byte) error {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return io.ErrClosedPipe
	}
	t.mu.RUnlock()

	// Make a copy to avoid data races
	cp := make([]byte, len(data))
	copy(cp, data)

	t.out <- cp
	return nil
}

// Close closes both channels.
func (t *ChannelTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.closed {
		t.closed = true
		close(t.out)
	}

	return nil
}

// Inject sends a message directly to the input channel (for testing).
func (t *ChannelTransport) Inject(data []byte) {
	t.in <- data
}

// Receive receives a message from the output channel (for testing).
func (t *ChannelTransport) Receive() ([]byte, bool) {
	data, ok := <-t.out
	return data, ok
}
