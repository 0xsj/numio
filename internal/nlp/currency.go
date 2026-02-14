// internal/nlp/currency.go

package nlp

import (
	"strings"
)

// ════════════════════════════════════════════════════════════════
// CURRENCY PATTERNS
// ════════════════════════════════════════════════════════════════

// RegisterCurrencyPatterns registers all currency patterns.
func RegisterCurrencyPatterns(r *PatternRegistry) {
	r.RegisterAll([]*Pattern{
		// Conversion patterns
		patternAmountToCode(),
		patternAmountInCode(),
		patternConvertAmountToCode(),
		patternHowMuchIn(),

		// Symbol-based
		patternSymbolAmountToCode(),
		patternSymbolAmountInCode(),

		// Word-based
		patternAmountCurrencyNameToCode(),
		patternAmountCurrencyNameInCode(),

		// Crypto patterns
		patternCryptoToFiat(),
		patternFiatToCrypto(),

		// Metal patterns
		patternMetalToFiat(),
		patternTurkishGold(),
		patternMetalPriceIn(),
	})
}

// ════════════════════════════════════════════════════════════════
// BASIC CONVERSION PATTERNS
// ════════════════════════════════════════════════════════════════

// patternAmountToCode: "100 USD to EUR", "50 usd to gbp"
func patternAmountToCode() *Pattern {
	return NewPattern("amount_to_code").
		Regex(`(\d+(?:[.,]\d+)?)\s*([a-zA-Z]{3})\s+to\s+([a-zA-Z]{3})`).
		Keywords("to").
		Priority(100).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			amount, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			from := strings.ToUpper(matches[2])
			to := strings.ToUpper(matches[3])

			if !isValidCurrencyCode(from) || !isValidCurrencyCode(to) {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// patternAmountInCode: "100 USD in EUR", "50 EUR in JPY"
func patternAmountInCode() *Pattern {
	return NewPattern("amount_in_code").
		Regex(`(\d+(?:[.,]\d+)?)\s*([a-zA-Z]{3})\s+in\s+([a-zA-Z]{3})`).
		Keywords("in").
		Priority(95).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			amount, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			from := strings.ToUpper(matches[2])
			to := strings.ToUpper(matches[3])

			if !isValidCurrencyCode(from) || !isValidCurrencyCode(to) {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// patternConvertAmountToCode: "convert 100 USD to EUR", "convert $50 to GBP"
func patternConvertAmountToCode() *Pattern {
	return NewPattern("convert_amount_to_code").
		Regex(`convert\s+([$€£¥]?)(\d+(?:[.,]\d+)?)\s*([a-zA-Z]{3})?\s+to\s+([a-zA-Z]{3})`).
		Keywords("convert", "to").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			symbol := matches[1]
			amount, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			from := strings.ToUpper(matches[3])
			if from == "" {
				from = symbolToCurrencyCode(symbol)
			}
			to := strings.ToUpper(matches[4])

			if from == "" || !isValidCurrencyCode(to) {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// patternHowMuchIn: "how much is 100 USD in EUR", "how much is $50 in pounds"
func patternHowMuchIn() *Pattern {
	return NewPattern("how_much_in").
		Regex(`how\s+much\s+(?:is\s+)?([$€£¥]?)(\d+(?:[.,]\d+)?)\s*([a-zA-Z]+)?\s+in\s+([a-zA-Z]+)`).
		Keywords("how", "much", "in").
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			symbol := matches[1]
			amount, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			from := matches[3]
			if from == "" {
				from = symbolToCurrencyCode(symbol)
			} else {
				from = normalizeCurrencyName(from)
			}

			to := normalizeCurrencyName(matches[4])

			if from == "" || to == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// SYMBOL-BASED PATTERNS
// ════════════════════════════════════════════════════════════════

// patternSymbolAmountToCode: "$100 to EUR", "€50 to USD"
func patternSymbolAmountToCode() *Pattern {
	return NewPattern("symbol_amount_to_code").
		Regex(`([$€£¥₿])(\d+(?:[.,]\d+)?[kmb]?)\s+to\s+([a-zA-Z]+)`).
		Keywords("to").
		Priority(95).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			symbol := matches[1]
			amount, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			from := symbolToCurrencyCode(symbol)
			to := normalizeCurrencyName(matches[3])

			if from == "" || to == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// patternSymbolAmountInCode: "$100 in EUR", "€50 in USD"
func patternSymbolAmountInCode() *Pattern {
	return NewPattern("symbol_amount_in_code").
		Regex(`([$€£¥₿])(\d+(?:[.,]\d+)?[kmb]?)\s+in\s+([a-zA-Z]+)`).
		Keywords("in").
		Priority(90).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			symbol := matches[1]
			amount, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			from := symbolToCurrencyCode(symbol)
			to := normalizeCurrencyName(matches[3])

			if from == "" || to == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// WORD-BASED PATTERNS
// ════════════════════════════════════════════════════════════════

// patternAmountCurrencyNameToCode: "100 dollars to euros", "50 pounds to yen"
func patternAmountCurrencyNameToCode() *Pattern {
	return NewPattern("amount_currency_name_to_code").
		Regex(`(\d+(?:[.,]\d+)?[kmb]?)\s+([a-zA-Z]+)\s+to\s+([a-zA-Z]+)`).
		Keywords("to").
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			amount, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			from := normalizeCurrencyName(matches[2])
			to := normalizeCurrencyName(matches[3])

			if from == "" || to == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// patternAmountCurrencyNameInCode: "100 dollars in euros", "50 pounds in yen"
func patternAmountCurrencyNameInCode() *Pattern {
	return NewPattern("amount_currency_name_in_code").
		Regex(`(\d+(?:[.,]\d+)?[kmb]?)\s+([a-zA-Z]+)\s+in\s+([a-zA-Z]+)`).
		Keywords("in").
		Priority(75).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			amount, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			from := normalizeCurrencyName(matches[2])
			to := normalizeCurrencyName(matches[3])

			if from == "" || to == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + from + " in " + to, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// CRYPTO PATTERNS
// ════════════════════════════════════════════════════════════════

// patternCryptoToFiat: "1 BTC to USD", "0.5 ETH to EUR", "1 bitcoin to dollars"
func patternCryptoToFiat() *Pattern {
	return NewPattern("crypto_to_fiat").
		Regex(`(\d+(?:\.\d+)?)\s*(btc|eth|bitcoin|ethereum|ether|xrp|ripple|ltc|litecoin|doge|dogecoin|sol|solana|ada|cardano)\s+(?:to|in)\s+([a-zA-Z]+)`).
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			amount, ok := ExtractNumber(matches[1])
			if !ok {
				return "", false
			}

			crypto := normalizeCryptoName(matches[2])
			fiat := normalizeCurrencyName(matches[3])

			if crypto == "" || fiat == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + crypto + " in " + fiat, true
		}).
		Build()
}

// patternFiatToCrypto: "$1000 to BTC", "1000 USD to bitcoin"
func patternFiatToCrypto() *Pattern {
	return NewPattern("fiat_to_crypto").
		Regex(`([$€£¥]?)(\d+(?:[.,]\d+)?[kmb]?)\s*([a-zA-Z]*)\s+(?:to|in)\s+(btc|eth|bitcoin|ethereum|ether|xrp|ripple|ltc|litecoin|doge|dogecoin|sol|solana|ada|cardano)`).
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			symbol := matches[1]
			amount, ok := ExtractNumber(matches[2])
			if !ok {
				return "", false
			}

			fiat := matches[3]
			if fiat == "" {
				fiat = symbolToCurrencyCode(symbol)
			} else {
				fiat = normalizeCurrencyName(fiat)
			}

			crypto := normalizeCryptoName(matches[4])

			if fiat == "" || crypto == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + fiat + " in " + crypto, true
		}).
		Build()
}

// ════════════════════════════════════════════════════════════════
// METAL PATTERNS
// ════════════════════════════════════════════════════════════════

// patternMetalToFiat: "1 oz gold to USD", "1 gram gold in USD", "1 kg silver to EUR"
func patternMetalToFiat() *Pattern {
	return NewPattern("metal_to_fiat").
		Regex(`(\d+(?:\.\d+)?)\s*(?:(gram|grams|g|kg|kilogram|kilograms|oz|ounce|ounces)\s+)?(gold|silver|platinum|palladium|xau|xag|xpt|xpd)\s+(?:to|in)\s+([a-zA-Z]+)`).
		Priority(80).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 5 {
				return "", false
			}

			amount, ok := ExtractNumber(matches[1])
			if !ok {
				amount = 1
			}

			weightUnit := strings.ToLower(matches[2])
			amount = convertWeightToTroyOz(amount, weightUnit)

			metal := normalizeMetalName(matches[3])
			fiat := normalizeCurrencyName(matches[4])

			if metal == "" || fiat == "" {
				return "", false
			}

			return FormatFloat(amount) + " " + metal + " in " + fiat, true
		}).
		Build()
}

// patternTurkishGold: "1 ceyrek in usd", "çeyrek in lira", "2 tam in usd"
func patternTurkishGold() *Pattern {
	return NewPattern("turkish_gold").
		Regex(`(\d+(?:\.\d+)?)?\s*(?:altin\s+)?(ceyrek|çeyrek|yarim|yarım|tam|ata|cumhuriyet)\s*(?:altin)?\s+(?:to|in)\s+([a-zA-Z]+)`).
		Priority(85).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 4 {
				return "", false
			}

			amount := 1.0
			if matches[1] != "" {
				var ok bool
				amount, ok = ExtractNumber(matches[1])
				if !ok {
					amount = 1
				}
			}

			denom := strings.ToLower(matches[2])
			troyOz := turkishGoldToTroyOz(denom)
			if troyOz == 0 {
				return "", false
			}

			totalOz := amount * troyOz

			fiat := normalizeCurrencyName(matches[3])
			if fiat == "" {
				return "", false
			}

			return FormatFloat(totalOz) + " XAU in " + fiat, true
		}).
		Build()
}

// patternMetalPriceIn: "gold in usd", "gold price in eur", "silver in lira"
func patternMetalPriceIn() *Pattern {
	return NewPattern("metal_price_in").
		Regex(`(gold|silver|platinum|palladium|xau|xag|xpt|xpd)\s+(?:price\s+)?(?:to|in)\s+([a-zA-Z]+)`).
		Priority(70).
		Handler(func(input string, matches []string) (string, bool) {
			if len(matches) < 3 {
				return "", false
			}

			metal := normalizeMetalName(matches[1])
			fiat := normalizeCurrencyName(matches[2])

			if metal == "" || fiat == "" {
				return "", false
			}

			return "1 " + metal + " in " + fiat, true
		}).
		Build()
}

// convertWeightToTroyOz converts a weight amount to troy ounces.
// Empty or oz/ounce units return amount as-is.
func convertWeightToTroyOz(amount float64, unit string) float64 {
	switch unit {
	case "g", "gram", "grams":
		return amount / 31.1035
	case "kg", "kilogram", "kilograms":
		return amount * 1000 / 31.1035
	default:
		return amount // oz, ounce, ounces, or empty
	}
}

// turkishGoldToTroyOz returns the troy ounce equivalent of a Turkish gold denomination.
// Pure gold weights (22K, multiply by 22/24 for purity):
//   - Çeyrek: 1.75g → ~1.606g pure → 0.05164 oz
//   - Yarım: 3.50g → ~3.212g pure → 0.10328 oz
//   - Tam (Ata/Cumhuriyet): 7.216g → ~6.615g pure → 0.21267 oz
func turkishGoldToTroyOz(denom string) float64 {
	switch denom {
	case "ceyrek", "çeyrek":
		return 0.05164
	case "yarim", "yarım":
		return 0.10328
	case "tam", "ata", "cumhuriyet":
		return 0.21267
	default:
		return 0
	}
}

// ════════════════════════════════════════════════════════════════
// HELPER FUNCTIONS
// ════════════════════════════════════════════════════════════════

// symbolToCurrencyCode maps currency symbols to codes.
func symbolToCurrencyCode(symbol string) string {
	switch symbol {
	case "$":
		return "USD"
	case "€":
		return "EUR"
	case "£":
		return "GBP"
	case "¥":
		return "JPY"
	case "₿":
		return "BTC"
	default:
		return ""
	}
}

// normalizeCurrencyName converts currency names to codes.
func normalizeCurrencyName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	// Check if already a valid code
	upper := strings.ToUpper(name)
	if isValidCurrencyCode(upper) {
		return upper
	}

	// Map common names to codes
	names := map[string]string{
		// USD
		"dollar":  "USD",
		"dollars": "USD",
		"usd":     "USD",
		"buck":    "USD",
		"bucks":   "USD",

		// EUR
		"euro":  "EUR",
		"euros": "EUR",
		"eur":   "EUR",

		// GBP
		"pound":    "GBP",
		"pounds":   "GBP",
		"gbp":      "GBP",
		"sterling": "GBP",
		"quid":     "GBP",

		// JPY
		"yen": "JPY",
		"jpy": "JPY",

		// CNY
		"yuan":     "CNY",
		"renminbi": "CNY",
		"rmb":      "CNY",
		"cny":      "CNY",

		// CHF
		"franc":      "CHF",
		"francs":     "CHF",
		"chf":        "CHF",
		"swissfranc": "CHF",

		// CAD
		"cad":    "CAD",
		"loonie": "CAD",

		// AUD
		"aud": "AUD",

		// INR
		"rupee":  "INR",
		"rupees": "INR",
		"inr":    "INR",

		// KRW
		"won": "KRW",
		"krw": "KRW",

		// MXN
		"peso":  "MXN",
		"pesos": "MXN",
		"mxn":   "MXN",

		// BRL
		"real":  "BRL",
		"reais": "BRL",
		"brl":   "BRL",

		// TRY
		"lira":        "TRY",
		"try":         "TRY",
		"tl":          "TRY",
		"turkishlira": "TRY",

		// RUB
		"ruble":  "RUB",
		"rubles": "RUB",
		"rouble": "RUB",
		"rub":    "RUB",

		// SEK
		"krona":  "SEK",
		"kronor": "SEK",
		"sek":    "SEK",

		// NOK
		"nok": "NOK",

		// DKK
		"dkk": "DKK",

		// PLN
		"zloty": "PLN",
		"pln":   "PLN",

		// THB
		"baht": "THB",
		"thb":  "THB",

		// SGD
		"sgd": "SGD",

		// HKD
		"hkd": "HKD",

		// NZD
		"nzd": "NZD",

		// ZAR
		"rand": "ZAR",
		"zar":  "ZAR",
	}

	if code, ok := names[name]; ok {
		return code
	}

	return ""
}

// normalizeCryptoName converts crypto names to codes.
func normalizeCryptoName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	cryptos := map[string]string{
		"btc":       "BTC",
		"bitcoin":   "BTC",
		"eth":       "ETH",
		"ethereum":  "ETH",
		"ether":     "ETH",
		"xrp":       "XRP",
		"ripple":    "XRP",
		"ltc":       "LTC",
		"litecoin":  "LTC",
		"doge":      "DOGE",
		"dogecoin":  "DOGE",
		"sol":       "SOL",
		"solana":    "SOL",
		"ada":       "ADA",
		"cardano":   "ADA",
		"dot":       "DOT",
		"polkadot":  "DOT",
		"avax":      "AVAX",
		"avalanche": "AVAX",
		"matic":     "MATIC",
		"polygon":   "MATIC",
		"link":      "LINK",
		"chainlink": "LINK",
		"uni":       "UNI",
		"uniswap":   "UNI",
		"atom":      "ATOM",
		"cosmos":    "ATOM",
	}

	if code, ok := cryptos[name]; ok {
		return code
	}

	// Check if already a valid code
	upper := strings.ToUpper(name)
	if isValidCryptoCode(upper) {
		return upper
	}

	return ""
}

// normalizeMetalName converts metal names to codes.
func normalizeMetalName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	metals := map[string]string{
		"gold":      "XAU",
		"xau":       "XAU",
		"silver":    "XAG",
		"xag":       "XAG",
		"platinum":  "XPT",
		"xpt":       "XPT",
		"palladium": "XPD",
		"xpd":       "XPD",
	}

	if code, ok := metals[name]; ok {
		return code
	}

	return ""
}

