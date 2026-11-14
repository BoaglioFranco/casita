package items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BoaglioFranco/casita/api/internal/database"
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

func (store *ItemStore) GetListNameTx(ctx context.Context, tx database.Executor, householdID, listID uuid.UUID) (string, error) {
	var name string
	err := tx.GetContext(ctx, &name, `SELECT name FROM lists WHERE id = $1 AND household_id = $2`, listID, householdID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrItemNotFound
	}
	if err != nil {
		return "", fmt.Errorf("unable to retrieve list name: %w", err)
	}
	return name, nil
}

func (store *ItemStore) InsertItem(ctx context.Context, params InsertItemParams) (uuid.UUID, error) {
	return store.insertItem(ctx, store.db, params)
}

func (store *ItemStore) InsertItemTx(ctx context.Context, tx database.Executor, params InsertItemParams) (uuid.UUID, error) {
	return store.insertItem(ctx, tx, params)
}

func (store *ItemStore) insertItem(ctx context.Context, exec database.Executor, params InsertItemParams) (uuid.UUID, error) {
	query := `
		INSERT INTO items (list_id, name, description, created_by_user_id)
		SELECT id, $2, $3, $4 FROM lists WHERE id = $1 AND household_id = $5
		RETURNING
		    id`

	var itemID uuid.UUID
	err := exec.GetContext(ctx, &itemID, query, params.ListID, params.Name, params.Description, params.CreatedByUserID, params.HouseholdID)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, ErrItemNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to insert item: %w", err)
	}

	return itemID, nil
}

func (store *ItemStore) DeleteItem(ctx context.Context, householdID, listID, itemID uuid.UUID) error {
	return store.deleteItem(ctx, store.db, householdID, listID, itemID)
}

func (store *ItemStore) DeleteItemTx(ctx context.Context, tx database.Executor, householdID, listID, itemID uuid.UUID) error {
	return store.deleteItem(ctx, tx, householdID, listID, itemID)
}

func (store *ItemStore) deleteItem(ctx context.Context, exec database.Executor, householdID, listID, itemID uuid.UUID) error {
	query := `
		DELETE FROM items i USING lists l
		WHERE i.id = $1 AND i.list_id = l.id AND l.id = $2 AND l.household_id = $3`

	result, err := exec.ExecContext(ctx, query, itemID, listID, householdID)
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

func (store *ItemStore) GetItemTx(ctx context.Context, tx database.Executor, householdID, listID, itemID uuid.UUID) (Item, error) {
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
		WHERE id = $1 AND list_id = $2
		    AND EXISTS (SELECT 1 FROM lists WHERE id = $2 AND household_id = $3)
		FOR UPDATE`

	rows := make([]itemRow, 0, 1)
	if err := tx.SelectContext(ctx, &rows, query, itemID, listID, householdID); err != nil {
		return Item{}, fmt.Errorf("unable to retrieve item: %w", err)
	}
	if len(rows) == 0 {
		return Item{}, ErrItemNotFound
	}

	return itemFromRow(rows[0]), nil
}

func (store *ItemStore) GetItemsFromList(ctx context.Context, householdID, listID uuid.UUID) ([]Item, error) {
	// Distinguish an empty authorized list from an absent or foreign list.
	var exists bool
	if err := store.db.GetContext(ctx, &exists, `SELECT EXISTS (SELECT 1 FROM lists WHERE id = $1 AND household_id = $2)`, listID, householdID); err != nil {
		return nil, fmt.Errorf("check list ownership: %w", err)
	}
	if !exists {
		return nil, ErrItemNotFound
	}
	query := `
		SELECT
		    i.id,
		    i.list_id,
		    i.name,
		    i.description,
		    i.created_by_user_id,
		    u.name AS created_by_user_name,
		    i.created_at,
		    i.updated_at
		FROM
		    items i
		    JOIN lists l ON l.id = i.list_id
		    LEFT JOIN users u ON u.id = i.created_by_user_id
		WHERE
		    i.list_id = $1 AND l.household_id = $2
		ORDER BY
		    i.created_at DESC`

	rows := make([]itemRow, 0)
	if err := store.db.SelectContext(ctx, &rows, query, listID, householdID); err != nil {
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
	if row.CreatedByUserID.Valid && row.CreatedByUserName != nil {
		item.CreatedBy = &ItemCreator{
			ID:   row.CreatedByUserID.UUID,
			Name: *row.CreatedByUserName,
		}
	}

	return item
}
