package lists

type ListService struct {
	listStore *ListStore
}

// Handlers

type CreateListPayload struct {
	Name string `json:"name" binding:"required"`
}

// Store
type InsertListParams struct {
	HouseholdID int64
	Name        string
}
