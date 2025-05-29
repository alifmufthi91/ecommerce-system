package internal

import (
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse"
)

type Modules struct {
	Warehouse *warehouse.WarehouseModule
	Stock     *stock.StockModule
}

type InitOptions struct {
	_options.DefaultOptions
}

func InitModules(opts InitOptions) *Modules {

	warehouseModule := warehouse.NewWarehouseModule(warehouse.Options{
		DefaultOptions: opts.DefaultOptions,
	})

	stockModule := stock.NewStockModule(stock.Options{
		DefaultOptions: opts.DefaultOptions,
	})

	return &Modules{
		Warehouse: warehouseModule,
		Stock:     stockModule,
	}
}
