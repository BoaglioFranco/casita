package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/BoaglioFranco/casita/api/internal/database"
	"github.com/BoaglioFranco/casita/api/migrations"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func integrationDatabase(t *testing.T) *database.Database {
	t.Helper()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	admin, err := database.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { admin.Close() })
	schema := "identity_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.ExecContext(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := admin.ExecContext(ctx, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Error(err)
		}
	})
	u, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	if err := migrations.Up(u.String()); err != nil {
		t.Fatal(err)
	}
	db, err := database.New(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestUserIdentityAndHouseholdIsolation(t *testing.T) {
	db := integrationDatabase(t)
	ctx := context.Background()
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			t.Fatal(err)
		}
	}
	scalar := func(query string, args ...any) int {
		t.Helper()
		var n int
		if err := db.GetContext(ctx, &n, query, args...); err != nil {
			t.Fatal(err)
		}
		return n
	}
	household, otherHousehold := uuid.New(), uuid.New()
	actor, peer, failedPeer, outsider := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	exec(`INSERT INTO households (id,name) VALUES ($1,'Home'),($2,'Other')`, household, otherHousehold)
	hash, err := bcrypt.GenerateFromPassword([]byte("test-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range []struct {
		id, household uuid.UUID
		name          string
	}{{actor, household, "Actor"}, {peer, household, "Peer"}, {failedPeer, household, "Failed peer"}, {outsider, otherHousehold, "Outsider"}} {
		exec(`INSERT INTO users (id,household_id,name,email,password_hash) VALUES ($1,$2,$3,$4,$5)`, u.id, u.household, u.name, u.id.String()+"@example.test", string(hash))
	}
	const secret = "integration-test-secret-at-least-32-bytes"
	t.Setenv("S2S_TOKEN", secret)
	type delivery struct {
		Recipient, Actor uuid.UUID
		Title, Body      string
	}
	var mu sync.Mutex
	var deliveries []delivery
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), claims, func(*jwt.Token) (any, error) { return []byte(secret), nil }, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("casita-api"), jwt.WithAudience("notifications"), jwt.WithExpirationRequired())
		if err != nil {
			t.Error(err)
			w.WriteHeader(401)
			return
		}
		if r.URL.Path == "/notifications/subscribe" {
			if claims.Subject != actor.String() {
				t.Errorf("subscription actor = %s", claims.Subject)
			}
			w.WriteHeader(200)
			return
		}
		if r.URL.Path != "/notifications/notify" {
			t.Errorf("unexpected upstream path %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		var payload struct {
			UserID       uuid.UUID `json:"user_id"`
			Notification struct {
				Title string `json:"title"`
				Body  string `json:"body"`
			} `json:"notification"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		subject, err := uuid.Parse(claims.Subject)
		if err != nil {
			t.Error(err)
		}
		mu.Lock()
		deliveries = append(deliveries, delivery{
			Recipient: payload.UserID, Actor: subject,
			Title: payload.Notification.Title, Body: payload.Notification.Body,
		})
		mu.Unlock()
		if payload.UserID == failedPeer {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	}))
	defer upstream.Close()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setupRoutes(router, db, upstream.URL)
	request := func(method, path, body string, cookie *http.Cookie, want int) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		if body != "" {
			r.Header.Set("Content-Type", "application/json")
		}
		if cookie != nil {
			r.AddCookie(cookie)
		}
		r.Header.Set("X-User-ID", outsider.String())
		r.Header.Set("X-Household-ID", otherHousehold.String())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if w.Code != want {
			t.Fatalf("%s %s: status %d want %d: %s", method, path, w.Code, want, w.Body.String())
		}
		return w
	}
	login := func(id uuid.UUID) *http.Cookie {
		t.Helper()
		w := request("POST", "/login", fmt.Sprintf(`{"email":%q,"password":"test-password"}`, id.String()+"@example.test"), nil, 200)
		for _, c := range w.Result().Cookies() {
			if c.Name == "session" {
				return c
			}
		}
		t.Fatal("missing session cookie")
		return nil
	}
	request("POST", "/login", fmt.Sprintf(`{"email":%q,"password":"wrong"}`, actor.String()+"@example.test"), nil, 401)
	cookie := login(actor)
	otherCookie := login(outsider)
	request("GET", "/me", "", nil, 401)
	request("GET", "/me", "", &http.Cookie{Name: "session", Value: "invalid"}, 401)
	expired := login(actor)
	tokenHash := sha256.Sum256([]byte(expired.Value))
	exec(`UPDATE sessions SET expires_at=now()-interval '1 second' WHERE token_hash=$1`, hex.EncodeToString(tokenHash[:]))
	request("GET", "/me", "", expired, 401)
	me := request("GET", "/me", "", cookie, 200)
	var current map[string]any
	if err := json.Unmarshal(me.Body.Bytes(), &current); err != nil {
		t.Fatal(err)
	}
	if len(current) != 3 || current["id"] != actor.String() || current["name"] != "Actor" || current["householdId"] != household.String() {
		t.Fatalf("unsafe or incorrect current user: %v", current)
	}
	// Session validation must read the current database values, not a cached household.
	exec(`UPDATE users SET name='Updated',household_id=$2 WHERE id=$1`, actor, otherHousehold)
	changed := request("GET", "/me", "", cookie, 200)
	if !strings.Contains(changed.Body.String(), otherHousehold.String()) || !strings.Contains(changed.Body.String(), "Updated") {
		t.Fatal(changed.Body.String())
	}
	exec(`UPDATE users SET name='Actor',household_id=$2 WHERE id=$1`, actor, household)
	for _, path := range []string{"/households", "/households/" + household.String() + "/me", "/households/" + household.String() + "/lists"} {
		request("GET", path, "", cookie, 404)
	}
	ownList, foreignList, foreignItem, foreignEvent := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	exec(`INSERT INTO lists(id,household_id,name) VALUES($1,$2,'Own list'),($3,$4,'Foreign list')`, ownList, household, foreignList, otherHousehold)
	exec(`INSERT INTO items(id,list_id,name,created_by_user_id) VALUES($1,$2,'Foreign item',$3)`, foreignItem, foreignList, outsider)
	exec(`INSERT INTO calendar_events(id,household_id,created_by_user_id,name,starts_at,notification_offset_minutes) VALUES($1,$2,$3,'Foreign event','2026-11-11T12:00:00Z',10)`, foreignEvent, otherHousehold, outsider)
	ownItems := "/lists/" + ownList.String() + "/items"
	request("GET", ownItems, "", cookie, 200)
	// Forged identity/household fields cannot change ownership or authorship.
	request("POST", "/lists", fmt.Sprintf(`{"name":"Created list","householdId":%q}`, otherHousehold), cookie, 201)
	if scalar(`SELECT count(*) FROM lists WHERE name='Created list' AND household_id=$1`, household) != 1 {
		t.Fatal("forged list household")
	}
	forged := fmt.Sprintf(`{"name":"Milk","userId":%q,"householdId":%q,"createdByUserId":%q,"listName":"Forged list","user":{"id":%q,"name":"Forged name","householdId":%q}}`, outsider, otherHousehold, outsider, outsider, otherHousehold)
	request("POST", ownItems, forged, cookie, 201)
	var itemID uuid.UUID
	if err := db.GetContext(ctx, &itemID, `SELECT id FROM items WHERE list_id=$1`, ownList); err != nil {
		t.Fatal(err)
	}
	got := request("GET", ownItems, "", cookie, 200)
	if !strings.Contains(got.Body.String(), actor.String()) || !strings.Contains(got.Body.String(), `"name":"Actor"`) {
		t.Fatal(got.Body.String())
	}
	if scalar(`SELECT count(*) FROM activity_events WHERE actor_user_id=$1 AND household_id=$2 AND item_id=$3`, actor, household, itemID) != 1 {
		t.Fatal("incorrect activity author")
	}
	mu.Lock()
	if len(deliveries) != 2 {
		t.Errorf("deliveries = %v", deliveries)
	}
	for _, d := range deliveries {
		if d.Actor != actor || (d.Recipient != peer && d.Recipient != failedPeer) {
			t.Errorf("incorrect delivery %+v", d)
		}
		if d.Title != "Item created" || d.Body != "Milk added to Own list by Actor" {
			t.Errorf("incorrect creation notification: %+v", d)
		}
	}
	mu.Unlock()
	eventBody := fmt.Sprintf(`{"name":"Own event","startsAt":"2026-11-11T12:00:00Z","durationMinutes":60,"householdId":%q,"createdByUserId":%q}`, otherHousehold, outsider)
	request("POST", "/calendar/events", eventBody, cookie, 201)
	calendar := request("GET", "/calendar?from=2026-11-11&to=2026-11-12", "", cookie, 200)
	if !strings.Contains(calendar.Body.String(), `"createdByUserId":"`+actor.String()+`"`) || strings.Contains(calendar.Body.String(), "Foreign event") || strings.Contains(calendar.Body.String(), "Member") {
		t.Fatal(calendar.Body.String())
	}
	var eventID uuid.UUID
	if err := db.GetContext(ctx, &eventID, `SELECT id FROM calendar_events WHERE household_id=$1`, household); err != nil {
		t.Fatal(err)
	}
	request("PUT", "/calendar/events/"+eventID.String(), `{"name":"Updated event","startsAt":"2026-11-11T13:00:00Z","durationMinutes":60}`, cookie, 204)
	beforeItems := scalar(`SELECT count(*) FROM items`)
	beforeActivity := scalar(`SELECT count(*) FROM activity_events`)
	beforeEvents := scalar(`SELECT count(*) FROM calendar_events`)
	for _, list := range []uuid.UUID{foreignList, uuid.New()} {
		path := "/lists/" + list.String() + "/items"
		request("GET", path, "", cookie, 404)
		request("POST", path, forged, cookie, 404)
		request("DELETE", path+"/"+foreignItem.String(), "", cookie, 404)
		request("DELETE", path+"/"+itemID.String(), "", cookie, 404)
	}
	request("DELETE", ownItems+"/"+foreignItem.String(), "", cookie, 404)
	for _, event := range []uuid.UUID{foreignEvent, uuid.New()} {
		path := "/calendar/events/" + event.String()
		request("PUT", path, eventBody, cookie, 404)
		request("DELETE", path, "", cookie, 404)
	}
	if scalar(`SELECT count(*) FROM items`) != beforeItems || scalar(`SELECT count(*) FROM activity_events`) != beforeActivity || scalar(`SELECT count(*) FROM calendar_events`) != beforeEvents {
		t.Fatal("rejected request mutated data")
	}
	mu.Lock()
	if len(deliveries) != 2 {
		t.Error("rejected request sent notifications")
	}
	mu.Unlock()
	lists := request("GET", "/lists", "", cookie, 200)
	if strings.Contains(lists.Body.String(), foreignList.String()) {
		t.Fatal("foreign list exposed")
	}
	request("GET", ownItems, "", otherCookie, 404)
	// A different household user removes the item. Use the current list and actor
	// names, rather than its creator's name or any names supplied by the client.
	peerCookie := login(peer)
	exec(`UPDATE users SET name='Updated peer' WHERE id=$1`, peer)
	exec(`UPDATE lists SET name='Renamed list' WHERE id=$1`, ownList)
	request("DELETE", ownItems+"/"+itemID.String(), "", peerCookie, 204)
	if scalar(`SELECT count(*) FROM activity_events WHERE actor_user_id=$1 AND event_type='ITEM_DELETE'`, peer) != 1 {
		t.Fatal("missing deletion activity")
	}
	mu.Lock()
	if len(deliveries) != 4 {
		t.Errorf("expected two create and two delete deliveries, got %v", deliveries)
	}
	for _, d := range deliveries[2:] {
		if d.Actor != peer || (d.Recipient != actor && d.Recipient != failedPeer) {
			t.Errorf("incorrect deletion recipient or actor: %+v", d)
		}
		if d.Title != "Item deleted" || d.Body != "Milk removed from Renamed list by Updated peer" {
			t.Errorf("incorrect deletion notification: %+v", d)
		}
	}
	mu.Unlock()
	request("DELETE", "/calendar/events/"+eventID.String(), "", cookie, 204)
	request("POST", "/notifications/subscribe", `{"endpoint":"https://push.example.test/device","keys":{"p256dh":"key","auth":"secret"},"user_id":"forged"}`, cookie, 200)
	// Deleting an account also revokes its sessions.
	exec(`DELETE FROM users WHERE id=$1`, outsider)
	request("GET", "/me", "", otherCookie, 401)
}
