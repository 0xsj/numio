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

	values := valuesToFloats(args)
	sparkline := graph.Sparkline(values)

	return types.StringValue(sparkline)
}

// FnSparklineStats generates a sparkline with statistics.
func FnSparklineStats(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("sparkstats requires at least one argument")
	}

	values := valuesToFloats(args)
	result := graph.SparklineWithStats(values)

	return types.StringValue(result.Full())
}

// FnSparklineTrend generates a sparkline with trend indicator.
func FnSparklineTrend(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("sparktrend requires at least one argument")
	}

	values := valuesToFloats(args)
	sparkline := graph.SparklineWithTrend(values)

	return types.StringValue(sparkline)
}

// FnSparklineStyled generates a sparkline with a specific style.
func FnSparklineStyled(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("sparkstyle requires at least 2 arguments: style and values")
	}

	style := graph.BlockStyle(int(args[0].AsFloat()))
	values := valuesToFloats(args[1:])
	sparkline := graph.SparklineStyled(values, style)

	return types.StringValue(sparkline)
}

// FnSparklineBounded generates a sparkline with explicit min/max bounds.
func FnSparklineBounded(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("sparkbound requires at least 3 arguments: min, max, and values")
	}

	minVal := args[0].AsFloat()
	maxVal := args[1].AsFloat()

	if minVal >= maxVal {
		return types.Error("sparkbound: min must be less than max")
	}

	values := valuesToFloats(args[2:])
	sparkline := graph.SparklineBounded(values, minVal, maxVal)

	return types.StringValue(sparkline)
}

// FnSparklineFixed generates a sparkline with fixed width.
func FnSparklineFixed(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("sparkwidth requires at least 2 arguments: width and values")
	}

	width := int(args[0].AsFloat())
	if width < 1 {
		return types.Error("sparkwidth: width must be at least 1")
	}
	if width > 100 {
		return types.Error("sparkwidth: width must be at most 100")
	}

	values := valuesToFloats(args[1:])
	sparkline := graph.SparklineFixed(values, width)

	return types.StringValue(sparkline)
}

// ════════════════════════════════════════════════════════════════
// HISTOGRAM FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnHistogram generates a histogram from values.
func FnHistogram(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("hist requires at least one argument")
	}

	values := valuesToFloats(args)
	histogram := graph.Histogram(values)

	return types.StringValue(histogram)
}

// FnHistogramBins generates a histogram with specified bin count.
func FnHistogramBins(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("histbins requires at least 2 arguments: bins and values")
	}

	bins := int(args[0].AsFloat())
	if bins < 2 {
		return types.Error("histbins: bins must be at least 2")
	}
	if bins > 50 {
		return types.Error("histbins: bins must be at most 50")
	}

	values := valuesToFloats(args[1:])
	histogram := graph.HistogramWithBins(values, bins)

	return types.StringValue(histogram)
}

// ════════════════════════════════════════════════════════════════
// GAUGE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnGauge generates a gauge visualization.
func FnGauge(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("gauge requires 1 or 2 arguments: value [, max]")
	}

	value := args[0].AsFloat()
	maxVal := 100.0

	if len(args) == 2 {
		maxVal = args[1].AsFloat()
	}

	if maxVal <= 0 {
		return types.Error("gauge: max must be positive")
	}

	result := graph.GaugeWithMax(value, maxVal)
	return types.StringValue(result)
}

// FnGaugeRange generates a gauge with min/max range.
func FnGaugeRange(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("gaugerange requires 3 arguments: value, min, max")
	}

	value := args[0].AsFloat()
	minVal := args[1].AsFloat()
	maxVal := args[2].AsFloat()

	if minVal >= maxVal {
		return types.Error("gaugerange: min must be less than max")
	}

	result := graph.GaugeWithRange(value, minVal, maxVal)
	return types.StringValue(result)
}

// FnBattery generates a battery-style gauge.
func FnBattery(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("battery requires exactly 1 argument: percent")
	}

	percent := args[0].AsFloat()
	result := graph.Battery(percent)

	return types.StringValue(result)
}

// FnMeter generates a meter visualization.
func FnMeter(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("meter requires 1 or 2 arguments: value [, max]")
	}

	value := args[0].AsFloat()
	maxVal := 100.0

	if len(args) == 2 {
		maxVal = args[1].AsFloat()
	}

	result := graph.Meter(value, maxVal)
	return types.StringValue(result)
}

// FnSignal generates a signal strength indicator.
func FnSignal(args []types.Value) types.Value {
	if len(args) != 1 {
		return types.Error("signal requires exactly 1 argument: strength (0-100)")
	}

	strength := args[0].AsFloat()
	result := graph.Signal(strength)

	return types.StringValue(result)
}

// FnStars generates a star rating.
func FnStars(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("stars requires 1 or 2 arguments: rating [, max]")
	}

	rating := args[0].AsFloat()
	maxRating := 5.0

	if len(args) == 2 {
		maxRating = args[1].AsFloat()
	}

	result := graph.Stars(rating, maxRating)
	return types.StringValue(result)
}

// FnHearts generates a hearts indicator.
func FnHearts(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("hearts requires 1 or 2 arguments: current [, max]")
	}

	current := args[0].AsFloat()
	maxVal := 5.0

	if len(args) == 2 {
		maxVal = args[1].AsFloat()
	}

	result := graph.Hearts(current, maxVal)
	return types.StringValue(result)
}

