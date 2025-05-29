package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/product/config"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/middleware"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/product/service"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	router         *gin.Engine
	config         *config.Config
	logger         *pkg.Logger
	productService service.ProductService
}

func NewHandler(rt *gin.Engine, cfg *config.Config, logger *pkg.Logger, productSvc service.ProductService) registry.Router {
	return &productHandler{
		productService: productSvc,
		router:         rt,
		config:         cfg,
		logger:         logger,
	}
}

func (h productHandler) RegisterRoutes(base *gin.RouterGroup) {
	g := base.Group("/products")

	g.Use(middleware.JwtMiddleware(h.config))

	g.POST("", h.CreateProduct)
	g.GET("", h.GetProducts)
	g.GET("/:id", h.GetProductByID)
}
