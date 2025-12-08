// internal/types/crypto.go

package types

import (
	"strings"
)

// Crypto represents a cryptocurrency.
type Crypto struct {
	Code        string   // Ticker symbol: "BTC", "ETH"
	Symbol      string   // Display symbol: "₿", "Ξ"
	Name        string   // Full name: "Bitcoin", "Ethereum"
	Aliases     []string // Natural language aliases
	CoingeckoID string   // CoinGecko API identifier
	Decimals    int      // Typical decimal places for display
}

// String returns the crypto code.
func (c Crypto) String() string {
	return c.Code
}

// HasSymbol returns true if the crypto has a unicode symbol.
func (c Crypto) HasSymbol() bool {
	return c.Symbol != "" && c.Symbol != c.Code
}

// CryptoRegistry holds all known cryptocurrencies.
type CryptoRegistry struct {
	byCode   map[string]*Crypto
	bySymbol map[string]*Crypto
	byAlias  map[string]*Crypto
}

// Global crypto registry.
var cryptos = newCryptoRegistry()

// newCryptoRegistry creates and populates the crypto registry.
func newCryptoRegistry() *CryptoRegistry {
	r := &CryptoRegistry{
		byCode:   make(map[string]*Crypto),
		bySymbol: make(map[string]*Crypto),
		byAlias:  make(map[string]*Crypto),
	}

	for i := range curatedCryptos {
		r.register(&curatedCryptos[i])
	}

	return r
}

// register adds a crypto to the registry.
func (r *CryptoRegistry) register(c *Crypto) {
	// By code (case-insensitive)
	r.byCode[strings.ToUpper(c.Code)] = c
	r.byCode[strings.ToLower(c.Code)] = c

	// By symbol
	if c.Symbol != "" && c.Symbol != c.Code {
		r.bySymbol[c.Symbol] = c
	}

	// By aliases (case-insensitive)
	for _, alias := range c.Aliases {
		r.byAlias[strings.ToLower(alias)] = c
	}
}

// Lookup finds a crypto by code, symbol, or alias.
func (r *CryptoRegistry) Lookup(s string) *Crypto {
	// Try exact symbol match first
	if c, ok := r.bySymbol[s]; ok {
		return c
	}

	// Try code (case-insensitive)
	if c, ok := r.byCode[strings.ToUpper(s)]; ok {
		return c
	}

	// Try alias (case-insensitive)
	if c, ok := r.byAlias[strings.ToLower(s)]; ok {
		return c
	}

	return nil
}

