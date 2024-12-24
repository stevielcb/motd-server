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
