package service

import (
	"context"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/order/config"
	productservice "github.com/alifmufthi91/ecommerce-system/services/order/external/product_service"
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/order/external/warehouse_service"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/constant"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/payload"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/repository"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/observ"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

//go:generate mockery --name=OrderService --case underscore
type OrderService interface {
	GetOrders(ctx context.Context, req payload.GetOrdersReq) ([]model.Order, error)
	CreateOrder(ctx context.Context, req payload.CreateOrderReq) (model.Order, error)
	CompleteOrder(ctx context.Context, req payload.CompleteOrderReq) (model.Order, error)
	ProcessExpiredOrders(ctx context.Context) error
}

type orderService struct {
	config        *config.Config
	db            *gorm.DB
	logger        *pkg.Logger
	orderRepo     repository.OrderRepository
	stockLockRepo repository.StockLockRepository
	warehouseSvc  warehouseservice.IWarehouseSvc
	productSvc    productservice.IProductSvc
}

func NewOrderService(config *config.Config, db *gorm.DB, logger *pkg.Logger, orderRepo repository.OrderRepository, stockLockRepo repository.StockLockRepository, warehouseSvc warehouseservice.IWarehouseSvc, productSvc productservice.IProductSvc) OrderService {
	return &orderService{
		config:        config,
		db:            db,
		logger:        logger,
		orderRepo:     orderRepo,
		stockLockRepo: stockLockRepo,
		warehouseSvc:  warehouseSvc,
		productSvc:    productSvc,
	}
}

func (s *orderService) GetOrders(ctx context.Context, req payload.GetOrdersReq) ([]model.Order, error) {
	ctx, span := observ.GetTracer().Start(ctx, "orderService.GetOrders")
	defer span.End()

	orders, err := s.orderRepo.GetOrders(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	return orders, nil
}

func (s *orderService) CreateOrder(ctx context.Context, req payload.CreateOrderReq) (res model.Order, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "orderService.CreateOrder")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}
	}()

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		return model.Order{}, apperr.WrapWithCode(err, apperr.CodeHTTPBadRequest, "invalid user ID")
	}

	resp, err := s.productSvc.GetProductByID(ctx, productservice.GetProductByIDReq{
		ProductID: req.ProductID.String(),
		Token:     req.Token,
	})
	if err != nil {
		return model.Order{}, err
	}
	if resp.Data.ID == uuid.Nil {
		return model.Order{}, apperr.NewWithCode(apperr.CodeHTTPNotFound, "product not found")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	order := model.Order{
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		UserID:     userId,
		Status:     constant.OrderStatusPending,
		ExpiresAt:  time.Now().Add(time.Second * constant.OrderExpirationTime),
		TotalPrice: resp.Data.Price * float64(req.Quantity),
	}

	err = s.orderRepo.WithTX(tx).WithReturning().CreateOrder(ctx, &order)
	if err != nil {
		return model.Order{}, err
	}

	reservedStocks, err := s.warehouseSvc.ReserveStocks(ctx, warehouseservice.ReserveStocksReq{
		Token: req.Token,
		Stocks: []warehouseservice.ReserveStocksReqData{
			{
				ProductID: req.ProductID.String(),
				Quantity:  req.Quantity,
			},
		},
	})

	for _, stock := range reservedStocks.Data {
		err := s.stockLockRepo.WithTX(tx).CreateStockLock(ctx, &model.StockLock{
			OrderID:     order.ID,
			ProductID:   stock.ProductID,
			Quantity:    stock.ReservedQuantity,
			WarehouseID: stock.WarehouseID,
		})
		if err != nil {
			return model.Order{}, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to create stock lock")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return model.Order{}, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
	}

	return order, nil
}

func (s *orderService) CompleteOrder(ctx context.Context, req payload.CompleteOrderReq) (result model.Order, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "orderService.CompleteOrder")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	tx := s.db.Begin()
	defer tx.Rollback()

	order, err := s.orderRepo.WithTX(tx).WithLockForUpdate().GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return model.Order{}, err
	}

	if order.Status != constant.OrderStatusPending {
		return model.Order{}, apperr.NewWithCode(apperr.CodeHTTPBadRequest, "order is not in pending status")
	}

	stockLocks, err := s.stockLockRepo.WithTX(tx).WithLockForUpdate().GetStockLocksByOrderID(ctx, order.ID.String())
	if err != nil {
		return model.Order{}, err
	}

	var commitStockLocks []warehouseservice.CommitReservesReqData
	for _, stockLock := range stockLocks {
		commitStockLocks = append(commitStockLocks, warehouseservice.CommitReservesReqData{
			ProductID:   stockLock.ProductID.String(),
			WarehouseID: stockLock.WarehouseID.String(),
			Quantity:    stockLock.Quantity,
		})
	}

	_, err = s.warehouseSvc.CommitReserves(ctx, warehouseservice.CommitReservesReq{
		Token:  req.Token,
		Stocks: commitStockLocks,
	})
	if err != nil {
		return model.Order{}, err
	}

	order.Status = constant.OrderStatusCompleted
	if err := s.orderRepo.WithTX(tx).UpdateOrder(ctx, &order); err != nil {
		return model.Order{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.Order{}, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
	}

	return order, nil
}

func (s *orderService) ProcessExpiredOrders(ctx context.Context) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "orderService.ProcessExpiredOrders")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	expiredOrders, err := s.orderRepo.GetOrders(ctx, payload.GetOrdersReq{
		StatusIN:      []string{constant.OrderStatusPending},
		ExpiresBefore: time.Now(),
	})
	if err != nil {
		return err
	}

	for _, order := range expiredOrders {
		tx := s.db.Begin()
		defer tx.Rollback()

		stockLocks, err := s.stockLockRepo.WithTX(tx).GetStockLocksByOrderID(ctx, order.ID.String())
		if err != nil {
			return err
		}
		var rollbackStockLocks []warehouseservice.RollbackReservesReqData
		for _, stockLock := range stockLocks {
			rollbackStockLocks = append(rollbackStockLocks, warehouseservice.RollbackReservesReqData{
				ProductID:   stockLock.ProductID.String(),
				WarehouseID: stockLock.WarehouseID.String(),
				Quantity:    stockLock.Quantity,
			})
		}

		err = s.warehouseSvc.RollbackReserves(ctx, warehouseservice.RollbackReservesReq{
			Token:  s.config.External.WarehouseServiceStaticToken,
			Stocks: rollbackStockLocks,
		})
		if err != nil {
			return err
		}

		order.Status = constant.OrderStatusCancelled
		if err := s.orderRepo.WithTX(tx).UpdateOrder(ctx, &order); err != nil {
			return err
		}

		if err := tx.Commit().Error; err != nil {
			return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
		}
	}

	return nil
}
