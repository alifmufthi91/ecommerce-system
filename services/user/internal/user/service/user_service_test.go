package service

import (
	"context"
	"errors"
	"testing"

	"github.com/alifmufthi91/ecommerce-system/services/user/config"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/payload"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userRepoMock "github.com/alifmufthi91/ecommerce-system/services/user/internal/user/repository/mocks"
)

func TestRegisterUser_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		userRepo *userRepoMock.UserRepository
	}

	tests := []struct {
		name  string
		req   payload.RegisterUserReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success",
			req: payload.RegisterUserReq{
				Name:     "Test User",
				Email:    "test@example.com",
				Phone:    "+6281234567890",
				Password: "password123",
			},
			setup: func(m dependencyMocks) {
				m.userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				userRepo: userRepoMock.NewUserRepository(t),
			}
			userSvc := userService{
				userRepo: mocks.userRepo,
			}

			tt.setup(mocks)

			err := userSvc.RegisterUser(context.Background(), tt.req)

			assert.Nil(t, err)
			mocks.userRepo.AssertExpectations(t)
		})
	}
}

func TestRegisterUser_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		userRepo *userRepoMock.UserRepository
	}

	tests := []struct {
		name  string
		req   payload.RegisterUserReq
		setup func(
			m dependencyMocks,
		)
		wantErr bool
	}{
		{
			name: "error - failed to create user",
			req: payload.RegisterUserReq{
				Name:     "Test User",
				Email:    "test@example.com",
				Phone:    "+6281234567890",
				Password: "password123",
			},
			setup: func(m dependencyMocks) {
				m.userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
					Return(assert.AnError)
			},
		},
		{
			name: "error - failed to hash password",
			req: payload.RegisterUserReq{
				Name:     "Test User",
				Email:    "test@example.com",
				Phone:    "+6281234567890",
				Password: "longpasswordthatwillfailthehashingprocessneedtobeverylongandcomplexanddifficulttopredictandshouldnotbeusedinproduction",
			},
			setup: func(m dependencyMocks) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				userRepo: userRepoMock.NewUserRepository(t),
			}
			userSvc := userService{
				userRepo: mocks.userRepo,
			}

			tt.setup(mocks)

			err := userSvc.RegisterUser(context.Background(), tt.req)

			assert.Error(t, err)
			mocks.userRepo.AssertExpectations(t)
		})
	}
}

func TestLoginUser_ShouldSuccess(t *testing.T) {
	type dependencyMocks struct {
		userRepo *userRepoMock.UserRepository
	}

	tests := []struct {
		name  string
		req   payload.LoginUserReq
		setup func(
			m dependencyMocks,
		)
	}{
		{
			name: "success - login with email",
			req: payload.LoginUserReq{
				EmailOrPhone: "test@example.com",
				Password:     "password123",
			},
			setup: func(m dependencyMocks) {
				m.userRepo.On("GetUserByEmailOrPhone", mock.Anything, mock.AnythingOfType("string")).
					Return(model.User{
						ID:           uuid.New(),
						Name:         "Test User",
						Email:        "test@example.com",
						Phone:        "+6281234567890",
						PasswordHash: "$2a$12$SS.CMkrCtNG/o5JXwuqFt.wZ1Wby2UTrdAUK.bcwg7NYHhW7pQ/b2",
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				userRepo: userRepoMock.NewUserRepository(t),
			}

			mockConfig := &config.Config{
				Token: config.Token{
					JWTSecret: []byte("secret"),
				},
			}

			userSvc := userService{
				userRepo: mocks.userRepo,
				config:   mockConfig,
			}

			tt.setup(mocks)

			token, err := userSvc.LoginUser(context.Background(), tt.req)

			assert.Nil(t, err)
			assert.NotEmpty(t, token)
			mocks.userRepo.AssertExpectations(t)
		})
	}
}

func TestLoginUser_ShouldReturnError(t *testing.T) {
	type dependencyMocks struct {
		userRepo *userRepoMock.UserRepository
	}

	tests := []struct {
		name  string
		req   payload.LoginUserReq
		setup func(
			m dependencyMocks,
		)
		errorMsg string
	}{
		{
			name: "error - user not found",
			req: payload.LoginUserReq{
				EmailOrPhone: "test@example.com",
				Password:     "password123",
			},
			setup: func(m dependencyMocks) {
				m.userRepo.On("GetUserByEmailOrPhone", mock.Anything, mock.AnythingOfType("string")).
					Return(model.User{}, errors.New("user not found"))
			},
			errorMsg: "user not found",
		},
		{
			name: "error - invalid credentials",
			req: payload.LoginUserReq{
				EmailOrPhone: "test@example.com",
				Password:     "wrongpassword",
			},
			setup: func(m dependencyMocks) {
				m.userRepo.On("GetUserByEmailOrPhone", mock.Anything, mock.AnythingOfType("string")).
					Return(model.User{
						ID:           uuid.New(),
						Name:         "Test User",
						Email:        "test@example.com",
						Phone:        "+6281234567890",
						PasswordHash: "$2a$12$SS.CMkrCtNG/o5JXwuqFt.wZ1Wby2UTrdAUK.bcwg7NYHhW7pQ/b2",
					}, nil)
			},
			errorMsg: "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mocks := dependencyMocks{
				userRepo: userRepoMock.NewUserRepository(t),
			}

			mockConfig := &config.Config{
				Token: config.Token{
					JWTSecret: []byte("secret"),
				},
			}

			userSvc := userService{
				userRepo: mocks.userRepo,
				config:   mockConfig,
			}

			tt.setup(mocks)

			token, err := userSvc.LoginUser(context.Background(), tt.req)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
			assert.Empty(t, token)
			mocks.userRepo.AssertExpectations(t)
		})
	}
}
