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

func (service *NotificationService) NotifyOtherMembers(ctx context.Context, householdID uuid.UUID, actorMemberID uuid.UUID, notification Notification) error {
	members, err := service.householdStore.MembersByHouseholdID(ctx, householdID)
	if err != nil {
		return err
	}
	var actorUserID uuid.UUID
	for _, member := range members {
		if member.ID == actorMemberID {
			actorUserID = member.UserID
			break
		}
	}
	if actorUserID == uuid.Nil {
		return errors.New("notification actor is not a member of the household")
	}
	otherMembers := slices.DeleteFunc(members, func(m households.Member) bool {
		return m.ID == actorMemberID
	})
	results := make(chan error, len(otherMembers))
	for _, member := range otherMembers {
		go func(m households.Member) {
			results <- service.notify(ctx, actorUserID, NotificationRequest{
				RequestID:    uuid.New(),
				UserID:       m.UserID,
				Notification: notification,
			})
		}(member)
	}

	var deliveryErrors []error
	for range otherMembers {
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
