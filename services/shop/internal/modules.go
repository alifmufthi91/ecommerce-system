package internal

import (
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop"
)

type Modules struct {
	Shop *shop.ShopModule
}

type InitOptions struct {
	_options.DefaultOptions
}

func InitModules(opts InitOptions) *Modules {

	shopModule := shop.NewShopModule(shop.Options{
		DefaultOptions: opts.DefaultOptions,
	})

	return &Modules{
		Shop: shopModule,
	}
}
