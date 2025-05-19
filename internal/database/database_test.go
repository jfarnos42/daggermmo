package database

import (
    "os"
    "testing"
)

func TestInitDB(t *testing.T) {
    // Crear archivo temporal para la base de datos SQLite
    tmpfile, err := os.CreateTemp("", "testdb-*.sqlite")
    if err != nil {
        t.Fatalf("failed to create temp db file: %v", err)
    }
    defer os.Remove(tmpfile.Name()) // limpiar archivo al final

    // Inicializar la base de datos con archivo temporal
    err = InitDB(tmpfile.Name())
    if err != nil {
        t.Fatalf("InitDB failed: %v", err)
    }

    // Comprobar que la variable DB no es nil
    if DB == nil {
        t.Fatal("DB variable is nil after InitDB")
    }
}
