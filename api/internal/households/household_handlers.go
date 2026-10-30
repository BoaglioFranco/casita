package households

import (
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/gin-gonic/gin"
)

func (service *HouseholdService) meHandler(c *gin.Context) {
	member := MemberFromContext(c)

	c.JSON(http.StatusOK, member)
}

func (service *HouseholdService) householdsHandler(c *gin.Context) {
	userID := auth.UserIDFromContext(c)

	households, err := service.householdStore.HouseholdsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve households"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"households": households})
}
