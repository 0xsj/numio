// internal/fetch/crypto.go

package fetch

import (
	"context"
	"os"
	"strings"
	"time"
)

// ════════════════════════════════════════════════════════════════
// COINGECKO PROVIDER (Free, no API key required for basic usage)
// ════════════════════════════════════════════════════════════════

const (
	coingeckoName       = "coingecko"
	coingeckoBaseURL    = "https://api.coingecko.com/api/v3"
	coingeckoProBaseURL = "https://pro-api.coingecko.com/api/v3"
	coingeckoEnvKey     = "COINGECKO_API_KEY"
)

// CoinGeckoProvider fetches crypto prices from CoinGecko.
// Free tier: 10-30 calls/minute, no API key required.
// Pro tier: Higher limits with API key.
type CoinGeckoProvider struct {
	*BaseProvider
	baseURL string
}

// NewCoinGeckoProvider creates a new CoinGecko provider.
func NewCoinGeckoProvider() *CoinGeckoProvider {
	base := NewBaseProvider(coingeckoName, ProviderTypeCrypto)
	base.SetAPIKeyEnv(coingeckoEnvKey)
	base.SetRequireKey(false)

	// Use lower rate limit for free tier
	base.SetClient(NewClientWithOptions(
		WithRateLimit(5), // Conservative for free tier
	))

	return &CoinGeckoProvider{
		BaseProvider: base,
		baseURL:      coingeckoBaseURL,
	}
}

// FetchRates fetches current crypto prices from CoinGecko.
func (p *CoinGeckoProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	// Build comma-separated list of CoinGecko IDs
	ids := strings.Join(coingeckoIDs(), ",")
	url := p.buildURL("/simple/price?ids=" + ids + "&vs_currencies=usd")

	var resp map[string]coingeckoPriceData
	if err := p.Client().GetJSON(ctx, url, &resp); err != nil {
		return nil, p.WrapError(err)
	}

	result := NewRatesResult(p.Name(), ProviderTypeCrypto).
		SetBase("USD").
		SetSource(url)

	// Map CoinGecko IDs back to symbols
	for id, data := range resp {
		if symbol, ok := coingeckoIDToSymbol[id]; ok {
			result.AddRate(symbol, data.USD)
		}
	}

	return result, nil
}

// buildURL constructs the API URL with optional API key.
func (p *CoinGeckoProvider) buildURL(path string) string {
	apiKey := p.APIKey()
	base := p.baseURL

	if apiKey != "" {
		base = coingeckoProBaseURL
		// Pro API uses header authentication, but we can also use query param
		if strings.Contains(path, "?") {
			path += "&x_cg_pro_api_key=" + apiKey
		} else {
			path += "?x_cg_pro_api_key=" + apiKey
		}
	}

	return base + path
}

// APIKey returns the API key from environment.
func (p *CoinGeckoProvider) APIKey() string {
	return os.Getenv(coingeckoEnvKey)
}

// coingeckoPriceData represents price data for a single coin.
type coingeckoPriceData struct {
	USD float64 `json:"usd"`
}

// ════════════════════════════════════════════════════════════════
// COINCAP PROVIDER (Free, no API key required)
// ════════════════════════════════════════════════════════════════

const (
	coincapName    = "coincap"
	coincapBaseURL = "https://api.coincap.io/v2"
)

// CoinCapProvider fetches crypto prices from CoinCap.
// Free API with no authentication required.
type CoinCapProvider struct {
	*BaseProvider
	baseURL string
}

// NewCoinCapProvider creates a new CoinCap provider.
func NewCoinCapProvider() *CoinCapProvider {
	base := NewBaseProvider(coincapName, ProviderTypeCrypto)
	base.SetRequireKey(false)

	return &CoinCapProvider{
		BaseProvider: base,
		baseURL:      coincapBaseURL,
	}
}

// FetchRates fetches current crypto prices from CoinCap.
func (p *CoinCapProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	// Build comma-separated list of CoinCap IDs
	ids := strings.Join(coincapIDs(), ",")
	url := p.baseURL + "/assets?ids=" + ids

	var resp coincapResponse
	if err := p.Client().GetJSON(ctx, url, &resp); err != nil {
		return nil, p.WrapError(err)
	}

	result := NewRatesResult(p.Name(), ProviderTypeCrypto).
		SetBase("USD").
		SetSource(url).
		SetTimestamp(resp.Timestamp.Time())

	// Map CoinCap IDs back to symbols
	for _, asset := range resp.Data {
		if symbol, ok := coincapIDToSymbol[asset.ID]; ok {
			if price := parseFloat(asset.PriceUSD); price > 0 {
				result.AddRate(symbol, price)
			}
		}
	}

	return result, nil
}

// coincapResponse is the API response structure.
type coincapResponse struct {
	Data      []coincapAsset   `json:"data"`
	Timestamp coincapTimestamp `json:"timestamp"`
}

