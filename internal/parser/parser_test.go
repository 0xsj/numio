// internal/parser/parser_test.go

package parser

import (
	"testing"

	"github.com/0xsj/numio/internal/ast"
	"github.com/0xsj/numio/internal/lexer"
	"github.com/0xsj/numio/internal/token"
)

// TestTokenization verifies the lexer produces expected tokens.
func TestTokenization(t *testing.T) {
	tests := []struct {
		input    string
		expected []token.Type
	}{
		{"100 usd", []token.Type{token.NUMBER, token.IDENTIFIER, token.EOF}},
		{"100 usd to lira", []token.Type{token.NUMBER, token.IDENTIFIER, token.IN, token.IDENTIFIER, token.EOF}},
		{"$100 to lira", []token.Type{token.DOLLAR, token.NUMBER, token.IN, token.IDENTIFIER, token.EOF}},
		{"100 dollars to euros", []token.Type{token.NUMBER, token.IDENTIFIER, token.IN, token.IDENTIFIER, token.EOF}},
		{"100 turkish lira", []token.Type{token.NUMBER, token.IDENTIFIER, token.EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := lexer.Tokenize(tt.input)

			if len(tokens) != len(tt.expected) {
				t.Errorf("token count mismatch: got %d, want %d", len(tokens), len(tt.expected))
				for i, tok := range tokens {
					t.Logf("  [%d] %s: %q", i, tok.Type, tok.Literal)
				}
				return
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i] {
					t.Errorf("token[%d] type: got %s, want %s (literal: %q)",
						i, tok.Type, tt.expected[i], tok.Literal)
				}
			}
		})
	}
}

// TestParseConversion tests conversion expression parsing.
func TestParseConversion(t *testing.T) {
	tests := []struct {
		input       string
		wantConv    bool   // expect ConversionExpr
		wantTarget  string // expected target
		wantErr     bool
		description string
	}{
		// Symbol-prefixed (should already work)
		{"$100 to EUR", true, "EUR", false, "symbol currency to code"},
		{"$100 to lira", true, "lira", false, "symbol currency to alias"},
		{"$100 in euros", true, "euros", false, "symbol currency in alias"},

		// Code-suffixed (the bug we're fixing)
		{"100 usd to EUR", true, "EUR", false, "code currency to code"},
		{"100 usd to lira", true, "lira", false, "code currency to alias"},
		{"100 USD to turkish lira", true, "turkish lira", false, "code to multi-word alias"},

		// Alias-suffixed
		{"100 dollars to euros", true, "euros", false, "alias to alias"},
		{"100 dollars to lira", true, "lira", false, "alias to alias 2"},

		// Crypto
		{"1 btc to usd", true, "usd", false, "crypto to currency"},
		{"1 bitcoin to dollars", true, "dollars", false, "crypto alias to currency alias"},

		// Metal
		{"1 xau to usd", true, "usd", false, "metal to currency"},
		{"1 gold to dollars", true, "dollars", false, "metal alias to currency alias"},

		// Units
		{"100 km to miles", true, "miles", false, "unit conversion"},
		{"5 kg to pounds", true, "pounds", false, "weight conversion"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			expr, errs := ParseExpr(tt.input)

			if tt.wantErr {
				if len(errs) == 0 {
					t.Errorf("expected error, got none")
				}
				return
			}

			if len(errs) > 0 {
				t.Errorf("unexpected errors: %v", errs)
				return
			}

			if expr == nil {
				t.Errorf("got nil expression")
				return
			}

			conv, ok := expr.(*ast.ConversionExpr)
			if tt.wantConv && !ok {
				t.Errorf("expected ConversionExpr, got %T: %v", expr, expr)
				return
			}

			if tt.wantConv && conv.Target != tt.wantTarget {
				t.Errorf("target: got %q, want %q", conv.Target, tt.wantTarget)
			}
		})
	}
}

