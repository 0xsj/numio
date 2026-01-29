// internal/eval/functions_string.go

package eval

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// BASE64 ENCODING
// ════════════════════════════════════════════════════════════════

// FnBase64Encode encodes a string to Base64.
func FnBase64Encode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("base64 requires exactly 1 argument")
	}

	input := args[0].AsString()
	encoded := base64.StdEncoding.EncodeToString([]byte(input))

	return types.StringValue(encoded)
}

// FnBase64Decode decodes a Base64 string.
func FnBase64Decode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("base64decode requires exactly 1 argument")
	}

	input := args[0].AsString()
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return types.Errorf("base64decode: invalid base64 string: %v", err)
	}

	return types.StringValue(string(decoded))
}

// FnBase64URLEncode encodes a string to URL-safe Base64.
func FnBase64URLEncode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("base64url requires exactly 1 argument")
	}

	input := args[0].AsString()
	encoded := base64.URLEncoding.EncodeToString([]byte(input))

	return types.StringValue(encoded)
}

// FnBase64URLDecode decodes a URL-safe Base64 string.
func FnBase64URLDecode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("base64urldecode requires exactly 1 argument")
	}

	input := args[0].AsString()
	decoded, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return types.Errorf("base64urldecode: invalid base64url string: %v", err)
	}

	return types.StringValue(string(decoded))
}

// ════════════════════════════════════════════════════════════════
// HASHING
// ════════════════════════════════════════════════════════════════

// FnMD5 computes the MD5 hash of a string.
func FnMD5(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("md5 requires exactly 1 argument")
	}

	input := args[0].AsString()
	hash := md5.Sum([]byte(input))

	return types.StringValue(hex.EncodeToString(hash[:]))
}

// FnSHA1 computes the SHA-1 hash of a string.
func FnSHA1(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sha1 requires exactly 1 argument")
	}

	input := args[0].AsString()
	hash := sha1.Sum([]byte(input))

	return types.StringValue(hex.EncodeToString(hash[:]))
}

// FnSHA256 computes the SHA-256 hash of a string.
func FnSHA256(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sha256 requires exactly 1 argument")
	}

	input := args[0].AsString()
	hash := sha256.Sum256([]byte(input))

	return types.StringValue(hex.EncodeToString(hash[:]))
}

// FnSHA512 computes the SHA-512 hash of a string.
func FnSHA512(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("sha512 requires exactly 1 argument")
	}

	input := args[0].AsString()
	hash := sha512.Sum512([]byte(input))

	return types.StringValue(hex.EncodeToString(hash[:]))
}

// FnHash computes a hash using the specified algorithm.
// Args: input, algorithm ("md5", "sha1", "sha256", "sha512")
func FnHash(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("hash requires 2 arguments: input, algorithm")
	}

	input := args[0].AsString()
	algo := strings.ToLower(args[1].AsString())

	var hashStr string
	switch algo {
	case "md5":
		hash := md5.Sum([]byte(input))
		hashStr = hex.EncodeToString(hash[:])
	case "sha1":
		hash := sha1.Sum([]byte(input))
		hashStr = hex.EncodeToString(hash[:])
	case "sha256":
		hash := sha256.Sum256([]byte(input))
		hashStr = hex.EncodeToString(hash[:])
	case "sha512":
		hash := sha512.Sum512([]byte(input))
		hashStr = hex.EncodeToString(hash[:])
	default:
		return types.Errorf("hash: unknown algorithm '%s'", algo)
	}

	return types.StringValue(hashStr)
}

// ════════════════════════════════════════════════════════════════
// URL ENCODING
// ════════════════════════════════════════════════════════════════

// FnURLEncode encodes a string for use in a URL.
func FnURLEncode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("urlencode requires exactly 1 argument")
	}

	input := args[0].AsString()
	encoded := url.QueryEscape(input)

	return types.StringValue(encoded)
}

// FnURLDecode decodes a URL-encoded string.
func FnURLDecode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("urldecode requires exactly 1 argument")
	}

	input := args[0].AsString()
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return types.Errorf("urldecode: invalid URL-encoded string: %v", err)
	}

	return types.StringValue(decoded)
}

// FnURLPathEncode encodes a string for use in a URL path.
func FnURLPathEncode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("urlpathencode requires exactly 1 argument")
	}

	input := args[0].AsString()
	encoded := url.PathEscape(input)

	return types.StringValue(encoded)
}

// FnURLPathDecode decodes a URL path-encoded string.
func FnURLPathDecode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("urlpathdecode requires exactly 1 argument")
	}

	input := args[0].AsString()
	decoded, err := url.PathUnescape(input)
	if err != nil {
		return types.Errorf("urlpathdecode: invalid URL path-encoded string: %v", err)
	}

	return types.StringValue(decoded)
}

