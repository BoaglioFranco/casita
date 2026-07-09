package lists

import "github.com/gin-gonic/gin"

func (s *ListService) SetupListRoutes(router gin.IRoutes) {
	router.POST("/households/:householdId/lists", s.createListHandler)
}
