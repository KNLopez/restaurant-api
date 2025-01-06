package service

import (
	"context"

	"github.com/KNLopez/restaurant-api/internal/models"
	"github.com/KNLopez/restaurant-api/internal/repository"
	"github.com/google/uuid"
)

type TableService struct {
	tableRepo repository.TableRepository
}

func NewTableService(tableRepo repository.TableRepository) *TableService {
	return &TableService{
		tableRepo: tableRepo,
	}
}

func (s *TableService) Create(ctx context.Context, table *models.Table) error {
	return s.tableRepo.Create(ctx, table)
}

func (s *TableService) GetByID(ctx context.Context, id uuid.UUID) (*models.Table, error) {
	return s.tableRepo.GetByID(ctx, id)
}

func (s *TableService) GetByQRCode(ctx context.Context, qrCode string) (*models.Table, error) {
	return s.tableRepo.GetByQRCode(ctx, qrCode)
}

func (s *TableService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TableStatus) error {
	return s.tableRepo.UpdateStatus(ctx, id, status)
}

func (s *TableService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.tableRepo.Delete(ctx, id)
}
