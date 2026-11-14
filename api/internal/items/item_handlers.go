package items

import (
	"errors"
	"net/http"

	"github.com/BoaglioFranco/casita/api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (service *ItemService) createItemHandler(c *gin.Context) {
	user := auth.CurrentUserFromContext(c)
	params := &CreateItemPayload{}

	listID, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	err = c.ShouldBindJSON(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	err = service.CreateItem(c.Request.Context(), CreateItemParams{
		User:        user,
		ListID:      listID,
		Name:        params.Name,
		Description: params.Description,
	})
	if errors.Is(err, ErrItemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item or list not found"})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (service *ItemService) removeItemHandler(c *gin.Context) {
	user := auth.CurrentUserFromContext(c)

	listID, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item id"})
		return
	}
	err = service.DeleteItem(c.Request.Context(), DeleteItemParams{
		User:   user,
		ListID: listID,
		ItemID: itemID,
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
	listID, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list id"})
		return
	}
	items, err := service.itemStore.GetItemsFromList(c.Request.Context(), auth.HouseholdIDFromContext(c), listID)
	if errors.Is(err, ErrItemNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item or list not found"})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
