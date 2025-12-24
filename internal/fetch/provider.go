// internal/fetch/provider.go

package fetch

import (
	"context"
	"time"
)

// ════════════════════════════════════════════════════════════════
// PROVIDER INTERFACE
// ════════════════════════════════════════════════════════════════

// Provider defines the interface for fetching exchange rates or asset prices.
type Provider interface {
	// Name returns the provider's identifier (e.g., "open-exchange-rates", "coingecko").
	Name() string

	// Type returns the type of assets this provider handles.
	Type() ProviderType

	// FetchRates fetches current rates and returns them as a map.
	// The map format depends on the provider type:
	//   - Fiat: base currency to rates (e.g., {"EUR": 0.92, "GBP": 0.79})
	//   - Crypto: symbol to USD price (e.g., {"BTC": 95000, "ETH": 3500})
	//   - Metal: symbol to USD price per oz (e.g., {"XAU": 2650, "XAG": 31.5})
	FetchRates(ctx context.Context) (*RatesResult, error)

	// IsAvailable checks if the provider is currently available.
	// This can be used to skip providers that require API keys if none is configured.
	IsAvailable() bool
}

// ProviderType identifies the category of a provider.
type ProviderType int

const (
	ProviderTypeFiat ProviderType = iota
	ProviderTypeCrypto
	ProviderTypeMetal
	ProviderTypeStock     // Future
	ProviderTypeCommodity // Future
)

// String returns the provider type name.
func (t ProviderType) String() string {
	switch t {
	case ProviderTypeFiat:
		return "fiat"
	case ProviderTypeCrypto:
		return "crypto"
	case ProviderTypeMetal:
		return "metal"
	case ProviderTypeStock:
		return "stock"
	case ProviderTypeCommodity:
		return "commodity"
	default:
		return "unknown"
	}
}

// ════════════════════════════════════════════════════════════════
// RATES RESULT
// ════════════════════════════════════════════════════════════════

// RatesResult holds the result of a rate fetch operation.
type RatesResult struct {
	// Provider name that returned this result
	Provider string

	// Type of rates (fiat, crypto, metal)
	Type ProviderType

	// BaseCurrency is the base currency for fiat rates (usually "USD")
	BaseCurrency string

	// Rates maps currency/asset code to rate/price
	// For fiat: 1 BaseCurrency = X target (e.g., {"EUR": 0.92} means 1 USD = 0.92 EUR)
	// For crypto: 1 token = X USD (e.g., {"BTC": 95000} means 1 BTC = 95000 USD)
	// For metal: 1 oz = X USD (e.g., {"XAU": 2650} means 1 oz gold = 2650 USD)
	Rates map[string]float64

	// Timestamp when the rates were fetched/updated
	Timestamp time.Time

	// Source URL or identifier (for debugging/attribution)
	Source string
}

// NewRatesResult creates a new RatesResult.
func NewRatesResult(provider string, providerType ProviderType) *RatesResult {
	return &RatesResult{
		Provider:  provider,
		Type:      providerType,
		Rates:     make(map[string]float64),
		Timestamp: time.Now(),
	}
}

// SetBase sets the base currency (for fiat providers).
func (r *RatesResult) SetBase(base string) *RatesResult {
	r.BaseCurrency = base
	return r
}

// SetSource sets the source URL.
func (r *RatesResult) SetSource(source string) *RatesResult {
	r.Source = source
	return r
}

// SetTimestamp sets the timestamp.
func (r *RatesResult) SetTimestamp(t time.Time) *RatesResult {
	r.Timestamp = t
	return r
}

// AddRate adds a rate to the result.
func (r *RatesResult) AddRate(code string, rate float64) *RatesResult {
	r.Rates[code] = rate
	return r
}

// Merge merges another RatesResult into this one.
// Existing rates are overwritten by the other result.
func (r *RatesResult) Merge(other *RatesResult) *RatesResult {
	if other == nil {
		return r
	}
	for code, rate := range other.Rates {
		r.Rates[code] = rate
	}
	return r
}

