package service

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/payload"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/repository"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockery --name=WarehouseService --case underscore
type WarehouseService interface {
	GetWarehouses(ctx context.Context) ([]model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, req payload.UpdateWarehouseReq) error
}

type warehouseService struct {
	config        *config.Config
	warehouseRepo repository.WarehouseRepository
}

func NewWarehouseService(config *config.Config, warehouseRepo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{
		config:        config,
		warehouseRepo: warehouseRepo,
	}
}

func (s *warehouseService) GetWarehouses(ctx context.Context) (result []model.Warehouse, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "warehouseService.GetWarehouses")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	warehouses, err := s.warehouseRepo.GetWarehouses(ctx)
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (s *warehouseService) UpdateWarehouse(ctx context.Context, req payload.UpdateWarehouseReq) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "warehouseService.UpdateWarehouse")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	warehouse, err := s.warehouseRepo.GetWarehouseByID(ctx, req.ID.String())
	if err != nil {
		return err
	}

	warehouse.Status = req.Status
	err = s.warehouseRepo.UpdateWarehouse(ctx, &warehouse)
	if err != nil {
		return err
	}

	return nil
}
