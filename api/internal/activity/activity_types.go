package activity

import "github.com/google/uuid"

type ActivityEventType string

const (
	ActivityItemCreate ActivityEventType = "ITEM_CREATE"
	ActivityItemDelete ActivityEventType = "ITEM_DELETE"
)

type RecordActivityParams struct {
	EventType        ActivityEventType
	HouseholdID      uuid.UUID
	ListID           *uuid.UUID
	ItemID           *uuid.UUID
	ActorUserID      *uuid.UUID
	ItemNameSnapshot *string
}
