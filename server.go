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
