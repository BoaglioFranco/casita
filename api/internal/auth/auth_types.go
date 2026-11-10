package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrSessionExpired = errors.New("session expired")

type AuthService struct {
	authStore *AuthStore
}

type CreateSessionParams struct {
	TokenHash string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type InsertUserParams struct {
	Email          string
	HashedPassword string
	Name           string
	HouseholdID    uuid.UUID
}
type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	HouseholdID  uuid.UUID `db:"household_id"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
type Session struct {
	TokenHash   string    `db:"token_hash"`
	UserID      uuid.UUID `db:"user_id"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
	Name        string    `db:"name"`
	HouseholdID uuid.UUID `db:"household_id"`
}

// CurrentUser is safe to return to clients; authentication credentials stay private.
type CurrentUser struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	HouseholdID uuid.UUID `json:"householdId"`
}
