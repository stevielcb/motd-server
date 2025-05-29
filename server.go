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
	"math/rand"
	"net"
	"os"
	"path/filepath"
)

func startServer() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.ListenHost, c.ListenPort))
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("Server started")

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
	defer conn.Close()

	var files []string

	err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		fmt.Println("Error when walking cache dir path")
		return
	}

	if len(files) == 0 {
		fmt.Println("No cached files found")
		return
	}

	randFile := files[rand.Intn(len(files))]
	dat, err := os.ReadFile(randFile)
	if err != nil {
		fmt.Println("Error when reading cached file from disk")
		return
	}
	conn.Write(dat)
}
