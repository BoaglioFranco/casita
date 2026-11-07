package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MintS2sJWT(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("user ID is required")
	}
	secret := os.Getenv("S2S_TOKEN")
	if secret == "" {
		return "", errors.New("S2S token not found")
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    "groceries-api",
		Audience:  jwt.ClaimStrings{"notifications"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
