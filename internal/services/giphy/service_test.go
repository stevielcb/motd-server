package giphy

import (
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name       string
		apiKeyFile string
		apiKey     string
		expectErr  bool
	}{
		{
			name:       "valid API key file",
			apiKeyFile: "test-api-key.txt",
			apiKey:     "test-api-key-123",
			expectErr:  false,
		},
		{
			name:       "non-existent API key file",
			apiKeyFile: "non-existent-file.txt",
			apiKey:     "",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary API key file if needed
			if tt.apiKey != "" {
				if err := os.WriteFile(tt.apiKeyFile, []byte(tt.apiKey), 0644); err != nil {
					t.Fatalf("failed to create test API key file: %v", err)
				}
				defer os.Remove(tt.apiKeyFile)
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service, err := NewService(tt.apiKeyFile, logger)

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && service == nil {
				t.Error("expected service but got nil")
			}
			if !tt.expectErr && service.apiKey != tt.apiKey {
				t.Errorf("expected API key %s, got %s", tt.apiKey, service.apiKey)
			}
		})
	}
}

func TestService_GetRandom(t *testing.T) {
	// Create a temporary API key file
	apiKey := "test-api-key-123"
	apiKeyFile := "test-api-key.txt"
	if err := os.WriteFile(apiKeyFile, []byte(apiKey), 0644); err != nil {
		t.Fatalf("failed to create test API key file: %v", err)
	}
	defer os.Remove(apiKeyFile)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service, err := NewService(apiKeyFile, logger)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	tests := []struct {
		name      string
		tag       string
		rating    string
		expectErr bool
	}{
		{
			name:      "valid tag and rating",
			tag:       "funny",
			rating:    "g",
			expectErr: true, // Will fail because we don't have a real API key
		},
		{
			name:      "empty tag",
			tag:       "",
			rating:    "g",
			expectErr: true,
		},
		{
			name:      "empty rating",
			tag:       "funny",
			rating:    "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.GetRandom(tt.tag, tt.rating)

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && url == "" {
				t.Error("expected URL but got empty string")
			}
		})
	}
}

func TestService_GetRandom_InvalidAPIKey(t *testing.T) {
	// Create a temporary API key file with invalid key
	apiKey := "invalid-api-key"
	apiKeyFile := "test-invalid-api-key.txt"
	if err := os.WriteFile(apiKeyFile, []byte(apiKey), 0644); err != nil {
		t.Fatalf("failed to create test API key file: %v", err)
	}
	defer os.Remove(apiKeyFile)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service, err := NewService(apiKeyFile, logger)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	// This should fail with an invalid API key
	_, err = service.GetRandom("funny", "g")
	if err == nil {
		t.Error("expected error with invalid API key but got none")
	}
}
