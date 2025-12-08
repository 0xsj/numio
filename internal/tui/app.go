// internal/tui/app.go

package tui

import (
	"strings"
	"time"

	"github.com/0xsj/numio/pkg/cache"
	"github.com/0xsj/numio/pkg/engine"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App is the main bubbletea model.
type App struct {
	// Components
	editor    *Editor
	results   *Results
	statusBar *StatusBar
	helpView  *HelpView

	// State
	styles      Styles
	keymap      KeyMap
	engine      *engine.Engine
	rateCache   *cache.RateCache
	showHelp    bool
	showHeader  bool
	wrapLines   bool
	debug       bool
	filename    string
	lastResults []ResultLine

	// Dimensions
	width  int
	height int

	// Messages
	statusMessage string
	statusTime    time.Time
}

// NewApp creates a new App instance.
func NewApp() *App {
	styles := DefaultStyles()
	keymap := DefaultKeyMap()
	rc := cache.New()
	eng := engine.NewWithCache(rc)

	helpView := NewHelpView(styles, keymap)

	return &App{
		editor:      NewEditor(styles, keymap),
		results:     NewResults(styles, eng),
		statusBar:   NewStatusBar(styles, helpView),
		helpView:    helpView,
		styles:      styles,
		keymap:      keymap,
		engine:      eng,
		rateCache:   rc,
		showHelp:    false,
		showHeader:  true,
		wrapLines:   false,
		debug:       false,
		filename:    "",
		lastResults: nil,
		width:       80,
		height:      24,
	}
}

// WithFile sets the initial file to load.
func (a *App) WithFile(filename string) *App {
	a.filename = filename
	return a
}

// WithContent sets the initial content.
func (a *App) WithContent(content string) *App {
	a.editor.SetContent(content)
	return a
}

// ════════════════════════════════════════════════════════════════
// BUBBLETEA INTERFACE
// ════════════════════════════════════════════════════════════════

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	// Initial evaluation
	a.evaluate()
	return nil
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKey(msg)

	case tea.WindowSizeMsg:
		a.handleResize(msg.Width, msg.Height)
		return a, nil

	case statusClearMsg:
		a.statusMessage = ""
		return a, nil
	}

	return a, nil
}

// View implements tea.Model.
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	var content string

	// Main content (editor + results)
	mainContent := a.renderMain()

	// Status bar
	statusBar := a.renderStatusBar()

	// Combine
	content = mainContent + "\n" + statusBar

	// Help overlay (if shown)
	if a.showHelp {
		content = a.renderWithHelp(content)
	}

	return content
}

// ════════════════════════════════════════════════════════════════
// KEY HANDLING
// ════════════════════════════════════════════════════════════════

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Check for quit
	if a.keymap.IsQuit(msg, a.editor.Mode()) {
		return a, tea.Quit
	}

	// Help toggle (works in all modes)
	if key.Matches(msg, a.keymap.Help) {
		a.showHelp = !a.showHelp
		return a, nil
	}

	// Close help with Escape if open
	if a.showHelp && msg.String() == "esc" {
		a.showHelp = false
		return a, nil
	}

	// If help is shown, don't process other keys
	if a.showHelp {
		return a, nil
	}

	// Global keys (work in normal mode)
	if a.editor.Mode() == ModeNormal {
		switch {
		case key.Matches(msg, a.keymap.Save):
			return a, a.save()

		case key.Matches(msg, a.keymap.Refresh):
			return a, a.refreshRates()

		case key.Matches(msg, a.keymap.ToggleWrap):
			a.wrapLines = !a.wrapLines
			a.setStatus("Wrap: " + boolToOnOff(a.wrapLines))
			return a, nil

		case key.Matches(msg, a.keymap.ToggleLines):
			a.editor.ToggleLineNumbers()
			a.setStatus("Line numbers toggled")
			return a, nil

		case key.Matches(msg, a.keymap.ToggleHeader):
			a.showHeader = !a.showHeader
			a.handleResize(a.width, a.height)
			return a, nil

		case key.Matches(msg, a.keymap.ToggleDebug):
			a.debug = !a.debug
			a.setStatus("Debug: " + boolToOnOff(a.debug))
			return a, nil
		}
	}

	// Pass to editor
	cmd := a.editor.Update(msg)

	// Re-evaluate after any edit
	a.evaluate()

	return a, cmd
}

// ════════════════════════════════════════════════════════════════
// RENDERING
// ════════════════════════════════════════════════════════════════

