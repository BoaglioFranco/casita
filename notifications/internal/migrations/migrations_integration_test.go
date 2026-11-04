package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

func TestMigrationLifecycle(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	schema := fmt.Sprintf("migration_test_%d", time.Now().UnixNano())
	if _, err := db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := db.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
			t.Errorf("clean up test schema: %v", err)
		}
	}()
	u, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	query := u.Query()
	query.Set("search_path", schema)
	u.RawQuery = query.Encode()
	testURL := u.String()

	// Concurrent first startups must coordinate, then repeated startup is a no-op.
	results := make(chan error, 2)
	for range 2 {
		go func() { results <- Up(testURL) }()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Fatalf("concurrent migration: %v", err)
		}
	}
	if err := Up(testURL); err != nil {
		t.Fatalf("repeat migration: %v", err)
	}

	var exists bool
	if err := db.QueryRow("SELECT to_regclass($1) IS NOT NULL", schema+".push_subscriptions").Scan(&exists); err != nil || !exists {
		t.Fatalf("expected migrated table: exists=%v, error=%v", exists, err)
	}
	m, err := New(testURL)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	version, dirty, err := m.Version()
	if err != nil || version != 1 || dirty {
		t.Fatalf("version=%d dirty=%v error=%v", version, dirty, err)
	}

	// A failed migration must block startup; it must never be auto-forced clean.
	if _, err := db.Exec("UPDATE " + schema + ".schema_migrations SET dirty = true"); err != nil {
		t.Fatal(err)
	}
	var dirtyErr migrate.ErrDirty
	if err := Up(testURL); !errors.As(err, &dirtyErr) {
		t.Fatalf("expected dirty migration error, got %v", err)
	}
	if err := m.Force(1); err != nil {
		t.Fatal(err)
	}
	if err := m.Down(); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if err := db.QueryRow("SELECT to_regclass($1) IS NOT NULL", schema+".push_subscriptions").Scan(&exists); err != nil || exists {
		t.Fatalf("expected table removed: exists=%v, error=%v", exists, err)
	}
	if err := Up(testURL); err != nil {
		t.Fatalf("migrate after rollback: %v", err)
	}
}
