package calendar

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getEventsStoreFunc func(context.Context, GetEventsParams) ([]CalendarEvent, error)

func (store getEventsStoreFunc) InsertEvent(context.Context, InsertEventParams) error {
	panic("unexpected insert")
}

func (store getEventsStoreFunc) UpdateEvent(context.Context, UpdateEventStoreParams) error {
	panic("unexpected update")
}

func (store getEventsStoreFunc) GetEvents(ctx context.Context, params GetEventsParams) ([]CalendarEvent, error) {
	return store(ctx, params)
}

func TestGetCalendar(t *testing.T) {
	householdID := uuid.New()
	// Date-only query parameters currently bind in the server's local timezone.
	from := time.Date(2026, 11, 1, 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 2)
	duration, invalidRule := 30, "FREQ=INVALID"
	event := CalendarEvent{ID: uuid.New(), StartsAt: from, DurationMinutes: &duration}
	badEvent := event
	badEvent.RecurrenceRule = &invalidRule
	validQuery := "?from=2026-11-01&to=2026-11-03"
	for _, tt := range []struct {
		name   string
		query  string
		events []CalendarEvent
		err    error
		status int
	}{
		{"success", validQuery, []CalendarEvent{event}, nil, http.StatusOK},
		{"empty calendar", validQuery, nil, nil, http.StatusOK},
		{"missing from", "?to=2026-11-03", nil, nil, http.StatusBadRequest},
		{"missing to", "?from=2026-11-01", nil, nil, http.StatusBadRequest},
		{"invalid date", "?from=invalid&to=2026-11-03", nil, nil, http.StatusBadRequest},
		{"reversed range", "?from=2026-11-03&to=2026-11-01", nil, nil, http.StatusBadRequest},
		{"equal boundaries", "?from=2026-11-01&to=2026-11-01", nil, nil, http.StatusBadRequest},
		{"store error", validQuery, nil, errors.New("read failed"), http.StatusInternalServerError},
		{"invalid stored rule", validQuery, []CalendarEvent{badEvent}, nil, http.StatusInternalServerError},
	} {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			service := &CalendarService{calendarStore: getEventsStoreFunc(func(ctx context.Context, params GetEventsParams) ([]CalendarEvent, error) {
				called = true
				if params.HouseholdID != householdID || !params.From.Equal(from) || !params.To.Equal(to) {
					t.Fatalf("unexpected query params: %+v", params)
				}
				return tt.events, tt.err
			})}
			router := gin.New()
			group := router.Group("/households/:householdId")
			group.Use(func(c *gin.Context) { c.Set("household_id", householdID) })
			service.SetupCalendarRoutes(group)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/households/"+householdID.String()+"/"+tt.query, nil))
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.status, response.Body.String())
			}
			if called != (tt.status != http.StatusBadRequest) {
				t.Fatalf("unexpected store call: %v", called)
			}
			if tt.status == http.StatusOK {
				var body struct {
					From        time.Time            `json:"from"`
					To          time.Time            `json:"to"`
					Occurrences []CalendarOccurrence `json:"occurrences"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if !body.From.Equal(from) || !body.To.Equal(to) || body.Occurrences == nil || len(body.Occurrences) != len(tt.events) {
					t.Fatalf("unexpected response: %s", response.Body.String())
				}
				if len(body.Occurrences) > 0 {
					o := body.Occurrences[0]
					if o.Event.ID != event.ID || !o.StartsAt.Equal(from) || !o.EndsAt.Equal(from.Add(30*time.Minute)) {
						t.Fatalf("unexpected occurrence: %+v", o)
					}
				}
			}
		})
	}
}

type updateEventStoreFunc func(context.Context, UpdateEventStoreParams) error

func (store updateEventStoreFunc) InsertEvent(context.Context, InsertEventParams) error {
	panic("unexpected insert")
}

func (store updateEventStoreFunc) UpdateEvent(ctx context.Context, params UpdateEventStoreParams) error {
	return store(ctx, params)
}

func (store updateEventStoreFunc) GetEvents(context.Context, GetEventsParams) ([]CalendarEvent, error) {
	panic("unexpected get")
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
		{"invalid recurrence", eventID.String(), `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","recurrenceRule":"FREQ=INVALID"}`, nil, http.StatusBadRequest, 0},
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
		{"invalid recurrence", `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","recurrenceRule":"FREQ=INVALID"}`},
		{"empty recurrence", `{"name":"Appointment","startsAt":"2026-10-31T12:00:00Z","recurrenceRule":""}`},
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
