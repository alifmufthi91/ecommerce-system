package service

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	stockRepoMock "github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/repository/mocks"
)

func TestGetStocks_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		stockRepo *stockRepoMock.StockRepository
	}

	tests := []struct {
		name  string
		req   payload.GetStocksReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			req: payload.GetStocksReq{
				WarehouseIDIN: []string{
					uuid.New().String(),
				},
				ProductIDIN: []string{
					uuid.New().String(),
				},
			},
			setup: func(m dependencyMocks) {
				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   uuid.New(),
							Quantity:    100,
							Reserved:    10,
						},
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   uuid.New(),
							Quantity:    200,
							Reserved:    20,
						},
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			stockSvc := stockService{
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := stockSvc.GetStocks(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resp))
			mocks.stockRepo.AssertExpectations(t)
		})
	}
}

func TestGetStocks_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		stockRepo *stockRepoMock.StockRepository
	}

	tests := []struct {
		name  string
		req   payload.GetStocksReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to get stocks",
			req: payload.GetStocksReq{
				WarehouseIDIN: []string{
					uuid.New().String(),
				},
				ProductIDIN: []string{
					uuid.New().String(),
				},
			},
			setup: func(m dependencyMocks) {
				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			stockSvc := stockService{
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := stockSvc.GetStocks(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Nil(t, resp)
			mocks.stockRepo.AssertExpectations(t)
		})
	}
}

func TestTransferStock_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	warehouseFromID := uuid.New()
	warehouseToID := uuid.New()
	warehouseProductID := uuid.New()

	tests := []struct {
		name  string
		req   payload.TransferStockReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success - exists into",
			req: payload.TransferStockReq{
				ProductID:       warehouseProductID,
				FromWarehouseID: warehouseFromID,
				ToWarehouseID:   warehouseToID,
				Quantity:        50,
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: warehouseFromID,
							ProductID:   warehouseProductID,
							Quantity:    100,
							Reserved:    10,
						},
						{
							ID:          uuid.New(),
							WarehouseID: warehouseToID,
							ProductID:   warehouseProductID,
							Quantity:    200,
							Reserved:    20,
						},
					}, nil)

				m.stockRepo.On("UpdateStock", mock.Anything, mock.Anything).
					Return(nil)

				m.stockRepo.On("CreateStockTransfer", mock.Anything, mock.Anything).
					Return(nil)

				m.db.ExpectCommit()
			},
		},
		{
			name: "success - not exists into",
			req: payload.TransferStockReq{
				ProductID:       warehouseProductID,
				FromWarehouseID: warehouseFromID,
				ToWarehouseID:   warehouseToID,
				Quantity:        50,
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: warehouseFromID,
							ProductID:   warehouseProductID,
							Quantity:    100,
							Reserved:    10,
						},
					}, nil)

				m.stockRepo.On("UpdateStock", mock.Anything, mock.Anything).
					Return(nil)

				m.stockRepo.On("CreateStock", mock.Anything, mock.Anything).
					Return(nil)

				m.stockRepo.On("CreateStockTransfer", mock.Anything, mock.Anything).
					Return(nil)

				m.db.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			err := stockSvc.TransferStock(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestTransferStock_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	warehouseFromID := uuid.New()
	warehouseToID := uuid.New()
	warehouseProductID := uuid.New()

	tests := []struct {
		name  string
		req   payload.TransferStockReq
		setup func(
			m dependencyMocks,
		)
		err string
	}{
		{
			name: "error - insufficient stock in from warehouse",
			req: payload.TransferStockReq{
				ProductID:       warehouseProductID,
				FromWarehouseID: warehouseFromID,
				ToWarehouseID:   warehouseToID,
				Quantity:        150,
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: warehouseFromID,
							ProductID:   warehouseProductID,
							Quantity:    100,
							Reserved:    10,
						},
					}, nil)

				m.db.ExpectRollback()
			},
		},
		{
			name: "error - from_warehouse_id and to_warehouse_id cannot be the same",
			req: payload.TransferStockReq{
				ProductID:       warehouseProductID,
				FromWarehouseID: warehouseFromID,
				ToWarehouseID:   warehouseFromID, // Same warehouse ID
				Quantity:        50,
			},
			setup: func(m dependencyMocks) {
				// No setup needed for this case, as it should return early
			},
			err: "from_warehouse_id and to_warehouse_id cannot be the same",
		},
		{
			name: "error - stock from warehouse_id not found",
			req: payload.TransferStockReq{
				ProductID:       warehouseProductID,
				FromWarehouseID: warehouseFromID,
				ToWarehouseID:   warehouseToID,
				Quantity:        50,
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{}, nil)

				m.db.ExpectRollback()
			},
			err: "stock from warehouse_id not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			if tt.setup != nil {
				tt.setup(mocks)
			}

			// When
			err := stockSvc.TransferStock(context.Background(), tt.req)

			// Then
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tt.err)
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestGetStockAvailablesByProduct_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		stockRepo *stockRepoMock.StockRepository
	}

	tests := []struct {
		name  string
		req   payload.GetStockAvailablesByProductReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			req: payload.GetStockAvailablesByProductReq{
				ProductIDIN: []string{
					uuid.New().String(),
				},
			},
			setup: func(m dependencyMocks) {
				m.stockRepo.On("GetAvailableStocksByProduct", mock.Anything, mock.Anything).
					Return([]model.GetStockAvailablesByProduct{
						{
							ProductID:      uuid.New(),
							AvailableStock: 100,
						},
						{
							ProductID:      uuid.New(),
							AvailableStock: 200,
						},
					}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			stockSvc := stockService{
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := stockSvc.GetStockAvailablesByProduct(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resp))
			mocks.stockRepo.AssertExpectations(t)
		})
	}
}

