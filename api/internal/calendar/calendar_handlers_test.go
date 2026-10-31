package calendar

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestCreateEventRejectsInvalidPayload(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"malformed JSON", `{`},
		{"missing name", `{"startsAt":"2026-10-31T12:00:00Z"}`},
		{"missing start", `{"name":"Appointment"}`},
		{"invalid start", `{"name":"Appointment","startsAt":"tomorrow"}`},
		{"zero duration", `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","durationMinutes":0}`},
		{"negative duration", `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","durationMinutes":-1}`},
		{"negative offset", `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","notificationOffsetMinutes":-1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Invalid input must return before accessing the store.
			service := &CalendarService{}
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("household_id", uuid.New())
				c.Set("member_id", uuid.New())
			})
			service.SetupCalendarRoutes(router)
			request := httptest.NewRequest(http.MethodPost, "/calendar/events", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestCreateEventReturnsServerErrorWhenInsertFails(t *testing.T) {
	for _, offset := range []string{"", `,"notificationOffsetMinutes":0`} {
		service := NewCalendarService(nil)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("household_id", uuid.New())
			c.Set("member_id", uuid.New())
		})
		service.SetupCalendarRoutes(router)
		body := `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z"` + offset + `}`
		request := httptest.NewRequest(http.MethodPost, "/calendar/events", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusInternalServerError, response.Body.String())
		}
	}
}
