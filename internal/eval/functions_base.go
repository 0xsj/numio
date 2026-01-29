// internal/eval/functions_base.go

package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// GENERAL BASE CONVERSION
// ════════════════════════════════════════════════════════════════

// FnToBase converts a number to a string representation in the given base.
// Args: value, base (2-36)
func FnToBase(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("tobase requires exactly 2 arguments: value, base")
	}

	val := int64(args[0].AsFloat())
	base := int(args[1].AsFloat())

	if base < 2 || base > 36 {
		return types.Error("tobase: base must be between 2 and 36")
	}

	result := strconv.FormatInt(val, base)

	return types.StringValue(strings.ToUpper(result))
}

// FnFromBase converts a string representation from the given base to decimal.
// Args: string value, base (2-36)
func FnFromBase(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("frombase requires exactly 2 arguments: value, base")
	}

	str := args[0].AsString()
	base := int(args[1].AsFloat())

	if base < 2 || base > 36 {
		return types.Error("frombase: base must be between 2 and 36")
	}

	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	str = strings.TrimPrefix(str, "0x")
	str = strings.TrimPrefix(str, "0b")
	str = strings.TrimPrefix(str, "0o")

	result, err := strconv.ParseInt(str, base, 64)
	if err != nil {
		return types.Errorf("frombase: invalid number '%s' for base %d", str, base)
	}

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// SPECIFIC BASE CONVERSIONS (TO)
// ════════════════════════════════════════════════════════════════

// FnToHex converts a number to hexadecimal string.
func FnToHex(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tohex requires exactly 1 argument")
	}

	val := int64(args[0].AsFloat())

	var result string
	if val < 0 {
		result = "-0x" + strings.ToUpper(strconv.FormatInt(-val, 16))
	} else {
		result = "0x" + strings.ToUpper(strconv.FormatInt(val, 16))
	}

	return types.StringValue(result)
}

// FnToBin converts a number to binary string.
func FnToBin(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tobin requires exactly 1 argument")
	}

	val := int64(args[0].AsFloat())

	var result string
	if val < 0 {
		result = "-0b" + strconv.FormatInt(-val, 2)
	} else {
		result = "0b" + strconv.FormatInt(val, 2)
	}

	return types.StringValue(result)
}

// FnToOct converts a number to octal string.
func FnToOct(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("tooct requires exactly 1 argument")
	}

	val := int64(args[0].AsFloat())

	var result string
	if val < 0 {
		result = "-0o" + strconv.FormatInt(-val, 8)
	} else {
		result = "0o" + strconv.FormatInt(val, 8)
	}

	return types.StringValue(result)
}

// FnToDec converts a number to decimal (useful for display consistency).
func FnToDec(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("todec requires exactly 1 argument")
	}

	val := int64(args[0].AsFloat())
	return types.Number(float64(val))
}

// ════════════════════════════════════════════════════════════════
// PARSING FROM SPECIFIC BASES
// ════════════════════════════════════════════════════════════════

// FnHex parses a hexadecimal value (with or without 0x prefix).
func FnHex(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("hex requires exactly 1 argument")
	}

	str := args[0].AsString()
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	str = strings.TrimPrefix(str, "0x")
	str = strings.TrimPrefix(str, "#")

	result, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		return types.Errorf("hex: invalid hexadecimal '%s'", str)
	}

	return types.Number(float64(result))
}

// FnBin parses a binary value (with or without 0b prefix).
func FnBin(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bin requires exactly 1 argument")
	}

	str := args[0].AsString()
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "0b")
	str = strings.TrimPrefix(str, "0B")

	result, err := strconv.ParseInt(str, 2, 64)
	if err != nil {
		return types.Errorf("bin: invalid binary '%s'", str)
	}

	return types.Number(float64(result))
}

// FnOct parses an octal value (with or without 0o prefix).
func FnOct(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("oct requires exactly 1 argument")
	}

	str := args[0].AsString()
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "0o")
	str = strings.TrimPrefix(str, "0O")

	result, err := strconv.ParseInt(str, 8, 64)
	if err != nil {
		return types.Errorf("oct: invalid octal '%s'", str)
	}

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// FORMATTED OUTPUT
// ════════════════════════════════════════════════════════════════

// FnHexPad converts to hex with zero-padding to specified width.
func FnHexPad(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("hexpad requires 1 or 2 arguments: value, [width]")
	}

	val := uint64(args[0].AsFloat())
	width := 0
	if len(args) == 2 {
		width = int(args[1].AsFloat())
	}

	var result string
	if width > 0 {
		result = fmt.Sprintf("0x%0*X", width, val)
	} else {
		result = fmt.Sprintf("0x%X", val)
	}

	return types.StringValue(result)
}

// FnBinPad converts to binary with zero-padding to specified width.
func FnBinPad(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("binpad requires 1 or 2 arguments: value, [width]")
	}

	val := uint64(args[0].AsFloat())
	width := 0
	if len(args) == 2 {
		width = int(args[1].AsFloat())
	}

	binStr := strconv.FormatUint(val, 2)
	if width > len(binStr) {
		binStr = strings.Repeat("0", width-len(binStr)) + binStr
	}

	return types.StringValue("0b" + binStr)
}

// FnBinGroup formats binary with grouping (e.g., 4 bits at a time).
func FnBinGroup(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bingroup requires 1 or 2 arguments: value, [groupSize]")
	}

	val := uint64(args[0].AsFloat())
	groupSize := 4
	if len(args) == 2 {
		groupSize = int(args[1].AsFloat())
		if groupSize < 1 {
			groupSize = 4
		}
	}

	binStr := strconv.FormatUint(val, 2)

	remainder := len(binStr) % groupSize
	if remainder != 0 {
		binStr = strings.Repeat("0", groupSize-remainder) + binStr
	}

	var groups []string
	for i := 0; i < len(binStr); i += groupSize {
		end := i + groupSize
		if end > len(binStr) {
			end = len(binStr)
		}
		groups = append(groups, binStr[i:end])
	}

	return types.StringValue("0b" + strings.Join(groups, " "))
}

