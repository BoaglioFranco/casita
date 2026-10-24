package database

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestPostgresConnectionAndTransactionRollback(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	name := "transaction-rollback-" + time.Now().UTC().Format("20060102150405.000000000")
	wantErr := errors.New("force rollback")
	err = db.WithTx(ctx, func(tx Executor) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO households (name) VALUES ($1)`, name); err != nil {
			return err
		}
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("WithTx() error = %v, want %v", err, wantErr)
	}

	var count int
	if err := db.GetContext(ctx, &count, `SELECT COUNT(*) FROM households WHERE name = $1`, name); err != nil {
		t.Fatalf("check rolled-back household: %v", err)
	}
	if count != 0 {
		t.Fatalf("rolled-back household count = %d, want 0", count)
	}
}
