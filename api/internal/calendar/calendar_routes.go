package calendar

import "github.com/gin-gonic/gin"

func (service *CalendarService) SetupCalendarRoutes(router gin.IRoutes) {
	router.GET("/calendar", service.getCalendar)
	router.POST("/calendar/events", service.createEvent)
	router.PUT("/calendar/events/:eventId", service.updateEvent)
}
