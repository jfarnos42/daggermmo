package server

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/jfarnos42/daggermmo/internal/database"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to Daggerfall MMO Server"))
}

func HTTPHealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("HTTP server is healthy"))
}

func BDHealthHandler(w http.ResponseWriter, r *http.Request) {
    if database.DB == nil {
        http.Error(w, "database not initialized", http.StatusInternalServerError)
        return
    }
    w.Write([]byte("Database connection is healthy"))
}

func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("AddPlayerHandler called")
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "username is required", http.StatusBadRequest)
        return
    }

    role := r.URL.Query().Get("role")
    if role == "" {
        role = "player"
    }

    err := database.AddPlayer(username, role)
    if err != nil {
        http.Error(w, "failed to add player: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("player added: " + username + " with role: " + role))
}

func ListPlayersHandler(w http.ResponseWriter, r *http.Request) {
    players, err := database.ListPlayers()
    if err != nil {
        http.Error(w, "failed to list players: "+err.Error(), http.StatusInternalServerError)
        return
    }

    WriteJSON(w, http.StatusOK, players)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("LoginHandler called")
    username := r.URL.Query().Get("username")
    if username == "" {
        WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username is required"})
        return
    }

    token, err := database.Login(username)
    if err != nil {
        WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
        return
    }

    role, err := database.GetPlayerRole(username)
    if err != nil {
        // Opcional: definir un role por defecto o manejar el error
        role = "player"
    }

    WriteJSON(w, http.StatusOK, map[string]string{
        "token": token,
        "role":  role,
    })
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("LogoutHandler called")
    token := r.URL.Query().Get("token")
    if token == "" {
        http.Error(w, "token is required", http.StatusBadRequest)
        return
    }

    err := database.Logout(token)
    if err != nil {
        if err.Error() == "invalid token or session not found" {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        http.Error(w, "failed to logout: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("user logged out"))
}

func GetPlayerRoleHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlayerRoleHandler called")
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "username is required", http.StatusBadRequest)
        return
    }

    role, err := database.GetPlayerRole(username)
    if err != nil {
        http.Error(w, "failed to get player role: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(role))
}

// WriteJSON is a helper function to write JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
