// internal/fetch/metals.go

package fetch

import (
	"context"
	"net/http"
	"os"
	"time"
)

// ════════════════════════════════════════════════════════════════
// METALS.LIVE PROVIDER (Free, no API key required)
// ════════════════════════════════════════════════════════════════

const (
	metalsLiveName    = "metals-live"
	metalsLiveBaseURL = "https://api.metals.live/v1"
)

// MetalsLiveProvider fetches precious metal prices from metals.live.
// Free API with no authentication required.
type MetalsLiveProvider struct {
	*BaseProvider
	baseURL string
}

// NewMetalsLiveProvider creates a new Metals.live provider.
func NewMetalsLiveProvider() *MetalsLiveProvider {
	base := NewBaseProvider(metalsLiveName, ProviderTypeMetal)
	base.SetRequireKey(false)

	return &MetalsLiveProvider{
		BaseProvider: base,
		baseURL:      metalsLiveBaseURL,
	}
}

// FetchRates fetches current metal prices from Metals.live.
func (p *MetalsLiveProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	url := p.baseURL + "/spot"

	var resp []metalsLiveSpot
	if err := p.Client().GetJSON(ctx, url, &resp); err != nil {
		return nil, p.WrapError(err)
	}

	result := NewRatesResult(p.Name(), ProviderTypeMetal).
		SetBase("USD").
		SetSource(url)

	// Map metal names to ISO codes
	for _, spot := range resp {
		if code, ok := metalsLiveToISO[spot.Metal]; ok {
			result.AddRate(code, spot.Price)
		}
	}

	return result, nil
}

// metalsLiveSpot represents a single metal spot price.
type metalsLiveSpot struct {
	Metal string  `json:"metal"`
	Price float64 `json:"price"`
}

// metalsLiveToISO maps Metals.live metal names to ISO 4217 codes.
var metalsLiveToISO = map[string]string{
	"gold":      "XAU",
	"silver":    "XAG",
	"platinum":  "XPT",
	"palladium": "XPD",
}

// ════════════════════════════════════════════════════════════════
// GOLDAPI PROVIDER (Free tier with API key)
// ════════════════════════════════════════════════════════════════

const (
	goldAPIName    = "goldapi"
	goldAPIBaseURL = "https://www.goldapi.io/api"
	goldAPIEnvKey  = "GOLDAPI_KEY"
)

// GoldAPIProvider fetches precious metal prices from GoldAPI.io.
// Requires API key (free tier available with registration).
type GoldAPIProvider struct {
	*BaseProvider
	baseURL string
}

// NewGoldAPIProvider creates a new GoldAPI provider.
func NewGoldAPIProvider() *GoldAPIProvider {
	base := NewBaseProvider(goldAPIName, ProviderTypeMetal)
	base.SetAPIKeyEnv(goldAPIEnvKey)
	base.SetRequireKey(true) // API key required

	return &GoldAPIProvider{
		BaseProvider: base,
		baseURL:      goldAPIBaseURL,
	}
}

// FetchRates fetches current metal prices from GoldAPI.
func (p *GoldAPIProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	apiKey := p.APIKey()
	if apiKey == "" {
		return nil, p.WrapError(ErrUnauthorized)
	}

	result := NewRatesResult(p.Name(), ProviderTypeMetal).
		SetBase("USD")

	// GoldAPI requires separate requests per metal
	for _, metal := range []string{"XAU", "XAG", "XPT", "XPD"} {
		price, err := p.fetchMetal(ctx, apiKey, metal)
		if err != nil {
			// Continue with other metals on error
			continue
		}
		result.AddRate(metal, price)
	}

	if result.IsEmpty() {
		return nil, p.WrapError(ErrRequestFailed)
	}

	return result, nil
}

// fetchMetal fetches a single metal price.
func (p *GoldAPIProvider) fetchMetal(ctx context.Context, apiKey, metal string) (float64, error) {
	url := p.baseURL + "/" + metal + "/USD"

	// Create custom request with API key header
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("x-access-token", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := p.Client().http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, ErrRequestFailed
	}

	var data goldAPIResponse
	response := newResponse(resp)
	if err := response.JSON(&data); err != nil {
		return 0, err
	}

	return data.Price, nil
}

// APIKey returns the API key from environment.
func (p *GoldAPIProvider) APIKey() string {
	return os.Getenv(goldAPIEnvKey)
}

