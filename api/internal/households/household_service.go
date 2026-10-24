package households

import "github.com/BoaglioFranco/groceriesList/api/internal/database"

func NewHouseholdService(db *database.Database) *HouseholdService {
	return &HouseholdService{
		householdStore: NewHouseholdStore(db),
	}
}
