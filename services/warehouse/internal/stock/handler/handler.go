package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/middleware"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/stock/service"
	"github.com/gin-gonic/gin"
)

type stockHandler struct {
	router       *gin.Engine
	config       *config.Config
	logger       *pkg.Logger
	stockService service.StockService
}

func NewHandler(rt *gin.Engine, cfg *config.Config, logger *pkg.Logger, stockSvc service.StockService) registry.Router {
	return &stockHandler{
		stockService: stockSvc,
		router:       rt,
		config:       cfg,
		logger:       logger,
	}
}

func (h stockHandler) RegisterRoutes(base *gin.RouterGroup) {
	g := base.Group("/stocks")

	g.Use(middleware.JwtMiddleware(h.config))

	g.GET("", h.GetStocks)
	g.POST("", h.CreateStock)
	g.POST("/transfer", h.TransferStock)
	g.GET("/availables", h.GetAvailableStocksByProduct)
	g.POST("/reserve", h.ReserveStocks)
	g.POST("/rollback", h.RollbackReserves)
	g.POST("/commit", h.CommitReserves)
}
