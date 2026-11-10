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
		itemStore:           NewItemStore(db),
		activityStore:       activity.NewActivityStore(db),
		notificationService: notificationService,
	}
}

func (service *ItemService) CreateItem(ctx context.Context, params CreateItemParams) error {
	err := service.db.WithTx(ctx, func(tx database.Executor) error {
		itemID, err := service.itemStore.InsertItemTx(ctx, tx, InsertItemParams{
			ListID:          params.ListID,
			HouseholdID:     params.HouseholdID,
			Name:            params.Name,
			Description:     params.Description,
			CreatedByUserID: params.ActorUserID,
		})
		if err != nil {
			return err
		}

		return service.activityStore.RecordActivityTx(ctx, tx, activity.RecordActivityParams{
			EventType:        activity.ActivityItemCreate,
			ActorUserID:      &params.ActorUserID,
			ListID:           &params.ListID,
			HouseholdID:      params.HouseholdID,
			ItemID:           &itemID,
			ItemNameSnapshot: &params.Name,
		})
	})
	if err != nil {
		return err
	}

	err = service.notificationService.NotifyOtherUsers(ctx, params.HouseholdID, params.ActorUserID, notifications.Notification{
		Title: "Item created",
		Body:  fmt.Sprintf("%s added to the list", params.Name),
	})
	if err != nil {
		log.Printf("notify household: %v", err)
	}

	return nil
}

func (service *ItemService) DeleteItem(ctx context.Context, params DeleteItemParams) error {
	var itemName string
	err := service.db.WithTx(ctx, func(tx database.Executor) error {
		item, err := service.itemStore.GetItemTx(ctx, tx, params.HouseholdID, params.ListID, params.ItemID)
		if err != nil {
			return err
		}
		itemName = item.Name

		if err := service.itemStore.DeleteItemTx(ctx, tx, params.HouseholdID, params.ListID, params.ItemID); err != nil {
			return err
		}

		return service.activityStore.RecordActivityTx(ctx, tx, activity.RecordActivityParams{
			EventType:        activity.ActivityItemDelete,
			ActorUserID:      &params.ActorUserID,
			ListID:           &params.ListID,
			HouseholdID:      params.HouseholdID,
			ItemNameSnapshot: &item.Name,
		})
	})

	if err != nil {
		return err
	}
	err = service.notificationService.NotifyOtherUsers(ctx, params.HouseholdID, params.ActorUserID, notifications.Notification{
		Title: "Item deleted",
		Body:  fmt.Sprintf("%s removed from the list", itemName),
	})
	if err != nil {
		log.Printf("notify household: %v", err)
	}
	return nil
}
