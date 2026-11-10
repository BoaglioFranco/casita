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
	"github.com/google/uuid"
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
	if err := db.QueryRow("SELECT to_regclass($1) IS NOT NULL", schema+".calendar_events").Scan(&exists); err != nil || !exists {
		t.Fatalf("expected migrated table: exists=%v, error=%v", exists, err)
	}
	m, err := New(testURL)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	version, dirty, err := m.Version()
	if err != nil || version != 3 || dirty {
		t.Fatalf("version=%d dirty=%v error=%v", version, dirty, err)
	}

	verifyIdentitySchema(t, testURL)
	// Accounts normally belong to the seeded household. Full rollback must also
	// work on a populated database despite restrictive household ownership.
	seedDB, err := sql.Open("postgres", testURL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := seedDB.Exec(`INSERT INTO users (email, password_hash, name, household_id)
		VALUES ('seeded@example.test', 'hash', 'Seeded', '019535d9-3df7-79fb-b466-fa907fa17f9e')`); err != nil {
		seedDB.Close()
		t.Fatal(err)
	}
	seedDB.Close()

	// A failed migration must block startup; it must never be auto-forced clean.
	if _, err := db.Exec("UPDATE " + schema + ".schema_migrations SET dirty = true"); err != nil {
		t.Fatal(err)
	}
	var dirtyErr migrate.ErrDirty
	if err := Up(testURL); !errors.As(err, &dirtyErr) {
		t.Fatalf("expected dirty migration error, got %v", err)
	}
	if err := m.Force(3); err != nil {
		t.Fatal(err)
	}
	if err := m.Down(); err != nil {
		t.Fatalf("rollback: %v", err)
	}
	if err := db.QueryRow("SELECT to_regclass($1) IS NOT NULL", schema+".calendar_events").Scan(&exists); err != nil || exists {
		t.Fatalf("expected table removed: exists=%v, error=%v", exists, err)
	}
	if err := Up(testURL); err != nil {
		t.Fatalf("migrate after rollback: %v", err)
	}
}

func verifyIdentitySchema(t *testing.T, databaseURL string) {
	t.Helper()
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatal(err)
		}
	}
	reject := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err == nil {
			t.Fatalf("expected constraint rejection: %s", query)
		}
	}
	count := func(query string, args ...any) int {
		t.Helper()
		var n int
		if err := db.QueryRow(query, args...).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	for _, table := range []string{"members", "push_subscriptions", "list_notification_subscriptions"} {
		var exists bool
		if err := db.QueryRow(`SELECT to_regclass($1) IS NOT NULL`, table).Scan(&exists); err != nil || exists {
			t.Fatalf("obsolete table %s: exists=%v err=%v", table, exists, err)
		}
	}
	household, user, list, item, event := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	exec(`INSERT INTO households(id,name) VALUES($1,'Constraint test')`, household)
	reject(`INSERT INTO users(email,password_hash,name) VALUES('missing@example.test','hash','Missing')`)
	reject(`INSERT INTO users(email,password_hash,name,household_id) VALUES('foreign@example.test','hash','Foreign',$1)`, uuid.New())
	exec(`INSERT INTO users(id,email,password_hash,name,household_id) VALUES($1,'valid@example.test','hash','Valid',$2)`, user, household)
	reject(`DELETE FROM households WHERE id=$1`, household)
	reject(`INSERT INTO users(id,email,password_hash,name,household_id) VALUES($1,'duplicate@example.test','hash','Duplicate',$2)`, user, household)
	exec(`INSERT INTO lists(id,household_id,name) VALUES($1,$2,'Test')`, list, household)
	reject(`INSERT INTO items(list_id,name,created_by_user_id) VALUES($1,'Invalid',$2)`, list, uuid.New())
	exec(`INSERT INTO items(id,list_id,name,created_by_user_id) VALUES($1,$2,'Test',$3)`, item, list, user)
	exec(`INSERT INTO activity_events(household_id,list_id,item_id,actor_user_id,event_type) VALUES($1,$2,$3,$4,'ITEM_CREATE')`, household, list, item, user)
	reject(`INSERT INTO activity_events(household_id,actor_user_id,event_type) VALUES($1,$2,'ITEM_CREATE')`, household, uuid.New())
	reject(`INSERT INTO calendar_events(household_id,created_by_user_id,name,starts_at,notification_offset_minutes) VALUES($1,$2,'Invalid',now(),10)`, household, uuid.New())
	exec(`INSERT INTO calendar_events(id,household_id,created_by_user_id,name,starts_at,notification_offset_minutes) VALUES($1,$2,$3,'Test',now(),10)`, event, household, user)
	reject(`INSERT INTO calendar_event_notification_settings(event_id,user_id) VALUES($1,$2)`, event, uuid.New())
	exec(`INSERT INTO calendar_event_notification_settings(event_id,user_id) VALUES($1,$2)`, event, user)
	reject(`INSERT INTO calendar_event_notification_settings(event_id,user_id) VALUES($1,$2)`, event, user)
	exec(`INSERT INTO sessions(token_hash,user_id,expires_at) VALUES('constraint-test',$1,now()+interval '1 day')`, user)
	exec(`DELETE FROM users WHERE id=$1`, user)
	if count(`SELECT count(*) FROM sessions WHERE user_id=$1`, user) != 0 || count(`SELECT count(*) FROM calendar_events WHERE id=$1`, event) != 0 || count(`SELECT count(*) FROM calendar_event_notification_settings WHERE event_id=$1`, event) != 0 {
		t.Fatal("user delete must cascade sessions, calendar authorship, settings")
	}
	if count(`SELECT count(*) FROM items WHERE id=$1 AND created_by_user_id IS NULL`, item) != 1 || count(`SELECT count(*) FROM activity_events WHERE household_id=$1 AND actor_user_id IS NULL`, household) != 1 {
		t.Fatal("user delete must retain items/activity with null authorship")
	}
	exec(`DELETE FROM households WHERE id=$1`, household)
}
