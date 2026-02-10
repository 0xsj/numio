// internal/explain/formatter.go

// Package explain provides step-by-step explanation of calculations.
package explain

import (
	"strings"

	"github.com/0xsj/numio/pkg/types"
)

// Formatter converts a Trace into human-readable output.
type Formatter struct {
	// UseUnicodeOps uses × instead of *, ÷ instead of /, etc.
	UseUnicodeOps bool

	// ShowIntermediateSteps shows all intermediate calculations.
	ShowIntermediateSteps bool

	// Verbose includes more detail in output.
	Verbose bool

	// Indent string for nested steps.
	Indent string
}

// DefaultFormatter returns a formatter with default settings.
func DefaultFormatter() *Formatter {
	return &Formatter{
		UseUnicodeOps:         true,
		ShowIntermediateSteps: true,
		Verbose:               false,
		Indent:                "  ",
	}
}

// CompactFormatter returns a formatter for single-line output.
func CompactFormatter() *Formatter {
	return &Formatter{
		UseUnicodeOps:         true,
		ShowIntermediateSteps: false,
		Verbose:               false,
		Indent:                "",
	}
}

// VerboseFormatter returns a formatter with maximum detail.
func VerboseFormatter() *Formatter {
	return &Formatter{
		UseUnicodeOps:         true,
		ShowIntermediateSteps: true,
		Verbose:               true,
		Indent:                "    ",
	}
}

// ════════════════════════════════════════════════════════════════
// MAIN FORMATTING METHODS
// ════════════════════════════════════════════════════════════════

// Format formats a trace into a human-readable string.
func (f *Formatter) Format(trace *Trace) string {
	if trace == nil || trace.IsEmpty() {
		return ""
	}

	if f.Verbose {
		return f.formatVerbose(trace)
	}

	return f.formatCompact(trace)
}

// FormatOneLine formats the trace as a single-line equation chain.
// Example: "(0.20 × $150) + $50 = $30 + $50 = $80"
func (f *Formatter) FormatOneLine(trace *Trace) string {
	if trace == nil || trace.IsEmpty() {
		return ""
	}

	return f.buildEquationChain(trace)
}

// FormatSteps formats just the steps as a list.
// It iterates all recorded steps, skipping pure literals/variables,
// so intermediate operations (conversions, percentages) aren't lost.
func (f *Formatter) FormatSteps(trace *Trace) []string {
	if trace == nil || trace.IsEmpty() {
		return nil
	}

	var lines []string

	// Walk all recorded steps, keeping only operations (skip literals/variables)
	for _, step := range trace.Steps {
		if step == nil || step.Op == "" {
			continue
		}
		f.collectStepLines(step, 0, &lines)
	}

	// If nothing was collected (e.g. only literals), fall back to root
	if len(lines) == 0 && trace.Root != nil {
		f.collectStepLines(trace.Root, 0, &lines)
	}

	return lines
}

// ════════════════════════════════════════════════════════════════
// COMPACT FORMAT
// ════════════════════════════════════════════════════════════════

// formatCompact produces a single-line or minimal output.
func (f *Formatter) formatCompact(trace *Trace) string {
	return f.buildEquationChain(trace)
}

// buildEquationChain builds a chain like: "100 + 8 = 108"
func (f *Formatter) buildEquationChain(trace *Trace) string {
	if trace.Root == nil {
		return trace.FinalResult.String()
	}

	var parts []string

	// Collect the evaluation chain
	f.collectChainParts(trace.Root, &parts)

	// Add final result if different from last part
	finalStr := f.formatValue(trace.FinalResult)
	if len(parts) == 0 || parts[len(parts)-1] != finalStr {
		parts = append(parts, finalStr)
	}

	// Remove duplicates in sequence
	parts = f.deduplicateSequence(parts)

	return strings.Join(parts, " = ")
}

// collectChainParts recursively collects expression parts for the chain.
func (f *Formatter) collectChainParts(step *Step, parts *[]string) {
	if step == nil {
		return
	}

	// For binary operations, show the expression and result
	if step.Op != "" && step.Left != "" && step.Right != "" {
		// Format the expression
		op := step.Op
		if f.UseUnicodeOps {
			op = OperatorSymbol(op)
		}

		exprStr := step.Left + " " + op + " " + step.Right
		*parts = append(*parts, exprStr)

		// If there are sub-steps, recurse into them first
		if step.HasSubSteps() {
			for _, sub := range step.SubSteps {
				f.collectChainParts(sub, parts)
			}
		}

		// Add the result
		resultStr := f.formatValue(step.Result)
		if resultStr != exprStr {
			*parts = append(*parts, resultStr)
		}
	} else if step.HasSubSteps() {
		// Group or complex expression - recurse
		for _, sub := range step.SubSteps {
			f.collectChainParts(sub, parts)
		}
		*parts = append(*parts, f.formatValue(step.Result))
	} else {
		// Simple literal or variable
		*parts = append(*parts, f.formatValue(step.Result))
	}
}

// deduplicateSequence removes adjacent duplicates from a slice.
func (f *Formatter) deduplicateSequence(parts []string) []string {
	if len(parts) <= 1 {
		return parts
	}

	result := []string{parts[0]}
	for i := 1; i < len(parts); i++ {
		if parts[i] != parts[i-1] {
			result = append(result, parts[i])
		}
	}
	return result
}

// ════════════════════════════════════════════════════════════════
// VERBOSE FORMAT
// ════════════════════════════════════════════════════════════════

