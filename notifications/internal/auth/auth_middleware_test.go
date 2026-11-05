package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestNewMiddlewareRejectsInvalidConfiguration(t *testing.T) {
	for _, secret := range []string{"", "short", strings.Repeat(" ", 64)} {
		if _, err := NewMiddleware(secret); err == nil {
			t.Fatal("expected invalid signing secret to be rejected")
		}
	}
}

func TestAuthentication(t *testing.T) {
	const secret = "test-only-signing-secret-at-least-32-bytes"
	userID := uuid.New()
	now := time.Now()
	validClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    Issuer,
		Audience:  jwt.ClaimStrings{Audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
	}

	tests := []struct {
		name      string
		change    func(*jwt.RegisteredClaims)
		method    jwt.SigningMethod
		key       any
		header    string
		duplicate bool
		want      int
	}{
		{name: "valid", want: http.StatusNoContent},
		{name: "case insensitive bearer", header: "bearer %s", want: http.StatusNoContent},
		{name: "audience list", change: func(c *jwt.RegisteredClaims) { c.Audience = jwt.ClaimStrings{"another", Audience} }, want: http.StatusNoContent},
		{name: "missing header", header: "missing"},
		{name: "missing bearer token", header: "Bearer"},
		{name: "wrong scheme", header: "Basic %s"},
		{name: "extra header fields", header: "Bearer %s extra"},
		{name: "duplicate headers", duplicate: true},
		{name: "malformed jwt", header: "Bearer invalid"},
		{name: "wrong signature", key: []byte("different-signing-secret-at-least-32-bytes")},
		{name: "wrong algorithm", method: jwt.SigningMethodHS384},
		{name: "unsigned", method: jwt.SigningMethodNone, key: jwt.UnsafeAllowNoneSignatureType},
		{name: "expired", change: func(c *jwt.RegisteredClaims) { c.ExpiresAt = jwt.NewNumericDate(now.Add(-time.Minute)) }},
		{name: "missing expiration", change: func(c *jwt.RegisteredClaims) { c.ExpiresAt = nil }},
		{name: "future not before", change: func(c *jwt.RegisteredClaims) { c.NotBefore = jwt.NewNumericDate(now.Add(time.Hour)) }},
		{name: "wrong issuer", change: func(c *jwt.RegisteredClaims) { c.Issuer = "other" }},
		{name: "missing issuer", change: func(c *jwt.RegisteredClaims) { c.Issuer = "" }},
		{name: "wrong audience", change: func(c *jwt.RegisteredClaims) { c.Audience = jwt.ClaimStrings{"other"} }},
		{name: "missing audience", change: func(c *jwt.RegisteredClaims) { c.Audience = nil }},
		{name: "missing subject", change: func(c *jwt.RegisteredClaims) { c.Subject = "" }},
		{name: "malformed subject", change: func(c *jwt.RegisteredClaims) { c.Subject = "not-a-uuid" }},
		{name: "nil subject", change: func(c *jwt.RegisteredClaims) { c.Subject = uuid.Nil.String() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := validClaims
			if tt.change != nil {
				tt.change(&claims)
			}
			method := tt.method
			if method == nil {
				method = jwt.SigningMethodHS256
			}
			key := tt.key
			if key == nil {
				key = []byte(secret)
			}
			token, err := jwt.NewWithClaims(method, claims).SignedString(key)
			if err != nil {
				t.Fatal(err)
			}
			middleware, err := NewMiddleware(secret)
			if err != nil {
				t.Fatal(err)
			}
			router := gin.New()
			router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
			called := false
			group := router.Group("/")
			group.Use(middleware)
			group.POST("/subscriptions", func(c *gin.Context) {
				called = true
				if got := c.MustGet("user_id"); got != userID {
					t.Errorf("user_id = %v, want %v", got, userID)
				}
				c.Status(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(`{"user_id":"untrusted"}`))
			req.Header.Set("X-User-ID", uuid.NewString())
			header := tt.header
			if header == "" {
				header = "Bearer %s"
			}
			if header != "missing" {
				req.Header.Set("Authorization", strings.ReplaceAll(header, "%s", token))
			}
			if tt.duplicate {
				req.Header.Add("Authorization", "Bearer "+token)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			want := tt.want
			if want == 0 {
				want = http.StatusUnauthorized
			}
			if response.Code != want || called != (want == http.StatusNoContent) {
				t.Fatalf("status=%d handler called=%v, want status=%d", response.Code, called, want)
			}
			if want == http.StatusUnauthorized && strings.Contains(response.Body.String(), token) {
				t.Fatal("response exposed the token")
			}
			health := httptest.NewRecorder()
			router.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
			if health.Code != http.StatusOK {
				t.Fatal("health endpoint should remain public")
			}
		})
	}
}
