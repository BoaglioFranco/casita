package items

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type ItemStore struct {
	db *database.Database
}

func NewStore(db *database.Database) *ItemStore {
	return &ItemStore{
		db: db,
	}
}

func (s *ItemStore) InsertItem(ctx context.Context, params InsertItemParams) (int64, error) {
	return s.insertItem(ctx, s.db, params)
}

func (s *ItemStore) InsertItemTx(ctx context.Context, tx database.Executor, params InsertItemParams) (int64, error) {
	return s.insertItem(ctx, tx, params)
}

func (s *ItemStore) insertItem(ctx context.Context, exec database.Executor, params InsertItemParams) (int64, error) {
	query := `INSERT INTO items (list_id, name, description) VALUES ($1, $2, $3) RETURNING id`

	var itemID int64
	err := exec.GetContext(ctx, &itemID, query, params.ListID, params.Name, params.Description)
	if err != nil {
		return 0, fmt.Errorf("unable to insert item: %w", err)
	}

	return itemID, nil
}

func (s *ItemStore) DeleteItem(ctx context.Context, itemID int64) error {
	return s.deleteItem(ctx, s.db, itemID)
}

func (s *ItemStore) DeleteItemTx(ctx context.Context, tx database.Executor, itemID int64) error {
	return s.deleteItem(ctx, tx, itemID)
}

func (s *ItemStore) deleteItem(ctx context.Context, exec database.Executor, itemID int64) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := exec.ExecContext(ctx, query, itemID)
	if err != nil {
		return fmt.Errorf("unable to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to check deleted item rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrItemNotFound
	}

	return nil
}

func (s *ItemStore) GetItemTx(ctx context.Context, tx database.Executor, itemID int64) (Item, error) {
	query := `
		SELECT id, list_id, name, description, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	items := make([]Item, 0, 1)
	if err := tx.SelectContext(ctx, &items, query, itemID); err != nil {
		return Item{}, fmt.Errorf("unable to retrieve item: %w", err)
	}
	if len(items) == 0 {
		return Item{}, ErrItemNotFound
	}

	return items[0], nil
}

func (s *ItemStore) GetItemsFromList(ctx context.Context, listID int64) ([]Item, error) {
	query := `
		SELECT id, list_id, name, description, created_at, updated_at
		FROM items
		WHERE list_id = $1
		ORDER BY created_at DESC
	`

	items := make([]Item, 0)
	if err := s.db.SelectContext(ctx, &items, query, listID); err != nil {
		return nil, fmt.Errorf("unable to retrieve items: %w", err)
	}

	return items, nil
}
