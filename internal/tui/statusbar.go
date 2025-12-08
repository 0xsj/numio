// internal/tui/statusbar.go

package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StatusBar renders the bottom status bar.
type StatusBar struct {
	styles   Styles
	helpView *HelpView
	width    int
}

// NewStatusBar creates a new status bar.
func NewStatusBar(styles Styles, helpView *HelpView) *StatusBar {
	return &StatusBar{
		styles:   styles,
		helpView: helpView,
		width:    80,
	}
}

// SetWidth updates the status bar width.
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// StatusInfo contains information to display in the status bar.
type StatusInfo struct {
	Mode         Mode
	Filename     string
	Modified     bool
	Line         int
	Col          int
	TotalLines   int
	Total        string
	LastCalc     string
	RatesAge     time.Duration
	ShowRatesAge bool
}

// Render returns the rendered status bar.
func (s *StatusBar) Render(info StatusInfo) string {
	// Left section: Mode + Filename
	left := s.renderLeft(info)

	// Center section: Quick help
	center := s.helpView.QuickHelp(info.Mode)

	// Right section: Total + Position
	right := s.renderRight(info)

	// Calculate spacing
	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)
	rightWidth := lipgloss.Width(right)

	availableSpace := s.width - leftWidth - rightWidth
	if availableSpace < centerWidth {
		// Not enough space, skip center
		center = ""
		availableSpace = s.width - leftWidth - rightWidth
	}

	// Build the status bar
	var result strings.Builder

	result.WriteString(left)

	if center != "" {
		// Center the help text
		leftPad := (availableSpace - centerWidth) / 2
		if leftPad > 0 {
			result.WriteString(strings.Repeat(" ", leftPad))
		}
		result.WriteString(center)
		rightPad := availableSpace - leftPad - centerWidth
		if rightPad > 0 {
			result.WriteString(strings.Repeat(" ", rightPad))
		}
	} else {
		// Just fill with spaces
		if availableSpace > 0 {
			result.WriteString(strings.Repeat(" ", availableSpace))
		}
	}

	result.WriteString(right)

	// Apply background style
	return s.styles.StatusBar.Width(s.width).Render(result.String())
}

// renderLeft renders the left section (mode + filename).
func (s *StatusBar) renderLeft(info StatusInfo) string {
	var parts []string

	// Mode indicator
	var modeStyle lipgloss.Style
	switch info.Mode {
	case ModeNormal:
		modeStyle = s.styles.ModeNormal
	case ModeInsert:
		modeStyle = s.styles.ModeInsert
	case ModeVisual:
		modeStyle = s.styles.ModeVisual
	}
	parts = append(parts, modeStyle.Render(info.Mode.String()))

	// Modified indicator
	if info.Modified {
		modStyle := lipgloss.NewStyle().
			Foreground(ColorResultNegative).
			Bold(true)
		parts = append(parts, modStyle.Render("•"))
	}

	// Filename or help hint
	if info.Filename != "" {
		fileStyle := s.styles.StatusText
		parts = append(parts, fileStyle.Render(info.Filename))
	}

	return strings.Join(parts, " ")
}

// renderRight renders the right section (total + rates + position).
func (s *StatusBar) renderRight(info StatusInfo) string {
	var parts []string

	// Total
	if info.Total != "" {
		totalLabel := s.styles.StatusHint.Render("total:")
		totalValue := s.styles.Total.Render(info.Total)
		parts = append(parts, totalLabel+" "+totalValue)
	}

	// Rates age
	if info.ShowRatesAge {
		ageStr := formatDuration(info.RatesAge)
		ageStyle := s.styles.StatusHint
		if info.RatesAge > time.Hour {
			ageStyle = ageStyle.Foreground(ColorResultNegative)
		}
		parts = append(parts, ageStyle.Render(ageStr))
	}

	// Position (line:col)
	if info.TotalLines > 0 {
		posStyle := s.styles.StatusHint
		pos := formatInt(info.Line) + ":" + formatInt(info.Col)
		parts = append(parts, posStyle.Render(pos))
	}

	return strings.Join(parts, "  ")
}

// formatDuration formats a duration for display.
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		secs := int(d.Seconds())
		return formatInt(secs) + "s"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		return formatInt(mins) + "m"
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		mins := int(d.Minutes()) % 60
		if mins > 0 {
			return formatInt(hours) + "h" + formatInt(mins) + "m"
		}
		return formatInt(hours) + "h"
	}
	days := int(d.Hours() / 24)
	return formatInt(days) + "d"
}

// formatInt formats an integer without fmt package.
func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + formatInt(-n)
	}

	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// formatFloat formats a float with given precision.
func formatFloat(f float64, precision int) string {
	if f < 0 {
		return "-" + formatFloat(-f, precision)
	}

	// Integer part
	intPart := int(f)
	result := formatInt(intPart)

	if precision <= 0 {
		return result
	}

	// Fractional part
	result += "."
	frac := f - float64(intPart)
	for i := 0; i < precision; i++ {
		frac *= 10
		digit := int(frac) % 10
		result += string(byte('0' + digit))
	}

	return result
}