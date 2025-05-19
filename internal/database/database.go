package database

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) error {
    var err error
    DB, err = sql.Open("sqlite3", filepath)
    if err != nil {
        return err
    }

    err = DB.Ping()
    if err != nil {
        return err
    }

    err = createTables()
    if err != nil {
        return err
    }

    log.Println("Database initialized and tables created.")
    return nil
}

func createTables() error {
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS players (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        level INTEGER DEFAULT 1,
        experience INTEGER DEFAULT 0,
        role TEXT DEFAULT 'Player'
    );

    CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        player_id INTEGER,
        token TEXT NOT NULL UNIQUE,
        FOREIGN KEY(player_id) REFERENCES players(id)
    );
    `
    _, err := DB.Exec(sqlStmt)
    if err == nil {
        log.Println("Database tables checked/created successfully.")
    }
    return err
}
