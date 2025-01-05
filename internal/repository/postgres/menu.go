package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/yourusername/restaurant-api/internal/models"
)

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{
		db: db,
	}
}

func (r *MenuRepository) Create(ctx context.Context, item *models.MenuItem) error {
	// TODO: Implement
	return nil
}

func (r *MenuRepository) List(ctx context.Context, restaurantID uuid.UUID) ([]*models.MenuItem, error) {
	// TODO: Implement
	return nil, nil
}
