package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
)

// getMotds fetches a new set of MOTDs from external sources (Giphy and XKCD)
// and writes them into the local cache directory.
//
// For each tag defined in the configuration (c.GiphyTags), a random Giphy is selected.
// Then, a random XKCD comic is also fetched.
//
// Files are saved locally for future serving to clients.
func getMotds() {
	for tag, rating := range c.GiphyTags {
		g, err := randomGiphy(tag, rating)
		if err != nil {
			slog.Error("failed to fetch giphy", "tag", tag, "rating", rating, "error", err)
			continue
		}
		writeToCache(g, "")
	}

	comic, err := randomXkcd()
	if err != nil {
		slog.Error("failed to fetch xkcd comic", "error", err)
		return
	}
	writeToCache(comic.ImageURL, comic.Alt)
}

// cleanupMotds ensures the cache directory does not exceed the maximum allowed
// number of files (c.CacheMaxFiles).
//
// It sorts the cached files by modification time (oldest first) and deletes
// the oldest files to maintain the configured cache size.
func cleanupMotds() {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		slog.Error("failed to read cache directory", "error", err)
		return
	}

	if len(entries) < c.CacheMaxFiles {
		return
	}

	// Sort entries by modification time (oldest first) using modern slices.SortFunc
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

	// Delete oldest files while keeping the newest of the
	// defined max cached files
	for _, entry := range entries[:len(entries)-c.CacheMaxFiles] {
		filePath := fmt.Sprintf("%s/%s", cacheDir, entry.Name())
		if err := os.Remove(filePath); err != nil {
			slog.Error("failed to remove old cache file", "file", filePath, "error", err)
		}
	}
}
