package handler

import (
	"github.com/alifmufthi91/ecommerce-system/services/user/config"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg/registry"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/user/service"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	router      *gin.Engine
	config      *config.Config
	logger      *pkg.Logger
	userService service.UserService
}

func NewHandler(rt *gin.Engine, cfg *config.Config, logger *pkg.Logger, userSvc service.UserService) registry.Router {
	return &userHandler{
		userService: userSvc,
		router:      rt,
		config:      cfg,
		logger:      logger,
	}
}

func (h userHandler) RegisterRoutes(base *gin.RouterGroup) {
	g := base.Group("/users")

	g.POST("", h.Register)
	g.POST("/login", h.Login)
}
