package repository

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=StockLockRepository --case underscore
type StockLockRepository interface {
	WithTX(tx *gorm.DB) StockLockRepository
	WithReturning() StockLockRepository
	WithLockForUpdate() StockLockRepository
	CreateStockLock(ctx context.Context, stockLock *model.StockLock) error
	GetStockLocksByOrderID(ctx context.Context, orderID string) ([]model.StockLock, error)
}

type stockLockRepository struct {
	db *gorm.DB
}

func NewStockLockRepository(db *gorm.DB) StockLockRepository {
	return &stockLockRepository{db: db}
}

func (r *stockLockRepository) WithTX(tx *gorm.DB) StockLockRepository {
	if tx == nil {
		return r
	}
	return &stockLockRepository{db: tx}
}

func (r *stockLockRepository) WithReturning() StockLockRepository {
	return &stockLockRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *stockLockRepository) WithLockForUpdate() StockLockRepository {
	return &stockLockRepository{
		db: r.db.Clauses(clause.Locking{Strength: "UPDATE"}),
	}
}

func (r *stockLockRepository) CreateStockLock(ctx context.Context, stockLock *model.StockLock) error {
	ctx, span := observ.GetTracer().Start(ctx, "stockLockRepository.CreateStockLock")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&stockLock).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (r *stockLockRepository) GetStockLocksByOrderID(ctx context.Context, orderID string) ([]model.StockLock, error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockLockRepository.GetStockLocksByOrderID")
	defer span.End()

	var stockLocks []model.StockLock
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&stockLocks).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.NewWithCode(apperr.CodeSQLRead, "failed to get stock locks", err)
	}
	return stockLocks, nil
}
