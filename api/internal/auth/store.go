package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type AuthStore struct {
	db *database.Database
}

func NewAuthStore(db *database.Database) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (s *AuthStore) CreateSession(ctx context.Context, params CreateSessionParams) error {
	query := `INSERT INTO sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, params.TokenHash, params.UserID, params.ExpiresAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *AuthStore) GetSession(ctx context.Context, tokenHash string) (*Session, error) {
	query := `SELECT token_hash, user_id, expires_at, created_at FROM sessions WHERE token_hash = ?`
	var session Session
	if err := s.db.GetContext(ctx, &session, query, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (s *AuthStore) UserIDFromEmail(ctx context.Context, email string) (int64, error) {
	query := `SELECT id FROM users WHERE email = ?`
	var userID int64
	if err := s.db.GetContext(ctx, &userID, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}

		return 0, fmt.Errorf("get user ID by email: %w", err)
	}

	return userID, nil
}

func (s *AuthStore) UserFromEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = ?`
	var user User
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (s *AuthStore) InsertUser(ctx context.Context, payload InsertUserParams) error {
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, query, payload.Email, payload.HashedPassword)
	if err != nil {
		return err
	}
	return nil
}
