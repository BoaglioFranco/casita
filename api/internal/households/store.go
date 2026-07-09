package households

import (
	"context"
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
