// internal/tui/app.go

package tui

import (
	"fmt"
	"strings"

	"github.com/0xsj/numio/pkg/engine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	lineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666"))
	commentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Italic(true)
	resultStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7ee787"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f85149"))
	cursorStyle  = lipgloss.NewStyle().Reverse(true)
	tildeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#444"))

	// Help styles
	helpBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#79c0ff")).
			Padding(1, 2)
	helpTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#79c0ff"))
	helpSectionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffa657")).MarginTop(1)
	helpKeyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#79c0ff")).Width(14)
	helpDescStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
	helpFooterStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Italic(true).MarginTop(1)
)

// App is the main model
type App struct {
	lines    []string
	row      int
	col      int
	mode     Mode
	width    int
	height   int
	engine   *engine.Engine
	showHelp bool
}

// NewApp creates a new app
func NewApp() *App {
	return &App{
		lines:    []string{""},
		row:      0,
		col:      0,
		mode:     ModeInsert,
		width:    80,
		height:   24,
		engine:   engine.New(),
		showHelp: false,
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		// Help toggle - works in any mode
		if msg.String() == "?" && a.mode == ModeNormal {
			a.showHelp = !a.showHelp
			return a, nil
		}
		if msg.String() == "f1" {
			a.showHelp = !a.showHelp
			return a, nil
		}

		// Close help with Esc or q
		if a.showHelp {
			if msg.String() == "esc" || msg.String() == "q" || msg.String() == "?" {
				a.showHelp = false
			}
			return a, nil
		}

		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return a, tea.Quit

		case "esc":
			if a.mode == ModeInsert {
				a.mode = ModeNormal
				if a.col > 0 {
					a.col--
				}
			}

		case "i":
			if a.mode == ModeNormal {
				a.mode = ModeInsert
			} else {
				a.insertChar('i')
			}

		case "a":
			if a.mode == ModeNormal {
				a.mode = ModeInsert
				if a.col < len(a.lines[a.row]) {
					a.col++
				}
			} else {
				a.insertChar('a')
			}

		case "o":
			if a.mode == ModeNormal {
				a.newLineBelow()
				a.mode = ModeInsert
			} else {
				a.insertChar('o')
			}

		case "O":
			if a.mode == ModeNormal {
				a.newLineAbove()
				a.mode = ModeInsert
			} else {
				a.insertChar('O')
			}

		case "q":
			if a.mode == ModeNormal {
				return a, tea.Quit
			} else {
				a.insertChar('q')
			}

		case "enter":
			if a.mode == ModeInsert {
				a.newLine()
			}

		case "backspace":
			if a.mode == ModeInsert {
				a.backspace()
			}

		case "delete":
			if a.mode == ModeInsert {
				a.deleteChar()
			}

		case "up", "k":
			if a.mode == ModeNormal || msg.String() == "up" {
				if a.row > 0 {
					a.row--
					a.clampCol()
				}
			} else {
				a.insertChar('k')
			}

		case "down", "j":
			if a.mode == ModeNormal || msg.String() == "down" {
				if a.row < len(a.lines)-1 {
					a.row++
					a.clampCol()
				}
			} else {
				a.insertChar('j')
			}

		case "left", "h":
			if a.mode == ModeNormal || msg.String() == "left" {
				if a.col > 0 {
					a.col--
				}
			} else {
				a.insertChar('h')
			}

		case "right", "l":
			if a.mode == ModeNormal || msg.String() == "right" {
				if a.col < len(a.lines[a.row]) {
					a.col++
				}
			} else {
				a.insertChar('l')
			}

		case "home", "0":
			if a.mode == ModeNormal || msg.String() == "home" {
				a.col = 0
			} else {
				a.insertChar('0')
			}

		case "end", "$":
			if a.mode == ModeNormal || msg.String() == "end" {
				a.col = len(a.lines[a.row])
			} else {
				a.insertChar('$')
			}

		case "G":
			if a.mode == ModeNormal {
				a.row = len(a.lines) - 1
				a.clampCol()
			} else {
				a.insertChar('G')
			}

		case "g":
			if a.mode == ModeNormal {
				a.row = 0
				a.col = 0
			} else {
				a.insertChar('g')
			}

		case "x":
			if a.mode == ModeNormal {
				a.deleteChar()
			} else {
				a.insertChar('x')
			}

		case "d":
			if a.mode == ModeNormal {
				a.deleteLine()
			} else {
				a.insertChar('d')
			}

		default:
			if a.mode == ModeInsert && len(msg.Runes) > 0 {
				for _, r := range msg.Runes {
					a.insertChar(r)
				}
			}
		}
	}

	return a, nil
}

func (a *App) insertChar(r rune) {
	line := a.lines[a.row]
	if a.col > len(line) {
		a.col = len(line)
	}
	a.lines[a.row] = line[:a.col] + string(r) + line[a.col:]
	a.col++
}

