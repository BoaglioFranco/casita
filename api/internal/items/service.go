package items

import (
	"context"

	"github.com/BoaglioFranco/groceriesList/api/internal/activity"
	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

func NewItemService(db *database.Database) *ItemService {
	return &ItemService{
		db:            db,
		store:         NewStore(db),
		activityStore: activity.NewStore(db),
	}
}

func (s *ItemService) CreateItem(ctx context.Context, params CreateItemParams) error {
	return s.db.WithTx(ctx, func(tx database.Executor) error {
		itemID, err := s.store.InsertItemTx(ctx, tx, InsertItemParams{
			ListID:      params.ListID,
			Name:        params.Name,
			Description: params.Description,
		})
		if err != nil {
			return err
		}

		return s.activityStore.RecordActivityTx(ctx, tx, activity.RecordActivityParams{
			EventType:        activity.ActivityItemCreate,
			ActorMemberID:    &params.ActorMemberID,
			ListID:           &params.ListID,
			HouseholdID:      params.HouseholdID,
			ItemID:           &itemID,
			ItemNameSnapshot: &params.Name,
		})
	})
}

func (s *ItemService) DeleteItem(ctx context.Context, params DeleteItemParams) error {
	return s.db.WithTx(ctx, func(tx database.Executor) error {
		item, err := s.store.GetItemTx(ctx, tx, params.ItemID)
		if err != nil {
			return err
		}

		if item.ListID != params.ListID {
			return ErrItemNotFound
		}

		if err := s.store.DeleteItemTx(ctx, tx, params.ItemID); err != nil {
			return err
		}

		return s.activityStore.RecordActivityTx(ctx, tx, activity.RecordActivityParams{
			EventType:        activity.ActivityItemDelete,
			ActorMemberID:    &params.ActorMemberID,
			ListID:           &params.ListID,
			HouseholdID:      params.HouseholdID,
			ItemNameSnapshot: &item.Name,
		})
	})
}
