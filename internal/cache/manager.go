package cache

import (
	"bytes"
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
)

const (
	// CacheFilePrefix is the prefix used for cached file content format
	CacheFilePrefix = "1337"
	// CacheFileFormat is the format string for cached file content
	CacheFileFormat = "%s;File=inline=1;size=%d;name=%s:%s"
	// CacheFileFormatWithMessage is the format string for cached file content with message
	CacheFileFormatWithMessage = "%s;File=inline=1;size=%d;name=%s:%s%s\n"
)

// Manager handles all cache-related operations
type Manager struct {
	cacheDir    string
	maxFiles    int
	maxFileSize int64
	logger      *slog.Logger
}

// NewManager creates a new cache manager
func NewManager(cacheDir string, maxFiles int, maxFileSize int64, logger *slog.Logger) (*Manager, error) {
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Manager{
		cacheDir:    cacheDir,
		maxFiles:    maxFiles,
		maxFileSize: maxFileSize,
		logger:      logger,
	}, nil
}

// WriteToCache downloads content from the specified URL and saves it into the local cache directory
func (m *Manager) WriteToCache(url string, msg string) error {
	m.logger.Info("caching content", "url", url, "message", msg)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download content: %w", err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	b64url := b64.StdEncoding.EncodeToString([]byte(url))
	cacheFile := filepath.Join(m.cacheDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), b64url))

	f, err := os.Create(cacheFile)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer f.Close()

	encoded := b64.StdEncoding.EncodeToString(buf.Bytes())
	content := m.formatCacheContent(buf.Len(), b64url, encoded, msg)

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to cache file: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync cache file: %w", err)
	}

	m.logger.Debug("successfully cached content", "file", cacheFile, "size", buf.Len())
	return nil
}

// formatCacheContent formats the cache content with optional message
func (m *Manager) formatCacheContent(size int, b64url, encoded, msg string) string {
	if msg != "" {
		return fmt.Sprintf(CacheFileFormatWithMessage, CacheFilePrefix, size, b64url, encoded, msg)
	}
	return fmt.Sprintf(CacheFileFormat, CacheFilePrefix, size, b64url, encoded)
}

// GetRandomFile returns a random file from the cache directory
func (m *Manager) GetRandomFile() ([]byte, error) {
	var files []string

	err := filepath.Walk(m.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk cache directory: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no cached files found")
	}

	// Select random file using cryptographically secure random number generation
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(files))))
	if err != nil {
		return nil, fmt.Errorf("failed to generate random index: %w", err)
	}
	randFile := files[randIndex.Int64()]
	dat, err := os.ReadFile(randFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cached file: %w", err)
	}

	return dat, nil
}

// Cleanup ensures the cache directory does not exceed the maximum allowed number of files
func (m *Manager) Cleanup() error {
	entries, err := os.ReadDir(m.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	if len(entries) < m.maxFiles {
		return nil
	}

	// Sort entries by modification time (oldest first)
	slices.SortFunc(entries, func(a, b os.DirEntry) int {
		infoA, errA := a.Info()
		infoB, errB := b.Info()
		if errA != nil || errB != nil {
			return 0
		}
		if infoA.ModTime().Before(infoB.ModTime()) {
			return -1
		}
		if infoA.ModTime().After(infoB.ModTime()) {
			return 1
		}
		return 0
	})

	// Delete oldest files while keeping the newest of the defined max cached files
	for _, entry := range entries[:len(entries)-m.maxFiles] {
		filePath := filepath.Join(m.cacheDir, entry.Name())
		if err := os.Remove(filePath); err != nil {
			m.logger.Error("failed to remove old cache file", "file", filePath, "error", err)
		}
	}

	return nil
}
