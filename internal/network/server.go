package network

import (
    "crypto/tls"
    "fmt"
    "log"
    "net"
    //"os"
)

type Server struct {
    addr string
}

func NewServer(addr string) *Server {
    return &Server{addr: addr}
}

// StartTLS inicia un servidor TLS escuchando en s.addr, usando cert y key archivos.
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

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    log.Printf("New client connected: %s", conn.RemoteAddr())

    // Mensaje de bienvenida
    _, err := conn.Write([]byte("Welcome to the Daggerfall MMO server over TLS!\n"))
    if err != nil {
        log.Printf("Failed to send welcome message: %v", err)
        return
    }

    // Aquí iría la lógica para manejar la conexión del cliente
    // Por ejemplo, leer comandos, enviar respuestas, etc.

    // Ejemplo básico: cerrar conexión después de bienvenida
    log.Printf("Closing connection for client: %s", conn.RemoteAddr())
}
