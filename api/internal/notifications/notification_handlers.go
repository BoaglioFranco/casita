package notifications

import (
	"net/http"

	"github.com/BoaglioFranco/casita/api/internal/auth"
	"github.com/gin-gonic/gin"
)

func (service *NotificationService) publicVapidKeyHandler(c *gin.Context) {
	publicKey, err := service.publicVapidKey(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Unable to retrieve VAPID public key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"publicKey": publicKey})
}

func (service *NotificationService) subscribeHandler(c *gin.Context) {
	var subscriptionInfo SubscriptionPayload
	if err := c.ShouldBindJSON(&subscriptionInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription data"})
		return
	}
	userID := auth.UserIDFromContext(c)
	if err := service.subscribe(c.Request.Context(), userID, subscriptionInfo); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Unable to save push subscription"})
		return
	}
	c.Status(http.StatusOK)
}