// ════════════════════════════════════════════════════════════════
// ASCII / CHARACTER CONVERSIONS
// ════════════════════════════════════════════════════════════════

// FnAscii returns the ASCII code of a character.
func FnAscii(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ascii requires exactly 1 argument")
	}

	str := args[0].AsString()
	if len(str) == 0 {
		return types.Error("ascii: empty string")
	}

	r := []rune(str)[0]
	return types.Number(float64(r))
}

// FnChr returns the character for an ASCII/Unicode code point.
func FnChr(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("chr requires exactly 1 argument")
	}

	code := int(args[0].AsFloat())
	if code < 0 || code > 0x10FFFF {
		return types.Error("chr: code point out of range")
	}

	return types.StringValue(string(rune(code)))
}

// FnOrd is an alias for ascii (Python-style).
func FnOrd(args []types.Value) types.Value {
	return FnAscii(args)
}

// ════════════════════════════════════════════════════════════════
// BYTE REPRESENTATION
// ════════════════════════════════════════════════════════════════

// FnBytes returns the byte representation of an integer.
func FnBytes(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 3 {
		return types.Error("bytes requires 1-3 arguments: value, [byteCount], [endian]")
	}

	val := uint64(args[0].AsFloat())

	byteCount := 0
	if len(args) >= 2 {
		byteCount = int(args[1].AsFloat())
	}

	bigEndian := true
	if len(args) >= 3 {
		endian := strings.ToLower(args[2].AsString())
		if endian == "little" || endian == "le" {
			bigEndian = false
		}
	}

	if byteCount == 0 {
		if val <= 0xFF {
			byteCount = 1
		} else if val <= 0xFFFF {
			byteCount = 2
		} else if val <= 0xFFFFFFFF {
			byteCount = 4
		} else {
			byteCount = 8
		}
	}

	bytes := make([]string, byteCount)
	for i := 0; i < byteCount; i++ {
		var b byte
		if bigEndian {
			b = byte(val >> (8 * (byteCount - 1 - i)))
		} else {
			b = byte(val >> (8 * i))
		}
		bytes[i] = fmt.Sprintf("%02X", b)
	}

	return types.StringValue(strings.Join(bytes, " "))
}

// FnFromBytes converts bytes to an integer.
func FnFromBytes(args []types.Value) types.Value {
	if len(args) < 1 {
		return types.Error("frombytes requires at least 1 argument")
	}

	if len(args) == 1 {
		str := args[0].AsString()
		str = strings.ReplaceAll(str, " ", "")
		str = strings.ReplaceAll(str, "0x", "")
		str = strings.ReplaceAll(str, "0X", "")

		result, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			return types.Errorf("frombytes: invalid hex string '%s'", str)
		}
		return types.Number(float64(result))
	}

	var result uint64
	for i, arg := range args {
		b := uint64(arg.AsFloat())
		if b > 255 {
			return types.Errorf("frombytes: byte %d out of range (0-255)", i)
		}
		result = (result << 8) | b
	}

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// NUMBER SYSTEM INFO
// ════════════════════════════════════════════════════════════════

// FnDigits returns the number of digits in a number in the given base.
func FnDigits(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("digits requires 1 or 2 arguments: value, [base]")
	}

	val := math.Abs(args[0].AsFloat())
	if val == 0 {
		return types.Number(1)
	}

	base := 10.0
	if len(args) == 2 {
		base = args[1].AsFloat()
		if base < 2 || base > 36 {
			return types.Error("digits: base must be between 2 and 36")
		}
	}

	digits := math.Floor(math.Log(val)/math.Log(base)) + 1
	return types.Number(digits)
}

// FnSumDigits returns the sum of digits in a number.
func FnSumDigits(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("sumdigits requires 1 or 2 arguments: value, [base]")
	}

	val := int64(math.Abs(args[0].AsFloat()))

	base := int64(10)
	if len(args) == 2 {
		base = int64(args[1].AsFloat())
		if base < 2 || base > 36 {
			return types.Error("sumdigits: base must be between 2 and 36")
		}
	}

	var sum int64
	for val > 0 {
		sum += val % base
		val /= base
	}

	return types.Number(float64(sum))
}

// FnReverseDigits reverses the digits of a number.
func FnReverseDigits(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("revdigits requires 1 or 2 arguments: value, [base]")
	}

	val := int64(args[0].AsFloat())
	negative := val < 0
	if negative {
		val = -val
	}

	base := int64(10)
	if len(args) == 2 {
		base = int64(args[1].AsFloat())
		if base < 2 || base > 36 {
			return types.Error("revdigits: base must be between 2 and 36")
		}
	}

	var result int64
	for val > 0 {
		result = result*base + val%base
		val /= base
	}

	if negative {
		result = -result
	}

	return types.Number(float64(result))
}

// FnIsPalindrome checks if a number is a palindrome in the given base.
func FnIsPalindrome(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("ispalindrome requires 1 or 2 arguments: value, [base]")
	}

	val := int64(math.Abs(args[0].AsFloat()))

	base := int64(10)
	if len(args) == 2 {
		base = int64(args[1].AsFloat())
		if base < 2 || base > 36 {
			return types.Error("ispalindrome: base must be between 2 and 36")
		}
	}

	original := val
	var reversed int64
	for val > 0 {
		reversed = reversed*base + val%base
		val /= base
	}

	if original == reversed {
		return types.Number(1)
	}
	return types.Number(0)
}
