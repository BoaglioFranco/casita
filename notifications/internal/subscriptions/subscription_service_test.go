package subscriptions

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/BoaglioFranco/casita/notifications/internal/database"
	"github.com/BoaglioFranco/casita/notifications/internal/migrations"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
)

func TestNotifyCleansUpOnlyExpiredSubscriptions(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	for _, tc := range []struct {
		name       string
		status     int
		failDelete bool
		wantError  bool
		wantRows   int
	}{
		{"not found", http.StatusNotFound, false, false, 1},
		{"gone", http.StatusGone, false, false, 1},
		{"rate limited", http.StatusTooManyRequests, false, true, 2},
		{"provider failure", http.StatusInternalServerError, false, true, 2},
		{"cleanup failure", http.StatusGone, true, true, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			admin, err := sql.Open("postgres", databaseURL)
			if err != nil {
				t.Fatal(err)
			}
			defer admin.Close()
			schema := fmt.Sprintf("cleanup_test_%d", time.Now().UnixNano())
			if _, err := admin.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if _, err := admin.Exec("DROP SCHEMA " + schema + " CASCADE"); err != nil {
					t.Error(err)
				}
			}()
			u, err := url.Parse(databaseURL)
			if err != nil {
				t.Fatal(err)
			}
			query := u.Query()
			query.Set("search_path", schema)
			u.RawQuery = query.Encode()
			if err := migrations.Up(u.String()); err != nil {
				t.Fatal(err)
			}
			db, err := database.New(ctx, u.String())
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			if tc.failDelete {
				_, err := db.ExecContext(ctx, `
					CREATE FUNCTION reject_delete() RETURNS trigger LANGUAGE plpgsql AS $$
					BEGIN RAISE EXCEPTION 'cleanup failed'; END; $$;
					CREATE TRIGGER reject_delete BEFORE DELETE ON push_subscriptions
					FOR EACH ROW EXECUTE FUNCTION reject_delete();`)
				if err != nil {
					t.Fatal(err)
				}
			}

			var attempts atomic.Int32
			provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts.Add(1)
				if r.URL.Path == "/problem" {
					w.WriteHeader(tc.status)
					return
				}
				w.WriteHeader(http.StatusCreated)
			}))
			defer provider.Close()
			store := NewSubscriptionStore(db)
			userID := uuid.New()
			for _, path := range []string{"/problem", "/healthy"} {
				if err := store.StoreSubscription(ctx, StoreSubscriptionParams{
					UserID: userID, Endpoint: provider.URL + path,
					P256dh: "BNNL5ZaTfK81qhXOx23-wewhigUeFb632jN6LvRWCFH1ubQr77FE_9qV1FuojuRmHP42zmf34rXgW80OvUVDgTk",
					Auth:   "zqbxT6JKstKSY9JKibZLSQ",
				}); err != nil {
					t.Fatal(err)
				}
			}
			privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
			if err != nil {
				t.Fatal(err)
			}
			service := NewNotificationService(db, publicKey, privateKey)
			err = service.Notify(ctx, NotificationRequest{UserID: userID, Notification: Notification{Title: "Test", Body: "Test"}})
			if (err != nil) != tc.wantError {
				t.Fatalf("Notify error=%v, wantError=%v", err, tc.wantError)
			}
			if tc.failDelete && !strings.Contains(err.Error(), "remove expired push subscription") {
				t.Fatalf("expected cleanup failure, got %v", err)
			}
			if attempts.Load() != 2 {
				t.Fatalf("expected both devices attempted, got %d", attempts.Load())
			}
			remaining, err := store.GetSubscriptionsByUser(ctx, userID)
			if err != nil {
				t.Fatal(err)
			}
			if len(remaining) != tc.wantRows {
				t.Fatalf("remaining subscriptions=%d, want=%d", len(remaining), tc.wantRows)
			}
			if tc.wantRows == 1 && remaining[0].Endpoint != provider.URL+"/healthy" {
				t.Fatal("deleted the healthy subscription instead of the expired one")
			}
		})
	}
}
