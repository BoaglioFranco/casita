package items

import (
	"context"
	"fmt"
	"log"

	"github.com/BoaglioFranco/casita/api/internal/activity"
	"github.com/BoaglioFranco/casita/api/internal/database"
	"github.com/BoaglioFranco/casita/api/internal/notifications"
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
	var listName string
	err := service.db.WithTx(ctx, func(tx database.Executor) error {
		var err error
		listName, err = service.itemStore.GetListNameTx(ctx, tx, params.User.HouseholdID, params.ListID)
		if err != nil {
			return err
		}
		itemID, err := service.itemStore.InsertItemTx(ctx, tx, InsertItemParams{
			ListID:          params.ListID,
			HouseholdID:     params.User.HouseholdID,
			Name:            params.Name,
			Description:     params.Description,
			CreatedByUserID: params.User.ID,
		})
		if err != nil {
			return err
		}

		return service.activityStore.RecordActivityTx(ctx, tx, activity.RecordActivityParams{
			EventType:        activity.ActivityItemCreate,
			ActorUserID:      &params.User.ID,
			ListID:           &params.ListID,
			HouseholdID:      params.User.HouseholdID,
			ItemID:           &itemID,
			ItemNameSnapshot: &params.Name,
		})
	})
	if err != nil {
		return err
	}

	err = service.notificationService.NotifyOtherUsers(ctx, params.User.HouseholdID, params.User.ID, notifications.Notification{
		Title: "Item created",
		Body:  fmt.Sprintf("%s added to %s by %s", params.Name, listName, params.User.Name),
	})
	if err != nil {
		log.Printf("notify household: %v", err)
	}

	return nil
}

func (service *ItemService) DeleteItem(ctx context.Context, params DeleteItemParams) error {
	var itemName, listName string
	err := service.db.WithTx(ctx, func(tx database.Executor) error {
		var err error
		listName, err = service.itemStore.GetListNameTx(ctx, tx, params.User.HouseholdID, params.ListID)
		if err != nil {
			return err
		}
		item, err := service.itemStore.GetItemTx(ctx, tx, params.User.HouseholdID, params.ListID, params.ItemID)
		if err != nil {
			return err
		}
		itemName = item.Name

		if err := service.itemStore.DeleteItemTx(ctx, tx, params.User.HouseholdID, params.ListID, params.ItemID); err != nil {
			return err
		}

		return service.activityStore.RecordActivityTx(ctx, tx, activity.RecordActivityParams{
			EventType:        activity.ActivityItemDelete,
			ActorUserID:      &params.User.ID,
			ListID:           &params.ListID,
			HouseholdID:      params.User.HouseholdID,
			ItemNameSnapshot: &item.Name,
		})
	})

	if err != nil {
		return err
	}
	err = service.notificationService.NotifyOtherUsers(ctx, params.User.HouseholdID, params.User.ID, notifications.Notification{
		Title: "Item deleted",
		Body:  fmt.Sprintf("%s removed from %s by %s", itemName, listName, params.User.Name),
	})
	if err != nil {
		log.Printf("notify household: %v", err)
	}
	return nil
}
