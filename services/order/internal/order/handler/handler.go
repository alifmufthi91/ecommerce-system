package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/order/config"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/order/service"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/middleware"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/registry"
	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	router       *gin.Engine
	config       *config.Config
	logger       *pkg.Logger
	orderService service.OrderService
}

func NewHandler(rt *gin.Engine, cfg *config.Config, logger *pkg.Logger, orderSvc service.OrderService) registry.Router {
	return &orderHandler{
		orderService: orderSvc,
		router:       rt,
		config:       cfg,
		logger:       logger,
	}
}

func (h orderHandler) RegisterRoutes(base *gin.RouterGroup) {
	g := base.Group("/orders")

	g.Use(middleware.JwtMiddleware(h.config))

	g.GET("", h.GetOrders)
	g.POST("", h.CreateOrder)
	g.PATCH("/:id/complete", h.CompleteOrder)
}
