package lists

import (
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/auth"
	"github.com/gin-gonic/gin"
)

func (service *ListService) createListHandler(c *gin.Context) {
	params := &CreateListPayload{}
	householdID := auth.HouseholdIDFromContext(c)

	err := c.ShouldBindJSON(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	err = service.listStore.InsertList(c.Request.Context(), InsertListParams{HouseholdID: householdID, Name: params.Name})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}

func (service *ListService) getListsHandler(c *gin.Context) {
	householdID := auth.HouseholdIDFromContext(c)

	lists, err := service.listStore.GetLists(c.Request.Context(), householdID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"lists": lists})
}
