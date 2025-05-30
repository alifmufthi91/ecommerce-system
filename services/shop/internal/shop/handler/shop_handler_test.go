package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/shop/config"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetShops_ShouldReturnExpectedStatusCode(t *testing.T) {
	testScenarios := []struct {
		testName           string
		mockResult         []model.Shop
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			statusCodeExpected: http.StatusOK,
			mockResult: []model.Shop{
				{Name: "Test Shop 1"},
				{Name: "Test Shop 2"},
			},
			mockError: nil,
		},
		{
			testName:           "failed - error handle get shops",
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

			mockShopSvc := &mocks.ShopService{}
			mockShopSvc.
				On("GetShops", mock.Anything).
				Return(scenario.mockResult, scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/shops", nil)
			ctx.Request.Header.Set("Authorization", "Bearer "+mockConfig.Token.JWTStatic)

			h := &shopHandler{
				router:      r,
				config:      mockConfig,
				shopService: mockShopSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}
