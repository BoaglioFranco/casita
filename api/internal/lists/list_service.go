package lists

import "github.com/BoaglioFranco/casita/api/internal/database"

func NewListService(db *database.Database) *ListService {
	return &ListService{
		listStore: NewListStore(db),
	}
}
