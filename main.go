package main

import (
	"time"
)

func main() {
	motdTicker := time.NewTicker(time.Duration(c.DownloadInterval) * time.Second)
	go func() {
		for range motdTicker.C {
			getMotds()
		}
	}()

	cleanupTicker := time.NewTicker(time.Duration(c.CleanupInterval) * time.Second)
	go func() {
		for range cleanupTicker.C {
			cleanupMotds()
		}
	}()

	startServer()
}
