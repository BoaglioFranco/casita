package items

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/google/uuid"
)

type ItemStore struct {
	db *database.Database
}

func NewItemStore(db *database.Database) *ItemStore {
	return &ItemStore{
		db: db,
	}
}

func (store *ItemStore) InsertItem(ctx context.Context, params InsertItemParams) (uuid.UUID, error) {
	return store.insertItem(ctx, store.db, params)
}

func (store *ItemStore) InsertItemTx(ctx context.Context, tx database.Executor, params InsertItemParams) (uuid.UUID, error) {
	return store.insertItem(ctx, tx, params)
}

func (store *ItemStore) insertItem(ctx context.Context, exec database.Executor, params InsertItemParams) (uuid.UUID, error) {
	query := `
		INSERT INTO items (list_id, name, description, created_by_member_id)
		    VALUES ($1, $2, $3, $4)
		RETURNING
		    id`

	var itemID uuid.UUID
	err := exec.GetContext(ctx, &itemID, query, params.ListID, params.Name, params.Description, params.CreatedByMemberID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to insert item: %w", err)
	}

	return itemID, nil
}

func (store *ItemStore) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	return store.deleteItem(ctx, store.db, itemID)
}

func (store *ItemStore) DeleteItemTx(ctx context.Context, tx database.Executor, itemID uuid.UUID) error {
	return store.deleteItem(ctx, tx, itemID)
}

func (store *ItemStore) deleteItem(ctx context.Context, exec database.Executor, itemID uuid.UUID) error {
	query := `
		DELETE FROM items
		WHERE id = $1`

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

func (store *ItemStore) GetItemTx(ctx context.Context, tx database.Executor, itemID uuid.UUID) (Item, error) {
	query := `
		SELECT
		    id,
		    list_id,
		    name,
		    description,
		    created_at,
		    updated_at
		FROM
		    items
		WHERE
		    id = $1`

	rows := make([]itemRow, 0, 1)
	if err := tx.SelectContext(ctx, &rows, query, itemID); err != nil {
		return Item{}, fmt.Errorf("unable to retrieve item: %w", err)
	}
	if len(rows) == 0 {
		return Item{}, ErrItemNotFound
	}

	return itemFromRow(rows[0]), nil
}

func (store *ItemStore) GetItemsFromList(ctx context.Context, listID uuid.UUID) ([]Item, error) {
	query := `
		SELECT
		    i.id,
		    i.list_id,
		    i.name,
		    i.description,
		    i.created_by_member_id,
		    m.name AS created_by_member_name,
		    i.created_at,
		    i.updated_at
		FROM
		    items i
		    LEFT JOIN members m ON m.id = i.created_by_member_id
		WHERE
		    i.list_id = $1
		ORDER BY
		    i.created_at DESC`

	rows := make([]itemRow, 0)
	if err := store.db.SelectContext(ctx, &rows, query, listID); err != nil {
		return nil, fmt.Errorf("unable to retrieve items: %w", err)
	}

	items := make([]Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, itemFromRow(row))
	}

	return items, nil
}

func itemFromRow(row itemRow) Item {
	item := Item{
		ID:          row.ID,
		ListID:      row.ListID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if row.CreatedByMemberID.Valid && row.CreatedByMemberName != nil {
		item.CreatedBy = &ItemCreator{
			ID:   row.CreatedByMemberID.UUID,
			Name: *row.CreatedByMemberName,
		}
	}

	return item
}
