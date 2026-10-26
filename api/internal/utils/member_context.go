package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserIDFromContext(c *gin.Context) uuid.UUID {
	userID := c.MustGet("user_id").(uuid.UUID)
	return userID
}

func MemberIDFromContext(c *gin.Context) uuid.UUID {
	memberID := c.MustGet("member_id").(uuid.UUID)
	return memberID
}

func HouseholdIDFromContext(c *gin.Context) uuid.UUID {
	householdID := c.MustGet("household_id").(uuid.UUID)
	return householdID
}
