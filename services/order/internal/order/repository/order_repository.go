package repository

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/payload"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=OrderRepository --case underscore
type OrderRepository interface {
	WithTX(tx *gorm.DB) OrderRepository
	WithReturning() OrderRepository
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrders(ctx context.Context, req payload.GetOrdersReq) ([]model.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	WithLockForUpdate() OrderRepository
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) WithTX(tx *gorm.DB) OrderRepository {
	if tx == nil {
		return r
	}
	return &orderRepository{db: tx}
}

func (r *orderRepository) WithReturning() OrderRepository {
	return &orderRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *orderRepository) WithLockForUpdate() OrderRepository {
	return &orderRepository{
		db: r.db.Clauses(clause.Locking{Strength: "UPDATE"}),
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	ctx, span := observ.GetTracer().Start(ctx, "orderRepository.CreateOrder")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&order).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.NewWithCode(apperr.CodeSQLCreate, "failed to create order", err)
	}
	return nil
}

func (r *orderRepository) GetOrders(ctx context.Context, req payload.GetOrdersReq) ([]model.Order, error) {
	ctx, span := observ.GetTracer().Start(ctx, "orderRepository.GetOrders")
	defer span.End()

	stmt := r.db.WithContext(ctx).Model(&model.Order{})
	if len(req.UserIDIN) > 0 {
		stmt = stmt.Where("user_id IN ?", req.UserIDIN)
	}

	if len(req.StatusIN) > 0 {
		stmt = stmt.Where("status IN ?", req.StatusIN)
	}

	if len(req.ProductIDIN) > 0 {
		stmt = stmt.Where("product_id IN ?", req.ProductIDIN)
	}

	if !req.ExpiresBefore.IsZero() {
		stmt = stmt.Where("expires_at < ?", req.ExpiresBefore)
	}

	var orders []model.Order
	if err := stmt.Find(&orders).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, apperr.NewWithCode(apperr.CodeSQLRead, "failed to get orders", err)
	}
	return orders, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	ctx, span := observ.GetTracer().Start(ctx, "orderRepository.GetOrderByID")
	defer span.End()

	var order model.Order
	if err := r.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		if err == gorm.ErrRecordNotFound {
			return model.Order{}, apperr.NewWithCode(apperr.CodeHTTPNotFound, "order not found", err)
		}
		return model.Order{}, apperr.NewWithCode(apperr.CodeSQLRead, "failed to get order", err)
	}
	return order, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	ctx, span := observ.GetTracer().Start(ctx, "orderRepository.UpdateOrder")
	defer span.End()

	if err := r.db.WithContext(ctx).Save(&order).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		return apperr.NewWithCode(apperr.CodeSQLUpdate, "failed to update order", err)
	}
	return nil
}
