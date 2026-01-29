// internal/eval/functions_bitwise.go

package eval

import (
	"math"
	"math/bits"

	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// BASIC BITWISE OPERATIONS
// ════════════════════════════════════════════════════════════════

// FnBitAnd performs bitwise AND on two integers.
func FnBitAnd(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("band requires at least 2 arguments")
	}

	result := int64(args[0].AsFloat())
	for _, arg := range args[1:] {
		result &= int64(arg.AsFloat())
	}

	return types.Number(float64(result))
}

// FnBitOr performs bitwise OR on two or more integers.
func FnBitOr(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("bor requires at least 2 arguments")
	}

	result := int64(args[0].AsFloat())
	for _, arg := range args[1:] {
		result |= int64(arg.AsFloat())
	}

	return types.Number(float64(result))
}

// FnBitXor performs bitwise XOR on two or more integers.
func FnBitXor(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("bxor requires at least 2 arguments")
	}

	result := int64(args[0].AsFloat())
	for _, arg := range args[1:] {
		result ^= int64(arg.AsFloat())
	}

	return types.Number(float64(result))
}

// FnBitNot performs bitwise NOT (complement) on an integer.
func FnBitNot(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bnot requires exactly 1 argument")
	}

	val := int64(args[0].AsFloat())
	result := ^val

	return types.Number(float64(result))
}

// FnBitNand performs bitwise NAND on two integers.
func FnBitNand(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bnand requires exactly 2 arguments")
	}

	a := int64(args[0].AsFloat())
	b := int64(args[1].AsFloat())
	result := ^(a & b)

	return types.Number(float64(result))
}

// FnBitNor performs bitwise NOR on two integers.
func FnBitNor(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bnor requires exactly 2 arguments")
	}

	a := int64(args[0].AsFloat())
	b := int64(args[1].AsFloat())
	result := ^(a | b)

	return types.Number(float64(result))
}

// FnBitXnor performs bitwise XNOR on two integers.
func FnBitXnor(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bxnor requires exactly 2 arguments")
	}

	a := int64(args[0].AsFloat())
	b := int64(args[1].AsFloat())
	result := ^(a ^ b)

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// SHIFT OPERATIONS
// ════════════════════════════════════════════════════════════════

// FnLeftShift performs left bit shift.
func FnLeftShift(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("lshift requires exactly 2 arguments: value, positions")
	}

	val := int64(args[0].AsFloat())
	shift := uint(args[1].AsFloat())

	if shift > 63 {
		return types.Error("lshift: shift amount too large (max 63)")
	}

	result := val << shift

	return types.Number(float64(result))
}

// FnRightShift performs right bit shift (arithmetic).
func FnRightShift(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rshift requires exactly 2 arguments: value, positions")
	}

	val := int64(args[0].AsFloat())
	shift := uint(args[1].AsFloat())

	if shift > 63 {
		return types.Error("rshift: shift amount too large (max 63)")
	}

	result := val >> shift

	return types.Number(float64(result))
}

// FnUnsignedRightShift performs unsigned right bit shift.
func FnUnsignedRightShift(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("urshift requires exactly 2 arguments: value, positions")
	}

	val := uint64(args[0].AsFloat())
	shift := uint(args[1].AsFloat())

	if shift > 63 {
		return types.Error("urshift: shift amount too large (max 63)")
	}

	result := val >> shift

	return types.Number(float64(result))
}

// FnRotateLeft performs left bit rotation.
func FnRotateLeft(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rotl requires exactly 2 arguments: value, positions")
	}

	val := uint64(args[0].AsFloat())
	k := int(args[1].AsFloat())

	result := bits.RotateLeft64(val, k)

	return types.Number(float64(result))
}

// FnRotateRight performs right bit rotation.
func FnRotateRight(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("rotr requires exactly 2 arguments: value, positions")
	}

	val := uint64(args[0].AsFloat())
	k := int(args[1].AsFloat())

	result := bits.RotateLeft64(val, -k)

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// BIT COUNTING & ANALYSIS
// ════════════════════════════════════════════════════════════════

// FnPopCount counts the number of set bits (1s) in an integer.
func FnPopCount(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("popcount requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	count := bits.OnesCount64(val)

	return types.Number(float64(count))
}

// FnLeadingZeros counts leading zero bits.
func FnLeadingZeros(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("clz requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	count := bits.LeadingZeros64(val)

	return types.Number(float64(count))
}

// FnTrailingZeros counts trailing zero bits.
func FnTrailingZeros(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ctz requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	count := bits.TrailingZeros64(val)

	return types.Number(float64(count))
}

// FnBitLength returns the minimum number of bits needed to represent the value.
func FnBitLength(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bitlen requires exactly 1 argument")
	}

	val := uint64(math.Abs(args[0].AsFloat()))
	if val == 0 {
		return types.Number(0)
	}

	length := bits.Len64(val)

	return types.Number(float64(length))
}

// FnByteLength returns the minimum number of bytes needed to represent the value.
func FnByteLength(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bytelen requires exactly 1 argument")
	}

	val := uint64(math.Abs(args[0].AsFloat()))
	if val == 0 {
		return types.Number(1)
	}

	bitLen := bits.Len64(val)
	byteLen := (bitLen + 7) / 8

	return types.Number(float64(byteLen))
}

// ════════════════════════════════════════════════════════════════
// BIT MANIPULATION
// ════════════════════════════════════════════════════════════════

// FnBitGet gets the bit at a specific position (0-indexed from right).
func FnBitGet(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bitget requires exactly 2 arguments: value, position")
	}

	val := uint64(args[0].AsFloat())
	pos := uint(args[1].AsFloat())

	if pos > 63 {
		return types.Error("bitget: position out of range (0-63)")
	}

	bit := (val >> pos) & 1

	return types.Number(float64(bit))
}