// curatedCryptos contains well-known cryptocurrencies.
var curatedCryptos = []Crypto{
	// ════════════════════════════════════════════════════════════
	// TIER 1 - Major cryptocurrencies
	// ════════════════════════════════════════════════════════════
	{
		Code:        "BTC",
		Symbol:      "₿",
		Name:        "Bitcoin",
		Aliases:     []string{"bitcoin", "btc", "xbt", "satoshi", "sats"},
		CoingeckoID: "bitcoin",
		Decimals:    8,
	},
	{
		Code:        "ETH",
		Symbol:      "Ξ",
		Name:        "Ethereum",
		Aliases:     []string{"ethereum", "eth", "ether"},
		CoingeckoID: "ethereum",
		Decimals:    6,
	},

	// ════════════════════════════════════════════════════════════
	// TIER 2 - Stablecoins
	// ════════════════════════════════════════════════════════════
	{
		Code:        "USDT",
		Symbol:      "₮",
		Name:        "Tether",
		Aliases:     []string{"tether", "usdt"},
		CoingeckoID: "tether",
		Decimals:    2,
	},
	{
		Code:        "USDC",
		Symbol:      "USDC",
		Name:        "USD Coin",
		Aliases:     []string{"usd coin", "usdc"},
		CoingeckoID: "usd-coin",
		Decimals:    2,
	},
	{
		Code:        "DAI",
		Symbol:      "DAI",
		Name:        "Dai",
		Aliases:     []string{"dai", "makerdao"},
		CoingeckoID: "dai",
		Decimals:    2,
	},
	{
		Code:        "BUSD",
		Symbol:      "BUSD",
		Name:        "Binance USD",
		Aliases:     []string{"binance usd", "busd"},
		CoingeckoID: "binance-usd",
		Decimals:    2,
	},

	// ════════════════════════════════════════════════════════════
	// TIER 3 - Major altcoins
	// ════════════════════════════════════════════════════════════
	{
		Code:        "BNB",
		Symbol:      "BNB",
		Name:        "BNB",
		Aliases:     []string{"bnb", "binance coin", "binance"},
		CoingeckoID: "binancecoin",
		Decimals:    4,
	},
	{
		Code:        "SOL",
		Symbol:      "◎",
		Name:        "Solana",
		Aliases:     []string{"solana", "sol"},
		CoingeckoID: "solana",
		Decimals:    4,
	},
	{
		Code:        "XRP",
		Symbol:      "XRP",
		Name:        "XRP",
		Aliases:     []string{"xrp", "ripple"},
		CoingeckoID: "ripple",
		Decimals:    4,
	},
	{
		Code:        "ADA",
		Symbol:      "₳",
		Name:        "Cardano",
		Aliases:     []string{"cardano", "ada"},
		CoingeckoID: "cardano",
		Decimals:    4,
	},
	{
		Code:        "DOGE",
		Symbol:      "Ð",
		Name:        "Dogecoin",
		Aliases:     []string{"dogecoin", "doge"},
		CoingeckoID: "dogecoin",
		Decimals:    4,
	},
	{
		Code:        "DOT",
		Symbol:      "DOT",
		Name:        "Polkadot",
		Aliases:     []string{"polkadot", "dot"},
		CoingeckoID: "polkadot",
		Decimals:    4,
	},
	{
		Code:        "MATIC",
		Symbol:      "MATIC",
		Name:        "Polygon",
		Aliases:     []string{"polygon", "matic"},
		CoingeckoID: "matic-network",
		Decimals:    4,
	},
	{
		Code:        "AVAX",
		Symbol:      "AVAX",
		Name:        "Avalanche",
		Aliases:     []string{"avalanche", "avax"},
		CoingeckoID: "avalanche-2",
		Decimals:    4,
	},
	{
		Code:        "LTC",
		Symbol:      "Ł",
		Name:        "Litecoin",
		Aliases:     []string{"litecoin", "ltc"},
		CoingeckoID: "litecoin",
		Decimals:    4,
	},
	{
		Code:        "LINK",
		Symbol:      "LINK",
		Name:        "Chainlink",
		Aliases:     []string{"chainlink", "link"},
		CoingeckoID: "chainlink",
		Decimals:    4,
	},
	{
		Code:        "ATOM",
		Symbol:      "ATOM",
		Name:        "Cosmos",
		Aliases:     []string{"cosmos", "atom"},
		CoingeckoID: "cosmos",
		Decimals:    4,
	},
	{
		Code:        "UNI",
		Symbol:      "UNI",
		Name:        "Uniswap",
		Aliases:     []string{"uniswap", "uni"},
		CoingeckoID: "uniswap",
		Decimals:    4,
	},
	{
		Code:        "XLM",
		Symbol:      "XLM",
		Name:        "Stellar",
		Aliases:     []string{"stellar", "xlm", "lumens"},
		CoingeckoID: "stellar",
		Decimals:    4,
	},
	{
		Code:        "ALGO",
		Symbol:      "ALGO",
		Name:        "Algorand",
		Aliases:     []string{"algorand", "algo"},
		CoingeckoID: "algorand",
		Decimals:    4,
	},
	{
		Code:        "TON",
		Symbol:      "TON",
		Name:        "Toncoin",
		Aliases:     []string{"toncoin", "ton", "telegram"},
		CoingeckoID: "the-open-network",
		Decimals:    4,
	},

	// ════════════════════════════════════════════════════════════
	// TIER 4 - DeFi & others
	// ════════════════════════════════════════════════════════════
	{
		Code:        "AAVE",
		Symbol:      "AAVE",
		Name:        "Aave",
		Aliases:     []string{"aave"},
		CoingeckoID: "aave",
		Decimals:    4,
	},
	{
		Code:        "MKR",
		Symbol:      "MKR",
		Name:        "Maker",
		Aliases:     []string{"maker", "mkr"},
		CoingeckoID: "maker",
		Decimals:    4,
	},
	{
		Code:        "CRV",
		Symbol:      "CRV",
		Name:        "Curve",
		Aliases:     []string{"curve", "crv"},
		CoingeckoID: "curve-dao-token",
		Decimals:    4,
	},
	{
		Code:        "NEAR",
		Symbol:      "NEAR",
		Name:        "NEAR Protocol",
		Aliases:     []string{"near", "near protocol"},
		CoingeckoID: "near",
		Decimals:    4,
	},
	{
		Code:        "APT",
		Symbol:      "APT",
		Name:        "Aptos",
		Aliases:     []string{"aptos", "apt"},
		CoingeckoID: "aptos",
		Decimals:    4,
	},
	{
		Code:        "ARB",
		Symbol:      "ARB",
		Name:        "Arbitrum",
		Aliases:     []string{"arbitrum", "arb"},
		CoingeckoID: "arbitrum",
		Decimals:    4,
	},
	{
		Code:        "OP",
		Symbol:      "OP",
		Name:        "Optimism",
		Aliases:     []string{"optimism", "op"},
		CoingeckoID: "optimism",
		Decimals:    4,
	},

	// ════════════════════════════════════════════════════════════
	// TIER 5 - Memecoins & trending
	// ════════════════════════════════════════════════════════════
	{
		Code:        "SHIB",
		Symbol:      "SHIB",
		Name:        "Shiba Inu",
		Aliases:     []string{"shiba", "shib", "shiba inu"},
		CoingeckoID: "shiba-inu",
		Decimals:    8,
	},
	{
		Code:        "PEPE",
		Symbol:      "PEPE",
		Name:        "Pepe",
		Aliases:     []string{"pepe"},
		CoingeckoID: "pepe",
		Decimals:    8,
	},
	{
		Code:        "WIF",
		Symbol:      "WIF",
		Name:        "dogwifhat",
		Aliases:     []string{"dogwifhat", "wif"},
		CoingeckoID: "dogwifcoin",
		Decimals:    4,
	},
	{
		Code:        "BONK",
		Symbol:      "BONK",
		Name:        "Bonk",
		Aliases:     []string{"bonk"},
		CoingeckoID: "bonk",
		Decimals:    8,
	},
}

