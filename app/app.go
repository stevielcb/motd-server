package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/stevielcb/motd-server/internal/cache"
	"github.com/stevielcb/motd-server/internal/config"
	"github.com/stevielcb/motd-server/internal/server"
	"github.com/stevielcb/motd-server/internal/services"
)

// App represents the main application with all its dependencies
type App struct {
	config   *config.Config
	cache    *cache.Manager
	server   *server.TCPServer
	services *services.Manager
	logger   *slog.Logger

	// Background workers
	downloadTicker *time.Ticker
	cleanupTicker  *time.Ticker

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates a new application instance with all dependencies
func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize cache manager
	cacheManager, err := cache.NewManager(cfg.CacheDir, cfg.CacheMaxFiles, cfg.MaxFileSize, logger)
	if err != nil {
		cancel()
		return nil, err
	}

	// Initialize services manager
	servicesManager, err := services.NewManager(cfg, logger)
	if err != nil {
		cancel()
		return nil, err
	}

	// Initialize TCP server
	tcpServer := server.NewTCPServer(cfg.ListenHost, cfg.ListenPort, cacheManager, logger)

	app := &App{
		config:   cfg,
		cache:    cacheManager,
		server:   tcpServer,
		services: servicesManager,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}

	return app, nil
}

// Start begins all application services
func (a *App) Start() error {
	a.logger.Info("starting motd-server")

	// Start background workers
	a.startBackgroundWorkers()

	// Start TCP server
	return a.server.Start()
}

// Stop gracefully shuts down the application
func (a *App) Stop() error {
	a.logger.Info("stopping motd-server")

	// Cancel context to stop background workers
	a.cancel()

	// Stop tickers
	if a.downloadTicker != nil {
		a.downloadTicker.Stop()
	}
	if a.cleanupTicker != nil {
		a.cleanupTicker.Stop()
	}

	// Wait for all goroutines to finish
	a.wg.Wait()

	// Stop server
	return a.server.Stop()
}

// startBackgroundWorkers starts the download and cleanup goroutines
func (a *App) startBackgroundWorkers() {
	// Download worker
	a.downloadTicker = time.NewTicker(time.Duration(a.config.DownloadInterval) * time.Second)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.logger.Info("starting MOTD download worker", "interval", a.config.DownloadInterval)

		for {
			select {
			case <-a.downloadTicker.C:
				if err := a.services.DownloadMOTDs(a.cache); err != nil {
					a.logger.Error("failed to download MOTDs", "error", err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()

	// Cleanup worker
	a.cleanupTicker = time.NewTicker(time.Duration(a.config.CleanupInterval) * time.Second)
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.logger.Info("starting cleanup worker", "interval", a.config.CleanupInterval)

		for {
			select {
			case <-a.cleanupTicker.C:
				if err := a.cache.Cleanup(); err != nil {
					a.logger.Error("failed to cleanup cache", "error", err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()
}