// FnBitSet sets the bit at a specific position to 1.
func FnBitSet(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bitset requires exactly 2 arguments: value, position")
	}

	val := uint64(args[0].AsFloat())
	pos := uint(args[1].AsFloat())

	if pos > 63 {
		return types.Error("bitset: position out of range (0-63)")
	}

	result := val | (1 << pos)

	return types.Number(float64(result))
}

// FnBitClear clears the bit at a specific position to 0.
func FnBitClear(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bitclear requires exactly 2 arguments: value, position")
	}

	val := uint64(args[0].AsFloat())
	pos := uint(args[1].AsFloat())

	if pos > 63 {
		return types.Error("bitclear: position out of range (0-63)")
	}

	result := val &^ (1 << pos)

	return types.Number(float64(result))
}

// FnBitToggle toggles the bit at a specific position.
func FnBitToggle(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bittoggle requires exactly 2 arguments: value, position")
	}

	val := uint64(args[0].AsFloat())
	pos := uint(args[1].AsFloat())

	if pos > 63 {
		return types.Error("bittoggle: position out of range (0-63)")
	}

	result := val ^ (1 << pos)

	return types.Number(float64(result))
}

// FnBitSlice extracts a range of bits [start:end] (inclusive).
func FnBitSlice(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("bitslice requires exactly 3 arguments: value, start, end")
	}

	val := uint64(args[0].AsFloat())
	start := uint(args[1].AsFloat())
	end := uint(args[2].AsFloat())

	if start > 63 || end > 63 {
		return types.Error("bitslice: positions out of range (0-63)")
	}

	if start > end {
		start, end = end, start
	}

	// Create mask for the range
	width := end - start + 1
	mask := (uint64(1) << width) - 1
	result := (val >> start) & mask

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// MASK GENERATION
// ════════════════════════════════════════════════════════════════

// FnBitMask creates a bitmask with n bits set to 1.
func FnBitMask(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bitmask requires exactly 1 argument: number of bits")
	}

	n := int(args[0].AsFloat())
	if n < 0 || n > 64 {
		return types.Error("bitmask: number of bits must be 0-64")
	}

	if n == 64 {
		return types.Number(float64(^uint64(0)))
	}

	mask := (uint64(1) << n) - 1

	return types.Number(float64(mask))
}

// FnBitMaskRange creates a bitmask from start to end positions.
func FnBitMaskRange(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("bitmaskrange requires exactly 2 arguments: start, end")
	}

	start := uint(args[0].AsFloat())
	end := uint(args[1].AsFloat())

	if start > 63 || end > 63 {
		return types.Error("bitmaskrange: positions out of range (0-63)")
	}

	if start > end {
		start, end = end, start
	}

	width := end - start + 1
	mask := ((uint64(1) << width) - 1) << start

	return types.Number(float64(mask))
}

// ════════════════════════════════════════════════════════════════
// BYTE OPERATIONS
// ════════════════════════════════════════════════════════════════

// FnByteSwap16 swaps bytes in a 16-bit value.
func FnByteSwap16(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bswap16 requires exactly 1 argument")
	}

	val := uint16(args[0].AsFloat())
	result := bits.ReverseBytes16(val)

	return types.Number(float64(result))
}

// FnByteSwap32 swaps bytes in a 32-bit value.
func FnByteSwap32(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bswap32 requires exactly 1 argument")
	}

	val := uint32(args[0].AsFloat())
	result := bits.ReverseBytes32(val)

	return types.Number(float64(result))
}

// FnByteSwap64 swaps bytes in a 64-bit value.
func FnByteSwap64(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bswap64 requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	result := bits.ReverseBytes64(val)

	return types.Number(float64(result))
}

// FnBitReverse reverses the bits in a value.
func FnBitReverse(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("bitrev requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	result := bits.Reverse64(val)

	return types.Number(float64(result))
}

// ════════════════════════════════════════════════════════════════
// POWER OF TWO UTILITIES
// ════════════════════════════════════════════════════════════════

// FnIsPowerOf2 checks if a number is a power of 2.
func FnIsPowerOf2(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ispow2 requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	if val == 0 {
		return types.Number(0) // false
	}

	isPow2 := (val & (val - 1)) == 0
	if isPow2 {
		return types.Number(1) // true
	}
	return types.Number(0) // false
}

// FnNextPowerOf2 returns the next power of 2 >= value.
func FnNextPowerOf2(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("nextpow2 requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	if val == 0 {
		return types.Number(1)
	}

	// If already power of 2, return it
	if (val & (val - 1)) == 0 {
		return types.Number(float64(val))
	}

	// Find next power of 2
	val--
	val |= val >> 1
	val |= val >> 2
	val |= val >> 4
	val |= val >> 8
	val |= val >> 16
	val |= val >> 32
	val++

	return types.Number(float64(val))
}

// FnPrevPowerOf2 returns the previous power of 2 <= value.
func FnPrevPowerOf2(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("prevpow2 requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	if val == 0 {
		return types.Number(0)
	}

	// Find highest set bit
	val |= val >> 1
	val |= val >> 2
	val |= val >> 4
	val |= val >> 8
	val |= val >> 16
	val |= val >> 32

	return types.Number(float64((val >> 1) + 1))
}

// FnLog2Int returns floor(log2(x)) for integers.
func FnLog2Int(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("ilog2 requires exactly 1 argument")
	}

	val := uint64(args[0].AsFloat())
	if val == 0 {
		return types.Error("ilog2: argument must be positive")
	}

	result := bits.Len64(val) - 1

	return types.Number(float64(result))
}
