package database

import (
    "crypto/rand"
    "database/sql"
    "encoding/hex"
    "errors"
    "fmt"
)

func Login(username string) (string, error) {
    if DB == nil {
        return "", errors.New("database not initialized")
    }

    var playerID int
    err := DB.QueryRow("SELECT id FROM players WHERE username = ?", username).Scan(&playerID)
    if err == sql.ErrNoRows {
        return "", errors.New("user does not exist")
    }
    if err != nil {
        return "", err
    }

    tokenBytes := make([]byte, 16)
    _, err = rand.Read(tokenBytes)
    if err != nil {
        return "", fmt.Errorf("failed to generate token: %w", err)
    }
    token := hex.EncodeToString(tokenBytes)

    _, err = DB.Exec("INSERT INTO sessions(player_id, token) VALUES (?, ?)", playerID, token)
    if err != nil {
        return "", fmt.Errorf("failed to create session: %w", err)
    }

    return token, nil
}

func Logout(token string) error {
    if DB == nil {
        return errors.New("database not initialized")
    }

    res, err := DB.Exec("DELETE FROM sessions WHERE token = ?", token)
    if err != nil {
        return fmt.Errorf("failed to delete session: %w", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rowsAffected == 0 {
        return errors.New("invalid token or session not found")
    }

    return nil
}
