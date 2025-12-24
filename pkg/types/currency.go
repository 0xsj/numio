// internal/types/currency.go

// Package types defines core value types for numio.
package types

import (
	"strings"
)

// Currency represents a fiat currency.
type Currency struct {
	Code        string   // ISO 4217 code: "USD", "EUR", "TRY"
	Symbol      string   // Symbol: "$", "€", "₺"
	Name        string   // Full name: "US Dollar", "Euro", "Turkish Lira"
	Aliases     []string // Natural language aliases: "dollars", "bucks"
	SymbolAfter bool     // true if symbol comes after amount (100₺ vs $100)
}

// String returns the currency code.
func (c Currency) String() string {
	return c.Code
}

// FormatAmount formats an amount with the currency symbol.
func (c Currency) FormatAmount(amount float64, precision int) string {
	// Format number with precision
	format := "%." + string(rune('0'+precision)) + "f"
	if precision == 0 {
		format = "%.0f"
	} else if precision == 2 {
		format = "%.2f"
	}

	var numStr string
	if amount < 0 {
		numStr = strings.Replace(sprintf(format, -amount), "-", "", 1)
		if c.SymbolAfter {
			return "-" + numStr + c.Symbol
		}
		return "-" + c.Symbol + numStr
	}

	numStr = sprintf(format, amount)
	if c.SymbolAfter {
		return numStr + c.Symbol
	}
	return c.Symbol + numStr
}

// sprintf is a simple float formatter to avoid importing fmt in hot path.
func sprintf(format string, v float64) string {
	// For now, use a simple implementation
	// In production, consider using strconv for performance
	switch format {
	case "%.0f":
		return formatFloat(v, 0)
	case "%.2f":
		return formatFloat(v, 2)
	default:
		return formatFloat(v, 2)
	}
}

// formatFloat formats a float with the given decimal places.
func formatFloat(v float64, decimals int) string {
	if decimals == 0 {
		return itoa(int64(v + 0.5))
	}

	// Multiply to shift decimals, round, then format
	shift := 1.0
	for i := 0; i < decimals; i++ {
		shift *= 10
	}

	rounded := int64(v*shift + 0.5)
	intPart := rounded / int64(shift)
	fracPart := rounded % int64(shift)

	// Ensure fracPart is positive
	if fracPart < 0 {
		fracPart = -fracPart
	}

	// Pad fractional part with leading zeros
	fracStr := itoa(fracPart)
	for len(fracStr) < decimals {
		fracStr = "0" + fracStr
	}

	return itoa(intPart) + "." + fracStr
}

// itoa converts an int64 to string without fmt package.
func itoa(n int64) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var buf [20]byte
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}

	if negative {
		i--
		buf[i] = '-'
	}

	return string(buf[i:])
}

// CurrencyRegistry holds all known currencies.
type CurrencyRegistry struct {
	byCode   map[string]*Currency
	bySymbol map[string]*Currency
	byAlias  map[string]*Currency
}

// Global currency registry.
var currencies = newCurrencyRegistry()

// newCurrencyRegistry creates and populates the currency registry.
func newCurrencyRegistry() *CurrencyRegistry {
	r := &CurrencyRegistry{
		byCode:   make(map[string]*Currency),
		bySymbol: make(map[string]*Currency),
		byAlias:  make(map[string]*Currency),
	}

	// Register all curated currencies
	for i := range curatedCurrencies {
		r.register(&curatedCurrencies[i])
	}

	return r
}

// register adds a currency to the registry.
func (r *CurrencyRegistry) register(c *Currency) {
	// By code (case-insensitive)
	r.byCode[strings.ToUpper(c.Code)] = c
	r.byCode[strings.ToLower(c.Code)] = c

	// By symbol
	if c.Symbol != "" {
		r.bySymbol[c.Symbol] = c
	}

	// By aliases (case-insensitive)
	for _, alias := range c.Aliases {
		r.byAlias[strings.ToLower(alias)] = c
	}
}

