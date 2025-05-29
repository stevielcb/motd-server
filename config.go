// config.go
//
// This file handles loading and initializing configuration for the motd-server application.
// It uses environment variables to populate a Config struct via the github.com/kelseyhightower/envconfig package.
package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

var (
	cacheDir     string
	giphyKeyFile string
	c            Config
)

// Config defines all configuration options for motd-server,
// populated from environment variables using envconfig.
type Config struct {
	CacheDir         string            `split_words:"true"`
	CacheMaxFiles    int               `split_words:"true" default:"50"`
	GiphyApiKeyFile  string            `split_words:"true"`
	GiphyTags        map[string]string `split_words:"true"`
	DownloadInterval int               `split_words:"true" default:"10"`
	CleanupInterval  int               `split_words:"true" default:"60"`
	ListenHost       string            `split_words:"true" default:"localhost"`
	ListenPort       int               `split_words:"true" default:"4200"`
}

// init loads environment variables into the Config struct,
// sets up paths for the Giphy API key file and cache directory,
// and ensures the cache directory exists.
func init() {
	err := envconfig.Process("motd", &c)
	if err != nil {
		panic(err)
	}

	home := os.Getenv("HOME")

	if c.GiphyApiKeyFile == "" {
		giphyKeyFile = home + "/.giphy-api"
	} else {
		giphyKeyFile = c.GiphyApiKeyFile
	}

	if c.CacheDir == "" {
		cacheDir = home + "/.motd"
	} else {
		cacheDir = c.CacheDir
	}

	os.Mkdir(cacheDir, 0700)
}
