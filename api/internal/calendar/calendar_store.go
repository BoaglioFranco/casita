package calendar

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/groceriesList/api/internal/database"
)

type CalendarStore struct {
	db *database.Database
}

func NewCalendarStore(db *database.Database) *CalendarStore {
	return &CalendarStore{db: db}
}

func (store *CalendarStore) InsertEvent(ctx context.Context, params InsertEventParams) error {
	query := `
		INSERT INTO calendar_events (
		    household_id, created_by_member_id, name, description, starts_at,
		    duration_minutes, recurrence_rule, notification_offset_minutes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := store.db.ExecContext(ctx, query,
		params.HouseholdID,
		params.CreatedByMemberID,
		params.Name,
		params.Description,
		params.StartsAt.UTC(),
		params.DurationMinutes,
		params.RecurrenceRule,
		params.NotificationOffsetMinutes,
	)
	if err != nil {
		return fmt.Errorf("unable to insert calendar event: %w", err)
	}

	return nil
}
