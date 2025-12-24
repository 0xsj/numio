// cmd/numio-cli/main.go

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/0xsj/numio/pkg/engine"
	"github.com/0xsj/numio/pkg/types"
)

const (
	appName    = "numio"
	appVersion = "0.1.0"
)

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		handleArgs(os.Args[1:])
		return
	}

	// Start REPL
	runREPL()
}

// handleArgs processes command line arguments.
func handleArgs(args []string) {
	switch args[0] {
	case "-h", "--help", "help":
		printHelp()

	case "-v", "--version", "version":
		fmt.Printf("%s v%s\n", appName, appVersion)

	case "-e", "--eval":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: -e requires an expression")
			os.Exit(1)
		}
		// Evaluate expression and print result
		result := engine.QuickEval(strings.Join(args[1:], " "))
		printResult(result)

	case "-f", "--file":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Error: -f requires a filename")
			os.Exit(1)
		}
		runFile(args[1])

	default:
		// Treat as expression
		result := engine.QuickEval(strings.Join(args, " "))
		printResult(result)
	}
}

// runFile evaluates a file.
func runFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	eng := engine.New()
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		result := eng.Eval(line)
		if !result.IsEmpty() {
			if result.IsError() {
				fmt.Fprintf(os.Stderr, "Line %d: %s\n", i+1, result.ErrorMessage())
			} else {
				fmt.Println(result.String())
			}
		}
	}
}

// runREPL starts the interactive REPL.
func runREPL() {
	printBanner()

	eng := engine.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			// EOF or error
			fmt.Println()
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for commands
		if handleCommand(line, eng) {
			continue
		}

		// Evaluate expression
		result := eng.Eval(line)
		printResult(result)
	}
}

// handleCommand processes REPL commands. Returns true if it was a command.
func handleCommand(input string, eng *engine.Engine) bool {
	lower := strings.ToLower(input)

	switch {
	case lower == "quit" || lower == "exit" || lower == "q":
		fmt.Println("Goodbye!")
		os.Exit(0)

	case lower == "help" || lower == "?":
		printREPLHelp()
		return true

	case lower == "clear" || lower == "cls":
		eng.Clear()
		fmt.Println("Cleared.")
		return true

	case lower == "vars" || lower == "variables":
		printVariables(eng)
		return true

	case lower == "total":
		result := eng.Total()
		fmt.Printf("Total: %s\n", result.String())
		return true

	case lower == "totals":
		printGroupedTotals(eng)
		return true

	case lower == "history" || lower == "lines":
		printHistory(eng)
		return true

	case lower == "rates":
		printRateInfo(eng)
		return true

	case strings.HasPrefix(lower, "set "):
		handleSet(input[4:], eng)
		return true

	case strings.HasPrefix(lower, "del ") || strings.HasPrefix(lower, "delete "):
		name := strings.TrimPrefix(lower, "del ")
		name = strings.TrimPrefix(name, "delete ")
		name = strings.TrimSpace(name)
		eng.DeleteVariable(name)
		fmt.Printf("Deleted: %s\n", name)
		return true
	}

	return false
}

// handleSet handles "set" commands.
func handleSet(args string, eng *engine.Engine) {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		fmt.Println("Usage: set <option> <value>")
		fmt.Println("Options: precision, strict")
		return
	}

	option := strings.ToLower(parts[0])
	value := strings.TrimSpace(parts[1])

	switch option {
	case "precision":
		var p int
		_, err := fmt.Sscanf(value, "%d", &p)
		if err != nil || p < 0 || p > 15 {
			fmt.Println("Precision must be 0-15")
			return
		}
		eng.SetPrecision(p)
		fmt.Printf("Precision set to %d\n", p)

	case "strict":
		switch strings.ToLower(value) {
		case "on", "true", "1":
			eng.SetStrict(true)
			fmt.Println("Strict mode enabled")
		case "off", "false", "0":
			eng.SetStrict(false)
			fmt.Println("Strict mode disabled")
		default:
			fmt.Println("Usage: set strict on|off")
		}

	default:
		fmt.Printf("Unknown option: %s\n", option)
	}
}

