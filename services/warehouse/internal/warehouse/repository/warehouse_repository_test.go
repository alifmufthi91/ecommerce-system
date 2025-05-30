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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateWarehouse(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Warehouse)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name string
		data model.Warehouse
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Warehouse{
				Name:    "Test Warehouse",
				Address: "123 Warehouse St, City, Country",
				Status:  constant.WarehouseStatusActive,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Warehouse) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "warehouses" ("name","address","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						data.Address,
						data.Status,
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
			name: "error - failed to create warehouse",
			data: model.Warehouse{
				Name:    "Test Warehouse",
				Address: "123 Warehouse St, City, Country",
				Status:  constant.WarehouseStatusActive,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Warehouse) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "warehouses" ("name","address","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						data.Address,
						data.Status,
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

			repo := NewWarehouseRepository(mockDb.Db)

			err := repo.CreateWarehouse(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetWarehouses(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data []model.Warehouse)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		data    []model.Warehouse
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: []model.Warehouse{
				{Name: "Test Warehouse 1", Address: "123 Warehouse St, City, Country", Status: constant.WarehouseStatusActive},
				{Name: "Test Warehouse 2", Address: "456 Warehouse Ave, City, Country", Status: constant.WarehouseStatusActive},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Warehouse) {
					rows := sqlmock.NewRows([]string{"id", "name", "address", "status", "created_at", "updated_at"})
					for _, warehouse := range data {
						rows.AddRow(uuid.New(), warehouse.Name, warehouse.Address, warehouse.Status, time.Now(), time.Now())
					}
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "warehouses"`),
					).WillReturnRows(rows)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to get warehouses",
			data: []model.Warehouse{},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Warehouse) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "warehouses"`),
					).WillReturnError(sqlmock.ErrCancelled)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewWarehouseRepository(mockDb.Db)

			result, err := repo.GetWarehouses(context.Background())

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, len(tt.data), len(result))
			for i, warehouse := range result {
				assert.Equal(t, tt.data[i].Name, warehouse.Name)
			}
		})
	}
}

func TestGetWarehouseByID(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, id string, warehouse model.Warehouse)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name     string
		id       string
		sqlMock  sqlMock
		wantData model.Warehouse
		wantErr  bool
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, id string, warehouse model.Warehouse) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE id = $1 ORDER BY "warehouses"."id" LIMIT $2`),
					).WithArgs(id, 1).WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "address", "status", "created_at", "updated_at"}).
							AddRow(uuid.New(), warehouse.Name, warehouse.Address, warehouse.Status, time.Now(), time.Now()),
					)
				},
			},
			wantData: model.Warehouse{
				Name:    "Test Warehouse",
				Address: "123 Warehouse St, City, Country",
				Status:  constant.WarehouseStatusActive,
			},
			wantErr: false,
		},
		{
			name: "error - warehouse not found",
			id:   uuid.New().String(),
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, id string, warehouse model.Warehouse) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE id = $1 ORDER BY "warehouses"."id" LIMIT $2`),
					).WithArgs(id, 1).WillReturnError(
						gorm.ErrRecordNotFound,
					)
				},
			},
			wantData: model.Warehouse{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.id, tt.wantData)

			repo := NewWarehouseRepository(mockDb.Db)

			result, err := repo.GetWarehouseByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.wantData.Name, result.Name)
		})
	}
}

func TestUpdateWarehouse(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Warehouse)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		data    model.Warehouse
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Warehouse{
				ID:      uuid.New(),
				Name:    "Updated Warehouse",
				Address: "789 Updated St, City, Country",
				Status:  constant.WarehouseStatusActive,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Warehouse) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(`UPDATE "warehouses" SET "name"=$1,"address"=$2,"status"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6`),
					).WithArgs(
						data.Name,
						data.Address,
						data.Status,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						data.ID,
					).WillReturnResult(
						sqlmock.NewResult(1, 1),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to update warehouse",
			data: model.Warehouse{
				ID:      uuid.New(),
				Name:    "Updated Warehouse",
				Address: "789 Updated St, City, Country",
				Status:  constant.WarehouseStatusActive,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Warehouse) {
					mockDB.ExpectExec(
						regexp.QuoteMeta(`UPDATE "warehouses" SET "name"=$1,"address"=$2,"status"=$3,"created_at"=$4,"updated_at"=$5 WHERE "id" = $6`),
					).WithArgs(
						data.Name,
						data.Address,
						data.Status,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						data.ID,
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

			repo := NewWarehouseRepository(mockDb.Db)

			err := repo.UpdateWarehouse(context.Background(), &tt.data)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