// ════════════════════════════════════════════════════════════════
// HEX ENCODING
// ════════════════════════════════════════════════════════════════

// FnHexEncode encodes a string to hexadecimal.
func FnHexEncode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hexencode requires exactly 1 argument")
	}

	input := args[0].AsString()
	encoded := hex.EncodeToString([]byte(input))

	return types.StringValue(encoded)
}

// FnHexDecode decodes a hexadecimal string.
func FnHexDecode(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hexdecode requires exactly 1 argument")
	}

	input := args[0].AsString()
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return types.Errorf("hexdecode: invalid hex string: %v", err)
	}

	return types.StringValue(string(decoded))
}

// ════════════════════════════════════════════════════════════════
// UUID GENERATION
// ════════════════════════════════════════════════════════════════

// FnUUID generates a random UUID (v4).
func FnUUID(args []types.Value) types.Value {
	uuid := generateUUIDv4()
	return types.StringValue(uuid)
}

// FnUUIDv4 is an alias for UUID.
func FnUUIDv4(args []types.Value) types.Value {
	return FnUUID(args)
}

// generateUUIDv4 generates a random UUID v4.
func generateUUIDv4() string {
	// Use crypto/rand for better randomness
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(pseudoRandByte())
	}

	// Set version (4) and variant (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC 4122

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// pseudoRandByte generates a pseudo-random byte using a simple LCG.
// Note: For production, use crypto/rand instead.
var pseudoRandState uint64 = uint64(1)

func pseudoRandByte() uint8 {
	// Simple LCG for demonstration; in production use crypto/rand
	pseudoRandState = pseudoRandState*6364136223846793005 + 1442695040888963407
	return uint8(pseudoRandState >> 56)
}

// ════════════════════════════════════════════════════════════════
// STRING MANIPULATION
// ════════════════════════════════════════════════════════════════

// FnStrLen returns the length of a string.
func FnStrLen(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("strlen requires exactly 1 argument")
	}

	input := args[0].AsString()
	length := utf8.RuneCountInString(input)

	return types.Number(float64(length))
}

// FnByteLen returns the byte length of a string.
func FnByteLen(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bytelen requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.Number(float64(len(input)))
}

// FnUpper converts a string to uppercase.
func FnUpper(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("upper requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.StringValue(strings.ToUpper(input))
}

// FnLower converts a string to lowercase.
func FnLower(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("lower requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.StringValue(strings.ToLower(input))
}

// FnTitle converts a string to title case.
func FnTitle(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("title requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.StringValue(strings.Title(input))
}

// FnTrim removes leading and trailing whitespace.
func FnTrim(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("trim requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.StringValue(strings.TrimSpace(input))
}

// FnTrimLeft removes leading whitespace.
func FnTrimLeft(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("trimleft requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.StringValue(strings.TrimLeftFunc(input, unicode.IsSpace))
}

// FnTrimRight removes trailing whitespace.
func FnTrimRight(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("trimright requires exactly 1 argument")
	}

	input := args[0].AsString()
	return types.StringValue(strings.TrimRightFunc(input, unicode.IsSpace))
}

// FnReverse reverses a string.
func FnReverse(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("reverse requires exactly 1 argument")
	}

	input := args[0].AsString()
	runes := []rune(input)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return types.StringValue(string(runes))
}

// FnRepeat repeats a string n times.
func FnRepeat(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("repeat requires 2 arguments: string, count")
	}

	input := args[0].AsString()
	count := int(args[1].AsFloat())

	if count < 0 {
		count = 0
	}
	if count > 10000 {
		count = 10000 // Limit to prevent memory issues
	}

	return types.StringValue(strings.Repeat(input, count))
}

// FnReplace replaces occurrences of a substring.
func FnReplace(args []types.Value) types.Value {
	if len(args) < 3 || len(args) > 4 {
		return types.Error("replace requires 3-4 arguments: string, old, new, [count]")
	}

	input := args[0].AsString()
	old := args[1].AsString()
	new := args[2].AsString()
	count := -1 // Replace all

	if len(args) == 4 {
		count = int(args[3].AsFloat())
	}

	return types.StringValue(strings.Replace(input, old, new, count))
}

// FnSubstr extracts a substring.
// Args: string, start, [length]
func FnSubstr(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("substr requires 2-3 arguments: string, start, [length]")
	}

	input := args[0].AsString()
	runes := []rune(input)
	start := int(args[1].AsFloat())

	// Handle negative start (from end)
	if start < 0 {
		start = len(runes) + start
	}
	if start < 0 {
		start = 0
	}
	if start >= len(runes) {
		return types.StringValue("")
	}

	end := len(runes)
	if len(args) == 3 {
		length := int(args[2].AsFloat())
		if length < 0 {
			return types.StringValue("")
		}
		end = start + length
		if end > len(runes) {
			end = len(runes)
		}
	}

	return types.StringValue(string(runes[start:end]))
}

// FnContains checks if a string contains a substring.
func FnContains(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("contains requires 2 arguments: string, substring")
	}

	input := args[0].AsString()
	substr := args[1].AsString()

	if strings.Contains(input, substr) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnStartsWith checks if a string starts with a prefix.
func FnStartsWith(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("startswith requires 2 arguments: string, prefix")
	}

	input := args[0].AsString()
	prefix := args[1].AsString()

	if strings.HasPrefix(input, prefix) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnEndsWith checks if a string ends with a suffix.
func FnEndsWith(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("endswith requires 2 arguments: string, suffix")
	}

	input := args[0].AsString()
	suffix := args[1].AsString()

	if strings.HasSuffix(input, suffix) {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIndexOf finds the index of a substring.
func FnIndexOf(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("indexof requires 2 arguments: string, substring")
	}

	input := args[0].AsString()
	substr := args[1].AsString()

	index := strings.Index(input, substr)
	return types.Number(float64(index))
}

// FnLastIndexOf finds the last index of a substring.
func FnLastIndexOf(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("lastindexof requires 2 arguments: string, substring")
	}

	input := args[0].AsString()
	substr := args[1].AsString()

	index := strings.LastIndex(input, substr)
	return types.Number(float64(index))
}

