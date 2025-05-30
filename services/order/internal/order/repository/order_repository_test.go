package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/payload"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateOrder(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Order)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	userID := uuid.New()
	productID := uuid.New()

	tests := []struct {
		name string
		data model.Order
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Order{
				UserID:     userID,
				ProductID:  productID,
				Quantity:   2,
				TotalPrice: 100.0,
				ExpiresAt:  time.Now().Add(24 * time.Hour),
				Status:     "pending",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Order) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "orders" ("user_id","product_id","quantity","total_price","status","expires_at","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`,
						),
					).WithArgs(
						data.UserID,
						data.ProductID,
						data.Quantity,
						data.TotalPrice,
						data.Status,
						data.ExpiresAt,
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
			name: "error - failed to create order",
			data: model.Order{
				UserID:     userID,
				ProductID:  productID,
				Quantity:   2,
				TotalPrice: 100.0,
				ExpiresAt:  time.Now().Add(24 * time.Hour),
				Status:     "pending",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Order) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "orders" ("user_id","product_id","quantity","total_price","status","expires_at","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`,
						),
					).WithArgs(
						data.UserID,
						data.ProductID,
						data.Quantity,
						data.TotalPrice,
						data.Status,
						data.ExpiresAt,
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

			repo := NewOrderRepository(mockDb.Db)

			err := repo.CreateOrder(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetOrders(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	userID := uuid.New()
	productID := uuid.New()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data []model.Order)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		req     payload.GetOrdersReq
		data    []model.Order
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			req:  payload.GetOrdersReq{},
			data: []model.Order{
				{
					UserID:     uuid.New(),
					ProductID:  uuid.New(),
					Quantity:   1,
					TotalPrice: 50.0,
					Status:     "completed",
				},
				{
					UserID:     uuid.New(),
					ProductID:  uuid.New(),
					Quantity:   2,
					TotalPrice: 100.0,
					Status:     "pending",
				},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Order) {
					rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity", "total_price", "status", "created_at", "updated_at"})
					for _, order := range data {
						rows.AddRow(uuid.New(), order.UserID, order.ProductID, order.Quantity, order.TotalPrice, order.Status, time.Now(), time.Now())
					}
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "orders"`),
					).WillReturnRows(rows)
				},
			},
			wantErr: false,
		},
		{
			name: "success - all filters",
			req: payload.GetOrdersReq{
				UserIDIN:      []string{userID.String()},
				ProductIDIN:   []string{productID.String()},
				ExpiresBefore: time.Now().Add(24 * time.Hour),
				StatusIN:      []string{"completed"},
			},
			data: []model.Order{
				{
					UserID:     userID,
					ProductID:  productID,
					Quantity:   1,
					TotalPrice: 50.0,
					Status:     "completed",
				},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Order) {
					rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity", "total_price", "status", "created_at", "updated_at"})
					for _, order := range data {
						rows.AddRow(uuid.New(), order.UserID, order.ProductID, order.Quantity, order.TotalPrice, order.Status, time.Now(), time.Now())
					}
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "orders" WHERE user_id IN ($1) AND status IN ($2) AND product_id IN ($3) AND expires_at < $4`),
					).WithArgs(
						data[0].UserID,
						data[0].Status,
						data[0].ProductID,
						sqlmock.AnyArg(),
					).WillReturnRows(rows)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to get orders",
			req:  payload.GetOrdersReq{},
			data: []model.Order{},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Order) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "orders"`),
					).WillReturnError(sqlmock.ErrCancelled)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewOrderRepository(mockDb.Db)

			result, err := repo.GetOrders(context.Background(), tt.req)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, len(tt.data), len(result))
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	orderID := uuid.New()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Order)
	}

	tests := []struct {
		name    string
		orderID string
		data    model.Order
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name:    "success",
			orderID: orderID.String(),
			data: model.Order{
				ID:         orderID,
				UserID:     uuid.New(),
				ProductID:  uuid.New(),
				Quantity:   1,
				TotalPrice: 50.0,
				Status:     "completed",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Order) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`),
					).WithArgs(data.ID, 1).WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "product_id", "quantity", "total_price", "status", "created_at", "updated_at"}).
							AddRow(data.ID, data.UserID, data.ProductID, data.Quantity, data.TotalPrice, data.Status, time.Now(), time.Now()),
					)
				},
			},
			wantErr: false,
		},
		{
			name:    "error - order not found",
			orderID: orderID.String(),
			data:    model.Order{},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Order) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "orders" WHERE id = $1 ORDER BY "orders"."id" LIMIT $2`),
					).WithArgs(data.ID, 1).WillReturnError(gorm.ErrRecordNotFound)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			order := model.Order{ID: orderID}
			if tt.data.ID != uuid.Nil {
				order = tt.data
			}

			tt.sqlMock.Setup(mockDb.Mock, order)

			repo := NewOrderRepository(mockDb.Db)
			result, err := repo.GetOrderByID(context.Background(), tt.orderID)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.data.ID, result.ID)
		})
	}
}

func TestUpdateOrder(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	orderID := uuid.New()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Order)
	}

	tests := []struct {
		name    string
		data    model.Order
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Order{
				ID:         orderID,
				UserID:     uuid.New(),
				ProductID:  uuid.New(),
				Quantity:   2,
				TotalPrice: 100.0,
				Status:     "pending",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Order) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(`UPDATE "orders" SET "user_id"=$1,"product_id"=$2,"quantity"=$3,"total_price"=$4,"status"=$5,"expires_at"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`),
					).WithArgs(
						data.UserID,
						data.ProductID,
						data.Quantity,
						data.TotalPrice,
						data.Status,
						data.ExpiresAt,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						data.ID,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to update order",
			data: model.Order{
				ID:         orderID,
				UserID:     uuid.New(),
				ProductID:  uuid.New(),
				Quantity:   2,
				TotalPrice: 100.0,
				Status:     "pending",
				ExpiresAt:  time.Now().Add(24 * time.Hour),
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Order) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(`UPDATE "orders" SET "user_id"=$1,"product_id"=$2,"quantity"=$3,"total_price"=$4,"status"=$5,"expires_at"=$6,"created_at"=$7,"updated_at"=$8 WHERE "id" = $9`),
					).WithArgs(
						data.UserID,
						data.ProductID,
						data.Quantity,
						data.TotalPrice,
						data.Status,
						data.ExpiresAt,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						data.ID,
					).WillReturnError(sqlmock.ErrCancelled)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewOrderRepository(mockDb.Db)

			err := repo.UpdateOrder(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