// printResult prints a value result.
func printResult(result types.Value) {
	if result.IsEmpty() {
		return
	}

	if result.IsError() {
		fmt.Fprintf(os.Stderr, "Error: %s\n", result.ErrorMessage())
		return
	}

	fmt.Printf("= %s\n", result.String())
}

// printVariables prints all variables.
func printVariables(eng *engine.Engine) {
	vars := eng.Variables()
	if len(vars) == 0 {
		fmt.Println("No variables defined.")
		return
	}

	fmt.Println("Variables:")
	for name, value := range vars {
		fmt.Printf("  %s = %s\n", name, value.String())
	}
}

// printGroupedTotals prints totals grouped by type.
func printGroupedTotals(eng *engine.Engine) {
	totals := eng.GroupedTotals()
	if len(totals) == 0 {
		fmt.Println("No totals.")
		return
	}

	fmt.Println("Totals:")
	for _, t := range totals {
		fmt.Printf("  %s\n", t.String())
	}
}

// printHistory prints line history.
func printHistory(eng *engine.Engine) {
	lines := eng.Lines()
	if len(lines) == 0 {
		fmt.Println("No history.")
		return
	}

	fmt.Println("History:")
	for i, lr := range lines {
		status := ""
		if lr.IsConsumed {
			status = " (consumed)"
		}
		if lr.IsContinuation {
			status = " (continuation)"
		}
		if lr.AssignedVar != "" {
			status = fmt.Sprintf(" -> %s", lr.AssignedVar)
		}

		if !lr.Value.IsEmpty() {
			fmt.Printf("  %d: %s = %s%s\n", i+1, lr.Input, lr.Value.String(), status)
		}
	}
}

// printRateInfo prints rate cache information.
func printRateInfo(eng *engine.Engine) {
	rc := eng.RateCache()
	stats := rc.Stats()

	fmt.Println("Rate Cache:")
	fmt.Printf("  Direct rates: %d\n", stats.DirectRates)
	fmt.Printf("  Cache file: %s\n", stats.CacheFile)
	fmt.Printf("  Has file cache: %v\n", stats.HasFileCache)
	fmt.Printf("  Is expired: %v\n", stats.IsExpired)

	if !stats.LastUpdate.IsZero() {
		fmt.Printf("  Last update: %s\n", stats.LastUpdate.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Age: %s\n", stats.Age.Round(1000000000).String())
	}
}

// printBanner prints the welcome banner.
func printBanner() {
	fmt.Printf(`
  ┌─────────────────────────────────────┐
  │           numio v%s              │
  │     Natural Math Calculator         │
  │                                     │
  │  Type 'help' for commands           │
  │  Type 'quit' to exit                │
  └─────────────────────────────────────┘

`, appVersion)
}

// printHelp prints command line help.
func printHelp() {
	fmt.Printf(`%s v%s - Natural Math Calculator

Usage:
  %s                    Start interactive REPL
  %s <expression>       Evaluate expression
  %s -e <expression>    Evaluate expression
  %s -f <file>          Evaluate file

Options:
  -h, --help      Show this help
  -v, --version   Show version
  -e, --eval      Evaluate expression
  -f, --file      Evaluate file

Examples:
  %s "100 + 50"
  %s "$100 in EUR"
  %s "20%% of 150"
  %s -f calculations.txt

`, appName, appVersion, appName, appName, appName, appName, appName, appName, appName, appName)
}

// printREPLHelp prints REPL help.
func printREPLHelp() {
	fmt.Println(`
Commands:
  help, ?          Show this help
  quit, exit, q    Exit the program
  clear, cls       Clear all state
  vars             Show all variables
  total            Show running total
  totals           Show grouped totals
  history          Show line history
  rates            Show rate cache info
  set <opt> <val>  Set option (precision, strict)
  del <name>       Delete a variable

Expressions:
  100 + 50                 Basic math
  20% of 150               Percentage
  $100 + 15%               Price with tax
  $100 in EUR              Currency conversion
  5 km to miles            Unit conversion
  tax = 15%                Variable assignment
  _ * 2                    Use previous result
  sum(1, 2, 3)             Functions

Supported:
  Currencies: USD, EUR, GBP, JPY, TRY, BTC, ETH, ...
  Units: km, miles, kg, lb, C, F, hours, ...
  Functions: sum, avg, min, max, sqrt, round, ...
`)
}
