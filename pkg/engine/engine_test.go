// pkg/engine/engine_test.go

package engine

import (
	"testing"
)

// TestConversionWithAliases tests that currency aliases work end-to-end.
func TestConversionWithAliases(t *testing.T) {
	eng := New()

	// Set up some test rates
	eng.SetRate("USD", "TRY", 32.5)
	eng.SetRate("USD", "EUR", 0.92)
	eng.SetRate("BTC", "USD", 95000.0)
	eng.SetRate("XAU", "USD", 2650.0)

	tests := []struct {
		input       string
		wantApprox  float64
		wantErr     bool
		description string
	}{
		// Symbol-prefixed (baseline - should work)
		{"$100 to EUR", 92.0, false, "USD symbol to EUR code"},
		{"$100 in EUR", 92.0, false, "USD symbol in EUR code"},

		// The bug: aliases should work
		{"$100 to lira", 3250.0, false, "USD symbol to TRY alias 'lira'"},
		{"$100 to turkish lira", 3250.0, false, "USD symbol to TRY alias 'turkish lira'"},
		{"$100 in euros", 92.0, false, "USD symbol to EUR alias 'euros'"},

		// Code-suffixed with alias target
		{"100 usd to lira", 3250.0, false, "USD code to TRY alias"},
		{"100 USD to turkish lira", 3250.0, false, "USD code to multi-word alias"},
		{"100 usd to euros", 92.0, false, "USD code to EUR alias"},

		// Alias-suffixed with alias target
		{"100 dollars to lira", 3250.0, false, "USD alias to TRY alias"},
		{"100 dollars to euros", 92.0, false, "USD alias to EUR alias"},

		// Crypto with aliases
		{"1 btc to usd", 95000.0, false, "BTC code to USD code"},
		{"1 bitcoin to dollars", 95000.0, false, "BTC alias to USD alias"},

		// Metal with aliases
		{"1 xau to usd", 2650.0, false, "XAU code to USD code"},
		{"1 gold to dollars", 2650.0, false, "XAU alias to USD alias"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := eng.Eval(tt.input)

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

			// Check approximate value (within 1% tolerance)
			got := result.AsFloat()
			tolerance := tt.wantApprox * 0.01
			if tolerance < 0.01 {
				tolerance = 0.01
			}
			if got < tt.wantApprox-tolerance || got > tt.wantApprox+tolerance {
				t.Errorf("value: got %v, want ~%v", got, tt.wantApprox)
			}
		})
	}
}

// TestSetRateWithAliases tests that SetRate works with aliases.
func TestSetRateWithAliases(t *testing.T) {
	eng := New()

	// Set rate using aliases
	eng.SetRate("dollars", "lira", 32.5)

	// Should be retrievable via codes
	rate, ok := eng.GetRate("USD", "TRY")
	if !ok {
		t.Error("expected to find USD->TRY rate")
		return
	}
	if rate != 32.5 {
		t.Errorf("rate: got %v, want 32.5", rate)
	}

	// Should also work via aliases
	rate, ok = eng.GetRate("dollars", "lira")
	if !ok {
		t.Error("expected to find dollars->lira rate")
		return
	}
	if rate != 32.5 {
		t.Errorf("rate: got %v, want 32.5", rate)
	}
}

// TestDebugConversion helps debug specific conversion issues.
func TestDebugConversion(t *testing.T) {
	eng := New()

	// Set a known rate
	eng.SetRate("USD", "TRY", 32.5)

	input := "$100 to lira"
	t.Logf("Input: %q", input)

	result := eng.Eval(input)
	t.Logf("Result: %v (kind: %v, error: %v)", result, result.Kind, result.IsError())

	if result.IsError() {
		t.Logf("This is an error result")
	}
}
