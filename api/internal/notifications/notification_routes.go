package notifications

import "github.com/gin-gonic/gin"

func (service *NotificationService) SetupPublicNotificationRoutes(router gin.IRoutes) {
	router.GET("/vapid-public-key", service.publicVapidKeyHandler)
}

func (service *NotificationService) SetupNotificationRoutes(router gin.IRoutes) {
	router.POST("/notifications/subscribe", service.subscribeHandler)
}
