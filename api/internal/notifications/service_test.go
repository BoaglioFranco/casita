package notifications

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/SherClockHolmes/webpush-go"
)

const (
	testP256dh = "BNNL5ZaTfK81qhXOx23-wewhigUeFb632jN6LvRWCFH1ubQr77FE_9qV1FuojuRmHP42zmf34rXgW80OvUVDgTk"
	testAuth   = "zqbxT6JKstKSY9JKibZLSQ"
)

func TestNotifyMembersContinuesAfterFailureAndDeletesExpiredSubscription(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, err := database.New(ctx, filepath.Join(t.TempDir(), "notifications.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.ExecContext(ctx, `
		CREATE TABLE push_subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			household_id INTEGER NOT NULL,
			member_id INTEGER NOT NULL,
			endpoint TEXT NOT NULL UNIQUE,
			p256dh TEXT NOT NULL,
			auth TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
			updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	requests := 0
	pushServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.URL.Path == "/expired" {
			http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(pushServer.Close)

	store := NewStore(db)
	for memberID, endpoint := range []string{pushServer.URL + "/expired", pushServer.URL + "/current"} {
		err := store.StoreSubscription(ctx, StoreSubscriptionParams{
			HouseholdId: 1,
			MemberId:    int64(memberID + 1),
			Endpoint:    endpoint,
			P256dh:      testP256dh,
			Auth:        testAuth,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		t.Fatal(err)
	}
	service := NewNotificationService(db, publicKey, privateKey)

	err = service.NotifyMembers(ctx, 1, Notification{Title: "Test", Body: "Test"})
	if err == nil || !strings.Contains(err.Error(), "410 Gone") {
		t.Fatalf("expected a 410 delivery error, got %v", err)
	}
	if requests != 2 {
		t.Fatalf("expected both subscriptions to be attempted, got %d requests", requests)
	}

	subscriptions, err := store.GetSubscriptionsByHousehold(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(subscriptions) != 1 || subscriptions[0].Endpoint != pushServer.URL+"/current" {
		t.Fatalf("expected only the current subscription to remain, got %#v", subscriptions)
	}
}

func TestStoreSubscriptionUpdatesExistingEndpoint(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, err := database.New(ctx, filepath.Join(t.TempDir(), "notifications.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.ExecContext(ctx, `
		CREATE TABLE push_subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			household_id INTEGER NOT NULL,
			member_id INTEGER NOT NULL,
			endpoint TEXT NOT NULL UNIQUE,
			p256dh TEXT NOT NULL,
			auth TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
			updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	store := NewStore(db)
	first := StoreSubscriptionParams{HouseholdId: 1, MemberId: 1, Endpoint: "https://push.example.test/one", P256dh: "old", Auth: "old"}
	if err := store.StoreSubscription(ctx, first); err != nil {
		t.Fatal(err)
	}
	updated := StoreSubscriptionParams{HouseholdId: 1, MemberId: 2, Endpoint: first.Endpoint, P256dh: "new", Auth: "new"}
	if err := store.StoreSubscription(ctx, updated); err != nil {
		t.Fatal(err)
	}

	subscriptions, err := store.GetSubscriptionsByHousehold(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(subscriptions) != 1 {
		t.Fatalf("expected one subscription, got %d", len(subscriptions))
	}
	if subscriptions[0].MemberID != 2 || subscriptions[0].P256dh != "new" || subscriptions[0].Auth != "new" {
		t.Fatalf("subscription was not updated: %#v", subscriptions[0])
	}
}
