package shop

import (
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/handler"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/repository"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/service"
)

type ShopModule struct {
	shopService service.ShopService
}

type Options struct {
	_options.DefaultOptions
}

func NewShopModule(opts Options) *ShopModule {

	shopRepo := repository.NewShopRepository(opts.Db)

	shopService := service.NewShopService(opts.Config, shopRepo)

	registry.RegisterRouter(handler.NewHandler(opts.Router, opts.Config, opts.Logger, shopService))

	return &ShopModule{
		shopService: shopService,
	}
}
