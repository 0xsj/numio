// internal/token/token.go

// Package token defines token types for the numio lexer.
package token

import "slices"

// Type represents the type of a token.
type Type int

const (
	// Special tokens
	EOF     Type = iota // End of input
	ILLEGAL             // Unknown/invalid token

	// Literals
	NUMBER     // 42, 3.14, 1,234.56, 1.5e6
	PERCENT    // 20%
	IDENTIFIER // variable names, unit names, currency codes

	// Operators
	PLUS     // +
	MINUS    // -
	STAR     // *
	SLASH    // /
	CARET    // ^
	POWER    // **
	LPAREN   // (
	RPAREN   // )
	EQUALS   // =
	COMMA    // ,

	// Keywords
	IN // in, to (for conversions)
	OF // of (for "20% of 150")

	// Currency symbols
	DOLLAR   // $
	EURO     // €
	POUND    // £
	YEN      // ¥
	BITCOIN  // ₿
	CURRENCY // Generic currency symbol (₹, ₩, ₽, ₪, etc.)

	// Comments
	COMMENT // # or // to end of line

	// Whitespace (usually skipped, but tracked for position)
	WHITESPACE
	NEWLINE
)

var typeNames = map[Type]string{
	EOF:        "EOF",
	ILLEGAL:    "ILLEGAL",
	NUMBER:     "NUMBER",
	PERCENT:    "PERCENT",
	IDENTIFIER: "IDENTIFIER",
	PLUS:       "PLUS",
	MINUS:      "MINUS",
	STAR:       "STAR",
	SLASH:      "SLASH",
	CARET:      "CARET",
	POWER:      "POWER",
	LPAREN:     "LPAREN",
	RPAREN:     "RPAREN",
	EQUALS:     "EQUALS",
	COMMA:      "COMMA",
	IN:         "IN",
	OF:         "OF",
	DOLLAR:     "DOLLAR",
	EURO:       "EURO",
	POUND:      "POUND",
	YEN:        "YEN",
	BITCOIN:    "BITCOIN",
	CURRENCY:   "CURRENCY",
	COMMENT:    "COMMENT",
	WHITESPACE: "WHITESPACE",
	NEWLINE:    "NEWLINE",
}

func (t Type) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

// Token represents a lexical token.
type Token struct {
	Type    Type   // Token type
	Literal string // Raw text of the token
	Pos     int    // Start position in input (byte offset)
}

// New creates a new token.
func New(typ Type, literal string, pos int) Token {
	return Token{
		Type:    typ,
		Literal: literal,
		Pos:     pos,
	}
}

// Is checks if the token is of a specific type.
func (t Token) Is(typ Type) bool {
	return t.Type == typ
}

// IsOneOf checks if the token is one of the given types.
func (t Token) IsOneOf(types ...Type) bool {
	return slices.Contains(types, t.Type)
}

// IsOperator checks if the token is a binary operator.
func (t Token) IsOperator() bool {
	return t.IsOneOf(PLUS, MINUS, STAR, SLASH, CARET, POWER)
}

// IsCurrencySymbol checks if the token is a currency symbol.
func (t Token) IsCurrencySymbol() bool {
	return t.IsOneOf(DOLLAR, EURO, POUND, YEN, BITCOIN, CURRENCY)
}

// IsKeyword checks if the token is a keyword.
func (t Token) IsKeyword() bool {
	return t.IsOneOf(IN, OF)
}

// Keywords maps keyword strings to token types.
var Keywords = map[string]Type{
	"in": IN,
	"to": IN, // "to" is an alias for "in"
	"of": OF,
}

// LookupIdentifier checks if an identifier is a keyword.
// Returns the keyword type if found, otherwise IDENTIFIER.
func LookupIdentifier(ident string) Type {
	if typ, ok := Keywords[ident]; ok {
		return typ
	}
	return IDENTIFIER
}

// CurrencySymbols maps currency runes to token types.
var CurrencySymbols = map[rune]Type{
	'$': DOLLAR,
	'€': EURO,
	'£': POUND,
	'¥': YEN,
	'₿': BITCOIN,
	'₹': CURRENCY, // Indian Rupee
	'₩': CURRENCY, // Korean Won
	'₽': CURRENCY, // Russian Ruble
	'₪': CURRENCY, // Israeli Shekel
	'₴': CURRENCY, // Ukrainian Hryvnia
	'₮': CURRENCY, // Tether (USDT symbol)
	'₳': CURRENCY, // Cardano (ADA symbol)
	'Ð': CURRENCY, // Dogecoin
	'Ł': CURRENCY, // Litecoin
	'Ξ': CURRENCY, // Ethereum
	'◎': CURRENCY, // Solana
}

// IsCurrencyRune checks if a rune is a currency symbol.
func IsCurrencyRune(r rune) bool {
	_, ok := CurrencySymbols[r]
	return ok
}

// LookupCurrencySymbol returns the token type for a currency rune.
func LookupCurrencySymbol(r rune) Type {
	if typ, ok := CurrencySymbols[r]; ok {
		return typ
	}
	return ILLEGAL
}