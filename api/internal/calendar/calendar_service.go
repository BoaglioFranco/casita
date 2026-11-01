package calendar

import (
	"context"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

const defaultNotificationOffsetMinutes = 10

func NewCalendarService(db *database.Database) *CalendarService {
	return &CalendarService{
		calendarStore: NewCalendarStore(db),
	}
}

func (service *CalendarService) CreateEvent(ctx context.Context, params CreateEventParams) error {
	offset := defaultNotificationOffsetMinutes
	if params.NotificationOffsetMinutes != nil {
		offset = *params.NotificationOffsetMinutes
	}

	return service.calendarStore.InsertEvent(ctx, InsertEventParams{
		HouseholdID:               params.HouseholdID,
		CreatedByMemberID:         params.CreatedByMemberID,
		Name:                      params.Name,
		Description:               params.Description,
		StartsAt:                  params.StartsAt,
		DurationMinutes:           params.DurationMinutes,
		NotificationOffsetMinutes: offset,
		RecurrenceRule:            params.RecurrenceRule,
	})
}

func (service *CalendarService) UpdateEvent(ctx context.Context, params UpdateEventParams) error {
	offset := defaultNotificationOffsetMinutes
	if params.NotificationOffsetMinutes != nil {
		offset = *params.NotificationOffsetMinutes
	}

	return service.calendarStore.UpdateEvent(ctx, UpdateEventStoreParams{
		HouseholdID:               params.HouseholdID,
		EventID:                   params.EventID,
		Name:                      params.Name,
		Description:               params.Description,
		StartsAt:                  params.StartsAt,
		DurationMinutes:           params.DurationMinutes,
		NotificationOffsetMinutes: offset,
		RecurrenceRule:            params.RecurrenceRule,
	})
}
