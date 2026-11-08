package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserIDFromContext requires the authentication middleware to have run.
func UserIDFromContext(c *gin.Context) uuid.UUID {
	return c.MustGet("user_id").(uuid.UUID)
}
