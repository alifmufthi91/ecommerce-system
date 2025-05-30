package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/warehouse/config"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/database"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/warehouse/internal/pkg"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

var (
	serverInstance *http.Server
	serverMutex    sync.Mutex
)

// Add this function for testing
func StartServerForTesting(config *config.Config) {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverInstance != nil {
		return // Server already running
	}

	defaultOpt := initDefaultOptionsForTesting(config)

	var (
		logger = defaultOpt.Logger
		db     = defaultOpt.Db
		router = defaultOpt.Router
	)

	RunMigration(db)

	// Initialize modules
	_ = internal.InitModules(internal.InitOptions{
		DefaultOptions: defaultOpt,
	})

	if sqlDb, err := db.DB(); err == nil {
		RegisterServiceRouter(config, router, sqlDb, logger)
	}

	logger.Info("Starting test server on port " + config.App.Port)

	serverInstance = &http.Server{
		Addr:    ":" + config.App.Port,
		Handler: router,
	}

	go func() {
		if err := serverInstance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start test server:", err)
		}
	}()
}

func initDefaultOptionsForTesting(config *config.Config) _options.DefaultOptions {
	logger := pkg.InitLogger(config)

	router := InitRest(config)

	db, _, err := database.InitDB(config, logger)
	if err != nil {
		logger.Fatal("Failed to connect to the database:", err)
	}

	defaultOpts := _options.DefaultOptions{
		Logger: logger,
		Router: router,
		Db:     db,
		Config: config,
	}

	return defaultOpts
}

func StopTestServer() error {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverInstance == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := serverInstance.Shutdown(ctx)
	serverInstance = nil
	return err
}

func RunMigration(gormDB *gorm.DB) {
	db, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to get db from gorm:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create migrate driver: %v", err)
	}

	absPath, err := filepath.Abs("../../../migrations") // adjust path as needed
	if err != nil {
		log.Fatalf("cannot resolve migration path: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("Migration succeeded")
}
