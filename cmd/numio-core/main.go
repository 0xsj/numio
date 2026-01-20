package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/numio/internal/rpc"
)

// Version information (set by ldflags at build time)
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	// Parse command line flags
	var (
		showVersion = flag.Bool("version", false, "Print version and exit")
		showHelp    = flag.Bool("help", false, "Print help and exit")
	)
	flag.Parse()

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// Run the server
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create transport over stdin/stdout
	transport := rpc.NewStdioTransport(os.Stdin, os.Stdout)

	// Create server
	server := rpc.NewServer(transport,
		rpc.WithMiddleware(
			rpc.RecoveryMiddleware(),
			rpc.ValidationMiddleware(),
		),
		rpc.WithOnStart(func() {
			logDebug("numio-core started")
		}),
		rpc.WithOnShutdown(func() {
			logDebug("numio-core shutting down")
		}),
		rpc.WithErrorHandler(func(err error) {
			logDebug("error: %v", err)
		}),
	)

	// Create handler and register methods
	handler := NewCoreHandler(version)
	handler.SetServer(server)
	handler.RegisterHandlers(server.Registry())

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logDebug("received shutdown signal")
		cancel()
		server.Shutdown(context.Background())
	}()

	// Start server (blocks until shutdown)
	err := server.Start()

	// Ignore EOF (normal shutdown)
	if err != nil && err.Error() != "EOF" {
		return err
	}

	_ = ctx // Silence unused warning

	return nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `numio-core - Desktop backend for numio calculator

Usage:
  numio-core [flags]

Description:
  numio-core is the backend process for the numio desktop application.
  It communicates via JSON-RPC 2.0 over stdin/stdout.

  This binary is typically spawned by a native UI application (macOS, Windows)
  and should not be run directly by users.

Flags:
  --version    Print version information and exit
  --help       Print this help message and exit

Protocol:
  Communication uses newline-delimited JSON-RPC 2.0 messages.

  Methods (UI → Core):
    initialize    Start session, get initial state
    shutdown      Graceful shutdown
    keypress      Send key event
    resize        Update viewport size
    getState      Request current render state
    setOption     Set configuration option

  Notifications (Core → UI):
    render        Push updated render state
    rateStatus    Rate fetch status update
    log           Debug/info logging

Example:
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"viewport":{"width":80,"height":24}}}' | numio-core

`)
}

func printVersion() {
	fmt.Printf("numio-core %s\n", version)
	fmt.Printf("  commit:  %s\n", commit)
	fmt.Printf("  built:   %s\n", buildTime)
}

// logDebug writes debug messages to stderr (not stdout, which is for RPC).
func logDebug(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[numio-core] "+format+"\n", args...)
}
