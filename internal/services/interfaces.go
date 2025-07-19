package services

import (
	"github.com/nishanths/go-xkcd/v2"
)

// GiphyProvider defines the interface for Giphy service
type GiphyProvider interface {
	GetRandom(tag string, rating string) (string, error)
}

// XKCDProvider defines the interface for XKCD service
type XKCDProvider interface {
	GetRandom() (xkcd.Comic, error)
}

// CacheManager defines the interface for cache operations
type CacheManager interface {
	WriteToCache(url string, msg string) error
	GetRandomFile() ([]byte, error)
	Cleanup() error
}
