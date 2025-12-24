// internal/fetch/registry.go

package fetch

import (
	"context"
	"sync"
)

// ════════════════════════════════════════════════════════════════
// REGISTRY
// ════════════════════════════════════════════════════════════════

// Registry manages providers and handles fetching with fallback support.
type Registry struct {
	mu        sync.RWMutex
	providers map[ProviderType][]Provider
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[ProviderType][]Provider),
	}
}

// DefaultRegistry creates a registry with all default providers.
func DefaultRegistry() *Registry {
	r := NewRegistry()

	// Register fiat providers
	for _, p := range NewFiatProviders() {
		r.Register(p)
	}

	// Register crypto providers
	for _, p := range NewCryptoProviders() {
		r.Register(p)
	}

	// Register metal providers
	for _, p := range NewMetalProviders() {
		r.Register(p)
	}

	return r
}

// Register adds a provider to the registry.
// Providers are stored in order of registration (first = highest priority).
func (r *Registry) Register(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()

	typ := p.Type()
	r.providers[typ] = append(r.providers[typ], p)
}

// RegisterFirst adds a provider at the beginning (highest priority).
func (r *Registry) RegisterFirst(p Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()

	typ := p.Type()
	r.providers[typ] = append([]Provider{p}, r.providers[typ]...)
}

// Providers returns all providers of a given type.
func (r *Registry) Providers(typ ProviderType) []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := r.providers[typ]
	result := make([]Provider, len(providers))
	copy(result, providers)
	return result
}

// AvailableProviders returns providers that are currently available.
func (r *Registry) AvailableProviders(typ ProviderType) []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []Provider
	for _, p := range r.providers[typ] {
		if p.IsAvailable() {
			available = append(available, p)
		}
	}
	return available
}

// AllProviders returns all registered providers across all types.
func (r *Registry) AllProviders() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []Provider
	for _, providers := range r.providers {
		all = append(all, providers...)
	}
	return all
}

// ════════════════════════════════════════════════════════════════
// FETCHING
// ════════════════════════════════════════════════════════════════

// Fetch fetches rates from the first available provider of the given type.
// Falls back to subsequent providers on failure.
func (r *Registry) Fetch(ctx context.Context, typ ProviderType) (*RatesResult, error) {
	providers := r.AvailableProviders(typ)
	if len(providers) == 0 {
		return nil, NewProviderError("registry", ErrNotFound)
	}

	var lastErr error
	for _, p := range providers {
		result, err := p.FetchRates(ctx)
		if err == nil && !result.IsEmpty() {
			return result, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, NewProviderError("registry", ErrRequestFailed)
}

// FetchFiat fetches fiat currency rates.
func (r *Registry) FetchFiat(ctx context.Context) (*RatesResult, error) {
	return r.Fetch(ctx, ProviderTypeFiat)
}

// FetchCrypto fetches cryptocurrency prices.
func (r *Registry) FetchCrypto(ctx context.Context) (*RatesResult, error) {
	return r.Fetch(ctx, ProviderTypeCrypto)
}

// FetchMetal fetches precious metal prices.
func (r *Registry) FetchMetal(ctx context.Context) (*RatesResult, error) {
	return r.Fetch(ctx, ProviderTypeMetal)
}

// FetchAll fetches rates from all provider types and merges results.
func (r *Registry) FetchAll(ctx context.Context) (*RatesResult, error) {
	result := NewRatesResult("combined", ProviderTypeFiat).
		SetBase("USD")

	var lastErr error
	var fetched int

	// Fetch fiat
	if fiat, err := r.FetchFiat(ctx); err == nil {
		result.Merge(fiat)
		fetched++
	} else {
		lastErr = err
	}

	// Fetch crypto
	if crypto, err := r.FetchCrypto(ctx); err == nil {
		result.Merge(crypto)
		fetched++
	} else {
		lastErr = err
	}

	// Fetch metals
	if metals, err := r.FetchMetal(ctx); err == nil {
		result.Merge(metals)
		fetched++
	} else {
		lastErr = err
	}

	// Return error only if nothing was fetched
	if fetched == 0 && lastErr != nil {
		return nil, lastErr
	}

	return result, nil
}

// FetchWithProvider fetches rates using a specific provider by name.
func (r *Registry) FetchWithProvider(ctx context.Context, name string) (*RatesResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, providers := range r.providers {
		for _, p := range providers {
			if p.Name() == name {
				if !p.IsAvailable() {
					return nil, NewProviderError(name, ErrUnauthorized)
				}
				return p.FetchRates(ctx)
			}
		}
	}

	return nil, NewProviderError(name, ErrNotFound)
}

// ════════════════════════════════════════════════════════════════
// PROVIDER INFO
// ════════════════════════════════════════════════════════════════

// ProviderInfo contains information about a registered provider.
type ProviderInfo struct {
	Name      string
	Type      ProviderType
	Available bool
}

// ListProviders returns information about all registered providers.
func (r *Registry) ListProviders() []ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var infos []ProviderInfo
	for _, providers := range r.providers {
		for _, p := range providers {
			infos = append(infos, ProviderInfo{
				Name:      p.Name(),
				Type:      p.Type(),
				Available: p.IsAvailable(),
			})
		}
	}
	return infos
}

// HasProvider checks if a provider with the given name is registered.
func (r *Registry) HasProvider(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, providers := range r.providers {
		for _, p := range providers {
			if p.Name() == name {
				return true
			}
		}
	}
	return false
}

// ════════════════════════════════════════════════════════════════
// DEFAULT REGISTRY
// ════════════════════════════════════════════════════════════════

// defaultRegistry is the package-level default registry.
var (
	defaultRegistry     *Registry
	defaultRegistryOnce sync.Once
)

// Default returns the default registry (lazily initialized).
func Default() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = DefaultRegistry()
	})
	return defaultRegistry
}

// ════════════════════════════════════════════════════════════════
// PACKAGE-LEVEL CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FetchFiatRates fetches fiat rates using the default registry.
func FetchFiatRates(ctx context.Context) (*RatesResult, error) {
	return Default().FetchFiat(ctx)
}

// FetchCryptoRates fetches crypto prices using the default registry.
func FetchCryptoRates(ctx context.Context) (*RatesResult, error) {
	return Default().FetchCrypto(ctx)
}

// FetchMetalRates fetches metal prices using the default registry.
func FetchMetalRates(ctx context.Context) (*RatesResult, error) {
	return Default().FetchMetal(ctx)
}

// FetchAllRates fetches all rates using the default registry.
func FetchAllRates(ctx context.Context) (*RatesResult, error) {
	return Default().FetchAll(ctx)
}
