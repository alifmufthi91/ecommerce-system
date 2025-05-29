package cmd

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/shop/config"
	"github.com/gin-gonic/gin"

	docs "github.com/alifmufthi91/ecommerce-system/services/shop/docs"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/shop/internal/pkg/registry"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yuseferi/zax/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitRest(config *config.Config) *gin.Engine {
	switch config.App.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodPost, http.MethodPut},
		AllowCredentials: true,
	}))

	return router
}

func RegisterServiceRouter(config *config.Config, router *gin.Engine, db *sql.DB, logger *pkg.Logger) {
	serviceRouter := router.Group("/api")
	serviceRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	serviceRouter.Use(requestid.New())
	serviceRouter.Use(RequestLogger())
	serviceRouter.Use(otelgin.Middleware(config.App.Name))

	serviceRouter.Use(ginzap.GinzapWithConfig(logger.Desugar(), &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			fields = append(fields, zap.String("request_id", requestid.Get(c)))

			// log trace and span ID
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
			}

			return fields
		}),
	}))

	docs.SwaggerInfo.BasePath = "/api"
	serviceRouter.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	for _, router := range registry.GetRouters() {
		router.RegisterRoutes(serviceRouter)
	}
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := zax.Set(c.Request.Context(), []zap.Field{
			zap.String("request_id", requestid.Get(c)),
		})

		c.Request = c.Request.WithContext(ctx)

		// Process the request
		c.Next()
	}
}
