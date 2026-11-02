package calendar

import (
	"context"
	"slices"
	"time"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
	"github.com/BoaglioFranco/groceriesList/api/internal/utils"
	"github.com/teambition/rrule-go"
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

func (service *CalendarService) GenerateEventTimeline(events []CalendarEvent, from time.Time, to time.Time) (occurrences []CalendarOccurrence, err error) {
	isRecurring := func(e CalendarEvent) bool {
		return e.RecurrenceRule != nil && *e.RecurrenceRule != ""
	}
	recurringEvents, oneTimeEvents := utils.Partition(events, isRecurring)

	// Transform these to actual RRULEs
	// Generate occurrences for each rrule (within the timeframe)
	occurrences = make([]CalendarOccurrence, 0)
	for _, e := range recurringEvents {
		rule, err := rrule.StrToRRule(*e.RecurrenceRule)
		if err != nil {
			return nil, err
		}
		rule.DTStart(e.StartsAt.UTC())
		occurrenceDates := rule.Between(from, to, true)
		for _, d := range occurrenceDates {
			o := CalendarOccurrence{
				Event:    e,
				StartsAt: d,
				EndsAt:   d.Add(time.Minute * time.Duration(*e.DurationMinutes)),
			}
			occurrences = append(occurrences, o)
		}
	}

	for _, e := range oneTimeEvents {
		o := CalendarOccurrence{
			Event:    e,
			StartsAt: e.StartsAt,
			EndsAt:   e.StartsAt.Add(time.Minute * time.Duration(*e.DurationMinutes)),
		}
		occurrences = append(occurrences, o)
	}

	slices.SortFunc(occurrences, func(a CalendarOccurrence, b CalendarOccurrence) int {
		return a.StartsAt.Compare(b.StartsAt)
	})

	return occurrences, nil
}
