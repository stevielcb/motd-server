package main

import (
  "fmt"
  "net"
  "os"
  "io/ioutil"
  "path/filepath"
  "math/rand"
  "time"
)

func startServer() {
  l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.ListenHost, c.ListenPort))
  if err != nil {
    os.Exit(1)
  }

  defer l.Close()
  for {
    conn, err := l.Accept()
    if err != nil {
      os.Exit(1)
    }
    go handleRequest(conn)
  }
}

func handleRequest(conn net.Conn) {
  var files []string
  err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
    if info.IsDir() {
      return nil
    }
    files = append(files, path)
    return nil
  })
  if err != nil {
    conn.Close()
  }
  rand.Seed(time.Now().UnixNano())
  randFile := files[rand.Intn(len(files))]
  dat, err := ioutil.ReadFile(randFile)
  if err != nil {
    conn.Close()
  }
  conn.Write(dat)
  conn.Close()
}
