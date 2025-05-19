package database

import (
    "testing"
)

func TestLoginLogout(t *testing.T) {
    cleanup := setupTestDB(t)
    defer cleanup()

    // First, add a player to be able to log in
    username := "testloginuser"
    role := "player"

    err := AddPlayer(username, role)
    if err != nil {
        t.Fatalf("AddPlayer failed: %v", err)
    }

    // Test login (should generate a token)
    token, err := Login(username)
    if err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    if token == "" {
        t.Fatal("Login returned empty token")
    }

    // Test logout with a valid token
    err = Logout(token)
    if err != nil {
        t.Fatalf("Logout failed: %v", err)
    }

    // Test logout with an invalid token (should fail)
    err = Logout("invalidtoken123")
    if err == nil {
        t.Fatal("Logout with invalid token should fail but didn't")
    }
}
