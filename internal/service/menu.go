package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/restaurant-api/internal/models"
	"github.com/yourusername/restaurant-api/internal/repository"
)

type MenuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
	}
}

func (s *MenuService) Create(ctx context.Context, item *models.MenuItem) error {
	return s.menuRepo.Create(ctx, item)
}

func (s *MenuService) List(ctx context.Context, restaurantID uuid.UUID) ([]*models.MenuItem, error) {
	return s.menuRepo.List(ctx, restaurantID)
}
