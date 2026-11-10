package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (service *AuthService) SetupAuthRoutes(router gin.IRoutes) {
	router.POST("/login", service.loginHandler)
}

func (service *AuthService) SetupCurrentUserRoutes(router gin.IRoutes) {
	router.GET("/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, CurrentUserFromContext(c))
	})
}
