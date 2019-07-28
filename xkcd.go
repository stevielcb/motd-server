package main

import (
	"github.com/nishanths/go-xkcd"
)

var (
	xkcdClient *xkcd.Client
)

func init() {
	xkcdClient = xkcd.NewClient()
}

func randomXkcd() (xkcd.Comic, error) {
	comic, err := xkcdClient.Random()
	if err != nil {
		return comic, err
	}
	return comic, nil
}
