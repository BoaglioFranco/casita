package auth

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrSessionExpired = errors.New("session expired")

type AuthService struct {
	store *AuthStore
}

type CreateSessionParams struct {
	TokenHash string
	UserID    int64
	ExpiresAt time.Time
}

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type InsertUserParams struct {
	Email          string
	HashedPassword string
}
type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
type Session struct {
	TokenHash string    `db:"token_hash"`
	UserID    int64     `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
