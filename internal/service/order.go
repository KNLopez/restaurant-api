package service

import (
	"context"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/KNLopez/restaurant-api/internal/repository"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (s *OrderService) Create(ctx context.Context, order *models.Order) error {
	return s.orderRepo.Create(ctx, order)
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	return s.orderRepo.UpdateStatus(ctx, id, status)
}

func (s *OrderService) Update(ctx context.Context, order *models.Order) error {
	return s.orderRepo.Update(ctx, order)
}

func (s *OrderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orderRepo.Delete(ctx, id)
}
