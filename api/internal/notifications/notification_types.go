package notifications

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