// TestParseCurrencyLiteral tests that currency suffixes are recognized.
func TestParseCurrencyLiteral(t *testing.T) {
	tests := []struct {
		input    string
		wantType string // "currency", "crypto", "metal", "unit", "number"
		wantCode string // expected currency/unit code
	}{
		// Currency codes
		{"100 usd", "currency", "USD"},
		{"100 USD", "currency", "USD"},
		{"100 eur", "currency", "EUR"},
		{"100 try", "currency", "TRY"},

		// Currency aliases
		{"100 dollars", "currency", "USD"},
		{"100 euros", "currency", "EUR"},
		{"100 lira", "currency", "TRY"},
		{"100 turkish lira", "currency", "TRY"},

		// Crypto
		{"1 btc", "crypto", "BTC"},
		{"1 bitcoin", "crypto", "BTC"},
		{"0.5 eth", "crypto", "ETH"},

		// Metal
		{"1 xau", "metal", "XAU"},
		{"1 gold", "metal", "XAU"},

		// Units
		{"100 km", "unit", "km"},
		{"5 miles", "unit", "mi"},

		// Plain number (no suffix)
		{"100", "number", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, errs := ParseExpr(tt.input)

			if len(errs) > 0 {
				t.Errorf("unexpected errors: %v", errs)
				return
			}

			switch tt.wantType {
			case "currency":
				curr, ok := expr.(*ast.CurrencyLit)
				if !ok {
					t.Errorf("expected CurrencyLit, got %T", expr)
					return
				}
				if curr.Currency == nil || curr.Currency.Code != tt.wantCode {
					code := ""
					if curr.Currency != nil {
						code = curr.Currency.Code
					}
					t.Errorf("currency code: got %q, want %q", code, tt.wantCode)
				}

			case "crypto":
				crypto, ok := expr.(*ast.CryptoLit)
				if !ok {
					t.Errorf("expected CryptoLit, got %T", expr)
					return
				}
				if crypto.Crypto == nil || crypto.Crypto.Code != tt.wantCode {
					code := ""
					if crypto.Crypto != nil {
						code = crypto.Crypto.Code
					}
					t.Errorf("crypto code: got %q, want %q", code, tt.wantCode)
				}

			case "metal":
				metal, ok := expr.(*ast.MetalLit)
				if !ok {
					t.Errorf("expected MetalLit, got %T", expr)
					return
				}
				if metal.Metal == nil || metal.Metal.Code != tt.wantCode {
					code := ""
					if metal.Metal != nil {
						code = metal.Metal.Code
					}
					t.Errorf("metal code: got %q, want %q", code, tt.wantCode)
				}

			case "unit":
				unit, ok := expr.(*ast.UnitLit)
				if !ok {
					t.Errorf("expected UnitLit, got %T", expr)
					return
				}
				if unit.Unit == nil || unit.Unit.Code != tt.wantCode {
					code := ""
					if unit.Unit != nil {
						code = unit.Unit.Code
					}
					t.Errorf("unit code: got %q, want %q", code, tt.wantCode)
				}

			case "number":
				_, ok := expr.(*ast.NumberLit)
				if !ok {
					t.Errorf("expected NumberLit, got %T", expr)
				}
			}
		})
	}
}

// TestDebugTokensAndParse is a helper to debug specific inputs.
func TestDebugTokensAndParse(t *testing.T) {
	input := "100 usd to lira"

	t.Logf("Input: %q", input)

	// Show tokens
	tokens := lexer.Tokenize(input)
	t.Log("Tokens:")
	for i, tok := range tokens {
		t.Logf("  [%d] %s: %q (pos=%d)", i, tok.Type, tok.Literal, tok.Pos)
	}

	// Parse
	expr, errs := ParseExpr(input)

	if len(errs) > 0 {
		t.Log("Parse errors:")
		for _, err := range errs {
			t.Logf("  %v", err)
		}
	}

	if expr != nil {
		t.Logf("Parsed: %T = %v", expr, expr)
	} else {
		t.Log("Parsed: nil")
	}
}
