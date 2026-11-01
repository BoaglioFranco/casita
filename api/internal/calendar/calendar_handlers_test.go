package calendar

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updateEventStoreFunc func(context.Context, UpdateEventStoreParams) error

func (store updateEventStoreFunc) InsertEvent(context.Context, InsertEventParams) error {
	panic("unexpected insert")
}

func (store updateEventStoreFunc) UpdateEvent(ctx context.Context, params UpdateEventStoreParams) error {
	return store(ctx, params)
}

func TestUpdateEvent(t *testing.T) {
	householdID, eventID := uuid.New(), uuid.New()
	validBody := `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z"}`
	for _, tt := range []struct {
		name   string
		id     string
		body   string
		err    error
		status int
		offset int
	}{
		{"updated", eventID.String(), validBody, nil, http.StatusNoContent, 10},
		{"zero offset", eventID.String(), `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","notificationOffsetMinutes":0}`, nil, http.StatusNoContent, 0},
		{"invalid id", "bad-id", validBody, nil, http.StatusBadRequest, 0},
		{"invalid payload", eventID.String(), `{}`, nil, http.StatusBadRequest, 0},
		{"not found", eventID.String(), validBody, ErrEventNotFound, http.StatusNotFound, 10},
		{"store error", eventID.String(), validBody, errors.New("update failed"), http.StatusInternalServerError, 10},
	} {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			service := &CalendarService{
				calendarStore: updateEventStoreFunc(func(ctx context.Context, params UpdateEventStoreParams) error {
					called = true
					if params.HouseholdID != householdID || params.EventID != eventID || params.NotificationOffsetMinutes != tt.offset {
						t.Fatalf("unexpected update params: %+v", params)
					}
					return tt.err
				}),
			}
			router := gin.New()
			router.Use(func(c *gin.Context) { c.Set("household_id", householdID) })
			service.SetupCalendarRoutes(router)
			request := httptest.NewRequest(http.MethodPut, "/calendar/events/"+tt.id, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.status, response.Body.String())
			}
			if called != (tt.status != http.StatusBadRequest) {
				t.Fatalf("unexpected store call: %v", called)
			}
		})
	}
}

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
