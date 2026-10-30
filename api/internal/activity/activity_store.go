package activity

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type ActivityStore struct {
	db *database.Database
}

func NewActivityStore(db *database.Database) *ActivityStore {
	return &ActivityStore{
		db: db,
	}
}

func (store *ActivityStore) RecordActivity(ctx context.Context, params RecordActivityParams) error {
	return store.recordActivity(ctx, store.db, params)
}

func (store *ActivityStore) RecordActivityTx(ctx context.Context, tx database.Executor, params RecordActivityParams) error {
	return store.recordActivity(ctx, tx, params)
}

func (store *ActivityStore) recordActivity(ctx context.Context, exec database.Executor, params RecordActivityParams) error {
	query := `
		INSERT INTO activity_events (event_type, household_id, list_id, actor_member_id, item_id, item_name_snapshot)
		    VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := exec.ExecContext(ctx, query, params.EventType, params.HouseholdID, params.ListID, params.ActorMemberID, params.ItemID, params.ItemNameSnapshot)
	if err != nil {
		return fmt.Errorf("failed to record activity %w", err)
	}

	return nil
}
