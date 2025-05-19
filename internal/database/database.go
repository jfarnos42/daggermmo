package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbFile string) error {
    var err error
    DB, err = sql.Open("sqlite3", dbFile)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    if err = createTables(); err != nil {
        return fmt.Errorf("failed to create tables: %w", err)
    }

    return nil
}

func createTables() error {
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS players (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        level INTEGER DEFAULT 1,
        experience INTEGER DEFAULT 0
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

func PingDB() error {
    if DB == nil {
        return fmt.Errorf("database not initialized")
    }
    return DB.Ping()
}
