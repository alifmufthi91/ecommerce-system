package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/alifmufthi91/ecommerce-system/services/product/config"
	"github.com/alifmufthi91/ecommerce-system/services/product/database"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/_options"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg"
	"github.com/alifmufthi91/ecommerce-system/services/product/internal/pkg/httpclient"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println(filepath.Base(os.Args[0]))
	rootCmd = &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: c.App.Name,
		Run: func(cmd *cobra.Command, args []string) {
			startServer(c)
		},
	}

	cobra.OnInitialize()
}

func startServer(config *config.Config) {
	defaultOpt := initDefaultOptions(config)

	var (
		logger = defaultOpt.Logger
		db     = defaultOpt.Db
		router = defaultOpt.Router
		dbConn = database.GetConnection()
	)

	// returned value can be used later when adding external connection like pubsub
	_ = internal.InitModules(internal.InitOptions{
		DefaultOptions: defaultOpt,
	})

	if sqlDb, err := db.DB(); err == nil {
		RegisterServiceRouter(config, router, sqlDb, logger)
	}

	logger.Info("Starting server on port " + config.App.Port)

	server := &http.Server{
		Addr:    ":" + config.App.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to shutdown server:", err)
	}

	_ = dbConn.Close()

	logger.Info("Server stopped")
}

func initDefaultOptions(config *config.Config) _options.DefaultOptions {

	logger := pkg.InitLogger(config)

	router := InitRest(config)

	db, _, err := database.InitDB(config, logger)
	if err != nil {
		logger.Fatal("Failed to connect to the database:", err)
	}

	httpClient := httpclient.Init(httpclient.Options{
		Config: config,
		Logger: logger,
	})

	defaultOpts := _options.DefaultOptions{
		Logger:     logger,
		Router:     router,
		Db:         db,
		Config:     config,
		HttpClient: httpClient,
	}

	return defaultOpts
}
