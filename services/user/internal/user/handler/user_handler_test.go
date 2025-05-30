package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/user/config"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser_ShouldReturnExpectedStatusCode(t *testing.T) {
	req := `{
		"name": "Alex",
		"email": "alex@mail.com",
		"phone": "+628512332112",
		"password": "TestPassword123"
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
			testName:           "failed - error handle create user",
			requestBody:        req,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			requestBody:        `{"name": "Item 1"}`,
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

			mockUserSvc := &mocks.UserService{}
			mockUserSvc.
				On("RegisterUser", mock.Anything, mock.Anything).
				Return(scenario.mockError)

			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(scenario.requestBody))

			h := &userHandler{
				router:      r,
				config:      mockConfig,
				userService: mockUserSvc,
			}
			h.RegisterRoutes(r.Group(""))

			// When
			r.ServeHTTP(rr, ctx.Request)

			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}

func TestLoginUser_ShouldReturnExpectedStatusCode(t *testing.T) {
	req := `{
		"email_or_phone": "test@example.com",
		"password": "password123"
	}`
	testScenarios := []struct {
		testName           string
		requestBody        string
		mockResult         string
		mockError          error
		statusCodeExpected int
	}{
		{
			testName:           "success",
			requestBody:        req,
			mockResult:         "mocked-token",
			statusCodeExpected: http.StatusOK,
			mockError:          nil,
		},
		{
			testName:           "failed - error handle login user",
			requestBody:        req,
			statusCodeExpected: http.StatusInternalServerError,
			mockError:          errors.New("something went wrong"),
		},
		{
			testName:           "failed - invalid request body",
			requestBody:        `{"emailOrPhone": ""}`,
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
			mockUserSvc := &mocks.UserService{}
			mockUserSvc.
				On("LoginUser", mock.Anything, mock.Anything).
				Return(scenario.mockResult, scenario.mockError)
			rr := httptest.NewRecorder()
			ctx := pkg.GetTestGinContext(rr)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(scenario.requestBody))

			h := &userHandler{
				router:      r,
				config:      mockConfig,
				userService: mockUserSvc,
			}
			h.RegisterRoutes(r.Group(""))
			// When
			r.ServeHTTP(rr, ctx.Request)
			// Then
			assert.Equal(t, scenario.statusCodeExpected, rr.Code)
		})
	}
}
