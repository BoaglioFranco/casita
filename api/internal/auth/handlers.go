package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *AuthService) loginHandler(c *gin.Context) {
	var loginInfo LoginPayload
	err := c.ShouldBindJSON(&loginInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}
	token, err := s.Login(c.Request.Context(), loginInfo.Email, loginInfo.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		c.Status(http.StatusUnauthorized)
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "prod",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().UTC().AddDate(1, 0, 0),
	})

	c.Status(http.StatusOK)
}
