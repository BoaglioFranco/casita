package subscriptions

import (
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (service *NotificationService) publicVapidKeyHandler(c *gin.Context) {
	if service.vapidPublicKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "VAPID public key not set",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"publicKey": service.vapidPublicKey,
	})
}

func (service *NotificationService) subscribeHandler(c *gin.Context) {
	// Authentication middleware must set this UUID after verifying the caller.
	value, exists := c.Get("user_id")
	userID, valid := value.(uuid.UUID)
	if !exists || !valid || userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var subscription webpush.Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data"})
		return
	}

	if err := service.subscriptionStore.StoreSubscription(c.Request.Context(), StoreSubscriptionParams{
		UserID:   userID,
		Endpoint: subscription.Endpoint,
		P256dh:   subscription.Keys.P256dh,
		Auth:     subscription.Keys.Auth,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save push subscription"})
		return
	}

	c.Status(http.StatusOK)
}

func (service *NotificationService) notifyHandler(c *gin.Context) {
	var notiRequest NotificationRequest
	if err := c.ShouldBindJSON(&notiRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification data"})
		return
	}
	if err := service.Notify(c.Request.Context(), notiRequest); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
