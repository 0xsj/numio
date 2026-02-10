// internal/fetch/history.go

package fetch

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// TYPES
// ════════════════════════════════════════════════════════════════

// AssetKind distinguishes asset types for history fetching.
type AssetKind int

const (
	AssetKindCrypto AssetKind = iota
	AssetKindFiat
	AssetKindMetal
)

// HistoryRange represents a time range for price history.
type HistoryRange int

const (
	HistoryRange7d HistoryRange = iota
	HistoryRange30d
	HistoryRange90d
	HistoryRange1y
)

// String returns the display label.
func (r HistoryRange) String() string {
	switch r {
	case HistoryRange7d:
		return "7d"
	case HistoryRange30d:
		return "30d"
	case HistoryRange90d:
		return "90d"
	case HistoryRange1y:
		return "1y"
	default:
		return "7d"
	}
}

// Days returns the number of days in this range.
func (r HistoryRange) Days() int {
	switch r {
	case HistoryRange7d:
		return 7
	case HistoryRange30d:
		return 30
	case HistoryRange90d:
		return 90
	case HistoryRange1y:
		return 365
	default:
		return 7
	}
}

// Label returns a label for the popup header.
func (r HistoryRange) Label() string {
	switch r {
	case HistoryRange7d:
		return "7 days"
	case HistoryRange30d:
		return "30 days"
	case HistoryRange90d:
		return "90 days"
	case HistoryRange1y:
		return "1 year"
	default:
		return "7 days"
	}
}

// Next cycles to the next range: 7d → 30d → 90d → 1y → 7d.
func (r HistoryRange) Next() HistoryRange {
	switch r {
	case HistoryRange7d:
		return HistoryRange30d
	case HistoryRange30d:
		return HistoryRange90d
	case HistoryRange90d:
		return HistoryRange1y
	case HistoryRange1y:
		return HistoryRange7d
	default:
		return HistoryRange7d
	}
}

// Prev cycles to the previous range: 7d → 1y → 90d → 30d → 7d.
func (r HistoryRange) Prev() HistoryRange {
	switch r {
	case HistoryRange7d:
		return HistoryRange1y
	case HistoryRange30d:
		return HistoryRange7d
	case HistoryRange90d:
		return HistoryRange30d
	case HistoryRange1y:
		return HistoryRange90d
	default:
		return HistoryRange7d
	}
}

// RangeFromKey maps a key string to a HistoryRange.
// Accepts "1"=7d, "2"=30d, "3"=90d, "4"=1y.
// Returns the range and true if matched, or zero and false.
func RangeFromKey(key string) (HistoryRange, bool) {
	switch key {
	case "1":
		return HistoryRange7d, true
	case "2":
		return HistoryRange30d, true
	case "3":
		return HistoryRange90d, true
	case "4":
		return HistoryRange1y, true
	default:
		return 0, false
	}
}

// PricePoint represents a single data point.
type PricePoint struct {
	Time  time.Time
	Price float64
}

// HistoryResult holds fetched historical price data.
type HistoryResult struct {
	Asset     string
	Base      string
	Kind      AssetKind
	Range     HistoryRange
	Points    []PricePoint
	FetchedAt time.Time
}

// Prices extracts just the price values from points.
func (h *HistoryResult) Prices() []float64 {
	if h == nil {
		return nil
	}
	prices := make([]float64, len(h.Points))
	for i, p := range h.Points {
		prices[i] = p.Price
	}
	return prices
}

// ════════════════════════════════════════════════════════════════
// ASSET DETECTION
// ════════════════════════════════════════════════════════════════

// DetectAsset inspects a Value and returns the asset code and kind.
// For multi-values, uses the last sub-value.
func DetectAsset(v types.Value) (code string, kind AssetKind, ok bool) {
	// Unwrap multi-value: use last sub-value
	if v.Kind == types.ValueMulti && len(v.Values) > 0 {
		v = v.Values[len(v.Values)-1]
	}

	switch v.Kind {
	case types.ValueCrypto:
		if v.Crypto != nil {
			return v.Crypto.Code, AssetKindCrypto, true
		}
	case types.ValueCurrency:
		if v.Curr != nil {
			return v.Curr.Code, AssetKindFiat, true
		}
	case types.ValueMetal:
		if v.Metal != nil {
			return v.Metal.Code, AssetKindMetal, true
		}
	}
	return "", 0, false
}

