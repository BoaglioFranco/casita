package notifications

import (
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/utils"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
)

func (s *NotificationService) publicVapidKeyHandler(c *gin.Context) {
	if s.vapidPublicKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "VAPID public key not set",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"publicKey": s.vapidPublicKey,
	})
}

func (s *NotificationService) subscribeHandler(c *gin.Context) {
	memberID := utils.MemberIDFromContext(c)
	householdID := utils.HouseholdIDFromContext(c)
	var subscription webpush.Subscription
	if err := c.ShouldBindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data"})
		return
	}

	if err := s.store.StoreSubscription(c.Request.Context(), StoreSubscriptionParams{
		HouseholdId: householdID,
		MemberId:    memberID,
		Endpoint:    subscription.Endpoint,
		P256dh:      subscription.Keys.P256dh,
		Auth:        subscription.Keys.Auth,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save push subscription"})
		return
	}

	c.Status(http.StatusOK)
}
