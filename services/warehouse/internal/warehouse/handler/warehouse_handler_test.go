package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWarehouses_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		mockResult         []model.Warehouse
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.Warehouse{
				{Name: "Test Warehouse 1"},
				{Name: "Test Warehouse 2"},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get warehouses",
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			// Given
			r := pkg.GinTest()
			mockConfig := &config.Config{
				Token: config.Token{
					JWTSecret: []byte("secret"),
					JWTStatic: "static-token",
				},
			}

			mockWarehouseSvc := &mocks.WarehouseService{}
			mockWarehouseSvc.
				On("GetWarehouses", mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/warehouses", nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &warehouseHandler{
				router:           r,
				config:           mockConfig,
				warehouseService: mockWarehouseSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestUpdateWarehouse_ShouldReturnExpectedStatusCode(t *testing.T) {
	payload := `{
		"status": "active"
	}`
	testScenarios := []struct {
		testName           string
		mockParam          string
		mockRequest        string
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			mockParam:          uuid.NewString(),
			mockRequest:        payload,
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
		},
		{
			testName:           "failed - error handle update warehouse",
			mockParam:          uuid.NewString(),
			mockRequest:        payload,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid warehouse ID",
			mockParam:          "invalid-id",
			mockRequest:        payload,
			statusCodeExpected: http.StatusBadRequest,
		},
		{
			testName:           "failed - invalid request body",
			mockParam:          uuid.NewString(),
			mockRequest:        `{"status": 123}`, // Invalid status type
			statusCodeExpected: http.StatusBadRequest,
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			// Given
			r := pkg.GinTest()
			mockConfig := &config.Config{
				Token: config.Token{
					JWTSecret: []byte("secret"),
					JWTStatic: "static-token",
				},
			}

			mockWarehouseSvc := &mocks.WarehouseService{}
			mockWarehouseSvc.
				On("UpdateWarehouse", mock.Anything, mock.Anything).
				Return(scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			fmt.Println("Mock Param:", scenario.mockParam)
			fmt.Println("Mock Request Body:", scenario.mockRequest)
			ctx.Request = httptest.NewRequest(http.MethodPut, "/warehouses/"+scenario.mockParam, strings.NewReader(scenario.mockRequest))
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)
			ctx.Request.Header.Set("Content-Type", "application/json")

			h := &warehouseHandler{
				router:           r,
				config:           mockConfig,
				warehouseService: mockWarehouseSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}
