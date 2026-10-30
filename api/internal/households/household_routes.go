package households

import "github.com/gin-gonic/gin"

func (service *HouseholdService) SetupHouseholdRoutes(router gin.IRouter) {
	router.GET("/me", service.meHandler)
}

func (service *HouseholdService) SetupRootHouseholdsRoutes(router gin.IRouter) {

	router.GET("/households", service.householdsHandler)
}
