package testutil

import (
    "database/sql"
    "testing"

    "github.com/crispy-2/backend/internal/dal"
    _ "github.com/mattn/go-sqlite3"
    "github.com/crispy-2/backend/internal/model"
    "github.com/pressly/goose/v3"
    "github.com/shopspring/decimal"
)

var DB *sql.DB

// NewTestDB creates a fresh in-memory SQLite database with all migrations applied.
func NewTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
    if err != nil {
        t.Fatalf("failed to open in-memory database: %v", err)
    }
    t.Cleanup(func() { db.Close() })

    goose.SetBaseFS(dal.Migrations)
    if err := goose.SetDialect("sqlite3"); err != nil {
        t.Fatalf("failed to set goose dialect: %v", err)
    }
    if err := goose.Up(db, "migrations"); err != nil {
        t.Fatalf("failed to run migrations: %v", err)
    }
    return db
}

// AssertDecimalEqual compares two decimals and fails the test if they are not equal.
func AssertDecimalEqual(t *testing.T, expected, actual decimal.Decimal) {
    t.Helper()
    if !expected.Equal(actual) {
        t.Errorf("decimal mismatch: expected %s, got %s", expected.String(), actual.String())
    }
}

// AssertStringDecimalEqual compares two StringDecimal values.
func AssertStringDecimalEqual(t *testing.T, expected, actual model.StringDecimal) {
    t.Helper()
    if !expected.Equal(actual.Decimal) {
        t.Errorf("StringDecimal mismatch: expected %s, got %s", expected.String(), actual.String())
    }
}
