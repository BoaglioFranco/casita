package households

import (
	"errors"

	"github.com/google/uuid"
)

var ErrMemberNotFound = errors.New("member not found")

type HouseholdService struct {
	householdStore *HouseholdStore
}

type InsertMemberParams struct {
	HouseholdID uuid.UUID
	Name        string
}

type Member struct {
	ID          uuid.UUID `db:"id"`
	HouseholdID uuid.UUID `db:"household_id"`
	Name        string    `db:"name"`
}
