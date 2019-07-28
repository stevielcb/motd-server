package main

import (
	"fmt"
	"io/ioutil"
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
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return
	}

	if len(files) < c.CacheMaxFiles {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	// Delete oldest files while keeping the newest of the
	// defined max cached files
	for file := range files[:len(files)-c.CacheMaxFiles] {
		os.Remove(fmt.Sprintf("%s/%s", cacheDir, files[file].Name()))
	}
}
