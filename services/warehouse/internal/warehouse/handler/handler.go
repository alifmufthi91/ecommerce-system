package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/middleware"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/warehouse/service"
	"github.com/gin-gonic/gin"
)

type warehouseHandler struct {
	router           *gin.Engine
	config           *config.Config
	logger           *pkg.Logger
	warehouseService service.WarehouseService
}

func NewHandler(rt *gin.Engine, cfg *config.Config, logger *pkg.Logger, warehouseSvc service.WarehouseService) registry.Router {
	return &warehouseHandler{
		warehouseService: warehouseSvc,
		router:           rt,
		config:           cfg,
		logger:           logger,
	}
}

func (h warehouseHandler) RegisterRoutes(base *gin.RouterGroup) {
	g := base.Group("/warehouses")

	g.Use(middleware.JwtMiddleware(h.config))

	g.GET("", h.GetWarehouses)
	g.PUT("/:id", h.UpdateWarehouse)
}
