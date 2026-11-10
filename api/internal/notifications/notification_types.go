package notifications

import "github.com/google/uuid"

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type SubscriptionPayload struct {
	Endpoint string                  `json:"endpoint" binding:"required"`
	Keys     SubscriptionPayloadKeys `json:"keys" binding:"required"`
}

type SubscriptionPayloadKeys struct {
	P256dh string `json:"p256dh" binding:"required"`
	Auth   string `json:"auth" binding:"required"`
}

type NotificationRequest struct {
	Notification Notification `json:"notification"`
	UserID       uuid.UUID    `json:"user_id"`
	RequestID    uuid.UUID    `json:"request_id"`
}
