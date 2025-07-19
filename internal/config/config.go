package config

import (
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
)

// Config defines all configuration options for motd-server,
// populated from environment variables using envconfig.
type Config struct {
	CacheDir         string            `split_words:"true"`
	CacheMaxFiles    int               `split_words:"true" default:"50"`
	MaxFileSize      int64             `split_words:"true" default:"10485760"` // 10MB in bytes
	GiphyApiKeyFile  string            `split_words:"true"`
	GiphyTags        map[string]string `split_words:"true"`
	DownloadInterval int               `split_words:"true" default:"10"`
	CleanupInterval  int               `split_words:"true" default:"60"`
	ListenHost       string            `split_words:"true" default:"localhost"`
	ListenPort       int               `split_words:"true" default:"4200"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config

	err := envconfig.Process("motd", &cfg)
	if err != nil {
		return nil, err
	}

	// Set default values for paths
	home := os.Getenv("HOME")

	if cfg.GiphyApiKeyFile == "" {
		cfg.GiphyApiKeyFile = home + "/.giphy-api"
	}

	if cfg.CacheDir == "" {
		cfg.CacheDir = home + "/.motd"
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cfg.CacheDir, 0700); err != nil {
		return nil, err
	}

	slog.Info("configuration loaded",
		"cacheDir", cfg.CacheDir,
		"giphyKeyFile", cfg.GiphyApiKeyFile,
		"listenHost", cfg.ListenHost,
		"listenPort", cfg.ListenPort,
		"downloadInterval", cfg.DownloadInterval,
		"cleanupInterval", cfg.CleanupInterval,
		"cacheMaxFiles", cfg.CacheMaxFiles,
		"maxFileSize", cfg.MaxFileSize,
	)

	return &cfg, nil
}
