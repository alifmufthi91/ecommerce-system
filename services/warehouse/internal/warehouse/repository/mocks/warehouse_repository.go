// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"

	repository "github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/repository"
)

// WarehouseRepository is an autogenerated mock type for the WarehouseRepository type
type WarehouseRepository struct {
	mock.Mock
}

// CreateWarehouse provides a mock function with given fields: ctx, warehouse
func (_m *WarehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	ret := _m.Called(ctx, warehouse)

	if len(ret) == 0 {
		panic("no return value specified for CreateWarehouse")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Warehouse) error); ok {
		r0 = rf(ctx, warehouse)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetWarehouseByID provides a mock function with given fields: ctx, id
func (_m *WarehouseRepository) GetWarehouseByID(ctx context.Context, id string) (model.Warehouse, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetWarehouseByID")
	}

	var r0 model.Warehouse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.Warehouse, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.Warehouse); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.Warehouse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWarehouses provides a mock function with given fields: ctx
func (_m *WarehouseRepository) GetWarehouses(ctx context.Context) ([]model.Warehouse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetWarehouses")
	}

	var r0 []model.Warehouse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.Warehouse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.Warehouse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Warehouse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateWarehouse provides a mock function with given fields: ctx, warehouse
func (_m *WarehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	ret := _m.Called(ctx, warehouse)

	if len(ret) == 0 {
		panic("no return value specified for UpdateWarehouse")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Warehouse) error); ok {
		r0 = rf(ctx, warehouse)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithReturning provides a mock function with no fields
func (_m *WarehouseRepository) WithReturning() repository.WarehouseRepository {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for WithReturning")
	}

	var r0 repository.WarehouseRepository
	if rf, ok := ret.Get(0).(func() repository.WarehouseRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.WarehouseRepository)
		}
	}

	return r0
}

// WithTX provides a mock function with given fields: tx
func (_m *WarehouseRepository) WithTX(tx *gorm.DB) repository.WarehouseRepository {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTX")
	}

	var r0 repository.WarehouseRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) repository.WarehouseRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.WarehouseRepository)
		}
	}

	return r0
}

// NewWarehouseRepository creates a new instance of WarehouseRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWarehouseRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *WarehouseRepository {
	mock := &WarehouseRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
