package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/stevielcb/motd-server/internal/services"
)

// TCPServer represents a TCP server that serves cached content
type TCPServer struct {
	host     string
	port     int
	cache    services.CacheManager
	logger   *slog.Logger
	listener net.Listener
}

// NewTCPServer creates a new TCP server instance
func NewTCPServer(host string, port int, cache services.CacheManager, logger *slog.Logger) *TCPServer {
	return &TCPServer{
		host:   host,
		port:   port,
		cache:  cache,
		logger: logger,
	}
}

// Start begins listening for connections
func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.listener = l
	s.logger.Info("server started", "address", addr)

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			s.logger.Error("failed to accept connection", "error", err)
			return fmt.Errorf("failed to accept connection: %w", err)
		}
		go s.handleRequest(conn)
	}
}

// Stop gracefully stops the server
func (s *TCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// handleRequest handles an individual client connection
func (s *TCPServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	data, err := s.cache.GetRandomFile()
	if err != nil {
		s.logger.Error("failed to get random file", "error", err)
		return
	}

	if _, err := conn.Write(data); err != nil {
		s.logger.Error("failed to write to connection", "error", err)
	}
}
