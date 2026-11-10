package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/google/uuid"
)

type AuthStore struct {
	db *database.Database
}

func NewAuthStore(db *database.Database) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (store *AuthStore) CreateSession(ctx context.Context, params CreateSessionParams) error {
	query := `
		INSERT INTO sessions (token_hash, user_id, expires_at)
		    VALUES ($1, $2, $3)`
	_, err := store.db.ExecContext(ctx, query, params.TokenHash, params.UserID, params.ExpiresAt)

	if err != nil {
		return err
	}
	return nil
}

func (store *AuthStore) GetSession(ctx context.Context, tokenHash string) (*Session, error) {
	query := `
		SELECT
		    s.token_hash,
		    s.user_id,
		    s.expires_at,
		    s.created_at,
		    u.name,
		    u.household_id
		FROM
		    sessions s
		    JOIN users u ON u.id = s.user_id
		WHERE
		    s.token_hash = $1`
	var session Session
	if err := store.db.GetContext(ctx, &session, query, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (store *AuthStore) UserIDFromEmail(ctx context.Context, email string) (uuid.UUID, error) {
	query := `
		SELECT
		    id
		FROM
		    users
		WHERE
		    email = $1`
	var userID uuid.UUID
	if err := store.db.GetContext(ctx, &userID, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrUserNotFound
		}

		return uuid.Nil, fmt.Errorf("get user ID by email: %w", err)
	}

	return userID, nil
}

func (store *AuthStore) UserFromEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT
		    *
		FROM
		    users
		WHERE
		    email = $1`
	var user User
	if err := store.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (store *AuthStore) InsertUser(ctx context.Context, payload InsertUserParams) error {
	query := `
		INSERT INTO users (email, password_hash, name, household_id)
		    VALUES ($1, $2, $3, $4)`

	_, err := store.db.ExecContext(ctx, query, payload.Email, payload.HashedPassword, payload.Name, payload.HouseholdID)
	if err != nil {
		return err
	}
	return nil
}
