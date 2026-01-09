// internal/tui/status.go

package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// RateStatus represents the current state of rate fetching.
type RateStatus int

const (
	RateStatusIdle RateStatus = iota
	RateStatusFetching
	RateStatusSuccess
	RateStatusError
)

// RateStatusInfo holds rate status information for display.
type RateStatusInfo struct {
	Status    RateStatus
	Message   string
	UpdatedAt time.Time
	Error     error
}

// SpinnerFrames for the fetching animation.
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ════════════════════════════════════════════════════════════════
// MESSAGES
// ════════════════════════════════════════════════════════════════

// RateFetchStartMsg signals that rate fetching has started.
type RateFetchStartMsg struct{}

// RateFetchDoneMsg signals that rate fetching completed.
type RateFetchDoneMsg struct {
	Count int
	Err   error
}

// RateStatusClearMsg signals to clear the status message.
type RateStatusClearMsg struct{}

// SpinnerTickMsg triggers spinner animation update.
type SpinnerTickMsg struct{}

// ════════════════════════════════════════════════════════════════
// COMMANDS
// ════════════════════════════════════════════════════════════════

// StartRateFetch returns a command that fetches rates in the background.
func StartRateFetch(fetchFunc func() (int, error)) tea.Cmd {
	return func() tea.Msg {
		count, err := fetchFunc()
		return RateFetchDoneMsg{Count: count, Err: err}
	}
}

// ClearStatusAfter returns a command that clears status after a delay.
func ClearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return RateStatusClearMsg{}
	})
}

// SpinnerTick returns a command for spinner animation.
func SpinnerTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// ════════════════════════════════════════════════════════════════
// STATUS DISPLAY
// ════════════════════════════════════════════════════════════════

// FormatRateStatus formats the rate status for display in the status bar.
func FormatRateStatus(info RateStatusInfo, spinnerFrame int) string {
	switch info.Status {
	case RateStatusFetching:
		frame := SpinnerFrames[spinnerFrame%len(SpinnerFrames)]
		return frame + " Fetching rates..."

	case RateStatusSuccess:
		return "✓ " + info.Message

	case RateStatusError:
		if info.Error != nil {
			return "✗ " + info.Error.Error()
		}
		return "✗ Fetch failed"

	case RateStatusIdle:
		if !info.UpdatedAt.IsZero() {
			age := time.Since(info.UpdatedAt)
			return "⏱ " + formatAge(age)
		}
		return ""

	default:
		return ""
	}
}

// formatAge formats a duration as a human-readable age.
func formatAge(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		if mins == 1 {
			return "1m ago"
		}
		return intToStr(mins) + "m ago"
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return intToStr(hours) + "h ago"
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1d ago"
	}
	return intToStr(days) + "d ago"
}

// intToStr converts int to string without fmt.
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
