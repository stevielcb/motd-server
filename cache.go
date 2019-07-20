package main

import (
  b64 "encoding/base64"
  "os"
  "io"
  "net/http"
  "time"
  "fmt"
  "bytes"
)

func writeToCache(url string, msg string) {
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

  cacheFile := fmt.Sprintf("%s/%d", cacheDir, time.Now().UnixNano())
  f, err := os.Create(cacheFile)
  if err != nil {
    return
  }
  defer f.Close()

  encoded := b64.StdEncoding.EncodeToString(buf.Bytes())
  f.WriteString(fmt.Sprintf("%s1337;File=;inline=1:%s%s\n", start, encoded, end))
  if msg != "" {
    f.WriteString(msg + "\n")
  }
  f.Sync()
}
