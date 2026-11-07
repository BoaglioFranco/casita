package notifications

import (
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/gin-gonic/gin"
)

func (service *NotificationService) publicVapidKeyHandler(c *gin.Context) {
	// TODO: Retrieve the public key from the notification service.
	c.JSON(http.StatusNotImplemented, gin.H{"error": ErrNotImplemented.Error()})
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
