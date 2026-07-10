package notifications

import (
	"context"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type NotificationStore struct {
	db *database.Database
}

func NewStore(db *database.Database) *NotificationStore {
	return &NotificationStore{db: db}
}

func (s *NotificationStore) StoreSubscription(ctx context.Context, params StoreSubscriptionParams) error {
	query := `INSERT INTO push_subscriptions (household_id, member_id, endpoint, p256dh, auth) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, params.HouseholdId, params.MemberId, params.Endpoint, params.P256dh, params.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationStore) GetSubscriptionsByHousehold(ctx context.Context, householdId int64) ([]PushSubscription, error) {
	query := `SELECT household_id, member_id, endpoint, p256dh, auth FROM push_subscriptions WHERE household_id = ?`
	householdSubscriptions := []PushSubscription{}
	err := s.db.SelectContext(ctx, &householdSubscriptions, query, householdId)
	if err != nil {
		return nil, err
	}

	return householdSubscriptions, nil
}
