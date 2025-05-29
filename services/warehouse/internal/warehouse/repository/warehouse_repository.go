package repository

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=WarehouseRepository --case underscore
type WarehouseRepository interface {
	WithTX(tx *gorm.DB) WarehouseRepository
	WithReturning() WarehouseRepository
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	GetWarehouses(ctx context.Context) ([]model.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id string) (model.Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) WithTX(tx *gorm.DB) WarehouseRepository {
	if tx == nil {
		return r
	}
	return &warehouseRepository{db: tx}
}

func (r *warehouseRepository) WithReturning() WarehouseRepository {
	return &warehouseRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *warehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	ctx, span := observ.GetTracer().Start(ctx, "warehouseRepository.CreateWarehouse")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&warehouse).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to create warehouse")
	}
	return nil
}

func (r *warehouseRepository) GetWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	ctx, span := observ.GetTracer().Start(ctx, "warehouseRepository.GetWarehouses")
	defer span.End()

	var warehouses []model.Warehouse
	if err := r.db.WithContext(ctx).Find(&warehouses).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to get warehouses")
	}
	return warehouses, nil
}

func (r *warehouseRepository) GetWarehouseByID(ctx context.Context, id string) (model.Warehouse, error) {
	ctx, span := observ.GetTracer().Start(ctx, "warehouseRepository.GetWarehouseByID")
	defer span.End()

	var warehouse model.Warehouse
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&warehouse).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		if err == gorm.ErrRecordNotFound {
			return model.Warehouse{}, apperr.WrapWithCode(err, apperr.CodeHTTPNotFound, "warehouse not found")
		}
		return model.Warehouse{}, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to get warehouse by ID")
	}
	return warehouse, nil
}

func (r *warehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	ctx, span := observ.GetTracer().Start(ctx, "warehouseRepository.UpdateWarehouse")
	defer span.End()

	if err := r.db.WithContext(ctx).Save(&warehouse).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to update warehouse")
	}
	return nil
}
