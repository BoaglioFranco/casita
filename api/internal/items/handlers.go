package items

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BoaglioFranco/groceriesList/api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *ItemService) createItemHandler(c *gin.Context) {
	memberId := utils.MemberIDFromContext(c)
	params := &CreateItemPayload{}

	listID, err := strconv.ParseInt(c.Param("listId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	householdID, err := strconv.ParseInt(c.Param("householdId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household id"})
		return
	}

	err = c.ShouldBindJSON(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	err = s.CreateItem(c.Request.Context(), CreateItemParams{
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

func (s *ItemService) removeItemHandler(c *gin.Context) {
	memberID := utils.MemberIDFromContext(c)

	listID, err := strconv.ParseInt(c.Param("listId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	householdID, err := strconv.ParseInt(c.Param("householdId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household id"})
		return
	}
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item id"})
		return
	}
	err = s.DeleteItem(c.Request.Context(), DeleteItemParams{
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

func (s *ItemService) getItemsHandler(c *gin.Context) {
	listID, err := strconv.ParseInt(c.Param("listId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	items, err := s.store.GetItemsFromList(c.Request.Context(), listID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