// Lookup finds a currency by code, symbol, or alias.
func (r *CurrencyRegistry) Lookup(s string) *Currency {
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

// curatedCurrencies contains well-known currencies with symbols and aliases.
var curatedCurrencies = []Currency{
	// ════════════════════════════════════════════════════════════
	// MAJOR CURRENCIES
	// ════════════════════════════════════════════════════════════
	{
		Code:    "USD",
		Symbol:  "$",
		Name:    "US Dollar",
		Aliases: []string{"dollar", "dollars", "usd", "bucks", "buck"},
	},
	{
		Code:    "EUR",
		Symbol:  "€",
		Name:    "Euro",
		Aliases: []string{"euro", "euros", "eur"},
	},
	{
		Code:    "GBP",
		Symbol:  "£",
		Name:    "British Pound",
		Aliases: []string{"pound", "pounds", "gbp", "quid", "sterling"},
	},
	{
		Code:    "JPY",
		Symbol:  "¥",
		Name:    "Japanese Yen",
		Aliases: []string{"yen", "jpy"},
	},
	{
		Code:    "CHF",
		Symbol:  "CHF",
		Name:    "Swiss Franc",
		Aliases: []string{"franc", "francs", "chf", "swiss franc", "swiss francs"},
	},

	// ════════════════════════════════════════════════════════════
	// AMERICAS
	// ════════════════════════════════════════════════════════════
	{
		Code:    "CAD",
		Symbol:  "C$",
		Name:    "Canadian Dollar",
		Aliases: []string{"cad", "canadian dollar", "canadian dollars", "loonie"},
	},
	{
		Code:    "MXN",
		Symbol:  "MX$",
		Name:    "Mexican Peso",
		Aliases: []string{"mxn", "peso", "pesos", "mexican peso"},
	},
	{
		Code:    "BRL",
		Symbol:  "R$",
		Name:    "Brazilian Real",
		Aliases: []string{"brl", "real", "reais", "brazilian real"},
	},
	{
		Code:    "ARS",
		Symbol:  "AR$",
		Name:    "Argentine Peso",
		Aliases: []string{"ars", "argentine peso"},
	},
	{
		Code:    "CLP",
		Symbol:  "CL$",
		Name:    "Chilean Peso",
		Aliases: []string{"clp", "chilean peso"},
	},
	{
		Code:    "COP",
		Symbol:  "CO$",
		Name:    "Colombian Peso",
		Aliases: []string{"cop", "colombian peso"},
	},

	// ════════════════════════════════════════════════════════════
	// EUROPE
	// ════════════════════════════════════════════════════════════
	{
		Code:        "RUB",
		Symbol:      "₽",
		Name:        "Russian Ruble",
		Aliases:     []string{"rub", "ruble", "rubles", "rouble", "roubles"},
		SymbolAfter: true,
	},
	{
		Code:        "UAH",
		Symbol:      "₴",
		Name:        "Ukrainian Hryvnia",
		Aliases:     []string{"uah", "hryvnia", "hryvnias"},
		SymbolAfter: true,
	},
	{
		Code:        "PLN",
		Symbol:      "zł",
		Name:        "Polish Zloty",
		Aliases:     []string{"pln", "zloty", "zlotys", "złoty"},
		SymbolAfter: true,
	},
	{
		Code:    "CZK",
		Symbol:  "Kč",
		Name:    "Czech Koruna",
		Aliases: []string{"czk", "koruna", "korunas", "czech koruna"},
	},
	{
		Code:    "SEK",
		Symbol:  "kr",
		Name:    "Swedish Krona",
		Aliases: []string{"sek", "swedish krona", "swedish kronor"},
	},
	{
		Code:    "NOK",
		Symbol:  "kr",
		Name:    "Norwegian Krone",
		Aliases: []string{"nok", "norwegian krone", "norwegian kroner"},
	},
	{
		Code:    "DKK",
		Symbol:  "kr",
		Name:    "Danish Krone",
		Aliases: []string{"dkk", "danish krone", "danish kroner"},
	},
	{
		Code:    "HUF",
		Symbol:  "Ft",
		Name:    "Hungarian Forint",
		Aliases: []string{"huf", "forint", "forints"},
	},
	{
		Code:    "RON",
		Symbol:  "lei",
		Name:    "Romanian Leu",
		Aliases: []string{"ron", "leu", "lei"},
	},

	// ════════════════════════════════════════════════════════════
	// MIDDLE EAST
	// ════════════════════════════════════════════════════════════
	{
		Code:        "TRY",
		Symbol:      "₺",
		Name:        "Turkish Lira",
		Aliases:     []string{"try", "tl", "lira", "liras", "turkish lira", "turk lirasi"},
		SymbolAfter: true,
	},
	{
		Code:    "ILS",
		Symbol:  "₪",
		Name:    "Israeli Shekel",
		Aliases: []string{"ils", "shekel", "shekels", "nis"},
	},
	{
		Code:    "AED",
		Symbol:  "د.إ",
		Name:    "UAE Dirham",
		Aliases: []string{"aed", "dirham", "dirhams", "emirati dirham"},
	},
	{
		Code:    "SAR",
		Symbol:  "﷼",
		Name:    "Saudi Riyal",
		Aliases: []string{"sar", "riyal", "riyals", "saudi riyal"},
	},
	{
		Code:    "QAR",
		Symbol:  "﷼",
		Name:    "Qatari Riyal",
		Aliases: []string{"qar", "qatari riyal"},
	},
	{
		Code:    "KWD",
		Symbol:  "د.ك",
		Name:    "Kuwaiti Dinar",
		Aliases: []string{"kwd", "kuwaiti dinar"},
	},
	{
		Code:    "EGP",
		Symbol:  "E£",
		Name:    "Egyptian Pound",
		Aliases: []string{"egp", "egyptian pound"},
	},

	// ════════════════════════════════════════════════════════════
	// ASIA PACIFIC
	// ════════════════════════════════════════════════════════════
	{
		Code:    "CNY",
		Symbol:  "¥",
		Name:    "Chinese Yuan",
		Aliases: []string{"cny", "yuan", "rmb", "renminbi", "chinese yuan"},
	},
	{
		Code:    "HKD",
		Symbol:  "HK$",
		Name:    "Hong Kong Dollar",
		Aliases: []string{"hkd", "hong kong dollar"},
	},
	{
		Code:    "TWD",
		Symbol:  "NT$",
		Name:    "Taiwan Dollar",
		Aliases: []string{"twd", "taiwan dollar", "nt dollar"},
	},
	{
		Code:    "KRW",
		Symbol:  "₩",
		Name:    "South Korean Won",
		Aliases: []string{"krw", "won", "korean won"},
	},
	{
		Code:    "INR",
		Symbol:  "₹",
		Name:    "Indian Rupee",
		Aliases: []string{"inr", "rupee", "rupees", "indian rupee"},
	},
	{
		Code:    "PKR",
		Symbol:  "₨",
		Name:    "Pakistani Rupee",
		Aliases: []string{"pkr", "pakistani rupee"},
	},
	{
		Code:    "BDT",
		Symbol:  "৳",
		Name:    "Bangladeshi Taka",
		Aliases: []string{"bdt", "taka", "bangladeshi taka"},
	},
	{
		Code:    "SGD",
		Symbol:  "S$",
		Name:    "Singapore Dollar",
		Aliases: []string{"sgd", "singapore dollar"},
	},
	{
		Code:    "MYR",
		Symbol:  "RM",
		Name:    "Malaysian Ringgit",
		Aliases: []string{"myr", "ringgit", "malaysian ringgit"},
	},
	{
		Code:    "THB",
		Symbol:  "฿",
		Name:    "Thai Baht",
		Aliases: []string{"thb", "baht", "thai baht"},
	},
	{
		Code:    "IDR",
		Symbol:  "Rp",
		Name:    "Indonesian Rupiah",
		Aliases: []string{"idr", "rupiah", "indonesian rupiah"},
	},
	{
		Code:    "VND",
		Symbol:  "₫",
		Name:    "Vietnamese Dong",
		Aliases: []string{"vnd", "dong", "vietnamese dong"},
	},
	{
		Code:    "PHP",
		Symbol:  "₱",
		Name:    "Philippine Peso",
		Aliases: []string{"php", "philippine peso"},
	},

	// ════════════════════════════════════════════════════════════
	// OCEANIA
	// ════════════════════════════════════════════════════════════
	{
		Code:    "AUD",
		Symbol:  "A$",
		Name:    "Australian Dollar",
		Aliases: []string{"aud", "australian dollar", "australian dollars", "aussie dollar"},
	},
	{
		Code:    "NZD",
		Symbol:  "NZ$",
		Name:    "New Zealand Dollar",
		Aliases: []string{"nzd", "new zealand dollar", "kiwi dollar"},
	},

	// ════════════════════════════════════════════════════════════
	// AFRICA
	// ════════════════════════════════════════════════════════════
	{
		Code:    "ZAR",
		Symbol:  "R",
		Name:    "South African Rand",
		Aliases: []string{"zar", "rand", "south african rand"},
	},
	{
		Code:    "NGN",
		Symbol:  "₦",
		Name:    "Nigerian Naira",
		Aliases: []string{"ngn", "naira", "nigerian naira"},
	},
	{
		Code:    "KES",
		Symbol:  "KSh",
		Name:    "Kenyan Shilling",
		Aliases: []string{"kes", "kenyan shilling"},
	},
}

// ════════════════════════════════════════════════════════════════
// PUBLIC API
// ════════════════════════════════════════════════════════════════

// LookupCurrency finds a currency by code, symbol, or alias.
// Returns nil if not found.
func LookupCurrency(s string) *Currency {
	return currencies.Lookup(s)
}

// ParseCurrency parses a string into a currency.
// Accepts codes ("USD"), symbols ("$"), or natural names ("dollars", "turkish lira").
// Returns nil if not found in curated list.
func ParseCurrency(s string) *Currency {
	return currencies.Lookup(strings.TrimSpace(s))
}

// CurrencyFromCode creates a dynamic currency from an ISO 4217 code.
// Use this for API-discovered currencies not in the curated list.
func CurrencyFromCode(code string) *Currency {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) != 3 {
		return nil
	}

	// Check if it's already curated
	if c := currencies.Lookup(code); c != nil {
		return c
	}

	// Create dynamic currency (no symbol, just code)
	return &Currency{
		Code:   code,
		Symbol: code,
		Name:   code,
	}
}

