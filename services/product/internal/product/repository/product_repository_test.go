package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Product)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name string
		data model.Product
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Product{
				Name: "Test Product",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Product) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "products" ("shop_id","name","description","price","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.ShopID,
						data.Name,
						data.Description,
						data.Price,
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
			name: "error - failed to create product",
			data: model.Product{
				Name: "Test Product",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Product) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "products" ("shop_id","name","description","price","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.ShopID,
						data.Name,
						data.Description,
						data.Price,
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

			repo := NewProductRepository(mockDb.Db)

			err := repo.CreateProduct(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetProducts(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data []model.Product)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		data    []model.Product
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: []model.Product{
				{
					Name:        "Test Product 1",
					ShopID:      uuid.New(),
					Description: "Description for product 1",
					Price:       100.0,
				},
				{
					Name:        "Test Product 2",
					ShopID:      uuid.New(),
					Description: "Description for product 2",
					Price:       200.0,
				},
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Product) {
					rows := sqlmock.NewRows([]string{"id", "name", "shop_id", "description", "price", "created_at", "updated_at"})
					for _, product := range data {
						rows.AddRow(uuid.New(), product.Name, product.ShopID, product.Description, product.Price, time.Now(), time.Now())
					}
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "products"`),
					).WillReturnRows(rows)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to get products",
			data: []model.Product{},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data []model.Product) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "products"`),
					).WillReturnError(sqlmock.ErrCancelled)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewProductRepository(mockDb.Db)

			result, err := repo.GetProducts(context.Background())

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, len(tt.data), len(result))
			for i, product := range result {
				assert.Equal(t, tt.data[i].Name, product.Name)
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.Product)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name    string
		data    model.Product
		sqlMock sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.Product{
				ID:          uuid.New(),
				Name:        "Test Product",
				ShopID:      uuid.New(),
				Description: "Description for test product",
				Price:       100.0,
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Product) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`),
					).WithArgs(data.ID, 1).WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "shop_id", "description", "price", "created_at", "updated_at"}).
							AddRow(data.ID, data.Name, data.ShopID, data.Description, data.Price, time.Now(), time.Now()),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - product not found",
			data: model.Product{
				ID: uuid.New(),
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.Product) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 ORDER BY "products"."id" LIMIT $2`),
					).WithArgs(data.ID, 1).WillReturnError(gorm.ErrRecordNotFound)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewProductRepository(mockDb.Db)

			result, err := repo.GetProductByID(context.Background(), tt.data.ID.String())

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.data.Name, result.Name)
		})
	}
}
