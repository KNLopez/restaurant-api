package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TableStatus string

const (
	TableStatusAvailable TableStatus = "available"
	TableStatusOccupied  TableStatus = "occupied"
	TableStatusReserved  TableStatus = "reserved"
)

type Table struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	RestaurantID uuid.UUID   `json:"restaurant_id" db:"restaurant_id"`
	Number       int         `json:"number" db:"number"`
	Capacity     int         `json:"capacity" db:"capacity"`
	Status       TableStatus `json:"status" db:"status"`
	QRCode       string      `json:"qr_code" db:"qr_code"`
	TableURL     string      `json:"table_url" db:"table_url"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// GenerateTableURL creates the URL for the table's storefront
func (t *Table) GenerateTableURL(baseURL string) {
	t.TableURL = fmt.Sprintf("%s/restaurants/%s/tables/%s", baseURL, t.RestaurantID, t.ID)
}
