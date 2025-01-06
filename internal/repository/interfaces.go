package repository

import (
	"context"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/google/uuid"
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
	GetByID(ctx context.Context, id uuid.UUID) (*models.MenuItem, error)
	List(ctx context.Context, restaurantID uuid.UUID) ([]*models.MenuItem, error)
	Update(ctx context.Context, item *models.MenuItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*models.Order, error)
	Update(ctx context.Context, order *models.Order) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TableRepository interface {
	Create(ctx context.Context, table *models.Table) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Table, error)
	GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]*models.Table, error)
	GetByQRCode(ctx context.Context, qrCode string) (*models.Table, error)
	Update(ctx context.Context, table *models.Table) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TableStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