// FnCount counts occurrences of a substring.
func FnCountStr(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("countstr requires 2 arguments: string, substring")
	}

	input := args[0].AsString()
	substr := args[1].AsString()

	count := strings.Count(input, substr)
	return types.Number(float64(count))
}

// FnSplit splits a string by a delimiter.
func FnSplit(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("split requires 2 arguments: string, delimiter")
	}

	input := args[0].AsString()
	delim := args[1].AsString()

	parts := strings.Split(input, delim)
	return types.StringValue(fmt.Sprintf("[%s]", strings.Join(parts, ", ")))
}

// FnJoin joins strings with a delimiter.
// Args: delimiter, strings...
func FnJoin(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("join requires at least 2 arguments: delimiter, strings...")
	}

	delim := args[0].AsString()
	parts := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		parts[i-1] = args[i].AsString()
	}

	return types.StringValue(strings.Join(parts, delim))
}

// FnPadLeft pads a string on the left.
func FnPadLeft(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("padleft requires 2-3 arguments: string, length, [char]")
	}

	input := args[0].AsString()
	length := int(args[1].AsFloat())
	padChar := " "
	if len(args) == 3 {
		padChar = args[2].AsString()
		if len(padChar) == 0 {
			padChar = " "
		}
	}

	runes := []rune(input)
	if len(runes) >= length {
		return types.StringValue(input)
	}

	padding := length - len(runes)
	padRunes := []rune(padChar)
	var result strings.Builder
	for i := 0; i < padding; i++ {
		result.WriteRune(padRunes[i%len(padRunes)])
	}
	result.WriteString(input)

	return types.StringValue(result.String())
}

// FnPadRight pads a string on the right.
func FnPadRight(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("padright requires 2-3 arguments: string, length, [char]")
	}

	input := args[0].AsString()
	length := int(args[1].AsFloat())
	padChar := " "
	if len(args) == 3 {
		padChar = args[2].AsString()
		if len(padChar) == 0 {
			padChar = " "
		}
	}

	runes := []rune(input)
	if len(runes) >= length {
		return types.StringValue(input)
	}

	padding := length - len(runes)
	padRunes := []rune(padChar)
	var result strings.Builder
	result.WriteString(input)
	for i := 0; i < padding; i++ {
		result.WriteRune(padRunes[i%len(padRunes)])
	}

	return types.StringValue(result.String())
}

// FnPadCenter centers a string.
func FnPadCenter(args []types.Value) types.Value {
	if len(args) < 2 || len(args) > 3 {
		return types.Error("padcenter requires 2-3 arguments: string, length, [char]")
	}

	input := args[0].AsString()
	length := int(args[1].AsFloat())
	padChar := " "
	if len(args) == 3 {
		padChar = args[2].AsString()
		if len(padChar) == 0 {
			padChar = " "
		}
	}

	runes := []rune(input)
	if len(runes) >= length {
		return types.StringValue(input)
	}

	totalPadding := length - len(runes)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	padRunes := []rune(padChar)
	var result strings.Builder

	for i := 0; i < leftPadding; i++ {
		result.WriteRune(padRunes[i%len(padRunes)])
	}
	result.WriteString(input)
	for i := 0; i < rightPadding; i++ {
		result.WriteRune(padRunes[i%len(padRunes)])
	}

	return types.StringValue(result.String())
}

