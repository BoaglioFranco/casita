package lists

import "github.com/gin-gonic/gin"

func (service *ListService) SetupListRoutes(router gin.IRoutes) {
	router.POST("/lists", service.createListHandler)
	router.GET("/lists", service.getListsHandler)
}
