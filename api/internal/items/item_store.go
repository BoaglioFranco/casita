package items

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type ItemStore struct {
	db *database.Database
}

func NewItemStore(db *database.Database) *ItemStore {
	return &ItemStore{
		db: db,
	}
}

func (store *ItemStore) InsertItem(ctx context.Context, params InsertItemParams) (int64, error) {
	return store.insertItem(ctx, store.db, params)
}

func (store *ItemStore) InsertItemTx(ctx context.Context, tx database.Executor, params InsertItemParams) (int64, error) {
	return store.insertItem(ctx, tx, params)
}

func (store *ItemStore) insertItem(ctx context.Context, exec database.Executor, params InsertItemParams) (int64, error) {
	query := `INSERT INTO items (list_id, name, description) VALUES ($1, $2, $3) RETURNING id`

	var itemID int64
	err := exec.GetContext(ctx, &itemID, query, params.ListID, params.Name, params.Description)
	if err != nil {
		return 0, fmt.Errorf("unable to insert item: %w", err)
	}

	return itemID, nil
}

func (store *ItemStore) DeleteItem(ctx context.Context, itemID int64) error {
	return store.deleteItem(ctx, store.db, itemID)
}

func (store *ItemStore) DeleteItemTx(ctx context.Context, tx database.Executor, itemID int64) error {
	return store.deleteItem(ctx, tx, itemID)
}

func (store *ItemStore) deleteItem(ctx context.Context, exec database.Executor, itemID int64) error {
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

func (store *ItemStore) GetItemTx(ctx context.Context, tx database.Executor, itemID int64) (Item, error) {
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

func (store *ItemStore) GetItemsFromList(ctx context.Context, listID int64) ([]Item, error) {
	query := `
		SELECT id, list_id, name, description, created_at, updated_at
		FROM items
		WHERE list_id = $1
		ORDER BY created_at DESC
	`

	items := make([]Item, 0)
	if err := store.db.SelectContext(ctx, &items, query, listID); err != nil {
		return nil, fmt.Errorf("unable to retrieve items: %w", err)
	}

	return items, nil
}
