package households

import (
	"errors"
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (service *HouseholdService) RequireMembership(c *gin.Context) {
	householdIDString := c.Param("householdId")
	if householdIDString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	householdID, err := uuid.Parse(householdIDString)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userID := auth.UserIDFromContext(c)
	member, err := service.householdStore.MemberByUserID(c.Request.Context(), userID, householdID)
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
	c.Set("member", member)
	c.Next()
}
