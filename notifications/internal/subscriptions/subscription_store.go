package subscriptions

import (
	"context"

	"github.com/BoaglioFranco/casita/notifications/internal/database"
	"github.com/google/uuid"
)

type SubscriptionStore struct {
	db *database.Database
}

func NewSubscriptionStore(db *database.Database) *SubscriptionStore {
	return &SubscriptionStore{db: db}
}

func (store *SubscriptionStore) StoreSubscription(ctx context.Context, params StoreSubscriptionParams) error {
	query := `
		INSERT INTO push_subscriptions (user_id, endpoint, p256dh, auth)
		    VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint)
		    DO UPDATE SET
		        user_id = excluded.user_id,
		        p256dh = excluded.p256dh,
		        auth = excluded.auth,
		        updated_at = CURRENT_TIMESTAMP`
	_, err := store.db.ExecContext(ctx, query, params.UserID, params.Endpoint, params.P256dh, params.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (store *SubscriptionStore) DeleteSubscription(ctx context.Context, endpoint string) error {
	_, err := store.db.ExecContext(ctx, `
		DELETE FROM push_subscriptions
		WHERE endpoint = $1`, endpoint)
	return err
}

func (store *SubscriptionStore) GetSubscriptionsByUser(ctx context.Context, userID uuid.UUID) ([]PushSubscription, error) {
	query := `
		SELECT
		    id,
		    user_id,
		    endpoint,
		    p256dh,
		    auth,
		    created_at,
		    updated_at
		FROM
		    push_subscriptions
		WHERE
		    user_id = $1`

	subscriptions := []PushSubscription{}
	err := store.db.SelectContext(ctx, &subscriptions, query, userID)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}
