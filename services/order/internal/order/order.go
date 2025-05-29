package order

import (
	productservice "github.com/alifmufthi91/ecommerce-system/services/order/external/product_service"
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/order/external/warehouse_service"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/handler"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/repository"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/service"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/registry"
)

type OrderModule struct {
	OrderService service.OrderService
}

type Options struct {
	_options.DefaultOptions
	WarehouseService warehouseservice.IWarehouseSvc
	ProductService   productservice.IProductSvc
}

func NewOrderModule(opts Options) *OrderModule {

	orderRepo := repository.NewOrderRepository(opts.Db)
	stockLockRepo := repository.NewStockLockRepository(opts.Db)
	orderService := service.NewOrderService(opts.Config, opts.Db, opts.Logger, orderRepo, stockLockRepo, opts.WarehouseService, opts.ProductService)

	registry.RegisterRouter(handler.NewHandler(opts.Router, opts.Config, opts.Logger, orderService))

	return &OrderModule{
		OrderService: orderService,
	}
}
