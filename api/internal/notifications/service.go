package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

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

	var deliveryErrors []error
	for _, sub := range subscriptions {
		if err := s.sendNotification(ctx, sub, notification); err != nil {
			var pushErr *PushServiceError
			if errors.As(err, &pushErr) && (pushErr.StatusCode == http.StatusNotFound || pushErr.StatusCode == http.StatusGone) {
				if deleteErr := s.store.DeleteSubscription(ctx, sub.Endpoint); deleteErr != nil {
					deliveryErrors = append(deliveryErrors, fmt.Errorf("remove expired push subscription: %w", deleteErr))
				}
			}
			deliveryErrors = append(deliveryErrors, err)
		}
	}

	return errors.Join(deliveryErrors...)
}

func (s *NotificationService) sendNotification(ctx context.Context, subscription PushSubscription, notification Notification) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	response, err := webpush.SendNotificationWithContext(
		ctx,
		payload,
		subscription.ToWebpushSubscription(),
		&webpush.Options{
			Subscriber:      "francoboaglio@gmail.com",
			VAPIDPublicKey:  s.vapidPublicKey,
			VAPIDPrivateKey: s.vapidPrivateKey,
			TTL:             24 * 60 * 60,
		},
	)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			return fmt.Errorf("send push request: %w", urlErr.Err)
		}
		return fmt.Errorf("send push request: %w", err)
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(io.LimitReader(response.Body, 4096))
	if readErr != nil {
		return fmt.Errorf("read push service response: %w", readErr)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return &PushServiceError{
			StatusCode: response.StatusCode,
			Status:     response.Status,
			Detail:     strings.TrimSpace(string(body)),
		}
	}

	return nil
}
