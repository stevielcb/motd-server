package main

import (
	"context"
	"math/rand"

	"github.com/nishanths/go-xkcd/v2"
)

var (
	xkcdClient *xkcd.Client
)

func init() {
	xkcdClient = xkcd.NewClient()
}

func randomXkcd() (xkcd.Comic, error) {
	latest, err := xkcdClient.Latest(context.Background())
	if err != nil {
		return xkcd.Comic{}, err
	}
	number := rand.Intn(latest.Number-1) + 1
	comic, err := xkcdClient.Get(context.Background(), number)
	if err != nil {
		return comic, err
	}
	return comic, nil
}
