// internal/fetch/fiat.go

package fetch

import (
	"context"
	"os"
	"strings"
	"time"
)

// ════════════════════════════════════════════════════════════════
// FRANKFURTER PROVIDER (Free, no API key required)
// ════════════════════════════════════════════════════════════════

const (
	frankfurterName    = "frankfurter"
	frankfurterBaseURL = "https://api.frankfurter.app"
)

// FrankfurterProvider fetches fiat rates from frankfurter.app.
// This is a free API with no authentication required.
// Data is sourced from the European Central Bank.
type FrankfurterProvider struct {
	*BaseProvider
	baseURL string
}

// NewFrankfurterProvider creates a new Frankfurter provider.
func NewFrankfurterProvider() *FrankfurterProvider {
	base := NewBaseProvider(frankfurterName, ProviderTypeFiat)
	base.SetRequireKey(false)

	return &FrankfurterProvider{
		BaseProvider: base,
		baseURL:      frankfurterBaseURL,
	}
}

// FetchRates fetches current fiat rates from Frankfurter.
func (p *FrankfurterProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	url := p.baseURL + "/latest?from=USD"

	var resp frankfurterResponse
	if err := p.Client().GetJSON(ctx, url, &resp); err != nil {
		return nil, p.WrapError(err)
	}

	result := NewRatesResult(p.Name(), ProviderTypeFiat).
		SetBase(resp.Base).
		SetSource(url)

	// Parse date if available
	if resp.Date != "" {
		if t, err := time.Parse("2006-01-02", resp.Date); err == nil {
			result.SetTimestamp(t)
		}
	}

	// Add USD = 1.0 (base currency)
	result.AddRate("USD", 1.0)

	// Add all rates
	for code, rate := range resp.Rates {
		result.AddRate(strings.ToUpper(code), rate)
	}

	return result, nil
}

// frankfurterResponse is the API response structure.
type frankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

// ════════════════════════════════════════════════════════════════
// EXCHANGE RATE API PROVIDER (Free tier with optional API key)
// ════════════════════════════════════════════════════════════════

const (
	exchangeRateAPIName       = "exchangerate-api"
	exchangeRateAPIBaseURL    = "https://open.er-api.com/v6/latest"
	exchangeRateAPIKeyBaseURL = "https://v6.exchangerate-api.com/v6"
	exchangeRateAPIEnvKey     = "EXCHANGERATE_API_KEY"
)

// ExchangeRateAPIProvider fetches fiat rates from exchangerate-api.com.
// Works without API key (open.er-api.com) or with key for higher limits.
type ExchangeRateAPIProvider struct {
	*BaseProvider
	baseURL string
}

// NewExchangeRateAPIProvider creates a new ExchangeRate-API provider.
func NewExchangeRateAPIProvider() *ExchangeRateAPIProvider {
	base := NewBaseProvider(exchangeRateAPIName, ProviderTypeFiat)
	base.SetAPIKeyEnv(exchangeRateAPIEnvKey)
	base.SetRequireKey(false) // Works without key, but with limits

	return &ExchangeRateAPIProvider{
		BaseProvider: base,
		baseURL:      exchangeRateAPIBaseURL,
	}
}

// FetchRates fetches current fiat rates from ExchangeRate-API.
func (p *ExchangeRateAPIProvider) FetchRates(ctx context.Context) (*RatesResult, error) {
	url := p.buildURL()

	var resp exchangeRateAPIResponse
	if err := p.Client().GetJSON(ctx, url, &resp); err != nil {
		return nil, p.WrapError(err)
	}

	// Check for API error
	if resp.Result != "success" {
		return nil, p.WrapError(ErrRequestFailed)
	}

	result := NewRatesResult(p.Name(), ProviderTypeFiat).
		SetBase(resp.BaseCode).
		SetSource(url)

	// Parse timestamp
	if resp.TimeLastUpdateUnix > 0 {
		result.SetTimestamp(time.Unix(resp.TimeLastUpdateUnix, 0))
	}

	// Add all rates
	for code, rate := range resp.Rates {
		result.AddRate(strings.ToUpper(code), rate)
	}

	return result, nil
}

// buildURL constructs the API URL, using API key if available.
func (p *ExchangeRateAPIProvider) buildURL() string {
	apiKey := p.APIKey()
	if apiKey != "" {
		return exchangeRateAPIKeyBaseURL + "/" + apiKey + "/latest/USD"
	}
	return p.baseURL + "/USD"
}

// APIKey returns the API key from environment.
func (p *ExchangeRateAPIProvider) APIKey() string {
	return os.Getenv(exchangeRateAPIEnvKey)
}

// exchangeRateAPIResponse is the API response structure.
type exchangeRateAPIResponse struct {
	Result             string             `json:"result"`
	Provider           string             `json:"provider"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	Rates              map[string]float64 `json:"rates"`
}

// ════════════════════════════════════════════════════════════════
// FIAT PROVIDER FACTORY
// ════════════════════════════════════════════════════════════════

// NewFiatProvider returns the default fiat provider (Frankfurter).
func NewFiatProvider() Provider {
	return NewFrankfurterProvider()
}

// NewFiatProviders returns all available fiat providers in priority order.
func NewFiatProviders() []Provider {
	return []Provider{
		NewFrankfurterProvider(),
		NewExchangeRateAPIProvider(),
	}
}

// ════════════════════════════════════════════════════════════════
// SUPPORTED CURRENCIES
// ════════════════════════════════════════════════════════════════

// FiatCurrencies returns a list of commonly supported fiat currency codes.
// Note: Actual availability depends on the provider.
var FiatCurrencies = []string{
	// Major currencies
	"USD", "EUR", "GBP", "JPY", "CHF", "CAD", "AUD", "NZD",
	// European
	"SEK", "NOK", "DKK", "PLN", "CZK", "HUF", "RON", "BGN",
	"HRK", "ISK", "RUB", "UAH", "TRY",
	// Americas
	"MXN", "BRL", "ARS", "CLP", "COP", "PEN",
	// Asia Pacific
	"CNY", "HKD", "TWD", "KRW", "SGD", "MYR", "THB", "IDR",
	"PHP", "VND", "INR", "PKR", "BDT",
	// Middle East & Africa
	"ILS", "AED", "SAR", "QAR", "KWD", "EGP", "ZAR", "NGN", "KES",
}
