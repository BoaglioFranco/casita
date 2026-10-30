package households

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrMemberNotFound = errors.New("member not found")

type HouseholdService struct {
	householdStore *HouseholdStore
}

type InsertMemberParams struct {
	HouseholdID uuid.UUID
	Name        string
	UserID      uuid.UUID
}

type Member struct {
	ID          uuid.UUID `db:"id" json:"id"`
	HouseholdID uuid.UUID `db:"household_id" json:"householdId"`
	Name        string    `db:"name" json:"name"`
	UserID      uuid.UUID `db:"user_id" json:"userId"`
}

type Household struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
