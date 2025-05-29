package service

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/shop/config"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/observ"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/repository"
	"go.opentelemetry.io/otel/codes"
)

//go:generate mockery --name=ShopService --case underscore
type ShopService interface {
	GetShops(ctx context.Context) ([]model.Shop, error)
}

type shopService struct {
	config   *config.Config
	shopRepo repository.ShopRepository
}

func NewShopService(config *config.Config, shopRepo repository.ShopRepository) ShopService {
	return &shopService{
		config:   config,
		shopRepo: shopRepo,
	}
}

func (s *shopService) GetShops(ctx context.Context) ([]model.Shop, error) {
	ctx, span := observ.GetTracer().Start(ctx, "shopService.GetShops")
	defer span.End()

	shops, err := s.shopRepo.GetShops(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return shops, nil
}