// coincapAsset represents a single asset in the response.
type coincapAsset struct {
	ID       string `json:"id"`
	Rank     string `json:"rank"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	PriceUSD string `json:"priceUsd"`
}

// coincapTimestamp handles CoinCap's millisecond timestamp.
type coincapTimestamp int64

// Time converts the timestamp to time.Time.
func (t coincapTimestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

// ════════════════════════════════════════════════════════════════
// CRYPTO PROVIDER FACTORY
// ════════════════════════════════════════════════════════════════

// NewCryptoProvider returns the default crypto provider (CoinGecko).
func NewCryptoProvider() Provider {
	return NewCoinGeckoProvider()
}

// NewCryptoProviders returns all available crypto providers in priority order.
func NewCryptoProviders() []Provider {
	return []Provider{
		NewCoinGeckoProvider(),
		NewCoinCapProvider(),
	}
}

// ════════════════════════════════════════════════════════════════
// COINGECKO ID MAPPINGS
// ════════════════════════════════════════════════════════════════

// coingeckoSymbolToID maps crypto symbols to CoinGecko IDs.
// These must match the IDs used by the CoinGecko API.
var coingeckoSymbolToID = map[string]string{
	// Tier 1 - Major
	"BTC": "bitcoin",
	"ETH": "ethereum",
	// Tier 2 - Stablecoins
	"USDT": "tether",
	"USDC": "usd-coin",
	"DAI":  "dai",
	"BUSD": "binance-usd",
	// Tier 3 - Major altcoins
	"BNB":   "binancecoin",
	"SOL":   "solana",
	"XRP":   "ripple",
	"ADA":   "cardano",
	"DOGE":  "dogecoin",
	"DOT":   "polkadot",
	"MATIC": "matic-network",
	"AVAX":  "avalanche-2",
	"LTC":   "litecoin",
	"LINK":  "chainlink",
	"ATOM":  "cosmos",
	"UNI":   "uniswap",
	"XLM":   "stellar",
	"ALGO":  "algorand",
	"TON":   "the-open-network",
	// Tier 4 - DeFi & others
	"AAVE": "aave",
	"MKR":  "maker",
	"CRV":  "curve-dao-token",
	"NEAR": "near",
	"APT":  "aptos",
	"ARB":  "arbitrum",
	"OP":   "optimism",
	// Tier 5 - Memecoins
	"SHIB": "shiba-inu",
	"PEPE": "pepe",
	"WIF":  "dogwifcoin",
	"BONK": "bonk",
}

// coingeckoIDToSymbol is the reverse mapping.
var coingeckoIDToSymbol = reverseMap(coingeckoSymbolToID)

// coingeckoIDs returns all CoinGecko IDs.
func coingeckoIDs() []string {
	ids := make([]string, 0, len(coingeckoSymbolToID))
	for _, id := range coingeckoSymbolToID {
		ids = append(ids, id)
	}
	return ids
}

// ════════════════════════════════════════════════════════════════
// COINCAP ID MAPPINGS
// ════════════════════════════════════════════════════════════════

// coincapSymbolToID maps crypto symbols to CoinCap IDs.
var coincapSymbolToID = map[string]string{
	// Tier 1 - Major
	"BTC": "bitcoin",
	"ETH": "ethereum",
	// Tier 2 - Stablecoins
	"USDT": "tether",
	"USDC": "usd-coin",
	"DAI":  "multi-collateral-dai",
	// Tier 3 - Major altcoins
	"BNB":   "binance-coin",
	"SOL":   "solana",
	"XRP":   "xrp",
	"ADA":   "cardano",
	"DOGE":  "dogecoin",
	"DOT":   "polkadot",
	"MATIC": "polygon",
	"AVAX":  "avalanche",
	"LTC":   "litecoin",
	"LINK":  "chainlink",
	"ATOM":  "cosmos",
	"UNI":   "uniswap",
	"XLM":   "stellar",
	"ALGO":  "algorand",
	// Tier 4 - DeFi
	"AAVE": "aave",
	"MKR":  "maker",
	"NEAR": "near-protocol",
	// Tier 5 - Memecoins
	"SHIB": "shiba-inu",
}

// coincapIDToSymbol is the reverse mapping.
var coincapIDToSymbol = reverseMap(coincapSymbolToID)

// coincapIDs returns all CoinCap IDs.
func coincapIDs() []string {
	ids := make([]string, 0, len(coincapSymbolToID))
	for _, id := range coincapSymbolToID {
		ids = append(ids, id)
	}
	return ids
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

// reverseMap creates a reverse mapping from value to key.
func reverseMap(m map[string]string) map[string]string {
	r := make(map[string]string, len(m))
	for k, v := range m {
		r[v] = k
	}
	return r
}

// parseFloat parses a string to float64.
// Simple implementation without strconv for consistency.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}

	var result float64
	var decimal float64 = 0
	var divisor float64 = 1
	inDecimal := false

	for _, c := range s {
		if c == '.' {
			inDecimal = true
			continue
		}
		if c >= '0' && c <= '9' {
			digit := float64(c - '0')
			if inDecimal {
				divisor *= 10
				decimal += digit / divisor
			} else {
				result = result*10 + digit
			}
		}
	}

	return result + decimal
}

// CryptoCurrencies returns a list of supported crypto symbols.
var CryptoCurrencies = []string{
	"BTC", "ETH", "USDT", "USDC", "DAI", "BUSD",
	"BNB", "SOL", "XRP", "ADA", "DOGE", "DOT",
	"MATIC", "AVAX", "LTC", "LINK", "ATOM", "UNI",
	"XLM", "ALGO", "TON", "AAVE", "MKR", "CRV",
	"NEAR", "APT", "ARB", "OP", "SHIB", "PEPE",
	"WIF", "BONK",
}
