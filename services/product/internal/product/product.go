package product

import (
	warehouseservice "github.com/alifmufthi91/ecommerce-system/services/product/external/warehouse_service"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/handler"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/repository"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/service"
)

type ProductModule struct {
	productService service.ProductService
}

type Options struct {
	_options.DefaultOptions
	WarehouseService warehouseservice.IWarehouseSvc
}

func NewProductModule(opts Options) *ProductModule {

	productRepo := repository.NewProductRepository(opts.Db)

	productService := service.NewProductService(opts.Config, opts.WarehouseService, productRepo)

	registry.RegisterRouter(handler.NewHandler(opts.Router, opts.Config, opts.Logger, productService))

	return &ProductModule{
		productService: productService,
	}
}
