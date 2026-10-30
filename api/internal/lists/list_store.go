package lists

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/google/uuid"
)

type ListStore struct {
	db *database.Database
}

func NewListStore(db *database.Database) *ListStore {
	return &ListStore{
		db: db,
	}
}

func (store *ListStore) InsertList(ctx context.Context, params InsertListParams) error {
	query := `
		INSERT INTO lists (household_id, name)
		    VALUES ($1, $2)`

	_, err := store.db.ExecContext(ctx, query, params.HouseholdID, params.Name)
	if err != nil {
		return fmt.Errorf("unable to insert list: %w", err)
	}

	return nil
}

func (store *ListStore) GetLists(ctx context.Context, householdID uuid.UUID) ([]List, error) {
	query := `
		SELECT
		    name,
		    icon,
		    household_id,
		    id,
		    created_at
		FROM
		    lists
		WHERE
		    household_id = $1
		ORDER BY
		    created_at ASC`

	var lists = make([]List, 0, 1)
	err := store.db.SelectContext(ctx, &lists, query, householdID)
	if err != nil {
		return nil, err
	}
	return lists, nil
}
