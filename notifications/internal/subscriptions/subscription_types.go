package subscriptions

import (
	"fmt"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
)

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PushServiceError struct {
	StatusCode int
	Status     string
	Detail     string
}

func (e *PushServiceError) Error() string {
	if e.Detail == "" {
		return fmt.Sprintf("push service returned %s", e.Status)
	}
	return fmt.Sprintf("push service returned %s: %s", e.Status, e.Detail)
}

type StoreSubscriptionParams struct {
	UserID   uuid.UUID
	Endpoint string
	P256dh   string
	Auth     string
}

type PushSubscription struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Endpoint  string    `db:"endpoint"`
	P256dh    string    `db:"p256dh"`
	Auth      string    `db:"auth"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (subscription *PushSubscription) ToWebpushSubscription() *webpush.Subscription {
	return &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.Auth,
			P256dh: subscription.P256dh,
		},
	}
}
