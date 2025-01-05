package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/restaurant-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *models.Restaurant) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Restaurant, error)
	GetByManagerID(ctx context.Context, managerID uuid.UUID) ([]*models.Restaurant, error)
	Update(ctx context.Context, restaurant *models.Restaurant) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MenuRepository interface {
	Create(ctx context.Context, item *models.MenuItem) error
	List(ctx context.Context, restaurantID uuid.UUID) ([]*models.MenuItem, error)
}