// formatVerbose produces detailed multi-line output.
func (f *Formatter) formatVerbose(trace *Trace) string {
	var sb strings.Builder

	sb.WriteString("Input: ")
	sb.WriteString(trace.Input)
	sb.WriteString("\n\n")

	sb.WriteString("Steps:\n")

	// Walk all recorded steps, keeping only operations
	rendered := false
	for _, step := range trace.Steps {
		if step == nil || step.Op == "" {
			continue
		}
		f.formatStepVerbose(&sb, step, 0)
		rendered = true
	}

	// Fall back to root if nothing was rendered
	if !rendered && trace.Root != nil {
		f.formatStepVerbose(&sb, trace.Root, 0)
	}

	sb.WriteString("\nResult: ")
	sb.WriteString(f.formatValue(trace.FinalResult))

	return sb.String()
}

// formatStepVerbose formats a single step with indentation.
func (f *Formatter) formatStepVerbose(sb *strings.Builder, step *Step, depth int) {
	if step == nil {
		return
	}

	indent := strings.Repeat(f.Indent, depth)

	// Format this step
	sb.WriteString(indent)
	sb.WriteString(f.formatStepLine(step))
	sb.WriteString("\n")

	// Format detail lines with extra indentation
	detailIndent := strings.Repeat(f.Indent, depth+1)
	for _, detail := range step.Details {
		sb.WriteString(detailIndent)
		sb.WriteString(detail)
		sb.WriteString("\n")
	}

	// Format sub-steps
	for _, sub := range step.SubSteps {
		f.formatStepVerbose(sb, sub, depth+1)
	}
}

// formatStepLine formats a single step as one line.
func (f *Formatter) formatStepLine(step *Step) string {
	if step.Op != "" && step.Left != "" && step.Right != "" {
		// Binary operation
		op := step.Op
		if f.UseUnicodeOps {
			op = OperatorSymbol(op)
		}
		return step.Left + " " + op + " " + step.Right + " = " + f.formatValue(step.Result)
	}

	if step.Op != "" && step.Left == "" {
		// Unary operation or function
		return step.ExprString + " = " + f.formatValue(step.Result)
	}

	// Literal or simple expression
	if step.ExprString != "" && step.ExprString != f.formatValue(step.Result) {
		return step.ExprString + " → " + f.formatValue(step.Result)
	}

	return f.formatValue(step.Result)
}

// collectStepLines collects all steps as formatted lines.
func (f *Formatter) collectStepLines(step *Step, depth int, lines *[]string) {
	if step == nil {
		return
	}

	indent := strings.Repeat(f.Indent, depth)
	*lines = append(*lines, indent+f.formatStepLine(step))

	// Collect detail lines with extra indentation
	detailIndent := strings.Repeat(f.Indent, depth+1)
	for _, detail := range step.Details {
		*lines = append(*lines, detailIndent+detail)
	}

	for _, sub := range step.SubSteps {
		f.collectStepLines(sub, depth+1, lines)
	}
}

// ════════════════════════════════════════════════════════════════
// VALUE FORMATTING
// ════════════════════════════════════════════════════════════════

// formatValue formats a types.Value for display.
func (f *Formatter) formatValue(v types.Value) string {
	return v.String()
}

// ════════════════════════════════════════════════════════════════
// CONVENIENCE FUNCTIONS
// ════════════════════════════════════════════════════════════════

// Format formats a trace with default settings.
func Format(trace *Trace) string {
	return DefaultFormatter().Format(trace)
}

// FormatOneLine formats a trace as a single line with default settings.
func FormatOneLine(trace *Trace) string {
	return DefaultFormatter().FormatOneLine(trace)
}

// FormatVerbose formats a trace with verbose settings.
func FormatVerbose(trace *Trace) string {
	return VerboseFormatter().Format(trace)
}

// FormatCompact formats a trace with minimal output.
func FormatCompact(trace *Trace) string {
	return CompactFormatter().Format(trace)
}

// ════════════════════════════════════════════════════════════════
// EXPRESSION STRING BUILDERS
// ════════════════════════════════════════════════════════════════

// BuildExpressionString builds a string representation of the full expression
// with values substituted for variables.
func BuildExpressionString(trace *Trace) string {
	if trace == nil || trace.Root == nil {
		return ""
	}
	return buildExprString(trace.Root)
}

// buildExprString recursively builds expression string.
func buildExprString(step *Step) string {
	if step == nil {
		return ""
	}

	// Binary operation
	if step.Op != "" && step.Left != "" && step.Right != "" {
		op := OperatorSymbol(step.Op)
		return step.Left + " " + op + " " + step.Right
	}

	// Has sub-steps - build from children
	if step.HasSubSteps() {
		var parts []string
		for _, sub := range step.SubSteps {
			parts = append(parts, buildExprString(sub))
		}
		if len(parts) == 1 {
			return "(" + parts[0] + ")"
		}
		return strings.Join(parts, " ")
	}

	// Leaf node
	return step.ExprString
}

// ════════════════════════════════════════════════════════════════
// SPECIAL FORMATS
// ════════════════════════════════════════════════════════════════

// FormatAsTree formats the trace as an ASCII tree.
func FormatAsTree(trace *Trace) string {
	if trace == nil || trace.Root == nil {
		return ""
	}

	var sb strings.Builder
	formatTreeNode(&sb, trace.Root, "", true)
	return sb.String()
}

// formatTreeNode formats a step as a tree node.
func formatTreeNode(sb *strings.Builder, step *Step, prefix string, isLast bool) {
	if step == nil {
		return
	}

	// Draw connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	sb.WriteString(prefix)
	sb.WriteString(connector)
	sb.WriteString(step.ExprString)
	sb.WriteString(" = ")
	sb.WriteString(step.Result.String())
	sb.WriteString("\n")

	// Prepare prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	// Format children
	for i, sub := range step.SubSteps {
		isLastChild := i == len(step.SubSteps)-1
		formatTreeNode(sb, sub, childPrefix, isLastChild)
	}
}
