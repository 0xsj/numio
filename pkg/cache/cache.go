// internal/cache/cache.go

// Package cache provides exchange rate caching with BFS path-finding.
package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/0xsj/numio/pkg/types"
)

// Default cache settings
const (
	DefaultTTL        = 1 * time.Hour
	DefaultMemoryTTL  = 5 * time.Minute
	DefaultCacheDir   = ".numio/cache"
	DefaultRatesFile  = "rates.json"
)

// RateCache stores exchange rates with multiple cache layers.
type RateCache struct {
	mu sync.RWMutex

	// Direct rates: (from, to) -> rate
	rates map[ratePair]float64

	// Raw rates from API (for persistence)
	rawRates map[string]float64

	// Timestamps
	lastUpdate time.Time
	ttl        time.Duration

	// File cache path
	cacheDir  string
	cacheFile string
}

// ratePair represents a currency pair for rate lookup.
type ratePair struct {
	From string
	To   string
}

// CachedRates represents the JSON structure for file persistence.
type CachedRates struct {
	Timestamp int64              `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
	BaseCurrency string          `json:"base_currency"`
}

// New creates a new RateCache with default settings.
func New() *RateCache {
	c := &RateCache{
		rates:     make(map[ratePair]float64),
		rawRates:  make(map[string]float64),
		ttl:       DefaultTTL,
		cacheDir:  getCacheDir(),
		cacheFile: DefaultRatesFile,
	}

	// Load defaults first
	c.loadDefaults()

	// Try to load from file cache
	c.LoadFromFile()

	return c
}

// NewWithTTL creates a RateCache with custom TTL.
func NewWithTTL(ttl time.Duration) *RateCache {
	c := New()
	c.ttl = ttl
	return c
}

// ════════════════════════════════════════════════════════════════
// RATE OPERATIONS
// ════════════════════════════════════════════════════════════════

// SetRate sets an exchange rate between two currencies.
func (c *RateCache) SetRate(from, to string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	c.rates[ratePair{From: from, To: to}] = rate

	// Also store inverse rate
	if rate != 0 {
		c.rates[ratePair{From: to, To: from}] = 1.0 / rate
	}
}

// GetRate gets the exchange rate between two currencies.
// Uses BFS to find conversion path if direct rate not available.
func (c *RateCache) GetRate(from, to string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	// Same currency
	if from == to {
		return 1.0, true
	}

	// Try direct rate
	if rate, ok := c.rates[ratePair{From: from, To: to}]; ok {
		return rate, true
	}

	// Try BFS to find conversion path
	return c.findRateBFS(from, to)
}

// findRateBFS uses breadth-first search to find a conversion path.
func (c *RateCache) findRateBFS(from, to string) (float64, bool) {
	// Queue entries: (currency, accumulated rate)
	type queueEntry struct {
		currency string
		rate     float64
	}

	visited := make(map[string]bool)
	queue := []queueEntry{{currency: from, rate: 1.0}}
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Find all currencies we can reach from current
		for pair, rate := range c.rates {
			if pair.From != current.currency {
				continue
			}

			nextCurrency := pair.To
			nextRate := current.rate * rate

			// Found target
			if nextCurrency == to {
				return nextRate, true
			}

			// Add to queue if not visited
			if !visited[nextCurrency] {
				visited[nextCurrency] = true
				queue = append(queue, queueEntry{currency: nextCurrency, rate: nextRate})
			}
		}
	}

	return 0, false
}

// HasRate checks if a rate exists (direct or via path).
func (c *RateCache) HasRate(from, to string) bool {
	_, ok := c.GetRate(from, to)
	return ok
}

// Clear removes all cached rates.
func (c *RateCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rates = make(map[ratePair]float64)
	c.rawRates = make(map[string]float64)
	c.lastUpdate = time.Time{}
}

// ════════════════════════════════════════════════════════════════
// BULK OPERATIONS
// ════════════════════════════════════════════════════════════════

// ApplyRawRates applies rates from an API response.
// Fiat rates: "1 USD = X currency"
// Crypto rates: "1 TOKEN = X USD"
func (c *RateCache) ApplyRawRates(rates map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Store raw rates for persistence
	c.rawRates = make(map[string]float64)
	for k, v := range rates {
		c.rawRates[k] = v
	}

	// Process rates
	for code, rate := range rates {
		code = strings.ToUpper(code)

		// Check if it's a crypto (rate is in USD)
		if types.IsCrypto(code) {
			// Crypto: 1 TOKEN = rate USD
			c.rates[ratePair{From: code, To: "USD"}] = rate
			if rate != 0 {
				c.rates[ratePair{From: "USD", To: code}] = 1.0 / rate
			}
		} else {
			// Fiat: 1 USD = rate CURRENCY
			c.rates[ratePair{From: "USD", To: code}] = rate
			if rate != 0 {
				c.rates[ratePair{From: code, To: "USD"}] = 1.0 / rate
			}
		}
	}

	c.lastUpdate = time.Now()
}

// RawRates returns the raw rates map (for persistence).
func (c *RateCache) RawRates() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]float64)
	for k, v := range c.rawRates {
		result[k] = v
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// CACHE VALIDITY
// ════════════════════════════════════════════════════════════════

// IsExpired returns true if the cache has expired.
func (c *RateCache) IsExpired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastUpdate.IsZero() {
		return true
	}
	return time.Since(c.lastUpdate) > c.ttl
}

// IsValid returns true if the cache is valid (not expired).
func (c *RateCache) IsValid() bool {
	return !c.IsExpired()
}

// LastUpdate returns the last update time.
func (c *RateCache) LastUpdate() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUpdate
}

// Age returns the age of the cache.
func (c *RateCache) Age() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastUpdate.IsZero() {
		return 0
	}
	return time.Since(c.lastUpdate)
}

// ════════════════════════════════════════════════════════════════
// FILE PERSISTENCE
// ════════════════════════════════════════════════════════════════

// LoadFromFile loads rates from the file cache.
func (c *RateCache) LoadFromFile() bool {
	path := c.getCachePath()
	if path == "" {
		return false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var cached CachedRates
	if err := json.Unmarshal(data, &cached); err != nil {
		return false
	}

	// Check if expired
	timestamp := time.Unix(cached.Timestamp, 0)
	if time.Since(timestamp) > c.ttl {
		return false
	}

	// Apply rates
	c.ApplyRawRates(cached.Rates)
	c.mu.Lock()
	c.lastUpdate = timestamp
	c.mu.Unlock()

	return true
}

// SaveToFile saves rates to the file cache.
func (c *RateCache) SaveToFile() error {
	path := c.getCachePath()
	if path == "" {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	c.mu.RLock()
	cached := CachedRates{
		Timestamp:    c.lastUpdate.Unix(),
		Rates:        c.rawRates,
		BaseCurrency: "USD",
	}
	c.mu.RUnlock()

	data, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// getCachePath returns the full path to the cache file.
func (c *RateCache) getCachePath() string {
	if c.cacheDir == "" {
		return ""
	}
	return filepath.Join(c.cacheDir, c.cacheFile)
}

// getCacheDir returns the cache directory path.
func getCacheDir() string {
	// Try XDG_CACHE_HOME first
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "numio")
	}

	// Fall back to ~/.numio/cache
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, DefaultCacheDir)
}

// IsCacheFileValid checks if the cache file exists and is not expired.
func IsCacheFileValid() bool {
	c := &RateCache{
		cacheDir:  getCacheDir(),
		cacheFile: DefaultRatesFile,
		ttl:       DefaultTTL,
	}
	return c.isCacheFileValid()
}

func (c *RateCache) isCacheFileValid() bool {
	path := c.getCachePath()
	if path == "" {
		return false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var cached CachedRates
	if err := json.Unmarshal(data, &cached); err != nil {
		return false
	}

	timestamp := time.Unix(cached.Timestamp, 0)
	return time.Since(timestamp) <= c.ttl
}

// ════════════════════════════════════════════════════════════════
// DEFAULT RATES (Fallback for offline mode)
// ════════════════════════════════════════════════════════════════

// loadDefaults loads hardcoded fallback rates.
func (c *RateCache) loadDefaults() {
	// Major fiat currencies (approximate rates vs USD)
	fiatDefaults := map[string]float64{
		"EUR": 0.92,
		"GBP": 0.79,
		"JPY": 149.50,
		"CHF": 0.88,
		"CAD": 1.36,
		"AUD": 1.53,
		"CNY": 7.24,
		"INR": 83.12,
		"KRW": 1320.0,
		"MXN": 17.15,
		"BRL": 4.97,
		"RUB": 92.50,
		"TRY": 32.50,
		"ZAR": 18.65,
		"SGD": 1.34,
		"HKD": 7.82,
		"NOK": 10.65,
		"SEK": 10.42,
		"DKK": 6.87,
		"PLN": 3.98,
		"THB": 35.20,
		"IDR": 15650.0,
		"MYR": 4.72,
		"PHP": 55.80,
		"CZK": 22.85,
		"ILS": 3.72,
		"AED": 3.67,
		"SAR": 3.75,
		"TWD": 31.50,
		"HUF": 355.0,
		"UAH": 41.0,
		"VND": 24500.0,
		"EGP": 30.90,
		"PKR": 285.0,
		"BDT": 110.0,
		"NGN": 800.0,
		"ARS": 850.0,
		"CLP": 880.0,
		"COP": 3950.0,
		"KES": 155.0,
		"QAR": 3.64,
		"KWD": 0.31,
		"RON": 4.57,
		"NZD": 1.64,
	}

	// Crypto rates (approximate USD values)
	cryptoDefaults := map[string]float64{
		"BTC":   95000.0,
		"ETH":   3500.0,
		"SOL":   180.0,
		"BNB":   650.0,
		"XRP":   2.20,
		"ADA":   0.95,
		"DOGE":  0.38,
		"DOT":   7.50,
		"MATIC": 0.55,
		"LTC":   105.0,
		"LINK":  22.0,
		"AVAX":  42.0,
		"ATOM":  9.50,
		"UNI":   12.0,
		"XLM":   0.42,
		"ALGO":  0.35,
		"TON":   6.50,
		"NEAR":  5.80,
		"APT":   12.50,
		"ARB":   1.10,
		"OP":    2.80,
		"AAVE":  350.0,
		"MKR":   3200.0,
		"CRV":   0.85,
		"SHIB":  0.000025,
		"PEPE":  0.000018,
		"WIF":   2.50,
		"BONK":  0.000032,
		"USDT":  1.0,
		"USDC":  1.0,
		"DAI":   1.0,
		"BUSD":  1.0,
	}

	// Metal rates (USD per troy oz)
	metalDefaults := map[string]float64{
		"XAU": 2650.0,  // Gold
		"XAG": 31.50,   // Silver
		"XPT": 1020.0,  // Platinum
		"XPD": 1100.0,  // Palladium
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// USD to itself
	c.rates[ratePair{From: "USD", To: "USD"}] = 1.0

	// Fiat: 1 USD = X currency
	for code, rate := range fiatDefaults {
		c.rates[ratePair{From: "USD", To: code}] = rate
		if rate != 0 {
			c.rates[ratePair{From: code, To: "USD"}] = 1.0 / rate
		}
	}

	// Crypto: 1 TOKEN = X USD
	for code, rate := range cryptoDefaults {
		c.rates[ratePair{From: code, To: "USD"}] = rate
		if rate != 0 {
			c.rates[ratePair{From: "USD", To: code}] = 1.0 / rate
		}
	}

	// Metals: 1 oz = X USD
	for code, rate := range metalDefaults {
		c.rates[ratePair{From: code, To: "USD"}] = rate
		if rate != 0 {
			c.rates[ratePair{From: "USD", To: code}] = 1.0 / rate
		}
	}
}

// ════════════════════════════════════════════════════════════════
// STATISTICS
// ════════════════════════════════════════════════════════════════

// Stats returns cache statistics.
type Stats struct {
	DirectRates  int
	LastUpdate   time.Time
	Age          time.Duration
	IsExpired    bool
	CacheFile    string
	HasFileCache bool
}

// Stats returns cache statistics.
func (c *RateCache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	path := c.getCachePath()
	_, err := os.Stat(path)

	return Stats{
		DirectRates:  len(c.rates),
		LastUpdate:   c.lastUpdate,
		Age:          time.Since(c.lastUpdate),
		IsExpired:    c.lastUpdate.IsZero() || time.Since(c.lastUpdate) > c.ttl,
		CacheFile:    path,
		HasFileCache: err == nil,
	}
}

// ════════════════════════════════════════════════════════════════
// CONVERSION HELPERS
// ════════════════════════════════════════════════════════════════

// Convert converts an amount from one currency to another.
func (c *RateCache) Convert(amount float64, from, to string) (float64, bool) {
	rate, ok := c.GetRate(from, to)
	if !ok {
		return 0, false
	}
	return amount * rate, true
}

// ConvertValue converts a types.Value to a target currency/unit.
func (c *RateCache) ConvertValue(v types.Value, target string) (types.Value, bool) {
	if v.IsError() || v.IsEmpty() {
		return v, false
	}

	target = strings.ToUpper(target)

	switch v.Kind {
	case types.ValueCurrency:
		if v.Curr == nil {
			return v, false
		}
		converted, ok := c.Convert(v.Num, v.Curr.Code, target)
		if !ok {
			return v, false
		}
		targetCurr := types.ParseCurrency(target)
		if targetCurr == nil {
			targetCurr = types.CurrencyFromCode(target)
		}
		return types.CurrencyValue(converted, targetCurr), true

	case types.ValueCrypto:
		if v.Crypto == nil {
			return v, false
		}
		converted, ok := c.Convert(v.Num, v.Crypto.Code, target)
		if !ok {
			return v, false
		}
		// Target could be currency or crypto
		if targetCrypto := types.ParseCrypto(target); targetCrypto != nil {
			return types.CryptoValue(converted, targetCrypto), true
		}
		if targetCurr := types.ParseCurrency(target); targetCurr != nil {
			return types.CurrencyValue(converted, targetCurr), true
		}
		return types.Number(converted), true

	case types.ValueMetal:
		if v.Metal == nil {
			return v, false
		}
		converted, ok := c.Convert(v.Num, v.Metal.Code, target)
		if !ok {
			return v, false
		}
		if targetCurr := types.ParseCurrency(target); targetCurr != nil {
			return types.CurrencyValue(converted, targetCurr), true
		}
		return types.Number(converted), true

	case types.ValueWithUnit:
		// Unit conversion handled elsewhere
		return v, false

	default:
		return v, false
	}
}