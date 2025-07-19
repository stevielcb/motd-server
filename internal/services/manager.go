package services

import (
	"log/slog"

	"github.com/stevielcb/motd-server/internal/config"
	"github.com/stevielcb/motd-server/internal/services/giphy"
	"github.com/stevielcb/motd-server/internal/services/xkcd"
)

// Manager coordinates all external service calls
type Manager struct {
	config *config.Config
	giphy  GiphyProvider
	xkcd   XKCDProvider
	logger *slog.Logger
}

// NewManager creates a new services manager
func NewManager(cfg *config.Config, logger *slog.Logger) (*Manager, error) {
	giphyService, err := giphy.NewService(cfg.GiphyApiKeyFile, cfg.MaxFileSize, logger)
	if err != nil {
		return nil, err
	}

	xkcdService := xkcd.NewService(logger)

	return &Manager{
		config: cfg,
		giphy:  giphyService,
		xkcd:   xkcdService,
		logger: logger,
	}, nil
}

// DownloadMOTDs fetches new MOTDs from all configured services
func (m *Manager) DownloadMOTDs(cache CacheManager) error {
	// Download from Giphy for each configured tag
	for tag, rating := range m.config.GiphyTags {
		url, err := m.giphy.GetRandom(tag, rating)
		if err != nil {
			m.logger.Error("failed to fetch giphy", "tag", tag, "rating", rating, "error", err)
			continue
		}

		if err := cache.WriteToCache(url, ""); err != nil {
			m.logger.Error("failed to cache giphy", "url", url, "error", err)
		}
	}

	// Download from XKCD
	comic, err := m.xkcd.GetRandom()
	if err != nil {
		m.logger.Error("failed to fetch xkcd comic", "error", err)
		return err
	}

	if err := cache.WriteToCache(comic.ImageURL, comic.Alt); err != nil {
		m.logger.Error("failed to cache xkcd", "url", comic.ImageURL, "error", err)
	}

	return nil
}
