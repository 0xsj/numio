// internal/tui/results.go

package tui

import (
	"strings"

	"github.com/0xsj/numio/pkg/engine"
	"github.com/0xsj/numio/pkg/types"
	"github.com/charmbracelet/lipgloss"
)

// Results renders the results pane.
type Results struct {
	styles      Styles
	highlighter *Highlighter
	engine      *engine.Engine
	width       int
	height      int
	scrollY     int
}

// NewResults creates a new results pane.
func NewResults(styles Styles, eng *engine.Engine) *Results {
	return &Results{
		styles:      styles,
		highlighter: NewHighlighter(styles),
		engine:      eng,
		width:       30,
		height:      24,
		scrollY:     0,
	}
}

// SetSize updates the results pane dimensions.
func (r *Results) SetSize(width, height int) {
	r.width = width
	r.height = height
}

// SetScrollY sets the vertical scroll position.
func (r *Results) SetScrollY(scrollY int) {
	r.scrollY = scrollY
}

// ResultLine represents a single result to display.
type ResultLine struct {
	Value    types.Value
	Text     string
	IsError  bool
	IsEmpty  bool
}

// Evaluate evaluates lines and returns results.
func (r *Results) Evaluate(lines []string) []ResultLine {
	// Clear previous state for fresh evaluation
	r.engine.Clear()

	results := make([]ResultLine, len(lines))

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Empty line
		if trimmed == "" {
			results[i] = ResultLine{
				Value:   types.Empty(),
				Text:    "",
				IsEmpty: true,
			}
			continue
		}

		// Comment-only line
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			results[i] = ResultLine{
				Value:   types.Empty(),
				Text:    "",
				IsEmpty: true,
			}
			continue
		}

		// Evaluate the line
		value := r.engine.Eval(line)

		if value.IsEmpty() {
			results[i] = ResultLine{
				Value:   value,
				Text:    "",
				IsEmpty: true,
			}
		} else if value.IsError() {
			results[i] = ResultLine{
				Value:   value,
				Text:    value.ErrorMessage(),
				IsError: true,
			}
		} else {
			results[i] = ResultLine{
				Value:   value,
				Text:    value.String(),
				IsError: false,
			}
		}
	}

	return results
}

// Render returns the rendered results pane.
func (r *Results) Render(results []ResultLine) string {
	var lines []string

	visibleLines := r.height - 1

	for i := 0; i < visibleLines; i++ {
		lineIdx := r.scrollY + i

		if lineIdx < len(results) {
			result := results[lineIdx]
			rendered := r.renderResult(result)
			lines = append(lines, rendered)
		} else {
			// Empty line
			lines = append(lines, "")
		}
	}

	// Join and apply right alignment
	content := strings.Join(lines, "\n")

	return r.styles.ResultPane.
		Width(r.width).
		Height(r.height - 1).
		Render(content)
}

// renderResult renders a single result line.
func (r *Results) renderResult(result ResultLine) string {
	if result.IsEmpty {
		return ""
	}

	if result.IsError {
		// Truncate error if too long
		text := result.Text
		maxLen := r.width - 4
		if len(text) > maxLen && maxLen > 3 {
			text = text[:maxLen-3] + "..."
		}
		return r.styles.ResultError.Render(text)
	}

	// Format the result value
	text := result.Text

	// Determine style based on value
	var style lipgloss.Style

	if result.Value.IsNumeric() {
		num := result.Value.AsFloat()
		if num < 0 {
			style = r.styles.ResultNegative
		} else {
			style = r.styles.ResultPositive
		}
	} else {
		style = r.styles.ResultNeutral
	}

	return style.Render(text)
}

// RenderTotal renders the total line for the status bar.
func (r *Results) RenderTotal() string {
	total := r.engine.Total()

	if total.IsEmpty() || (total.IsNumeric() && total.AsFloat() == 0) {
		return ""
	}

	return total.String()
}

// RenderGroupedTotals renders grouped totals.
func (r *Results) RenderGroupedTotals() []string {
	totals := r.engine.GroupedTotals()

	var results []string
	for _, t := range totals {
		if !t.IsEmpty() && !(t.IsNumeric() && t.AsFloat() == 0) {
			results = append(results, t.String())
		}
	}

	return results
}

// GetLastResult returns the last non-empty result.
func (r *Results) GetLastResult(results []ResultLine) *ResultLine {
	for i := len(results) - 1; i >= 0; i-- {
		if !results[i].IsEmpty && !results[i].IsError {
			return &results[i]
		}
	}
	return nil
}

// Engine returns the engine instance.
func (r *Results) Engine() *engine.Engine {
	return r.engine
}

// Variables returns all defined variables.
func (r *Results) Variables() map[string]types.Value {
	return r.engine.Variables()
}