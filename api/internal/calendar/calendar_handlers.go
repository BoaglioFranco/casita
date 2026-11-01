package calendar

import (
	"errors"
	"net/http"

	"github.com/BoaglioFranco/groceriesList/api/internal/households"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (service *CalendarService) updateEvent(c *gin.Context) {
	householdID := households.HouseholdIDFromContext(c)
	params := &UpdateEventPayload{}

	eventID, err := uuid.Parse(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
		return
	}
	if err := c.ShouldBindJSON(params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	err = service.UpdateEvent(c.Request.Context(), UpdateEventParams{
		HouseholdID:               householdID,
		EventID:                   eventID,
		Name:                      params.Name,
		Description:               params.Description,
		StartsAt:                  params.StartsAt,
		DurationMinutes:           params.DurationMinutes,
		NotificationOffsetMinutes: params.NotificationOffsetMinutes,
		RecurrenceRule:            params.RecurrenceRule,
	})
	if err != nil {
		if errors.Is(err, ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (service *CalendarService) getCalendar(c *gin.Context) {
	householdID := households.HouseholdIDFromContext()

	var query GetCalendarQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	events, err := service.calendarStore.GetEvents(c.Request.Context(), GetEventsParams{
		From:        query.From,
		To:          query.To,
		HouseholdID: householdID,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	service.GenerateEventTimeline(events, query.From, query.To)

}
