package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/google/uuid"
)

var ErrNotImplemented = errors.New("notification service integration is not implemented")

type NotificationService struct {
	httpClient http.Client
	baseURL    string
}

func NewNotificationService(notificationServerURL string) *NotificationService {
	return &NotificationService{
		httpClient: http.Client{Timeout: 5 * time.Second},
		baseURL:    notificationServerURL,
	}
}

func (service *NotificationService) NotifyOtherMembers(ctx context.Context, householdID uuid.UUID, actorMemberID uuid.UUID, notification Notification) error {
	// TODO: Resolve recipients and submit delivery requests to the notification service.
	return ErrNotImplemented
}

func (service *NotificationService) subscribe(ctx context.Context, userID uuid.UUID, payload SubscriptionPayload) error {
	token, err := auth.MintS2sJWT(userID)
	if err != nil {
		return err
	}
	endpoint, err := url.JoinPath(service.baseURL, "/notifications/subscribe")
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := service.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("notification service returned status %d", res.StatusCode)
	}

	return nil
}
