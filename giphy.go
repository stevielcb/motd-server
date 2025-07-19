// giphy.go
//
// This file handles fetching random Giphy images to be used as MOTDs (messages of the day).
// It interacts with the Giphy API to fetch either an original or a downsized image based on size limits.
package main

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

const (
	// MaxFileSizeMB is the maximum file size in bytes (10MB)
	MaxFileSizeMB = 10 * 1024 * 1024
)

var (
	apiKey string
)

// init reads the Giphy API key from the file specified by giphyKeyFile
// and initializes the global apiKey variable.
func init() {
	dat, err := os.ReadFile(giphyKeyFile)
	if err != nil {
		slog.Error("failed to read giphy API key file", "file", giphyKeyFile, "error", err)
		panic(err)
	}

	apiKey = string(dat)
}

// randomGiphy fetches a random Giphy URL matching the given tag and rating.
//
// It selects an image and checks the size of the "original" version.
// If the original exceeds MaxFileSizeMB, it falls back to the downsized large version.
// Returns the URL of the selected Giphy image or an error if encountered.
func randomGiphy(tag string, rating string) (string, error) {
	url := fmt.Sprintf(
		"http://api.giphy.com/v1/gifs/random?api_key=%s&tag=%s&rating=%s",
		apiKey,
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

	size, _ := strconv.Atoi(sizeResp.Header.Get("Content-Length"))
	// If the original gif is larger than MaxFileSizeMB,
	// get the downsized image instead.
	downsizedURL, ok := downsized["url"].(string)
	if !ok || downsizedURL == "" {
		return "", fmt.Errorf("no downsized image URL in giphy API response")
	}

	if int64(size) > MaxFileSizeMB {
		return downsizedURL, nil
	}

	return originalURL, nil
}
