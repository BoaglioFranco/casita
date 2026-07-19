package households

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type HouseholdStore struct {
	db *database.Database
}

func NewStore(db *database.Database) *HouseholdStore {
	return &HouseholdStore{
		db: db,
	}
}

func (s *HouseholdStore) InsertHousehold(ctx context.Context, name string) error {
	query := `INSERT INTO households (name) VALUES (?)`

	_, err := s.db.ExecContext(ctx, query, name)

	if err != nil {
		return fmt.Errorf("unable to insert household: %w", err)
	}
	return nil
}

func (s *HouseholdStore) InsertMember(ctx context.Context, params InsertMemberParams) error {
	query := `INSERT INTO members (household_id, name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, query, params.HouseholdID, params.Name)
	if err != nil {
		return fmt.Errorf("unable to insert member: %w", err)
	}

	return nil
}

func (s *HouseholdStore) GetMember(ctx context.Context, memberID int64) (*Member, error) {
	query := `
		SELECT id, household_id, name
		FROM members
		WHERE id = ?
	`

	members := make([]Member, 0, 1)
	if err := s.db.SelectContext(ctx, &members, query, memberID); err != nil {
		return nil, fmt.Errorf("unable to retrieve member: %w", err)
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return &members[0], nil
}

func (s *HouseholdStore) MemberByUserID(ctx context.Context, userID int64, househouldID int64) (*Member, error) {
	query := `SELECT id, household_id, name FROM members WHERE user_id = ? AND household_id = ?`
	var member Member
	if err := s.db.GetContext(ctx, &member, query, userID, househouldID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}

		return nil, fmt.Errorf("unable to retrieve member: %w", err)
	}

	return &member, nil
}
