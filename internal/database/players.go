package database

import (
    //"database/sql"
    "errors"
    "fmt"
)

func AddPlayer(username string) error {
    if DB == nil {
        return errors.New("database not initialized")
    }

    stmt, err := DB.Prepare("INSERT INTO players(username) VALUES(?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(username)
    if err != nil {
        return fmt.Errorf("failed to add player: %w", err)
    }
    return nil
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

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return players, nil
}
