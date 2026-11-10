package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/households"
	"github.com/google/uuid"
)

var ErrNotImplemented = errors.New("notification service integration is not implemented")

type NotificationService struct {
	httpClient     http.Client
	baseURL        string
	householdStore *households.HouseholdStore
}

func NewNotificationService(db *database.Database, notificationServerURL string) *NotificationService {
	return &NotificationService{
		httpClient:     http.Client{Timeout: 5 * time.Second},
		baseURL:        notificationServerURL,
		householdStore: households.NewHouseholdStore(db),
	}
}

func (service *NotificationService) NotifyOtherUsers(ctx context.Context, householdID uuid.UUID, actorUserID uuid.UUID, notification Notification) error {
	userIDs, err := service.householdStore.UserIDsByHouseholdID(ctx, householdID)
	if err != nil {
		return err
	}
	if !slices.Contains(userIDs, actorUserID) {
		return errors.New("notification actor does not belong to the household")
	}
	results := make(chan error, len(userIDs)-1)
	for _, userID := range userIDs {
		if userID == actorUserID {
			continue
		}
		go func(userID uuid.UUID) {
			results <- service.notify(ctx, actorUserID, NotificationRequest{
				RequestID: uuid.New(), UserID: userID, Notification: notification,
			})
		}(userID)
	}
	var deliveryErrors []error
	for range len(userIDs) - 1 {
		if err := <-results; err != nil {
			deliveryErrors = append(deliveryErrors, err)
		}
	}
	return errors.Join(deliveryErrors...)
}

func (service *NotificationService) publicVapidKey(ctx context.Context) (string, error) {
	endpoint, err := url.JoinPath(service.baseURL, "/vapid-public-key")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	res, err := service.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("notification service returned status %d", res.StatusCode)
	}
	var payload struct {
		PublicKey string `json:"publicKey"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.PublicKey == "" {
		return "", errors.New("notification service returned an empty VAPID public key")
	}
	return payload.PublicKey, nil
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

func (service *NotificationService) notify(ctx context.Context, userID uuid.UUID, nr NotificationRequest) error {
	token, err := auth.MintS2sJWT(userID)
	if err != nil {
		return err
	}
	endpoint, err := url.JoinPath(service.baseURL, "/notifications/notify")
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(nr)
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
