package network

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// Server represents a TLS game server.
type Server struct {
	addr    string
	clients map[string]net.Conn
	mu      sync.Mutex
}

// NewServer creates a new Server listening on addr.
func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		clients: make(map[string]net.Conn),
	}
}

// StartTLS starts the TLS server listening on s.addr using certFile and keyFile.
func (s *Server) StartTLS(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("loading TLS certificate failed: %w", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", s.addr, config)
	if err != nil {
		return fmt.Errorf("failed to start TLS listener: %w", err)
	}
	defer listener.Close()

	log.Printf("TLS game server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				return nil // listener closed
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection manages a connected client, keeps connection alive to receive simple commands.
func (s *Server) handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()

	// Add client to map
	s.mu.Lock()
	s.clients[addr] = conn
	s.mu.Unlock()

	log.Printf("New client connected: %s", addr)

	// Send welcome message
	_, err := conn.Write([]byte("Welcome to the Daggerfall MMO server over TLS!\nType 'quit' to disconnect.\n"))
	if err != nil {
		log.Printf("Failed to send welcome message to %s: %v", addr, err)
		s.removeClient(addr)
		return
	}

	// Read loop to keep connection alive and optionally handle simple commands
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected: %s", addr)
			} else {
				log.Printf("Read error from %s: %v", addr, err)
			}
			break
		}

		line = strings.TrimSpace(line)
		if line == "quit" || line == "exit" {
			log.Printf("Client requested disconnect: %s", addr)
			break
		}

		// For now, just echo back the command
		_, err = conn.Write([]byte("You said: " + line + "\n"))
		if err != nil {
			log.Printf("Write error to %s: %v", addr, err)
			break
		}
	}

	log.Printf("Closing connection for client: %s", addr)
	s.removeClient(addr)
}

// removeClient safely removes and closes a client connection.
func (s *Server) removeClient(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conn, ok := s.clients[addr]; ok {
		conn.Close()
		delete(s.clients, addr)
	}
}

// ListClients returns a slice of connected client addresses.
func (s *Server) ListClients() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := make([]string, 0, len(s.clients))
	for addr := range s.clients {
		clients = append(clients, addr)
	}
	return clients
}

// Addr returns the server listen address.
func (s *Server) Addr() string {
	return s.addr
}