func TestGetStockAvailablesByProduct_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		stockRepo *stockRepoMock.StockRepository
	}

	tests := []struct {
		name  string
		req   payload.GetStockAvailablesByProductReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to get available stocks by product",
			req: payload.GetStockAvailablesByProductReq{
				ProductIDIN: []string{
					uuid.New().String(),
				},
			},
			setup: func(m dependencyMocks) {
				m.stockRepo.On("GetAvailableStocksByProduct", mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			stockSvc := stockService{
				stockRepo: mocks.stockRepo,
			}

			if tt.setup != nil {
				tt.setup(mocks)
			}

			// When
			resp, err := stockSvc.GetStockAvailablesByProduct(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Nil(t, resp)
			mocks.stockRepo.AssertExpectations(t)
		})
	}
}

func TestReserveStocks_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	productID := uuid.New()
	productID2 := uuid.New()

	tests := []struct {
		name  string
		req   payload.ReserveStocksReq
		setup func(
			m dependencyMocks,
			req payload.ReserveStocksReq,
		)
		expectedLen int
	}{
		{
			name: "success",
			req: payload.ReserveStocksReq{
				Stocks: []payload.ReserveStocksData{
					{
						ProductID: productID.String(),
						Quantity:  50,
					},
				},
			},
			setup: func(m dependencyMocks, req payload.ReserveStocksReq) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    100,
							Reserved:    10,
						},
					}, nil)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), mock.Anything, 0, req.Stocks[0].Quantity).
					Return(nil)

				m.db.ExpectCommit()
			},
			expectedLen: 1,
		},
		{
			name: "success - multiple sources",
			req: payload.ReserveStocksReq{
				Stocks: []payload.ReserveStocksData{
					{
						ProductID: productID.String(),
						Quantity:  50,
					},
				},
			},
			setup: func(m dependencyMocks, req payload.ReserveStocksReq) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    30,
							Reserved:    10,
						},
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    70,
							Reserved:    20,
						},
					}, nil)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), mock.Anything, 0, req.Stocks[0].Quantity-20).
					Return(nil)
				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), mock.Anything, 0, 20).
					Return(nil)

				m.db.ExpectCommit()
			},
			expectedLen: 2,
		},
		{
			name: "success - multiple products",
			req: payload.ReserveStocksReq{
				Stocks: []payload.ReserveStocksData{
					{
						ProductID: productID.String(),
						Quantity:  50,
					},
					{
						ProductID: productID2.String(),
						Quantity:  30,
					},
				},
			},
			setup: func(m dependencyMocks, req payload.ReserveStocksReq) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    30,
							Reserved:    10,
						},
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    50,
							Reserved:    10,
						},
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID2,
							Quantity:    50,
							Reserved:    5,
						},
					}, nil)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), mock.Anything, 0, req.Stocks[0].Quantity-20).
					Return(nil)
				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), mock.Anything, 0, 20).
					Return(nil)
				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID2.String(), mock.Anything, 0, req.Stocks[1].Quantity).
					Return(nil)

				m.db.ExpectCommit()
			},
			expectedLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			if tt.setup != nil {
				tt.setup(mocks, tt.req)
			}

			resp, err := stockSvc.ReserveStocks(context.Background(), tt.req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedLen, len(resp))
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestReserveStocks_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	productID := uuid.New()

	tests := []struct {
		name  string
		req   payload.ReserveStocksReq
		setup func(
			m dependencyMocks,
			req payload.ReserveStocksReq,
		)
		err string
	}{
		{
			name: "error - failed to get stocks",
			req: payload.ReserveStocksReq{
				Stocks: []payload.ReserveStocksData{
					{
						ProductID: productID.String(),
						Quantity:  50,
					},
				},
			},
			setup: func(m dependencyMocks, req payload.ReserveStocksReq) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return(nil, errors.New("failed to get stocks"))

				m.db.ExpectRollback()
			},
			err: "failed to get stocks",
		},
		{
			name: "error - insufficient stock",
			req: payload.ReserveStocksReq{
				Stocks: []payload.ReserveStocksData{
					{
						ProductID: productID.String(),
						Quantity:  150,
					},
				},
			},
			setup: func(m dependencyMocks, req payload.ReserveStocksReq) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    100,
							Reserved:    10,
						},
					}, nil)

				m.db.ExpectRollback()
			},
			err: "insufficient stock",
		},
		{
			name: "error - failed to add stock quantity and reserve quantity",
			req: payload.ReserveStocksReq{
				Stocks: []payload.ReserveStocksData{
					{
						ProductID: productID.String(),
						Quantity:  50,
					},
				},
			},
			setup: func(m dependencyMocks, req payload.ReserveStocksReq) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)
				m.stockRepo.On("WithLockForUpdate", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("GetStocks", mock.Anything, mock.Anything).
					Return([]model.WarehouseStock{
						{
							ID:          uuid.New(),
							WarehouseID: uuid.New(),
							ProductID:   productID,
							Quantity:    100,
							Reserved:    10,
						},
					}, nil)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), mock.Anything, 0, req.Stocks[0].Quantity).
					Return(errors.New("failed to add stock quantity and reserve quantity"))

				m.db.ExpectRollback()
			},
			err: "failed to add stock quantity and reserve quantity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			if tt.setup != nil {
				tt.setup(mocks, tt.req)
			}

			resp, err := stockSvc.ReserveStocks(context.Background(), tt.req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.err)

			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestRollbackReserves_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	productID := uuid.New()
	productID2 := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name  string
		req   payload.RollbackReservesReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success - single stock rollback",
			req: payload.RollbackReservesReq{
				Stocks: []payload.RollbackReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), 0, -20).
					Return(nil)

				m.db.ExpectCommit()
			},
		},
		{
			name: "success - multiple stocks rollback",
			req: payload.RollbackReservesReq{
				Stocks: []payload.RollbackReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
					{
						ProductID:   productID2.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    15,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), 0, -20).
					Return(nil)
				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID2.String(), warehouseID.String(), 0, -15).
					Return(nil)

				m.db.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			err := stockSvc.RollbackReserves(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestRollbackReserves_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	productID := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name  string
		req   payload.RollbackReservesReq
		setup func(
			m dependencyMocks,
		)
		err string
	}{
		{
			name: "error - failed to add stock quantity and reserve quantity",
			req: payload.RollbackReservesReq{
				Stocks: []payload.RollbackReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), 0, -20).
					Return(errors.New("failed to add stock quantity and reserve quantity"))

				m.db.ExpectRollback()
			},
			err: "failed to add stock quantity and reserve quantity",
		},
		{
			name: "error - transaction commit failure",
			req: payload.RollbackReservesReq{
				Stocks: []payload.RollbackReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), 0, -20).
					Return(nil)

				m.db.ExpectCommit().WillReturnError(errors.New("failed to commit transaction"))
			},
			err: "failed to commit transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given

			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			err := stockSvc.RollbackReserves(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.err)
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestCommitReserves_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	productID := uuid.New()
	productID2 := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name  string
		req   payload.CommitReservesReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success - single stock commit",
			req: payload.CommitReservesReq{
				Stocks: []payload.CommitReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), -20, -20).
					Return(nil)

				m.db.ExpectCommit()
			},
		},
		{
			name: "success - multiple stocks commit",
			req: payload.CommitReservesReq{
				Stocks: []payload.CommitReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
					{
						ProductID:   productID2.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    15,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), -20, -20).
					Return(nil)
				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID2.String(), warehouseID.String(), -15, -15).
					Return(nil)

				m.db.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			err := stockSvc.CommitReserves(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}

