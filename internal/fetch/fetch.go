// internal/fetch/fetch.go

// Package fetch provides HTTP client utilities for fetching exchange rates
// and asset prices from external APIs.
package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

// Default client settings
const (
	DefaultTimeout     = 10 * time.Second
	DefaultMaxRetries  = 3
	DefaultRateLimit   = 10 // requests per second
	DefaultUserAgent   = "numio/1.0"
	DefaultBackoffBase = 500 * time.Millisecond
	DefaultBackoffMax  = 5 * time.Second
)

// Common errors
var (
	ErrRequestFailed   = errors.New("request failed")
	ErrRateLimited     = errors.New("rate limited")
	ErrTimeout         = errors.New("request timeout")
	ErrInvalidResponse = errors.New("invalid response")
	ErrNotFound        = errors.New("resource not found")
	ErrUnauthorized    = errors.New("unauthorized (check API key)")
)

// ════════════════════════════════════════════════════════════════
// CLIENT
// ════════════════════════════════════════════════════════════════

// Client is an HTTP client wrapper with retry, rate limiting, and timeout.
type Client struct {
	http        *http.Client
	userAgent   string
	maxRetries  int
	rateLimiter *rateLimiter
	backoffBase time.Duration
	backoffMax  time.Duration
}

// NewClient creates a new Client with default settings.
func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: DefaultTimeout,
		},
		userAgent:   DefaultUserAgent,
		maxRetries:  DefaultMaxRetries,
		rateLimiter: newRateLimiter(DefaultRateLimit),
		backoffBase: DefaultBackoffBase,
		backoffMax:  DefaultBackoffMax,
	}
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// NewClientWithOptions creates a Client with custom options.
func NewClientWithOptions(opts ...ClientOption) *Client {
	c := NewClient()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.http.Timeout = d
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(n int) ClientOption {
	return func(c *Client) {
		if n >= 0 {
			c.maxRetries = n
		}
	}
}

// WithRateLimit sets the rate limit (requests per second).
func WithRateLimit(rps int) ClientOption {
	return func(c *Client) {
		if rps > 0 {
			c.rateLimiter = newRateLimiter(rps)
		}
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithBackoff sets the backoff parameters for retries.
func WithBackoff(base, max time.Duration) ClientOption {
	return func(c *Client) {
		if base > 0 {
			c.backoffBase = base
		}
		if max > 0 {
			c.backoffMax = max
		}
	}
}

// ════════════════════════════════════════════════════════════════
// HTTP METHODS
// ════════════════════════════════════════════════════════════════

// Get performs an HTTP GET request with retries and rate limiting.
func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	return c.Do(ctx, http.MethodGet, url, nil)
}

// GetJSON performs a GET request and decodes the JSON response.
func (c *Client) GetJSON(ctx context.Context, url string, v any) error {
	resp, err := c.Get(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Close()

	return resp.JSON(v)
}

// Do performs an HTTP request with retries and rate limiting.
func (c *Client) Do(ctx context.Context, method, url string, body io.Reader) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Wait for rate limiter
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		// Apply backoff on retry
		if attempt > 0 {
			backoff := c.calculateBackoff(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Accept", "application/json")

		// Execute request
		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = c.wrapError(err)
			// Retry on timeout or temporary errors
			if isRetryable(err) {
				continue
			}
			return nil, lastErr
		}

		// Check status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return newResponse(resp), nil
		}

		// Handle error responses
		resp.Body.Close()
		lastErr = c.statusError(resp.StatusCode)

		// Retry on server errors (5xx) or rate limiting (429)
		if isRetryableStatus(resp.StatusCode) {
			continue
		}

		return nil, lastErr
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ErrRequestFailed
}

// calculateBackoff returns the backoff duration for a retry attempt.
// Uses exponential backoff with jitter.
func (c *Client) calculateBackoff(attempt int) time.Duration {
	backoff := c.backoffBase
	for i := 1; i < attempt; i++ {
		backoff *= 2
		if backoff > c.backoffMax {
			backoff = c.backoffMax
			break
		}
	}
	return backoff
}

// wrapError wraps an HTTP error with a more descriptive error.
func (c *Client) wrapError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	if errors.Is(err, context.Canceled) {
		return err
	}
	return err
}

// statusError returns an appropriate error for an HTTP status code.
func (c *Client) statusError(status int) error {
	switch status {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrRateLimited
	default:
		return ErrRequestFailed
	}
}

// isRetryable returns true if the error is retryable.
func isRetryable(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// Could add more checks for network errors
	return false
}

// isRetryableStatus returns true if the HTTP status code is retryable.
func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests ||
		status == http.StatusServiceUnavailable ||
		status == http.StatusGatewayTimeout ||
		status == http.StatusBadGateway
}

// ════════════════════════════════════════════════════════════════
// RESPONSE
// ════════════════════════════════════════════════════════════════

// Response wraps an HTTP response with helper methods.
type Response struct {
	StatusCode int
	Header     http.Header
	body       io.ReadCloser
}

// newResponse creates a Response from an http.Response.
func newResponse(r *http.Response) *Response {
	return &Response{
		StatusCode: r.StatusCode,
		Header:     r.Header,
		body:       r.Body,
	}
}

// Close closes the response body.
func (r *Response) Close() error {
	if r.body != nil {
		return r.body.Close()
	}
	return nil
}

// Body returns the response body reader.
func (r *Response) Body() io.ReadCloser {
	return r.body
}

// Bytes reads and returns the response body as bytes.
func (r *Response) Bytes() ([]byte, error) {
	if r.body == nil {
		return nil, nil
	}
	return io.ReadAll(r.body)
}

// JSON decodes the response body as JSON into v.
func (r *Response) JSON(v any) error {
	if r.body == nil {
		return ErrInvalidResponse
	}

	decoder := json.NewDecoder(r.body)
	if err := decoder.Decode(v); err != nil {
		return ErrInvalidResponse
	}
	return nil
}

// ════════════════════════════════════════════════════════════════
// RATE LIMITER (Token Bucket)
// ════════════════════════════════════════════════════════════════

// rateLimiter implements a simple token bucket rate limiter.
type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// newRateLimiter creates a rate limiter with the given requests per second.
func newRateLimiter(rps int) *rateLimiter {
	if rps <= 0 {
		rps = DefaultRateLimit
	}
	return &rateLimiter{
		tokens:     float64(rps),
		maxTokens:  float64(rps),
		refillRate: float64(rps),
		lastRefill: time.Now(),
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (r *rateLimiter) Wait(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens += elapsed * r.refillRate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
	r.lastRefill = now

	// If we have a token, consume it
	if r.tokens >= 1 {
		r.tokens--
		return nil
	}

	// Calculate wait time for next token
	waitTime := time.Duration((1 - r.tokens) / r.refillRate * float64(time.Second))

	r.mu.Unlock()
	select {
	case <-ctx.Done():
		r.mu.Lock()
		return ctx.Err()
	case <-time.After(waitTime):
		r.mu.Lock()
		r.tokens = 0 // We consumed the token we waited for
		return nil
	}
}

// ════════════════════════════════════════════════════════════════
// DEFAULT CLIENT
// ════════════════════════════════════════════════════════════════

// defaultClient is a package-level default client.
var defaultClient = NewClient()

// Get performs a GET request using the default client.
func Get(ctx context.Context, url string) (*Response, error) {
	return defaultClient.Get(ctx, url)
}

// GetJSON performs a GET request and decodes JSON using the default client.
func GetJSON(ctx context.Context, url string, v any) error {
	return defaultClient.GetJSON(ctx, url, v)
}
