package xkcd

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/nishanths/go-xkcd/v2"
)

// Service handles XKCD API interactions
type Service struct {
	client *xkcd.Client
	logger *slog.Logger
}

// NewService creates a new XKCD service
func NewService(logger *slog.Logger) *Service {
	return &Service{
		client: xkcd.NewClient(),
		logger: logger,
	}
}

// GetRandom fetches a random XKCD comic
func (s *Service) GetRandom() (xkcd.Comic, error) {
	latest, err := s.client.Latest(context.Background())
	if err != nil {
		return xkcd.Comic{}, fmt.Errorf("failed to fetch latest xkcd comic: %w", err)
	}

	number := rand.Intn(latest.Number-1) + 1
	comic, err := s.client.Get(context.Background(), number)
	if err != nil {
		return comic, fmt.Errorf("failed to fetch xkcd comic %d: %w", number, err)
	}

	s.logger.Debug("fetched xkcd comic", "number", comic.Number, "title", comic.Title)
	return comic, nil
}