// ════════════════════════════════════════════════════════════════
// STRING ANALYSIS
// ════════════════════════════════════════════════════════════════

// FnWordCount counts words in a string.
func FnWordCount(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("wordcount requires exactly 1 argument")
	}

	input := args[0].AsString()
	words := strings.Fields(input)

	return types.Number(float64(len(words)))
}

// FnLineCount counts lines in a string.
func FnLineCount(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("linecount requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(0)
	}

	lines := strings.Count(input, "\n") + 1
	return types.Number(float64(lines))
}

// FnIsAlpha checks if a string contains only letters.
func FnIsAlpha(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isalpha requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(0)
	}

	for _, r := range input {
		if !unicode.IsLetter(r) {
			return types.Number(0)
		}
	}
	return types.Number(1)
}

// FnIsNumeric checks if a string contains only digits.
func FnIsNumeric(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isnumeric requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(0)
	}

	for _, r := range input {
		if !unicode.IsDigit(r) {
			return types.Number(0)
		}
	}
	return types.Number(1)
}

// FnIsAlphaNum checks if a string contains only letters and digits.
func FnIsAlphaNum(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isalphanum requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(0)
	}

	for _, r := range input {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return types.Number(0)
		}
	}
	return types.Number(1)
}

// FnIsWhitespace checks if a string contains only whitespace.
func FnIsWhitespace(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("iswhitespace requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(1)
	}

	for _, r := range input {
		if !unicode.IsSpace(r) {
			return types.Number(0)
		}
	}
	return types.Number(1)
}

// FnIsUpper checks if a string is all uppercase.
func FnIsUpper(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("isupper requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(0)
	}

	hasLetter := false
	for _, r := range input {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return types.Number(0)
			}
		}
	}

	if hasLetter {
		return types.Number(1)
	}
	return types.Number(0)
}

// FnIsLower checks if a string is all lowercase.
func FnIsLower(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("islower requires exactly 1 argument")
	}

	input := args[0].AsString()
	if input == "" {
		return types.Number(0)
	}

	hasLetter := false
	for _, r := range input {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsLower(r) {
				return types.Number(0)
			}
		}
	}

	if hasLetter {
		return types.Number(1)
	}
	return types.Number(0)
}

// ════════════════════════════════════════════════════════════════
// FORMATTING
// ════════════════════════════════════════════════════════════════

// FnSlug creates a URL-friendly slug from a string.
func FnSlug(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("slug requires exactly 1 argument")
	}

	input := args[0].AsString()
	input = strings.ToLower(input)

	var result strings.Builder
	lastWasDash := false

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
			lastWasDash = false
		} else if !lastWasDash {
			result.WriteRune('-')
			lastWasDash = true
		}
	}

	slug := strings.Trim(result.String(), "-")
	return types.StringValue(slug)
}

// FnCamelCase converts a string to camelCase.
func FnCamelCase(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("camelcase requires exactly 1 argument")
	}

	input := args[0].AsString()
	words := splitIntoWords(input)

	if len(words) == 0 {
		return types.StringValue("")
	}

	var result strings.Builder
	result.WriteString(strings.ToLower(words[0]))

	for _, word := range words[1:] {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}

	return types.StringValue(result.String())
}

// FnPascalCase converts a string to PascalCase.
func FnPascalCase(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("pascalcase requires exactly 1 argument")
	}

	input := args[0].AsString()
	words := splitIntoWords(input)

	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])))
			result.WriteString(strings.ToLower(word[1:]))
		}
	}

	return types.StringValue(result.String())
}

// FnSnakeCase converts a string to snake_case.
func FnSnakeCase(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("snakecase requires exactly 1 argument")
	}

	input := args[0].AsString()
	words := splitIntoWords(input)

	var lowered []string
	for _, word := range words {
		lowered = append(lowered, strings.ToLower(word))
	}

	return types.StringValue(strings.Join(lowered, "_"))
}

// FnKebabCase converts a string to kebab-case.
func FnKebabCase(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("kebabcase requires exactly 1 argument")
	}

	input := args[0].AsString()
	words := splitIntoWords(input)

	var lowered []string
	for _, word := range words {
		lowered = append(lowered, strings.ToLower(word))
	}

	return types.StringValue(strings.Join(lowered, "-"))
}

// FnConstantCase converts a string to CONSTANT_CASE.
func FnConstantCase(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("constantcase requires exactly 1 argument")
	}

	input := args[0].AsString()
	words := splitIntoWords(input)

	var uppered []string
	for _, word := range words {
		uppered = append(uppered, strings.ToUpper(word))
	}

	return types.StringValue(strings.Join(uppered, "_"))
}

// splitIntoWords splits a string into words for case conversion.
func splitIntoWords(s string) []string {
	var words []string
	var current strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			// Start new word on uppercase
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}
