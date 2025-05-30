package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/constant"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateStock(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.WarehouseStock)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name string
		data model.WarehouseStock
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.WarehouseStock{
				WarehouseID: uuid.New(),
				ProductID:   uuid.New(),
				Quantity:    100,
				Reserved:    0,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.WarehouseStock) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "warehouse_stocks" ("warehouse_id","product_id","quantity","reserved","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.WarehouseID,
						data.ProductID,
						data.Quantity,
						data.Reserved,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnRows(
						sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to create stock",
			data: model.WarehouseStock{
				WarehouseID: uuid.New(),
				ProductID:   uuid.New(),
				Quantity:    100,
				Reserved:    0,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.WarehouseStock) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "warehouse_stocks" ("warehouse_id","product_id","quantity","reserved","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.WarehouseID,
						data.ProductID,
						data.Quantity,
						data.Reserved,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnError(
						sqlmock.ErrCancelled,
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewStockRepository(mockDb.Db)

			err := repo.CreateStock(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestCreateStockTransfer(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.StockTransfer)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name string
		data model.StockTransfer
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.StockTransfer{
				ProductID:       uuid.New(),
				FromWarehouseID: uuid.New(),
				ToWarehouseID:   uuid.New(),
				Quantity:        50,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.StockTransfer) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "stock_transfers" ("from_warehouse_id","to_warehouse_id","product_id","quantity","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.FromWarehouseID,
						data.ToWarehouseID,
						data.ProductID,
						data.Quantity,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnRows(
						sqlmock.NewRows([]string{"id"}).AddRow(data.ID),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to create stock transfer",
			data: model.StockTransfer{
				ProductID:       uuid.New(),
				FromWarehouseID: uuid.New(),
				ToWarehouseID:   uuid.New(),
				Quantity:        50,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.StockTransfer) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "stock_transfers" ("from_warehouse_id","to_warehouse_id","product_id","quantity","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.FromWarehouseID,
						data.ToWarehouseID,
						data.ProductID,
						data.Quantity,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnError(
						sqlmock.ErrCancelled,
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewStockRepository(mockDb.Db)

			err := repo.CreateStockTransfer(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetStocks(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, req payload.GetStocksReq)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		req     payload.GetStocksReq
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success - get stocks",
			req: payload.GetStocksReq{
				WarehouseIDIN: []string{uuid.New().String()},
				ProductIDIN:   []string{uuid.New().String()},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, req payload.GetStocksReq) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`SELECT "warehouse_stocks"."id","warehouse_stocks"."warehouse_id","warehouse_stocks"."product_id","warehouse_stocks"."quantity","warehouse_stocks"."reserved","warehouse_stocks"."created_at","warehouse_stocks"."updated_at" FROM "warehouse_stocks" JOIN warehouses w ON w.id = warehouse_id WHERE warehouse_id IN ($1) AND product_id IN ($2) AND w.status = $3`,
						),
					).WithArgs(req.WarehouseIDIN[0], req.ProductIDIN[0], constant.WarehouseStatusActive).WillReturnRows(
						sqlmock.NewRows([]string{"id", "warehouse_id", "product_id", "quantity", "reserved", "created_at", "updated_at"}).
							AddRow(uuid.New(), req.WarehouseIDIN[0], req.ProductIDIN[0], 100, 0, time.Now(), time.Now()),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to get stocks",
			req: payload.GetStocksReq{
				WarehouseIDIN: []string{uuid.New().String()},
				ProductIDIN:   []string{uuid.New().String()},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, req payload.GetStocksReq) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`SELECT "warehouse_stocks"."id","warehouse_stocks"."warehouse_id","warehouse_stocks"."product_id","warehouse_stocks"."quantity","warehouse_stocks"."reserved","warehouse_stocks"."created_at","warehouse_stocks"."updated_at" FROM "warehouse_stocks" JOIN warehouses w ON w.id = warehouse_id WHERE warehouse_id IN ($1) AND product_id IN ($2) AND w.status = $3`,
						),
					).WithArgs(req.WarehouseIDIN[0], req.ProductIDIN[0], constant.WarehouseStatusActive).WillReturnError(
						sqlmock.ErrCancelled,
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.req)

			repo := NewStockRepository(mockDb.Db)

			stocks, err := repo.GetStocks(context.Background(), tt.req)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.NotEmpty(t, stocks)
		})
	}
}

func TestUpdateStock(t *testing.T) {
	dmockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, stock model.WarehouseStock)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	tests := []struct {
		name    string
		stock   model.WarehouseStock
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success - update stock",
			stock: model.WarehouseStock{
				ID:          uuid.New(),
				WarehouseID: uuid.New(),
				ProductID:   uuid.New(),
				Quantity:    100,
				Reserved:    0,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, stock model.WarehouseStock) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(
							`UPDATE "warehouse_stocks" SET "warehouse_id"=$1,"product_id"=$2,"quantity"=$3,"reserved"=$4,"created_at"=$5,"updated_at"=$6 WHERE "id" = $7`,
						),
					).WithArgs(
						stock.WarehouseID,
						stock.ProductID,
						stock.Quantity,
						stock.Reserved,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						stock.ID,
					).WillReturnResult(
						sqlmock.NewResult(1, 1),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to update stock",
			stock: model.WarehouseStock{
				ID:          uuid.New(),
				WarehouseID: uuid.New(),
				ProductID:   uuid.New(),
				Quantity:    100,
				Reserved:    0,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, stock model.WarehouseStock) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(
							`UPDATE "warehouse_stocks" SET "warehouse_id"=$1,"product_id"=$2,"quantity"=$3,"reserved"=$4,"created_at"=$5,"updated_at"=$6 WHERE "id" = $7`,
						),
					).WithArgs(
						stock.WarehouseID,
						stock.ProductID,
						stock.Quantity,
						stock.Reserved,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						stock.ID,
					).WillReturnError(
						sqlmock.ErrCancelled,
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(dmockDb.Mock, tt.stock)

			repo := NewStockRepository(dmockDb.Db)

			err := repo.UpdateStock(context.Background(), &tt.stock)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetAvailableStocksByProduct(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, req payload.GetStockAvailablesByProductReq)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		req     payload.GetStockAvailablesByProductReq
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success - get available stocks by product",
			req: payload.GetStockAvailablesByProductReq{
				ProductIDIN: []string{uuid.New().String()},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, req payload.GetStockAvailablesByProductReq) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`SELECT "product_id",sum(quantity) - sum(reserved) as available_stock FROM warehouse_stocks ws JOIN warehouses w ON ws.warehouse_id = w.id WHERE w.status = $1 AND product_id IN ($2) GROUP BY "product_id"`,
						),
					).WithArgs(constant.WarehouseStatusActive, req.ProductIDIN[0]).WillReturnRows(
						sqlmock.NewRows([]string{"product_id", "available_stock"}).
							AddRow(req.ProductIDIN[0], 100),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to get available stocks by product",
			req: payload.GetStockAvailablesByProductReq{
				ProductIDIN: []string{uuid.New().String()},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, req payload.GetStockAvailablesByProductReq) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`SELECT "product_id",sum(quantity) - sum(reserved) as available_stock FROM warehouse_stocks ws JOIN warehouses w ON ws.warehouse_id = w.id WHERE w.status = $1 AND product_id IN ($2) GROUP BY "product_id"`,
						),
					).WithArgs(constant.WarehouseStatusActive, req.ProductIDIN[0]).WillReturnError(
						sqlmock.ErrCancelled,
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.req)

			repo := NewStockRepository(mockDb.Db)

			stocks, err := repo.GetAvailableStocksByProduct(context.Background(), tt.req)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.NotEmpty(t, stocks)
		})
	}
}

