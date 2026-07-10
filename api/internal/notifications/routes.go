package notifications

import "github.com/gin-gonic/gin"

func (s *NotificationService) SetupPublicNotificationRoutes(router gin.IRoutes) {
	router.GET("/vapid-public-key", s.publicVapidKeyHandler)
}

func (s *NotificationService) SetupNotificationRoutes(router gin.IRoutes) {
	router.POST("/subscribe", s.subscribeHandler)
}