func TestCommitReserves_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		db        sqlmock.Sqlmock
		stockRepo *stockRepoMock.StockRepository
	}

	mockDb, err := pkg.SetupMockDB()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	productID := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name  string
		req   payload.CommitReservesReq
		setup func(
			m dependencyMocks,
		)
		err string
	}{
		{
			name: "error - failed to add stock quantity and reserve quantity",
			req: payload.CommitReservesReq{
				Stocks: []payload.CommitReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), -20, -20).
					Return(errors.New("failed to add stock quantity and reserve quantity"))

				m.db.ExpectRollback()
			},
			err: "failed to add stock quantity and reserve quantity",
		},
		{
			name: "error - transaction commit failure",
			req: payload.CommitReservesReq{
				Stocks: []payload.CommitReservesData{
					{
						ProductID:   productID.String(),
						WarehouseID: warehouseID.String(),
						Quantity:    20,
					},
				},
			},
			setup: func(m dependencyMocks) {
				m.db.ExpectBegin()

				m.stockRepo.On("WithTX", mock.Anything).
					Return(m.stockRepo)

				m.stockRepo.On("AddStockQtyAndReserveQty", mock.Anything, productID.String(), warehouseID.String(), -20, -20).
					Return(nil)

				m.db.ExpectCommit().WillReturnError(errors.New("failed to commit transaction"))
			},
			err: "failed to commit transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				db:        mockDb.Mock,
				stockRepo: stockRepoMock.NewStockRepository(t),
			}
			logger := pkg.InitLogger(&config.Config{})
			stockSvc := stockService{
				logger:    logger,
				db:        mockDb.Db,
				stockRepo: mocks.stockRepo,
			}

			tt.setup(mocks)

			// When
			err := stockSvc.CommitReserves(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.err)
			mocks.stockRepo.AssertExpectations(t)
			mockDb.Mock.ExpectationsWereMet()
		})
	}
}
