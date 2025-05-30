package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/product/config"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/payload"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct_ShouldReturnExpectedStatusCode(t *testing.T) {
	req := `{
		"name": "Item 1",
		"description": "description A",
		"price": 100,
		"shop_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	}`
	testScenarios := []struct {
		testName           string
		requestBody        string
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			requestBody:        req,
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
		},
		{
			testName:           "failed - error handle create product",
			requestBody:        req,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			requestBody:        `{"name": "Item 1", "description": "description A", "price": 100}`,
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

			mockProductSvc := &mocks.ProductService{}
			mockProductSvc.
				On("CreateProduct", mock.Anything, mock.Anything).
				Return(scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(scenario.requestBody))
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &productHandler{
				router:         r,
				config:         mockConfig,
				productService: mockProductSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestGetProducts_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		mockResult         []payload.GetProductsResp
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			statusCodeExpected: http.StatusOK,
			mockResult: []payload.GetProductsResp{
				{ID: uuid.New(), Name: "Test Product 1", AvailableStock: 10, Price: 100.0},
				{ID: uuid.New(), Name: "Test Product 2", AvailableStock: 5, Price: 50.0},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get products",
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

			mockProductSvc := &mocks.ProductService{}
			mockProductSvc.
				On("GetProducts", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/products", nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &productHandler{
				router:         r,
				config:         mockConfig,
				productService: mockProductSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestGetProductByID_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		productID          string
		mockResult         model.Product
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			productID:          uuid.New().String(),
			statusCodeExpected: http.StatusOK,
			mockResult: model.Product{
				ID:          uuid.New(),
				Name:        "Test Product",
				Description: "Test Description",
				Price:       100.0,
				ShopID:      uuid.New(),
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get product by ID",
			productID:          uuid.New().String(),
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid param",
			productID:          "invalid-uuid",
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

			mockProductSvc := &mocks.ProductService{}
			mockProductSvc.
				On("GetProductByID", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/products/"+scenario.productID, nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &productHandler{
				router:         r,
				config:         mockConfig,
				productService: mockProductSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}
