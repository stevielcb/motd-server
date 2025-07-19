// motd-server
//
// A lightweight TCP server that serves a random file from a specified cache directory
// upon connection. Intended for use in scenarios where a "message of the day" (MOTD)
// or similar rotating message needs to be served.
//
// How it works:
// - Listens on the address and port specified in the configuration (c.ListenHost and c.ListenPort).
// - On each incoming connection, randomly selects a file from the cache directory.
// - Sends the contents of the selected file to the client and then closes the connection.
//
// Main Functions:
// - startServer(): Initializes the TCP listener and accepts incoming connections.
// - handleRequest(conn net.Conn): Handles an individual client connection, selecting and sending a file.
//
// Configuration assumptions:
// - `cacheDir`: Path to the directory containing cached message files.
// - `c`: Global configuration object containing `ListenHost` and `ListenPort`.

package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"path/filepath"
)

func startServer() {
	addr := fmt.Sprintf("%s:%d", c.ListenHost, c.ListenPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to start server", "address", addr, "error", err)
		os.Exit(1)
	}

	slog.Info("server started", "address", addr)

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	var files []string

	err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		slog.Error("failed to walk cache directory", "cacheDir", cacheDir, "error", err)
		return
	}

	if len(files) == 0 {
		slog.Warn("no cached files found", "cacheDir", cacheDir)
		return
	}

	randFile := files[rand.Intn(len(files))]
	dat, err := os.ReadFile(randFile)
	if err != nil {
		slog.Error("failed to read cached file", "file", randFile, "error", err)
		return
	}

	if _, err := conn.Write(dat); err != nil {
		slog.Error("failed to write to connection", "error", err)
	}
}
