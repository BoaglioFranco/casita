package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/BoaglioFranco/casita/api/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthService(db *database.Database) *AuthService {
	return &AuthService{
		authStore: NewAuthStore(db),
	}
}

func generateSessionToken() (rawToken string, tokenHash string, err error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", err
	}

	rawToken = base64.RawURLEncoding.EncodeToString(tokenBytes)
	return rawToken, hashRawSessionToken(rawToken), nil
}

func hashRawSessionToken(rawToken string) (tokenHash string) {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

func (service *AuthService) Login(ctx context.Context, email string, password string) (token string, err error) {
	user, err := service.authStore.UserFromEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	rawToken, hashedToken, err := generateSessionToken()
	if err != nil {
		return "", err
	}
	err = service.authStore.CreateSession(ctx, CreateSessionParams{TokenHash: hashedToken, UserID: user.ID, ExpiresAt: time.Now().UTC().AddDate(1, 0, 0)})
	if err != nil {
		return "", err
	}

	return rawToken, nil
}

func (service *AuthService) AuthenticateSession(ctx context.Context, token string) (*Session, error) {
	tokenHash := hashRawSessionToken(token)
	userSession, err := service.authStore.GetSession(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if userSession.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionExpired
	}
	return userSession, nil
}
