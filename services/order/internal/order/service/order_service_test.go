package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/order/config"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/constant"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	productservice "github.com/alifmufthi91/ecommerce-system/services/order/external/product_service"
	productSvcMock "github.com/alifmufthi91/ecommerce-system/services/order/external/product_service/mocks"
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/order/external/warehouse_service"
	warehouseSvcMock "github.com/alifmufthi91/ecommerce-system/services/order/external/warehouse_service/mocks"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/payload"
	orderRepoMock "github.com/alifmufthi91/ecommerce-system/services/order/internal/order/repository/mocks"
	stockLockRepoMock "github.com/alifmufthi91/ecommerce-system/services/order/internal/order/repository/mocks"
)

func TestGetOrders_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		orderRepo *orderRepoMock.OrderRepository
	}

	tests := []struct {
		name  string
		req   payload.GetOrdersReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			req:  payload.GetOrdersReq{},
			setup: func(m dependencyMocks) {
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).
					Return([]model.Order{
						{
							ID:         uuid.New(),
							UserID:     uuid.New(),
							ProductID:  uuid.New(),
							Quantity:   2,
							TotalPrice: 100.0,
							Status:     "Pending",
						},
						{
							ID:         uuid.New(),
							UserID:     uuid.New(),
							ProductID:  uuid.New(),
							Quantity:   1,
							TotalPrice: 50.0,
							Status:     "Completed",
						},
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				orderRepo: orderRepoMock.NewOrderRepository(t),
			}
			orderSvc := orderService{
				orderRepo: mocks.orderRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := orderSvc.GetOrders(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resp))
			mocks.orderRepo.AssertExpectations(t)
		})
	}
}

func TestGetOrders_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		orderRepo *orderRepoMock.OrderRepository
	}

	tests := []struct {
		name  string
		req   payload.GetOrdersReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to get orders",
			req:  payload.GetOrdersReq{},
			setup: func(m dependencyMocks) {
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				orderRepo: orderRepoMock.NewOrderRepository(t),
			}
			orderSvc := orderService{
				orderRepo: mocks.orderRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := orderSvc.GetOrders(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Nil(t, resp)
			mocks.orderRepo.AssertExpectations(t)
		})
	}
}

func TestCreateOrder_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db            sqlmock.Sqlmock
		orderRepo     *orderRepoMock.OrderRepository
		stockLockRepo *stockLockRepoMock.StockLockRepository
		warehouseSvc  *warehouseSvcMock.IWarehouseSvc
		productSvc    *productSvcMock.IProductSvc
	}

	// Mock DB and transaction
	mockDB, err := pkg.SetupMockDB()
	assert.NoError(t, err)

	productID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name  string
		req   payload.CreateOrderReq
		setup func(m dependencyMocks)
	}{
		{
			name: "success",
			req: payload.CreateOrderReq{
				UserID:    userID.String(),
				ProductID: productID,
				Quantity:  2,
				Token:     "test-token",
			},
			setup: func(m dependencyMocks) {
				m.productSvc.On("GetProductByID", mock.Anything, productservice.GetProductByIDReq{
					ProductID: productID.String(),
					Token:     "test-token",
				}).Return(productservice.GetProductByIDResp{
					Data: productservice.GetProductByIDRespData{
						ID:    productID,
						Price: 50.0,
					},
				}, nil)

				m.db.ExpectBegin()
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithReturning").Return(m.orderRepo)
				m.orderRepo.On("CreateOrder", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						order := args.Get(1).(*model.Order)
						order.ID = uuid.New() // Simulate DB setting ID
					}).Return(nil)

				// Mock warehouse service
				m.warehouseSvc.On("ReserveStocks", mock.Anything, warehouseservice.ReserveStocksReq{
					Token: "test-token",
					Stocks: []warehouseservice.ReserveStocksReqData{
						{
							ProductID: productID.String(),
							Quantity:  2,
						},
					},
				}).Return(warehouseservice.ReserveStocksResp{
					Data: []warehouseservice.ReserveStocksRespData{
						{
							ProductID:        productID,
							ReservedQuantity: 2,
							WarehouseID:      uuid.New(),
						},
					},
				}, nil)

				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("CreateStockLock", mock.Anything, mock.Anything).
					Return(nil)
				m.db.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:            mockDB.Mock,
				orderRepo:     orderRepoMock.NewOrderRepository(t),
				stockLockRepo: stockLockRepoMock.NewStockLockRepository(t),
				productSvc:    productSvcMock.NewIProductSvc(t),
				warehouseSvc:  warehouseSvcMock.NewIWarehouseSvc(t),
			}

			orderSvc := orderService{
				db:            mockDB.Db,
				orderRepo:     mocks.orderRepo,
				stockLockRepo: mocks.stockLockRepo,
				productSvc:    mocks.productSvc,
				warehouseSvc:  mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			result, err := orderSvc.CreateOrder(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, result.ID)
			assert.Equal(t, productID, result.ProductID)
			assert.Equal(t, userID, result.UserID)
			assert.Equal(t, 2, result.Quantity)
			assert.Equal(t, 100.0, result.TotalPrice) // 50.0 * 2
			assert.Equal(t, constant.OrderStatusPending, result.Status)

			mocks.orderRepo.AssertExpectations(t)
			mocks.stockLockRepo.AssertExpectations(t)
			mocks.productSvc.AssertExpectations(t)
			mocks.warehouseSvc.AssertExpectations(t)
		})
	}
}

