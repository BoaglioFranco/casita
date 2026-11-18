package schedules

import "github.com/BoaglioFranco/casita/notifications/internal/database"

type ScheduleStore struct {
	db *database.Database
}

func NewScheduleStore(db *database.Database) *ScheduleStore {
	return &ScheduleStore{db: db}
}
