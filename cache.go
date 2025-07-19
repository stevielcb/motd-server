// cache.go
//
// This file provides functionality for caching downloaded MOTD content.
// It downloads a file from a given URL, encodes it, and writes it to the cache directory
// in a structured format for later retrieval by the server.

package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	// CacheFilePrefix is the prefix used for cached file content format
	CacheFilePrefix = "1337"
	// CacheFileFormat is the format string for cached file content
	CacheFileFormat = "%s;File=inline=1;size=%d;name=%s:%s"
)

// writeToCache downloads content from the specified URL and saves it into the local cache directory.
//
// The downloaded content is base64-encoded and written alongside optional metadata (msg).
// Files are named using a timestamp and a base64-encoded version of the URL.
//
// Arguments:
// - url: the URL to download content from.
// - msg: an optional message string appended after the encoded content.
func writeToCache(url string, msg string) {
	slog.Info("caching content", "url", url, "message", msg)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("failed to download content", "url", url, "error", err)
		return
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		slog.Error("failed to read response body", "url", url, "error", err)
		return
	}

	b64url := b64.StdEncoding.EncodeToString([]byte(url))

	cacheFile := fmt.Sprintf("%s/%d_%s", cacheDir, time.Now().UnixNano(), b64url)
	f, err := os.Create(cacheFile)
	if err != nil {
		slog.Error("failed to create cache file", "file", cacheFile, "error", err)
		return
	}
	defer f.Close()

	encoded := b64.StdEncoding.EncodeToString(buf.Bytes())
	var content string
	if msg != "" {
		content = fmt.Sprintf(CacheFileFormat+"%s\n", CacheFilePrefix, buf.Len(), b64url, encoded, msg)
	} else {
		content = fmt.Sprintf(CacheFileFormat, CacheFilePrefix, buf.Len(), b64url, encoded)
	}

	if _, err := f.WriteString(content); err != nil {
		slog.Error("failed to write to cache file", "file", cacheFile, "error", err)
		return
	}

	if err := f.Sync(); err != nil {
		slog.Error("failed to sync cache file", "file", cacheFile, "error", err)
		return
	}

	slog.Debug("successfully cached content", "file", cacheFile, "size", buf.Len())
}
