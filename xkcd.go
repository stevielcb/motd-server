// xkcd.go
//
// This file handles fetching random XKCD comics to be used as MOTDs (messages of the day).
// It uses the github.com/nishanths/go-xkcd/v2 client for interacting with the XKCD API.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/nishanths/go-xkcd/v2"
)

var (
	xkcdClient *xkcd.Client
)

// init initializes the global XKCD client used to fetch comics.
func init() {
	xkcdClient = xkcd.NewClient()
}

// randomXkcd fetches a random XKCD comic.
//
// It first retrieves the latest comic to determine the highest comic number,
// then randomly selects a number within the valid range and fetches that comic.
// Returns the comic and any error encountered.
func randomXkcd() (xkcd.Comic, error) {
	latest, err := xkcdClient.Latest(context.Background())
	if err != nil {
		return xkcd.Comic{}, fmt.Errorf("failed to fetch latest xkcd comic: %w", err)
	}

	number := rand.Intn(latest.Number-1) + 1
	comic, err := xkcdClient.Get(context.Background(), number)
	if err != nil {
		return comic, fmt.Errorf("failed to fetch xkcd comic %d: %w", number, err)
	}

	slog.Debug("fetched xkcd comic", "number", comic.Number, "title", comic.Title)
	return comic, nil
}
