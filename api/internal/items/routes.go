package items

import "github.com/gin-gonic/gin"

func (s *ItemService) SetupItemRoutes(router gin.IRoutes) {
	router.POST("/households/:householdId/lists/:listId/items", s.createItemHandler)
	router.GET("/households/:householdId/lists/:listId/items", s.getItemsHandler)
	router.DELETE("/households/:householdId/lists/:listId/items/:itemId", s.removeItemHandler)
}