func TestCreateOrder_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		orderRepo     *orderRepoMock.OrderRepository
		stockLockRepo *stockLockRepoMock.StockLockRepository
		productSvc    *productSvcMock.IProductSvc
		warehouseSvc  *warehouseSvcMock.IWarehouseSvc
		db            sqlmock.Sqlmock
	}

	mockDB, err := pkg.SetupMockDB()
	assert.NoError(t, err)

	productID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		req     payload.CreateOrderReq
		setup   func(m dependencyMocks)
		wantErr string
	}{
		{
			name: "error - invalid user ID",
			req: payload.CreateOrderReq{
				UserID:    "invalid-uuid",
				ProductID: productID,
				Quantity:  2,
				Token:     "test-token",
			},
			setup: func(m dependencyMocks) {
				// No mocks needed as validation fails early
			},
			wantErr: "invalid user ID",
		},
		{
			name: "error - product service failure",
			req: payload.CreateOrderReq{
				UserID:    userID.String(),
				ProductID: productID,
				Quantity:  2,
				Token:     "test-token",
			},
			setup: func(m dependencyMocks) {
				m.productSvc.On("GetProductByID", mock.Anything, mock.Anything).
					Return(productservice.GetProductByIDResp{}, assert.AnError)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:            mockDB.Mock,
				orderRepo:     orderRepoMock.NewOrderRepository(t),
				stockLockRepo: stockLockRepoMock.NewStockLockRepository(t),
				productSvc:    productSvcMock.NewIProductSvc(t),
				warehouseSvc:  warehouseSvcMock.NewIWarehouseSvc(t),
			}

			orderSvc := orderService{
				db:            mockDB.Db,
				orderRepo:     mocks.orderRepo,
				stockLockRepo: mocks.stockLockRepo,
				productSvc:    mocks.productSvc,
				warehouseSvc:  mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			result, err := orderSvc.CreateOrder(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Equal(t, model.Order{}, result)
			if tt.wantErr != "" {
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCompleteOrder_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db            sqlmock.Sqlmock
		orderRepo     *orderRepoMock.OrderRepository
		stockLockRepo *stockLockRepoMock.StockLockRepository
		warehouseSvc  *warehouseSvcMock.IWarehouseSvc
	}

	// Mock DB and transaction
	mockDB, err := pkg.SetupMockDB()
	assert.NoError(t, err)

	orderID := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name  string
		req   payload.CompleteOrderReq
		setup func(m dependencyMocks)
	}{
		{
			name: "success",
			req: payload.CompleteOrderReq{
				OrderID: orderID.String(),
				Token:   "test-token",
			},
			setup: func(m dependencyMocks) {
				// Mock transaction
				m.db.ExpectBegin()

				// Mock get order by ID
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithLockForUpdate").Return(m.orderRepo)
				m.orderRepo.On("GetOrderByID", mock.Anything, orderID.String()).Return(model.Order{
					ID:         orderID,
					UserID:     uuid.New(),
					ProductID:  productID,
					Quantity:   2,
					TotalPrice: 100.0,
					Status:     constant.OrderStatusPending,
				}, nil)

				// Mock get stock locks
				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("WithLockForUpdate").Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    2,
						WarehouseID: warehouseID,
					},
				}, nil)

				// Mock warehouse service commit
				m.warehouseSvc.On("CommitReserves", mock.Anything, warehouseservice.CommitReservesReq{
					Token: "test-token",
					Stocks: []warehouseservice.CommitReservesReqData{
						{
							ProductID:   productID.String(),
							WarehouseID: warehouseID.String(),
							Quantity:    2,
						},
					},
				}).Return(warehouseservice.CommitReservesResp{}, nil)

				// Mock update order
				m.orderRepo.On("UpdateOrder", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.ID == orderID && order.Status == constant.OrderStatusCompleted
				})).Return(nil)

				m.db.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:            mockDB.Mock,
				orderRepo:     orderRepoMock.NewOrderRepository(t),
				stockLockRepo: stockLockRepoMock.NewStockLockRepository(t),
				warehouseSvc:  warehouseSvcMock.NewIWarehouseSvc(t),
			}

			orderSvc := orderService{
				db:            mockDB.Db,
				orderRepo:     mocks.orderRepo,
				stockLockRepo: mocks.stockLockRepo,
				warehouseSvc:  mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			result, err := orderSvc.CompleteOrder(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, orderID, result.ID)
			assert.Equal(t, constant.OrderStatusCompleted, result.Status)

			mocks.orderRepo.AssertExpectations(t)
			mocks.stockLockRepo.AssertExpectations(t)
			mocks.warehouseSvc.AssertExpectations(t)
		})
	}
}

