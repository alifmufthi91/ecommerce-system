package database

import (
	"database/sql"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/product/config"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	defaultConnection *sql.DB
)

func InitDB(config *config.Config, logger *pkg.Logger) (db *gorm.DB, sqlDB *sql.DB, err error) {

	zlogger := pkg.NewZapGormLogger(logger, gormLogger.Info, 1*time.Second)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: config.DB.DSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 zlogger,
	})

	if err != nil {
		return
	}

	sqlDB, err = db.DB()
	if err != nil {
		return
	}

	sqlDB.SetConnMaxIdleTime(time.Minute * 3)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	defaultConnection = sqlDB

	return
}

func GetConnection() *sql.DB {
	return defaultConnection
}
