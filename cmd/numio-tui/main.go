// cmd/numio-tui/main.go

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/0xsj/numio/internal/tui"
)

const (
	appName    = "numio"
	appVersion = "0.1.0"
)

func main() {
	// Parse command line arguments
	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help", "help":
			printHelp()
			return

		case "-v", "--version", "version":
			fmt.Printf("%s v%s\n", appName, appVersion)
			return
		}
	}

	// Check if a file was provided
	var filename string
	var content string

	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		filename = args[0]

		// Try to read the file
		data, err := os.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				// New file, start with empty content
				content = ""
			} else {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}
		} else {
			content = string(data)
		}
	}

	// Run the TUI
	var err error
	if filename != "" {
		err = tui.RunWithFile(filename, content)
	} else {
		err = tui.Run()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`%s v%s - Natural Math Calculator (TUI)

Usage:
  %s                    Start with empty buffer
  %s <file>             Open file for editing
  %s -h, --help         Show this help
  %s -v, --version      Show version

Keyboard Shortcuts:
  Navigation:
    h/j/k/l or arrows   Move cursor
    0 / $               Start / End of line
    gg / G              Top / Bottom of file
    w / b               Next / Previous word
    PgUp / PgDn         Page up / Page down

  Editing:
    i                   Insert mode
    a                   Append mode
    o / O               New line below / above
    dd                  Delete line
    x                   Delete character
    yy / p              Yank / Paste line
    u / Ctrl+r          Undo / Redo

  General:
    Esc                 Normal mode
    ? / F1              Toggle help
    Ctrl+s              Save file
    Ctrl+r              Refresh rates
    W                   Toggle wrap
    N                   Toggle line numbers
    H                   Toggle header
    q                   Quit

Examples:
  %s                        Start fresh
  %s budget.calc            Open budget.calc
  %s ~/finances/taxes.calc  Open with path

`, appName, appVersion, appName, appName, appName, appName, appName, appName, appName)
}