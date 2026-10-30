package lists

import (
	"time"

	"github.com/google/uuid"
)

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

type List struct {
	ID          uuid.UUID `db:"id" json:"id"`
	HouseholdID uuid.UUID `db:"household_id" json:"householdId"`
	Name        string    `db:"name" json:"name"`
	Icon        *string   `db:"icon" json:"icon"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}
