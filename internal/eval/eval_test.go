// internal/eval/eval_test.go

package eval

import (
	"strings"
	"testing"

	"github.com/0xsj/numio/internal/parser"
	"github.com/0xsj/numio/pkg/types"
)

// TestEvalConversion tests currency conversion evaluation.
func TestEvalConversion(t *testing.T) {
	// Create evaluator with a mock rate cache adapter
	ctx := NewContext()
	ctx.SetRateCacheAdapter(&mockRateCache{
		rates: map[string]map[string]float64{
			"USD": {"TRY": 32.5, "EUR": 0.92},
			"TRY": {"USD": 1.0 / 32.5},
			"EUR": {"USD": 1.0 / 0.92},
			"BTC": {"USD": 95000.0},
			"XAU": {"USD": 2650.0},
		},
	})
	eval := NewWithContext(ctx)

	tests := []struct {
		input       string
		wantKind    types.ValueKind
		wantApprox  float64 // approximate expected value
		wantErr     bool
		description string
	}{
		// Symbol-prefixed conversions
		{"$100 to EUR", types.ValueCurrency, 92.0, false, "USD to EUR via symbol"},
		{"$100 to lira", types.ValueCurrency, 3250.0, false, "USD to TRY via alias"},
		{"$100 in turkish lira", types.ValueCurrency, 3250.0, false, "USD to TRY via multi-word alias"},

		// Code-suffixed conversions (the original bug)
		{"100 usd to EUR", types.ValueCurrency, 92.0, false, "USD code to EUR code"},
		{"100 usd to lira", types.ValueCurrency, 3250.0, false, "USD code to TRY alias"},
		{"100 USD to turkish lira", types.ValueCurrency, 3250.0, false, "USD code to TRY multi-word"},

		// Alias-suffixed conversions
		{"100 dollars to euros", types.ValueCurrency, 92.0, false, "USD alias to EUR alias"},
		{"100 dollars to lira", types.ValueCurrency, 3250.0, false, "USD alias to TRY alias"},

		// Crypto conversions
		{"1 btc to usd", types.ValueCurrency, 95000.0, false, "BTC to USD"},
		{"1 bitcoin to dollars", types.ValueCurrency, 95000.0, false, "BTC alias to USD alias"},

		// Metal conversions
		{"1 xau to usd", types.ValueCurrency, 2650.0, false, "XAU to USD"},
		{"1 gold to dollars", types.ValueCurrency, 2650.0, false, "gold alias to USD alias"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Parse
			line, errs := parser.ParseLine(tt.input)
			if len(errs) > 0 {
				t.Fatalf("parse errors: %v", errs)
			}

			// Evaluate
			result := eval.EvalLine(line)

			if tt.wantErr {
				if !result.IsError() {
					t.Errorf("expected error, got %v", result)
				}
				return
			}

			if result.IsError() {
				t.Errorf("unexpected error: %v", result)
				return
			}

			if result.Kind != tt.wantKind {
				t.Errorf("kind: got %v, want %v", result.Kind, tt.wantKind)
			}

			// Check approximate value (within 1% tolerance)
			got := result.AsFloat()
			tolerance := tt.wantApprox * 0.01
			if got < tt.wantApprox-tolerance || got > tt.wantApprox+tolerance {
				t.Errorf("value: got %v, want ~%v", got, tt.wantApprox)
			}
		})
	}
}

// mockRateCache implements RateCacheAdapter for testing.
type mockRateCache struct {
	rates map[string]map[string]float64
}

func (m *mockRateCache) GetRate(from, to string) (float64, bool) {
	// Resolve aliases to codes
	from = resolveTestCode(from)
	to = resolveTestCode(to)

	if from == to {
		return 1.0, true
	}

	if fromRates, ok := m.rates[from]; ok {
		if rate, ok := fromRates[to]; ok {
			return rate, true
		}
	}

	// Try via USD
	if from != "USD" && to != "USD" {
		if fromRates, ok := m.rates[from]; ok {
			if toUSD, ok := fromRates["USD"]; ok {
				if usdRates, ok := m.rates["USD"]; ok {
					if usdToTarget, ok := usdRates[to]; ok {
						return toUSD * usdToTarget, true
					}
				}
			}
		}
	}

	return 0, false
}

func (m *mockRateCache) Convert(amount float64, from, to string) (float64, bool) {
	rate, ok := m.GetRate(from, to)
	if !ok {
		return 0, false
	}
	return amount * rate, true
}

func (m *mockRateCache) ConvertValue(v types.Value, target string) (types.Value, bool) {
	if v.IsError() || v.IsEmpty() {
		return v, false
	}

	// Resolve target to canonical code
	targetCode := resolveTestCode(target)

	switch v.Kind {
	case types.ValueCurrency:
		if v.Curr == nil {
			return v, false
		}
		converted, ok := m.Convert(v.Num, v.Curr.Code, targetCode)
		if !ok {
			return v, false
		}
		targetCurr := types.ParseCurrency(targetCode)
		if targetCurr == nil {
			targetCurr = types.CurrencyFromCode(targetCode)
		}
		return types.CurrencyValue(converted, targetCurr), true

	case types.ValueCrypto:
		if v.Crypto == nil {
			return v, false
		}
		converted, ok := m.Convert(v.Num, v.Crypto.Code, targetCode)
		if !ok {
			return v, false
		}
		// Target could be currency or crypto
		if targetCrypto := types.ParseCrypto(targetCode); targetCrypto != nil {
			return types.CryptoValue(converted, targetCrypto), true
		}
		if targetCurr := types.ParseCurrency(targetCode); targetCurr != nil {
			return types.CurrencyValue(converted, targetCurr), true
		}
		return types.Number(converted), true

	case types.ValueMetal:
		if v.Metal == nil {
			return v, false
		}
		converted, ok := m.Convert(v.Num, v.Metal.Code, targetCode)
		if !ok {
			return v, false
		}
		if targetCurr := types.ParseCurrency(targetCode); targetCurr != nil {
			return types.CurrencyValue(converted, targetCurr), true
		}
		return types.Number(converted), true

	default:
		return v, false
	}
}

// resolveTestCode resolves common aliases for testing.
func resolveTestCode(s string) string {
	aliases := map[string]string{
		"usd": "USD", "dollars": "USD", "dollar": "USD",
		"eur": "EUR", "euros": "EUR", "euro": "EUR",
		"try": "TRY", "lira": "TRY", "turkish lira": "TRY",
		"btc": "BTC", "bitcoin": "BTC",
		"xau": "XAU", "gold": "XAU",
	}

	if code, ok := aliases[s]; ok {
		return code
	}
	// Try lowercase
	if code, ok := aliases[strings.ToLower(s)]; ok {
		return code
	}
	return strings.ToUpper(s)
}
