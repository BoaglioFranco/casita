package households

import "errors"

var ErrMemberNotFound = errors.New("member not found")

type HouseholdService struct {
	store *HouseholdStore
}

type InsertMemberParams struct {
	HouseholdID int64
	Name        string
}

type Member struct {
	ID          int64  `db:"id"`
	HouseholdID int64  `db:"household_id"`
	Name        string `db:"name"`
}
