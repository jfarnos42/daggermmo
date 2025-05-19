package server

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/jfarnos42/daggermmo/internal/database"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("root ok\n"))
}

func HTTPHealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("http: ok"))
}

func BDHealthHandler(w http.ResponseWriter, r *http.Request) {
    err := database.PingDB()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("bd: error"))
        return
    }
    w.Write([]byte("bd: ok"))
}

func AddPlayerHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("AddPlayerHandler called")
    username := r.URL.Query().Get("username")
    if username == "" {
        http.Error(w, "username is required", http.StatusBadRequest)
        return
    }

    err := database.AddPlayer(username)
    if err != nil {
        http.Error(w, "failed to add player: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("player added: " + username))
}

func ListPlayersHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("ListPlayersHandler called")
    players, err := database.ListPlayers()
    if err != nil {
        http.Error(w, "failed to list players: "+err.Error(), http.StatusInternalServerError)
        return
    }

    for _, p := range players {
        w.Write([]byte(p + "\n"))
    }
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

    WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("logoutHandler called")
    token := r.URL.Query().Get("token")
    if token == "" {
        WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "token is required"})
        return
    }

    err := database.Logout(token)
    if err != nil {
        WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
        return
    }

    WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