// ════════════════════════════════════════════════════════════════
// BAR FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnBar generates a simple horizontal bar.
func FnBar(args []types.Value) types.Value {
	if len(args) < 1 || len(args) > 2 {
		return types.Error("bar requires 1 or 2 arguments: value [, max]")
	}

	value := args[0].AsFloat()
	maxVal := 100.0

	if len(args) == 2 {
		maxVal = args[1].AsFloat()
	}

	if maxVal <= 0 {
		return types.Error("bar: max must be positive")
	}

	ratio := value / maxVal
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	barWidth := 20
	filled := int(ratio * float64(barWidth))

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

	ratio := value / maxVal
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	percent := int(ratio * 100)

	barWidth := 15
	filled := int(ratio * float64(barWidth))

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
	sb.WriteString(intToStr(percent))
	sb.WriteString("%")

	return types.StringValue(sb.String())
}

// ════════════════════════════════════════════════════════════════
// TREND FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnTrend generates a trend indicator.
func FnTrend(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("trend requires at least 2 arguments")
	}

	values := valuesToFloats(args)
	result := graph.TrendIndicator(values)

	return types.StringValue(result)
}

// FnTrendStyled generates a trend indicator with style.
func FnTrendStyled(args []types.Value) types.Value {
	if len(args) < 3 {
		return types.Error("trendstyle requires at least 3 arguments: style and values")
	}

	style := graph.TrendStyle(int(args[0].AsFloat()))
	values := valuesToFloats(args[1:])
	result := graph.TrendIndicatorStyled(values, style)

	return types.StringValue(result)
}

// FnChange generates a change indicator.
func FnChange(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("change requires at least 2 arguments")
	}

	values := valuesToFloats(args)
	result := graph.ChangeWithArrow(values)

	return types.StringValue(result)
}

// FnMiniTrend generates a mini trend line.
func FnMiniTrend(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("minitrend requires at least 2 arguments")
	}

	values := valuesToFloats(args)
	result := graph.MiniTrend(values)

	return types.StringValue(result)
}

// FnSlope calculates the slope of data.
func FnSlope(args []types.Value) types.Value {
	if len(args) < 2 {
		return types.Error("slope requires at least 2 arguments")
	}

	values := valuesToFloats(args)
	result := graph.Slope(values)

	return types.StringValue(result)
}

// FnDelta shows the difference between two values.
func FnDelta(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("delta requires exactly 2 arguments")
	}

	a := args[0].AsFloat()
	b := args[1].AsFloat()
	result := graph.Delta(a, b)

	return types.StringValue(result)
}

// FnDeltaPercent shows the percentage difference.
func FnDeltaPercent(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("deltapct requires exactly 2 arguments")
	}

	a := args[0].AsFloat()
	b := args[1].AsFloat()
	result := graph.DeltaPercent(a, b)

	return types.StringValue(result)
}

// ════════════════════════════════════════════════════════════════
// DOT PLOT FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnDotPlot generates a dot plot.
func FnDotPlot(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("dots requires at least one argument")
	}

	values := valuesToFloats(args)
	result := graph.DotPlot(values)

	return types.StringValue(result)
}

// FnDistDots generates a distribution dot plot.
func FnDistDots(args []types.Value) types.Value {
	if len(args) == 0 {
		return types.Error("distdots requires at least one argument")
	}

	values := valuesToFloats(args)
	result := graph.DistributionDots(values, 10)

	return types.StringValue(result)
}

// FnNumberLine shows a value on a number line.
func FnNumberLine(args []types.Value) types.Value {
	if len(args) != 3 {
		return types.Error("numline requires 3 arguments: value, min, max")
	}

	value := args[0].AsFloat()
	min := args[1].AsFloat()
	max := args[2].AsFloat()

	if min >= max {
		return types.Error("numline: min must be less than max")
	}

	result := graph.NumberLineLabeled(value, min, max, 20)
	return types.StringValue(result)
}

// FnScatter generates a scatter plot (returns multiline).
func FnScatter(args []types.Value) types.Value {
	if len(args) < 4 {
		return types.Error("scatter requires at least 4 arguments (x,y pairs)")
	}

	if len(args)%2 != 0 {
		return types.Error("scatter requires even number of arguments (x,y pairs)")
	}

	values := valuesToFloats(args)
	result := graph.ScatterPlot(values, 20, 10)

	return types.StringValue(result)
}

// ════════════════════════════════════════════════════════════════
// COMPARISON FUNCTIONS
// ════════════════════════════════════════════════════════════════

// FnCompare generates a comparison indicator.
func FnCompare(args []types.Value) types.Value {
	if len(args) != 2 {
		return types.Error("compare requires exactly 2 arguments")
	}

	a := args[0].AsFloat()
	b := args[1].AsFloat()
	result := graph.CompareValues(a, b)

	return types.StringValue(result)
}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

// valuesToFloats converts a slice of Values to float64.
func valuesToFloats(args []types.Value) []float64 {
	values := make([]float64, len(args))
	for i, arg := range args {
		values[i] = arg.AsFloat()
	}
	return values
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
