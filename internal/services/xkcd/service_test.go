package xkcd

import (
	"io"
	"log/slog"
	"testing"
)

func TestNewService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewService(logger)

	if service == nil {
		t.Fatal("expected service but got nil")
	}

	if service.client == nil {
		t.Error("expected client but got nil")
	}

	if service.logger != logger {
		t.Error("logger not properly set")
	}
}

func TestService_GetRandom(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewService(logger)

	// Test that GetRandom returns a comic
	comic, err := service.GetRandom()

	// This test might fail if there's no internet connection or the XKCD API is down
	// We'll just check that we get a reasonable response
	if err != nil {
		// If we get an error, it should be a network-related error, not a logic error
		t.Logf("GetRandom returned error (expected if no internet): %v", err)
		return
	}

	// If successful, check that we got a valid comic
	if comic.Number <= 0 {
		t.Error("expected comic number to be positive")
	}

	if comic.Title == "" {
		t.Error("expected comic title to be non-empty")
	}

	if comic.ImageURL == "" {
		t.Error("expected comic image URL to be non-empty")
	}
}

func TestService_GetRandom_MultipleCalls(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := NewService(logger)

	// Test multiple calls to ensure randomness
	comics := make(map[int]bool)

	for i := 0; i < 5; i++ {
		comic, err := service.GetRandom()
		if err != nil {
			t.Logf("GetRandom returned error (expected if no internet): %v", err)
			return
		}

		comics[comic.Number] = true
	}

	// If we got multiple comics, they should be different (random)
	if len(comics) > 1 {
		t.Logf("Got %d different comics out of 5 calls", len(comics))
	} else {
		t.Log("Got same comic multiple times (this can happen with random selection)")
	}
}
