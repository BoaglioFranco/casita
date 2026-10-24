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
	query := `
		INSERT INTO push_subscriptions (household_id, member_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(endpoint) DO UPDATE SET
			household_id = excluded.household_id,
			member_id = excluded.member_id,
			p256dh = excluded.p256dh,
			auth = excluded.auth,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.ExecContext(ctx, query, params.HouseholdId, params.MemberId, params.Endpoint, params.P256dh, params.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationStore) DeleteSubscription(ctx context.Context, endpoint string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM push_subscriptions WHERE endpoint = $1`, endpoint)
	return err
}

func (s *NotificationStore) GetSubscriptionsByHousehold(ctx context.Context, householdID int64) ([]PushSubscription, error) {
	query := `
		SELECT household_id, member_id, endpoint, p256dh, auth
		FROM push_subscriptions
		WHERE household_id = $1
	`
	householdSubscriptions := []PushSubscription{}
	err := s.db.SelectContext(ctx, &householdSubscriptions, query, householdID)
	if err != nil {
		return nil, err
	}

	return householdSubscriptions, nil
}

func (s *NotificationStore) GetSubscriptionsForOtherMembers(ctx context.Context, householdID int64, actorMemberID int64) ([]PushSubscription, error) {
	query := `
		SELECT household_id, member_id, endpoint, p256dh, auth
		FROM push_subscriptions
		WHERE household_id = $1 AND member_id != $2
	`
	householdSubscriptions := []PushSubscription{}
	err := s.db.SelectContext(ctx, &householdSubscriptions, query, householdID, actorMemberID)
	if err != nil {
		return nil, err
	}

	return householdSubscriptions, nil
}
