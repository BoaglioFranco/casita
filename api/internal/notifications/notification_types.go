package notifications

import (
	"fmt"
	"time"

	"github.com/SherClockHolmes/webpush-go"
)

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
	HouseholdId int64
	MemberId    int64
	Endpoint    string
	P256dh      string
	Auth        string
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PushSubscriptionPayload struct {
	Endpoint string                      `json:"endpoint"`
	Keys     PushSubscriptionPayloadKeys `json:"keys"`
}

type PushSubscriptionPayloadKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type PushSubscription struct {
	ID          int64     `db:"id"`
	HouseholdID int64     `db:"household_id"`
	MemberID    int64     `db:"member_id"`
	Endpoint    string    `db:"endpoint"`
	P256dh      string    `db:"p256dh"`
	Auth        string    `db:"auth"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
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
