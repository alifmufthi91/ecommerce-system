package pkg

import (
	"context"

	"github.com/alifmufthi91/ecommerce-system/services/product/config"
	"github.com/yuseferi/zax/v2"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func InitLogger(config *config.Config) *Logger {
	logger, err := zap.NewProduction(
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(
			zap.String("app_name", config.App.Name),
		),
		zap.WithCaller(true),
	)
	if err != nil {
		panic(err)
	}

	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	var fields []interface{}

	contextFields := zax.Get(ctx)

	for _, field := range contextFields {
		fields = append(fields, field.Key)
		fields = append(fields, field.String)
	}

	return &Logger{
		SugaredLogger: l.With(fields...),
	}
}
