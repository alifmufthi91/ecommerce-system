package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/shop/config"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/middleware"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/shop/service"
	"github.com/gin-gonic/gin"
)

type shopHandler struct {
	router      *gin.Engine
	config      *config.Config
	logger      *pkg.Logger
	shopService service.ShopService
}

func NewHandler(rt *gin.Engine, cfg *config.Config, logger *pkg.Logger, shopSvc service.ShopService) registry.Router {
	return &shopHandler{
		shopService: shopSvc,
		router:      rt,
		config:      cfg,
		logger:      logger,
	}
}

func (h shopHandler) RegisterRoutes(base *gin.RouterGroup) {
	g := base.Group("/shops")

	g.Use(middleware.JwtMiddleware(h.config))

	g.GET("", h.GetShops)
}
