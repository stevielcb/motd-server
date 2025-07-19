package giphy

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Service handles Giphy API interactions
type Service struct {
	apiKey      string
	maxFileSize int64
	logger      *slog.Logger
}

// NewService creates a new Giphy service
func NewService(apiKeyFile string, maxFileSize int64, logger *slog.Logger) (*Service, error) {
	dat, err := os.ReadFile(apiKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read giphy API key file: %w", err)
	}

	return &Service{
		apiKey:      string(dat),
		maxFileSize: maxFileSize,
		logger:      logger,
	}, nil
}

// GetRandom fetches a random Giphy URL matching the given tag and rating
func (s *Service) GetRandom(tag string, rating string) (string, error) {
	url := fmt.Sprintf(
		"http://api.giphy.com/v1/gifs/random?api_key=%s&tag=%s&rating=%s",
		s.apiKey,
		url.QueryEscape(tag),
		rating,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch giphy API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read giphy API response: %w", err)
	}

	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal giphy API response: %w", err)
	}

	data, ok := result["data"].(map[string]any)
	if !ok || data == nil {
		return "", fmt.Errorf("no data in giphy API response")
	}

	images, ok := data["images"].(map[string]any)
	if !ok || images == nil {
		return "", fmt.Errorf("no images data in giphy API response")
	}

	original, ok := images["original"].(map[string]any)
	if !ok || original == nil {
		return "", fmt.Errorf("no original image data in giphy API response")
	}

	downsized, ok := images["downsized_large"].(map[string]any)
	if !ok || downsized == nil {
		return "", fmt.Errorf("no downsized image data in giphy API response")
	}

	originalURL, ok := original["url"].(string)
	if !ok || originalURL == "" {
		return "", fmt.Errorf("no original image URL in giphy API response")
	}

	sizeResp, err := http.Head(originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to check original image size: %w", err)
	}

	if sizeResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to check image size, status: %d", sizeResp.StatusCode)
	}

	sizeStr := sizeResp.Header.Get("Content-Length")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse content length: %w", err)
	}

	// If the original gif is larger than MaxFileSizeMB, get the downsized image instead
	downsizedURL, ok := downsized["url"].(string)
	if !ok || downsizedURL == "" {
		return "", fmt.Errorf("no downsized image URL in giphy API response")
	}

	if int64(size) > s.maxFileSize {
		return downsizedURL, nil
	}

	return originalURL, nil
}
