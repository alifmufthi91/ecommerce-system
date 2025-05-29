package pkg

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type ZapGormLogger struct {
	ZapLogger     *Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

// NewZapGormLogger creates a new instance of ZapGormLogger
func NewZapGormLogger(zapLogger *Logger, logLevel logger.LogLevel, slowThreshold time.Duration) *ZapGormLogger {
	return &ZapGormLogger{
		ZapLogger:     zapLogger,
		LogLevel:      logLevel,
		SlowThreshold: slowThreshold,
	}
}

// LogMode implements the LogMode method of gorm.Logger
func (zl *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &ZapGormLogger{
		ZapLogger:     zl.ZapLogger,
		LogLevel:      level,
		SlowThreshold: zl.SlowThreshold,
	}
}

// Info logs general information
func (zl *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if zl.LogLevel >= logger.Info {
		zl.ZapLogger.WithContext(ctx).Infof(msg, data...)
	}
}

// Warn logs warnings
func (zl *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if zl.LogLevel >= logger.Warn {
		zl.ZapLogger.WithContext(ctx).Warnf(msg, data...)
	}
}

// Error logs errors
func (zl *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if zl.LogLevel >= logger.Error {
		zl.ZapLogger.WithContext(ctx).Errorf(msg, data...)
	}
}

// Trace logs SQL queries and execution time
func (zl *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		if zl.LogLevel >= logger.Error {
			zl.ZapLogger.WithContext(ctx).WithOptions(zap.AddCallerSkip(3)).With("sql", sql, "rows_affected", rows, "elapsed", elapsed).Error(err.Error())
		}
		return
	}

	switch {
	case elapsed > zl.SlowThreshold && zl.SlowThreshold != 0:
		if zl.LogLevel >= logger.Warn {
			zl.ZapLogger.WithContext(ctx).WithOptions(zap.AddCallerSkip(3)).With("sql", sql, "rows_affected", rows, "elapsed", elapsed).Warn("Slow SQL query")
		}
	default:
		if zl.LogLevel >= logger.Info {
			zl.ZapLogger.WithContext(ctx).WithOptions(zap.AddCallerSkip(3)).With("sql", sql, "rows_affected", rows, "elapsed", elapsed).Info("SQL executed")
		}
	}
}
