package households

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BoaglioFranco/groceriesList/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *HouseholdService) RequireMembership(c *gin.Context) {
	householdIDString := c.Param("householdId")
	if householdIDString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	householdID, err := strconv.ParseInt(householdIDString, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := utils.UserIDFromContext(c)
	member, err := s.store.MemberByUserID(c.Request.Context(), userID, householdID)
	if err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set("member_id", member.ID)
	c.Set("household_id", householdID)
	c.Next()
}
