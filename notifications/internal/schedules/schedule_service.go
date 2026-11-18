package schedules

import "github.com/BoaglioFranco/casita/notifications/internal/database"

type ScheduleService struct {
	scheduleStore *ScheduleStore
}

func NewScheduleService(db *database.Database) *ScheduleService {
	return &ScheduleService{scheduleStore: NewScheduleStore(db)}
}
