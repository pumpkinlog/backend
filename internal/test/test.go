package test

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func Ptr[T any](v T) *T {
	return &v
}

func NewPgxConn(t *testing.T) *pgx.Conn {
	t.Helper()

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		t.Skip("skipping due to missing environment variable DATABASE_DSN")
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse DSN: %v", err)
	}

	conn, err := pgx.ConnectConfig(t.Context(), config)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(t.Context()); err != nil {
			t.Fatalf("failed to close connection: %v", err)
		}
	})

	return conn
}
