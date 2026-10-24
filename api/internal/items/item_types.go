package items

import (
	"errors"
	"time"

	"github.com/BoaglioFranco/groceriesList/api/internal/activity"
	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/notifications"
)

var ErrItemNotFound = errors.New("item not found")

type ItemService struct {
	db                  *database.Database
	itemStore           *ItemStore
	activityStore       *activity.ActivityStore
	notificationService *notifications.NotificationService
}

type Item struct {
	ID          int64     `db:"id" json:"id"`
	ListID      int64     `db:"list_id" json:"listId"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// Handlers

type CreateItemPayload struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type CreateItemParams struct {
	HouseholdID   int64
	ListID        int64
	ActorMemberID int64
	Name          string
	Description   *string
}

type DeleteItemParams struct {
	HouseholdID   int64
	ListID        int64
	ItemID        int64
	ActorMemberID int64
}

// Store
type InsertItemParams struct {
	ListID      int64
	Name        string
	Description *string
}
