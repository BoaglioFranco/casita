package notifications

import (
	"context"
	"encoding/json"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/SherClockHolmes/webpush-go"
)

type NotificationService struct {
	store           *NotificationStore
	vapidPublicKey  string
	vapidPrivateKey string
}

func NewNotificationService(db *database.Database, vapidPublicKey string, vapidPrivateKey string) *NotificationService {
	return &NotificationService{
		store:           NewStore(db),
		vapidPublicKey:  vapidPublicKey,
		vapidPrivateKey: vapidPrivateKey,
	}
}

func (s *NotificationService) NotifyMembers(ctx context.Context, householdId int64, notification Notification) error {

	subscriptions, err := s.store.GetSubscriptionsByHousehold(ctx, householdId)
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		err = s.sendNotification(sub, notification)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) sendNotification(subscription PushSubscription, notification Notification) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	response, err := webpush.SendNotification(
		payload,
		subscription.ToWebpushSubscription(),
		&webpush.Options{
			Subscriber:      "mailto:test@example.com",
			VAPIDPublicKey:  s.vapidPublicKey,
			VAPIDPrivateKey: s.vapidPrivateKey,
			TTL:             60,
		},
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
