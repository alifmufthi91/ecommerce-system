package internal

import (
	productservice "github.com/alifmufthi91/ecommerce-system/services/order/external/product_service"
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/order/external/warehouse_service"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order"
)

type Modules struct {
	Order *order.OrderModule
}

type InitOptions struct {
	_options.DefaultOptions
}

func InitModules(opts InitOptions) *Modules {

	warehouseSvc := warehouseservice.Init(opts.DefaultOptions)
	productSvc := productservice.Init(opts.DefaultOptions)

	orderModule := order.NewOrderModule(order.Options{
		DefaultOptions:   opts.DefaultOptions,
		WarehouseService: warehouseSvc,
		ProductService:   productSvc,
	})

	return &Modules{
		Order: orderModule,
	}
}
