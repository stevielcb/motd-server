package server

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"
)

// Mock cache manager for testing
type mockCacheManager struct {
	shouldError bool
	returnData  []byte
}

func (m *mockCacheManager) WriteToCache(url string, msg string) error {
	if m.shouldError {
		return fmt.Errorf("mock write error")
	}
	return nil
}

func (m *mockCacheManager) GetRandomFile() ([]byte, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock get error")
	}
	return m.returnData, nil
}

func (m *mockCacheManager) Cleanup() error {
	if m.shouldError {
		return fmt.Errorf("mock cleanup error")
	}
	return nil
}

func TestNewTCPServer(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cacheManager := &mockCacheManager{}

	server := NewTCPServer("localhost", 8080, cacheManager, logger)

	if server == nil {
		t.Fatal("expected server but got nil")
	}

	if server.host != "localhost" {
		t.Errorf("expected host localhost, got %s", server.host)
	}

	if server.port != 8080 {
		t.Errorf("expected port 8080, got %d", server.port)
	}

	if server.cache != cacheManager {
		t.Error("cache manager not properly set")
	}

	if server.logger != logger {
		t.Error("logger not properly set")
	}
}

func TestTCPServer_Start(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cacheManager := &mockCacheManager{
		returnData: []byte("test data"),
	}

	server := NewTCPServer("localhost", 0, cacheManager, logger)

	// Test that server creation works
	if server == nil {
		t.Fatal("expected server but got nil")
	}

	// Test that server can be stopped even if not started
	if err := server.Stop(); err != nil {
		t.Errorf("failed to stop server: %v", err)
	}
}

func TestTCPServer_Start_InvalidPort(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cacheManager := &mockCacheManager{}

	// Test with a clearly invalid port (negative)
	server := NewTCPServer("localhost", -1, cacheManager, logger)

	// This should fail when trying to start
	err := server.Start()
	if err == nil {
		t.Error("expected error when starting server on invalid port, but got none")
	}
}

func TestTCPServer_Stop(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cacheManager := &mockCacheManager{}

	server := NewTCPServer("localhost", 8080, cacheManager, logger)

	// Test stopping server that hasn't been started
	err := server.Stop()
	if err != nil {
		t.Errorf("unexpected error when stopping unstarted server: %v", err)
	}
}

func TestTCPServer_HandleRequest_CacheError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cacheManager := &mockCacheManager{
		shouldError: true, // Simulate cache error
	}

	server := NewTCPServer("localhost", 8080, cacheManager, logger)

	// Create a mock connection
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Handle request in a goroutine
	go server.handleRequest(serverConn)

	// Give it a moment to process
	time.Sleep(100 * time.Millisecond)

	// The connection should be closed due to the cache error
	// Try to read from the connection - it should be closed
	_, err := clientConn.Read([]byte{})
	if err == nil {
		t.Error("expected connection to be closed due to cache error, but it's still open")
	}
}

func TestTCPServer_HandleRequest_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	expectedData := []byte("test response data")
	cacheManager := &mockCacheManager{
		returnData: expectedData,
	}

	server := NewTCPServer("localhost", 8080, cacheManager, logger)

	// Create a mock connection
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Handle request in a goroutine
	go server.handleRequest(serverConn)

	// Read response from client side
	response := make([]byte, len(expectedData))
	n, err := clientConn.Read(response)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if n != len(expectedData) {
		t.Errorf("expected %d bytes, got %d", len(expectedData), n)
	}

	if string(response) != string(expectedData) {
		t.Errorf("expected response %s, got %s", string(expectedData), string(response))
	}
}

func TestTCPServer_HandleRequest_WriteError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cacheManager := &mockCacheManager{
		returnData: []byte("test data"),
	}

	server := NewTCPServer("localhost", 8080, cacheManager, logger)

	// Create a mock connection that will fail on write
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	// Close the server connection immediately to simulate write error
	serverConn.Close()

	// Handle request - should not panic
	server.handleRequest(serverConn)
}
