package internal

import (
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/product/external/warehouse_service"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product"
)

type Modules struct {
	Product *product.ProductModule
}

type InitOptions struct {
	_options.DefaultOptions
}

func InitModules(opts InitOptions) *Modules {

	warehouseSvc := warehouseservice.Init(opts.DefaultOptions)

	productModule := product.NewProductModule(product.Options{
		DefaultOptions:   opts.DefaultOptions,
		WarehouseService: warehouseSvc,
	})

	return &Modules{
		Product: productModule,
	}
}
