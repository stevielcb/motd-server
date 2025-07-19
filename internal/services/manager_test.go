package services

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/nishanths/go-xkcd/v2"
	"github.com/stevielcb/motd-server/internal/config"
)

// Mock implementations for testing
type mockGiphyProvider struct {
	shouldError bool
}

func (m *mockGiphyProvider) GetRandom(tag string, rating string) (string, error) {
	if m.shouldError {
		return "", errors.New("mock giphy error")
	}
	return "https://example.com/giphy.gif", nil
}

type mockXKCDProvider struct {
	shouldError bool
}

func (m *mockXKCDProvider) GetRandom() (xkcd.Comic, error) {
	if m.shouldError {
		return xkcd.Comic{}, errors.New("mock xkcd error")
	}
	return xkcd.Comic{
		ImageURL: "https://example.com/xkcd.png",
		Alt:      "Mock XKCD comic",
	}, nil
}

type mockCacheManager struct {
	writeError bool
}

func (m *mockCacheManager) WriteToCache(url string, msg string) error {
	if m.writeError {
		return errors.New("mock cache write error")
	}
	return nil
}

func (m *mockCacheManager) GetRandomFile() ([]byte, error) {
	return []byte("mock content"), nil
}

func (m *mockCacheManager) Cleanup() error {
	return nil
}

func TestManager_DownloadMOTDs(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := &config.Config{
		GiphyTags: map[string]string{
			"funny": "g",
		},
	}

	tests := []struct {
		name          string
		giphyError    bool
		xkcdError     bool
		cacheError    bool
		expectedError bool
	}{
		{
			name:          "successful download",
			expectedError: false,
		},
		{
			name:          "giphy error",
			giphyError:    true,
			expectedError: false, // Giphy errors are logged but don't stop the process
		},
		{
			name:          "xkcd error",
			xkcdError:     true,
			expectedError: true,
		},
		{
			name:          "cache write error",
			cacheError:    true,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				config: cfg,
				giphy:  &mockGiphyProvider{shouldError: tt.giphyError},
				xkcd:   &mockXKCDProvider{shouldError: tt.xkcdError},
				logger: logger,
			}

			cache := &mockCacheManager{writeError: tt.cacheError}

			err := manager.DownloadMOTDs(cache)

			if tt.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
