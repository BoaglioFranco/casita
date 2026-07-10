package utils

import "github.com/gin-gonic/gin"

func MemberIDFromContext(c *gin.Context) int64 {
	memberID := c.MustGet("member_id").(int64)
	return memberID
}

func HouseholdIDFromContext(c *gin.Context) int64 {
	householdID := c.MustGet("household_id").(int64)
	return householdID
}
