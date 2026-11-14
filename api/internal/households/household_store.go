package households

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/casita/api/internal/database"
	"github.com/google/uuid"
)

type HouseholdStore struct {
	db *database.Database
}

func NewHouseholdStore(db *database.Database) *HouseholdStore {
	return &HouseholdStore{db: db}
}

func (store *HouseholdStore) UserIDsByHouseholdID(ctx context.Context, householdID uuid.UUID) ([]uuid.UUID, error) {
	users := make([]uuid.UUID, 0)
	if err := store.db.SelectContext(ctx, &users, `SELECT id FROM users WHERE household_id = $1`, householdID); err != nil {
		return nil, fmt.Errorf("retrieve household users: %w", err)
	}
	return users, nil
}
