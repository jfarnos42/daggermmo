package server

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/jfarnos42/daggermmo/internal/database"
)

func setupTestDB(t *testing.T) func() {
    err := database.InitDB(":memory:")
    if err != nil {
        t.Fatalf("Failed to init test DB: %v", err)
    }
    return func() {
        database.DB.Close()
        database.DB = nil
    }
}

func TestRootHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    w := httptest.NewRecorder()

    RootHandler(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }

    body := w.Body.String()
    expected := "Welcome to Daggerfall MMO Server"
    if body != expected {
        t.Errorf("Expected body '%s', got '%s'", expected, body)
    }
}

func TestHTTPHealthHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    w := httptest.NewRecorder()

    HTTPHealthHandler(w, req)

    if w.Body.String() != "HTTP server is healthy" {
        t.Errorf("Unexpected response body: %s", w.Body.String())
    }
}

func TestBDHealthHandler(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    req := httptest.NewRequest(http.MethodGet, "/dbhealth", nil)
    w := httptest.NewRecorder()

    BDHealthHandler(w, req)

    if w.Body.String() != "Database connection is healthy" {
        t.Errorf("Unexpected response body: %s", w.Body.String())
    }
}

func TestAddPlayerHandler(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    // Valid add player request
    req := httptest.NewRequest(http.MethodGet, "/addplayer?username=testuser&role=admin", nil)
    w := httptest.NewRecorder()

    AddPlayerHandler(w, req)

    if w.Result().StatusCode != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Result().StatusCode)
    }

    if !strings.Contains(w.Body.String(), "player added: testuser with role: admin") {
        t.Errorf("Unexpected body: %s", w.Body.String())
    }

    // Missing username
    req2 := httptest.NewRequest(http.MethodGet, "/addplayer", nil)
    w2 := httptest.NewRecorder()
    AddPlayerHandler(w2, req2)
    if w2.Result().StatusCode != http.StatusBadRequest {
        t.Errorf("Expected 400 Bad Request for missing username, got %d", w2.Result().StatusCode)
    }
}

func TestListPlayersHandler(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    // Add players first
    _ = database.AddPlayer("user1", "player")
    _ = database.AddPlayer("user2", "admin")

    req := httptest.NewRequest(http.MethodGet, "/listplayers", nil)
    w := httptest.NewRecorder()

    ListPlayersHandler(w, req)

    if w.Result().StatusCode != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Result().StatusCode)
    }

    body := w.Body.String()
    if !strings.Contains(body, "user1") || !strings.Contains(body, "user2") {
        t.Errorf("Response body does not contain expected players: %s", body)
    }
}

func TestLoginHandler(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    username := "loginuser"
    _ = database.AddPlayer(username, "player")

    // Valid login
    req := httptest.NewRequest(http.MethodGet, "/login?username="+username, nil)
    w := httptest.NewRecorder()
    LoginHandler(w, req)

    if w.Result().StatusCode != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Result().StatusCode)
    }

    if !strings.Contains(w.Body.String(), "token") || !strings.Contains(w.Body.String(), "role") {
        t.Errorf("Response missing token or role: %s", w.Body.String())
    }

    // Missing username
    req2 := httptest.NewRequest(http.MethodGet, "/login", nil)
    w2 := httptest.NewRecorder()
    LoginHandler(w2, req2)
    if w2.Result().StatusCode != http.StatusBadRequest {
        t.Errorf("Expected 400 Bad Request for missing username, got %d", w2.Result().StatusCode)
    }
}

func TestLogoutHandler(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    username := "logoutuser"
    _ = database.AddPlayer(username, "player")
    token, _ := database.Login(username)

    // Valid logout
    req := httptest.NewRequest(http.MethodGet, "/logout?token="+token, nil)
    w := httptest.NewRecorder()
    LogoutHandler(w, req)

    if w.Result().StatusCode != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Result().StatusCode)
    }

    // Invalid logout token
    req2 := httptest.NewRequest(http.MethodGet, "/logout?token=invalidtoken", nil)
    w2 := httptest.NewRecorder()
    LogoutHandler(w2, req2)

    if w2.Result().StatusCode != http.StatusUnauthorized {
        t.Errorf("Expected 401 Unauthorized for invalid token, got %d", w2.Result().StatusCode)
    }

    // Missing token param
    req3 := httptest.NewRequest(http.MethodGet, "/logout", nil)
    w3 := httptest.NewRecorder()
    LogoutHandler(w3, req3)
    if w3.Result().StatusCode != http.StatusBadRequest {
        t.Errorf("Expected 400 Bad Request for missing token, got %d", w3.Result().StatusCode)
    }
}

func TestGetPlayerRoleHandler(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    username := "roleuser"
    _ = database.AddPlayer(username, "admin")

    // Valid request
    req := httptest.NewRequest(http.MethodGet, "/getplayerrole?username="+username, nil)
    w := httptest.NewRecorder()
    GetPlayerRoleHandler(w, req)

    if w.Result().StatusCode != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Result().StatusCode)
    }

    body := w.Body.String()
    if body != "admin" {
        t.Errorf("Expected role 'admin', got '%s'", body)
    }

    // Missing username param
    req2 := httptest.NewRequest(http.MethodGet, "/getplayerrole", nil)
    w2 := httptest.NewRecorder()
    GetPlayerRoleHandler(w2, req2)
    if w2.Result().StatusCode != http.StatusBadRequest {
        t.Errorf("Expected 400 Bad Request for missing username, got %d", w2.Result().StatusCode)
    }
}
