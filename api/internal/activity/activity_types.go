package activity

type ActivityEventType string

const (
	ActivityItemCreate ActivityEventType = "ITEM_CREATE"
	ActivityItemDelete ActivityEventType = "ITEM_DELETE"
)

type RecordActivityParams struct {
	EventType        ActivityEventType
	HouseholdID      int64
	ListID           *int64
	ItemID           *int64
	ActorMemberID    *int64
	ItemNameSnapshot *string
}
