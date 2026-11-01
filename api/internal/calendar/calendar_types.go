package calendar

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrEventNotFound = errors.New("event not found")

type CalendarService struct {
	calendarStore eventStore
}

type eventStore interface {
	InsertEvent(ctx context.Context, params InsertEventParams) error
	UpdateEvent(ctx context.Context, params UpdateEventStoreParams) error
}

type CreateEventPayload struct {
	Name                      string    `json:"name" binding:"required"`
	Description               *string   `json:"description"`
	StartsAt                  time.Time `json:"startsAt" binding:"required"`
	DurationMinutes           *int      `json:"durationMinutes" binding:"omitempty,gt=0"`
	NotificationOffsetMinutes *int      `json:"notificationOffsetMinutes" binding:"omitempty,gte=0"`
	RecurrenceRule            *string   `json:"recurrenceRule"`
}

type CreateEventParams struct {
	HouseholdID               uuid.UUID
	CreatedByMemberID         uuid.UUID
	Name                      string
	Description               *string
	StartsAt                  time.Time
	DurationMinutes           *int
	NotificationOffsetMinutes *int
	RecurrenceRule            *string
}

type InsertEventParams struct {
	HouseholdID               uuid.UUID
	CreatedByMemberID         uuid.UUID
	Name                      string
	Description               *string
	StartsAt                  time.Time
	DurationMinutes           *int
	NotificationOffsetMinutes int
	RecurrenceRule            *string
}

type UpdateEventPayload CreateEventPayload

type UpdateEventParams struct {
	HouseholdID               uuid.UUID
	EventID                   uuid.UUID
	Name                      string
	Description               *string
	StartsAt                  time.Time
	DurationMinutes           *int
	NotificationOffsetMinutes *int
	RecurrenceRule            *string
}

type UpdateEventStoreParams struct {
	HouseholdID               uuid.UUID
	EventID                   uuid.UUID
	Name                      string
	Description               *string
	StartsAt                  time.Time
	DurationMinutes           *int
	NotificationOffsetMinutes int
	RecurrenceRule            *string
}

type GetCalendarQuery struct {
	From time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To   time.Time `form:"to" binding:"required,gtfield=From" time_format:"2006-01-02"`
}