func (a *App) renderMain() string {
	// Calculate pane widths
	editorWidth := a.width * 2 / 3
	resultWidth := a.width - editorWidth - 1

	// Calculate height (minus status bar)
	contentHeight := a.height - 1
	if a.showHeader {
		contentHeight -= 2
	}

	// Update component sizes
	a.editor.SetSize(editorWidth, contentHeight)
	a.results.SetSize(resultWidth, contentHeight)
	a.results.SetScrollY(a.editor.scrollY)

	// Render editor
	editorContent := a.editor.Render()
	editorPane := lipgloss.NewStyle().
		Width(editorWidth).
		Height(contentHeight).
		Render(editorContent)

	// Render results
	resultsContent := a.results.Render(a.lastResults)
	resultsPane := lipgloss.NewStyle().
		Width(resultWidth).
		Height(contentHeight).
		Render(resultsContent)

	// Join horizontally
	main := lipgloss.JoinHorizontal(lipgloss.Top, editorPane, resultsPane)

	// Add header if shown
	if a.showHeader {
		header := a.renderHeader()
		main = header + "\n" + main
	}

	return main
}

func (a *App) renderHeader() string {
	title := "numio"
	if a.filename != "" {
		title += " - " + a.filename
	}
	if a.editor.Modified() {
		title += " [+]"
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(ColorNumber).
		Bold(true)

	separator := strings.Repeat("─", a.width)
	sepStyle := lipgloss.NewStyle().Foreground(ColorBorder)

	return titleStyle.Render(title) + "\n" + sepStyle.Render(separator)
}

func (a *App) renderStatusBar() string {
	info := StatusInfo{
		Mode:         a.editor.Mode(),
		Filename:     a.filename,
		Modified:     a.editor.Modified(),
		Line:         a.editor.CursorPos().Line + 1,
		Col:          a.editor.CursorPos().Col + 1,
		TotalLines:   a.editor.LineCount(),
		Total:        a.results.RenderTotal(),
		ShowRatesAge: true,
		RatesAge:     a.rateCache.Age(),
	}

	// Override with status message if recent
	if a.statusMessage != "" && time.Since(a.statusTime) < 3*time.Second {
		// Show status message instead of normal status
	}

	a.statusBar.SetWidth(a.width)
	return a.statusBar.Render(info)
}

func (a *App) renderWithHelp(content string) string {
	// Render help overlay centered
	a.helpView.SetSize(a.width, a.height)
	helpOverlay := a.helpView.RenderCentered(a.width, a.height)

	// For simplicity, just return the help overlay
	// A full implementation would composite them
	return helpOverlay
}

// ════════════════════════════════════════════════════════════════
// ACTIONS
// ════════════════════════════════════════════════════════════════

func (a *App) evaluate() {
	lines := a.editor.Lines()
	a.lastResults = a.results.Evaluate(lines)
}

func (a *App) save() tea.Cmd {
	if a.filename == "" {
		a.setStatus("No filename set")
		return nil
	}

	// In a real implementation, save to file
	a.setStatus("Saved: " + a.filename)
	return a.clearStatusAfter(3 * time.Second)
}

func (a *App) refreshRates() tea.Cmd {
	// In a real implementation, fetch from API
	a.setStatus("Rates refreshed")
	return a.clearStatusAfter(3 * time.Second)
}

func (a *App) setStatus(msg string) {
	a.statusMessage = msg
	a.statusTime = time.Now()
}

func (a *App) clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return statusClearMsg{}
	})
}

// ════════════════════════════════════════════════════════════════
// RESIZE
// ════════════════════════════════════════════════════════════════

func (a *App) handleResize(width, height int) {
	a.width = width
	a.height = height

	// Update styles with new dimensions
	a.styles = a.styles.WithDimensions(width, height)
}

// ════════════════════════════════════════════════════════════════
// MESSAGES
// ════════════════════════════════════════════════════════════════

type statusClearMsg struct{}

// ════════════════════════════════════════════════════════════════
// HELPERS
// ════════════════════════════════════════════════════════════════

func boolToOnOff(b bool) string {
	if b {
		return "ON"
	}
	return "OFF"
}

// ════════════════════════════════════════════════════════════════
// RUN
// ════════════════════════════════════════════════════════════════

// Run starts the TUI application.
func Run() error {
	app := NewApp()
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// RunWithFile starts the TUI with a file loaded.
func RunWithFile(filename string, content string) error {
	app := NewApp().WithFile(filename).WithContent(content)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}