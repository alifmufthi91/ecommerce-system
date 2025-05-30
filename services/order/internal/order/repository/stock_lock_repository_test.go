package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateStockLock(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.StockLock)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	orderID := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name string
		data model.StockLock
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.StockLock{
				OrderID:     orderID,
				ProductID:   productID,
				WarehouseID: warehouseID,
				Quantity:    2,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.StockLock) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "stock_locks" ("order_id","product_id","warehouse_id","quantity","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.OrderID,
						data.ProductID,
						data.WarehouseID,
						data.Quantity,
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
			name: "error - failed to create stock_lock",
			data: model.StockLock{
				OrderID:     orderID,
				ProductID:   productID,
				Quantity:    2,
				WarehouseID: warehouseID,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.StockLock) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "stock_locks" ("order_id","product_id","warehouse_id","quantity","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.OrderID,
						data.ProductID,
						data.WarehouseID,
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

			repo := NewStockLockRepository(mockDb.Db)

			err := repo.CreateStockLock(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetStockLocksByOrderID(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	orderID := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()

	tests := []struct {
		name    string
		orderID uuid.UUID
		sqlMock func(mockDB sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name:    "success",
			orderID: orderID,
			sqlMock: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT * FROM "stock_locks" WHERE order_id = $1`,
					),
				).WithArgs(orderID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "order_id", "product_id", "warehouse_id", "quantity", "created_at", "updated_at"}).
						AddRow(uuid.New(), orderID, productID, warehouseID, 2, time.Now(), time.Now()),
				)
			},
			wantErr: false,
		},
		{
			name:    "error - failed to get stock locks by order ID",
			orderID: orderID,
			sqlMock: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT * FROM "stock_locks" WHERE order_id = $1`,
					),
				).WithArgs(orderID).WillReturnError(
					sqlmock.ErrCancelled,
				)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock(mockDb.Mock)

			repo := NewStockLockRepository(mockDb.Db)

			locks, err := repo.GetStockLocksByOrderID(context.Background(), tt.orderID.String())

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, 1, len(locks))
		})
	}
}