// DetectAssetFromLine combines eval-based detection with text scanning.
// First tries the eval result. If that yields USD (base currency) or nothing,
// falls back to scanning the input text for known crypto/fiat/metal identifiers.
// Priority: crypto > metal > non-USD fiat (scans left-to-right, first match wins).
func DetectAssetFromLine(v types.Value, line string) (code string, kind AssetKind, ok bool) {
	// Try eval result first
	code, kind, ok = DetectAsset(v)
	if ok {
		// If we got a non-USD fiat or any crypto/metal, use it
		if kind != AssetKindFiat || strings.ToUpper(code) != "USD" {
			return code, kind, true
		}
		// Got USD — fallback to text scanning to find the source asset
	}

	// Text scanning fallback: tokenize and check each word
	words := tokenizeLine(line)

	// First pass: look for crypto (highest priority for charts)
	for _, w := range words {
		upper := strings.ToUpper(w)
		if _, cgOK := coingeckoSymbolToID[upper]; cgOK {
			return upper, AssetKindCrypto, true
		}
		// Also check via types registry (handles aliases like "bitcoin")
		if c := types.LookupCrypto(w); c != nil {
			return c.Code, AssetKindCrypto, true
		}
	}

	// Second pass: look for metals
	for _, w := range words {
		if m := types.LookupMetal(w); m != nil {
			return m.Code, AssetKindMetal, true
		}
	}

	// Third pass: look for non-USD fiat
	for _, w := range words {
		if c := types.LookupCurrency(w); c != nil && strings.ToUpper(c.Code) != "USD" {
			return c.Code, AssetKindFiat, true
		}
	}

	return "", 0, false
}

// tokenizeLine splits a line into word-like tokens for asset scanning.
// Strips currency symbols ($, €, £, etc.) and common noise words.
func tokenizeLine(line string) []string {
	// Remove common symbols and punctuation
	cleaned := strings.Map(func(r rune) rune {
		switch r {
		case '$', '€', '£', '¥', '₺', '₿', '₹', ',', '(', ')', '=':
			return ' '
		default:
			return r
		}
	}, line)

	raw := strings.Fields(cleaned)

	// Filter out noise words and pure numbers
	noise := map[string]bool{
		"in": true, "to": true, "as": true, "of": true,
		"price": true, "rate": true, "value": true,
	}

	words := make([]string, 0, len(raw))
	for _, w := range raw {
		lower := strings.ToLower(w)
		if noise[lower] {
			continue
		}
		// Skip pure numbers
		if isNumeric(w) {
			continue
		}
		words = append(words, w)
	}
	return words
}

// isNumeric checks if a string is a plain number (digits, dots, minus).
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && c != '.' && c != '-' {
			return false
		}
	}
	return true
}

// ════════════════════════════════════════════════════════════════
// IN-MEMORY CACHE (10-minute TTL)
// ════════════════════════════════════════════════════════════════

type historyCacheEntry struct {
	result    *HistoryResult
	fetchedAt time.Time
}

var (
	historyCache    = make(map[string]*historyCacheEntry)
	historyCacheMu  sync.RWMutex
	historyCacheTTL = 10 * time.Minute
)

// cacheKey builds a cache key like "BTC:30".
func cacheKey(asset string, days int) string {
	return strings.ToUpper(asset) + ":" + itoa(days)
}

// itoa converts int to string without fmt.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func getCached(key string) (*HistoryResult, bool) {
	historyCacheMu.RLock()
	defer historyCacheMu.RUnlock()

	entry, ok := historyCache[key]
	if !ok {
		return nil, false
	}
	if time.Since(entry.fetchedAt) > historyCacheTTL {
		return nil, false
	}
	return entry.result, true
}

func putCache(key string, result *HistoryResult) {
	historyCacheMu.Lock()
	defer historyCacheMu.Unlock()

	historyCache[key] = &historyCacheEntry{
		result:    result,
		fetchedAt: time.Now(),
	}
}

