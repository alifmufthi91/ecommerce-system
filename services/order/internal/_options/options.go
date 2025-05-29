package _options

import (
	"github.com/alifmufthi91/ecommerce-system/services/order/config"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/order/internal/pkg/httpclient"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DefaultOptions struct {
	Config     *config.Config
	Db         *gorm.DB
	Router     *gin.Engine
	Logger     *pkg.Logger
	HttpClient httpclient.IHTTPClient
}
