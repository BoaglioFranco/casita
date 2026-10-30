package households

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequireAuth populates the context with these parameters and guarantees that they will be valid UUIDs.
func MemberIDFromContext(c *gin.Context) uuid.UUID {
	memberID := c.MustGet("member_id").(uuid.UUID)
	return memberID
}

func HouseholdIDFromContext(c *gin.Context) uuid.UUID {
	householdID := c.MustGet("household_id").(uuid.UUID)
	return householdID
}

func MemberFromContext(c *gin.Context) *Member {
	member := c.MustGet("member").(*Member)
	return member
}