func TestCompleteOrder_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		db            sqlmock.Sqlmock
		orderRepo     *orderRepoMock.OrderRepository
		stockLockRepo *stockLockRepoMock.StockLockRepository
		warehouseSvc  *warehouseSvcMock.IWarehouseSvc
	}

	mockDB, err := pkg.SetupMockDB()
	assert.NoError(t, err)

	orderID := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name    string
		req     payload.CompleteOrderReq
		setup   func(m dependencyMocks)
		wantErr string
	}{
		{
			name: "error - order not found",
			req: payload.CompleteOrderReq{
				OrderID: orderID.String(),
				Token:   "test-token",
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithLockForUpdate").Return(m.orderRepo)
				m.orderRepo.On("GetOrderByID", mock.Anything, orderID.String()).Return(model.Order{}, assert.AnError)
			},
			wantErr: "",
		},
		{
			name: "error - order not in pending status",
			req: payload.CompleteOrderReq{
				OrderID: orderID.String(),
				Token:   "test-token",
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithLockForUpdate").Return(m.orderRepo)
				m.orderRepo.On("GetOrderByID", mock.Anything, orderID.String()).Return(model.Order{
					ID:         orderID,
					UserID:     uuid.New(),
					ProductID:  productID,
					Quantity:   2,
					TotalPrice: 100.0,
					Status:     constant.OrderStatusCompleted, // Already completed
				}, nil)
			},
			wantErr: "order is not in pending status",
		},
		{
			name: "error - failed to get stock locks",
			req: payload.CompleteOrderReq{
				OrderID: orderID.String(),
				Token:   "test-token",
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithLockForUpdate").Return(m.orderRepo)
				m.orderRepo.On("GetOrderByID", mock.Anything, orderID.String()).Return(model.Order{
					ID:         orderID,
					UserID:     uuid.New(),
					ProductID:  productID,
					Quantity:   2,
					TotalPrice: 100.0,
					Status:     constant.OrderStatusPending,
				}, nil)

				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("WithLockForUpdate").Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).Return(nil, assert.AnError)
			},
			wantErr: "",
		},
		{
			name: "error - warehouse service commit failure",
			req: payload.CompleteOrderReq{
				OrderID: orderID.String(),
				Token:   "test-token",
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithLockForUpdate").Return(m.orderRepo)
				m.orderRepo.On("GetOrderByID", mock.Anything, orderID.String()).Return(model.Order{
					ID:         orderID,
					UserID:     uuid.New(),
					ProductID:  productID,
					Quantity:   2,
					TotalPrice: 100.0,
					Status:     constant.OrderStatusPending,
				}, nil)

				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("WithLockForUpdate").Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    2,
						WarehouseID: warehouseID,
					},
				}, nil)

				m.warehouseSvc.On("CommitReserves", mock.Anything, mock.Anything).
					Return(warehouseservice.CommitReservesResp{}, assert.AnError)
			},
			wantErr: "",
		},
		{
			name: "error - failed to update order",
			req: payload.CompleteOrderReq{
				OrderID: orderID.String(),
				Token:   "test-token",
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("WithLockForUpdate").Return(m.orderRepo)
				m.orderRepo.On("GetOrderByID", mock.Anything, orderID.String()).Return(model.Order{
					ID:         orderID,
					UserID:     uuid.New(),
					ProductID:  productID,
					Quantity:   2,
					TotalPrice: 100.0,
					Status:     constant.OrderStatusPending,
				}, nil)

				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("WithLockForUpdate").Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    2,
						WarehouseID: warehouseID,
					},
				}, nil)

				m.warehouseSvc.On("CommitReserves", mock.Anything, mock.Anything).
					Return(warehouseservice.CommitReservesResp{}, nil)

				m.orderRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:            mockDB.Mock,
				orderRepo:     orderRepoMock.NewOrderRepository(t),
				stockLockRepo: stockLockRepoMock.NewStockLockRepository(t),
				warehouseSvc:  warehouseSvcMock.NewIWarehouseSvc(t),
			}

			orderSvc := orderService{
				db:            mockDB.Db,
				orderRepo:     mocks.orderRepo,
				stockLockRepo: mocks.stockLockRepo,
				warehouseSvc:  mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			result, err := orderSvc.CompleteOrder(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Equal(t, model.Order{}, result)
			if tt.wantErr != "" {
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestProcessExpiredOrders_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db            sqlmock.Sqlmock
		orderRepo     *orderRepoMock.OrderRepository
		stockLockRepo *stockLockRepoMock.StockLockRepository
		warehouseSvc  *warehouseSvcMock.IWarehouseSvc
	}

	// Mock DB and transaction
	mockDB, err := pkg.SetupMockDB()
	assert.NoError(t, err)

	orderID1 := uuid.New()
	orderID2 := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()
	warehouseID1 := uuid.New()
	warehouseID2 := uuid.New()
	expiredTime := time.Now().Add(-time.Hour) // 1 hour ago

	tests := []struct {
		name  string
		setup func(m dependencyMocks)
	}{
		{
			name: "success - process multiple expired orders",
			setup: func(m dependencyMocks) {
				// Mock get expired orders
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).Return([]model.Order{
					{
						ID:         orderID1,
						UserID:     uuid.New(),
						ProductID:  productID1,
						Quantity:   2,
						TotalPrice: 100.0,
						Status:     constant.OrderStatusPending,
						ExpiresAt:  expiredTime,
					},
					{
						ID:         orderID2,
						UserID:     uuid.New(),
						ProductID:  productID2,
						Quantity:   1,
						TotalPrice: 50.0,
						Status:     constant.OrderStatusPending,
						ExpiresAt:  expiredTime,
					},
				}, nil)

				// Mock transaction for first order
				m.db.ExpectBegin()
				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID1.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID1,
						ProductID:   productID1,
						Quantity:    2,
						WarehouseID: warehouseID1,
					},
				}, nil)

				// Mock warehouse service rollback for first order
				m.warehouseSvc.On("RollbackReserves", mock.Anything, warehouseservice.RollbackReservesReq{
					Token: "static-token", // from config
					Stocks: []warehouseservice.RollbackReservesReqData{
						{
							ProductID:   productID1.String(),
							WarehouseID: warehouseID1.String(),
							Quantity:    2,
						},
					},
				}).Return(nil)

				// Mock update order for first order
				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("UpdateOrder", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.ID == orderID1 && order.Status == constant.OrderStatusCancelled
				})).Return(nil)
				m.db.ExpectCommit()

				// Mock transaction for second order
				m.db.ExpectBegin()
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID2.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID2,
						ProductID:   productID2,
						Quantity:    1,
						WarehouseID: warehouseID2,
					},
				}, nil)

				// Mock warehouse service rollback for second order
				m.warehouseSvc.On("RollbackReserves", mock.Anything, warehouseservice.RollbackReservesReq{
					Token: "static-token", // from config
					Stocks: []warehouseservice.RollbackReservesReqData{
						{
							ProductID:   productID2.String(),
							WarehouseID: warehouseID2.String(),
							Quantity:    1,
						},
					},
				}).Return(nil)

				// Mock update order for second order
				m.orderRepo.On("UpdateOrder", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
					return order.ID == orderID2 && order.Status == constant.OrderStatusCancelled
				})).Return(nil)
				m.db.ExpectCommit()
			},
		},
		{
			name: "success - no expired orders",
			setup: func(m dependencyMocks) {
				// Mock get expired orders returns empty
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).Return([]model.Order{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:            mockDB.Mock,
				orderRepo:     orderRepoMock.NewOrderRepository(t),
				stockLockRepo: stockLockRepoMock.NewStockLockRepository(t),
				warehouseSvc:  warehouseSvcMock.NewIWarehouseSvc(t),
			}

			orderSvc := orderService{
				config: &config.Config{
					External: config.External{
						WarehouseServiceStaticToken: "static-token",
					},
				},
				db:            mockDB.Db,
				orderRepo:     mocks.orderRepo,
				stockLockRepo: mocks.stockLockRepo,
				warehouseSvc:  mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			err := orderSvc.ProcessExpiredOrders(context.Background())

			// Then
			assert.NoError(t, err)
			mocks.orderRepo.AssertExpectations(t)
			mocks.stockLockRepo.AssertExpectations(t)
			mocks.warehouseSvc.AssertExpectations(t)
		})
	}
}

