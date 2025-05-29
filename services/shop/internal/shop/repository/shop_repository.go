package repository

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=ShopRepository --case underscore
type ShopRepository interface {
	WithTX(tx *gorm.DB) ShopRepository
	WithReturning() ShopRepository
	CreateShop(ctx context.Context, shop *model.Shop) error
	GetShops(ctx context.Context) ([]model.Shop, error)
}

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) WithTX(tx *gorm.DB) ShopRepository {
	if tx == nil {
		return r
	}
	return &shopRepository{db: tx}
}

func (r *shopRepository) WithReturning() ShopRepository {
	return &shopRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *shopRepository) CreateShop(ctx context.Context, shop *model.Shop) error {
	ctx, span := observ.GetTracer().Start(ctx, "shopRepository.CreateShop")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&shop).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeSQLCreate, "failed to create shop")
	}
	return nil
}

func (r *shopRepository) GetShops(ctx context.Context) ([]model.Shop, error) {
	ctx, span := observ.GetTracer().Start(ctx, "shopRepository.GetShops")
	defer span.End()

	var shops []model.Shop
	if err := r.db.WithContext(ctx).Find(&shops).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.WrapWithCode(err, apperr.CodeSQLRead, "failed to get shops")
	}
	return shops, nil
}
