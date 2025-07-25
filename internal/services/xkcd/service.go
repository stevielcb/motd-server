package xkcd

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"

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

	// Generate a cryptographically secure random number from 1 to latest.Number (inclusive)
	// rand.Int generates [0, latest.Number), so we add 1 to get [1, latest.Number]
	randNum, err := rand.Int(rand.Reader, big.NewInt(int64(latest.Number)))
	if err != nil {
		return xkcd.Comic{}, fmt.Errorf("failed to generate random number: %w", err)
	}
	number := int(randNum.Int64()) + 1 // Convert from 0-based to 1-based comic numbering
	comic, err := s.client.Get(context.Background(), number)
	if err != nil {
		return comic, fmt.Errorf("failed to fetch xkcd comic %d: %w", number, err)
	}

	s.logger.Debug("fetched xkcd comic", "number", comic.Number, "title", comic.Title)
	return comic, nil
}
