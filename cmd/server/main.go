package main

import (
    "fmt"
    "log"
    "daggerfall-mmo-server/internal/network"
)

func main() {
    fmt.Println("Daggerfall MMO Server started.")

    server := network.NewServer(":7777")
    err := server.Start()
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