// ════════════════════════════════════════════════════════════════
// PUBLIC API
// ════════════════════════════════════════════════════════════════

// LookupCrypto finds a crypto by code, symbol, or alias.
// Returns nil if not found.
func LookupCrypto(s string) *Crypto {
	return cryptos.Lookup(s)
}

// ParseCrypto parses a string into a crypto.
// Accepts codes ("BTC"), symbols ("₿"), or natural names ("bitcoin").
// Returns nil if not found.
func ParseCrypto(s string) *Crypto {
	return cryptos.Lookup(strings.TrimSpace(s))
}

// IsCrypto checks if a string refers to a known cryptocurrency.
func IsCrypto(s string) bool {
	return cryptos.Lookup(s) != nil
}

// IsCryptoCode checks if a string is a crypto ticker code.
func IsCryptoCode(code string) bool {
	return cryptos.byCode[strings.ToUpper(code)] != nil
}

// IsCryptoSymbol checks if a string is a crypto symbol (₿, Ξ, etc).
func IsCryptoSymbol(s string) bool {
	return cryptos.bySymbol[s] != nil
}

// IsCryptoSymbolRune checks if a rune is a crypto symbol.
func IsCryptoSymbolRune(r rune) bool {
	return cryptos.bySymbol[string(r)] != nil
}

// AllCryptos returns all curated cryptocurrencies.
func AllCryptos() []Crypto {
	return curatedCryptos
}

// CryptoCodes returns all crypto ticker codes.
func CryptoCodes() []string {
	codes := make([]string, len(curatedCryptos))
	for i, c := range curatedCryptos {
		codes[i] = c.Code
	}
	return codes
}

// CryptoSymbols returns a map of symbol → crypto code.
func CryptoSymbols() map[string]string {
	symbols := make(map[string]string)
	for _, c := range curatedCryptos {
		if c.Symbol != "" && c.Symbol != c.Code {
			symbols[c.Symbol] = c.Code
		}
	}
	return symbols
}

// Stablecoins returns only stablecoin cryptos.
func Stablecoins() []Crypto {
	stables := make([]Crypto, 0, 4)
	for _, c := range curatedCryptos {
		switch c.Code {
		case "USDT", "USDC", "DAI", "BUSD":
			stables = append(stables, c)
		}
	}
	return stables
}

// CoingeckoIDs returns a map of code → CoinGecko API ID.
func CoingeckoIDs() map[string]string {
	ids := make(map[string]string)
	for _, c := range curatedCryptos {
		if c.CoingeckoID != "" {
			ids[c.Code] = c.CoingeckoID
		}
	}
	return ids
}