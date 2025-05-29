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
	"net/http"
	"os"
	"time"
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
	fmt.Printf("Caching url, %s, with message, %s\n", url, msg)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return
	}

	b64url := b64.StdEncoding.EncodeToString([]byte(url))

	cacheFile := fmt.Sprintf("%s/%d_%s", cacheDir, time.Now().UnixNano(), b64url)
	f, err := os.Create(cacheFile)
	if err != nil {
		return
	}
	defer f.Close()

	encoded := b64.StdEncoding.EncodeToString(buf.Bytes())
	f.WriteString(fmt.Sprintf("1337;File=inline=1;size=%d;name=%s:%s", buf.Len(), b64url, encoded))
	if msg != "" {
		f.WriteString(msg + "\n")
	}
	f.Sync()
}