func (a *App) newLine() {
	line := a.lines[a.row]
	if a.col > len(line) {
		a.col = len(line)
	}
	before := line[:a.col]
	after := line[a.col:]
	a.lines[a.row] = before

	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row+1]...)
	newLines = append(newLines, after)
	if a.row+1 < len(a.lines) {
		newLines = append(newLines, a.lines[a.row+1:]...)
	}
	a.lines = newLines

	a.row++
	a.col = 0
}

func (a *App) newLineBelow() {
	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row+1]...)
	newLines = append(newLines, "")
	if a.row+1 < len(a.lines) {
		newLines = append(newLines, a.lines[a.row+1:]...)
	}
	a.lines = newLines
	a.row++
	a.col = 0
}

func (a *App) newLineAbove() {
	newLines := make([]string, 0, len(a.lines)+1)
	newLines = append(newLines, a.lines[:a.row]...)
	newLines = append(newLines, "")
	newLines = append(newLines, a.lines[a.row:]...)
	a.lines = newLines
	a.col = 0
}

func (a *App) backspace() {
	if a.col > 0 {
		line := a.lines[a.row]
		a.lines[a.row] = line[:a.col-1] + line[a.col:]
		a.col--
	} else if a.row > 0 {
		prevLen := len(a.lines[a.row-1])
		a.lines[a.row-1] += a.lines[a.row]
		a.lines = append(a.lines[:a.row], a.lines[a.row+1:]...)
		a.row--
		a.col = prevLen
	}
}

func (a *App) deleteChar() {
	line := a.lines[a.row]
	if a.col < len(line) {
		a.lines[a.row] = line[:a.col] + line[a.col+1:]
	} else if a.row < len(a.lines)-1 {
		a.lines[a.row] = line + a.lines[a.row+1]
		a.lines = append(a.lines[:a.row+1], a.lines[a.row+2:]...)
	}
}

func (a *App) deleteLine() {
	if len(a.lines) == 1 {
		a.lines[0] = ""
		a.col = 0
	} else {
		a.lines = append(a.lines[:a.row], a.lines[a.row+1:]...)
		if a.row >= len(a.lines) {
			a.row = len(a.lines) - 1
		}
		a.clampCol()
	}
}

func (a *App) clampCol() {
	if a.col > len(a.lines[a.row]) {
		a.col = len(a.lines[a.row])
	}
}

// View implements tea.Model
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	// Show help overlay
	if a.showHelp {
		return a.renderHelp()
	}

	var b strings.Builder

	contentHeight := a.height - 2
	if contentHeight < 1 {
		contentHeight = 20
	}

	lineNumWidth := 5
	resultWidth := 20
	editorWidth := a.width - lineNumWidth - resultWidth - 4

	if editorWidth < 20 {
		editorWidth = 20
	}

	a.engine.Clear()

	for i := 0; i < contentHeight; i++ {
		if i < len(a.lines) {
			b.WriteString(lineNumStyle.Render(fmt.Sprintf("%3d ", i+1)))
		} else {
			b.WriteString(lineNumStyle.Render("    "))
		}

		b.WriteString("│")

		var editorContent string
		var resultContent string

		if i < len(a.lines) {
			line := a.lines[i]

			if i == a.row {
				editorContent = a.renderLineWithCursor(line)
			} else {
				editorContent = a.highlightLine(line)
			}

			resultContent = a.evaluateLine(line)
		} else {
			editorContent = tildeStyle.Render("~")
			resultContent = ""
		}

		editorLen := lipgloss.Width(editorContent)
		if editorLen < editorWidth {
			editorContent += strings.Repeat(" ", editorWidth-editorLen)
		} else if editorLen > editorWidth {
			editorContent = editorContent[:editorWidth]
		}

		resultContent = fmt.Sprintf("%*s", resultWidth, resultContent)

		b.WriteString(editorContent)
		b.WriteString("│")
		b.WriteString(resultContent)
		b.WriteString("\n")
	}

	b.WriteString(a.renderStatusBar())

	return b.String()
}

