package stock

import (
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/handler"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/repository"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/service"
)

type StockModule struct {
	stockService service.StockService
}

type Options struct {
	_options.DefaultOptions
}

func NewStockModule(opts Options) *StockModule {

	stockRepo := repository.NewStockRepository(opts.Db)

	stockService := service.NewStockService(opts.Config, opts.Db, stockRepo)

	registry.RegisterRouter(handler.NewHandler(opts.Router, opts.Config, opts.Logger, stockService))

	return &StockModule{
		stockService: stockService,
	}
}
