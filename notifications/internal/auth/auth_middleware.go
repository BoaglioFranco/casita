package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	Issuer   = "casita-api"
	Audience = "notifications"
)

// NewMiddleware verifies user assertions signed by the main backend.
func NewMiddleware(secret string) (gin.HandlerFunc, error) {
	if len(strings.TrimSpace(secret)) < 32 {
		return nil, errors.New("S2S_TOKEN must contain at least 32 bytes")
	}
	key := []byte(secret)
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(Issuer),
		jwt.WithAudience(Audience),
	)

	return func(c *gin.Context) {
		unauthorized := func() {
			c.Header("WWW-Authenticate", "Bearer")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		}

		headers := c.Request.Header.Values("Authorization")
		if len(headers) != 1 {
			unauthorized()
			return
		}
		parts := strings.Fields(headers[0])
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized()
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := parser.ParseWithClaims(parts[1], claims, func(_ *jwt.Token) (any, error) {
			return key, nil
		})
		if err != nil || !token.Valid {
			unauthorized()
			return
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil || userID == uuid.Nil {
			unauthorized()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}, nil
}
