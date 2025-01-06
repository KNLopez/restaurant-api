package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusAccepted OrderStatus = "accepted"
	OrderStatusReady    OrderStatus = "ready"
	OrderStatusComplete OrderStatus = "complete"
	OrderStatusCanceled OrderStatus = "canceled"
)

type Order struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	UserID       uuid.UUID   `json:"user_id" db:"user_id"`
	RestaurantID uuid.UUID   `json:"restaurant_id" db:"restaurant_id"`
	TableID      uuid.UUID   `json:"table_id" db:"table_id"`
	Status       OrderStatus `json:"status" db:"status"`
	TotalAmount  float64     `json:"total_amount" db:"total_amount"`
	Items        []OrderItem `json:"items"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
	ID         uuid.UUID `json:"id" db:"id"`
	OrderID    uuid.UUID `json:"order_id" db:"order_id"`
	MenuItemID uuid.UUID `json:"menu_item_id" db:"menu_item_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	Price      float64   `json:"price" db:"price"`
}