// goldAPIResponse is the API response structure.
type goldAPIResponse struct {
	Timestamp      int64   `json:"timestamp"`
	Metal          string  `json:"metal"`
	Currency       string  `json:"currency"`
	Exchange       string  `json:"exchange"`
	Symbol         string  `json:"symbol"`
	PrevClosePrice float64 `json:"prev_close_price"`
	OpenPrice      float64 `json:"open_price"`
	LowPrice       float64 `json:"low_price"`
	HighPrice      float64 `json:"high_price"`
	OpenTime       int64   `json:"open_time"`
	Price          float64 `json:"price"`
	Ch             float64 `json:"ch"`
	Chp            float64 `json:"chp"`
	Ask            float64 `json:"ask"`
	Bid            float64 `json:"bid"`
	PriceGram24K   float64 `json:"price_gram_24k"`
	PriceGram22K   float64 `json:"price_gram_22k"`
	PriceGram21K   float64 `json:"price_gram_21k"`
	PriceGram20K   float64 `json:"price_gram_20k"`
	PriceGram18K   float64 `json:"price_gram_18k"`
	PriceGram16K   float64 `json:"price_gram_16k"`
	PriceGram14K   float64 `json:"price_gram_14k"`
	PriceGram10K   float64 `json:"price_gram_10k"`
}

// ════════════════════════════════════════════════════════════════
// METAL PRICE API PROVIDER (Free tier with API key)
// ════════════════════════════════════════════════════════════════

const (
	metalPriceAPIName    = "metalpriceapi"
	metalPriceAPIBaseURL = "https://api.metalpriceapi.com/v1"
	metalPriceAPIEnvKey  = "METALPRICEAPI_KEY"
)

// MetalPriceAPIProvider fetches precious metal prices from MetalpriceAPI.
// Requires API key (free tier available with registration).
type MetalPriceAPIProvider struct {
	*BaseProvider
	baseURL string
}

// NewMetalPriceAPIProvider creates a new MetalpriceAPI provider.
func NewMetalPriceAPIProvider() *MetalPriceAPIProvider {
	base := NewBaseProvider(metalPriceAPIName, ProviderTypeMetal)
	base.SetAPIKeyEnv(metalPriceAPIEnvKey)
	base.SetRequireKey(true) // API key required

	return &MetalPriceAPIProvider{
		BaseProvider: base,
		baseURL:      metalPriceAPIBaseURL,
	}
}

// FetchRates fetches current metal prices from MetalpriceAPI.
func (p *MetalPriceAPIProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	apiKey := p.APIKey()
	if apiKey == "" {
		return nil, p.WrapError(ErrUnauthorized)
	}

	url := p.baseURL + "/latest?api_key=" + apiKey + "&base=USD&currencies=XAU,XAG,XPT,XPD"

	var resp metalPriceAPIResponse
	if err := p.Client().GetJSON(ctx, url, &resp); err != nil {
		return nil, p.WrapError(err)
	}

	if !resp.Success {
		return nil, p.WrapError(ErrRequestFailed)
	}

	result := NewRatesResult(p.Name(), ProviderTypeMetal).
		SetBase("USD").
		SetSource(url)

	// Parse timestamp
	if resp.Timestamp > 0 {
		result.SetTimestamp(time.Unix(resp.Timestamp, 0))
	}

	// MetalpriceAPI returns rates as "1 USD = X oz metal"
	// We need to invert to get "1 oz = X USD"
	for code, rate := range resp.Rates {
		if rate > 0 {
			result.AddRate(code, 1.0/rate)
		}
	}

	return result, nil
}

// APIKey returns the API key from environment.
func (p *MetalPriceAPIProvider) APIKey() string {
	return os.Getenv(metalPriceAPIEnvKey)
}

// metalPriceAPIResponse is the API response structure.
type metalPriceAPIResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Date      string             `json:"date"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

// ════════════════════════════════════════════════════════════════
// METAL PROVIDER FACTORY
// ════════════════════════════════════════════════════════════════

// NewMetalProvider returns the default metal provider (Metals.live).
func NewMetalProvider() Provider {
	return NewMetalsLiveProvider()
}

// NewMetalProviders returns all available metal providers in priority order.
func NewMetalProviders() []Provider {
	providers := []Provider{
		NewMetalsLiveProvider(), // Free, no key required
	}

	// Add providers that require API keys only if configured
	if goldAPI := NewGoldAPIProvider(); goldAPI.IsAvailable() {
		providers = append(providers, goldAPI)
	}

	if metalPriceAPI := NewMetalPriceAPIProvider(); metalPriceAPI.IsAvailable() {
		providers = append(providers, metalPriceAPI)
	}

	return providers
}

// ════════════════════════════════════════════════════════════════
// SUPPORTED METALS
// ════════════════════════════════════════════════════════════════

// MetalCodes returns a list of supported precious metal codes.
var MetalCodes = []string{
	"XAU", // Gold
	"XAG", // Silver
	"XPT", // Platinum
	"XPD", // Palladium
}
