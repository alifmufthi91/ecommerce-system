package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/order/config"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/service/mocks"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrders_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		queries            string
		mockResult         []model.Order
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			queries:            "",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.Order{
				{ID: uuid.New(),
					Status:     "pending",
					UserID:     uuid.New(),
					ProductID:  uuid.New(),
					Quantity:   2,
					TotalPrice: 100.0,
				},
			},
			mockError: nil,
		},
		{
			testName:           "success with query",
			queries:            "?status=pending",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.Order{
				{ID: uuid.New(),
					Status:     "pending",
					UserID:     uuid.New(),
					ProductID:  uuid.New(),
					Quantity:   2,
					TotalPrice: 100.0,
				},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get orders",
			queries:            "",
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - error binding query",
			queries:            "?expires_before=invalid-date",
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

			mockOrderSvc := &mocks.OrderService{}
			mockOrderSvc.
				On("GetOrders", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/orders"+scenario.queries, nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &orderHandler{
				router:       r,
				config:       mockConfig,
				orderService: mockOrderSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestCreateOrder_ShouldReturnExpectedStatusCode(t *testing.T) {
	payload := `{
		"product_id": "9a2b7c93-7c27-4e20-842f-24bf4df95bf0",
		"quantity": 3
	}`

	testScenarios := []struct {
		testName           string
		mockRequest        string
		mockResult         model.Order
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			statusCodeExpected: http.StatusOK,
			mockRequest:        payload,
			mockError:          nil,
			mockResult: model.Order{
				ID: uuid.New(),
			},
		},
		{
			testName:           "failed - error handle create",
			mockRequest:        payload,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
			mockResult:         model.Order{},
		},
		{
			testName: "failed - error binding payload",
			mockRequest: `{
				"product_id": "",
				"quantity": 0
			}`,
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

			mockOrderSvc := &mocks.OrderService{}
			mockOrderSvc.
				On("CreateOrder", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(scenario.mockRequest))
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &orderHandler{
				router:       r,
				config:       mockConfig,
				orderService: mockOrderSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestCompleteOrder_ShouldReturnExpectedStatusCode(t *testing.T) {
	orderID := uuid.New()

	testScenarios := []struct {
		testName           string
		orderID            string
		mockResult         model.Order
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			orderID:            orderID.String(),
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
			mockResult: model.Order{
				ID:         orderID,
				Status:     "completed",
				UserID:     uuid.New(),
				ProductID:  uuid.New(),
				Quantity:   2,
				TotalPrice: 100.0,
			},
		},
		{
			testName:           "failed - error handle complete order",
			orderID:            orderID.String(),
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
			mockResult:         model.Order{},
		},
		{
			testName:           "failed - missing order ID",
			orderID:            "",
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

			mockOrderSvc := &mocks.OrderService{}
			mockOrderSvc.
				On("CompleteOrder", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)

			// Build URL with order ID parameter
			url := "/orders/" + scenario.orderID + "/complete"
			ctx.Request = httptest.NewRequest(http.MethodPatch, url, nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &orderHandler{
				router:       r,
				config:       mockConfig,
				orderService: mockOrderSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}
