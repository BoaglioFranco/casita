package calendar

import (
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/households"
	"github.com/gin-gonic/gin"
)

func (service *CalendarService) createEvent(c *gin.Context) {
	householdID := households.HouseholdIDFromContext(c)
	memberID := households.MemberIDFromContext(c)
	params := &CreateEventPayload{}

	err := c.ShouldBindJSON(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	err = service.CreateEvent(c.Request.Context(), CreateEventParams{
		HouseholdID:               householdID,
		CreatedByMemberID:         memberID,
		Name:                      params.Name,
		Description:               params.Description,
		StartsAt:                  params.StartsAt,
		DurationMinutes:           params.DurationMinutes,
		NotificationOffsetMinutes: params.NotificationOffsetMinutes,
		RecurrenceRule:            params.RecurrenceRule,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
