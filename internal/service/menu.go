package service

import (
	"context"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/KNLopez/restaurant-api/internal/repository"
	"github.com/google/uuid"
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

func (s *MenuService) GetByID(ctx context.Context, id uuid.UUID) (*models.MenuItem, error) {
	return s.menuRepo.GetByID(ctx, id)
}

func (s *MenuService) Update(ctx context.Context, item *models.MenuItem) error {
	return s.menuRepo.Update(ctx, item)
}

func (s *MenuService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.menuRepo.Delete(ctx, id)
}
