package items

import "github.com/gin-gonic/gin"

func (s *ItemService) SetupItemRoutes(router gin.IRoutes) {
	router.POST("/lists/:listId/items", s.createItemHandler)
	router.GET("/lists/:listId/items", s.getItemsHandler)
	router.DELETE("/lists/:listId/items/:itemId", s.removeItemHandler)
}