// LookupCurrencyBySymbol finds a currency by its symbol.
func LookupCurrencyBySymbol(symbol string) *Currency {
	return currencies.bySymbol[symbol]
}

// IsKnownCurrencyCode checks if a code is a known curated currency.
func IsKnownCurrencyCode(code string) bool {
	return currencies.byCode[strings.ToUpper(code)] != nil
}

// AllCurrencies returns all curated currencies.
func AllCurrencies() []Currency {
	return curatedCurrencies
}

// CurrencyCodes returns all curated currency codes.
func CurrencyCodes() []string {
	codes := make([]string, len(curatedCurrencies))
	for i, c := range curatedCurrencies {
		codes[i] = c.Code
	}
	return codes
}

// CurrencySymbols returns a map of symbol → currency code.
func CurrencySymbols() map[string]string {
	symbols := make(map[string]string)
	for _, c := range curatedCurrencies {
		if c.Symbol != "" && c.Symbol != c.Code {
			symbols[c.Symbol] = c.Code
		}
	}
	return symbols
}

// IsCurrencySymbol checks if a string is a known currency symbol.
func IsCurrencySymbol(s string) bool {
	return currencies.bySymbol[s] != nil
}

// IsCurrencySymbolRune checks if a rune is a known currency symbol.
func IsCurrencySymbolRune(r rune) bool {
	return currencies.bySymbol[string(r)] != nil
}
