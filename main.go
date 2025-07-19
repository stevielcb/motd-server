package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/stevielcb/motd-server/app"
	"github.com/stevielcb/motd-server/internal/config"
)

// version will be set during build time via ldflags
var version = "dev"

func main() {
	// Parse command line flags
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("motd-server version %s\n", version)
		os.Exit(0)
	}

	// Initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create application
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("failed to create application", "error", err)
		os.Exit(1)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("received shutdown signal")
		if err := application.Stop(); err != nil {
			logger.Error("failed to stop application gracefully", "error", err)
		}
		os.Exit(0)
	}()

	// Start application
	if err := application.Start(); err != nil {
		logger.Error("failed to start application", "error", err)
		os.Exit(1)
	}
}
