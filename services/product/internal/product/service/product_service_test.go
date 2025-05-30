package service

import (
	"context"
	"testing"

	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/product/external/warehouse_service"
	warehouseSvcMock "github.com/alifmufthi91/ecommerce-system/services/product/external/warehouse_service/mocks"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/payload"
	productRepoMock "github.com/alifmufthi91/ecommerce-system/services/product/internal/product/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		productRepo *productRepoMock.ProductRepository
	}

	tests := []struct {
		name  string
		req   payload.CreateProductReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			req: payload.CreateProductReq{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100.0,
				ShopID:      uuid.New(),
			},
			setup: func(m dependencyMocks) {
				m.productRepo.On("CreateProduct", mock.Anything, mock.Anything).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				productRepo: productRepoMock.NewProductRepository(t),
			}
			productSvc := productService{
				productRepo: mocks.productRepo,
			}

			tt.setup(mocks)

			// When
			err := productSvc.CreateProduct(context.Background(), tt.req)

			// Then
			assert.NoError(t, err)
			mocks.productRepo.AssertExpectations(t)
		})
	}
}

func TestCreateProduct_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		productRepo *productRepoMock.ProductRepository
	}

	tests := []struct {
		name  string
		req   payload.CreateProductReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to create product",
			req: payload.CreateProductReq{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100.0,
				ShopID:      uuid.New(),
			},
			setup: func(m dependencyMocks) {
				m.productRepo.On("CreateProduct", mock.Anything, mock.Anything).
					Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				productRepo: productRepoMock.NewProductRepository(t),
			}
			productSvc := productService{
				productRepo: mocks.productRepo,
			}

			tt.setup(mocks)

			// When
			err := productSvc.CreateProduct(context.Background(), tt.req)

			// Then
			assert.Error(t, err)
			mocks.productRepo.AssertExpectations(t)
		})
	}
}

func TestGetProducts_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		productRepo  *productRepoMock.ProductRepository
		warehouseSvc *warehouseSvcMock.IWarehouseSvc
	}

	productID1 := uuid.New()
	productID2 := uuid.New()

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			setup: func(m dependencyMocks) {
				m.productRepo.On("GetProducts", mock.Anything).
					Return([]model.Product{
						{
							ID:          productID1,
							Name:        "Product One",
							ShopID:      uuid.New(),
							Description: "Description for product one",
							Price:       100.0,
						},
						{
							ID:          productID2,
							Name:        "Product Two",
							ShopID:      uuid.New(),
							Description: "Description for product two",
							Price:       200.0,
						},
					}, nil)

				m.warehouseSvc.On("GetStockAvailables", mock.Anything, mock.Anything).
					Return(warehouseservice.GetStockAvailablesResp{
						Data: []warehouseservice.GetStockAvailablesData{
							{
								ProductID:      productID1.String(),
								AvailableStock: 50,
							},
							{
								ProductID:      productID2.String(),
								AvailableStock: 30,
							},
						},
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				productRepo:  productRepoMock.NewProductRepository(t),
				warehouseSvc: warehouseSvcMock.NewIWarehouseSvc(t),
			}
			productSvc := productService{
				productRepo:  mocks.productRepo,
				warehouseSvc: mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			resp, err := productSvc.GetProducts(context.Background(), "test-token")

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resp))
			mocks.productRepo.AssertExpectations(t)
		})
	}
}

func TestGetProducts_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		productRepo  *productRepoMock.ProductRepository
		warehouseSvc *warehouseSvcMock.IWarehouseSvc
	}

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to get products",
			setup: func(m dependencyMocks) {
				m.productRepo.On("GetProducts", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
		{
			name: "error - failed to get stock availables",
			setup: func(m dependencyMocks) {
				m.productRepo.On("GetProducts", mock.Anything).
					Return([]model.Product{
						{
							ID:          uuid.New(),
							Name:        "Product One",
							ShopID:      uuid.New(),
							Description: "Description for product one",
							Price:       100.0,
						},
					}, nil)

				m.warehouseSvc.On("GetStockAvailables", mock.Anything, mock.Anything).
					Return(warehouseservice.GetStockAvailablesResp{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				productRepo:  productRepoMock.NewProductRepository(t),
				warehouseSvc: warehouseSvcMock.NewIWarehouseSvc(t),
			}
			productSvc := productService{
				productRepo:  mocks.productRepo,
				warehouseSvc: mocks.warehouseSvc,
			}

			tt.setup(mocks)

			// When
			resp, err := productSvc.GetProducts(context.Background(), "test-token")

			// Then
			assert.Error(t, err)
			assert.Nil(t, resp)
			mocks.productRepo.AssertExpectations(t)
		})
	}
}

func TestGetProductByID_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		productRepo *productRepoMock.ProductRepository
	}

	productID := uuid.New()

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			setup: func(m dependencyMocks) {
				m.productRepo.On("GetProductByID", mock.Anything, productID.String()).
					Return(model.Product{
						ID:          productID,
						Name:        "Test Product",
						Description: "Test Description",
						Price:       100.0,
						ShopID:      uuid.New(),
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				productRepo: productRepoMock.NewProductRepository(t),
			}
			productSvc := productService{
				productRepo: mocks.productRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := productSvc.GetProductByID(context.Background(), productID.String())

			// Then
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, productID, resp.ID)
			mocks.productRepo.AssertExpectations(t)
		})
	}
}

func TestGetProductByID_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		productRepo *productRepoMock.ProductRepository
	}

	productID := uuid.New()

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - product not found",
			setup: func(m dependencyMocks) {
				m.productRepo.On("GetProductByID", mock.Anything, productID.String()).
					Return(model.Product{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				productRepo: productRepoMock.NewProductRepository(t),
			}
			productSvc := productService{
				productRepo: mocks.productRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := productSvc.GetProductByID(context.Background(), productID.String())

			// Then
			assert.Error(t, err)
			assert.Equal(t, model.Product{}, resp)
			mocks.productRepo.AssertExpectations(t)
		})
	}
}
