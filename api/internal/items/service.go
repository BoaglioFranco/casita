package items

import (
	"context"
	"fmt"
	"log"

	"github.com/BoaglioFranco/groceriesList/api/internal/activity"
	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/notifications"
)

func NewItemService(db *database.Database, notificationService *notifications.NotificationService) *ItemService {
	return &ItemService{
		db:                  db,
		store:               NewStore(db),
		activityStore:       activity.NewStore(db),
		notificationService: notificationService,
	}
}

func (s *ItemService) CreateItem(ctx context.Context, params CreateItemParams) error {
	err := s.db.WithTx(ctx, func(tx database.Executor) error {
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
	if err != nil {
		return err
	}

	err = s.notificationService.NotifyMembers(ctx, params.HouseholdID, notifications.Notification{
		Title: "Item created",
		Body:  fmt.Sprintf("%s added to the list", params.Name),
	})
	if err != nil {
		log.Printf("notify household: %v", err)
	}

	return nil
}

func (s *ItemService) DeleteItem(ctx context.Context, params DeleteItemParams) error {
	var itemName string
	err := s.db.WithTx(ctx, func(tx database.Executor) error {
		item, err := s.store.GetItemTx(ctx, tx, params.ItemID)
		if err != nil {
			return err
		}

		if item.ListID != params.ListID {
			return ErrItemNotFound
		}
		itemName = item.Name

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

	if err != nil {
		return err
	}
	err = s.notificationService.NotifyMembers(ctx, params.HouseholdID, notifications.Notification{
		Title: "Item deleted",
		Body:  fmt.Sprintf("%s removed from the list", itemName),
	})
	if err != nil {
		log.Printf("notify household: %v", err)
	}
	return nil
}
