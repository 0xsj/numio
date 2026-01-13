// internal/graph/dotplot.go

package graph

import (
	"math"
	"strings"
)

// DotStyle represents different dot rendering styles.
type DotStyle int

const (
	DotStyleCircle  DotStyle = iota // ○ ●
	DotStyleDiamond                 // ◇ ◆
	DotStyleSquare                  // □ ■
	DotStyleStar                    // ☆ ★
	DotStyleBraille                 // ⠀ ⣿
	DotStyleASCII                   // . *
)

// DotChars returns the empty and filled characters for a style.
func DotChars(style DotStyle) (empty, filled rune) {
	switch style {
	case DotStyleDiamond:
		return '◇', '◆'
	case DotStyleSquare:
		return '□', '■'
	case DotStyleStar:
		return '☆', '★'
	case DotStyleBraille:
		return '⠀', '⣿'
	case DotStyleASCII:
		return '.', '*'
	default:
		return '○', '●'
	}
}

// DotPlotOptions configures dot plot rendering.
type DotPlotOptions struct {
	// Width is the plot width
	Width int

	// Height is the plot height (for 2D plots)
	Height int

	// Style determines the character set
	Style DotStyle

	// Min/Max bounds (NaN = auto-detect)
	Min float64
	Max float64

	// ShowAxis adds axis markers
	ShowAxis bool
}

// DefaultDotPlotOptions returns default dot plot options.
func DefaultDotPlotOptions() DotPlotOptions {
	return DotPlotOptions{
		Width:    20,
		Height:   1,
		Style:    DotStyleCircle,
		Min:      math.NaN(),
		Max:      math.NaN(),
		ShowAxis: false,
	}
}

// ════════════════════════════════════════════════════════════════
// SINGLE ROW DOT PLOT
// ════════════════════════════════════════════════════════════════

// DotPlot generates a single-row dot plot showing value positions.
func DotPlot(values []float64) string {
	return DotPlotWithOptions(values, DefaultDotPlotOptions())
}

// DotPlotWithOptions generates a dot plot with custom options.
func DotPlotWithOptions(values []float64, opts DotPlotOptions) string {
	if len(values) == 0 {
		return ""
	}

	// Determine bounds
	min, max := opts.Min, opts.Max
	if math.IsNaN(min) || math.IsNaN(max) {
		detectedMin, detectedMax := detectBounds(values)
		if math.IsNaN(min) {
			min = detectedMin
		}
		if math.IsNaN(max) {
			max = detectedMax
		}
	}

	// Handle flat data
	if min >= max {
		max = min + 1
	}

	// Get characters
	empty, filled := DotChars(opts.Style)

	// Initialize plot row
	width := opts.Width
	if width <= 0 {
		width = 20
	}

	plot := make([]rune, width)
	for i := range plot {
		plot[i] = empty
	}

	// Place dots
	dataRange := max - min
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}

		// Calculate position
		pos := int(((v - min) / dataRange) * float64(width-1))
		if pos < 0 {
			pos = 0
		}
		if pos >= width {
			pos = width - 1
		}

		plot[pos] = filled
	}

	// Build result
	var sb strings.Builder

	if opts.ShowAxis {
		sb.WriteString("[")
	}

	for _, r := range plot {
		sb.WriteRune(r)
	}

	if opts.ShowAxis {
		sb.WriteString("]")
	}

	return sb.String()
}

// ════════════════════════════════════════════════════════════════
// SCATTER PLOT (2D)
// ════════════════════════════════════════════════════════════════

