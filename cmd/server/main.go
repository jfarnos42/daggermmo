package main

import (
    "flag"
    "log"
    "net/http"

    "github.com/jfarnos42/daggermmo/internal/database"
    "github.com/jfarnos42/daggermmo/internal/network"
    "github.com/jfarnos42/daggermmo/internal/server"
    "github.com/jfarnos42/daggermmo/internal/commands"
)

func main() {
    initDB := flag.Bool("initdb", false, "Initialize the database and exit")
    dbPath := flag.String("dbpath", "/home/daggeruser/daggerfall.db", "Path to the database file")
    flag.Parse()

    log.Println("Daggerfall MMO Server started.")

    // Usa dbPath que viene por flag para inicializar DB
    err := database.InitDB(*dbPath)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    if *initDB {
        log.Println("Database initialized successfully (initdb mode).")
        return
    }

    // Crear servidor de juego
    gameServer := network.NewServer(":7777")

    // Start subsystems
    go startHTTPServer()
    go startGameServer(gameServer)
    go commands.StartCommandPrompt(gameServer)

    // Bloqueo infinito
    select {}
}

func startGameServer(s *network.Server) {
    err := s.StartTLS("cert.pem", "key.pem")
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
    mux.HandleFunc("/getplayerrole", server.GetPlayerRoleHandler)

    log.Println("HTTP health server listening on :8080")
    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        log.Fatalf("HTTP server failed: %v", err)
    }
}
