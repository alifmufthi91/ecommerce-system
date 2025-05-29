package repository

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/product/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=ProductRepository --case underscore
type ProductRepository interface {
	WithTX(tx *gorm.DB) ProductRepository
	WithReturning() ProductRepository
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProducts(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, productID string) (model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) WithTX(tx *gorm.DB) ProductRepository {
	if tx == nil {
		return r
	}
	return &productRepository{db: tx}
}

func (r *productRepository) WithReturning() ProductRepository {
	return &productRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	ctx, span := observ.GetTracer().Start(ctx, "productRepository.CreateProduct")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&product).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.WrapWithCode(err, apperr.CodeSQLCreate, "failed to create product")
	}
	return nil
}

func (r *productRepository) GetProducts(ctx context.Context) ([]model.Product, error) {
	ctx, span := observ.GetTracer().Start(ctx, "productRepository.GetProducts")
	defer span.End()

	var products []model.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.WrapWithCode(err, apperr.CodeSQLRead, "failed to get products")
	}
	return products, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID string) (model.Product, error) {
	ctx, span := observ.GetTracer().Start(ctx, "productRepository.GetProductByID")
	defer span.End()

	var product model.Product
	if err := r.db.WithContext(ctx).Where("id = ?", productID).First(&product).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		if err == gorm.ErrRecordNotFound {
			return model.Product{}, apperr.NewWithCode(apperr.CodeHTTPNotFound, "product not found")
		}
		return model.Product{}, apperr.WrapWithCode(err, apperr.CodeSQLRead, "failed to get product by ID")
	}
	return product, nil
}
