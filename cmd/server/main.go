package main

import (
    "log"
    "net/http"

    "github.com/jfarnos42/daggermmo/internal/database"
    "github.com/jfarnos42/daggermmo/internal/network"
    "github.com/jfarnos42/daggermmo/internal/server"
)

func main() {
    log.Println("Daggerfall MMO Server started.")

    err := database.InitDB("/home/daggeruser/daggerfall.db")
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    log.Println("Database initialized successfully.")

    go startHTTPServer()
    go startGameServer()

    select {}
}

func startGameServer() {
    server := network.NewServer(":7777")

    err := server.StartTLS("cert.pem", "key.pem")
    if err != nil {
        log.Fatalf("Failed to start TLS game server: %v", err)
    }
}

func startHTTPServer() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", server.RootHandler)
    mux.HandleFunc("/httphealth", server.HTTPHealthHandler)
    mux.HandleFunc("/bdhealth", server.BDHealthHandler)
    mux.HandleFunc("/addplayer", server.AddPlayerHandler)
    mux.HandleFunc("/listplayers", server.ListPlayersHandler)
    mux.HandleFunc("/login", server.LoginHandler)
    mux.HandleFunc("/logout", server.LogoutHandler)

    log.Println("HTTP health server listening on :8080")
    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        log.Fatalf("HTTP server failed: %v", err)
    }
}
