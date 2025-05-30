package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/payload"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetStocks_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		queries            string
		mockResult         []model.WarehouseStock
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			queries:            "",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.WarehouseStock{
				{
					ID:          uuid.New(),
					WarehouseID: uuid.New(),
					ProductID:   uuid.New(),
					Quantity:    10,
					Reserved:    2,
				},
			},
			mockError: nil,
		},
		{
			testName:           "success with query",
			queries:            "?warehouse_id_in=some-warehouse-id&product_id_in=some-product-id",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.WarehouseStock{
				{
					ID:          uuid.New(),
					WarehouseID: uuid.New(),
					ProductID:   uuid.New(),
					Quantity:    2,
					Reserved:    1,
				},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get stocks",
			queries:            "",
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

			mockStockSvc := &mocks.StockService{}
			mockStockSvc.
				On("GetStocks", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/stocks"+scenario.queries, nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &stockHandler{
				router:       r,
				config:       mockConfig,
				stockService: mockStockSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestTransferStock_ShouldReturnExpectedStatusCode(t *testing.T) {
	payload := `{
		"from_warehouse_id": "8f1cc115-4434-4829-81c4-23fb01aa0dc0",
		"to_warehouse_id": "14c0374f-0fa3-4a02-baff-04e226910d3b",
		"product_id": "9a2b7c93-7c27-4e20-842f-24bf4df95bf0",
		"quantity": 1
	}`
	testScenarios := []struct {
		testName           string
		mockReq            string
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			mockReq:            payload,
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
		},
		{
			testName:           "failed - error handle transfer stock",
			mockReq:            payload,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			mockReq:            `{"from_warehouse_id": "invalid-uuid"}`,
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

			mockStockSvc := &mocks.StockService{}
			mockStockSvc.
				On("TransferStock", mock.Anything, mock.Anything).
				Return(scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/stocks/transfer", strings.NewReader(scenario.mockReq))
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &stockHandler{
				router:       r,
				config:       mockConfig,
				stockService: mockStockSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestGetAvailableStocksByProduct_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		queries            string
		mockResult         []model.GetStockAvailablesByProduct
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			queries:            "",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.GetStockAvailablesByProduct{
				{
					ProductID:      uuid.New(),
					AvailableStock: 10,
				},
			},
			mockError: nil,
		},
		{
			testName:           "success with query",
			queries:            "?product_id_in=some-product-id",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.GetStockAvailablesByProduct{
				{
					ProductID:      uuid.New(),
					AvailableStock: 5,
				},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get available stocks by product",
			queries:            "",
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

			mockStockSvc := &mocks.StockService{}
			if scenario.statusCodeExpected != http.StatusBadRequest {
				mockStockSvc.
					On("GetStockAvailablesByProduct", mock.Anything, mock.Anything).
					Return(scenario.mockResult, scenario.mockError)
			}

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/stocks/availables"+scenario.queries, nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &stockHandler{
				router:       r,
				config:       mockConfig,
				stockService: mockStockSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestReserveStocks_ShouldReturnExpectedStatusCode(t *testing.T) {
	payloadBody := `{
        "stocks": [
            {
                "product_id": "8f1cc115-4434-4829-81c4-23fb01aa0dc0",
                "quantity": 5
            }
        ]
    }`

	testScenarios := []struct {
		testName           string
		mockReq            string
		mockResult         []payload.ReserveStocksResp
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			mockReq:            payloadBody,
			statusCodeExpected: http.StatusOK,
			mockResult: []payload.ReserveStocksResp{
				{
					ProductID:        "8f1cc115-4434-4829-81c4-23fb01aa0dc0",
					WarehouseID:      uuid.New().String(),
					ReservedQuantity: 5,
				},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle reserve stocks",
			mockReq:            payloadBody,
			statusCodeExpected: http.StatusInternalServerError,
			mockResult:         nil,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			mockReq:            `{"stocks": [{"product_id": "invalid-uuid"}]}`,
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

			mockStockSvc := &mocks.StockService{}
			if scenario.statusCodeExpected != http.StatusBadRequest {
				mockStockSvc.
					On("ReserveStocks", mock.Anything, mock.Anything).
					Return(scenario.mockResult, scenario.mockError)
			}

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/stocks/reserve", strings.NewReader(scenario.mockReq))
			ctx.Request.Header.Set("Content-Type", "application/json")
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &stockHandler{
				router:       r,
				config:       mockConfig,
				stockService: mockStockSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestRollbackReserves_ShouldReturnExpectedStatusCode(t *testing.T) {
	payloadBody := `{
        "stocks": [
            {
                "product_id": "8f1cc115-4434-4829-81c4-23fb01aa0dc0",
                "warehouse_id": "14c0374f-0fa3-4a02-baff-04e226910d3b",
                "quantity": 5
            }
        ]
    }`

	testScenarios := []struct {
		testName           string
		mockReq            string
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			mockReq:            payloadBody,
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
		},
		{
			testName:           "failed - error handle rollback reserves",
			mockReq:            payloadBody,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			mockReq:            `{"stocks": [{"product_id": "invalid-uuid"}]}`,
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

			mockStockSvc := &mocks.StockService{}
			if scenario.statusCodeExpected != http.StatusBadRequest {
				mockStockSvc.
					On("RollbackReserves", mock.Anything, mock.Anything).
					Return(scenario.mockError)
			}

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/stocks/rollback", strings.NewReader(scenario.mockReq))
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &stockHandler{
				router:       r,
				config:       mockConfig,
				stockService: mockStockSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestCommitReserves_ShouldReturnExpectedStatusCode(t *testing.T) {
	payloadBody := `{
        "stocks": [
            {
                "product_id": "8f1cc115-4434-4829-81c4-23fb01aa0dc0",
                "warehouse_id": "14c0374f-0fa3-4a02-baff-04e226910d3b",
                "quantity": 5
            }
        ]
    }`

	testScenarios := []struct {
		testName           string
		mockReq            string
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			mockReq:            payloadBody,
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
		},
		{
			testName:           "failed - error handle commit reserves",
			mockReq:            payloadBody,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			mockReq:            `{"stocks": [{"product_id": "invalid-uuid"}]}`,
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

			mockStockSvc := &mocks.StockService{}
			if scenario.statusCodeExpected != http.StatusBadRequest {
				mockStockSvc.
					On("CommitReserves", mock.Anything, mock.Anything).
					Return(scenario.mockError)
			}

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/stocks/commit", strings.NewReader(scenario.mockReq))
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &stockHandler{
				router:       r,
				config:       mockConfig,
				stockService: mockStockSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}
