package lists

import "github.com/BoaglioFranco/groceriesList/api/internal/database"

func NewListService(db *database.Database) *ListService {
	return &ListService{
		store: NewListStore(db),
	}
}
