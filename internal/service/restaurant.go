package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/restaurant-api/internal/models"
	"github.com/yourusername/restaurant-api/internal/repository"
)

type RestaurantService struct {
	restaurantRepo repository.RestaurantRepository
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
	}
}

func (s *RestaurantService) Create(ctx context.Context, restaurant *models.Restaurant) error {
	return s.restaurantRepo.Create(ctx, restaurant)
}

func (s *RestaurantService) GetByID(ctx context.Context, id uuid.UUID) (*models.Restaurant, error) {
	return s.restaurantRepo.GetByID(ctx, id)
}