// isValidCurrencyCode checks if a string is a valid currency code.
func isValidCurrencyCode(code string) bool {
	// Common fiat currency codes
	validCodes := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true, "CNY": true,
		"CHF": true, "CAD": true, "AUD": true, "NZD": true, "HKD": true,
		"SGD": true, "INR": true, "KRW": true, "MXN": true, "BRL": true,
		"TRY": true, "RUB": true, "ZAR": true, "SEK": true, "NOK": true,
		"DKK": true, "PLN": true, "THB": true, "MYR": true, "IDR": true,
		"PHP": true, "VND": true, "CZK": true, "HUF": true, "ILS": true,
		"CLP": true, "COP": true, "PEN": true, "ARS": true, "EGP": true,
		"AED": true, "SAR": true, "QAR": true, "KWD": true, "BHD": true,
		"OMR": true, "JOD": true, "LBP": true, "PKR": true, "BDT": true,
		"LKR": true, "NPR": true, "MMK": true, "KHR": true, "TWD": true,
	}

	return validCodes[code]
}

// isValidCryptoCode checks if a string is a valid crypto code.
func isValidCryptoCode(code string) bool {
	validCodes := map[string]bool{
		"BTC": true, "ETH": true, "XRP": true, "LTC": true, "BCH": true,
		"DOGE": true, "SOL": true, "ADA": true, "DOT": true, "AVAX": true,
		"MATIC": true, "LINK": true, "UNI": true, "ATOM": true, "XLM": true,
		"ALGO": true, "VET": true, "FIL": true, "THETA": true, "XMR": true,
		"EOS": true, "AAVE": true, "MKR": true, "COMP": true, "SNX": true,
		"YFI": true, "SUSHI": true, "CRV": true, "1INCH": true, "BAT": true,
	}

	return validCodes[code]
}
