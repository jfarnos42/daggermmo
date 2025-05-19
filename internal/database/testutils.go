// internal/database/testutils.go
package database

import (
    "os"
    "testing"
)

func setupTestDB(t *testing.T) func() {
    testDBPath := "/tmp/test.db"
    err := InitDB(testDBPath)
    if err != nil {
        t.Fatalf("Failed to initialize test database: %v", err)
    }

    return func() {
        DB.Close()
        os.Remove(testDBPath)
    }
}
