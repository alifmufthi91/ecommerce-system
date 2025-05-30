package service

import (
	"context"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/constant"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/payload"
	warehouseRepoMock "github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWarehouses_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		warehouseRepo *warehouseRepoMock.WarehouseRepository
	}

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			setup: func(m dependencyMocks) {
				m.warehouseRepo.On("GetWarehouses", mock.Anything).
					Return([]model.Warehouse{
						{
							ID:   uuid.New(),
							Name: "Warehouse One",
						},
						{
							ID:   uuid.New(),
							Name: "Warehouse Two",
						},
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				warehouseRepo: warehouseRepoMock.NewWarehouseRepository(t),
			}
			warehouseSvc := warehouseService{
				warehouseRepo: mocks.warehouseRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := warehouseSvc.GetWarehouses(context.Background())

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resp))
			mocks.warehouseRepo.AssertExpectations(t)
		})
	}
}

func TestGetWarehouses_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		warehouseRepo *warehouseRepoMock.WarehouseRepository
	}

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to get warehouses",
			setup: func(m dependencyMocks) {
				m.warehouseRepo.On("GetWarehouses", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				warehouseRepo: warehouseRepoMock.NewWarehouseRepository(t),
			}
			warehouseSvc := warehouseService{
				warehouseRepo: mocks.warehouseRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := warehouseSvc.GetWarehouses(context.Background())

			// Then
			assert.Error(t, err)
			assert.Nil(t, resp)
			mocks.warehouseRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateWarehouse_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		warehouseRepo *warehouseRepoMock.WarehouseRepository
	}

	tests := []struct {
		name  string
		req   payload.UpdateWarehouseReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			req: payload.UpdateWarehouseReq{
				ID:     uuid.New(),
				Status: constant.WarehouseStatusInactive,
			},
			setup: func(m dependencyMocks) {
				m.warehouseRepo.On("GetWarehouseByID", mock.Anything, mock.AnythingOfType("string")).
					Return(model.Warehouse{
						ID:     uuid.New(),
						Name:   "Warehouse One",
						Status: constant.WarehouseStatusActive,
					}, nil)
				m.warehouseRepo.On("UpdateWarehouse", mock.Anything, mock.AnythingOfType("*model.Warehouse")).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				warehouseRepo: warehouseRepoMock.NewWarehouseRepository(t),
			}
			warehouseSvc := warehouseService{
				warehouseRepo: mocks.warehouseRepo,
			}

			tt.setup(mocks)

			err := warehouseSvc.UpdateWarehouse(context.Background(), tt.req)

			assert.NoError(t, err)
			mocks.warehouseRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateWarehouse_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		warehouseRepo *warehouseRepoMock.WarehouseRepository
	}

	tests := []struct {
		name  string
		req   payload.UpdateWarehouseReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - warehouse not found",
			req: payload.UpdateWarehouseReq{
				ID:     uuid.New(),
				Status: constant.WarehouseStatusInactive,
			},
			setup: func(m dependencyMocks) {
				m.warehouseRepo.On("GetWarehouseByID", mock.Anything, mock.AnythingOfType("string")).
					Return(model.Warehouse{}, assert.AnError)
			},
		},
		{
			name: "error - failed to update warehouse",
			req: payload.UpdateWarehouseReq{
				ID:     uuid.New(),
				Status: constant.WarehouseStatusInactive,
			},
			setup: func(m dependencyMocks) {
				m.warehouseRepo.On("GetWarehouseByID", mock.Anything, mock.AnythingOfType("string")).
					Return(model.Warehouse{
						ID:     uuid.New(),
						Name:   "Warehouse One",
						Status: constant.WarehouseStatusActive,
					}, nil)
				m.warehouseRepo.On("UpdateWarehouse", mock.Anything, mock.AnythingOfType("*model.Warehouse")).
					Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				warehouseRepo: warehouseRepoMock.NewWarehouseRepository(t),
			}
			warehouseSvc := warehouseService{
				warehouseRepo: mocks.warehouseRepo,
			}

			tt.setup(mocks)

			err := warehouseSvc.UpdateWarehouse(context.Background(), tt.req)

			assert.Error(t, err)
			mocks.warehouseRepo.AssertExpectations(t)
		})
	}
}
