package lists

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
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
	query := `INSERT INTO lists (household_id, name) VALUES ($1, $2)`

	_, err := store.db.ExecContext(ctx, query, params.HouseholdID, params.Name)
	if err != nil {
		return fmt.Errorf("unable to insert list: %w", err)
	}

	return nil
}
