package households

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/google/uuid"
)

type HouseholdStore struct {
	db *database.Database
}

func NewHouseholdStore(db *database.Database) *HouseholdStore {
	return &HouseholdStore{
		db: db,
	}
}

func (store *HouseholdStore) InsertHousehold(ctx context.Context, name string) error {
	query := `
		INSERT INTO households (name)
		    VALUES ($1)`

	_, err := store.db.ExecContext(ctx, query, name)

	if err != nil {
		return fmt.Errorf("unable to insert household: %w", err)
	}
	return nil
}

func (store *HouseholdStore) InsertMember(ctx context.Context, params InsertMemberParams) error {
	query := `
		INSERT INTO members (household_id, name, user_id)
		    VALUES ($1, $2, $3)`

	_, err := store.db.ExecContext(ctx, query, params.HouseholdID, params.Name, params.UserID)
	if err != nil {
		return fmt.Errorf("unable to insert member: %w", err)
	}

	return nil
}

func (store *HouseholdStore) GetMember(ctx context.Context, memberID uuid.UUID) (*Member, error) {

	query := `
		SELECT
		    id,
		    household_id,
		    name,
		    user_id
		FROM
		    members
		WHERE
		    id = $1`

	members := make([]Member, 0, 1)
	if err := store.db.SelectContext(ctx, &members, query, memberID); err != nil {
		return nil, fmt.Errorf("unable to retrieve member: %w", err)
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return &members[0], nil
}

func (store *HouseholdStore) MemberByUserID(ctx context.Context, userID uuid.UUID, householdID uuid.UUID) (*Member, error) {
	query := `
		SELECT
		    id,
		    household_id,
		    name,
		    user_id
		FROM
		    members
		WHERE
		    user_id = $1
		    AND household_id = $2`
	var member Member
	if err := store.db.GetContext(ctx, &member, query, userID, householdID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}

		return nil, fmt.Errorf("unable to retrieve member: %w", err)
	}

	return &member, nil
}

func (store *HouseholdStore) MembersByHouseholdID(ctx context.Context, householdID uuid.UUID) ([]Member, error) {
	query := `
		SELECT
		    id,
		    household_id,
		    name,
		    user_id
		FROM
		    members
		WHERE
		    household_id = $1`

	members := make([]Member, 0)
	if err := store.db.SelectContext(ctx, &members, query, householdID); err != nil {
		return nil, fmt.Errorf("unable to retrieve household members: %w", err)
	}
	return members, nil
}

func (store *HouseholdStore) HouseholdsByUserID(ctx context.Context, userID uuid.UUID) ([]Household, error) {
	query := `
		SELECT
		    h.id,
		    h.created_at,
		    h.name
		FROM
		    households h
		    INNER JOIN members m ON m.household_id = h.id
		WHERE
		    m.user_id = $1
		ORDER BY
		    h.created_at ASC`

	households := make([]Household, 0, 1)
	err := store.db.SelectContext(ctx, &households, query, userID)
	if err != nil {
		return nil, err
	}
	return households, nil

}
