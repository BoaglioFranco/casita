package schedules

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ScheduleID            uuid.UUID `db:"schedule_id"`
	ScheduleVersion       int       `db:"schedule_version"`
	Deleted               bool      `db:"deleted"`
	StartsAt              time.Time `db:"starts_at"`
	RecurrenceRule        *string   `db:"recurrence_rule"`
	ReminderOffsetMinutes int       `db:"reminder_offset_minutes"`
	NotificationID        uuid.UUID `db:"notification_id"`
}
