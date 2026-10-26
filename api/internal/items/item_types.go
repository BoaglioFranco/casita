package items

import (
	"errors"
	"time"

	"github.com/BoaglioFranco/groceriesList/api/internal/activity"
	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/notifications"
	"github.com/google/uuid"
)

var ErrItemNotFound = errors.New("item not found")

type ItemService struct {
	db                  *database.Database
	itemStore           *ItemStore
	activityStore       *activity.ActivityStore
	notificationService *notifications.NotificationService
}

type Item struct {
	ID          uuid.UUID    `json:"id"`
	ListID      uuid.UUID    `json:"listId"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	CreatedBy   *ItemCreator `json:"createdBy"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type ItemCreator struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type itemRow struct {
	ID                  uuid.UUID     `db:"id"`
	ListID              uuid.UUID     `db:"list_id"`
	Name                string        `db:"name"`
	Description         *string       `db:"description"`
	CreatedByMemberID   uuid.NullUUID `db:"created_by_member_id"`
	CreatedByMemberName *string       `db:"created_by_member_name"`
	CreatedAt           time.Time     `db:"created_at"`
	UpdatedAt           time.Time     `db:"updated_at"`
}

// Handlers

type CreateItemPayload struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type CreateItemParams struct {
	HouseholdID   uuid.UUID
	ListID        uuid.UUID
	ActorMemberID uuid.UUID
	Name          string
	Description   *string
}

type DeleteItemParams struct {
	HouseholdID   uuid.UUID
	ListID        uuid.UUID
	ItemID        uuid.UUID
	ActorMemberID uuid.UUID
}

// Store
type InsertItemParams struct {
	ListID            uuid.UUID
	Name              string
	Description       *string
	CreatedByMemberID uuid.UUID
}
