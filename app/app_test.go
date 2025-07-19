package app

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stevielcb/motd-server/internal/config"
)

func TestNew(t *testing.T) {
	// Create a temporary directory for cache
	tempDir := t.TempDir()

	// Create a temporary API key file
	apiKeyFile := tempDir + "/giphy-api"
	if err := os.WriteFile(apiKeyFile, []byte("test-api-key"), 0644); err != nil {
		t.Fatalf("failed to create test API key file: %v", err)
	}

	cfg := &config.Config{
		CacheDir:         tempDir,
		CacheMaxFiles:    50,
		GiphyApiKeyFile:  apiKeyFile,
		GiphyTags:        map[string]string{"funny": "g"},
		DownloadInterval: 10,
		CleanupInterval:  60,
		ListenHost:       "localhost",
		ListenPort:       4200,
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app, err := New(cfg, logger)

	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	if app == nil {
		t.Fatal("expected app but got nil")
	}

	if app.config != cfg {
		t.Error("config not properly set")
	}

	if app.cache == nil {
		t.Error("cache manager not properly initialized")
	}

	if app.server == nil {
		t.Error("server not properly initialized")
	}

	if app.services == nil {
		t.Error("services manager not properly initialized")
	}

	if app.logger != logger {
		t.Error("logger not properly set")
	}
}

func TestNew_InvalidCacheDir(t *testing.T) {
	// Create a temporary API key file
	tempDir := t.TempDir()
	apiKeyFile := tempDir + "/giphy-api"
	if err := os.WriteFile(apiKeyFile, []byte("test-api-key"), 0644); err != nil {
		t.Fatalf("failed to create test API key file: %v", err)
	}

	cfg := &config.Config{
		CacheDir:         "/nonexistent/path/that/should/fail",
		CacheMaxFiles:    50,
		GiphyApiKeyFile:  apiKeyFile,
		GiphyTags:        map[string]string{"funny": "g"},
		DownloadInterval: 10,
		CleanupInterval:  60,
		ListenHost:       "localhost",
		ListenPort:       4200,
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app, err := New(cfg, logger)

	if err == nil {
		t.Error("expected error with invalid cache directory but got none")
	}

	if app != nil {
		t.Error("expected nil app but got one")
	}
}

func TestNew_InvalidGiphyKeyFile(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		CacheDir:         tempDir,
		CacheMaxFiles:    50,
		GiphyApiKeyFile:  "/nonexistent/giphy-api",
		GiphyTags:        map[string]string{"funny": "g"},
		DownloadInterval: 10,
		CleanupInterval:  60,
		ListenHost:       "localhost",
		ListenPort:       4200,
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app, err := New(cfg, logger)

	if err == nil {
		t.Error("expected error with invalid giphy API key file but got none")
	}

	if app != nil {
		t.Error("expected nil app but got one")
	}
}

func TestApp_StartAndStop(t *testing.T) {
	// Create a temporary directory for cache
	tempDir := t.TempDir()

	// Create a temporary API key file
	apiKeyFile := tempDir + "/giphy-api"
	if err := os.WriteFile(apiKeyFile, []byte("test-api-key"), 0644); err != nil {
		t.Fatalf("failed to create test API key file: %v", err)
	}

	cfg := &config.Config{
		CacheDir:         tempDir,
		CacheMaxFiles:    50,
		GiphyApiKeyFile:  apiKeyFile,
		GiphyTags:        map[string]string{"funny": "g"},
		DownloadInterval: 10,
		CleanupInterval:  60,
		ListenHost:       "localhost",
		ListenPort:       0, // Use port 0 to avoid binding issues
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	// Start the app in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Start()
	}()

	// Give the app a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop the app
	if err := app.Stop(); err != nil {
		t.Errorf("failed to stop app: %v", err)
	}

	// Check for app errors - ignore connection errors since we're stopping the server
	select {
	case err := <-errChan:
		if err != nil && !isConnectionError(err) {
			t.Errorf("app returned unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		// App should have stopped by now
	}
}

// isConnectionError checks if the error is related to connection closure
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "failed to accept connection") &&
		strings.Contains(errMsg, "use of closed network connection")
}

func TestApp_BackgroundWorkers(t *testing.T) {
	// Create a temporary directory for cache
	tempDir := t.TempDir()

	// Create a temporary API key file
	apiKeyFile := tempDir + "/giphy-api"
	if err := os.WriteFile(apiKeyFile, []byte("test-api-key"), 0644); err != nil {
		t.Fatalf("failed to create test API key file: %v", err)
	}

	cfg := &config.Config{
		CacheDir:         tempDir,
		CacheMaxFiles:    50,
		GiphyApiKeyFile:  apiKeyFile,
		GiphyTags:        map[string]string{"funny": "g"},
		DownloadInterval: 1, // Very short interval for testing
		CleanupInterval:  1, // Very short interval for testing
		ListenHost:       "localhost",
		ListenPort:       0, // Use port 0 to avoid binding issues
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	app, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	// Start background workers
	app.startBackgroundWorkers()

	// Give workers a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop the app
	if err := app.Stop(); err != nil {
		t.Errorf("failed to stop app: %v", err)
	}

	// Wait for workers to stop
	app.wg.Wait()
}
