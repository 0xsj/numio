// internal/tui/help.go

package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpView renders the help overlay.
type HelpView struct {
	styles Styles
	keymap KeyMap
	width  int
	height int
}

// NewHelpView creates a new help view.
func NewHelpView(styles Styles, keymap KeyMap) *HelpView {
	return &HelpView{
		styles: styles,
		keymap: keymap,
		width:  60,
		height: 24,
	}
}

// SetSize updates the help view dimensions.
func (h *HelpView) SetSize(width, height int) {
	h.width = min(60, width-4)
	h.height = min(24, height-4)
}

// Render returns the help overlay content.
func (h *HelpView) Render() string {
	var content strings.Builder

	// Title
	title := h.styles.HelpTitle.Render("Help")
	content.WriteString(title)
	content.WriteString("\n")

	// Header row
	headerStyle := lipgloss.NewStyle().
		Foreground(ColorFgDim).
		Bold(true)
	header := headerStyle.Render("Key") +
		strings.Repeat(" ", 14) +
		headerStyle.Render("Action")
	content.WriteString(header)
	content.WriteString("\n")

	// Separator
	content.WriteString(strings.Repeat("─", h.width-4))
	content.WriteString("\n")

	// Sections
	sections := h.keymap.GetHelpBindings()

	for i, section := range sections {
		// Section header
		sectionHeader := h.styles.HelpSection.Render(section.Section)
		content.WriteString(sectionHeader)
		content.WriteString("\n")

		// Bindings
		for _, binding := range section.Bindings {
			help := binding.Help()
			keyStr := h.styles.HelpKey.Render(help.Key)
			descStr := h.styles.HelpDesc.Render(help.Desc)
			content.WriteString(keyStr + descStr)
			content.WriteString("\n")
		}

		// Add spacing between sections (except last)
		if i < len(sections)-1 {
			content.WriteString("\n")
		}
	}

	// Footer hint
	content.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(ColorFgDim).
		Italic(true).
		Render("Press ? or Esc to close")
	content.WriteString(footer)

	// Wrap in modal style
	modal := h.styles.HelpModal.
		Width(h.width).
		Render(content.String())

	return modal
}

// RenderCentered returns the help overlay centered in the given dimensions.
func (h *HelpView) RenderCentered(totalWidth, totalHeight int) string {
	modal := h.Render()

	// Calculate centering
	modalHeight := strings.Count(modal, "\n") + 1
	modalWidth := h.width + 4 // Account for border/padding

	padTop := (totalHeight - modalHeight) / 2
	padLeft := (totalWidth - modalWidth) / 2

	if padTop < 0 {
		padTop = 0
	}
	if padLeft < 0 {
		padLeft = 0
	}

	// Build centered output
	var result strings.Builder

	// Top padding
	for i := 0; i < padTop; i++ {
		result.WriteString("\n")
	}

	// Add left padding to each line
	lines := strings.Split(modal, "\n")
	leftPad := strings.Repeat(" ", padLeft)
	for i, line := range lines {
		result.WriteString(leftPad)
		result.WriteString(line)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// QuickHelp returns a short help string for the status bar.
func (h *HelpView) QuickHelp(mode Mode) string {
	var hints []string

	switch mode {
	case ModeNormal:
		hints = []string{
			"? help",
			"^s save",
			"^r rates",
		}
	case ModeInsert:
		hints = []string{
			"Esc normal",
			"^s save",
		}
	case ModeVisual:
		hints = []string{
			"Esc normal",
			"y yank",
			"d delete",
		}
	}

	style := lipgloss.NewStyle().Foreground(ColorFgDim)
	return style.Render(strings.Join(hints, "  "))
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}