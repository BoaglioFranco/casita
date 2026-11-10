package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (service *AuthService) RequireAuth(c *gin.Context) {

	token, err := c.Cookie("session")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	session, err := service.AuthenticateSession(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrSessionExpired) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Set("user_id", session.UserID)
	c.Set("household_id", session.HouseholdID)
	c.Set("current_user", CurrentUser{ID: session.UserID, Name: session.Name, HouseholdID: session.HouseholdID})
	c.Next()
}