// ════════════════════════════════════════════════════════════════
// FETCH HISTORY
// ════════════════════════════════════════════════════════════════

// FetchHistory fetches historical price data for an asset.
// Uses in-memory cache with 10-minute TTL.
func FetchHistory(ctx context.Context, asset string, kind AssetKind, hr HistoryRange) (*HistoryResult, error) {
	key := cacheKey(asset, hr.Days())

	// Check cache
	if cached, ok := getCached(key); ok {
		return cached, nil
	}

	var result *HistoryResult
	var err error

	switch kind {
	case AssetKindCrypto:
		result, err = FetchCryptoHistory(ctx, asset, hr)
	case AssetKindFiat:
		result, err = FetchFiatHistory(ctx, asset, hr)
	default:
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	putCache(key, result)
	return result, nil
}

// ════════════════════════════════════════════════════════════════
// CRYPTO HISTORY (CoinGecko /market_chart)
// ════════════════════════════════════════════════════════════════

// coingeckoMarketChartResponse is the CoinGecko /market_chart response.
type coingeckoMarketChartResponse struct {
	Prices [][]float64 `json:"prices"` // [[timestamp_ms, price], ...]
}

// CoinGeckoID returns the CoinGecko API ID for a crypto symbol.
func CoinGeckoID(symbol string) (string, bool) {
	id, ok := coingeckoSymbolToID[strings.ToUpper(symbol)]
	return id, ok
}

// FetchCryptoHistory fetches historical crypto prices from CoinGecko.
func FetchCryptoHistory(ctx context.Context, symbol string, hr HistoryRange) (*HistoryResult, error) {
	cgID, ok := CoinGeckoID(symbol)
	if !ok {
		return nil, ErrNotFound
	}

	url := coingeckoBaseURL + "/coins/" + cgID + "/market_chart?vs_currency=usd&days=" + itoa(hr.Days())

	client := NewClientWithOptions(
		WithRateLimit(5),
	)

	var resp coingeckoMarketChartResponse
	if err := client.GetJSON(ctx, url, &resp); err != nil {
		return nil, err
	}

	points := make([]PricePoint, 0, len(resp.Prices))
	for _, p := range resp.Prices {
		if len(p) < 2 {
			continue
		}
		points = append(points, PricePoint{
			Time:  time.UnixMilli(int64(p[0])),
			Price: p[1],
		})
	}

	return &HistoryResult{
		Asset:     strings.ToUpper(symbol),
		Base:      "USD",
		Kind:      AssetKindCrypto,
		Range:     hr,
		Points:    points,
		FetchedAt: time.Now(),
	}, nil
}

// ════════════════════════════════════════════════════════════════
// FIAT HISTORY (Frankfurter time-series)
// ════════════════════════════════════════════════════════════════

// frankfurterTimeSeriesResponse is the Frankfurter time-series response.
type frankfurterTimeSeriesResponse struct {
	Base    string                        `json:"base"`
	StartAt string                        `json:"start_date"`
	EndAt   string                        `json:"end_date"`
	Rates   map[string]map[string]float64 `json:"rates"` // {"2024-01-15": {"EUR": 0.92}}
}

// datePrice pairs a date string with a price for sorting.
type datePrice struct {
	date  string
	price float64
}

// FetchFiatHistory fetches historical fiat rates from Frankfurter.
func FetchFiatHistory(ctx context.Context, code string, hr HistoryRange) (*HistoryResult, error) {
	code = strings.ToUpper(code)

	// USD is the base currency — can't chart it
	if code == "USD" {
		return nil, ErrNotFound
	}

	end := time.Now()
	start := end.AddDate(0, 0, -hr.Days())

	url := frankfurterBaseURL + "/" + start.Format("2006-01-02") + ".." + end.Format("2006-01-02") + "?from=USD&to=" + code

	client := NewClient()

	var resp frankfurterTimeSeriesResponse
	if err := client.GetJSON(ctx, url, &resp); err != nil {
		return nil, err
	}

	// Collect and sort by date
	var datePrices []datePrice
	for date, rates := range resp.Rates {
		if rate, ok := rates[code]; ok {
			datePrices = append(datePrices, datePrice{date: date, price: rate})
		}
	}

	// Sort by date string (ISO format sorts correctly)
	sortDatePrices(datePrices)

	points := make([]PricePoint, 0, len(datePrices))
	for _, dp := range datePrices {
		t, err := time.Parse("2006-01-02", dp.date)
		if err != nil {
			continue
		}
		points = append(points, PricePoint{
			Time:  t,
			Price: dp.price,
		})
	}

	return &HistoryResult{
		Asset:     code,
		Base:      "USD",
		Kind:      AssetKindFiat,
		Range:     hr,
		Points:    points,
		FetchedAt: time.Now(),
	}, nil
}

// ════════════════════════════════════════════════════════════════
// FORECASTING (Linear Regression)
// ════════════════════════════════════════════════════════════════

// Forecast holds price projection data.
type Forecast struct {
	// Projected prices for future periods
	Points []float64
	// Projected final price
	Target float64
	// Change from last known price to projected price
	ChangePct float64
	// R² (coefficient of determination) — 0..1, measures fit quality
	RSquared float64
	// Number of periods projected forward
	Periods int
	// Label for the forecast horizon (e.g. "7d")
	Label string
}

// Confidence returns a human-readable confidence level based on R².
func (f *Forecast) Confidence() string {
	if f.RSquared >= 0.8 {
		return "high"
	}
	if f.RSquared >= 0.5 {
		return "moderate"
	}
	return "low"
}

// LinearForecast projects future prices using linear regression.
// Projects forward ~20% of the data length (e.g. 7 points for 30-day data).
func LinearForecast(prices []float64) *Forecast {
	n := len(prices)
	if n < 3 {
		return nil
	}

	// Linear regression: y = slope*x + intercept
	// Using least squares
	var sumX, sumY, sumXY, sumX2 float64
	for i, y := range prices {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	nf := float64(n)
	denom := nf*sumX2 - sumX*sumX
	if denom == 0 {
		return nil
	}

	slope := (nf*sumXY - sumX*sumY) / denom
	intercept := (sumY - slope*sumX) / nf

	// R² calculation
	meanY := sumY / nf
	var ssTot, ssRes float64
	for i, y := range prices {
		predicted := slope*float64(i) + intercept
		ssRes += (y - predicted) * (y - predicted)
		ssTot += (y - meanY) * (y - meanY)
	}
	rSquared := 0.0
	if ssTot > 0 {
		rSquared = 1 - ssRes/ssTot
	}
	if rSquared < 0 {
		rSquared = 0
	}

	// Project forward ~20% of data length, min 3, max 10
	periods := n / 5
	if periods < 3 {
		periods = 3
	}
	if periods > 10 {
		periods = 10
	}

	points := make([]float64, periods)
	for i := range points {
		points[i] = slope*float64(n+i) + intercept
		// Don't let prices go negative
		if points[i] < 0 {
			points[i] = 0
		}
	}

	lastPrice := prices[n-1]
	target := points[periods-1]
	changePct := 0.0
	if lastPrice > 0 {
		changePct = ((target - lastPrice) / lastPrice) * 100
	}

	// Label based on data length
	label := itoa(periods) + " periods"
	switch {
	case n <= 10:
		label = itoa(periods) + "d"
	case n <= 35:
		label = itoa(periods) + "d"
	case n <= 100:
		label = itoa(periods) + "d"
	default:
		label = itoa(periods*365/n) + "mo"
	}

	return &Forecast{
		Points:    points,
		Target:    target,
		ChangePct: changePct,
		RSquared:  rSquared,
		Periods:   periods,
		Label:     label,
	}
}

// sortDatePrices sorts by date string (insertion sort, small datasets).
func sortDatePrices(dp []datePrice) {
	for i := 1; i < len(dp); i++ {
		key := dp[i]
		j := i - 1
		for j >= 0 && dp[j].date > key.date {
			dp[j+1] = dp[j]
			j--
		}
		dp[j+1] = key
	}
}
