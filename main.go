package main

import (
	"time"
)

func main() {
	motdTicker := time.NewTicker(time.Duration(c.DownloadInterval) * time.Second)
	go func() {
		for {
			select {
			case <-motdTicker.C:
				getMotds()
			}
		}
	}()

	cleanupTicker := time.NewTicker(time.Duration(c.CleanupInterval) * time.Second)
	go func() {
		for {
			select {
			case <-cleanupTicker.C:
				cleanupMotds()
			}
		}
	}()

	startServer()
}