func TestProcessExpiredOrders_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		db            sqlmock.Sqlmock
		orderRepo     *orderRepoMock.OrderRepository
		stockLockRepo *stockLockRepoMock.StockLockRepository
		warehouseSvc  *warehouseSvcMock.IWarehouseSvc
	}

	mockDB, err := pkg.SetupMockDB()
	assert.NoError(t, err)

	orderID := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()
	expiredTime := time.Now().Add(-time.Hour)

	tests := []struct {
		name    string
		setup   func(m dependencyMocks)
		wantErr string
	}{
		{
			name: "error - failed to get expired orders",
			setup: func(m dependencyMocks) {
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
			wantErr: "",
		},
		{
			name: "error - failed to get stock locks",
			setup: func(m dependencyMocks) {
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).Return([]model.Order{
					{
						ID:         orderID,
						UserID:     uuid.New(),
						ProductID:  productID,
						Quantity:   2,
						TotalPrice: 100.0,
						Status:     constant.OrderStatusPending,
						ExpiresAt:  expiredTime,
					},
				}, nil)

				m.db.ExpectBegin()
				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).
					Return(nil, assert.AnError)
			},
			wantErr: "",
		},
		{
			name: "error - warehouse service rollback failure",
			setup: func(m dependencyMocks) {
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).Return([]model.Order{
					{
						ID:         orderID,
						UserID:     uuid.New(),
						ProductID:  productID,
						Quantity:   2,
						TotalPrice: 100.0,
						Status:     constant.OrderStatusPending,
						ExpiresAt:  expiredTime,
					},
				}, nil)

				m.db.ExpectBegin()
				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    2,
						WarehouseID: warehouseID,
					},
				}, nil)

				m.warehouseSvc.On("RollbackReserves", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
			wantErr: "",
		},
		{
			name: "error - failed to update order",
			setup: func(m dependencyMocks) {
				m.orderRepo.On("GetOrders", mock.Anything, mock.Anything).Return([]model.Order{
					{
						ID:         orderID,
						UserID:     uuid.New(),
						ProductID:  productID,
						Quantity:   2,
						TotalPrice: 100.0,
						Status:     constant.OrderStatusPending,
						ExpiresAt:  expiredTime,
					},
				}, nil)

				m.db.ExpectBegin()
				m.stockLockRepo.On("WithTX", mock.Anything).Return(m.stockLockRepo)
				m.stockLockRepo.On("GetStockLocksByOrderID", mock.Anything, orderID.String()).Return([]model.StockLock{
					{
						ID:          uuid.New(),
						OrderID:     orderID,
						ProductID:   productID,
						Quantity:    2,
						WarehouseID: warehouseID,
					},
				}, nil)

				m.warehouseSvc.On("RollbackReserves", mock.Anything, mock.Anything).Return(nil)

				m.orderRepo.On("WithTX", mock.Anything).Return(m.orderRepo)
				m.orderRepo.On("UpdateOrder", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:            mockDB.Mock,
				orderRepo:     orderRepoMock.NewOrderRepository(t),
				stockLockRepo: stockLockRepoMock.NewStockLockRepository(t),
				warehouseSvc:  warehouseSvcMock.NewIWarehouseSvc(t),
			}

			orderSvc := orderService{
				config: &config.Config{
					External: config.External{
						WarehouseServiceStaticToken: "static-token",
					},
				},
				db:            mockDB.Db,
				orderRepo:     mocks.orderRepo,
				stockLockRepo: mocks.stockLockRepo,
				warehouseSvc:  mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			err := orderSvc.ProcessExpiredOrders(context.Background())

			// Then
			assert.Error(t, err)
			if tt.wantErr != "" {
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
