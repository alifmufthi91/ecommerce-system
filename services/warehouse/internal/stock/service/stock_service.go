package service

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/repository"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

//go:generate mockery --name=StockService --case underscore
type StockService interface {
	GetStocks(ctx context.Context, req payload.GetStocksReq) ([]model.WarehouseStock, error)
	TransferStock(ctx context.Context, req payload.TransferStockReq) error
	GetStockAvailablesByProduct(ctx context.Context, req payload.GetStockAvailablesByProductReq) ([]model.GetStockAvailablesByProduct, error)
	ReserveStocks(ctx context.Context, req payload.ReserveStocksReq) ([]payload.ReserveStocksResp, error)
	RollbackReserves(ctx context.Context, req payload.RollbackReservesReq) error
	CommitReserves(ctx context.Context, req payload.CommitReservesReq) error
}

type stockService struct {
	config    *config.Config
	stockRepo repository.StockRepository
	db        *gorm.DB
}

func NewStockService(config *config.Config, db *gorm.DB, stockRepo repository.StockRepository) StockService {
	return &stockService{
		config:    config,
		stockRepo: stockRepo,
		db:        db,
	}
}

func (s *stockService) GetStocks(ctx context.Context, req payload.GetStocksReq) (result []model.WarehouseStock, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockService.TransferStock")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	stocks, err := s.stockRepo.GetStocks(ctx, req)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (s *stockService) TransferStock(ctx context.Context, req payload.TransferStockReq) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockService.TransferStock")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if req.FromWarehouseID == req.ToWarehouseID {
		return apperr.NewWithCode(apperr.CodeHTTPBadRequest, "from_warehouse_id and to_warehouse_id cannot be the same")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	stocks, err := s.stockRepo.WithTX(tx).WithLockForUpdate().GetStocks(ctx, payload.GetStocksReq{
		WarehouseIDIN: []string{req.FromWarehouseID.String(), req.ToWarehouseID.String()},
		ProductIDIN:   []string{req.ProductID.String()},
	})

	if err != nil {
		return err
	}

	var fromStock, toStock *model.WarehouseStock
	for _, stock := range stocks {
		if stock.WarehouseID == req.FromWarehouseID {
			fromStock = &stock
		}
		if stock.WarehouseID == req.ToWarehouseID {
			toStock = &stock
		}
	}

	if fromStock == nil {
		return apperr.NewWithCode(apperr.CodeHTTPNotFound, "stock from warehouse_id not found")
	}

	fromStock.Quantity -= req.Quantity
	if fromStock.Quantity < 0 || fromStock.Quantity < fromStock.Reserved {
		return apperr.NewWithCode(apperr.CodeHTTPBadRequest, "insufficient stock in from warehouse")
	}

	err = s.stockRepo.WithTX(tx).UpdateStock(ctx, fromStock)
	if err != nil {
		return err
	}

	if toStock == nil {
		toStock = &model.WarehouseStock{
			WarehouseID: req.ToWarehouseID,
			ProductID:   req.ProductID,
			Quantity:    req.Quantity,
		}
		err = s.stockRepo.WithTX(tx).CreateStock(ctx, toStock)
		if err != nil {
			return err
		}
	} else {
		toStock.Quantity += req.Quantity
		err = s.stockRepo.WithTX(tx).UpdateStock(ctx, toStock)
		if err != nil {
			return err
		}
	}

	err = s.stockRepo.WithTX(tx).CreateStockTransfer(ctx, &model.StockTransfer{
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
	})
	if err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to create stock transfer")
	}

	if err := tx.Commit().Error; err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
	}

	return nil
}

func (s *stockService) GetStockAvailablesByProduct(ctx context.Context, req payload.GetStockAvailablesByProductReq) (result []model.GetStockAvailablesByProduct, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockService.GetStockAvailablesByProduct")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	availableStocks, err := s.stockRepo.GetAvailableStocksByProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return availableStocks, nil
}

func (s *stockService) ReserveStocks(ctx context.Context, req payload.ReserveStocksReq) (result []payload.ReserveStocksResp, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockService.ReserveStocks")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	tx := s.db.Begin()
	defer tx.Rollback()

	var productIDs []string
	for _, stock := range req.Stocks {
		productIDs = append(productIDs, stock.ProductID)
	}

	stocks, err := s.stockRepo.WithTX(tx).WithLockForUpdate().GetStocks(ctx, payload.GetStocksReq{
		ProductIDIN: productIDs,
	})
	if err != nil {
		return result, err
	}

	var reservedStocks []model.WarehouseStock
	for _, stock := range req.Stocks {
		demandQty := stock.Quantity
		for _, s := range stocks {
			if s.ProductID.String() == stock.ProductID {
				reserveQty := min(demandQty, s.Quantity)
				if reserveQty == 0 {
					continue
				}
				demandQty -= reserveQty

				reservedStocks = append(reservedStocks, model.WarehouseStock{
					ProductID:   s.ProductID,
					WarehouseID: s.WarehouseID,
					Reserved:    reserveQty,
				})

				if demandQty == 0 {
					break
				}
			}
		}
		if demandQty > 0 {
			return result, apperr.NewWithCode(apperr.CodeHTTPBadRequest, "insufficient stock for product "+stock.ProductID)
		}
	}

	for _, reservedStock := range reservedStocks {
		err = s.stockRepo.WithTX(tx).AddStockQtyAndReserveQty(ctx, reservedStock.ProductID.String(), reservedStock.WarehouseID.String(), 0, reservedStock.Reserved)
		if err != nil {
			return result, err
		}

		result = append(result, payload.ReserveStocksResp{
			ProductID:        reservedStock.ProductID.String(),
			WarehouseID:      reservedStock.WarehouseID.String(),
			ReservedQuantity: reservedStock.Reserved,
		})
	}

	if err := tx.Commit().Error; err != nil {
		return result, apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
	}

	return result, nil
}

func (s *stockService) RollbackReserves(ctx context.Context, req payload.RollbackReservesReq) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockService.RollbackReserve")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	tx := s.db.Begin()
	defer tx.Rollback()

	for _, stock := range req.Stocks {
		err = s.stockRepo.WithTX(tx).AddStockQtyAndReserveQty(ctx, stock.ProductID, stock.WarehouseID, 0, -stock.Quantity)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
	}

	return nil
}

func (s *stockService) CommitReserves(ctx context.Context, req payload.CommitReservesReq) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "stockService.CommitReserves")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	tx := s.db.Begin()
	defer tx.Rollback()

	for _, stock := range req.Stocks {
		err = s.stockRepo.WithTX(tx).AddStockQtyAndReserveQty(ctx, stock.ProductID, stock.WarehouseID, -stock.Quantity, -stock.Quantity)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "failed to commit transaction")
	}

	return nil
}
