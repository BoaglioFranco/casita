package households

import "github.com/BoaglioFranco/groceriesList/api/internal/database"

func NewHouseholdService(db *database.Database) *HouseholdService {
	return &HouseholdService{
		store: NewStore(db),
	}
}
