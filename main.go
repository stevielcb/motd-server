package main

import (
	"log/slog"
	"os"
	"time"
)

// main
//
// Entry point of the motd-server application.
//
// This function sets up two background goroutines:
// - One that periodically downloads or refreshes the MOTDs based on c.DownloadInterval.
// - Another that periodically cleans up old MOTD files based on c.CleanupInterval.
//
// After starting the background processes, it launches the TCP server by calling startServer().
//
// Goroutines:
// - getMotds(): Refreshes MOTD cache.
// - cleanupMotds(): Cleans up outdated MOTD files.
//
// Functions called:
// - startServer(): Starts listening for incoming TCP connections and serving messages.
func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting motd-server")

	motdTicker := time.NewTicker(time.Duration(c.DownloadInterval) * time.Second)
	go func() {
		slog.Info("starting MOTD download goroutine", "interval", c.DownloadInterval)
		for range motdTicker.C {
			getMotds()
		}
	}()

	cleanupTicker := time.NewTicker(time.Duration(c.CleanupInterval) * time.Second)
	go func() {
		slog.Info("starting cleanup goroutine", "interval", c.CleanupInterval)
		for range cleanupTicker.C {
			cleanupMotds()
		}
	}()

	startServer()
}