func (a *App) renderHelp() string {
	var content strings.Builder

	// Title
	content.WriteString(helpTitleStyle.Render("Help"))
	content.WriteString("\n\n")

	// Navigation section
	content.WriteString(helpSectionStyle.Render("Navigation"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("↑/k") + helpDescStyle.Render("Move up") + "\n")
	content.WriteString(helpKeyStyle.Render("↓/j") + helpDescStyle.Render("Move down") + "\n")
	content.WriteString(helpKeyStyle.Render("←/h") + helpDescStyle.Render("Move left") + "\n")
	content.WriteString(helpKeyStyle.Render("→/l") + helpDescStyle.Render("Move right") + "\n")
	content.WriteString(helpKeyStyle.Render("0 / Home") + helpDescStyle.Render("Start of line") + "\n")
	content.WriteString(helpKeyStyle.Render("$ / End") + helpDescStyle.Render("End of line") + "\n")
	content.WriteString(helpKeyStyle.Render("g") + helpDescStyle.Render("Go to first line") + "\n")
	content.WriteString(helpKeyStyle.Render("G") + helpDescStyle.Render("Go to last line") + "\n")

	// Editing section
	content.WriteString(helpSectionStyle.Render("Editing"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("i") + helpDescStyle.Render("Insert mode") + "\n")
	content.WriteString(helpKeyStyle.Render("a") + helpDescStyle.Render("Append (insert after)") + "\n")
	content.WriteString(helpKeyStyle.Render("o") + helpDescStyle.Render("New line below") + "\n")
	content.WriteString(helpKeyStyle.Render("O") + helpDescStyle.Render("New line above") + "\n")
	content.WriteString(helpKeyStyle.Render("x") + helpDescStyle.Render("Delete character") + "\n")
	content.WriteString(helpKeyStyle.Render("d") + helpDescStyle.Render("Delete line") + "\n")
	content.WriteString(helpKeyStyle.Render("Esc") + helpDescStyle.Render("Normal mode") + "\n")

	// General section
	content.WriteString(helpSectionStyle.Render("General"))
	content.WriteString("\n")
	content.WriteString(helpKeyStyle.Render("?") + helpDescStyle.Render("Toggle help") + "\n")
	content.WriteString(helpKeyStyle.Render("F1") + helpDescStyle.Render("Toggle help") + "\n")
	content.WriteString(helpKeyStyle.Render("q") + helpDescStyle.Render("Quit (normal mode)") + "\n")
	content.WriteString(helpKeyStyle.Render("Ctrl+C") + helpDescStyle.Render("Force quit") + "\n")

	// Expressions section
	content.WriteString(helpSectionStyle.Render("Expressions"))
	content.WriteString("\n")
	content.WriteString(helpDescStyle.Render("100 + 50        → 150") + "\n")
	content.WriteString(helpDescStyle.Render("20% of 150      → 30") + "\n")
	content.WriteString(helpDescStyle.Render("$100 + 15%      → $115.00") + "\n")
	content.WriteString(helpDescStyle.Render("price = 100     → 100") + "\n")
	content.WriteString(helpDescStyle.Render("price * 2       → 200") + "\n")

	// Footer
	content.WriteString(helpFooterStyle.Render("\nPress ? or Esc to close"))

	// Wrap in border and center
	helpBox := helpBorderStyle.Render(content.String())

	return lipgloss.Place(a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		helpBox)
}

func (a *App) renderLineWithCursor(line string) string {
	col := a.col
	if col > len(line) {
		col = len(line)
	}

	if col == len(line) {
		return a.highlightLine(line) + cursorStyle.Render(" ")
	}

	before := line[:col]
	cursorChar := string(line[col])
	after := line[col+1:]

	return a.highlightLine(before) + cursorStyle.Render(cursorChar) + a.highlightLine(after)
}

func (a *App) highlightLine(line string) string {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return commentStyle.Render(line)
	}

	return line
}

func (a *App) evaluateLine(line string) string {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return ""
	}

	result := a.engine.Eval(line)

	if result.IsEmpty() {
		return ""
	}

	if result.IsError() {
		return errorStyle.Render("err")
	}

	return resultStyle.Render(result.String())
}

func (a *App) renderStatusBar() string {
	var modeStyle lipgloss.Style
	if a.mode == ModeInsert {
		modeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000")).
			Background(lipgloss.Color("#7ee787")).
			Padding(0, 1)
	} else {
		modeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000")).
			Background(lipgloss.Color("#79c0ff")).
			Padding(0, 1)
	}

	mode := modeStyle.Render(a.mode.String())

	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Render("  ? help  ^s save  ^r rates")

	pos := fmt.Sprintf("%d:%d", a.row+1, a.col+1)

	total := a.engine.Total()
	totalStr := ""
	if !total.IsEmpty() && total.AsFloat() != 0 {
		totalStr = resultStyle.Render(fmt.Sprintf("total: %s", total.String())) + "  "
	}

	left := mode + hint
	right := totalStr + pos

	spaces := a.width - lipgloss.Width(left) - lipgloss.Width(right)
	if spaces < 0 {
		spaces = 1
	}

	statusBg := lipgloss.NewStyle().Background(lipgloss.Color("#1a1a2e"))
	return statusBg.Render(left + strings.Repeat(" ", spaces) + right)
}

// Run starts the TUI
func Run() error {
	p := tea.NewProgram(NewApp(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// RunWithFile starts with file content
func RunWithFile(filename, content string) error {
	app := NewApp()
	if content != "" {
		app.lines = strings.Split(content, "\n")
	}
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}