// ScatterPlot generates a 2D scatter plot from x,y pairs.
// Values should be alternating x,y pairs: x1,y1,x2,y2,...
func ScatterPlot(values []float64, width, height int) string {
	if len(values) < 2 {
		return ""
	}

	// Extract x,y pairs
	n := len(values) / 2
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i := 0; i < n; i++ {
		xs[i] = values[i*2]
		ys[i] = values[i*2+1]
	}

	// Find bounds
	minX, maxX := detectBounds(xs)
	minY, maxY := detectBounds(ys)

	if minX >= maxX {
		maxX = minX + 1
	}
	if minY >= maxY {
		maxY = minY + 1
	}

	// Create grid
	if width <= 0 {
		width = 20
	}
	if height <= 0 {
		height = 10
	}

	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Plot points
	rangeX := maxX - minX
	rangeY := maxY - minY

	for i := 0; i < n; i++ {
		x := int(((xs[i] - minX) / rangeX) * float64(width-1))
		y := int(((ys[i] - minY) / rangeY) * float64(height-1))

		// Invert y (top = high)
		y = height - 1 - y

		if x >= 0 && x < width && y >= 0 && y < height {
			grid[y][x] = '●'
		}
	}

	// Build result
	var sb strings.Builder
	for i, row := range grid {
		for _, c := range row {
			sb.WriteRune(c)
		}
		if i < len(grid)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ════════════════════════════════════════════════════════════════
// DISTRIBUTION DOT PLOT
// ════════════════════════════════════════════════════════════════

// DistributionDots shows value distribution as stacked dots.
func DistributionDots(values []float64, bins int) string {
	if len(values) == 0 {
		return ""
	}

	if bins <= 0 {
		bins = 10
	}

	min, max := detectBounds(values)
	if min >= max {
		max = min + 1
	}

	// Count values in each bin
	counts := make([]int, bins)
	binWidth := (max - min) / float64(bins)

	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}

		idx := int((v - min) / binWidth)
		if idx < 0 {
			idx = 0
		}
		if idx >= bins {
			idx = bins - 1
		}
		counts[idx]++
	}

	// Find max count
	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}

	if maxCount == 0 {
		return strings.Repeat("○", bins)
	}

	// Generate output using dot density
	var sb strings.Builder

	for _, c := range counts {
		if c == 0 {
			sb.WriteString("○")
		} else if c == maxCount {
			sb.WriteString("●")
		} else if float64(c)/float64(maxCount) > 0.5 {
			sb.WriteString("◉")
		} else {
			sb.WriteString("◎")
		}
	}

	return sb.String()
}

// ════════════════════════════════════════════════════════════════
// NUMBER LINE
// ════════════════════════════════════════════════════════════════

// NumberLine shows a single value's position on a number line.
func NumberLine(value, min, max float64, width int) string {
	if width <= 0 {
		width = 20
	}

	if min >= max {
		max = min + 1
	}

	// Calculate position
	ratio := (value - min) / (max - min)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	pos := int(ratio * float64(width-1))

	// Build number line
	var sb strings.Builder
	sb.WriteString("[")

	for i := 0; i < width; i++ {
		if i == pos {
			sb.WriteString("●")
		} else if i == 0 || i == width-1 || i == width/2 {
			sb.WriteString("┼")
		} else {
			sb.WriteString("─")
		}
	}

	sb.WriteString("]")

	return sb.String()
}

// NumberLineLabeled shows a number line with min/max labels.
func NumberLineLabeled(value, min, max float64, width int) string {
	line := NumberLine(value, min, max, width)

	var sb strings.Builder
	sb.WriteString(formatFloat(min))
	sb.WriteString(" ")
	sb.WriteString(line)
	sb.WriteString(" ")
	sb.WriteString(formatFloat(max))

	return sb.String()
}

// ════════════════════════════════════════════════════════════════
// RANGE INDICATOR
// ════════════════════════════════════════════════════════════════

// RangeIndicator shows a value within a range.
func RangeIndicator(value, min, max, rangeMin, rangeMax float64) string {
	width := 20

	// Build range indicator
	var sb strings.Builder
	sb.WriteString("[")

	lineMin := min
	lineMax := max

	dataRange := lineMax - lineMin
	if dataRange <= 0 {
		dataRange = 1
	}

	// Calculate positions
	rangeStartPos := int(((rangeMin - lineMin) / dataRange) * float64(width-1))
	rangeEndPos := int(((rangeMax - lineMin) / dataRange) * float64(width-1))
	valuePos := int(((value - lineMin) / dataRange) * float64(width-1))

	// Clamp positions
	clamp := func(v, min, max int) int {
		if v < min {
			return min
		}
		if v > max {
			return max
		}
		return v
	}

	rangeStartPos = clamp(rangeStartPos, 0, width-1)
	rangeEndPos = clamp(rangeEndPos, 0, width-1)
	valuePos = clamp(valuePos, 0, width-1)

	for i := 0; i < width; i++ {
		if i == valuePos {
			sb.WriteString("●")
		} else if i >= rangeStartPos && i <= rangeEndPos {
			sb.WriteString("█")
		} else {
			sb.WriteString("░")
		}
	}

	sb.WriteString("]")

	return sb.String()
}
