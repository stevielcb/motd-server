package cache

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name      string
		cacheDir  string
		maxFiles  int
		expectErr bool
	}{
		{
			name:      "valid configuration",
			cacheDir:  t.TempDir(),
			maxFiles:  50,
			expectErr: false,
		},
		{
			name:      "invalid cache directory",
			cacheDir:  "/nonexistent/path/that/should/fail",
			maxFiles:  50,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			manager, err := NewManager(tt.cacheDir, tt.maxFiles, 10*1024*1024, logger) // 10MB default

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && manager == nil {
				t.Error("expected manager but got nil")
			}
		})
	}
}

func TestManager_WriteToCache(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	manager, err := NewManager(tempDir, 50, 10*1024*1024, logger) // 10MB default
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	tests := []struct {
		name      string
		url       string
		msg       string
		expectErr bool
	}{
		{
			name:      "valid URL with message",
			url:       "https://httpbin.org/image/png",
			msg:       "test message",
			expectErr: false,
		},
		{
			name:      "valid URL without message",
			url:       "https://httpbin.org/image/png",
			msg:       "",
			expectErr: false,
		},
		{
			name:      "invalid URL",
			url:       "not-a-valid-url",
			msg:       "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.WriteToCache(tt.url, tt.msg)

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check if file was created (for successful cases)
			if !tt.expectErr {
				entries, err := os.ReadDir(tempDir)
				if err != nil {
					t.Fatalf("failed to read cache directory: %v", err)
				}
				if len(entries) == 0 {
					t.Error("expected cache file to be created but none found")
				}
			}
		})
	}
}

func TestManager_GetRandomFile(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	manager, err := NewManager(tempDir, 50, 10*1024*1024, logger) // 10MB default
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	tests := []struct {
		name        string
		setupFiles  bool
		expectErr   bool
		expectEmpty bool
	}{
		{
			name:        "empty cache directory",
			setupFiles:  false,
			expectErr:   true,
			expectEmpty: true,
		},
		{
			name:        "cache directory with files",
			setupFiles:  true,
			expectErr:   false,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up directory for each test
			entries, err := os.ReadDir(tempDir)
			if err != nil {
				t.Fatalf("failed to read cache directory: %v", err)
			}
			for _, entry := range entries {
				os.Remove(filepath.Join(tempDir, entry.Name()))
			}

			if tt.setupFiles {
				// Create some test files
				for i := 0; i < 3; i++ {
					testFile := filepath.Join(tempDir, fmt.Sprintf("test_%d.txt", i))
					if err := os.WriteFile(testFile, []byte(fmt.Sprintf("test content %d", i)), 0644); err != nil {
						t.Fatalf("failed to create test file: %v", err)
					}
				}
			}

			data, err := manager.GetRandomFile()

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectEmpty && len(data) > 0 {
				t.Error("expected empty data but got some")
			}
			if !tt.expectEmpty && len(data) == 0 {
				t.Error("expected data but got empty")
			}
		})
	}
}

func TestManager_Cleanup(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	manager, err := NewManager(tempDir, 3, 10*1024*1024, logger) // Set max files to 3, 10MB default
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	tests := []struct {
		name              string
		numFiles          int
		expectCleanup     bool
		expectedRemaining int
	}{
		{
			name:              "no cleanup needed",
			numFiles:          2,
			expectCleanup:     false,
			expectedRemaining: 2,
		},
		{
			name:              "cleanup needed",
			numFiles:          5,
			expectCleanup:     true,
			expectedRemaining: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up directory for each test
			entries, err := os.ReadDir(tempDir)
			if err != nil {
				t.Fatalf("failed to read cache directory: %v", err)
			}
			for _, entry := range entries {
				os.Remove(filepath.Join(tempDir, entry.Name()))
			}

			// Create test files with different timestamps
			for i := 0; i < tt.numFiles; i++ {
				testFile := filepath.Join(tempDir, fmt.Sprintf("test_%d.txt", i))
				if err := os.WriteFile(testFile, []byte(fmt.Sprintf("test content %d", i)), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}

				// Set different modification times to ensure proper ordering
				modTime := time.Now().Add(time.Duration(-i) * time.Hour)
				if err := os.Chtimes(testFile, modTime, modTime); err != nil {
					t.Fatalf("failed to set file time: %v", err)
				}
			}

			// Run cleanup
			if err := manager.Cleanup(); err != nil {
				t.Fatalf("cleanup failed: %v", err)
			}

			// Check remaining files
			entries, err = os.ReadDir(tempDir)
			if err != nil {
				t.Fatalf("failed to read cache directory after cleanup: %v", err)
			}

			if len(entries) != tt.expectedRemaining {
				t.Errorf("expected %d files after cleanup, got %d", tt.expectedRemaining, len(entries))
			}
		})
	}
}

func TestManager_CacheFileFormat(t *testing.T) {
	tempDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	manager, err := NewManager(tempDir, 50, 10*1024*1024, logger) // 10MB default
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	// Test that cache files are written in the correct format
	err = manager.WriteToCache("https://httpbin.org/image/png", "test message")
	if err != nil {
		t.Fatalf("failed to write to cache: %v", err)
	}

	// Find the cache file
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read cache directory: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no cache file created")
	}

	cacheFile := filepath.Join(tempDir, entries[0].Name())
	content, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("failed to read cache file: %v", err)
	}

	// Check that the content starts with the expected prefix
	if !bytes.HasPrefix(content, []byte(CacheFilePrefix)) {
		t.Errorf("cache file content does not start with expected prefix %s", CacheFilePrefix)
	}

	// Check that the message is included
	if !bytes.Contains(content, []byte("test message")) {
		t.Error("cache file content does not contain the expected message")
	}
}
