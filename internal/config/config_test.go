package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment variables
	originalHome := os.Getenv("HOME")
	originalCacheDir := os.Getenv("MOTD_CACHE_DIR")
	originalGiphyKeyFile := os.Getenv("MOTD_GIPHY_API_KEY_FILE")
	originalMaxFiles := os.Getenv("MOTD_CACHE_MAX_FILES")
	originalDownloadInterval := os.Getenv("MOTD_DOWNLOAD_INTERVAL")
	originalCleanupInterval := os.Getenv("MOTD_CLEANUP_INTERVAL")
	originalListenHost := os.Getenv("MOTD_LISTEN_HOST")
	originalListenPort := os.Getenv("MOTD_LISTEN_PORT")

	// Restore environment variables after test
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
		if originalCacheDir != "" {
			os.Setenv("MOTD_CACHE_DIR", originalCacheDir)
		} else {
			os.Unsetenv("MOTD_CACHE_DIR")
		}
		if originalGiphyKeyFile != "" {
			os.Setenv("MOTD_GIPHY_API_KEY_FILE", originalGiphyKeyFile)
		} else {
			os.Unsetenv("MOTD_GIPHY_API_KEY_FILE")
		}
		if originalMaxFiles != "" {
			os.Setenv("MOTD_CACHE_MAX_FILES", originalMaxFiles)
		} else {
			os.Unsetenv("MOTD_CACHE_MAX_FILES")
		}
		if originalDownloadInterval != "" {
			os.Setenv("MOTD_DOWNLOAD_INTERVAL", originalDownloadInterval)
		} else {
			os.Unsetenv("MOTD_DOWNLOAD_INTERVAL")
		}
		if originalCleanupInterval != "" {
			os.Setenv("MOTD_CLEANUP_INTERVAL", originalCleanupInterval)
		} else {
			os.Unsetenv("MOTD_CLEANUP_INTERVAL")
		}
		if originalListenHost != "" {
			os.Setenv("MOTD_LISTEN_HOST", originalListenHost)
		} else {
			os.Unsetenv("MOTD_LISTEN_HOST")
		}
		if originalListenPort != "" {
			os.Setenv("MOTD_LISTEN_PORT", originalListenPort)
		} else {
			os.Unsetenv("MOTD_LISTEN_PORT")
		}
	}()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectErr      bool
		expectedConfig *Config
	}{
		{
			name: "default configuration",
			envVars: map[string]string{
				"HOME": "/tmp/test-home",
			},
			expectErr: false,
			expectedConfig: &Config{
				CacheDir:         "/tmp/test-home/.motd",
				CacheMaxFiles:    50,
				GiphyApiKeyFile:  "/tmp/test-home/.giphy-api",
				GiphyTags:        nil,
				DownloadInterval: 10,
				CleanupInterval:  60,
				ListenHost:       "localhost",
				ListenPort:       4200,
			},
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"HOME":                    "/tmp/test-home",
				"MOTD_CACHE_DIR":          "/tmp/test-cache",
				"MOTD_GIPHY_API_KEY_FILE": "/tmp/test-giphy.key",
				"MOTD_CACHE_MAX_FILES":    "100",
				"MOTD_DOWNLOAD_INTERVAL":  "30",
				"MOTD_CLEANUP_INTERVAL":   "120",
				"MOTD_LISTEN_HOST":        "0.0.0.0",
				"MOTD_LISTEN_PORT":        "8080",
			},
			expectErr: false,
			expectedConfig: &Config{
				CacheDir:         "/tmp/test-cache",
				CacheMaxFiles:    100,
				GiphyApiKeyFile:  "/tmp/test-giphy.key",
				GiphyTags:        nil,
				DownloadInterval: 30,
				CleanupInterval:  120,
				ListenHost:       "0.0.0.0",
				ListenPort:       8080,
			},
		},
		{
			name: "invalid cache directory",
			envVars: map[string]string{
				"HOME":           "/tmp/test-home",
				"MOTD_CACHE_DIR": "/nonexistent/path/that/should/fail",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables for this test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up environment variables not set in this test
			envVarsToUnset := []string{
				"MOTD_CACHE_DIR", "MOTD_GIPHY_API_KEY_FILE", "MOTD_CACHE_MAX_FILES",
				"MOTD_DOWNLOAD_INTERVAL", "MOTD_CLEANUP_INTERVAL", "MOTD_LISTEN_HOST", "MOTD_LISTEN_PORT",
			}
			for _, key := range envVarsToUnset {
				if _, exists := tt.envVars[key]; !exists {
					os.Unsetenv(key)
				}
			}

			cfg, err := Load()

			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && cfg == nil {
				t.Error("expected config but got nil")
			}

			if tt.expectedConfig != nil && cfg != nil {
				if cfg.CacheDir != tt.expectedConfig.CacheDir {
					t.Errorf("expected CacheDir %s, got %s", tt.expectedConfig.CacheDir, cfg.CacheDir)
				}
				if cfg.CacheMaxFiles != tt.expectedConfig.CacheMaxFiles {
					t.Errorf("expected CacheMaxFiles %d, got %d", tt.expectedConfig.CacheMaxFiles, cfg.CacheMaxFiles)
				}
				if cfg.GiphyApiKeyFile != tt.expectedConfig.GiphyApiKeyFile {
					t.Errorf("expected GiphyApiKeyFile %s, got %s", tt.expectedConfig.GiphyApiKeyFile, cfg.GiphyApiKeyFile)
				}
				if cfg.DownloadInterval != tt.expectedConfig.DownloadInterval {
					t.Errorf("expected DownloadInterval %d, got %d", tt.expectedConfig.DownloadInterval, cfg.DownloadInterval)
				}
				if cfg.CleanupInterval != tt.expectedConfig.CleanupInterval {
					t.Errorf("expected CleanupInterval %d, got %d", tt.expectedConfig.CleanupInterval, cfg.CleanupInterval)
				}
				if cfg.ListenHost != tt.expectedConfig.ListenHost {
					t.Errorf("expected ListenHost %s, got %s", tt.expectedConfig.ListenHost, cfg.ListenHost)
				}
				if cfg.ListenPort != tt.expectedConfig.ListenPort {
					t.Errorf("expected ListenPort %d, got %d", tt.expectedConfig.ListenPort, cfg.ListenPort)
				}
			}
		})
	}
}

func TestLoad_WithGiphyTags(t *testing.T) {
	// Save original environment variables
	originalHome := os.Getenv("HOME")
	originalGiphyTags := os.Getenv("MOTD_GIPHY_TAGS")

	// Restore environment variables after test
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
		if originalGiphyTags != "" {
			os.Setenv("MOTD_GIPHY_TAGS", originalGiphyTags)
		} else {
			os.Unsetenv("MOTD_GIPHY_TAGS")
		}
	}()

	// Set up test environment
	os.Setenv("HOME", "/tmp/test-home")
	os.Setenv("MOTD_GIPHY_TAGS", "funny:g,cats:pg")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	expectedTags := map[string]string{
		"funny": "g",
		"cats":  "pg",
	}

	if len(cfg.GiphyTags) != len(expectedTags) {
		t.Errorf("expected %d giphy tags, got %d", len(expectedTags), len(cfg.GiphyTags))
	}

	for expectedKey, expectedValue := range expectedTags {
		if actualValue, exists := cfg.GiphyTags[expectedKey]; !exists {
			t.Errorf("expected giphy tag key %s not found", expectedKey)
		} else if actualValue != expectedValue {
			t.Errorf("expected giphy tag value %s for key %s, got %s", expectedValue, expectedKey, actualValue)
		}
	}
}
