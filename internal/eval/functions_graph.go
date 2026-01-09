// internal/eval/functions_graph.go

package eval

import (
	"strings"

	"github.com/0xsj/numio/internal/graph"
	"github.com/0xsj/numio/pkg/types"
)

// ════════════════════════════════════════════════════════════════
// SPARKLINE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnSparkline generates a sparkline from numeric arguments.
func FnSparkline(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("spark requires at least one argument")
	}

	// Convert arguments to float64 slice
	values := make([]float64, len(args))
	for i, arg := range args {
		if arg.IsError() {
			return arg
		}
		values[i] = arg.AsFloat()
	}

	// Generate sparkline
	sparkline := graph.Sparkline(values)

	return types.StringValue(sparkline)
}

// FnSparklineStats generates a sparkline with statistics.
func FnSparklineStats(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("sparkstats requires at least one argument")
	}

	// Convert arguments to float64 slice
	values := make([]float64, len(args))
	for i, arg := range args {
		if arg.IsError() {
			return arg
		}
		values[i] = arg.AsFloat()
	}

	// Generate sparkline with stats
	result := graph.SparklineWithStats(values)

	return types.StringValue(result.Full())
}

// FnSparklineTrend generates a sparkline with trend indicator.
func FnSparklineTrend(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("sparktrend requires at least one argument")
	}

	// Convert arguments to float64 slice
	values := make([]float64, len(args))
	for i, arg := range args {
		if arg.IsError() {
			return arg
		}
		values[i] = arg.AsFloat()
	}

	// Generate sparkline with trend
	sparkline := graph.SparklineWithTrend(values)

	return types.StringValue(sparkline)
}

// FnSparklineStyled generates a sparkline with a specific style.
func FnSparklineStyled(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("sparkstyle requires at least 2 arguments: style and values")
	}

	// First argument is style name (as a number mapping or we parse from context)
	// For simplicity, we'll use numeric style selection:
	// 0 = blocks (default), 1 = braille, 2 = bars, 3 = dots, 4 = ascii
	styleArg := args[0]
	style := graph.BlockStyle(int(styleArg.AsFloat()))

	// Rest are values
	values := make([]float64, len(args)-1)
	for i, arg := range args[1:] {
		if arg.IsError() {
			return arg
		}
		values[i] = arg.AsFloat()
	}

	// Generate sparkline with style
	sparkline := graph.SparklineStyled(values, style)

	return types.StringValue(sparkline)
}

// FnSparklineBounded generates a sparkline with explicit min/max bounds.
func FnSparklineBounded(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("sparkbound requires at least 3 arguments: min, max, and values")
	}

	// First two arguments are min and max
	minVal := args[0].AsFloat()
	maxVal := args[1].AsFloat()

	if minVal >= maxVal {
		return types.Error("sparkbound: min must be less than max")
	}

	// Rest are values
	values := make([]float64, len(args)-2)
	for i, arg := range args[2:] {
		if arg.IsError() {
			return arg
		}
		values[i] = arg.AsFloat()
	}

	// Generate sparkline with bounds
	sparkline := graph.SparklineBounded(values, minVal, maxVal)

	return types.StringValue(sparkline)
}

// FnSparklineFixed generates a sparkline with fixed width.
func FnSparklineFixed(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("sparkwidth requires at least 2 arguments: width and values")
	}

	// First argument is width
	width := int(args[0].AsFloat())
	if width < 1 {
		return types.Error("sparkwidth: width must be at least 1")
	}
	if width > 100 {
		return types.Error("sparkwidth: width must be at most 100")
	}

	// Rest are values
	values := make([]float64, len(args)-1)
	for i, arg := range args[1:] {
		if arg.IsError() {
			return arg
		}
		values[i] = arg.AsFloat()
	}

	// Generate sparkline with fixed width
	sparkline := graph.SparklineFixed(values, width)

	return types.StringValue(sparkline)
}

// ════════════════════════════════════════════════════════════════
// BAR CHART FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnBar generates a simple horizontal bar.
func FnBar(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bar requires 1 or 2 arguments: value [, max]")
	}

	value := args[0].AsFloat()
	maxVal := 100.0 // Default max

	if len(args) == 2 {
		maxVal = args[1].AsFloat()
	}

	if maxVal <= 0 {
		return types.Error("bar: max must be positive")
	}

	// Calculate bar length (max 20 chars)
	barWidth := 20
	ratio := value / maxVal
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * float64(barWidth))

	// Build bar
	var sb strings.Builder
	sb.WriteString("[")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			sb.WriteString("█")
		} else {
			sb.WriteString("░")
		}
	}
	sb.WriteString("]")

	return types.StringValue(sb.String())
}

// FnProgress generates a progress bar with percentage.
func FnProgress(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("progress requires 1 or 2 arguments: value [, max]")
	}

	value := args[0].AsFloat()
	maxVal := 100.0

	if len(args) == 2 {
		maxVal = args[1].AsFloat()
	}

	if maxVal <= 0 {
		return types.Error("progress: max must be positive")
	}

	// Calculate percentage
	ratio := value / maxVal
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	percent := int(ratio * 100)

	// Calculate bar length (max 15 chars to leave room for percentage)
	barWidth := 15
	filled := int(ratio * float64(barWidth))

	// Build progress bar
	var sb strings.Builder
	sb.WriteString("[")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			sb.WriteString("█")
		} else {
			sb.WriteString("░")
		}
	}
	sb.WriteString("] ")

	// Add percentage
	sb.WriteString(intToStr(percent))
	sb.WriteString("%")

	return types.StringValue(sb.String())
}

// intToStr converts int to string without fmt.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
