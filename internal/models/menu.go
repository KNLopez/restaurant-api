package models

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID `json:"id" db:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id" db:"restaurant_id"`
	CategoryID   uuid.UUID `json:"category_id" db:"category_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Price        float64   `json:"price" db:"price"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
