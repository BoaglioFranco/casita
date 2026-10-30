package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserIDFromContext(c *gin.Context) uuid.UUID {
	userID := c.MustGet("user_id").(uuid.UUID)
	return userID
}
