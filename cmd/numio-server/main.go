// cmd/numio-server/main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/0xsj/numio/internal/socket"
	"github.com/0xsj/numio/pkg/engine"
)

func main() {
	// Parse flags
	socketPath := flag.String("socket", socket.DefaultSocketPath, "Unix socket path")
	refreshRates := flag.Bool("refresh", true, "Refresh exchange rates on startup")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *verbose {
		fmt.Println("numio-server starting...")
		fmt.Printf("  Socket: %s\n", *socketPath)
	}

	// Create engine
	eng := engine.New()

	// Load cached rates
	if eng.LoadRatesFromFile() {
		if *verbose {
			stats := eng.RateCacheStats()
			fmt.Printf("  Loaded %d cached rates\n", stats.DirectRates)
		}
	}

	// Refresh rates if requested
	if *refreshRates {
		if *verbose {
			fmt.Println("  Refreshing exchange rates...")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		n, err := eng.RefreshRatesIfExpired(ctx)
		cancel()

		if err != nil {
			if *verbose {
				fmt.Printf("  Warning: failed to refresh rates: %v\n", err)
			}
		} else if n > 0 {
			if *verbose {
				fmt.Printf("  Fetched %d rates\n", n)
			}
			// Save to cache
			if err := eng.SaveRatesToFile(); err != nil && *verbose {
				fmt.Printf("  Warning: failed to save rates: %v\n", err)
			}
		} else if *verbose {
			fmt.Println("  Rates cache is valid, no refresh needed")
		}
	}

	// Create server
	server := socket.NewServerWithPath(eng, *socketPath)

	if *verbose {
		fmt.Printf("  Listening on %s\n", *socketPath)
		fmt.Println("  Press Ctrl+C to stop")
		fmt.Println()
	}

	// Run server with signal handling
	if err := server.ListenForSignals(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Server stopped")
	}
}
