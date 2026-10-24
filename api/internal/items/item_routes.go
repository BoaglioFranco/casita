package items

import "github.com/gin-gonic/gin"

func (service *ItemService) SetupItemRoutes(router gin.IRoutes) {
	router.POST("/lists/:listId/items", service.createItemHandler)
	router.GET("/lists/:listId/items", service.getItemsHandler)
	router.DELETE("/lists/:listId/items/:itemId", service.removeItemHandler)
}
