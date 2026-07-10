package households

type InsertMemberParams struct {
	HouseholdID int64
	Name        string
}

type Member struct {
	ID          int64  `db:"id"`
	HouseholdID int64  `db:"household_id"`
	Name        string `db:"name"`
}
