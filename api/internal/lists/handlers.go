package lists

import (
	"net/http"
	"strconv"

	"github.com/BoaglioFranco/groceriesList/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *ListService) createListHandler(c *gin.Context) {
	params := &CreateListPayload{}
	householdID := utils.HouseholdIDFromContext(c)

	householdIDParam, err := strconv.ParseInt(c.Param("householdId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household id"})
		return
	}
	if householdIDParam != householdID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err = c.ShouldBindJSON(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	err = s.store.InsertList(c.Request.Context(), InsertListParams{HouseholdID: householdID, Name: params.Name})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}
