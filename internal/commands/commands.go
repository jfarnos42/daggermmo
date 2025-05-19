package commands

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/jfarnos42/daggermmo/internal/network"
)

var server *network.Server

func StartCommandPrompt(s *network.Server) {
    server = s
    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Command prompt started. Type 'help' for commands.")
    for {
        fmt.Print("> ")
        input, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Error reading input: %v", err)
            continue
        }
        input = strings.TrimSpace(input)

        switch input {
        case "help":
            printHelp()
        case "status":
            printStatus()
        case "who":
            listClients()
        case "quit", "exit":
            fmt.Println("Exiting command prompt.")
            return
        default:
            fmt.Println("Unknown command. Type 'help' for a list of commands.")
        }
    }
}

func printHelp() {
    fmt.Println("Available commands:")
    fmt.Println("  help        			- Show this help message")
    fmt.Println("  status      			- Show server status")
    fmt.Println("  who					- List connected clients")
    fmt.Println("  quit/exit   			- Exit command prompt")
}

func printStatus() {
    fmt.Println("Server running on", server.Addr())
}

func listClients() {
    clients := server.ListClients()
    if len(clients) == 0 {
        fmt.Println("No clients connected.")
        return
    }
    fmt.Println("Connected clients:")
    for _, addr := range clients {
        fmt.Println("-", addr)
    }
}
