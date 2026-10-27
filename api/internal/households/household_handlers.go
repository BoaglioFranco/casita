package households

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (service *HouseholdService) meHandler(c *gin.Context) {
	member := MemberFromContext(c)

	c.JSON(http.StatusOK, member)
}
