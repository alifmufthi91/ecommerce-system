package warehouse

import (
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/handler"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/repository"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/service"
)

type WarehouseModule struct {
	warehouseService service.WarehouseService
}

type Options struct {
	_options.DefaultOptions
}

func NewWarehouseModule(opts Options) *WarehouseModule {

	warehouseRepo := repository.NewWarehouseRepository(opts.Db)

	warehouseService := service.NewWarehouseService(opts.Config, warehouseRepo)

	registry.RegisterRouter(handler.NewHandler(opts.Router, opts.Config, opts.Logger, warehouseService))

	return &WarehouseModule{
		warehouseService: warehouseService,
	}
}
