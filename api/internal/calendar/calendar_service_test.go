package calendar

import (
	"context"
	"errors"
	"testing"
)

type eventStoreFunc func(context.Context, InsertEventParams) error

func (store eventStoreFunc) InsertEvent(ctx context.Context, params InsertEventParams) error {
	return store(ctx, params)
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
