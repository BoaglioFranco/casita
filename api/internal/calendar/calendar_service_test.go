package calendar

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateEventTimeline(t *testing.T) {
	from := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 2)
	duration, rule := 30, "FREQ=DAILY"
	recurring := CalendarEvent{
		ID: uuid.New(), StartsAt: from.AddDate(0, 0, -10),
		DurationMinutes: &duration, RecurrenceRule: &rule,
	}
	oneTime := CalendarEvent{ID: uuid.New(), StartsAt: from.Add(time.Hour), DurationMinutes: &duration}
	service := &CalendarService{}
	occurrences, err := service.GenerateEventTimeline([]CalendarEvent{oneTime, recurring}, from, to)
	if err != nil {
		t.Fatal(err)
	}
	wantStarts := []time.Time{from, oneTime.StartsAt, from.AddDate(0, 0, 1), to}
	if len(occurrences) != len(wantStarts) {
		t.Fatalf("occurrence count = %d, want %d", len(occurrences), len(wantStarts))
	}
	for i, o := range occurrences {
		if !o.StartsAt.Equal(wantStarts[i]) || !o.EndsAt.Equal(wantStarts[i].Add(30*time.Minute)) {
			t.Fatalf("unexpected occurrence %d: %+v", i, o)
		}
		wantEvent := recurring
		if i == 1 {
			wantEvent = oneTime
		}
		if o.Event.ID != wantEvent.ID || !o.Event.StartsAt.Equal(wantEvent.StartsAt) {
			t.Fatalf("source event was not preserved: %+v", o.Event)
		}
	}
}

type eventStoreFunc func(context.Context, InsertEventParams) error

func (store eventStoreFunc) InsertEvent(ctx context.Context, params InsertEventParams) error {
	return store(ctx, params)
}

func (store eventStoreFunc) UpdateEvent(context.Context, UpdateEventStoreParams) error {
	panic("unexpected update")
}

func (store eventStoreFunc) GetEvents(context.Context, GetEventsParams) ([]CalendarEvent, error) {
	panic("unexpected get")
}

func TestValidateRecurrenceRule(t *testing.T) {
	valid, empty, invalid := "FREQ=DAILY;INTERVAL=2", "", "FREQ=INVALID"
	for _, tt := range []struct {
		name    string
		rule    *string
		invalid bool
	}{
		{"one-time event", nil, false},
		{"valid rule", &valid, false},
		{"empty rule", &empty, true},
		{"invalid frequency", &invalid, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecurrenceRule(tt.rule)
			if errors.Is(err, ErrInvalidRecurrenceRule) != tt.invalid {
				t.Fatalf("validation error = %v, want invalid = %v", err, tt.invalid)
			}
		})
	}
}

func TestCreateEventResolvesNotificationOffset(t *testing.T) {
	zero, custom := 0, 30
	for _, tt := range []struct {
		name   string
		offset *int
		want   int
	}{
		{name: "omitted", want: 10},
		{name: "explicit zero", offset: &zero, want: 0},
		{name: "custom offset", offset: &custom, want: 30},
	} {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			service := &CalendarService{
				calendarStore: eventStoreFunc(func(ctx context.Context, params InsertEventParams) error {
					calls++
					if params.NotificationOffsetMinutes != tt.want {
						t.Fatalf("offset = %d, want %d", params.NotificationOffsetMinutes, tt.want)
					}
					return nil
				}),
			}
			if err := service.CreateEvent(context.Background(), CreateEventParams{NotificationOffsetMinutes: tt.offset}); err != nil {
				t.Fatalf("CreateEvent: %v", err)
			}
			if calls != 1 {
				t.Fatalf("insert calls = %d, want 1", calls)
			}
		})
	}
}

func TestCreateEventPropagatesStoreError(t *testing.T) {
	wantErr := errors.New("insert failed")
	service := &CalendarService{
		calendarStore: eventStoreFunc(func(context.Context, InsertEventParams) error {
			return wantErr
		}),
	}
	if err := service.CreateEvent(context.Background(), CreateEventParams{}); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}
