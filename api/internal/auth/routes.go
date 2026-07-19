package auth

import "github.com/gin-gonic/gin"

func (s *AuthService) SetupAuthRoutes(router gin.IRoutes) {
	router.POST("/login", s.loginHandler)
}
