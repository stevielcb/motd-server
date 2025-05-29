package main

import (
	"fmt"
	"os"
	"sort"
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
		if err == nil {
			writeToCache(g, "")
		}
	}

	c, err := randomXkcd()
	if err == nil {
		writeToCache(c.ImageURL, c.Alt)
	}
}

// cleanupMotds ensures the cache directory does not exceed the maximum allowed
// number of files (c.CacheMaxFiles).
//
// It sorts the cached files by modification time (oldest first) and deletes
// the oldest files to maintain the configured cache size.
func cleanupMotds() {
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}

	if len(entries) < c.CacheMaxFiles {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		infoI, errI := entries[i].Info()
		infoJ, errJ := entries[j].Info()
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Delete oldest files while keeping the newest of the
	// defined max cached files
	for _, entry := range entries[:len(entries)-c.CacheMaxFiles] {
		os.Remove(fmt.Sprintf("%s/%s", cacheDir, entry.Name()))
	}
}