// Count returns the number of rates.
func (r *RatesResult) Count() int {
	return len(r.Rates)
}

// IsEmpty returns true if there are no rates.
func (r *RatesResult) IsEmpty() bool {
	return len(r.Rates) == 0
}

// Get returns a rate by code, or 0 and false if not found.
func (r *RatesResult) Get(code string) (float64, bool) {
	rate, ok := r.Rates[code]
	return rate, ok
}

// ════════════════════════════════════════════════════════════════
// PROVIDER ERRORS
// ════════════════════════════════════════════════════════════════

// ProviderError represents an error from a specific provider.
type ProviderError struct {
	Provider string
	Err      error
}

// Error implements the error interface.
func (e *ProviderError) Error() string {
	return e.Provider + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new ProviderError.
func NewProviderError(provider string, err error) *ProviderError {
	return &ProviderError{
		Provider: provider,
		Err:      err,
	}
}

// ════════════════════════════════════════════════════════════════
// BASE PROVIDER (for embedding)
// ════════════════════════════════════════════════════════════════

// BaseProvider provides common functionality for providers.
// Embed this in concrete provider implementations.
type BaseProvider struct {
	name       string
	typ        ProviderType
	client     *Client
	apiKey     string
	apiKeyEnv  string // Environment variable name for API key
	requireKey bool   // Whether API key is required
}

// NewBaseProvider creates a new BaseProvider.
func NewBaseProvider(name string, typ ProviderType) *BaseProvider {
	return &BaseProvider{
		name:   name,
		typ:    typ,
		client: NewClient(),
	}
}

// Name returns the provider name.
func (p *BaseProvider) Name() string {
	return p.name
}

// Type returns the provider type.
func (p *BaseProvider) Type() ProviderType {
	return p.typ
}

// Client returns the HTTP client.
func (p *BaseProvider) Client() *Client {
	return p.client
}

// SetClient sets a custom HTTP client.
func (p *BaseProvider) SetClient(c *Client) {
	if c != nil {
		p.client = c
	}
}

// SetAPIKey sets the API key directly.
func (p *BaseProvider) SetAPIKey(key string) {
	p.apiKey = key
}

// SetAPIKeyEnv sets the environment variable name for the API key.
func (p *BaseProvider) SetAPIKeyEnv(envVar string) {
	p.apiKeyEnv = envVar
}

// SetRequireKey sets whether an API key is required.
func (p *BaseProvider) SetRequireKey(required bool) {
	p.requireKey = required
}

// APIKey returns the API key (from direct setting or environment).
func (p *BaseProvider) APIKey() string {
	if p.apiKey != "" {
		return p.apiKey
	}
	if p.apiKeyEnv != "" {
		return getEnv(p.apiKeyEnv)
	}
	return ""
}

// IsAvailable returns true if the provider can be used.
// Returns false if an API key is required but not configured.
func (p *BaseProvider) IsAvailable() bool {
	if p.requireKey {
		return p.APIKey() != ""
	}
	return true
}

// WrapError wraps an error with the provider name.
func (p *BaseProvider) WrapError(err error) error {
	return NewProviderError(p.name, err)
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

// getEnv gets an environment variable value.
// Implemented without os package to keep imports minimal in this file.
// The actual implementation will use os.Getenv.
var getEnv = defaultGetEnv

// defaultGetEnv is a placeholder that will be replaced.
// We'll use os.Getenv in the actual providers.
func defaultGetEnv(key string) string {
	// This will be properly implemented when we import "os"
	// For now, return empty (providers will import os themselves)
	return ""
}

// SetEnvFunc allows setting a custom env lookup function (useful for testing).
func SetEnvFunc(fn func(string) string) {
	if fn != nil {
		getEnv = fn
	}
}
