package calendar

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CalendarService struct {
	calendarStore eventStore
}

type eventStore interface {
	InsertEvent(ctx context.Context, params InsertEventParams) error
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
