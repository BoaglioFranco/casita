package lists

import "github.com/google/uuid"

type ListService struct {
	listStore *ListStore
}

// Handlers

type CreateListPayload struct {
	Name string `json:"name" binding:"required"`
}

// Store
type InsertListParams struct {
	HouseholdID uuid.UUID
	Name        string
}
