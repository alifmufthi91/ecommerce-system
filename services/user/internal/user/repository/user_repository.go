package repository

import (
	"context"
	"strings"

	"github.com/alifmufthi91/ecommerce-system/services/user/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/apperr"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/observ"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=UserRepository --case underscore
type UserRepository interface {
	WithTX(tx *gorm.DB) UserRepository
	WithReturning() UserRepository
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmailOrPhone(ctx context.Context, emailOrPhone string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) WithTX(tx *gorm.DB) UserRepository {
	if tx == nil {
		return r
	}
	return &userRepository{db: tx}
}

func (r *userRepository) WithReturning() UserRepository {
	return &userRepository{
		db: r.db.Clauses(clause.Returning{}),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	ctx, span := observ.GetTracer().Start(ctx, "userRepository.CreateUser")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		if strings.Contains(err.Error(), "users_phone_key") {
			return apperr.NewWithCode(apperr.CodeHTTPBadRequest, "Phone number already exists")
		}
		if strings.Contains(err.Error(), "users_email_key") {
			return apperr.NewWithCode(apperr.CodeHTTPBadRequest, "Email already exists")
		}
		return apperr.WrapWithCode(err, apperr.CodeSQLCreate, "failed to create user")
	}
	return nil
}

func (r *userRepository) GetUserByEmailOrPhone(ctx context.Context, emailOrPhone string) (model.User, error) {
	ctx, span := observ.GetTracer().Start(ctx, "userRepository.GetUserByEmailOrPhone")
	defer span.End()

	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, apperr.NewWithCode(apperr.CodeHTTPNotFound, "User not found")
		}
		span.SetStatus(codes.Error, err.Error())
		return user, apperr.WrapWithCode(err, apperr.CodeSQLRead, "failed to get user by email or phone")
	}
	return user, nil
}
