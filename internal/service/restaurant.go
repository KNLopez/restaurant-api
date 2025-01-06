package service

import (
	"context"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/KNLopez/restaurant-api/internal/repository"
	"github.com/google/uuid"
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

func (s *RestaurantService) Update(ctx context.Context, restaurant *models.Restaurant) error {
	return s.restaurantRepo.Update(ctx, restaurant)
}

func (s *RestaurantService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.restaurantRepo.Delete(ctx, id)
}
