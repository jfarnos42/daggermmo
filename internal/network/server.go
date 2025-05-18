package network

import (
    "fmt"
    "net"
)

type Server struct {
    addr string
}

func NewServer(addr string) *Server {
    return &Server{addr: addr}
}

func (s *Server) Start() error {
    ln, err := net.Listen("tcp", s.addr)
    if err != nil {
        return err
    }
    defer ln.Close()

    fmt.Printf("Server listening on %s\n", s.addr)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

    // Here we would read/write to the client
    // For now, just send a welcome message
    conn.Write([]byte("Welcome to the Daggerfall MMO server\n"))
}
