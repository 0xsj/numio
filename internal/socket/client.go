// internal/socket/client.go

package socket

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ════════════════════════════════════════════════════════════════
// CLIENT
// ════════════════════════════════════════════════════════════════

// Client is a Unix domain socket client for numio IPC.
type Client struct {
	socketPath string
	conn       net.Conn
	reader     *bufio.Reader

	// Configuration
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration

	// Request ID counter
	nextID atomic.Int64

	// Mutex for thread-safe operations
	mu sync.Mutex
}

// Dial connects to a numio server at the default socket path.
func Dial() (*Client, error) {
	return DialPath(DefaultSocketPath)
}

// DialPath connects to a numio server at the specified socket path.
func DialPath(socketPath string) (*Client, error) {
	return DialPathWithTimeout(socketPath, 5*time.Second)
}

// DialPathWithTimeout connects with a custom timeout.
func DialPathWithTimeout(socketPath string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("unix", socketPath, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		socketPath:     socketPath,
		conn:           conn,
		reader:         bufio.NewReader(conn),
		connectTimeout: timeout,
		readTimeout:    30 * time.Second,
		writeTimeout:   10 * time.Second,
	}, nil
}

// ════════════════════════════════════════════════════════════════
// CONFIGURATION
// ════════════════════════════════════════════════════════════════

// SetReadTimeout sets the read timeout.
func (c *Client) SetReadTimeout(d time.Duration) {
	c.readTimeout = d
}

// SetWriteTimeout sets the write timeout.
func (c *Client) SetWriteTimeout(d time.Duration) {
	c.writeTimeout = d
}

// ════════════════════════════════════════════════════════════════
// LIFECYCLE
// ════════════════════════════════════════════════════════════════

// Close closes the connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Reconnect attempts to reconnect to the server.
func (c *Client) Reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close existing connection
	if c.conn != nil {
		c.conn.Close()
	}

	// Reconnect
	conn, err := net.DialTimeout("unix", c.socketPath, c.connectTimeout)
	if err != nil {
		return err
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	return nil
}

// ════════════════════════════════════════════════════════════════
// LOW-LEVEL API
// ════════════════════════════════════════════════════════════════

// Send sends a request and returns the response.
func (c *Client) Send(req *Request) (*Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Assign ID if not set
	if req.ID == nil {
		req.ID = c.nextID.Add(1)
	}

	// Encode request
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')

	// Set write deadline
	if c.writeTimeout > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	// Send request
	if _, err := c.conn.Write(data); err != nil {
		return nil, err
	}

	// Set read deadline
	if c.readTimeout > 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}

	// Read response
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Decode response
	var resp Response
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Call is a convenience method that creates a request and sends it.
func (c *Client) Call(method Method, params any) (*Response, error) {
	var paramsJSON json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		paramsJSON = data
	}

	req := &Request{
		Method: method,
		Params: paramsJSON,
	}

	return c.Send(req)
}

// ════════════════════════════════════════════════════════════════
// HIGH-LEVEL API
// ════════════════════════════════════════════════════════════════

// Eval evaluates an expression and returns the result.
func (c *Client) Eval(expression string) (*EvalResult, error) {
	resp, err := c.Call(MethodEval, EvalParams{Expression: expression})
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	// Parse result
	data, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	var result EvalResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EvalString evaluates an expression and returns the string result.
func (c *Client) EvalString(expression string) (string, error) {
	result, err := c.Eval(expression)
	if err != nil {
		return "", err
	}
	return result.Value, nil
}

// EvalLines evaluates multiple expressions.
func (c *Client) EvalLines(lines []string) ([]EvalResult, error) {
	resp, err := c.Call(MethodEvalLines, EvalLinesParams{Lines: lines})
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	// Parse result
	data, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	var result EvalLinesResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// GetVars returns all variables.
func (c *Client) GetVars() (map[string]any, error) {
	resp, err := c.Call(MethodGetVars, nil)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	// Parse result
	data, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	var result VarsResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Vars, nil
}

// SetVar sets a variable.
func (c *Client) SetVar(name, value string) error {
	resp, err := c.Call(MethodSetVar, SetVarParams{Name: name, Value: value})
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// Clear clears the engine state.
func (c *Client) Clear() error {
	resp, err := c.Call(MethodClear, nil)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// Ping checks if the server is alive.
func (c *Client) Ping() error {
	resp, err := c.Call(MethodPing, nil)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// Shutdown requests the server to shut down.
func (c *Client) Shutdown() error {
	resp, err := c.Call(MethodShutdown, nil)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// QuickEval connects, evaluates, and disconnects in one call.
// Useful for one-off evaluations.
func QuickEval(expression string) (string, error) {
	client, err := Dial()
	if err != nil {
		return "", err
	}
	defer client.Close()

	return client.EvalString(expression)
}

// IsServerRunning checks if a numio server is running.
func IsServerRunning() bool {
	return IsServerRunningAt(DefaultSocketPath)
}

// IsServerRunningAt checks if a numio server is running at the given path.
func IsServerRunningAt(socketPath string) bool {
	client, err := DialPathWithTimeout(socketPath, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer client.Close()

	return client.Ping() == nil
}
