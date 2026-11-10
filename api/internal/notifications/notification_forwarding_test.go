package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestSubscribeHandlerForwardsAuthenticatedIdentityAndPayload(t *testing.T) {
	const secret = "test-only-signing-secret-at-least-32-bytes"
	t.Setenv("S2S_TOKEN", secret)
	userID := uuid.New()

	for _, upstreamStatus := range []int{http.StatusOK, http.StatusUnauthorized, http.StatusInternalServerError} {
		t.Run(http.StatusText(upstreamStatus), func(t *testing.T) {
			called := false
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				if r.Method != http.MethodPost || r.URL.Path != "/notifications/subscribe" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing JSON content type")
				}
				claims := &jwt.RegisteredClaims{}
				_, err := jwt.ParseWithClaims(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), claims,
					func(_ *jwt.Token) (any, error) { return []byte(secret), nil },
					jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired(),
					jwt.WithIssuer("groceries-api"), jwt.WithAudience("notifications"))
				if err != nil || claims.Subject != userID.String() {
					t.Errorf("incorrect authenticated identity: subject=%s, error=%v", claims.Subject, err)
				}
				if claims.IssuedAt == nil || claims.ExpiresAt == nil {
					t.Error("missing token timestamps")
				} else if claims.ExpiresAt.Sub(claims.IssuedAt.Time) != 5*time.Minute {
					t.Error("token should expire five minutes after issuance")
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Error(err)
				}
				keys, ok := payload["keys"].(map[string]any)
				if !ok || payload["endpoint"] != "https://push.example.test/device" || keys["p256dh"] != "public-key" || keys["auth"] != "auth-secret" {
					t.Error("subscription JSON shape or values changed")
				}
				if _, exists := payload["user_id"]; exists {
					t.Error("identity should only be conveyed through the JWT")
				}
				w.WriteHeader(upstreamStatus)
			}))
			defer upstream.Close()
			service := NewNotificationService(nil, upstream.URL)
			router := gin.New()
			// Stand in for the existing session middleware's verified identity.
			router.Use(func(c *gin.Context) { c.Set("user_id", userID); c.Next() })
			service.SetupNotificationRoutes(router)
			body := `{"endpoint":"https://push.example.test/device","keys":{"p256dh":"public-key","auth":"auth-secret"},"user_id":"untrusted"}`
			req := httptest.NewRequest(http.MethodPost, "/notifications/subscribe", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", uuid.NewString())
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			want := http.StatusOK
			if upstreamStatus != http.StatusOK {
				want = http.StatusBadGateway
			}
			if !called || response.Code != want {
				t.Fatalf("upstream called=%v status=%d, want=%d", called, response.Code, want)
			}

			called = false
			bad := httptest.NewRequest(http.MethodPost, "/notifications/subscribe", strings.NewReader(`{"endpoint":"missing-keys"}`))
			bad.Header.Set("Content-Type", "application/json")
			response = httptest.NewRecorder()
			router.ServeHTTP(response, bad)
			if called || response.Code != http.StatusBadRequest {
				t.Fatal("invalid payload should be rejected before forwarding")
			}
		})
	}
}

func TestSubscribePropagatesRequestCancellation(t *testing.T) {
	t.Setenv("S2S_TOKEN", "test-only-signing-secret-at-least-32-bytes")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	service := NewNotificationService(nil, "http://127.0.0.1:1")
	err := service.subscribe(ctx, uuid.New(), SubscriptionPayload{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}
