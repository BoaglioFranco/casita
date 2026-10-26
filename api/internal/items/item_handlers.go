package items

import (
	"errors"
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (service *ItemService) createItemHandler(c *gin.Context) {
	memberId := utils.MemberIDFromContext(c)
	householdID := utils.HouseholdIDFromContext(c)
	params := &CreateItemPayload{}

	listID, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	householdIDParam, err := uuid.Parse(c.Param("householdId"))
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

	err = service.CreateItem(c.Request.Context(), CreateItemParams{
		HouseholdID:   householdID,
		ListID:        listID,
		ActorMemberID: memberId,
		Name:          params.Name,
		Description:   params.Description,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (service *ItemService) removeItemHandler(c *gin.Context) {
	memberID := utils.MemberIDFromContext(c)
	householdID := utils.HouseholdIDFromContext(c)

	listID, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	householdIDParam, err := uuid.Parse(c.Param("householdId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household id"})
		return
	}
	if householdIDParam != householdID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item id"})
		return
	}
	err = service.DeleteItem(c.Request.Context(), DeleteItemParams{
		HouseholdID:   householdID,
		ListID:        listID,
		ItemID:        itemID,
		ActorMemberID: memberID,
	})
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (service *ItemService) getItemsHandler(c *gin.Context) {
	householdID := utils.HouseholdIDFromContext(c)

	householdIDParam, err := uuid.Parse(c.Param("householdId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household id"})
		return
	}
	if householdIDParam != householdID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	listID, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	items, err := service.itemStore.GetItemsFromList(c.Request.Context(), listID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
