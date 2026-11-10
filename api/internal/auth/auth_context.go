package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserIDFromContext(c *gin.Context) uuid.UUID {
	userID := c.MustGet("user_id").(uuid.UUID)
	return userID
}

func HouseholdIDFromContext(c *gin.Context) uuid.UUID {
	return c.MustGet("household_id").(uuid.UUID)
}

func CurrentUserFromContext(c *gin.Context) CurrentUser {
	return c.MustGet("current_user").(CurrentUser)
}
