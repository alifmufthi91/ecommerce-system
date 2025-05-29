package _options

import (
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DefaultOptions struct {
	Config *config.Config
	Db     *gorm.DB
	Router *gin.Engine
	Logger *pkg.Logger
}
