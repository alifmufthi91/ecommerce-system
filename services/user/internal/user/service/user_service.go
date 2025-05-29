package service

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/user/config"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/auth"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/payload"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/repository"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockery --name=UserService --case underscore
type UserService interface {
	RegisterUser(ctx context.Context, req payload.RegisterUserReq) error
	LoginUser(ctx context.Context, req payload.LoginUserReq) (string, error)
}

type userService struct {
	config   *config.Config
	userRepo repository.UserRepository
}

func NewUserService(config *config.Config, userRepo repository.UserRepository) UserService {
	return &userService{
		config:   config,
		userRepo: userRepo,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req payload.RegisterUserReq) (err error) {
	ctx, span := observ.GetTracer().Start(ctx, "userService.RegisterUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	passwordHash, err := pkg.HashPassword(req.Password)
	if err != nil {
		return apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "Failed to hash password")
	}

	user := &model.User{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	err = s.userRepo.CreateUser(ctx, user)

	return err
}

func (s *userService) LoginUser(ctx context.Context, req payload.LoginUserReq) (token string, err error) {
	ctx, span := observ.GetTracer().Start(ctx, "userService.LoginUser")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	user, err := s.userRepo.GetUserByEmailOrPhone(ctx, req.EmailOrPhone)
	if err != nil {
		return "", err
	}

	if !pkg.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", apperr.NewWithCode(apperr.CodeHTTPUnauthorized, "Invalid credentials")
	}
	token, err = auth.GenerateToken(s.config.Token.JWTSecret, user.ID.String(), user.Email, user.Name, user.Phone)
	if err != nil {
		return "", apperr.WrapWithCode(err, apperr.CodeHTTPInternalServerError, "Failed to generate token")
	}

	return token, nil
}
