package service

import (
	"context"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/model"
	shopRepoMock "github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/repository/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetShops_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		shopRepo *shopRepoMock.ShopRepository
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
				m.shopRepo.On("GetShops", mock.Anything).
					Return([]model.Shop{
						{
							ID:   uuid.New(),
							Name: "Shop One",
						},
						{
							ID:   uuid.New(),
							Name: "Shop Two",
						},
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				shopRepo: shopRepoMock.NewShopRepository(t),
			}
			shopSvc := shopService{
				shopRepo: mocks.shopRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := shopSvc.GetShops(context.Background())

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 2, len(resp))
			mocks.shopRepo.AssertExpectations(t)
		})
	}
}

func TestGetShops_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		shopRepo *shopRepoMock.ShopRepository
	}

	tests := []struct {
		name  string
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "error - failed to get shops",
			setup: func(m dependencyMocks) {
				m.shopRepo.On("GetShops", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				shopRepo: shopRepoMock.NewShopRepository(t),
			}
			shopSvc := shopService{
				shopRepo: mocks.shopRepo,
			}

			tt.setup(mocks)

			// When
			resp, err := shopSvc.GetShops(context.Background())

			// Then
			assert.Error(t, err)
			assert.Nil(t, resp)
			mocks.shopRepo.AssertExpectations(t)
		})
	}
}
