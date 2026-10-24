package auth

import "github.com/gin-gonic/gin"

func (service *AuthService) SetupAuthRoutes(router gin.IRoutes) {
	router.POST("/login", service.loginHandler)
}
