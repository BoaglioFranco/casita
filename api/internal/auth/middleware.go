package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *AuthService) RequireAuth(c *gin.Context) {

	token, err := c.Cookie("session")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	session, err := s.AuthenticateSession(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrSessionExpired) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Set("user_id", session.UserID)
	c.Next()
}
