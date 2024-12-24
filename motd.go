package main

import (
	"fmt"
	"os"
	"sort"
)

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
