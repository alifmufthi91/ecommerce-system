package repository

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/constant"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=StockRepository --case underscore
type StockRepository interface {
	WithTX(tx *gorm.DB) StockRepository
	WithReturning() StockRepository
	WithLockForUpdate() StockRepository
	CreateStock(ctx context.Context, stock *model.WarehouseStock) error
	CreateStockTransfer(ctx context.Context, transferStock *model.StockTransfer) error
	GetStocks(ctx context.Context, req payload.GetStocksReq) ([]model.WarehouseStock, error)
	UpdateStock(ctx context.Context, stock *model.WarehouseStock) error
	GetAvailableStocksByProduct(ctx context.Context, req payload.GetStockAvailablesByProductReq) ([]model.GetStockAvailablesByProduct, error)
	AddStockQtyAndReserveQty(ctx context.Context, productID string, warehouseID string, quantity int, reserved int) error
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) WithTX(tx *gorm.DB) StockRepository {
	if tx == nil {
		return r
	}
	return &stockRepository{db: tx}
}

func (r *stockRepository) WithReturning() StockRepository {
	return &stockRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *stockRepository) WithLockForUpdate() StockRepository {
	return &stockRepository{
		db: r.db.Clauses(clause.Locking{Strength: "UPDATE"}),
	}
}

func (r *stockRepository) CreateStock(ctx context.Context, stock *model.WarehouseStock) error {
	ctx, span := observ.GetTracer().Start(ctx, "stockRepository.CreateStock")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&stock).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to create stock")
	}
	return nil
}

func (r *stockRepository) CreateStockTransfer(ctx context.Context, transferStock *model.StockTransfer) error {
	ctx, span := observ.GetTracer().Start(ctx, "stockRepository.CreateStockTransfer")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&transferStock).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to create stock transfer")
	}
	return nil
}

func (r *stockRepository) GetStocks(ctx context.Context, req payload.GetStocksReq) ([]model.WarehouseStock, error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockRepository.GetStocks")
	defer span.End()

	stmt := r.db.WithContext(ctx)
	if len(req.WarehouseIDIN) > 0 {
		stmt = stmt.Where("warehouse_id IN ?", req.WarehouseIDIN)
	}

	if len(req.ProductIDIN) > 0 {
		stmt = stmt.Where("product_id IN ?", req.ProductIDIN)
	}

	stmt = stmt.Joins("JOIN warehouses w ON w.id = warehouse_id").
		Where("w.status = ?", constant.WarehouseStatusActive)

	var stocks []model.WarehouseStock
	if err := stmt.Find(&stocks).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to get stocks")
	}
	return stocks, nil
}

func (r *stockRepository) UpdateStock(ctx context.Context, stock *model.WarehouseStock) error {
	ctx, span := observ.GetTracer().Start(ctx, "stockRepository.UpdateStock")
	defer span.End()

	if err := r.db.WithContext(ctx).Save(&stock).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to update stock")
	}
	return nil
}

func (r *stockRepository) GetAvailableStocksByProduct(ctx context.Context, req payload.GetStockAvailablesByProductReq) ([]model.GetStockAvailablesByProduct, error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockRepository.GetAvailableStocksByProduct")
	defer span.End()

	stmt := r.db.WithContext(ctx)

	stmt = stmt.
		Table("warehouse_stocks ws").
		Joins("JOIN warehouses w ON ws.warehouse_id = w.id").
		Where("w.status = ?", constant.WarehouseStatusActive)

	stmt = stmt.Select("product_id", "sum(quantity) - sum(reserved) as available_stock").
		Group("product_id")

	if len(req.ProductIDIN) > 0 {
		stmt = stmt.Where("product_id IN ?", req.ProductIDIN)
	}

	var stocks []model.GetStockAvailablesByProduct
	if err := stmt.Find(&stocks).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to get available stocks by product")
	}
	return stocks, nil
}

func (r *stockRepository) AddStockQtyAndReserveQty(ctx context.Context, productID string, warehouseID string, quantity int, reserved int) error {
	ctx, span := observ.GetTracer().Start(ctx, "stockRepository.AddStockQtyAndReserveQty")
	defer span.End()

	if err := r.db.WithContext(ctx).Model(&model.WarehouseStock{}).
		Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
		Updates(map[string]any{
			"quantity": gorm.Expr("quantity + ?", quantity),
			"reserved": gorm.Expr("reserved + ?", reserved),
		}).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to update stock quantity and reserved quantity")
	}
	return nil
}
