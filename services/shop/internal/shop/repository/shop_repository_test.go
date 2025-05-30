package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateShop(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Shop)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name string
		data model.Shop
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Shop{
				Name: "Test Shop",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Shop) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "shops" ("name","address","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						sqlmock.AnyArg(),
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
			name: "error - failed to create shop",
			data: model.Shop{
				Name: "Test Shop",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Shop) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "shops" ("name","address","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						sqlmock.AnyArg(),
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

			repo := NewShopRepository(mockDb.Db)

			err := repo.CreateShop(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetShops(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data []model.Shop)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		data    []model.Shop
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: []model.Shop{
				{Name: "Test Shop 1"},
				{Name: "Test Shop 2"},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Shop) {
					rows := sqlmock.NewRows([]string{"id", "name", "address", "created_at", "updated_at"})
					for _, shop := range data {
						rows.AddRow(uuid.New(), shop.Name, nil, time.Now(), time.Now())
					}
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "shops"`),
					).WillReturnRows(rows)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to get shops",
			data: []model.Shop{},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Shop) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "shops"`),
					).WillReturnError(sqlmock.ErrCancelled)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewShopRepository(mockDb.Db)

			result, err := repo.GetShops(context.Background())

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, len(tt.data), len(result))
			for i, shop := range result {
				assert.Equal(t, tt.data[i].Name, shop.Name)
			}
		})
	}
}
