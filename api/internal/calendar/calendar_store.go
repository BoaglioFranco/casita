package calendar

import (
	"context"
	"fmt"

	"github.com/BoaglioFranco/casita/api/internal/database"
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
		    household_id, created_by_user_id, name, description, starts_at,
		    duration_minutes, recurrence_rule, notification_offset_minutes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := store.db.ExecContext(ctx, query,
		params.HouseholdID,
		params.CreatedByUserID,
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

func (store *CalendarStore) UpdateEvent(ctx context.Context, params UpdateEventStoreParams) error {
	query := `
		UPDATE calendar_events
		SET name = $1, description = $2, starts_at = $3,
		    duration_minutes = $4, notification_offset_minutes = $5,
		    recurrence_rule = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND household_id = $8`

	result, err := store.db.ExecContext(ctx, query,
		params.Name,
		params.Description,
		params.StartsAt.UTC(),
		params.DurationMinutes,
		params.NotificationOffsetMinutes,
		params.RecurrenceRule,
		params.EventID,
		params.HouseholdID,
	)
	if err != nil {
		return fmt.Errorf("unable to update calendar event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to check updated event rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (store *CalendarStore) DeleteEvent(ctx context.Context, params DeleteEventParams) error {
	query := `
		DELETE FROM calendar_events
		WHERE id = $1 AND household_id = $2`

	result, err := store.db.ExecContext(ctx, query, params.EventID, params.HouseholdID)
	if err != nil {
		return fmt.Errorf("unable to delete calendar event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to check deleted event rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (store *CalendarStore) GetEvents(ctx context.Context, params GetEventsParams) ([]CalendarEvent, error) {
	query := `
		SELECT * FROM calendar_events
		WHERE household_id = $1
		  AND (
		    (
		      recurrence_rule IS NULL
		      AND starts_at >= $2
		      AND starts_at <= $3
		    )
		    OR (
		      recurrence_rule IS NOT NULL
		      AND starts_at <= $3
		    )
		  )`
	events := make([]CalendarEvent, 0, 1)
	err := store.db.SelectContext(ctx, &events, query, params.HouseholdID, params.From, params.To)
	if err != nil {
		return nil, err
	}
	return events, nil
}
