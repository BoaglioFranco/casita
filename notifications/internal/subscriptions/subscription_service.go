package subscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/BoaglioFranco/groceriesList/notifications/internal/database"
	"github.com/SherClockHolmes/webpush-go"
)

type NotificationService struct {
	subscriptionStore *SubscriptionStore
	vapidPublicKey    string
	vapidPrivateKey   string
}

func NewNotificationService(db *database.Database, vapidPublicKey, vapidPrivateKey string) *NotificationService {
	return &NotificationService{
		subscriptionStore: NewSubscriptionStore(db),
		vapidPublicKey:    vapidPublicKey,
		vapidPrivateKey:   vapidPrivateKey,
	}
}

func (service *NotificationService) Notify(ctx context.Context, request NotificationRequest) error {
	subs, err := service.subscriptionStore.GetSubscriptionsByUser(ctx, request.UserID)
	if err != nil {
		return err
	}

	results := make(chan error, len(subs))
	for _, sub := range subs {
		go func(s PushSubscription) {
			err := service.sendNotification(ctx, s, request.Notification)
			results <- err
		}(sub)
	}

	deliveryErrors := make([]error, 0)
	for range subs {
		if err := <-results; err != nil {
			deliveryErrors = append(deliveryErrors, err)
		}
	}
	return errors.Join(deliveryErrors...)
}

func (service *NotificationService) sendNotification(ctx context.Context, subscription PushSubscription, notification Notification) error {
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
			VAPIDPublicKey:  service.vapidPublicKey,
			VAPIDPrivateKey: service.vapidPrivateKey,
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
