package database

import (
    //"database/sql"
    "errors"
    "fmt"
)

func AddPlayer(username, role string) error {
    if DB == nil {
        return errors.New("database not initialized")
    }

    if role == "" {
        role = "player"
    }

    stmt, err := DB.Prepare("INSERT INTO players(username, role) VALUES(?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(username, role)
    if err != nil {
        return fmt.Errorf("failed to add player: %w", err)
    }
    return nil
}

func GetPlayerRole(username string) (string, error) {
    if DB == nil {
        return "", errors.New("database not initialized")
    }

    var role string
    err := DB.QueryRow("SELECT role FROM players WHERE username = ?", username).Scan(&role)
    if err != nil {
        return "", err
    }

    return role, nil
}

func ListPlayers() ([]string, error) {
    if DB == nil {
        return nil, errors.New("database not initialized")
    }

    rows, err := DB.Query("SELECT username FROM players")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var players []string
    for rows.Next() {
        var username string
        if err := rows.Scan(&username); err != nil {
            return nil, err
        }
        players = append(players, username)
    }
    return players, nil
}