func TestAddStockQtyAndReserveQty(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, productID string, warehouseID string, quantity int, reserved int)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name        string
		productID   string
		warehouseID string
		quantity    int
		reserved    int
		sqlMock     sqlMock
		wantErr     bool
	}{
		{
			name:        "success - add stock quantity and reserve quantity",
			productID:   uuid.New().String(),
			warehouseID: uuid.New().String(),
			quantity:    50,
			reserved:    10,
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, productID string, warehouseID string, quantity int, reserved int) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(
							`UPDATE "warehouse_stocks" SET "quantity"=quantity + $1,"reserved"=reserved + $2,"updated_at"=$3 WHERE warehouse_id = $4 AND product_id = $5`,
						),
					).WithArgs(quantity, reserved, sqlmock.AnyArg(), warehouseID, productID).WillReturnResult(
						sqlmock.NewResult(1, 1),
					)
				},
			},
			wantErr: false,
		},
		{
			name:        "error - failed to add stock quantity and reserve quantity",
			productID:   uuid.New().String(),
			warehouseID: uuid.New().String(),
			quantity:    50,
			reserved:    10,
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, productID string, warehouseID string, quantity int, reserved int) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(
							`UPDATE "warehouse_stocks" SET "quantity"=quantity + $1,"reserved"=reserved + $2,"updated_at"=$3 WHERE warehouse_id = $4 AND product_id = $5`,
						),
					).WithArgs(quantity, reserved, sqlmock.AnyArg(), warehouseID, productID).WillReturnError(
						sqlmock.ErrCancelled,
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.productID, tt.warehouseID, tt.quantity, tt.reserved)

			repo := NewStockRepository(mockDb.Db)

			err := repo.AddStockQtyAndReserveQty(context.Background(), tt.productID, tt.warehouseID, tt.quantity, tt.reserved)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
