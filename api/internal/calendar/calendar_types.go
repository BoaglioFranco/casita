package calendar

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrEventNotFound = errors.New("event not found")
var ErrInvalidRecurrenceRule = errors.New("invalid recurrence rule")

type CalendarEvent struct {
	ID                        uuid.UUID  `db:"id" json:"id"`
	HouseholdID               *uuid.UUID `db:"household_id" json:"householdId"`
	CreatedByMemberID         *uuid.UUID `db:"created_by_member_id" json:"createdByMemberId"`
	Name                      string     `db:"name" json:"name"`
	Description               *string    `db:"description" json:"description"`
	StartsAt                  time.Time  `db:"starts_at" json:"startsAt"`
	DurationMinutes           *int       `db:"duration_minutes" json:"durationMinutes"`
	NotificationOffsetMinutes int        `db:"notification_offset_minutes" json:"notificationOffsetMinutes"`
	RecurrenceRule            *string    `db:"recurrence_rule" json:"recurrenceRule"`
	CreatedAt                 time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt                 time.Time  `db:"updated_at" json:"updatedAt"`
}

type CalendarOccurrence struct {
	Event    CalendarEvent `json:"event"`
	StartsAt time.Time     `json:"startsAt"`
	EndsAt   time.Time     `json:"endsAt"`
}

type CalendarService struct {
	calendarStore eventStore
}

type eventStore interface {
	InsertEvent(ctx context.Context, params InsertEventParams) error
	UpdateEvent(ctx context.Context, params UpdateEventStoreParams) error
	GetEvents(ctx context.Context, params GetEventsParams) ([]CalendarEvent, error)
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

type GetEventsParams struct {
	HouseholdID uuid.UUID
	From        time.Time
	To          time.Time
}
